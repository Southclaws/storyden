package pluginbuilder

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider"
	localworkspace "github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider/local"
)

func TestGoDiscoveryListsPackagesAndSymbols(t *testing.T) {
	ctx := context.Background()
	agent := newDiscoveryTestAgent(t, ctx)

	packages, err := agent.ListGoPackages(ctx, PackageListInput{Pattern: "./..."})
	require.NoError(t, err)
	require.False(t, packages.Truncated)
	requirePackage(t, packages.Packages, "example.com/plugin")
	requirePackage(t, packages.Packages, "example.com/plugin/events")

	symbols, err := agent.GoPackageSymbols(ctx, PackageSymbolsInput{
		ImportPath: "example.com/plugin/events",
		MaxSymbols: 20,
	})
	require.NoError(t, err)
	require.Equal(t, "example.com/plugin/events", symbols.Package.ImportPath)
	requireSymbol(t, symbols.Symbols, "ReplyCreated", "type")
	requireSymbol(t, symbols.Symbols, "HandleReply", "func")
}

func TestGoDiscoveryReturnsSymbolDetails(t *testing.T) {
	ctx := context.Background()
	agent := newDiscoveryTestAgent(t, ctx)

	detail, err := agent.GoSymbolDetail(ctx, SymbolDetailInput{
		ImportPath: "example.com/plugin/events",
		Symbol:     "ReplyCreated",
	})
	require.NoError(t, err)
	require.Equal(t, "ReplyCreated", detail.Symbol.Name)
	require.Contains(t, detail.Symbol.Doc, "ReplyCreated is emitted")
	requireField(t, detail.Fields, "ID", "string")
	requireField(t, detail.Fields, "Body", "string")
	requireMethod(t, detail.Methods, "ThreadID")
}

func TestGoDiscoverySearchesSymbols(t *testing.T) {
	ctx := context.Background()
	agent := newDiscoveryTestAgent(t, ctx)

	result, err := agent.GoSymbolSearch(ctx, SymbolSearchInput{
		Query:      "reply",
		Pattern:    "./...",
		MaxResults: 20,
	})
	require.NoError(t, err)
	requireSymbol(t, result.Results, "ReplyCreated", "type")
	requireSymbol(t, result.Results, "HandleReply", "func")
}

func TestGoDiscoveryRejectsRegexStyleSymbolSearch(t *testing.T) {
	ctx := context.Background()
	agent := newDiscoveryTestAgent(t, ctx)

	_, err := agent.GoSymbolSearch(ctx, SymbolSearchInput{
		Query:   "Event.*Reply",
		Pattern: "./...",
	})
	require.ErrorContains(t, err, "literal substring")
}

func newDiscoveryTestAgent(t *testing.T, ctx context.Context) *Agent {
	t.Helper()

	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	writeWorkspaceFile(t, ctx, workspace, "go.mod", "module example.com/plugin\n\ngo 1.24\n")
	writeWorkspaceFile(t, ctx, workspace, "main.go", `package main

import "example.com/plugin/events"

func main() {
	events.HandleReply(events.ReplyCreated{ID: "reply-1", Body: "hello"})
}
`)
	writeWorkspaceFile(t, ctx, workspace, "events/reply.go", `package events

// ReplyCreated is emitted when a reply is created.
type ReplyCreated struct {
	ID       string
	Body     string
	ThreadIDValue string
}

// ThreadID returns the containing thread ID.
func (r ReplyCreated) ThreadID() string {
	return r.ThreadIDValue
}

// HandleReply handles a reply event.
func HandleReply(event ReplyCreated) string {
	return event.ID
}
`)

	return &Agent{workspace: workspace}
}

func requirePackage(t *testing.T, packages []PackageInfo, importPath string) {
	t.Helper()
	for _, pkg := range packages {
		if pkg.ImportPath == importPath {
			return
		}
	}
	require.Failf(t, "missing package", "expected package %s in %#v", importPath, packages)
}

func requireSymbol(t *testing.T, symbols []GoSymbolSummary, name, kind string) {
	t.Helper()
	for _, symbol := range symbols {
		if symbol.Name == name && symbol.Kind == kind {
			return
		}
	}
	require.Failf(t, "missing symbol", "expected %s %s in %#v", kind, name, symbols)
}

func requireField(t *testing.T, fields []GoFieldInfo, name, typ string) {
	t.Helper()
	for _, field := range fields {
		if field.Name == name && field.Type == typ {
			return
		}
	}
	require.Failf(t, "missing field", "expected field %s %s in %#v", name, typ, fields)
}

func requireMethod(t *testing.T, methods []GoMethodInfo, name string) {
	t.Helper()
	for _, method := range methods {
		if method.Name == name {
			return
		}
	}
	require.Failf(t, "missing method", "expected method %s in %#v", name, methods)
}

func writeWorkspaceFile(t *testing.T, ctx context.Context, workspace workspaceprovider.Workspace, path, content string) {
	t.Helper()
	_, err := workspace.WriteFile(ctx, path, []byte(content))
	require.NoError(t, err)
}
