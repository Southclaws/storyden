package pluginbuilder

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
	adkagent "google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool/functiontool"

	"github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider"
	plugindev "github.com/Southclaws/storyden/lib/plugin/dev"
)

const (
	maxDiscoveryFiles    = 1000
	maxDiscoveryFileSize = 2_000_000
)

type PackageListResult = plugindev.PackageListResult
type PackageInfo = plugindev.PackageInfo
type PackageSymbolsResult = plugindev.PackageSymbolsResult
type SymbolDetailResult = plugindev.SymbolDetailResult
type SymbolSearchResult = plugindev.SymbolSearchResult
type GoSymbolSummary = plugindev.GoSymbolSummary
type GoFieldInfo = plugindev.GoFieldInfo
type GoMethodInfo = plugindev.GoMethodInfo

type PackageListInput struct {
	Pattern     string `json:"pattern,omitempty" jsonschema:"Go package pattern. Use ./... for workspace packages or a full import path such as github.com/Southclaws/storyden/lib/plugin/rpc/..."`
	IncludeDeps bool   `json:"include_deps,omitempty" jsonschema:"Include transitive dependencies returned by go list -deps"`
	MaxPackages int    `json:"max_packages,omitempty" jsonschema:"Maximum packages to return"`
}

type PackageSymbolsInput struct {
	ImportPath        string `json:"import_path" jsonschema:"Full Go package import path to inspect"`
	IncludeUnexported bool   `json:"include_unexported,omitempty" jsonschema:"Include unexported package symbols"`
	MaxSymbols        int    `json:"max_symbols,omitempty" jsonschema:"Maximum symbols to return"`
}

type SymbolDetailInput struct {
	ImportPath string `json:"import_path" jsonschema:"Full Go package import path that owns the symbol"`
	Symbol     string `json:"symbol" jsonschema:"Exported or unexported symbol name to inspect"`
}

type SymbolSearchInput struct {
	Query             string `json:"query" jsonschema:"Case-insensitive literal substring to search for. This is not regex or wildcard search."`
	Pattern           string `json:"pattern,omitempty" jsonschema:"Package pattern to search. Defaults to ./... for workspace packages."`
	IncludeDeps       bool   `json:"include_deps,omitempty" jsonschema:"Also search imported dependency packages"`
	IncludeUnexported bool   `json:"include_unexported,omitempty" jsonschema:"Include unexported symbols"`
	MaxResults        int    `json:"max_results,omitempty" jsonschema:"Maximum symbol matches to return"`
}

func (a *Agent) addGoDiscoveryTools(add toolAdder) error {
	if err := add(functiontool.New(functiontool.Config{
		Name: "plugin_go_packages",
		Description: `List Go packages visible from the plugin workspace.

Use this as a directory listing for Go code. Start with pattern "./..." to see
workspace packages. Use a full package prefix with "/..." to discover subpackages,
for example "github.com/Southclaws/storyden/lib/plugin/..." or
"github.com/bwmarrin/discordgo/...".

Every result includes import_path. Pass that exact import_path to
plugin_go_package_symbols or plugin_go_symbol_detail to dive deeper.`,
	}, func(ctx adkagent.Context, args PackageListInput) (PackageListResult, error) {
		return a.ListGoPackages(ctx, args)
	})); err != nil {
		return err
	}

	if err := add(functiontool.New(functiontool.Config{
		Name: "plugin_go_package_symbols",
		Description: `List symbols exported by one Go package.

Use this after plugin_go_packages or when you know the import path. This is the
main discovery tool for finding SDK capabilities without guessing method names.
For example, inspect "github.com/Southclaws/storyden/lib/plugin/rpc" to discover
available event payload types, or inspect a third-party package after adding it
to go.mod and running plugin_go_tidy.

The response includes symbol names, kinds, signatures, docs, imports, and the
owning import_path so you can recursively inspect related packages.`,
	}, func(ctx adkagent.Context, args PackageSymbolsInput) (PackageSymbolsResult, error) {
		return a.GoPackageSymbols(ctx, args)
	})); err != nil {
		return err
	}

	if err := add(functiontool.New(functiontool.Config{
		Name: "plugin_go_symbol_detail",
		Description: `Inspect one symbol in one Go package.

Use this when package_symbols shows a promising type or function. For types, this
returns fields, methods, and interface methods where available. This is the
closest tool to editor "go to definition" plus "show methods" for the plugin
builder.`,
	}, func(ctx adkagent.Context, args SymbolDetailInput) (SymbolDetailResult, error) {
		return a.GoSymbolDetail(ctx, args)
	})); err != nil {
		return err
	}

	return add(functiontool.New(functiontool.Config{
		Name: "plugin_go_symbol_search",
		Description: `Search symbols by name or documentation across Go packages.

Use this when you know the concept but not the API name, such as "reply",
"reaction", "discord", "message", or "event". Start with pattern "./..." for
workspace code. Use a full package prefix with "/..." to search a dependency
family. Set include_deps only when you need broader dependency search because it
can return a lot of packages.

The query is a literal case-insensitive substring, not a regex, glob, or
wildcard. Do not send queries such as "Event.*Reply"; search for "Event",
"Reply", or another plain word instead.`,
	}, func(ctx adkagent.Context, args SymbolSearchInput) (SymbolSearchResult, error) {
		return a.GoSymbolSearch(ctx, args)
	}))
}

func (a *Agent) ListGoPackages(ctx context.Context, in PackageListInput) (PackageListResult, error) {
	return withAgentWorkspaceMirror(ctx, a, func(dir string) (PackageListResult, error) {
		return plugindev.ListGoPackages(ctx, dir, plugindev.PackageListOptions{
			Pattern:     in.Pattern,
			IncludeDeps: in.IncludeDeps,
			MaxPackages: in.MaxPackages,
		})
	})
}

func (a *Agent) GoPackageSymbols(ctx context.Context, in PackageSymbolsInput) (PackageSymbolsResult, error) {
	return withAgentWorkspaceMirror(ctx, a, func(dir string) (PackageSymbolsResult, error) {
		return plugindev.GoPackageSymbols(ctx, dir, plugindev.PackageSymbolsOptions{
			ImportPath:        in.ImportPath,
			IncludeUnexported: in.IncludeUnexported,
			MaxSymbols:        in.MaxSymbols,
		})
	})
}

func (a *Agent) GoSymbolDetail(ctx context.Context, in SymbolDetailInput) (SymbolDetailResult, error) {
	return withAgentWorkspaceMirror(ctx, a, func(dir string) (SymbolDetailResult, error) {
		return plugindev.GoSymbolDetail(ctx, dir, plugindev.SymbolDetailOptions{
			ImportPath: in.ImportPath,
			Symbol:     in.Symbol,
		})
	})
}

func (a *Agent) GoSymbolSearch(ctx context.Context, in SymbolSearchInput) (SymbolSearchResult, error) {
	return withAgentWorkspaceMirror(ctx, a, func(dir string) (SymbolSearchResult, error) {
		return plugindev.GoSymbolSearch(ctx, dir, plugindev.SymbolSearchOptions{
			Query:             in.Query,
			Pattern:           in.Pattern,
			IncludeDeps:       in.IncludeDeps,
			IncludeUnexported: in.IncludeUnexported,
			MaxResults:        in.MaxResults,
		})
	})
}

func (a *Agent) loadGoPackage(ctx context.Context, importPath string) (*packages.Package, error) {
	return withAgentWorkspaceMirror(ctx, a, func(dir string) (*packages.Package, error) {
		return plugindev.LoadGoPackage(ctx, dir, importPath)
	})
}

func withAgentWorkspaceMirror[T any](ctx context.Context, a *Agent, fn func(dir string) (T, error)) (T, error) {
	var zero T
	workspace, err := a.Workspace(ctx)
	if err != nil {
		return zero, err
	}

	root, err := mirrorWorkspace(ctx, workspace)
	if err != nil {
		return zero, err
	}
	defer os.RemoveAll(root)

	return fn(root)
}

func mirrorWorkspace(ctx context.Context, workspace workspaceprovider.Workspace) (string, error) {
	root, err := os.MkdirTemp("", "storyden-plugin-discovery-*")
	if err != nil {
		return "", err
	}

	files, err := workspace.List(ctx, workspaceprovider.ListOptions{MaxFiles: maxDiscoveryFiles})
	if err != nil {
		_ = os.RemoveAll(root)
		return "", err
	}
	for _, file := range files {
		if file.Size > maxDiscoveryFileSize {
			continue
		}
		if strings.HasSuffix(file.Path, ".zip") {
			continue
		}
		data, err := workspace.ReadFile(ctx, file.Path, maxDiscoveryFileSize)
		if err != nil {
			_ = os.RemoveAll(root)
			return "", err
		}
		if data.Truncated {
			continue
		}
		dst := filepath.Join(root, filepath.FromSlash(data.Path))
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			_ = os.RemoveAll(root)
			return "", err
		}
		if err := os.WriteFile(dst, data.Content, 0o644); err != nil {
			_ = os.RemoveAll(root)
			return "", err
		}
	}

	return root, nil
}
