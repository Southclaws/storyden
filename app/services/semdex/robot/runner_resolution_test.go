package robot

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry"
	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry/denbot"
	"github.com/Southclaws/storyden/app/services/semdex/robot/toolsets/system_documents"
	"github.com/Southclaws/storyden/app/services/semdex/robot/toolsets/system_memory"
	"github.com/Southclaws/storyden/app/services/semdex/robot/toolsets/system_memory_management"
)

func TestResolveAgentSpecResolvesDenbot(t *testing.T) {
	registry := agent_registry.New(slog.Default())
	require.NoError(t, denbot.Register(registry))

	agent := &Agent{
		logger: slog.Default(),
		agents: registry,
	}

	spec, err := agent.resolveAgentSpec(context.Background(), "denbot")
	require.NoError(t, err)

	assert.Equal(t, "denbot", spec.RobotRef)
	assert.Equal(t, "denbot", spec.AppName)
	assert.Equal(t, "denbot", spec.AgentName)
	assert.Equal(t, "Denbot", spec.DisplayName)
	assert.False(t, spec.DatabaseRobotID.Ok())
	assert.Contains(t, spec.ToolsetRefs, system_documents.ID)
	assert.Contains(t, spec.ToolsetRefs, system_memory.ID)
	assert.NotContains(t, spec.ToolsetRefs, system_memory_management.ID)
}

func TestResolveAgentSpecRejectsUnknownBuiltin(t *testing.T) {
	agent := &Agent{
		logger: slog.Default(),
		agents: agent_registry.New(slog.Default()),
	}

	_, err := agent.resolveAgentSpec(context.Background(), "missing_builtin")
	require.Error(t, err)
	assert.Contains(t, err.Error(), `unknown robot "missing_builtin"`)
}
