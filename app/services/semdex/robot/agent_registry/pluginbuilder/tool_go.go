package pluginbuilder

import (
	"context"
	"fmt"
	"strings"

	adkagent "google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/tool/functiontool"

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
		Description: "Run gofmt -w . in the managed plugin workspace.",
	}, func(ctx adkagent.Context, args struct{}) (CommandResult, error) {
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
		Description: "Run go vet ./... and Plugin Builder semantic lint checks in the managed plugin workspace.",
	}, func(ctx adkagent.Context, args struct{}) (CommandResult, error) {
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
		Description: "Run go mod tidy in the managed plugin workspace to resolve module dependencies and write go.sum.",
	}, func(ctx adkagent.Context, args struct{}) (CommandResult, error) {
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
		Description: "Run go test in the managed plugin workspace. The optional pattern defaults to ./...",
	}, func(ctx adkagent.Context, args GoTestInput) (CommandResult, error) {
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
