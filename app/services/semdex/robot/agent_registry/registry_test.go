package agent_registry

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistryRegisterValidatesRequiredDefinitionFields(t *testing.T) {
	registry := New(slog.Default())

	err := registry.Register(Definition{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ID is required")

	err = registry.Register(Definition{ID: "missing_name"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")

	err = registry.Register(Definition{ID: "missing_app", Name: "Missing App"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "app name is required")

	err = registry.Register(Definition{ID: "missing_agent", Name: "Missing Agent", AppName: "storyden"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "agent name is required")

	err = registry.Register(Definition{ID: "missing_instruction", Name: "Missing Instruction", AppName: "storyden", AgentName: "agent"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "instruction is required")
}

func TestRegistryRegisterRejectsDuplicates(t *testing.T) {
	registry := New(slog.Default())
	def := Definition{
		ID:          "robot_builder",
		Name:        "Robot Builder",
		AppName:     "storyden",
		AgentName:   "storyden",
		Instruction: "Build robots.",
	}

	require.NoError(t, registry.Register(def))
	err := registry.Register(def)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")
}
