package rpc_handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	adksession "google.golang.org/adk/v2/session"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/rbac"
	robotresource "github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	robotservice "github.com/Southclaws/storyden/app/services/semdex/robot"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func (h *Handler) handleRobotRun(ctx context.Context, req *rpc.RPCRequestRobotRun) (rpc.RPCResponseRobotRun, error) {
	result := rpc.RPCResponseRobotRun{Method: "robot_run", Mode: req.Params.Mode}

	contents, err := validateRobotRunRequest(req.Params)
	if err != nil {
		result.Error = opt.New(err.Error())
		return result, nil
	}

	accessConfig, ok := h.manifest.Metadata.Access.Get()
	if !ok {
		result.Error = opt.New("manifest does not request access; define access permissions to run robots")
		return result, nil
	}

	pluginAccount, err := h.ensureAccessAccount(ctx, accessConfig)
	if err != nil {
		result.Error = opt.New(err.Error())
		return result, nil
	}
	if err := h.ensureAccessRole(ctx, pluginAccount, accessConfig); err != nil {
		result.Error = opt.New(err.Error())
		return result, nil
	}

	pluginAccount, err = h.accountQuerier.GetByID(ctx, pluginAccount.ID)
	if err != nil {
		result.Error = opt.New(err.Error())
		return result, nil
	}

	runCtx := session.WithAccessKey(ctx, pluginAccount.Account, pluginAccount.Roles.Roles())
	if err := session.Authorise(runCtx, nil, rbac.PermissionUseRobots); err != nil {
		result.Error = opt.New("plugin account does not have USE_ROBOTS permission")
		return result, nil
	}

	sessionID := xid.New()
	if requestedSessionID, ok := req.Params.SessionID.Get(); ok {
		if err := h.validateRobotRunContinuation(runCtx, requestedSessionID, pluginAccount.ID, req.Params.RobotID); err != nil {
			result.Error = opt.New(err.Error())
			return result, nil
		}
		sessionID = requestedSessionID
	}

	result.SessionID = opt.New(sessionID)

	runMode := robotservice.ModeHeadless
	if req.Params.Mode == rpc.RobotRunModeAutomation {
		runMode = robotservice.ModeUnattended
	}

	stream := h.robotAgent.RunMessages(
		runCtx,
		req.Params.RobotID,
		pluginAccount.ID.String(),
		sessionID.String(),
		contents,
		nil,
		robotservice.RunOptions{
			Mode:      runMode,
			Source:    robotservice.SourcePluginRPC,
			Workspace: robotRunWorkspaceSpec(req.Params.Workspace),
		},
	)

	var automationText strings.Builder
	finalizationReady := false
	for event, streamErr := range stream {
		if streamErr != nil {
			if req.Params.Mode == rpc.RobotRunModeAutomation && finalizationReady {
				return result, nil
			}
			setRobotRunFailure(&result, friendlyRobotRunError(streamErr), automationText.String())
			return result, nil
		}
		if event == nil || event.LLMResponse.Content == nil || event.Author == "user" {
			continue
		}

		if req.Params.Mode == rpc.RobotRunModeConversation {
			appendRobotRunEvents(&result, event)
			continue
		}

		for _, part := range event.LLMResponse.Content.Parts {
			if part == nil || part.Thought {
				continue
			}
			if part.Text != "" {
				automationText.WriteString(part.Text)
			}
			if !finalizationReady && part.FunctionResponse != nil && part.FunctionResponse.Name == robotservice.UnattendedFinishToolName() {
				finalization, err := robotRunOutputFromMap(part.FunctionResponse.Response)
				if err != nil {
					setRobotRunFailure(&result, "robot_run finish tool produced invalid structured output: "+err.Error(), automationText.String())
					finalizationReady = true
					continue
				}
				result.Finalization = opt.New(finalization)
				finalizationReady = true
			}
		}
	}

	if req.Params.Mode == rpc.RobotRunModeConversation || finalizationReady {
		return result, nil
	}

	setRobotRunFailure(&result, "robot_run did not call the unattended finish tool", automationText.String())
	return result, nil
}

func validateRobotRunRequest(params rpc.RPCRequestRobotRunParams) ([]*genai.Content, error) {
	if strings.TrimSpace(params.RobotID) == "" {
		return nil, errors.New("robot_id must not be empty")
	}
	if len(params.Messages) == 0 {
		return nil, errors.New("messages must contain at least one message")
	}
	if params.Mode != rpc.RobotRunModeConversation && params.Mode != rpc.RobotRunModeAutomation {
		return nil, fmt.Errorf("unsupported robot_run mode %q", params.Mode)
	}
	if params.Messages[len(params.Messages)-1].Role != rpc.RobotRunMessageRoleUser {
		return nil, errors.New("the final message must have role user")
	}
	if params.Mode == rpc.RobotRunModeAutomation {
		if len(params.Messages) != 1 {
			return nil, errors.New("automation mode requires exactly one user message")
		}
		if params.SessionID.Ok() {
			return nil, errors.New("automation mode does not accept session_id")
		}
	}

	contents := make([]*genai.Content, len(params.Messages))
	for i, message := range params.Messages {
		if strings.TrimSpace(message.Content) == "" {
			return nil, fmt.Errorf("messages[%d].content must not be empty", i)
		}
		author, hasAuthor := message.Author.Get()
		if hasAuthor && strings.TrimSpace(author) == "" {
			return nil, fmt.Errorf("messages[%d].author must not be empty", i)
		}

		role := genai.Role(genai.RoleUser)
		switch message.Role {
		case rpc.RobotRunMessageRoleUser:
		case rpc.RobotRunMessageRoleAssistant:
			role = genai.RoleModel
		default:
			return nil, fmt.Errorf("messages[%d].role is invalid", i)
		}
		contents[i] = robotservice.ContentWithSpeaker(role, message.Content, author)
	}

	return contents, nil
}

func (h *Handler) validateRobotRunContinuation(ctx context.Context, sessionID xid.ID, accountID account.AccountID, robotRef string) error {
	sess, _, err := h.robotSessions.Get(ctx, robotresource.SessionID(sessionID), robotresource.NewMessageCursorParams(opt.NewEmpty[robotresource.MessageID](), 1))
	if err != nil || sess.Human.ID != accountID {
		return errors.New("session is not available to this plugin")
	}
	rootRobotRef, ok := robotservice.SessionRootRobotRef(sess.State).Get()
	if !ok || rootRobotRef != strings.TrimSpace(robotRef) {
		return errors.New("session is rooted at a different Robot")
	}
	return nil
}

func appendRobotRunEvents(result *rpc.RPCResponseRobotRun, event *adksession.Event) {
	if event.LLMResponse.Content.Role == genai.RoleModel {
		result.FinalText = opt.NewEmpty[string]()
	}
	var finalText strings.Builder
	for _, part := range event.LLMResponse.Content.Parts {
		if part == nil || part.Thought {
			continue
		}
		switch {
		case part.Text != "":
			result.Events = append(result.Events, rpc.RobotRunEvent{RobotRunEventUnion: &rpc.RobotRunTextEvent{Type: "text", Text: part.Text}})
			if event.LLMResponse.Content.Role == genai.RoleModel {
				finalText.WriteString(part.Text)
			}
		case part.FunctionCall != nil:
			result.Events = append(result.Events, rpc.RobotRunEvent{RobotRunEventUnion: &rpc.RobotRunToolCallEvent{
				Type: "tool_call", Name: part.FunctionCall.Name, Arguments: ensureRobotRunMap(part.FunctionCall.Args),
				CallID: opt.NewIf(part.FunctionCall.ID, func(value string) bool { return value != "" }),
			}})
		case part.FunctionResponse != nil:
			result.Events = append(result.Events, rpc.RobotRunEvent{RobotRunEventUnion: &rpc.RobotRunToolResultEvent{
				Type: "tool_result", Name: part.FunctionResponse.Name, Result: ensureRobotRunMap(part.FunctionResponse.Response),
				CallID: opt.NewIf(part.FunctionResponse.ID, func(value string) bool { return value != "" }),
			}})
		}
	}
	if finalText.Len() > 0 {
		result.FinalText = opt.New(finalText.String())
	}
}

func ensureRobotRunMap(value map[string]any) map[string]any {
	if value == nil {
		return map[string]any{}
	}
	return value
}

func friendlyRobotRunError(err error) string {
	if message := fmsg.GetIssue(err); message != "" {
		return message
	}
	return err.Error()
}

func robotRunWorkspaceSpec(in opt.Optional[rpc.RPCRequestRobotRunParamsWorkspace]) opt.Optional[robotservice.WorkspaceMountSpec] {
	workspace, ok := in.Get()
	if !ok {
		return opt.NewEmpty[robotservice.WorkspaceMountSpec]()
	}

	spec := robotservice.WorkspaceMountSpec{Metadata: map[string]any{}}
	if id, ok := workspace.WorkspaceID.Get(); ok {
		spec.WorkspaceID = opt.New(robotresource.WorkspaceID(id))
	}
	if id, ok := workspace.WorkspaceInstanceID.Get(); ok {
		spec.WorkspaceInstanceID = opt.New(robotresource.WorkspaceInstanceID(id))
	}
	return opt.New(spec)
}

func robotRunOutputFromMap(data map[string]any) (rpc.RobotRunOutput, error) {
	var output rpc.RobotRunOutput
	if data == nil {
		return output, errors.New("empty final output")
	}
	raw, err := json.Marshal(data)
	if err != nil {
		return rpc.RobotRunOutput{}, err
	}
	if err := json.Unmarshal(raw, &output); err != nil {
		return rpc.RobotRunOutput{}, err
	}
	if output.Status == "" {
		return rpc.RobotRunOutput{}, errors.New("missing status")
	}
	if strings.TrimSpace(output.Summary) == "" {
		return rpc.RobotRunOutput{}, errors.New("missing summary")
	}
	return output, nil
}

func setRobotRunFailure(result *rpc.RPCResponseRobotRun, errmsg, summary string) {
	result.Error = opt.New(errmsg)
	if result.Mode != rpc.RobotRunModeAutomation {
		return
	}
	if strings.TrimSpace(summary) == "" {
		summary = errmsg
	}
	result.Finalization = opt.New(rpc.RobotRunOutput{
		Status:  rpc.RobotRunStatusFailed,
		Summary: summary,
		Attention: opt.New(rpc.RobotRunAttention{
			Reason:  rpc.RobotRunAttentionReasonError,
			Message: errmsg,
		}),
	})
}
