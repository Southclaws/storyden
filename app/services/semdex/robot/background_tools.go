package robot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	adkagent "google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry"
	"github.com/Southclaws/storyden/app/services/semdex/robot/tools"
)

const checkBackLaterToolName = "check_back_later"

const (
	minimumCheckBackDelay = time.Second
	maximumCheckBackDelay = 30 * 24 * time.Hour
)

type checkBackLaterInput struct {
	DurationSeconds int    `json:"duration_seconds"`
	Task            string `json:"task"`
}

type deferredInvocationOutput struct {
	Status      string `json:"status"`
	TaskID      string `json:"task_id"`
	Message     string `json:"message"`
	NextAction  string `json:"next_action"`
	ScheduledAt string `json:"scheduled_at,omitempty"`
}

func newBackgroundToolset() (tool.Toolset, error) {
	checkBackLater, err := functiontool.New(
		functiontool.Config{
			Name:        checkBackLaterToolName,
			Description: "Pause a bounded task and resume this Robot session after an approximate delay. Use this when the task is useful later, such as checking whether a thread or page changed. The current turn does not wait; a new turn starts after the delay with the task and authoritative completion of this asynchronous call. Do not use it merely to delay an answer you can provide now.",
			InputSchema: checkBackLaterInputSchema(),
		},
		func(ctx adkagent.Context, input checkBackLaterInput) (deferredInvocationOutput, error) {
			delay := time.Duration(input.DurationSeconds) * time.Second
			if delay < minimumCheckBackDelay || delay > maximumCheckBackDelay {
				return deferredInvocationOutput{}, fmt.Errorf("duration_seconds must be between %d and %d", int(minimumCheckBackDelay/time.Second), int(maximumCheckBackDelay/time.Second))
			}
			task := strings.TrimSpace(input.Task)
			if task == "" {
				return deferredInvocationOutput{}, fmt.Errorf("task is required")
			}
			if ctx.FunctionCallID() == "" {
				return deferredInvocationOutput{}, fmt.Errorf("check_back_later requires a function call ID")
			}

			runAt := time.Now().UTC().Add(delay)
			inputID := InternalInvocationID("scheduled", ctx.SessionID(), ctx.FunctionCallID())
			content := genai.NewContentFromText(fmt.Sprintf(
				"A scheduled check is now due. Task ID: %s\nTask: %s",
				inputID.String(),
				task,
			), genai.RoleUser)

			if err := enqueueInternalInvocation(ctx, InternalInvocation{
				InputID: inputID,
				Content: content,
				Options: RunOptions{
					Mode:      ModeInteractive,
					Source:    SourceScheduled,
					NotBefore: runAt,
				},
			}); err != nil {
				return deferredInvocationOutput{}, err
			}

			return deferredInvocationOutput{
				Status:      "pending",
				TaskID:      inputID.String(),
				ScheduledAt: runAt.Format(time.RFC3339),
				Message:     "The session will resume after the requested delay.",
				NextAction:  "Finish the current response without waiting or calling check_back_later again. The completion will arrive in a later turn.",
			}, nil
		},
	)
	if err != nil {
		return nil, err
	}
	return &tools.Toolset{ToolList: []tool.Tool{checkBackLater}}, nil
}

func checkBackLaterInputSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:     "object",
		Required: []string{"duration_seconds", "task"},
		Properties: map[string]*jsonschema.Schema{
			"duration_seconds": {
				Type:        "integer",
				Description: "Approximate delay in seconds, from 1 second to 30 days.",
				Minimum:     jsonschema.Ptr(float64(minimumCheckBackDelay / time.Second)),
				Maximum:     jsonschema.Ptr(float64(maximumCheckBackDelay / time.Second)),
			},
			"task": {
				Type:        "string",
				Description: "Self-contained instruction describing what to do when the session resumes.",
			},
		},
	}
}

type delegationInput struct {
	Request string `json:"request"`
}

func (s *Agent) buildDelegationToolset(ctx context.Context) (tool.Toolset, error) {
	robots, err := s.db.Robot.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	delegationTools := make([]tool.Tool, 0, len(robots))
	for _, stored := range robots {
		spec, err := s.resolveAgentSpec(ctx, stored.ID.String())
		if err != nil {
			s.logger.Warn("skipping unavailable delegated Robot", "robot_id", stored.ID.String(), "error", err)
			continue
		}
		if _, ok := spec.ModelRef.Get(); !ok {
			continue
		}

		robotRef := spec.RobotRef
		toolName := spec.AgentName
		description := fmt.Sprintf("Start %s as an asynchronous specialist for a bounded task. %s The specialist runs later in this shared session and its result returns in a new turn. This call returns pending immediately; do not claim the specialist has finished until the completion arrives.", spec.DisplayName, strings.TrimSpace(spec.Description))
		delegationTool, err := functiontool.New(
			functiontool.Config{
				Name:        toolName,
				Description: strings.TrimSpace(description),
			},
			func(toolCtx adkagent.Context, input delegationInput) (deferredInvocationOutput, error) {
				request := strings.TrimSpace(input.Request)
				if request == "" {
					return deferredInvocationOutput{}, fmt.Errorf("request is required")
				}
				callID := toolCtx.FunctionCallID()
				if callID == "" {
					return deferredInvocationOutput{}, fmt.Errorf("delegation requires a function call ID")
				}
				inputID := InternalInvocationID("delegation", toolCtx.SessionID(), callID)
				if err := enqueueInternalInvocation(toolCtx, InternalInvocation{
					InputID:  inputID,
					RobotRef: robotRef,
					Content:  genai.NewContentFromText(request, genai.RoleUser),
					Options: RunOptions{
						Mode:   ModeUnattended,
						Source: SourceDelegation,
						Delegation: &agent_registry.DelegationRun{
							CallID:   callID,
							ToolName: toolName,
							Request:  request,
						},
					},
				}); err != nil {
					return deferredInvocationOutput{}, err
				}

				return deferredInvocationOutput{
					Status:     "pending",
					TaskID:     inputID.String(),
					Message:    spec.DisplayName + " accepted the delegated task.",
					NextAction: "Finish the current response without inventing a result. The specialist result will arrive in a later turn.",
				}, nil
			},
		)
		if err != nil {
			return nil, err
		}
		delegationTools = append(delegationTools, delegationTool)
	}

	return &tools.Toolset{ToolList: delegationTools}, nil
}
