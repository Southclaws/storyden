package pluginbuilder

import (
	"context"
	"fmt"
	"strings"

	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider"
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
	return commandResult(workspace.Run(ctx, workspaceprovider.CommandSpec{Command: "gofmt", Args: []string{"-w", "."}}))
}

func (a *Agent) GoVet(ctx context.Context) (CommandResult, error) {
	workspace, err := a.Workspace(ctx)
	if err != nil {
		return CommandResult{}, err
	}
	result, err := commandResult(workspace.Run(ctx, workspaceprovider.CommandSpec{Command: "go", Args: []string{"vet", "./..."}}))
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
	return commandResult(workspace.Run(ctx, workspaceprovider.CommandSpec{Command: "go", Args: []string{"mod", "tidy"}}))
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
	return commandResult(workspace.Run(ctx, workspaceprovider.CommandSpec{Command: "go", Args: []string{"test", pattern}}))
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
