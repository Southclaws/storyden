package robot_ref

import "strings"

const agentNamePrefix = "robot_"

// AgentName returns the stable ADK agent name for a custom Robot.
func AgentName(id ID) string {
	return agentNamePrefix + id.String()
}

// IDFromAgentName parses a custom Robot ID from its ADK agent name.
func IDFromAgentName(name string) (ID, bool) {
	raw := strings.TrimPrefix(name, agentNamePrefix)
	if raw == name {
		return ID{}, false
	}

	id, err := NewID(raw)
	return id, err == nil
}
