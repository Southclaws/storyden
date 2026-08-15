package robot_session

import (
	"context"
	"errors"
	"maps"
	"time"

	"github.com/google/uuid"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/predicate"
	ent_robot_session "github.com/Southclaws/storyden/internal/ent/robotsession"
)

var (
	ErrSessionBusy = errors.New("robot session is already running")
	ErrLeaseLost   = errors.New("robot session execution lease was lost")
)

type ExecutionLease struct {
	SessionID  robot.SessionID
	TurnID     xid.ID
	Token      string
	Generation uint64
	ExpiresAt  time.Time
}

type executionLeaseContextKey struct{}

func WithExecutionLease(ctx context.Context, lease ExecutionLease) context.Context {
	return context.WithValue(ctx, executionLeaseContextKey{}, lease)
}

func ExecutionLeaseFromContext(ctx context.Context) (ExecutionLease, bool) {
	lease, ok := ctx.Value(executionLeaseContextKey{}).(ExecutionLease)
	return lease, ok
}

func (q *Repository) AcquireExecution(
	ctx context.Context,
	sessionID robot.SessionID,
	turnID xid.ID,
	allowBlocked bool,
	duration time.Duration,
) (*ExecutionLease, error) {
	now := time.Now()
	lease := &ExecutionLease{
		SessionID: sessionID,
		TurnID:    turnID,
		Token:     uuid.NewString(),
		ExpiresAt: now.Add(duration),
	}

	eligible := []predicate.RobotSession{
		ent_robot_session.ExecutionStatusEQ(ent_robot_session.ExecutionStatusIdle),
		ent_robot_session.And(
			ent_robot_session.ExecutionStatusEQ(ent_robot_session.ExecutionStatusRunning),
			ent_robot_session.LeaseExpiresAtLT(now),
		),
	}
	if allowBlocked {
		eligible = append(eligible, ent_robot_session.ExecutionStatusEQ(ent_robot_session.ExecutionStatusBlocked))
	}

	updated, err := q.db.RobotSession.Update().
		Where(
			ent_robot_session.IDEQ(xid.ID(sessionID)),
			ent_robot_session.Or(eligible...),
		).
		SetExecutionStatus(ent_robot_session.ExecutionStatusRunning).
		SetActiveTurnID(turnID).
		SetLeaseToken(lease.Token).
		AddLeaseGeneration(1).
		SetLeaseExpiresAt(lease.ExpiresAt).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	if updated != 1 {
		return nil, ErrSessionBusy
	}

	session, err := q.db.RobotSession.Query().
		Where(
			ent_robot_session.IDEQ(xid.ID(sessionID)),
			ent_robot_session.ActiveTurnIDEQ(turnID),
			ent_robot_session.LeaseTokenEQ(lease.Token),
		).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	lease.Generation = session.LeaseGeneration

	if allowBlocked {
		state := maps.Clone(session.State)
		if state == nil {
			state = map[string]any{}
		}
		state = robot.ClearPendingClientTools(state)
		cleared, err := q.db.RobotSession.Update().
			Where(executionLeasePredicates(*lease)...).
			SetState(state).
			Save(ctx)
		if err != nil {
			return nil, err
		}
		if cleared != 1 {
			return nil, ErrLeaseLost
		}
	}
	return lease, nil
}

func (q *Repository) HeartbeatExecution(ctx context.Context, lease ExecutionLease, duration time.Duration) (bool, error) {
	expiresAt := time.Now().Add(duration)
	updated, err := q.db.RobotSession.Update().
		Where(executionLeasePredicates(lease)...).
		SetLeaseExpiresAt(expiresAt).
		Save(ctx)
	if err != nil {
		return false, err
	}
	return updated == 1, nil
}

func (q *Repository) ReleaseExecution(ctx context.Context, lease ExecutionLease) error {
	updated, err := q.db.RobotSession.Update().
		Where(executionLeasePredicates(lease)...).
		SetExecutionStatus(ent_robot_session.ExecutionStatusIdle).
		ClearActiveTurnID().
		ClearLeaseToken().
		ClearLeaseExpiresAt().
		Save(ctx)
	if err != nil {
		return err
	}
	if updated != 1 {
		return ErrLeaseLost
	}
	return nil
}

func (q *Repository) BlockExecution(ctx context.Context, lease ExecutionLease) error {
	updated, err := q.db.RobotSession.Update().
		Where(executionLeasePredicates(lease)...).
		SetExecutionStatus(ent_robot_session.ExecutionStatusBlocked).
		ClearLeaseToken().
		ClearLeaseExpiresAt().
		Save(ctx)
	if err != nil {
		return err
	}
	if updated != 1 {
		return ErrLeaseLost
	}
	return nil
}

func (q *Repository) StorePendingClientTools(ctx context.Context, sessionID robot.SessionID, ids []string) error {
	lease, ok := ExecutionLeaseFromContext(ctx)
	if !ok || lease.SessionID != sessionID {
		return ErrLeaseLost
	}
	return ent.WithTx(ctx, q.db, func(tx *ent.Tx) error {
		if err := validateExecutionLease(ctx, tx, lease); err != nil {
			return err
		}
		session, err := tx.RobotSession.Query().
			Where(ent_robot_session.IDEQ(xid.ID(sessionID))).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return ErrLeaseLost
			}
			return err
		}
		state := maps.Clone(session.State)
		if state == nil {
			state = map[string]any{}
		}
		seen := map[string]struct{}{}
		pending := robot.PendingClientToolIDs(state)
		for _, id := range pending {
			seen[id] = struct{}{}
		}
		for _, id := range ids {
			if id == "" {
				continue
			}
			if _, exists := seen[id]; exists {
				continue
			}
			seen[id] = struct{}{}
			pending = append(pending, id)
		}
		state[robot.PendingClientToolsStateKey] = pending
		updated, err := tx.RobotSession.Update().
			Where(activeExecutionLeasePredicates(lease)...).
			SetState(state).
			Save(ctx)
		if err != nil {
			return err
		}
		if updated != 1 {
			return ErrLeaseLost
		}
		return nil
	})
}

func executionLeasePredicates(lease ExecutionLease) []predicate.RobotSession {
	return []predicate.RobotSession{
		ent_robot_session.IDEQ(xid.ID(lease.SessionID)),
		ent_robot_session.ExecutionStatusEQ(ent_robot_session.ExecutionStatusRunning),
		ent_robot_session.ActiveTurnIDEQ(lease.TurnID),
		ent_robot_session.LeaseTokenEQ(lease.Token),
		ent_robot_session.LeaseGenerationEQ(lease.Generation),
	}
}

func activeExecutionLeasePredicates(lease ExecutionLease) []predicate.RobotSession {
	return append(executionLeasePredicates(lease), ent_robot_session.LeaseExpiresAtGT(time.Now()))
}

func validateExecutionLease(ctx context.Context, tx *ent.Tx, lease ExecutionLease) error {
	// The conditional no-op update both validates the fencing token at mutation
	// time and takes the row write lock for the remainder of the transaction.
	// A takeover therefore cannot interleave between validation and a message or
	// state write by the old execution.
	updated, err := tx.RobotSession.Update().
		Where(activeExecutionLeasePredicates(lease)...).
		SetLeaseToken(lease.Token).
		Save(ctx)
	if err != nil {
		return err
	}
	if updated != 1 {
		return ErrLeaseLost
	}
	return nil
}
