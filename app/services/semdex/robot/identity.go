package robot

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/agent/llmagent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/robot/robot_ref"
	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry"
	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry/robotbuilder"
	ent_robot_session_message "github.com/Southclaws/storyden/internal/ent/robotsessionmessage"
	"github.com/Southclaws/storyden/lib/mcp"
)

type robotIdentity struct {
	ID               opt.Optional[xid.ID]
	Name             string
	Description      string
	Capabilities     []string
	UnavailableTools []string
}

type robotParticipant struct {
	Name   string
	Active bool
}

type robotIdentityContext struct {
	Current      robotIdentity
	Participants []robotParticipant
}

type robotNameResolver func(context.Context, string) string

func robotCapabilityNames(toolNames []string) []string {
	if len(toolNames) == 0 {
		return nil
	}

	capabilities := append([]string(nil), toolNames...)
	sort.Strings(capabilities)
	return capabilities
}

func robotIdentityInstruction(ctx robotIdentityContext) string {
	var b strings.Builder

	b.WriteString("## Current Robot\n\n")
	b.WriteString("You are currently running as this Robot:\n\n")
	b.WriteString(fmt.Sprintf("Name: %s\n", ctx.Current.Name))
	if strings.TrimSpace(ctx.Current.Description) != "" {
		b.WriteString(fmt.Sprintf("Description: %s\n", strings.TrimSpace(ctx.Current.Description)))
	}
	if len(ctx.Current.Capabilities) > 0 {
		b.WriteString("Capabilities available to this Robot:\n")
		for _, capability := range ctx.Current.Capabilities {
			b.WriteString(fmt.Sprintf("- %s\n", capability))
		}
	}
	if len(ctx.Current.UnavailableTools) > 0 {
		b.WriteString("\nConfigured tools currently unavailable:\n")
		for _, tool := range ctx.Current.UnavailableTools {
			b.WriteString(fmt.Sprintf("- %s\n", tool))
		}
		b.WriteString("\nIf these tools are relevant to the user's request, briefly explain that this Robot's toolset has changed and continue with the available capabilities.\n")
	}
	b.WriteString("\nThis identity is authoritative for the current turn. Do not infer your current identity from earlier assistant messages.\n")

	if len(ctx.Participants) > 1 {
		b.WriteString("\n## Session Robot Context\n\n")
		b.WriteString("This session has involved multiple Robots.\n\n")
		b.WriteString("Current Robot:\n")
		b.WriteString(fmt.Sprintf("- %s - active for this turn\n\n", ctx.Current.Name))

		var others []robotParticipant
		for _, participant := range ctx.Participants {
			if participant.Active || participant.Name == ctx.Current.Name {
				continue
			}
			others = append(others, participant)
		}

		if len(others) > 0 {
			b.WriteString("Other Robots seen in this session:\n")
			for _, participant := range others {
				b.WriteString(fmt.Sprintf("- %s - previously active\n", participant.Name))
			}
			b.WriteString("\n")
		}

		b.WriteString("Robot-switch markers in the conversation history indicate where the active Robot changed.\n")
		b.WriteString("Messages before a robot-switch marker were produced by the Robot active at that time, not necessarily the current Robot.\n")
		b.WriteString("Use prior messages as conversation context, but do not adopt another Robot's identity, role, tools, permissions, or responsibilities unless the runtime context says that Robot is active now.\n")
	}

	return b.String()
}

func (s *Agent) buildRobotIdentityContext(ctx context.Context, sessionID string, current robotIdentity) robotIdentityContext {
	participants := []robotParticipant{{Name: current.Name, Active: true}}

	sessionXID, err := xid.FromString(sessionID)
	if err != nil {
		return robotIdentityContext{Current: current, Participants: participants}
	}

	messages, err := s.db.RobotSessionMessage.Query().
		Where(ent_robot_session_message.SessionIDEQ(sessionXID)).
		WithRobot().
		All(ctx)
	if err != nil {
		s.logger.Debug("failed to load robot session participants",
			slog.String("session_id", sessionID),
			slog.String("error", err.Error()))
		return robotIdentityContext{Current: current, Participants: participants}
	}

	seen := map[string]bool{current.Name: true}
	resolveName := s.robotNameResolver()
	for _, message := range messages {
		if message.RobotID == nil {
			if message.BuiltinRobot != nil {
				name := resolveName(ctx, *message.BuiltinRobot)
				if !seen[name] {
					seen[name] = true
					participants = append(participants, robotParticipant{Name: name})
				}
				continue
			}
			if author, _ := message.EventData["Author"].(string); author == robotbuilder.AgentName && !seen[robotbuilder.DisplayName] {
				seen[robotbuilder.DisplayName] = true
				participants = append(participants, robotParticipant{Name: robotbuilder.DisplayName})
			}
			continue
		}

		robot := message.Edges.Robot
		if robot == nil {
			continue
		}
		if seen[robot.Name] {
			continue
		}

		seen[robot.Name] = true
		participants = append(participants, robotParticipant{Name: robot.Name})
	}

	sort.SliceStable(participants, func(i, j int) bool {
		if participants[i].Active != participants[j].Active {
			return participants[i].Active
		}
		return participants[i].Name < participants[j].Name
	})

	return robotIdentityContext{Current: current, Participants: participants}
}

func (s *Agent) robotNameResolver() robotNameResolver {
	cache := map[string]string{
		"":                            robotbuilder.DisplayName,
		agent_registry.RobotBuilderID: robotbuilder.DisplayName,
	}

	return func(ctx context.Context, id string) string {
		if name, ok := cache[id]; ok {
			return name
		}
		if def, ok := s.agents.Get(id); ok {
			cache[id] = def.Name
			return def.Name
		}

		robotID, err := robot_ref.NewID(id)
		if err != nil {
			cache[id] = id
			return id
		}

		robot, err := s.db.Robot.Get(ctx, xid.ID(robotID))
		if err != nil {
			cache[id] = id
			return id
		}

		cache[id] = robot.Name
		return robot.Name
	}
}

func projectRobotSwitchesBeforeModel(logger *slog.Logger, resolve robotNameResolver) llmagent.BeforeModelCallback {
	return func(ctx agent.Context, req *model.LLMRequest) (*model.LLMResponse, error) {
		rewritten := projectRobotSwitches(ctx, req, resolve)
		if rewritten > 0 {
			logger.Info("projected robot switch history",
				slog.String("agent", ctx.AgentName()),
				slog.String("invocation", ctx.InvocationID()),
				slog.String("session", ctx.SessionID()),
				slog.Int("count", rewritten),
			)
		}
		return nil, nil
	}
}

func projectRobotSwitches(ctx context.Context, req *model.LLMRequest, resolve robotNameResolver) int {
	if req == nil {
		return 0
	}

	completed := completedRobotSwitches(req)

	var count int
	projected := make([]*genai.Content, 0, len(req.Contents))

	for _, content := range req.Contents {
		if content == nil {
			continue
		}

		parts := make([]*genai.Part, 0, len(content.Parts))
		removedSwitchCall := false
		for _, part := range content.Parts {
			if part == nil {
				continue
			}

			if call := part.FunctionCall; call != nil && call.Name == mcp.GetRobotSwitchTool().Name {
				if _, ok := completed[call.ID]; ok {
					count++
					removedSwitchCall = true
					continue
				}
			}

			if response := part.FunctionResponse; response != nil && response.Name == mcp.GetRobotSwitchTool().Name {
				targetID, ok := completed[response.ID]
				if ok {
					parts = append(parts, genai.NewPartFromText(robotSwitchMarker(ctx, targetID, resolve)))
					continue
				}
			}

			if targetID, ok := robotSwitchTargetFromContextText(part.Text); ok {
				part = genai.NewPartFromText(robotSwitchMarker(ctx, targetID, resolve))
				count++
			}

			parts = append(parts, part)
		}

		if len(parts) == 0 && removedSwitchCall {
			parts = append(parts, genai.NewPartFromText("Robot switch requested."))
		}

		if len(parts) == 0 {
			continue
		}

		projected = append(projected, &genai.Content{
			Role:  content.Role,
			Parts: parts,
		})
	}

	req.Contents = projected

	return count
}

func robotSwitchTargetFromContextText(text string) (string, bool) {
	const marker = "called tool `robot_switch` with parameters:"

	idx := strings.Index(text, marker)
	if idx < 0 {
		return "", false
	}

	raw := strings.TrimSpace(text[idx+len(marker):])
	if raw == "" {
		return "", false
	}

	var args map[string]any
	if err := json.Unmarshal([]byte(raw), &args); err != nil {
		return "", false
	}

	robotID, _ := args["robot_id"].(string)
	if robotID == "" {
		return "", false
	}

	return robotID, true
}

func completedRobotSwitches(req *model.LLMRequest) map[string]string {
	completed := map[string]string{}

	for _, content := range req.Contents {
		if content == nil {
			continue
		}
		for _, part := range content.Parts {
			if part == nil || part.FunctionResponse == nil {
				continue
			}

			response := part.FunctionResponse
			if response.Name != mcp.GetRobotSwitchTool().Name {
				continue
			}
			if approved, _ := response.Response["success"].(bool); !approved {
				continue
			}
			robotID, _ := response.Response["robot_id"].(string)
			if robotID == "" {
				continue
			}

			completed[response.ID] = robotID
		}
	}

	return completed
}

func robotSwitchMarker(ctx context.Context, targetRobotID string, resolve robotNameResolver) string {
	targetName := targetRobotID
	if resolve != nil {
		targetName = resolve(ctx, targetRobotID)
	}

	return fmt.Sprintf(`ROBOT SWITCH

The active Robot changed at this point.

Current Robot after this switch: %s

Messages after this marker were handled by %s until the next robot-switch marker. Messages before this marker may have been produced by a different Robot with different instructions, tools, permissions, and responsibilities.`, targetName, targetName)
}
