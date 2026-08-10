package robot_session

import (
	"context"
	"maps"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	adksession "google.golang.org/adk/v2/session"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/internal/ent"
	ent_robot_session "github.com/Southclaws/storyden/internal/ent/robotsession"
	entschema "github.com/Southclaws/storyden/internal/ent/schema"
)

func (q *Repository) Create(
	ctx context.Context,
	sessionID robot.SessionID,
	name string,
	accountID account.AccountID,
	state map[string]any,
) (*robot.Session, error) {
	_, err := q.db.RobotSession.Create().
		SetID(xid.ID(sessionID)).
		SetName(name).
		SetAccountID(xid.ID(accountID)).
		SetState(state).
		Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sess, err := q.db.RobotSession.Query().
		Where(ent_robot_session.IDEQ(xid.ID(sessionID))).
		WithUser().
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return robot.MapSession(sess, nil)
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
	err := q.db.RobotSession.UpdateOneID(xid.ID(sessionID)).
		SetState(state).
		Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
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
	if len(stateDelta) == 0 {
		return saveMessage(ctx, q.db.RobotSessionMessage.Create(), sessionID, accountID, actor, storedEvent)
	}

	err := ent.WithTx(ctx, q.db, func(tx *ent.Tx) error {
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

		return saveMessage(ctx, tx.RobotSessionMessage.Create(), sessionID, accountID, actor, storedEvent)
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
) error {
	create = create.
		SetSessionID(xid.ID(sessionID)).
		SetInvocationID(event.InvocationID).
		SetEventData(entschema.NewRobotSessionEvent(*event))
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
