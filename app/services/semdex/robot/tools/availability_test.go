package tools

import (
	"context"
	"encoding/json"
	"io"
	"iter"
	"log/slog"
	"testing"

	"github.com/Southclaws/storyden/lib/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/model"
	"google.golang.org/adk/v2/session"
	"google.golang.org/adk/v2/tool/toolconfirmation"
)

func TestDiscoveryAvailabilityTracksActualModelTools(t *testing.T) {
	registry := NewRegistry(slog.New(slog.NewTextHandler(io.Discard, nil)))
	definition := *mcp.GetWebFetchTool()
	definition.Name = "external:fetch"
	selected := &Tool{Definition: &definition, CallableName: "external_fetch"}
	require.NoError(t, registry.Register(selected))
	discovery := newToolDiscoveryTools(registry)
	ctx := &toolTestContext{state: toolTestState{}}
	capture := CaptureCallableTools()
	_, err := capture(ctx, &model.LLMRequest{Tools: map[string]any{"tool_load": nil, "tool_get": nil}})
	require.NoError(t, err)
	output, err := discovery.get(ctx, mcp.ToolToolGetInput{Id: "external:fetch"})
	require.NoError(t, err)
	require.NotNil(t, output.Availability)
	assert.Equal(t, mcp.RobotToolAvailabilityYamlLoadRequired, *output.Availability)
	assert.Equal(t, "external_fetch", output.CallableName)
	search := discovery.search(ctx, mcp.ToolToolSearchInput{Query: ""})
	require.Len(t, search.Tools, 3)
	var found bool
	for _, item := range search.Tools {
		if item.Id == "external:fetch" {
			found = true
			require.NotNil(t, item.Availability)
			assert.Equal(t, mcp.RobotToolAvailabilityYamlLoadRequired, *item.Availability)
			assert.Equal(t, "external_fetch", item.CallableName)
		}
	}
	require.True(t, found)

	_, err = capture(ctx, &model.LLMRequest{Tools: map[string]any{"tool_load": nil, "external_fetch": nil}})
	require.NoError(t, err)
	output, err = discovery.get(ctx, mcp.ToolToolGetInput{Id: "external:fetch"})
	require.NoError(t, err)
	assert.Equal(t, mcp.RobotToolAvailabilityYamlCallable, *output.Availability)

	_, err = capture(ctx, &model.LLMRequest{Tools: map[string]any{}})
	require.NoError(t, err)
	output, err = discovery.get(ctx, mcp.ToolToolGetInput{Id: "external:fetch"})
	require.NoError(t, err)
	assert.Equal(t, mcp.RobotToolAvailabilityYamlBlocked, *output.Availability)
}

func TestAvailabilityReportsBlockersAndToolsetActivation(t *testing.T) {
	ctx := &toolTestContext{state: toolTestState{}}
	_, err := CaptureCallableTools()(ctx, &model.LLMRequest{Tools: map[string]any{"tool_load": nil, "toolset_load": nil}})
	require.NoError(t, err)
	selected := &Tool{Definition: mcp.GetDocumentGetTool()}
	result := toolAvailability(ctx, selected)
	assert.Equal(t, mcp.RobotToolAvailabilityYamlLoadRequired, *result)
	definition := *mcp.GetWebFetchTool()
	definition.RequiresWorkspace = true
	selected.Definition = &definition
	result = toolAvailability(ctx, selected)
	assert.Equal(t, mcp.RobotToolAvailabilityYamlBlocked, *result)
	definition.RequiresWorkspace = false
	require.NoError(t, ctx.state.Set(ToolBlockersStateKey, map[string]string{"web_fetch": "Repair unavailable plugin."}))
	result = toolAvailability(ctx, selected)
	assert.Equal(t, mcp.RobotToolAvailabilityYamlBlocked, *result)
	assert.Nil(t, toolAvailability(context.Background(), selected))
}

type toolTestContext struct {
	agent.Context
	state toolTestState
}

func (c *toolTestContext) State() session.State                 { return c.state }
func (c *toolTestContext) ReadonlyState() session.ReadonlyState { return c.state }
func (c *toolTestContext) Value(any) any                        { return nil }

type toolTestState map[string]any

func (s toolTestState) Get(key string) (any, error) {
	value, ok := s[key]
	if !ok {
		return nil, session.ErrStateKeyNotExist
	}
	return value, nil
}
func (s toolTestState) Set(key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	var stored any
	if err := json.Unmarshal(data, &stored); err != nil {
		return err
	}
	s[key] = stored
	return nil
}
func (s toolTestState) All() iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		for key, value := range s {
			if !yield(key, value) {
				return
			}
		}
	}
}

func (c *toolTestContext) ToolConfirmation() *toolconfirmation.ToolConfirmation { return nil }
