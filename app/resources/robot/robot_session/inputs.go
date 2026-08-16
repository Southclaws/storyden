package robot_session

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"entgo.io/ent/dialect/sql"
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
	ent_robot_session_turn "github.com/Southclaws/storyden/internal/ent/robotsessionturn"
)

var ErrNoQueuedInputs = errors.New("robot session has no queued inputs")

type EnqueueInputParams struct {
	ID           robot.InputID
	SessionID    robot.SessionID
	AccountID    account.AccountID
	SourceKind   string
	BatchKey     string
	InputData    json.RawMessage
	NotBefore    opt.Optional[time.Time]
	VisibleEvent opt.Optional[*adksession.Event]
}

// EnqueueInput durably records one submission and, for human messages, appends
// its individually visible representation to the ordered session log.
func (q *Repository) EnqueueInput(ctx context.Context, params EnqueueInputParams) error {
	notBefore := opt.Map(params.NotBefore, func(value time.Time) time.Time {
		return value.UTC()
	})
	err := ent.WithTx(ctx, q.db, func(tx *ent.Tx) error {
		sequence, err := allocateSessionEventSequence(ctx, tx, params.SessionID, nil)
		if err != nil {
			return err
		}
		if _, err := tx.RobotSessionInput.Create().
			SetID(xid.ID(params.ID)).
			SetSessionID(xid.ID(params.SessionID)).
			SetAccountID(xid.ID(params.AccountID)).
			SetSequence(sequence).
			SetSourceKind(params.SourceKind).
			SetBatchKey(params.BatchKey).
			SetInputData(params.InputData).
			SetNillableNotBefore(notBefore.Ptr()).
			Save(ctx); err != nil {
			return err
		}

		visible, ok := params.VisibleEvent.Get()
		if !ok {
			return nil
		}
		return saveMessage(
			ctx,
			tx.RobotSessionMessage.Create().SetID(xid.ID(params.ID)),
			params.SessionID,
			opt.New(params.AccountID),
			opt.NewEmpty[robot.Actor](),
			visible,
			sequence,
			nil,
			ent_robot_session_message.EventKindInputQueued,
			false,
		)
	})
	if err != nil {
		if ent.IsConstraintError(err) {
			exists, lookupErr := q.db.RobotSessionInput.Query().
				Where(ent_robot_session_input.IDEQ(xid.ID(params.ID))).
				Exist(ctx)
			if lookupErr != nil {
				return fault.Wrap(lookupErr, fctx.With(ctx))
			}
			if exists {
				return nil
			}
		}
		return fault.Wrap(err, fctx.With(ctx))
	}
	return nil
}

func (q *Repository) QueuedInputs(ctx context.Context, sessionID robot.SessionID, limit int) ([]robot.SessionInput, error) {
	if limit <= 0 {
		limit = robot.SessionEventPageSize
	}
	rows, err := q.db.RobotSessionInput.Query().
		Where(
			ent_robot_session_input.SessionIDEQ(xid.ID(sessionID)),
			ent_robot_session_input.StatusEQ(ent_robot_session_input.StatusQueued),
			ent_robot_session_input.Or(
				ent_robot_session_input.NotBeforeIsNil(),
				ent_robot_session_input.NotBeforeLTE(time.Now().UTC()),
			),
		).
		Order(ent_robot_session_input.BySequence(sql.OrderAsc())).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	inputs := make([]robot.SessionInput, len(rows))
	for i, row := range rows {
		inputs[i] = robot.SessionInput{
			ID: robot.InputID(row.ID), SessionID: robot.SessionID(row.SessionID), AccountID: row.AccountID,
			Sequence: row.Sequence, SourceKind: row.SourceKind, BatchKey: row.BatchKey,
			InputData: row.InputData, NotBefore: opt.NewPtr(row.NotBefore), CreatedAt: row.CreatedAt,
		}
	}
	return inputs, nil
}

type MaterialiseTurnParams struct {
	ID             robot.TurnID
	SessionID      robot.SessionID
	InputIDs       []robot.InputID
	InitiatorID    account.AccountID
	SourceKind     string
	RobotRef       string
	InputData      json.RawMessage
	ContinuationOf opt.Optional[robot.TurnID]
}

// MaterialiseTurn atomically assigns queued inputs to one runtime turn. If any
// input was already claimed, the transaction is rolled back and the caller can
// retry from the current queue head.
func (q *Repository) MaterialiseTurn(ctx context.Context, params MaterialiseTurnParams) error {
	if len(params.InputIDs) == 0 {
		return ErrNoQueuedInputs
	}
	ids := make([]xid.ID, len(params.InputIDs))
	for i, id := range params.InputIDs {
		ids[i] = xid.ID(id)
	}

	err := ent.WithTx(ctx, q.db, func(tx *ent.Tx) error {
		create := tx.RobotSessionTurn.Create().
			SetID(xid.ID(params.ID)).
			SetSessionID(xid.ID(params.SessionID)).
			SetInitiatedByAccountID(xid.ID(params.InitiatorID)).
			SetSourceKind(params.SourceKind).
			SetRobotRef(params.RobotRef).
			SetInputData(params.InputData)
		if parent, ok := params.ContinuationOf.Get(); ok {
			create.SetContinuationOfTurnID(xid.ID(parent))
		}
		if _, err := create.Save(ctx); err != nil {
			return err
		}

		claimed, err := tx.RobotSessionInput.Update().
			Where(
				ent_robot_session_input.IDIn(ids...),
				ent_robot_session_input.SessionIDEQ(xid.ID(params.SessionID)),
				ent_robot_session_input.StatusEQ(ent_robot_session_input.StatusQueued),
			).
			SetStatus(ent_robot_session_input.StatusClaimed).
			SetTurnID(xid.ID(params.ID)).
			Save(ctx)
		if err != nil {
			return err
		}
		if claimed != len(ids) {
			return ErrNoQueuedInputs
		}

		if _, err := tx.RobotSessionMessage.Update().
			Where(ent_robot_session_message.IDIn(ids...), ent_robot_session_message.EventKindEQ(ent_robot_session_message.EventKindInputQueued)).
			SetTurnID(xid.ID(params.ID)).
			Save(ctx); err != nil {
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
			SetInputIds(ids).
			Save(ctx)
		return err
	})
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	return nil
}

func (q *Repository) NextQueuedTurn(ctx context.Context, sessionID robot.SessionID) (*robot.Turn, error) {
	row, err := q.db.RobotSessionTurn.Query().
		Where(ent_robot_session_turn.SessionIDEQ(xid.ID(sessionID)), ent_robot_session_turn.StatusEQ(ent_robot_session_turn.StatusQueued)).
		Order(ent_robot_session_turn.ByCreatedAt(sql.OrderAsc())).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrTurnNotFound
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return mapTurn(row), nil
}

// RecoverableTurn returns the active running turn after its execution lease has
// expired. Any replica may reclaim it; AcquireQueuedTurnExecution provides the
// generation fence that prevents a stale worker from writing afterward.
func (q *Repository) RecoverableTurn(ctx context.Context, sessionID robot.SessionID) (*robot.Turn, error) {
	session, err := q.db.RobotSession.Query().Where(
		ent_robot_session.IDEQ(xid.ID(sessionID)),
		ent_robot_session.ExecutionStatusEQ(ent_robot_session.ExecutionStatusRunning),
		ent_robot_session.ActiveTurnIDNotNil(),
		ent_robot_session.Or(
			ent_robot_session.LeaseExpiresAtIsNil(),
			ent_robot_session.LeaseExpiresAtLT(time.Now()),
		),
	).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrTurnNotFound
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	row, err := q.db.RobotSessionTurn.Get(ctx, *session.ActiveTurnID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrTurnNotFound
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return mapTurn(row), nil
}

// RunnableSessionIDs returns sessions that need a coordinator signal. This is
// intentionally replica-agnostic: duplicate signals are harmless because turn
// materialisation and execution leases are transactional.
func (q *Repository) RunnableSessionIDs(ctx context.Context, limit int) ([]robot.SessionID, error) {
	if limit <= 0 {
		limit = robot.SessionEventPageSize
	}
	var inputSessionIDs []xid.ID
	err := q.db.RobotSessionInput.Query().
		Where(
			ent_robot_session_input.StatusEQ(ent_robot_session_input.StatusQueued),
			ent_robot_session_input.Or(
				ent_robot_session_input.NotBeforeIsNil(),
				ent_robot_session_input.NotBeforeLTE(time.Now().UTC()),
			),
		).
		Modify(func(selector *sql.Selector) {
			selector.OrderExpr(sql.Expr(sql.Min(selector.C(ent_robot_session_input.FieldSequence))))
		}).
		Limit(limit).
		GroupBy(ent_robot_session_input.FieldSessionID).
		Scan(ctx, &inputSessionIDs)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	var turnSessionIDs []xid.ID
	err = q.db.RobotSessionTurn.Query().
		Where(ent_robot_session_turn.StatusEQ(ent_robot_session_turn.StatusQueued)).
		Modify(func(selector *sql.Selector) {
			selector.OrderExpr(sql.Expr(sql.Min(selector.C(ent_robot_session_turn.FieldCreatedAt))))
		}).
		Limit(limit).
		GroupBy(ent_robot_session_turn.FieldSessionID).
		Scan(ctx, &turnSessionIDs)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	recoverableSessions, err := q.db.RobotSession.Query().
		Where(
			ent_robot_session.ExecutionStatusEQ(ent_robot_session.ExecutionStatusRunning),
			ent_robot_session.ActiveTurnIDNotNil(),
			ent_robot_session.Or(
				ent_robot_session.LeaseExpiresAtIsNil(),
				ent_robot_session.LeaseExpiresAtLT(time.Now()),
			),
		).
		Limit(limit).
		IDs(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	unique := make(map[xid.ID]struct{}, len(inputSessionIDs)+len(turnSessionIDs)+len(recoverableSessions))
	capacity := len(inputSessionIDs) + len(turnSessionIDs) + len(recoverableSessions)
	if capacity > limit {
		capacity = limit
	}
	result := make([]robot.SessionID, 0, capacity)
	appendID := func(id xid.ID) {
		if _, exists := unique[id]; exists || len(result) >= limit {
			return
		}
		unique[id] = struct{}{}
		result = append(result, robot.SessionID(id))
	}
	for _, id := range inputSessionIDs {
		appendID(id)
	}
	for _, id := range turnSessionIDs {
		appendID(id)
	}
	for _, id := range recoverableSessions {
		appendID(id)
	}
	return result, nil
}

func mapTurn(row *ent.RobotSessionTurn) *robot.Turn {
	return &robot.Turn{
		ID: robot.TurnID(row.ID), SessionID: robot.SessionID(row.SessionID),
		InitiatorID: opt.Map(opt.NewPtr(row.InitiatedByAccountID), func(id xid.ID) account.AccountID { return account.AccountID(id) }),
		SourceKind:  row.SourceKind, RobotRef: row.RobotRef, InputData: row.InputData,
		Status: robot.TurnStatus(row.Status), CreatedAt: row.CreatedAt,
		StartedAt: opt.NewPtr(row.StartedAt), FinishedAt: opt.NewPtr(row.FinishedAt),
		CancelRequestedAt: opt.NewPtr(row.CancelRequestedAt), ErrorText: opt.NewPtr(row.ErrorText),
	}
}

func (q *Repository) HasRunningExecution(ctx context.Context, sessionID robot.SessionID) (bool, error) {
	row, err := q.db.RobotSession.Get(ctx, xid.ID(sessionID))
	if err != nil {
		return false, fault.Wrap(err, fctx.With(ctx))
	}
	if row.ExecutionStatus == ent_robot_session.ExecutionStatusIdle {
		return false, nil
	}
	if row.ExecutionStatus == ent_robot_session.ExecutionStatusRunning && row.LeaseExpiresAt != nil && row.LeaseExpiresAt.Before(time.Now()) {
		return false, nil
	}
	return true, nil
}
