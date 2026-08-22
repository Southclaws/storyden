package robot_session

import (
	"context"
	"maps"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	adksession "google.golang.org/adk/v2/session"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/internal/ent"
	ent_robot_session "github.com/Southclaws/storyden/internal/ent/robotsession"
	ent_robot_session_input "github.com/Southclaws/storyden/internal/ent/robotsessioninput"
	ent_robot_session_message "github.com/Southclaws/storyden/internal/ent/robotsessionmessage"
	ent_robot_session_view "github.com/Southclaws/storyden/internal/ent/robotsessionview"
	entschema "github.com/Southclaws/storyden/internal/ent/schema"
)

type CreateOption func(*ent.RobotSessionMutation)

func WithName(name string) CreateOption {
	return func(m *ent.RobotSessionMutation) {
		m.SetName(name)
	}
}

func WithTrailOrigin(actionRunID xid.ID) CreateOption {
	return func(m *ent.RobotSessionMutation) {
		m.SetOriginKind(ent_robot_session.OriginKindTrailAction)
		m.SetOriginID(actionRunID)
	}
}

func (q *Repository) Create(
	ctx context.Context,
	sessionID robot.SessionID,
	name string,
	accountID account.AccountID,
	state map[string]any,
	options ...CreateOption,
) (*robot.Session, error) {
	err := ent.WithTx(ctx, q.db, func(tx *ent.Tx) error {
		create := tx.RobotSession.Create().
			SetID(xid.ID(sessionID)).
			SetName(name).
			SetAccountID(xid.ID(accountID)).
			SetState(state)
		mutation := create.Mutation()
		for _, option := range options {
			option(mutation)
		}
		_, err := create.Save(ctx)
		if err != nil {
			return err
		}
		_, err = tx.RobotSessionView.Create().
			SetSessionID(xid.ID(sessionID)).
			SetAccountID(xid.ID(accountID)).
			Save(ctx)
		return err
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sess, err := withCreator(q.db.RobotSession.Query().
		Where(ent_robot_session.IDEQ(xid.ID(sessionID)))).Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return robot.MapSession(sess, nil)
}

func (q *Repository) EnsureView(ctx context.Context, sessionID robot.SessionID, accountID account.AccountID) error {
	now := time.Now()
	err := q.db.RobotSessionView.Create().
		SetSessionID(xid.ID(sessionID)).
		SetAccountID(xid.ID(accountID)).
		SetLastAccessedAt(now).
		OnConflictColumns(ent_robot_session_view.FieldSessionID, ent_robot_session_view.FieldAccountID).
		SetLastAccessedAt(now).
		Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	return nil
}

func (q *Repository) Delete(ctx context.Context, sessionID robot.SessionID) error {
	err := q.db.RobotSession.DeleteOneID(xid.ID(sessionID)).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil
		}
		return fault.Wrap(err, fctx.With(ctx))
	}
	return nil
}

func (q *Repository) UpdateName(
	ctx context.Context,
	sessionID robot.SessionID,
	name string,
) error {
	err := q.db.RobotSession.UpdateOneID(xid.ID(sessionID)).
		SetName(name).
		Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (q *Repository) UpdateState(
	ctx context.Context,
	sessionID robot.SessionID,
	state map[string]any,
) error {
	update := q.db.RobotSession.Update().
		Where(ent_robot_session.IDEQ(xid.ID(sessionID))).
		SetState(state)
	if lease, ok := ExecutionLeaseFromContext(ctx); ok {
		if lease.SessionID != sessionID {
			return fault.Wrap(ErrLeaseLost, fctx.With(ctx))
		}
		update.Where(activeExecutionLeasePredicates(lease)...)
	}
	updated, err := update.Save(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	if updated != 1 {
		return fault.Wrap(ErrLeaseLost, fctx.With(ctx))
	}
	return nil
}

func (q *Repository) AppendMessage(
	ctx context.Context,
	sessionID robot.SessionID,
	accountID opt.Optional[account.AccountID],
	actor opt.Optional[robot.Actor],
	event *adksession.Event,
) error {
	if event == nil {
		return fault.New("robot session message event is required", fctx.With(ctx))
	}

	storedEvent, stateDelta := eventForPersistence(event)
	lease, hasLease := ExecutionLeaseFromContext(ctx)
	if hasLease && lease.SessionID != sessionID {
		return fault.Wrap(ErrLeaseLost, fctx.With(ctx))
	}
	if len(stateDelta) == 0 {
		err := ent.WithTx(ctx, q.db, func(tx *ent.Tx) error {
			var leasePtr *ExecutionLease
			if hasLease {
				leasePtr = &lease
			}
			sequence, err := allocateSessionEventSequence(ctx, tx, sessionID, leasePtr)
			if err != nil {
				return err
			}
			hidden, err := shouldHideRuntimeInput(ctx, tx, leasePtr, storedEvent)
			if err != nil {
				return err
			}
			return saveMessage(ctx, tx.RobotSessionMessage.Create(), sessionID, accountID, actor, storedEvent, sequence, leasePtr, ent_robot_session_message.EventKindMessage, hidden)
		})
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
		return nil
	}

	err := ent.WithTx(ctx, q.db, func(tx *ent.Tx) error {
		var leasePtr *ExecutionLease
		if hasLease {
			leasePtr = &lease
		}
		sequence, err := allocateSessionEventSequence(ctx, tx, sessionID, leasePtr)
		if err != nil {
			return err
		}
		session, err := tx.RobotSession.Query().
			Where(ent_robot_session.IDEQ(xid.ID(sessionID))).
			Only(ctx)
		if err != nil {
			return err
		}

		state := maps.Clone(session.State)
		if state == nil {
			state = make(map[string]any, len(stateDelta))
		}
		maps.Copy(state, stateDelta)
		if err := tx.RobotSession.UpdateOneID(xid.ID(sessionID)).SetState(state).Exec(ctx); err != nil {
			return err
		}

		hidden, err := shouldHideRuntimeInput(ctx, tx, leasePtr, storedEvent)
		if err != nil {
			return err
		}
		return saveMessage(ctx, tx.RobotSessionMessage.Create(), sessionID, accountID, actor, storedEvent, sequence, leasePtr, ent_robot_session_message.EventKindMessage, hidden)
	})
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	return nil
}

func saveMessage(
	ctx context.Context,
	create *ent.RobotSessionMessageCreate,
	sessionID robot.SessionID,
	accountID opt.Optional[account.AccountID],
	actor opt.Optional[robot.Actor],
	event *adksession.Event,
	sequence uint64,
	lease *ExecutionLease,
	eventKind ent_robot_session_message.EventKind,
	hiddenFromProjection bool,
) error {
	create = create.
		SetSessionID(xid.ID(sessionID)).
		SetSequence(sequence).
		SetEventKind(eventKind).
		SetHiddenFromProjection(hiddenFromProjection).
		SetInvocationID(event.InvocationID).
		SetEventData(entschema.NewRobotSessionEvent(*event))
	if lease != nil {
		create.SetTurnID(lease.TurnID)
	}
	if event.Branch != "" {
		create.SetBranch(event.Branch)
	}
	if event.IsolationScope != "" {
		create.SetIsolationScope(event.IsolationScope)
	}

	if aid, ok := accountID.Get(); ok {
		create.SetAccountID(xid.ID(aid))
	}

	if err := applyMessageActor(create, accountID, actor); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	_, err := create.Save(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

// ADK persists the effective input at the start of every run. The visible
// inputs were already appended individually when they were queued, so exposing
// this event would duplicate them as one merged message. Tool responses also
// use the user role, but remain part of the visible conversation history.
func shouldHideRuntimeInput(ctx context.Context, tx *ent.Tx, lease *ExecutionLease, event *adksession.Event) (bool, error) {
	if lease == nil || event.Author != "user" || event.Content == nil {
		return false, nil
	}
	for _, part := range event.Content.Parts {
		if part != nil && part.FunctionResponse != nil {
			return false, nil
		}
	}
	return tx.RobotSessionInput.Query().Where(ent_robot_session_input.TurnIDEQ(lease.TurnID)).Exist(ctx)
}

// eventForPersistence applies ADK's session-state contract without mutating
// the event yielded to the active invocation. Temporary keys are invocation
// scoped and must not be written to either the session state or event history.
func eventForPersistence(event *adksession.Event) (*adksession.Event, map[string]any) {
	if len(event.Actions.StateDelta) == 0 {
		return event, nil
	}

	stateDelta := make(map[string]any, len(event.Actions.StateDelta))
	for key, value := range event.Actions.StateDelta {
		if strings.HasPrefix(key, adksession.KeyPrefixTemp) {
			continue
		}
		stateDelta[key] = value
	}
	if len(stateDelta) == len(event.Actions.StateDelta) {
		return event, stateDelta
	}

	stored := *event
	stored.Actions.StateDelta = stateDelta
	return &stored, stateDelta
}

func applyMessageActor(create *ent.RobotSessionMessageCreate, accountID opt.Optional[account.AccountID], actorOpt opt.Optional[robot.Actor]) error {
	actor, ok := actorOpt.Get()
	if !ok {
		return nil
	}
	if _, hasAccount := accountID.Get(); hasAccount {
		return fault.New("robot session message cannot have both human author and robot actor")
	}
	if err := actor.Validate(); err != nil {
		return err
	}

	if rid, ok := actor.DatabaseRobotID.Get(); ok {
		create.SetRobotID(rid)
	}
	if bid, ok := actor.BuiltinRobotID.Get(); ok {
		create.SetBuiltinRobot(bid.String())
	}

	return nil
}
