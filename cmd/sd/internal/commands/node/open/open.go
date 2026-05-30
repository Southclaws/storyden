// Package open implements `sd node open <slug>`. It fetches a node and emits
// its attached link URL. The default is print-only so agents in a non-TTY
// pipeline don't accidentally launch a browser; pass --launch to actually
// open the URL.
package open

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/nodeapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
)

type OpenCommand *cobra.Command

func New(store *config.Store) OpenCommand {
	var launch bool
	var force bool

	command := &cobra.Command{
		Use:   "open <slug>",
		Short: "Print or launch a node's attached link URL",
		Long: `# Open a Node

Fetch a node and print the URL of its attached link. Pass ` + "`--launch`" + ` to actually open the URL in your default browser; the default just prints it so agents in non-TTY pipelines stay quiet and predictable.

Use this during triage to quickly inspect what a link points to before deciding where it belongs.

## Examples

Print a node's link URL:
~~~bash
sd node open my-page
~~~

Launch in the browser (only if a terminal is attached):
~~~bash
sd node open my-page --launch
~~~

Launch even from a non-TTY environment:
~~~bash
sd node open my-page --launch --force
~~~

If the node has no attached link, the command errors with a clear message.
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			slug := args[0]

			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			node, err := nodeapi.Fetch(cmd.Context(), client.OpenAPI, slug)
			if err != nil {
				return err
			}

			if node.Link == nil {
				return fmt.Errorf("node %s has no attached link", slug)
			}

			url := string(node.Link.Url)
			fmt.Fprintln(cmd.OutOrStdout(), url)

			if !launch {
				return nil
			}

			if !output.IsTerminal(cmd.OutOrStdout()) && !force {
				return fmt.Errorf("refusing to launch browser from non-tty; pass --force to override")
			}

			return launchURL(url)
		},
	}

	command.Flags().BoolVar(&launch, "launch", false, "Open the URL in the system browser")
	command.Flags().BoolVar(&force, "force", false, "Launch even when stdout is not a TTY")

	help.SetupMarkdownHelp(command)

	return OpenCommand(command)
}

func launchURL(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("launch URL: %w", err)
	}
	return nil
}
