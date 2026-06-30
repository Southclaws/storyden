package robot

import (
	"context"

	"github.com/Southclaws/opt"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/rs/xid"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool"
	"google.golang.org/adk/v2/tool/functiontool"

	robotresource "github.com/Southclaws/storyden/app/resources/robot"
)

const unattendedFinishToolName = "robot_run_finish"

func UnattendedFinishToolName() string {
	return unattendedFinishToolName
}

const unattendedInstruction = `## Unattended Invocation

You are running unattended. No user is available for questions, elicitations, or confirmations. Complete the task using the tools available to this Robot. If required information is missing, a required user interaction cannot be completed, a tool is unavailable, or an action is blocked by permissions or policy, stop with a blocked or failed status.

When the unattended run is complete or cannot continue, your final action must be calling the robot_run_finish tool exactly once. Do not ask the user to respond in chat. Do not provide a normal final text answer instead of calling robot_run_finish.`

type unattendedFinishInput struct {
	Status    string                     `json:"status"`
	Summary   string                     `json:"summary"`
	Attention *unattendedFinishAttention `json:"attention,omitempty"`
}

type unattendedFinishAttention struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

func newUnattendedFinishTool() (tool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        unattendedFinishToolName,
			Description: "Finish an unattended Robot invocation with its structured final status. This must be the final action in every unattended run.",
			InputSchema: unattendedFinishInputSchema(),
		},
		func(ctx agent.Context, args unattendedFinishInput) (unattendedFinishInput, error) {
			return args, nil
		},
	)
}

func unattendedFinishInputSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "object",
		Title:       "RobotRunOutput",
		Description: "Final structured status for an unattended Storyden Robot invocation.",
		Required:    []string{"status", "summary"},
		Properties: map[string]*jsonschema.Schema{
			"status": {
				Type:        "string",
				Description: "Overall invocation status.",
				Enum:        []any{"completed", "blocked", "failed"},
			},
			"summary": {
				Type:        "string",
				Description: "Concise user-facing summary of what happened.",
			},
			"attention": {
				Type:        "object",
				Description: "Details when the run needs human attention.",
				Required:    []string{"reason", "message"},
				Properties: map[string]*jsonschema.Schema{
					"reason": {
						Type:        "string",
						Description: "Reason the run needs attention.",
						Enum:        []any{"missing_input", "needs_approval", "policy_blocked", "tool_unavailable", "error"},
					},
					"message": {
						Type:        "string",
						Description: "Actionable explanation for the user.",
					},
				},
			},
		},
	}
}

func (s *Agent) markSessionUnattended(ctx context.Context, sessionID string, options RunOptions) error {
	id, err := xid.FromString(sessionID)
	if err != nil {
		return err
	}

	sess, _, err := s.sessionRepo.Get(ctx, robotresource.SessionID(id), robotresource.NewMessageCursorParams(opt.NewEmpty[robotresource.MessageID](), 1))
	if err != nil {
		return err
	}

	state := sess.State
	if state == nil {
		state = make(map[string]any)
	}
	state["runtime_mode"] = string(options.Mode)
	state["invocation_source"] = string(options.Source)

	return s.sessionRepo.UpdateState(ctx, robotresource.SessionID(id), state)
}
