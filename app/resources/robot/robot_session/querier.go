package robot_session

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/session_ref"
	"github.com/Southclaws/storyden/internal/ent"
	ent_robot_session "github.com/Southclaws/storyden/internal/ent/robotsession"
	ent_robot_session_message "github.com/Southclaws/storyden/internal/ent/robotsessionmessage"
)

func (q *Repository) List(
	ctx context.Context,
	params pagination.Parameters,
	accountID opt.Optional[account.AccountID],
) (*pagination.Result[*session_ref.Ref], error) {
	query := q.db.RobotSession.Query().
		WithUser()

	if aid, ok := accountID.Get(); ok {
		query.Where(ent_robot_session.AccountIDEQ(xid.ID(aid)))
	}

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	query.Limit(params.Limit()).Offset(params.Offset())
	query.Order(ent_robot_session.ByCreatedAt(sql.OrderDesc()))

	sessions, err := query.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	refs, err := dt.MapErr(sessions, robot.MapSessionRef)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := pagination.NewPageResult(params, total, refs)
	return &result, nil
}

func (q *Repository) Get(
	ctx context.Context,
	sessionID robot.SessionID,
	messageParams pagination.Parameters,
) (*robot.Session, *pagination.Result[*robot.Message], error) {
	session, err := q.db.RobotSession.Query().
		Where(ent_robot_session.IDEQ(xid.ID(sessionID))).
		WithUser().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	messageQuery := q.db.RobotSessionMessage.Query().
		Where(ent_robot_session_message.SessionIDEQ(xid.ID(sessionID))).
		WithRobot(func(rq *ent.RobotQuery) {
			rq.WithAuthor()
		}).
		WithAuthor().
		Order(ent_robot_session_message.ByCreatedAt(sql.OrderAsc()))

	totalMessages, err := messageQuery.Clone().Count(ctx)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	messageQuery.Limit(messageParams.Limit()).Offset(messageParams.Offset())

	messages, err := messageQuery.All(ctx)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	sess, err := robot.MapSession(session, messages)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	paginationResult := pagination.NewPageResult(messageParams, totalMessages, sess.Messages)

	return sess, &paginationResult, nil
}

func (q *Repository) GetWithMessageFilters(
	ctx context.Context,
	sessionID robot.SessionID,
	accountID account.AccountID,
	numRecentEvents int,
	after time.Time,
) (*robot.Session, error) {
	session, err := q.db.RobotSession.Query().
		Where(
			ent_robot_session.IDEQ(xid.ID(sessionID)),
			ent_robot_session.AccountIDEQ(xid.ID(accountID)),
		).
		WithUser().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	messageQuery := q.db.RobotSessionMessage.Query().
		Where(ent_robot_session_message.SessionIDEQ(xid.ID(sessionID))).
		WithRobot(func(rq *ent.RobotQuery) {
			rq.WithAuthor()
		}).
		WithAuthor().
		Order(ent_robot_session_message.ByCreatedAt(sql.OrderAsc()))

	if numRecentEvents > 0 {
		messageQuery.Limit(numRecentEvents)
	}
	if !after.IsZero() {
		messageQuery.Where(ent_robot_session_message.CreatedAtGTE(after))
	}

	messages, err := messageQuery.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sess, err := robot.MapSession(session, messages)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return sess, nil
}

func (q *Repository) ListAll(ctx context.Context, accountID opt.Optional[account.AccountID]) ([]*robot.Session, error) {
	query := q.db.RobotSession.Query().WithUser()

	if aid, ok := accountID.Get(); ok {
		query.Where(ent_robot_session.AccountIDEQ(xid.ID(aid)))
	}

	sessions, err := query.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return dt.MapErr(sessions, func(s *ent.RobotSession) (*robot.Session, error) {
		user, err := account.MapRef(s.Edges.User)
		if err != nil {
			return nil, err
		}

		return &robot.Session{
			Ref: session_ref.Ref{
				ID:        session_ref.ID(s.ID),
				CreatedAt: s.CreatedAt,
				UpdatedAt: s.UpdatedAt,
				Human:     *user,
			},
			State: s.State,
		}, nil
	})
}
