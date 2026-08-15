package robot_ref

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomRobotAgentName(t *testing.T) {
	id, err := NewID("d9n7eodo2dtgv230rpng")
	require.NoError(t, err)

	assert.Equal(t, "robot_d9n7eodo2dtgv230rpng", AgentName(id))
}

func TestCustomRobotIDFromAgentName(t *testing.T) {
	id, ok := IDFromAgentName("robot_d9n7eodo2dtgv230rpng")
	require.True(t, ok)
	assert.Equal(t, "d9n7eodo2dtgv230rpng", id.String())

	_, ok = IDFromAgentName("denbot")
	assert.False(t, ok)
}
