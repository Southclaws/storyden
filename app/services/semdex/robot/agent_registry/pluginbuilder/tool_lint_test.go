package pluginbuilder

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	localworkspace "github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider/local"
)

func TestPluginLintFlagsSwallowedEventHandlerFailures(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	writeWorkspaceFile(t, ctx, workspace, "main.go", `package main

import (
	"context"
	"log"
)

type Plugin struct{}
type Event struct{}
type Response struct {
	StatusCode int
}

func (p Plugin) OnThreadReplyCreated(fn func(context.Context, *Event) error) {}

func register(pl Plugin) {
	pl.OnThreadReplyCreated(func(ctx context.Context, event *Event) error {
		resp := &Response{StatusCode: 500}
		if resp.StatusCode >= 400 {
			log.Printf("failed to add reaction")
			return nil
		}
		return nil
	})
}
`)

	agent := &Agent{workspace: workspace}

	result, err := agent.PluginLint(ctx)
	require.NoError(t, err)
	require.False(t, result.Success)
	require.Len(t, result.Issues, 1)
	require.Equal(t, "main.go", result.Issues[0].Path)
	require.Contains(t, result.Issues[0].Message, "return an error instead of nil")
}

func TestPluginLintAllowsSuccessfulEventHandlerNilReturn(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	writeWorkspaceFile(t, ctx, workspace, "main.go", `package main

import (
	"context"
	"fmt"
)

type Plugin struct{}
type Event struct{}
type Response struct {
	StatusCode int
}

func (p Plugin) OnThreadReplyCreated(fn func(context.Context, *Event) error) {}

func register(pl Plugin) {
	pl.OnThreadReplyCreated(func(ctx context.Context, event *Event) error {
		resp := &Response{StatusCode: 500}
		if resp.StatusCode >= 400 {
			return fmt.Errorf("failed to add reaction: %d", resp.StatusCode)
		}
		return nil
	})
}
`)

	agent := &Agent{workspace: workspace}

	result, err := agent.PluginLint(ctx)
	require.NoError(t, err)
	require.True(t, result.Success, result.Format())
}

func TestPluginLintFlagsUnsupportedStorydenEventRPCMethod(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	writeWorkspaceFile(t, ctx, workspace, "main.go", `package main

import "context"

type Plugin struct{}
type EventActivityCreated struct{}

func register(pl Plugin) {
	pl.HandleEventRPC("EventActivityCreated", func(ctx context.Context, ev *EventActivityCreated) error {
		return nil
	})
}
`)

	agent := &Agent{workspace: workspace}

	result, err := agent.PluginLint(ctx)
	require.NoError(t, err)
	require.False(t, result.Success)
	require.Len(t, result.Issues, 1)
	require.Contains(t, result.Issues[0].Message, "no HandleEventRPC method")
	require.Contains(t, result.Issues[0].Message, "pl.OnActivityCreated")
}

func TestPluginLintFlagsRobotChatSSEForPluginRobotRuns(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	writeWorkspaceFile(t, ctx, workspace, "main.go", `package main

import "context"

type Client struct{}

func (c Client) RobotChatSSEWithResponse(context.Context, string) error { return nil }

func run(ctx context.Context, client Client) error {
	return client.RobotChatSSEWithResponse(ctx, "robot")
}
`)

	agent := &Agent{workspace: workspace}

	result, err := agent.PluginLint(ctx)
	require.NoError(t, err)
	require.False(t, result.Success)
	require.Len(t, result.Issues, 1)
	require.Contains(t, result.Issues[0].Message, "UI streaming endpoint")
	require.Contains(t, result.Issues[0].Message, "pl.RunRobot")
}

func TestPluginLintFlagsBareErrReturnedForHTTPStatusFailure(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	writeWorkspaceFile(t, ctx, workspace, "main.go", `package main

import "context"

type Plugin struct{}
type Event struct{}
type Response struct {
	StatusCode int
}

func (p Plugin) OnThreadReplyCreated(fn func(context.Context, *Event) error) {}

func register(pl Plugin) {
	pl.OnThreadReplyCreated(func(ctx context.Context, event *Event) error {
		resp, err := addReaction(ctx)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 400 {
			return err
		}
		return nil
	})
}

func addReaction(ctx context.Context) (*Response, error) {
	return &Response{StatusCode: 500}, nil
}
`)

	agent := &Agent{workspace: workspace}

	result, err := agent.PluginLint(ctx)
	require.NoError(t, err)
	require.False(t, result.Success)
	require.Len(t, result.Issues, 1)
	require.Contains(t, result.Issues[0].Message, "return a descriptive error")
}
