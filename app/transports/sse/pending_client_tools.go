package sse

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
)

const pendingClientToolsStateKey = "pending_client_tools"

type pendingClientTools struct {
	IDs []string
}

type chatError struct {
	Code    string
	Message string
	Status  int
}

type pendingClientToolReconciliation struct {
	Pending             pendingClientTools
	Provided            map[string]struct{}
	BlockingInteraction opt.Optional[chatError]
}

func readPendingClientTools(state map[string]any) pendingClientTools {
	return pendingClientTools{IDs: readPendingToolIDs(state)}
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
	}

	return reconciliation
}

func clearPendingClientTools(state map[string]any) map[string]any {
	if state == nil {
		state = make(map[string]any)
	}
	delete(state, pendingClientToolsStateKey)
	return state
}

func existingSessState(sess *robot.Session) map[string]any {
	if sess == nil {
		return nil
	}
	return sess.State
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
	switch value := existing.(type) {
	case []string:
		result = append(result, value...)
	case []any:
		for _, id := range value {
			if value, ok := id.(string); ok {
				result = append(result, value)
			}
		}
	}
	return result
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

func storePendingToolID(ctx context.Context, sessionRepo *robot_session.Repository, sessionID robot.SessionID, toolCallID string) error {
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
