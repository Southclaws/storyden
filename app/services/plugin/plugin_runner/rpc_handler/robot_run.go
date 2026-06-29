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
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/rbac"
	robotresource "github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	robotservice "github.com/Southclaws/storyden/app/services/semdex/robot"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func (h *Handler) handleRobotRun(ctx context.Context, req *rpc.RPCRequestRobotRun) (rpc.RPCResponseRobotRun, error) {
	result := rpc.RPCResponseRobotRun{
		Method: "robot_run",
	}
	sessionID := xid.New()
	result.SessionID = opt.New(sessionID)

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

	// Refresh role edges after potential role assignment updates.
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
	if strings.TrimSpace(req.Params.Message) == "" {
		result.Error = opt.New("message must not be empty")
		return result, nil
	}

	stream := h.robotAgent.Run(
		runCtx,
		req.Params.RobotID,
		pluginAccount.ID.String(),
		sessionID.String(),
		&genai.Content{
			Role: "user",
			Parts: []*genai.Part{
				{Text: req.Params.Message},
			},
		},
		nil,
		robotservice.RunOptions{
			Mode:      robotservice.ModeUnattended,
			Source:    robotservice.SourcePluginRPC,
			Workspace: robotRunWorkspaceSpec(req.Params.Workspace),
		},
	)

	var output strings.Builder
	for event, streamErr := range stream {
		if streamErr != nil {
			message := fmsg.GetIssue(streamErr)
			if message == "" {
				message = streamErr.Error()
			}
			setRobotRunFailure(&result, message, output.String())
			return result, nil
		}
		if event == nil || event.LLMResponse.Content == nil {
			continue
		}

		var eventText strings.Builder
		for _, part := range event.LLMResponse.Content.Parts {
			if part == nil {
				continue
			}
			if part.Text != "" {
				eventText.WriteString(part.Text)
			}
			if part.FunctionResponse != nil && part.FunctionResponse.Name == robotservice.UnattendedFinishToolName() {
				runOutput, err := robotRunOutputFromMap(part.FunctionResponse.Response)
				if err != nil {
					message := "robot_run finish tool produced invalid structured output: " + err.Error()
					setRobotRunFailure(&result, message, output.String())
					return result, nil
				}
				result.Output = opt.New(runOutput)
				return result, nil
			}
		}

		text := eventText.String()
		if text == "" {
			continue
		}
		output.WriteString(text)
	}

	message := "robot_run did not call the unattended finish tool"
	setRobotRunFailure(&result, message, output.String())
	return result, nil
}

func robotRunWorkspaceSpec(in opt.Optional[rpc.RPCRequestRobotRunParamsWorkspace]) opt.Optional[robotservice.WorkspaceMountSpec] {
	workspace, ok := in.Get()
	if !ok {
		return opt.NewEmpty[robotservice.WorkspaceMountSpec]()
	}

	spec := robotservice.WorkspaceMountSpec{
		Metadata: map[string]any{},
	}
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
		return rpc.RobotRunOutput{}, fmt.Errorf("missing status")
	}
	if strings.TrimSpace(output.Summary) == "" {
		return rpc.RobotRunOutput{}, fmt.Errorf("missing summary")
	}
	return output, nil
}

func setRobotRunFailure(result *rpc.RPCResponseRobotRun, errmsg, summary string) {
	result.Error = opt.New(errmsg)
	result.Output = opt.New(rpc.RobotRunOutput{
		Status:  rpc.RobotRunStatusFailed,
		Summary: summary,
		Attention: opt.New(rpc.RobotRunAttention{
			Reason:  rpc.RobotRunAttentionReasonError,
			Message: summary,
		}),
	})
}
