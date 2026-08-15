package robot

const PendingClientToolsStateKey = "pending_client_tools"

func PendingClientToolIDs(state map[string]any) []string {
	if state == nil {
		return nil
	}
	value, ok := state[PendingClientToolsStateKey]
	if !ok {
		return nil
	}

	var result []string
	switch ids := value.(type) {
	case []string:
		result = append(result, ids...)
	case []any:
		for _, id := range ids {
			if value, ok := id.(string); ok {
				result = append(result, value)
			}
		}
	}
	return result
}

func ClearPendingClientTools(state map[string]any) map[string]any {
	delete(state, PendingClientToolsStateKey)
	return state
}
