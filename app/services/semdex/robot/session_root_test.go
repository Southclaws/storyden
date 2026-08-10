package robot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionRootRobotRefRoundTrip(t *testing.T) {
	state := sessionInitialState("  custom-robot  ")

	ref, ok := SessionRootRobotRef(state).Get()
	require.True(t, ok)
	assert.Equal(t, "custom-robot", ref)
}

func TestSessionRootRobotRefRejectsMissingOrInvalidState(t *testing.T) {
	assert.False(t, SessionRootRobotRef(nil).Ok())
	assert.False(t, SessionRootRobotRef(map[string]any{sessionRootRobotRefStateKey: 42}).Ok())
	assert.False(t, SessionRootRobotRef(map[string]any{sessionRootRobotRefStateKey: "  "}).Ok())
}
