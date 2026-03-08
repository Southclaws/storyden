package bindings

import (
	"context"
	"fmt"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	adksession "google.golang.org/adk/session"
	"google.golang.org/genai"

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

func (r *Robots) RobotChatSSE(ctx context.Context, request openapi.RobotChatSSERequestObject) (openapi.RobotChatSSEResponseObject, error) {
	return nil, fault.New("bindings layer should not be hit, this is a bug.", fctx.With(ctx))
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
			Name:      s.Name,
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
			Name:      session.Name,
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

// serialiseRobotSessionMessage converts a robot.Message (containing ADK Event)
// into the Vercel AI SDK UIMessage format for the frontend.
func serialiseRobotSessionMessage(m *robot.Message) (openapi.RobotSessionMessage, error) {
	parts, err := serialiseADKEventToParts(m.Event)
	if err != nil {
		return openapi.RobotSessionMessage{}, err
	}

	role := serialiseMessageRole(m.Event)

	msg := openapi.RobotSessionMessage{
		Id:        m.ID.String(),
		Role:      openapi.RobotSessionMessageRole(role),
		Parts:     parts,
		CreatedAt: m.CreatedAt,
		Robot:     opt.Map(m.Robot, func(r *robot.Robot) openapi.Robot { return serialiseRobot(r) }).Ptr(),
		Author:    opt.Map(m.Author, func(a *account.Account) openapi.ProfileReference { return serialiseProfileReferenceFromAccount(*a) }).Ptr(),
	}

	return msg, nil
}

func serialiseMessageRole(event adksession.Event) string {
	if event.Author == "user" {
		return "user"
	}
	return "assistant"
}

// serialiseADKEventToParts converts ADK Event.LLMResponse.Content.Parts into UIMessagePart[]
// Since UIMessagePart is a discriminated union represented as json.RawMessage,
// we need to marshal each concrete type to JSON.
func serialiseADKEventToParts(event adksession.Event) ([]openapi.UIMessagePart, error) {
	if event.LLMResponse.Content == nil {
		return []openapi.UIMessagePart{}, nil
	}

	var parts []openapi.UIMessagePart

	for _, adkPart := range event.LLMResponse.Content.Parts {
		if adkPart == nil {
			continue
		}

		if adkPart.Text != "" {
			textPart := openapi.TextUIPart{
				Type:  openapi.TextUIPartType("text"),
				Text:  adkPart.Text,
				State: ptr(openapi.TextUIPartState("done")), // Historical messages are always "done"
			}
			var uiPart openapi.UIMessagePart
			if err := uiPart.FromTextUIPart(textPart); err != nil {
				return nil, fmt.Errorf("create text part: %w", err)
			}
			parts = append(parts, uiPart)
		}

		if adkPart.FunctionCall != nil {
			uiPart, err := serialiseFunctionCallToPart(adkPart.FunctionCall)
			if err != nil {
				return nil, err
			}
			parts = append(parts, uiPart)
		}

		if adkPart.FunctionResponse != nil {
			uiPart, err := serialiseFunctionResponseToPart(adkPart.FunctionResponse)
			if err != nil {
				return nil, err
			}
			parts = append(parts, uiPart)
		}
	}

	return parts, nil
}

// serialiseFunctionCallToPart converts an ADK FunctionCall into a UIMessagePart (input-available state)
func serialiseFunctionCallToPart(fc *genai.FunctionCall) (openapi.UIMessagePart, error) {
	inputAvailable := openapi.ToolUIPartInputAvailable{
		ToolCallId: fc.ID,
		ToolName:   fc.Name,
		State:      openapi.InputAvailable,
		Input:      fc.Args,
	}

	var toolPart openapi.ToolUIPart
	if err := toolPart.FromToolUIPartInputAvailable(inputAvailable); err != nil {
		return openapi.UIMessagePart{}, fmt.Errorf("create tool input part: %w", err)
	}

	var uiPart openapi.UIMessagePart
	if err := uiPart.FromToolUIPart(toolPart); err != nil {
		return openapi.UIMessagePart{}, fmt.Errorf("create UI message part from tool part: %w", err)
	}

	// NOTE: This is sort of a hack around OpenAPI's inability to express the
	// weird Vercel AI SDK types for UIMessagePart tools. For some reason they
	// overload the "type" field to also include the tool's name, prefixing it
	// with "tool-". This means we have to manually set it here which breaks the
	// OpenAPI contract a bit. Couldn't think of a better way to express it.
	uiPart.Type = openapi.UIMessagePartType("tool-" + fc.Name)

	return uiPart, nil
}

// serialiseFunctionResponseToPart converts an ADK FunctionResponse into a UIMessagePart (output-available state)
func serialiseFunctionResponseToPart(fr *genai.FunctionResponse) (openapi.UIMessagePart, error) {
	outputAvailable := openapi.ToolUIPartOutputAvailable{
		ToolCallId: fr.ID,
		ToolName:   fr.Name,
		State:      openapi.OutputAvailable,
		Input:      fr.Response, // ADK stores the original input separately if needed
		Output:     fr.Response,
	}

	var toolPart openapi.ToolUIPart
	if err := toolPart.FromToolUIPartOutputAvailable(outputAvailable); err != nil {
		return openapi.UIMessagePart{}, fmt.Errorf("create tool output part: %w", err)
	}

	var uiPart openapi.UIMessagePart
	if err := uiPart.FromToolUIPart(toolPart); err != nil {
		return openapi.UIMessagePart{}, fmt.Errorf("create UI message part from tool part: %w", err)
	}

	// NOTE: Same as above.
	uiPart.Type = openapi.UIMessagePartType("tool-" + fr.Name)

	return openapi.UIMessagePart(uiPart), nil
}

// ptr is a helper to get a pointer to a value
func ptr[T any](v T) *T {
	return &v
}
