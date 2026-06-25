package sse

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/google/uuid"
	adkmodel "google.golang.org/adk/model"
	adksession "google.golang.org/adk/session"
	"google.golang.org/genai"
)

const (
	pendingClientToolsStateKey      = "pending_client_tools"
	pendingClientToolRobotsStateKey = "pending_client_tool_robots"
)

type pendingClientTools struct {
	IDs    []string
	Robots map[string]string
}

type robotSwitchRecovery struct {
	ToolCallID string
	RobotID    string
}

type chatError struct {
	Code    string
	Message string
	Status  int
}

type pendingClientToolReconciliation struct {
	Pending             pendingClientTools
	Provided            map[string]struct{}
	OwnerRobotID        opt.Optional[string]
	StaleRobotSwitch    opt.Optional[robotSwitchRecovery]
	BlockingInteraction opt.Optional[chatError]
}

func readPendingClientTools(state map[string]any) pendingClientTools {
	return pendingClientTools{
		IDs:    readPendingToolIDs(state),
		Robots: readPendingToolRobots(state),
	}
}

func reconcilePendingClientTools(messages []chatMessage, pending pendingClientTools) pendingClientToolReconciliation {
	reconciliation := pendingClientToolReconciliation{
		Pending:  pending,
		Provided: map[string]struct{}{},
	}

	if len(pending.IDs) == 0 {
		return reconciliation
	}

	reconciliation.Provided = getProvidedPendingToolIDs(messages, pending.IDs)
	if len(reconciliation.Provided) == 0 {
		if targetRobotID, ok := getPendingRobotSwitchInputTargetID(messages, pending.IDs).Get(); ok {
			reconciliation.StaleRobotSwitch = opt.New(robotSwitchRecovery{
				ToolCallID: pending.IDs[0],
				RobotID:    targetRobotID,
			})
			return reconciliation
		}

		reconciliation.BlockingInteraction = opt.New(chatError{
			Code:    "pending_tool_interaction",
			Message: "pending tool interaction must be resolved before continuing",
			Status:  http.StatusConflict,
		})
		return reconciliation
	}

	if len(reconciliation.Provided) != len(pending.IDs) {
		reconciliation.BlockingInteraction = opt.New(chatError{
			Code:    "pending_tool_interaction",
			Message: "all pending tool interactions from the assistant turn must be resolved together",
			Status:  http.StatusConflict,
		})
		return reconciliation
	}

	reconciliation.OwnerRobotID = getPendingToolRobotID(pending.IDs, pending.Robots)

	return reconciliation
}

func clearPendingClientTools(state map[string]any) map[string]any {
	if state == nil {
		state = make(map[string]any)
	}

	delete(state, pendingClientToolsStateKey)
	delete(state, pendingClientToolRobotsStateKey)

	return state
}

func existingSessState(sess *robot.Session) map[string]any {
	if sess == nil {
		return nil
	}
	return sess.State
}

func persistClientToolResult(
	ctx context.Context,
	sessionRepo *robot_session.Repository,
	sessionID robot.SessionID,
	accountID account.AccountID,
	content *genai.Content,
) error {
	event := adksession.NewEvent(uuid.NewString())
	event.Author = "user"
	event.LLMResponse = adkmodel.LLMResponse{Content: content}

	eventData, err := marshalMap(event)
	if err != nil {
		return err
	}

	return sessionRepo.AppendMessage(
		ctx,
		sessionID,
		event.InvocationID,
		opt.New(accountID),
		opt.NewEmpty[robot.Actor](),
		eventData,
	)
}

func robotSwitchToolResultContent(toolCallID string, robotID string) *genai.Content {
	return &genai.Content{
		Role: genai.RoleUser,
		Parts: []*genai.Part{
			{
				FunctionResponse: &genai.FunctionResponse{
					ID:   toolCallID,
					Name: "robot_switch",
					Response: map[string]any{
						"success":  true,
						"robot_id": robotID,
					},
				},
			},
		},
	}
}

func finishEmptyStream(w http.ResponseWriter) error {
	emitter, err := newStreamEmitter(w)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return err
	}
	defer emitter.Done()

	if err := emitter.Headers(); err != nil {
		return err
	}

	startPart := openapi.StreamPart{}
	if err := startPart.FromStartPart(openapi.StartPart{MessageId: uuid.NewString()}); err != nil {
		return err
	}
	if err := emitter.Send(startPart); err != nil {
		return err
	}

	finishPart := openapi.StreamPart{}
	if err := finishPart.FromFinishMessagePart(openapi.FinishMessagePart{}); err != nil {
		return err
	}
	return emitter.Send(finishPart)
}

func marshalMap(v any) (map[string]any, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}

	return out, nil
}

func getRobotSwitchTargetID(messages []chatMessage, pendingToolIDs []string) opt.Optional[string] {
	if len(messages) == 0 {
		return opt.NewEmpty[string]()
	}

	pending := make(map[string]struct{}, len(pendingToolIDs))
	for _, id := range pendingToolIDs {
		pending[id] = struct{}{}
	}

	lastMessage := messages[len(messages)-1]
	if !strings.EqualFold(lastMessage.Role, "assistant") {
		return opt.NewEmpty[string]()
	}

	for _, part := range lastMessage.Parts {
		toolName := part.ToolName
		if toolName == "" && strings.HasPrefix(part.Type, "tool-") {
			toolName = strings.TrimPrefix(part.Type, "tool-")
		}
		if part.State != "output-available" || toolName != "robot_switch" {
			continue
		}
		if len(pending) > 0 {
			if _, ok := pending[part.ToolCallId]; !ok {
				continue
			}
		}

		var output struct {
			RobotID string `json:"robot_id"`
		}
		if err := json.Unmarshal(part.Output, &output); err != nil || output.RobotID == "" {
			continue
		}
		return opt.New(output.RobotID)
	}

	return opt.NewEmpty[string]()
}

func readPendingToolRobots(state map[string]any) map[string]string {
	result := map[string]string{}
	existing, ok := state[pendingClientToolRobotsStateKey]
	if !ok {
		return result
	}

	switch v := existing.(type) {
	case map[string]string:
		for id, robotID := range v {
			result[id] = robotID
		}
	case map[string]any:
		for id, robotID := range v {
			if s, ok := robotID.(string); ok {
				result[id] = s
			}
		}
	}

	return result
}

func readPendingToolIDs(state map[string]any) []string {
	if state == nil {
		return nil
	}

	existing, ok := state[pendingClientToolsStateKey]
	if !ok {
		return nil
	}

	var result []string
	switch v := existing.(type) {
	case []string:
		result = append(result, v...)
	case []any:
		for _, id := range v {
			if s, ok := id.(string); ok {
				result = append(result, s)
			}
		}
	}

	return result
}

func getPendingToolRobotID(pendingToolIDs []string, pendingToolRobots map[string]string) opt.Optional[string] {
	for _, pendingID := range pendingToolIDs {
		robotID, ok := pendingToolRobots[pendingID]
		if !ok || robotID == "" {
			continue
		}
		return opt.New(robotID)
	}

	return opt.NewEmpty[string]()
}

func getPendingRobotSwitchInputTargetID(messages []chatMessage, pendingToolIDs []string) opt.Optional[string] {
	if len(messages) == 0 || len(pendingToolIDs) != 1 {
		return opt.NewEmpty[string]()
	}

	pendingToolID := pendingToolIDs[0]
	for i := len(messages) - 1; i >= 0; i-- {
		for _, part := range messages[i].Parts {
			if part.ToolCallId != pendingToolID {
				continue
			}

			toolName := part.ToolName
			if toolName == "" && strings.HasPrefix(part.Type, "tool-") {
				toolName = strings.TrimPrefix(part.Type, "tool-")
			}
			if toolName != "robot_switch" {
				continue
			}

			var input struct {
				RobotID string `json:"robot_id"`
			}
			if len(part.Input) == 0 {
				continue
			}
			if err := json.Unmarshal(part.Input, &input); err != nil {
				continue
			}
			if input.RobotID == "" {
				continue
			}

			return opt.New(input.RobotID)
		}
	}

	return opt.NewEmpty[string]()
}

func getProvidedPendingToolIDs(messages []chatMessage, pendingToolIDs []string) map[string]struct{} {
	provided := map[string]struct{}{}
	if len(messages) == 0 || len(pendingToolIDs) == 0 {
		return provided
	}

	pending := make(map[string]struct{}, len(pendingToolIDs))
	for _, id := range pendingToolIDs {
		pending[id] = struct{}{}
	}

	lastMessage := messages[len(messages)-1]
	if !strings.EqualFold(lastMessage.Role, "assistant") {
		return provided
	}

	for _, part := range lastMessage.Parts {
		if part.ToolCallId == "" {
			continue
		}
		if part.State != "output-available" && part.State != "approval-responded" {
			continue
		}
		if _, ok := pending[part.ToolCallId]; ok {
			provided[part.ToolCallId] = struct{}{}
		}
	}

	return provided
}

func storePendingToolID(ctx context.Context, sessionRepo *robot_session.Repository, sessionID robot.SessionID, toolCallID string, robotID opt.Optional[string]) error {
	sess, _, err := sessionRepo.Get(ctx, sessionID, robot.NewMessageCursorParams(opt.NewEmpty[robot.MessageID](), 1))
	if err != nil {
		return err
	}

	state := sess.State
	if state == nil {
		state = make(map[string]any)
	}

	pendingIDs := readPendingToolIDs(state)
	pendingIDs = append(pendingIDs, toolCallID)
	state[pendingClientToolsStateKey] = pendingIDs
	pendingRobots := readPendingToolRobots(state)
	if id, ok := robotID.Get(); ok {
		pendingRobots[toolCallID] = id
		state["current_robot_id"] = id
	}
	state[pendingClientToolRobotsStateKey] = pendingRobots

	return sessionRepo.UpdateState(ctx, sessionID, state)
}

func writeChatError(w http.ResponseWriter, err chatError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status)

	_ = json.NewEncoder(w).Encode(map[string]any{
		"error":   err.Code,
		"message": err.Message,
		"status":  err.Status,
	})
}
