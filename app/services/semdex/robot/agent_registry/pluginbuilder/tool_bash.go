package pluginbuilder

import (
	"context"
	"errors"
	"fmt"
	"time"

	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"

	"github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider"
)

const maxBashTimeoutSeconds = 300

type RunBashInput struct {
	Command        string `json:"command" jsonschema:"Bash command to run in the managed plugin workspace"`
	Stdin          string `json:"stdin,omitempty" jsonschema:"Optional stdin to pass to the command"`
	TimeoutSeconds int    `json:"timeout_seconds,omitempty" jsonschema:"Optional timeout in seconds, capped at 300"`
}

func (a *Agent) addBashTools(add toolAdder) error {
	return add(functiontool.New(functiontool.Config{
		Name:        "plugin_run_bash",
		Description: "Run one synchronous Bash command in the managed plugin workspace. Only available when the workspace template allows untrusted commands. Use for one-shot inspection or commands not covered by focused tools. Do not use for long-running dev servers, background processes, secrets, deployment, packaging, or install; prefer focused plugin_* tools whenever they exist.",
	}, func(ctx adktool.Context, args RunBashInput) (CommandResult, error) {
		return a.RunBash(ctx, args)
	}))
}

func (a *Agent) RunBash(ctx context.Context, in RunBashInput) (CommandResult, error) {
	if !pluginBuilderAllowUntrustedCommandsFromContext(ctx) {
		return CommandResult{}, errors.New("active workspace does not allow untrusted commands")
	}
	if in.Command == "" {
		return CommandResult{}, errors.New("command is required")
	}

	workspace, err := a.Workspace(ctx)
	if err != nil {
		return CommandResult{}, err
	}

	timeout := time.Duration(in.TimeoutSeconds) * time.Second
	if in.TimeoutSeconds < 0 {
		return CommandResult{}, fmt.Errorf("timeout_seconds must be non-negative")
	}
	if in.TimeoutSeconds > maxBashTimeoutSeconds {
		timeout = maxBashTimeoutSeconds * time.Second
	}

	return commandResult(workspace.Run(ctx, workspaceprovider.CommandSpec{
		Command: "bash",
		Args:    []string{"-lc", in.Command},
		Stdin:   in.Stdin,
		Timeout: timeout,
	}))
}
