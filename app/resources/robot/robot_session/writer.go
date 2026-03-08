package robot_session

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/internal/ent"
	ent_robot_session "github.com/Southclaws/storyden/internal/ent/robotsession"
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
	invocationID string,
	accountID opt.Optional[account.AccountID],
	robotID opt.Optional[xid.ID],
	eventData map[string]any,
) error {
	create := q.db.RobotSessionMessage.Create().
		SetSessionID(xid.ID(sessionID)).
		SetInvocationID(invocationID).
		SetEventData(eventData)

	if aid, ok := accountID.Get(); ok {
		create.SetAccountID(xid.ID(aid))
	}

	if rid, ok := robotID.Get(); ok {
		create.SetRobotID(rid)
	}

	_, err := create.Save(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
