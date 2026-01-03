package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_querier"
	"github.com/Southclaws/storyden/app/resources/robot/robot_ref"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/app/resources/robot/robot_writer"
	"github.com/Southclaws/storyden/app/resources/robot/session_ref"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Robots struct {
	robotQuerier *robot_querier.Querier
	robotWriter  *robot_writer.Writer
	sessionRepo  *robot_session.Repository
}

func NewRobots(
	robotQuerier *robot_querier.Querier,
	robotWriter *robot_writer.Writer,
	sessionRepo *robot_session.Repository,
) Robots {
	return Robots{
		robotQuerier: robotQuerier,
		robotWriter:  robotWriter,
		sessionRepo:  sessionRepo,
	}
}

func (r *Robots) RobotsList(ctx context.Context, request openapi.RobotsListRequestObject) (openapi.RobotsListResponseObject, error) {
	pageParams := deserialisePageParams(request.Params.Page, 20)

	result, err := r.robotQuerier.List(ctx, pageParams)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotsList200JSONResponse{
		RobotsListOKJSONResponse: openapi.RobotsListOKJSONResponse(openapi.RobotsListResult{
			CurrentPage: result.CurrentPage,
			NextPage:    result.NextPage.Ptr(),
			PageSize:    result.Size,
			Results:     result.Results,
			Robots:      serialiseRobots(result.Items),
			TotalPages:  result.TotalPages,
		}),
	}, nil
}

func (r *Robots) RobotCreate(ctx context.Context, request openapi.RobotCreateRequestObject) (openapi.RobotCreateResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	opts := []robot_writer.Option{}
	if meta := request.Body.Meta; meta != nil {
		opts = append(opts, robot_writer.WithMeta(*meta))
	}

	created, err := r.robotWriter.Create(ctx,
		request.Body.Name,
		request.Body.Description,
		request.Body.Playbook,
		accountID,
		opts...,
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotCreate200JSONResponse{
		RobotCreateOKJSONResponse: openapi.RobotCreateOKJSONResponse(serialiseRobot(created)),
	}, nil
}

func (r *Robots) RobotGet(ctx context.Context, request openapi.RobotGetRequestObject) (openapi.RobotGetResponseObject, error) {
	robotID, err := robot_ref.NewID(request.RobotId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	rb, err := r.robotQuerier.Get(ctx, robotID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotGet200JSONResponse{
		RobotGetOKJSONResponse: openapi.RobotGetOKJSONResponse(serialiseRobot(rb)),
	}, nil
}

func (r *Robots) RobotUpdate(ctx context.Context, request openapi.RobotUpdateRequestObject) (openapi.RobotUpdateResponseObject, error) {
	robotID, err := robot_ref.NewID(request.RobotId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	opts := []robot_writer.Option{}
	if name := request.Body.Name; name != nil {
		opts = append(opts, robot_writer.WithName(*name))
	}
	if description := request.Body.Description; description != nil {
		opts = append(opts, robot_writer.WithDescription(*description))
	}
	if playbook := request.Body.Playbook; playbook != nil {
		opts = append(opts, robot_writer.WithPlaybook(*playbook))
	}
	if meta := request.Body.Meta; meta != nil {
		opts = append(opts, robot_writer.WithMeta(*meta))
	}

	updated, err := r.robotWriter.Update(ctx, robotID, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotUpdate200JSONResponse{
		RobotGetOKJSONResponse: openapi.RobotGetOKJSONResponse(serialiseRobot(updated)),
	}, nil
}

func (r *Robots) RobotSessionsList(ctx context.Context, request openapi.RobotSessionsListRequestObject) (openapi.RobotSessionsListResponseObject, error) {
	pageParams := deserialisePageParams(request.Params.Page, 20)

	var accountID opt.Optional[account.AccountID]
	if request.Params.AccountId != nil {
		id, err := xid.FromString(string(*request.Params.AccountId))
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		aid := account.AccountID(id)

		accountID = opt.New(aid)
	}

	result, err := r.sessionRepo.List(ctx, pageParams, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sessions := dt.Map(result.Items, func(s *session_ref.Ref) openapi.RobotSessionRef {
		return openapi.RobotSessionRef{
			Id:        openapi.Identifier(s.ID.String()),
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
			CreatedBy: serialiseProfileReferenceFromAccount(s.Human),
		}
	})

	return openapi.RobotSessionsList200JSONResponse{
		RobotSessionsListOKJSONResponse: openapi.RobotSessionsListOKJSONResponse(openapi.RobotSessionsListResult{
			CurrentPage: result.CurrentPage,
			NextPage:    result.NextPage.Ptr(),
			PageSize:    result.Size,
			Results:     result.Results,
			Sessions:    sessions,
			TotalPages:  result.TotalPages,
		}),
	}, nil
}

func (r *Robots) RobotSessionGet(ctx context.Context, request openapi.RobotSessionGetRequestObject) (openapi.RobotSessionGetResponseObject, error) {
	sessionID, err := robot.NewSessionID(request.SessionId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	messageParams := deserialisePageParams(request.Params.Page, 50)

	session, pagination, err := r.sessionRepo.Get(ctx, robot.SessionID(sessionID), messageParams)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// Convert messages to Vercel AI SDK UIMessage format
	messageDTOs, err := dt.MapErr(pagination.Items, serialiseRobotSessionMessage)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RobotSessionGet200JSONResponse{
		RobotSessionGetOKJSONResponse: openapi.RobotSessionGetOKJSONResponse(openapi.RobotSession{
			Id:        openapi.Identifier(session.ID.String()),
			CreatedAt: session.CreatedAt,
			UpdatedAt: session.UpdatedAt,
			CreatedBy: serialiseProfileReferenceFromAccount(session.Human),
			MessageList: openapi.PaginatedRobotMessageList{
				CurrentPage: pagination.CurrentPage,
				NextPage:    pagination.NextPage.Ptr(),
				PageSize:    pagination.Size,
				Results:     pagination.Results,
				TotalPages:  pagination.TotalPages,
				Messages:    messageDTOs,
			},
		}),
	}, nil
}

func serialiseRobot(r *robot.Robot) openapi.Robot {
	return openapi.Robot{
		Id:          openapi.Identifier(r.ID.String()),
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		Name:        r.Name,
		Description: r.Description,
		Playbook:    r.Playbook,
		Author:      serialiseProfileReferenceFromAccount(r.Author),
		Tools:       r.Tools,
		Meta:        (*openapi.Metadata)(&r.Metadata),
	}
}

func serialiseRobots(robots []*robot.Robot) []openapi.Robot {
	return dt.Map(robots, func(r *robot.Robot) openapi.Robot {
		return serialiseRobot(r)
	})
}
