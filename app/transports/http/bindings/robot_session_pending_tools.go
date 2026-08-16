package bindings

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/robot"
)

const pendingClientToolsStateKey = robot.PendingClientToolsStateKey

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
	return robot.ClearPendingClientTools(state)
}

func existingSessState(sess *robot.Session) map[string]any {
	if sess == nil {
		return nil
	}
	return sess.State
}

func readPendingToolIDs(state map[string]any) []string {
	return robot.PendingClientToolIDs(state)
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

func writeChatError(w http.ResponseWriter, err chatError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status)

	_ = json.NewEncoder(w).Encode(map[string]any{
		"error":   err.Code,
		"message": err.Message,
		"status":  err.Status,
	})
}
