package robot_session

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/internal/ent"
	ent_robot_session "github.com/Southclaws/storyden/internal/ent/robotsession"
	ent_robot_session_message "github.com/Southclaws/storyden/internal/ent/robotsessionmessage"
	ent_robot_session_turn "github.com/Southclaws/storyden/internal/ent/robotsessionturn"
	ent_robot_session_view "github.com/Southclaws/storyden/internal/ent/robotsessionview"
)

var ErrTurnNotFound = errors.New("robot session turn not found")

type EnqueueTurnParams struct {
	ID             robot.TurnID
	SessionID      robot.SessionID
	InitiatorID    opt.Optional[account.AccountID]
	SourceKind     string
	RobotRef       string
	InputData      json.RawMessage
	ContinuationOf opt.Optional[robot.TurnID]
}

func (q *Repository) EnqueueTurn(ctx context.Context, params EnqueueTurnParams) error {
	err := ent.WithTx(ctx, q.db, func(tx *ent.Tx) error {
		create := tx.RobotSessionTurn.Create().
			SetID(xid.ID(params.ID)).
			SetSessionID(xid.ID(params.SessionID)).
			SetSourceKind(params.SourceKind).
			SetRobotRef(params.RobotRef).
			SetInputData(params.InputData)
		if initiator, ok := params.InitiatorID.Get(); ok {
			create.SetInitiatedByAccountID(xid.ID(initiator))
		}
		if parent, ok := params.ContinuationOf.Get(); ok {
			create.SetContinuationOfTurnID(xid.ID(parent))
		}
		if _, err := create.Save(ctx); err != nil {
			return err
		}

		sequence, err := allocateSessionEventSequence(ctx, tx, params.SessionID, nil)
		if err != nil {
			return err
		}
		_, err = tx.RobotSessionMessage.Create().
			SetSessionID(xid.ID(params.SessionID)).
			SetTurnID(xid.ID(params.ID)).
			SetSequence(sequence).
			SetEventKind(ent_robot_session_message.EventKindTurnQueued).
			Save(ctx)
		return err
	})
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	return nil
}

func (q *Repository) FinishTurn(ctx context.Context, sessionID robot.SessionID, turnID robot.TurnID, status robot.TurnStatus, errorText string) error {
	turnStatus, eventKind, err := terminalKinds(status)
	if err != nil {
		return err
	}
	now := time.Now()
	err = ent.WithTx(ctx, q.db, func(tx *ent.Tx) error {
		eligible := []ent_robot_session_turn.Status{ent_robot_session_turn.StatusRunning}
		if status == robot.TurnStatusFailed || status == robot.TurnStatusCancelled {
			eligible = append(eligible, ent_robot_session_turn.StatusQueued)
		}
		update := tx.RobotSessionTurn.Update().
			Where(
				ent_robot_session_turn.IDEQ(xid.ID(turnID)),
				ent_robot_session_turn.StatusIn(eligible...),
			).
			SetStatus(turnStatus).
			SetFinishedAt(now)
		if errorText != "" {
			update.SetErrorText(errorText)
		}
		updated, err := update.Save(ctx)
		if err != nil {
			return err
		}
		if updated != 1 {
			stored, err := tx.RobotSessionTurn.Get(ctx, xid.ID(turnID))
			if err == nil && stored.Status == turnStatus {
				return nil
			}
			return ErrTurnNotFound
		}

		sequence, err := allocateSessionEventSequence(ctx, tx, sessionID, nil)
		if err != nil {
			return err
		}
		create := tx.RobotSessionMessage.Create().
			SetSessionID(xid.ID(sessionID)).
			SetTurnID(xid.ID(turnID)).
			SetSequence(sequence).
			SetEventKind(eventKind)
		if errorText != "" {
			create.SetErrorText(errorText)
		}
		_, err = create.Save(ctx)
		return err
	})
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	return nil
}

func terminalKinds(status robot.TurnStatus) (ent_robot_session_turn.Status, ent_robot_session_message.EventKind, error) {
	switch status {
	case robot.TurnStatusCompleted:
		return ent_robot_session_turn.StatusCompleted, ent_robot_session_message.EventKindTurnCompleted, nil
	case robot.TurnStatusBlocked:
		return ent_robot_session_turn.StatusBlocked, ent_robot_session_message.EventKindTurnBlocked, nil
	case robot.TurnStatusFailed:
		return ent_robot_session_turn.StatusFailed, ent_robot_session_message.EventKindTurnFailed, nil
	case robot.TurnStatusCancelled:
		return ent_robot_session_turn.StatusCancelled, ent_robot_session_message.EventKindTurnCancelled, nil
	default:
		return "", "", errors.New("turn status is not terminal")
	}
}

func (q *Repository) ReadTurnEvents(ctx context.Context, sessionID robot.SessionID, turnID robot.TurnID, after uint64, limit int) ([]robot.SessionEvent, uint64, bool, error) {
	if limit <= 0 {
		limit = robot.SessionEventPageSize
	}
	rows, err := q.db.RobotSessionMessage.Query().
		Where(
			ent_robot_session_message.SessionIDEQ(xid.ID(sessionID)),
			ent_robot_session_message.TurnIDEQ(xid.ID(turnID)),
			ent_robot_session_message.SequenceGT(after),
		).
		WithRobot(func(query *ent.RobotQuery) { query.WithAuthor() }).
		WithAuthor().
		Order(ent.Asc(ent_robot_session_message.FieldSequence)).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, after, false, fault.Wrap(err, fctx.With(ctx))
	}

	events := make([]robot.SessionEvent, 0, len(rows))
	next := after
	closed := false
	for _, row := range rows {
		next = row.Sequence
		event := robot.SessionEvent{
			Sequence:  row.Sequence,
			Kind:      robot.SessionEventKind(row.EventKind),
			TurnID:    turnID,
			ErrorText: opt.NewPtr(row.ErrorText),
		}
		if row.EventKind == ent_robot_session_message.EventKindMessage {
			message, err := robot.MapMessage(row)
			if err != nil {
				return nil, after, false, fault.Wrap(err, fctx.With(ctx))
			}
			event.Message = opt.New(message)
		}
		if row.EventKind != ent_robot_session_message.EventKindMessage && row.EventKind != ent_robot_session_message.EventKindTurnQueued {
			closed = true
		}
		events = append(events, event)
	}
	return events, next, closed, nil
}

func (q *Repository) CanReadTurn(ctx context.Context, sessionID robot.SessionID, turnID robot.TurnID, accountID account.AccountID) (bool, error) {
	exists, err := q.db.RobotSessionTurn.Query().
		Where(
			ent_robot_session_turn.IDEQ(xid.ID(turnID)),
			ent_robot_session_turn.SessionIDEQ(xid.ID(sessionID)),
		).
		Exist(ctx)
	if err != nil || !exists {
		return false, err
	}
	return q.db.RobotSessionView.Query().
		Where(
			ent_robot_session_view.SessionIDEQ(xid.ID(sessionID)),
			ent_robot_session_view.AccountIDEQ(xid.ID(accountID)),
		).
		Exist(ctx)
}

func (q *Repository) ActiveTurn(ctx context.Context, sessionID robot.SessionID) (robot.TurnID, error) {
	turn, err := q.db.RobotSessionTurn.Query().
		Where(
			ent_robot_session_turn.SessionIDEQ(xid.ID(sessionID)),
			ent_robot_session_turn.StatusIn(ent_robot_session_turn.StatusQueued, ent_robot_session_turn.StatusRunning),
		).
		Order(ent.Desc(ent_robot_session_turn.FieldCreatedAt)).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return robot.TurnID{}, ErrTurnNotFound
		}
		return robot.TurnID{}, fault.Wrap(err, fctx.With(ctx))
	}
	return robot.TurnID(turn.ID), nil
}

func (q *Repository) BlockedTurn(ctx context.Context, sessionID robot.SessionID) (robot.TurnID, error) {
	session, err := q.db.RobotSession.Query().
		Where(
			ent_robot_session.IDEQ(xid.ID(sessionID)),
			ent_robot_session.ExecutionStatusEQ(ent_robot_session.ExecutionStatusBlocked),
			ent_robot_session.ActiveTurnIDNotNil(),
		).
		Only(ctx)
	if err != nil || session.ActiveTurnID == nil {
		return robot.TurnID{}, ErrTurnNotFound
	}
	return robot.TurnID(*session.ActiveTurnID), nil
}

func (q *Repository) SetResumeTurn(ctx context.Context, sessionID robot.SessionID, accountID account.AccountID, turnID *robot.TurnID) error {
	update := q.db.RobotSessionView.Update().Where(
		ent_robot_session_view.SessionIDEQ(xid.ID(sessionID)),
		ent_robot_session_view.AccountIDEQ(xid.ID(accountID)),
	)
	if turnID == nil {
		update.ClearResumeTurnID()
	} else {
		update.SetResumeTurnID(xid.ID(*turnID))
	}
	updated, err := update.Save(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	if updated != 1 {
		return ErrTurnNotFound
	}
	return nil
}

func (q *Repository) ResumeTurn(ctx context.Context, sessionID robot.SessionID, accountID account.AccountID) (robot.TurnID, error) {
	view, err := q.db.RobotSessionView.Query().
		Where(
			ent_robot_session_view.SessionIDEQ(xid.ID(sessionID)),
			ent_robot_session_view.AccountIDEQ(xid.ID(accountID)),
		).
		Only(ctx)
	if err != nil {
		return robot.TurnID{}, ErrTurnNotFound
	}
	if view.ResumeTurnID == nil {
		return robot.TurnID{}, ErrTurnNotFound
	}
	return robot.TurnID(*view.ResumeTurnID), nil
}

func allocateSessionEventSequence(ctx context.Context, tx *ent.Tx, sessionID robot.SessionID, lease *ExecutionLease) (uint64, error) {
	update := tx.RobotSession.Update().Where(ent_robot_session.IDEQ(xid.ID(sessionID)))
	if lease != nil {
		update.Where(activeExecutionLeasePredicates(*lease)...)
	}
	updated, err := update.AddNextEventSequence(1).Save(ctx)
	if err != nil {
		return 0, err
	}
	if updated != 1 {
		if lease != nil {
			return 0, ErrLeaseLost
		}
		return 0, ErrTurnNotFound
	}
	session, err := tx.RobotSession.Get(ctx, xid.ID(sessionID))
	if err != nil {
		return 0, err
	}
	return session.NextEventSequence, nil
}
