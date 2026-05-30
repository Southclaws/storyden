package get

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/render"
	"github.com/Southclaws/storyden/cmd/sd/internal/threadapi"
)

type GetCommand *cobra.Command

const (
	formatJSON     = "json"
	formatMarkdown = "markdown"
	formatYAML     = "yaml"
)

func New(store *config.Store) GetCommand {
	var format string

	command := &cobra.Command{
		Use:   "get <thread-mark>",
		Short: "Get a thread by its mark",
		Long: `# Get a Thread

Retrieve and display a discussion thread.

The default output renders the thread body and metadata richly in your terminal.

## Examples

View a thread:
~~~bash
sd thread get d8277oeot5p4b8gbvm60-check-out-the-new-storyden-command-line-interface
~~~

Get thread as JSON:
~~~bash
sd thread get my-thread --format json
~~~

Get thread as YAML:
~~~bash
sd thread get my-thread --format yaml
~~~
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mark := args[0]

			if err := validateFormat(format); err != nil {
				return err
			}

			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			thread, err := threadapi.Fetch(cmd.Context(), client.OpenAPI, mark)
			if err != nil {
				return err
			}

			switch format {
			case formatJSON:
				return render.ThreadJSON(cmd.OutOrStdout(), thread)
			case formatMarkdown:
				return render.ThreadMarkdown(cmd.OutOrStdout(), thread)
			case formatYAML:
				return render.ThreadYAML(cmd.OutOrStdout(), thread)
			default:
				return fmt.Errorf("unsupported format %q", format)
			}
		},
	}

	command.Flags().StringVarP(&format, "format", "f", formatMarkdown, "Output format (markdown, json, yaml)")

	help.SetupMarkdownHelp(command)

	return GetCommand(command)
}

func validateFormat(format string) error {
	switch format {
	case formatJSON, formatMarkdown, formatYAML:
		return nil
	default:
		return fmt.Errorf("--format must be one of: json, markdown, yaml")
	}
}
