package rpc_handler

import (
	"context"
	"strings"

	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func (h *Handler) handleRobotRun(ctx context.Context, req *rpc.RPCRequestRobotRun) (rpc.RPCResponseRobotRun, error) {
	result := rpc.RPCResponseRobotRun{
		Method: "robot_run",
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
		opt.New(req.Params.RobotID),
		pluginAccount.ID.String(),
		xid.New().String(),
		&genai.Content{
			Role: "user",
			Parts: []*genai.Part{
				{Text: req.Params.Message},
			},
		},
		nil,
	)

	var output strings.Builder
	for event, streamErr := range stream {
		if streamErr != nil {
			message := fmsg.GetIssue(streamErr)
			if message == "" {
				message = streamErr.Error()
			}
			result.Error = opt.New(message)
			return result, nil
		}
		if event == nil || event.LLMResponse.Content == nil {
			continue
		}

		for _, part := range event.LLMResponse.Content.Parts {
			if part == nil || part.Text == "" {
				continue
			}
			output.WriteString(part.Text)
		}
	}

	result.Response = opt.New(output.String())
	return result, nil
}
