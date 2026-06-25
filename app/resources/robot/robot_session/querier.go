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
	messageParams robot.MessageCursorParams,
) (*robot.Session, *robot.MessageCursorResult, error) {
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
		Order(
			ent_robot_session_message.ByCreatedAt(sql.OrderDesc()),
			ent_robot_session_message.ByID(sql.OrderDesc()),
		)

	if before, ok := messageParams.Before.Get(); ok {
		cursor, err := q.db.RobotSessionMessage.Query().
			Where(
				ent_robot_session_message.IDEQ(xid.ID(before)),
				ent_robot_session_message.SessionIDEQ(xid.ID(sessionID)),
			).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				err = fault.Wrap(err, ftag.With(ftag.InvalidArgument))
			}
			return nil, nil, fault.Wrap(err, fctx.With(ctx))
		}

		messageQuery.Where(
			ent_robot_session_message.Or(
				ent_robot_session_message.CreatedAtLT(cursor.CreatedAt),
				ent_robot_session_message.And(
					ent_robot_session_message.CreatedAtEQ(cursor.CreatedAt),
					ent_robot_session_message.IDLT(cursor.ID),
				),
			),
		)
	}

	messageQuery.Limit(messageParams.Limit() + 1)

	messages, err := messageQuery.All(ctx)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	moreResults := len(messages) > messageParams.Limit()
	if moreResults {
		messages = messages[:messageParams.Limit()]
	}

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	nextBefore := opt.NewEmpty[robot.MessageID]()
	if moreResults && len(messages) > 0 {
		nextBefore = opt.New(robot.MessageID(messages[0].ID))
	}

	mappedMessages, err := dt.MapErr(messages, robot.MapMessage)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	sess, err := robot.MapSession(session, nil)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}
	sess.Messages = mappedMessages

	return sess, &robot.MessageCursorResult{
		Size:       messageParams.Size,
		Results:    len(mappedMessages),
		NextBefore: nextBefore,
		Items:      mappedMessages,
	}, nil
}

func (q *Repository) GetWithMessageFilters(
	ctx context.Context,
	sessionID robot.SessionID,
	accountID opt.Optional[account.AccountID],
	numRecentEvents int,
	after time.Time,
) (*robot.Session, error) {
	sessionQuery := q.db.RobotSession.Query().
		Where(ent_robot_session.IDEQ(xid.ID(sessionID))).
		WithUser()

	if aid, ok := accountID.Get(); ok {
		sessionQuery.Where(ent_robot_session.AccountIDEQ(xid.ID(aid)))
	}

	session, err := sessionQuery.Only(ctx)
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
		WithAuthor()

	if numRecentEvents > 0 {
		// Fetch the N most recent by ordering DESC then reversing, so callers
		// always receive messages in ascending chronological order.
		messageQuery.Order(ent_robot_session_message.ByCreatedAt(sql.OrderDesc())).Limit(numRecentEvents)
	} else {
		messageQuery.Order(ent_robot_session_message.ByCreatedAt(sql.OrderAsc()))
	}

	if !after.IsZero() {
		messageQuery.Where(ent_robot_session_message.CreatedAtGTE(after))
	}

	messages, err := messageQuery.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if numRecentEvents > 0 {
		for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
			messages[i], messages[j] = messages[j], messages[i]
		}
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
