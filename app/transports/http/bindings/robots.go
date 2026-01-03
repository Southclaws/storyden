package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_querier"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Robots struct {
	robotQuerier *robot_querier.Querier
}

func NewRobots(robotQuerier *robot_querier.Querier) Robots {
	return Robots{
		robotQuerier: robotQuerier,
	}
}

func (r *Robots) RobotsList(ctx context.Context, request openapi.RobotsListRequestObject) (openapi.RobotsListResponseObject, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

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

func serialiseRobot(r *robot.Robot) openapi.Robot {
	return openapi.Robot{
		Id:          openapi.Identifier(r.ID.String()),
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		Name:        r.Name,
		Description: r.Description,
		Playbook:    r.Playbook,
		Author:      serialiseProfileReferenceFromAccount(r.Author),
		Meta:        (*openapi.Metadata)(&r.Metadata),
	}
}

func serialiseRobots(robots []*robot.Robot) []openapi.Robot {
	return dt.Map(robots, func(r *robot.Robot) openapi.Robot {
		return serialiseRobot(r)
	})
}
