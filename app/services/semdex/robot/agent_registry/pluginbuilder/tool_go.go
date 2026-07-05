package pluginbuilder

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/mod/modfile"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider"
)

const storydenModulePath = "github.com/Southclaws/storyden"
const storydenModuleCurrentRef = storydenModulePath + "@main"

const (
	goFormatTimeout = 30 * time.Second
	goTidyTimeout   = 5 * time.Minute
	goVetTimeout    = 5 * time.Minute
	goTestTimeout   = 5 * time.Minute
)

type GoTestInput struct {
	Pattern string `json:"pattern" jsonschema:"Go test package pattern, defaults to ./..."`
}

type CommandResult struct {
	Command    string `json:"command"`
	Success    bool   `json:"success"`
	Output     string `json:"output,omitempty"`
	Error      string `json:"error,omitempty"`
	Truncated  bool   `json:"truncated"`
	DurationMS int64  `json:"duration_ms"`
}

func (a *Agent) addGoTools(add toolAdder) error {
	if err := add(functiontool.New(functiontool.Config{
		Name:        "plugin_go_fmt",
		Description: "Format all Go files in the managed plugin workspace with gofmt. Use after editing Go code or when plugin_validate reports go_fmt. Side effect: rewrites Go files.",
	}, func(ctx adktool.Context, args struct{}) (CommandResult, error) {
		result, err := a.GoFormat(ctx)
		if err != nil {
			return CommandResult{}, err
		}
		return result, nil
	})); err != nil {
		return err
	}

	if err := add(functiontool.New(functiontool.Config{
		Name:        "plugin_go_vet",
		Description: "Run go vet ./... plus Plugin Builder semantic lint checks. Use to diagnose compile/lint readiness after code edits. Does not install or package the plugin. If this fails, fix the reported Go or plugin-lifecycle issue before plugin_install.",
	}, func(ctx adktool.Context, args struct{}) (CommandResult, error) {
		result, err := a.GoVet(ctx)
		if err != nil {
			return CommandResult{}, err
		}
		return result, nil
	})); err != nil {
		return err
	}

	if err := add(functiontool.New(functiontool.Config{
		Name:        "plugin_go_tidy",
		Description: "Run go mod tidy to resolve imports and update go.mod/go.sum. Use after adding, removing, or changing imports or dependencies. Side effect: rewrites module files.",
	}, func(ctx adktool.Context, args struct{}) (CommandResult, error) {
		result, err := a.GoTidy(ctx)
		if err != nil {
			return CommandResult{}, err
		}
		return result, nil
	})); err != nil {
		return err
	}

	return add(functiontool.New(functiontool.Config{
		Name:        "plugin_go_test",
		Description: "Run go test for the managed plugin workspace. The optional pattern defaults to ./.... Use after implementation changes to catch compile errors and test failures. Does not compile the final supervised binary; plugin_install does that once.",
	}, func(ctx adktool.Context, args GoTestInput) (CommandResult, error) {
		result, err := a.GoTest(ctx, args)
		if err != nil {
			return CommandResult{}, err
		}
		return result, nil
	}))
}

func (a *Agent) GoFormat(ctx context.Context) (CommandResult, error) {
	workspace, err := a.Workspace(ctx)
	if err != nil {
		return CommandResult{}, err
	}
	return commandResult(workspace.Run(ctx, workspaceprovider.CommandSpec{Command: "gofmt", Args: []string{"-w", "."}, Timeout: goFormatTimeout}))
}

func (a *Agent) GoVet(ctx context.Context) (CommandResult, error) {
	workspace, err := a.Workspace(ctx)
	if err != nil {
		return CommandResult{}, err
	}
	result, err := commandResult(workspace.Run(ctx, workspaceprovider.CommandSpec{Command: "go", Args: []string{"vet", "./..."}, Timeout: goVetTimeout}))
	if err != nil || !result.Success {
		return result, err
	}

	lint, err := a.PluginLint(ctx)
	if err != nil {
		return CommandResult{}, err
	}
	if !lint.Success {
		return CommandResult{
			Command: "plugin lint",
			Success: false,
			Output:  lint.Format(),
		}, nil
	}

	return result, nil
}

func (a *Agent) GoTidy(ctx context.Context) (CommandResult, error) {
	workspace, err := a.Workspace(ctx)
	if err != nil {
		return CommandResult{}, err
	}
	if err := ensureStorydenModuleRequirement(ctx, workspace); err != nil {
		return CommandResult{}, err
	}
	if result, err := ensureCurrentStorydenModule(ctx, workspace); err != nil || !result.Success {
		return result, err
	}
	return commandResult(workspace.Run(ctx, workspaceprovider.CommandSpec{Command: "go", Args: []string{"mod", "tidy"}, Timeout: goTidyTimeout}))
}

func (a *Agent) GoTest(ctx context.Context, in GoTestInput) (CommandResult, error) {
	workspace, err := a.Workspace(ctx)
	if err != nil {
		return CommandResult{}, err
	}
	pattern := strings.TrimSpace(in.Pattern)
	if pattern == "" {
		pattern = "./..."
	}
	if strings.HasPrefix(pattern, "-") || strings.ContainsAny(pattern, ";&|`$<>") {
		return CommandResult{}, fmt.Errorf("invalid go test pattern %q", pattern)
	}
	return commandResult(workspace.Run(ctx, workspaceprovider.CommandSpec{Command: "go", Args: []string{"test", pattern}, Timeout: goTestTimeout}))
}

func commandResult(result workspaceprovider.CommandResult, err error) (CommandResult, error) {
	return CommandResult{
		Command:    result.Command,
		Success:    result.Success,
		Output:     result.Output,
		Error:      result.Error,
		Truncated:  result.Truncated,
		DurationMS: result.DurationMS,
	}, err
}

func ensureStorydenModuleRequirement(ctx context.Context, workspace workspaceprovider.Workspace) error {
	data, err := workspace.ReadFile(ctx, "go.mod", -1)
	if err != nil {
		return err
	}

	file, err := modfile.Parse("go.mod", data.Content, nil)
	if err != nil {
		return fmt.Errorf("parse go.mod: %w", err)
	}

	if hasStorydenReplace(file) {
		return nil
	}

	changed := false
	for _, req := range file.Require {
		if req.Mod.Path != storydenModulePath || req.Mod.Version != "v0.0.0" {
			continue
		}
		if err := file.DropRequire(storydenModulePath); err != nil {
			return err
		}
		changed = true
	}
	if !changed {
		return nil
	}

	formatted, err := file.Format()
	if err != nil {
		return fmt.Errorf("format go.mod: %w", err)
	}
	_, err = workspace.WriteFile(ctx, "go.mod", formatted)
	return err
}

func ensureCurrentStorydenModule(ctx context.Context, workspace workspaceprovider.Workspace) (CommandResult, error) {
	needed, err := workspaceNeedsStorydenModule(ctx, workspace)
	if err != nil {
		return CommandResult{}, err
	}
	if !needed {
		return CommandResult{Command: "go get " + storydenModuleCurrentRef, Success: true}, nil
	}

	replaced, err := workspaceHasStorydenReplace(ctx, workspace)
	if err != nil {
		return CommandResult{}, err
	}
	if replaced {
		return CommandResult{Command: "go get " + storydenModuleCurrentRef, Success: true}, nil
	}

	return commandResult(workspace.Run(ctx, workspaceprovider.CommandSpec{
		Command: "go",
		Args:    []string{"get", storydenModuleCurrentRef},
		Timeout: goTidyTimeout,
	}))
}

func workspaceNeedsStorydenModule(ctx context.Context, workspace workspaceprovider.Workspace) (bool, error) {
	data, err := workspace.ReadFile(ctx, "go.mod", -1)
	if err != nil {
		return false, err
	}
	if strings.Contains(string(data.Content), storydenModulePath) {
		return true, nil
	}

	files, err := workspace.List(ctx, workspaceprovider.ListOptions{MaxFiles: 500})
	if err != nil {
		return false, err
	}
	for _, file := range files {
		if !strings.HasSuffix(file.Path, ".go") {
			continue
		}
		source, err := workspace.ReadFile(ctx, file.Path, -1)
		if err != nil {
			return false, err
		}
		if strings.Contains(string(source.Content), storydenModulePath) {
			return true, nil
		}
	}

	return false, nil
}

func workspaceHasStorydenReplace(ctx context.Context, workspace workspaceprovider.Workspace) (bool, error) {
	data, err := workspace.ReadFile(ctx, "go.mod", -1)
	if err != nil {
		return false, err
	}

	file, err := modfile.Parse("go.mod", data.Content, nil)
	if err != nil {
		return false, fmt.Errorf("parse go.mod: %w", err)
	}

	return hasStorydenReplace(file), nil
}

func hasStorydenReplace(file *modfile.File) bool {
	for _, replace := range file.Replace {
		if replace.Old.Path == storydenModulePath {
			return true
		}
	}
	return false
}
