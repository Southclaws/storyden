package pluginbuilder

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	localworkspace "github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider/local"
)

func TestStorydenSDKEventsIncludesIDUsageHints(t *testing.T) {
	ctx := context.Background()
	agent := newStorydenSDKTestAgent(t, ctx)

	result, err := agent.StorydenSDKEvents(ctx, StorydenSDKEventsInput{
		Query:     "reply",
		MaxEvents: 20,
	})
	require.NoError(t, err)

	for _, event := range result.Events {
		if event.Event != "EventThreadReplyCreated" {
			continue
		}
		require.Equal(t, "OnThreadReplyCreated", event.HandlerMethod)
		for _, usage := range event.FieldUsages {
			if usage.Field == "ReplyID" {
				require.Equal(t, "event.ReplyID.String()", usage.Expression)
				return
			}
		}
		require.Fail(t, "missing ReplyID field usage", "%#v", event.FieldUsages)
	}
	require.Fail(t, "missing EventThreadReplyCreated", "%#v", result.Events)
}

func TestStorydenSDKEventsExplainActivityBoundary(t *testing.T) {
	ctx := context.Background()
	agent := newStorydenSDKTestAgent(t, ctx)

	result, err := agent.StorydenSDKEvents(ctx, StorydenSDKEventsInput{
		Query:     "activity",
		MaxEvents: 20,
	})
	require.NoError(t, err)

	found := false
	for _, event := range result.Events {
		if event.Event == "EventActivityCreated" {
			require.Equal(t, "OnActivityCreated", event.HandlerMethod)
			found = true
		}
	}
	require.True(t, found, "missing EventActivityCreated in %#v", result.Events)
	requireSDKHint(t, result.Hints, "not Discord gateway messages")
	requireSDKHint(t, result.Hints, "HandleEventRPC")
}

func TestStorydenSDKSearchNormalisesReactionQuery(t *testing.T) {
	ctx := context.Background()
	agent := newStorydenSDKTestAgent(t, ctx)

	result, err := agent.StorydenSDKSearch(ctx, StorydenSDKSearchInput{
		Area:       "http_api",
		Query:      "reaction",
		MaxResults: 50,
	})
	require.NoError(t, err)
	requireSymbol(t, result.Symbols, "PostReactAddJSONRequestBody", "type")
	requireSymbol(t, result.Symbols, "PostReactAddWithResponse", "method")
	requireNoSymbol(t, result.Symbols, "ClientInterface")
}

func TestStorydenSDKSearchFindsPluginRuntimeMethods(t *testing.T) {
	ctx := context.Background()
	agent := newStorydenSDKTestAgent(t, ctx)

	result, err := agent.StorydenSDKSearch(ctx, StorydenSDKSearchInput{
		Area:       "plugin",
		Query:      "buildapiclient",
		MaxResults: 20,
	})
	require.NoError(t, err)
	requireSymbol(t, result.Symbols, "BuildAPIClient", "method")
	require.Contains(t, result.Hints[0].Message, "plugin_storyden_sdk_events")
	requireSDKHint(t, result.Hints, "manifest.yaml must include access")
	requireSDKHint(t, result.Hints, "stable bot account handle")
}

func TestStorydenSDKSearchFindsRobotRunWithoutHostGo(t *testing.T) {
	ctx := context.Background()
	agent := newStorydenSDKTestAgent(t, ctx)

	result, err := agent.StorydenSDKSearch(ctx, StorydenSDKSearchInput{
		Area:       "plugin",
		Query:      "robot_run",
		MaxResults: 20,
	})
	require.NoError(t, err)
	requireSymbol(t, result.Symbols, "RunRobot", "method")
	requireSDKHint(t, result.Hints, "USE_ROBOTS")
}

func TestStorydenSDKSearchFindsRobotRunFromRPCAndAllAreas(t *testing.T) {
	ctx := context.Background()
	agent := newStorydenSDKTestAgent(t, ctx)

	for _, area := range []string{"rpc", "all"} {
		result, err := agent.StorydenSDKSearch(ctx, StorydenSDKSearchInput{
			Area:       area,
			Query:      "robot_run",
			MaxResults: 20,
		})
		require.NoError(t, err)
		requireSymbol(t, result.Symbols, "RunRobot", "method")
		requireSDKHint(t, result.Hints, "USE_ROBOTS")
	}
}

func TestStorydenSDKSearchFindsRobotRunFromNaturalQueries(t *testing.T) {
	ctx := context.Background()
	agent := newStorydenSDKTestAgent(t, ctx)

	cases := []StorydenSDKSearchInput{
		{Area: "rpc", Query: "run", MaxResults: 20},
		{Area: "http_api", Query: "robot", MaxResults: 20},
		{Area: "operations", Query: "robot run", MaxResults: 20},
	}

	for _, tc := range cases {
		result, err := agent.StorydenSDKSearch(ctx, tc)
		require.NoError(t, err)
		requireSymbol(t, result.Symbols, "RunRobot", "method")
		requireSDKHint(t, result.Hints, "USE_ROBOTS")
		requireSDKHint(t, result.Hints, "Do not use generated HTTP RobotChatSSE")
	}
}

func TestStorydenSDKSearchHandlesNaturalMultiTermQueries(t *testing.T) {
	ctx := context.Background()
	agent := newStorydenSDKTestAgent(t, ctx)

	result, err := agent.StorydenSDKSearch(ctx, StorydenSDKSearchInput{
		Area:       "http_api",
		Query:      "react reply thread",
		MaxResults: 80,
	})
	require.NoError(t, err)
	requireSymbol(t, result.Symbols, "PostReactAddWithResponse", "method")
	requireSymbol(t, result.Symbols, "ReplyCreateWithResponse", "method")
}

func newStorydenSDKTestAgent(t *testing.T, ctx context.Context) *Agent {
	t.Helper()

	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	repoRoot := storydenRepoRoot(t)
	writeWorkspaceFile(t, ctx, workspace, "go.mod", "module example.com/plugin\n\ngo 1.24\n\nrequire github.com/Southclaws/storyden v0.0.0\n\nreplace github.com/Southclaws/storyden => "+repoRoot+"\n")
	writeWorkspaceFile(t, ctx, workspace, "main.go", "package main\n")

	return &Agent{workspace: workspace}
}

func requireSDKHint(t *testing.T, hints []StorydenSDKHint, fragment string) {
	t.Helper()

	for _, hint := range hints {
		if strings.Contains(hint.Message, fragment) {
			return
		}
	}
	require.Fail(t, "missing SDK hint", "fragment %q in %#v", fragment, hints)
}

func storydenRepoRoot(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	require.NoError(t, err)
	root, err := filepath.Abs(filepath.Join(wd, "../../../../../.."))
	require.NoError(t, err)
	return filepath.ToSlash(root)
}

func requireNoSymbol(t *testing.T, symbols []GoSymbolSummary, name string) {
	t.Helper()
	for _, symbol := range symbols {
		require.NotEqual(t, name, symbol.Name, "unexpected noisy symbol in result: %#v", symbol)
	}
}
