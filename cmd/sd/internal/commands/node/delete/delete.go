package delete

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/batch"
	"github.com/Southclaws/storyden/cmd/sd/internal/commands/listflags"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
)

type DeleteCommand *cobra.Command

func New(store *config.Store) DeleteCommand {
	var target string
	var fromStdin bool
	var dryRun bool
	var format string

	command := &cobra.Command{
		Use:   "delete <slug>...",
		Short: "Delete one or more nodes",
		Long: `# Delete a Node

Permanently delete one or more nodes. Children of each deleted node are moved to a new parent.

**Warning**: This operation cannot be undone! Use ` + "`--dry-run`" + ` to preview the plan.

## What Happens to Children?

When you delete a node, its children are automatically moved:
- By default, they move to the deleted node's parent
- Or use ` + "`--target`" + ` to specify a different parent
- Root-level nodes' children become root-level

## Examples

Delete a node:
~~~bash
sd node delete old-page
~~~

Delete several nodes at once:
~~~bash
sd node delete old-page archive-1 archive-2
~~~

Preview without making any changes:
~~~bash
sd node delete old-page --dry-run
~~~

Pipe identifiers from a list command:
~~~bash
sd node list --visibility draft --format jsonl | sd node delete --from-stdin
~~~

Delete and move children to a specific parent:
~~~bash
sd node delete deprecated --target archive
~~~
`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ids, err := resolveIDs(args, fromStdin)
			if err != nil {
				return err
			}
			if len(ids) == 0 {
				return fmt.Errorf("no nodes specified; pass slugs as arguments or use --from-stdin")
			}
			if format != "" && format != listflags.FormatPlain && format != listflags.FormatJSONL {
				return fmt.Errorf("--format must be plain or jsonl")
			}

			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			emit, summary := emitters(cmd, format)

			ok := batch.Run(cmd.Context(), ids,
				func(ctx context.Context, id string) (string, error) {
					return id, deleteNode(ctx, client.OpenAPI, id, target)
				},
				batch.Options{DryRun: dryRun},
				emit,
				summary,
			)
			if !ok {
				return fmt.Errorf("one or more deletions failed")
			}
			return nil
		},
	}

	command.Flags().StringVar(&target, "target", "", "Target node slug for orphaned children")
	command.Flags().BoolVar(&fromStdin, "from-stdin", false, "Read identifiers from stdin (one per line, plain or JSONL)")
	command.Flags().BoolVar(&dryRun, "dry-run", false, "Show the plan without making any changes")
	command.Flags().StringVar(&format, "format", "", "Per-item result format: plain (default), jsonl")

	help.SetupMarkdownHelp(command)

	return DeleteCommand(command)
}

func resolveIDs(args []string, fromStdin bool) ([]string, error) {
	if fromStdin {
		ids, err := batch.ReadIdentifiers(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("read stdin: %w", err)
		}
		return append(ids, args...), nil
	}
	return args, nil
}

// emitters returns (per-item emit, end-of-run summary). JSONL routes per-item
// results to stdout so they can be piped; plain routes to stderr so stdout
// stays empty for shell composition.
func emitters(cmd *cobra.Command, format string) (func(batch.Result), func(int, int)) {
	if format == listflags.FormatJSONL {
		return batch.JSONLEmitter(cmd.OutOrStdout()), nil
	}
	return batch.PlainEmitter(cmd.ErrOrStderr()), batch.PlainSummary(cmd.ErrOrStderr())
}

func deleteNode(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	slug string,
	target string,
) error {
	params := &openapi.NodeDeleteParams{}
	if target != "" {
		targetQuery := openapi.TargetNodeSlugQuery(target)
		params.TargetNode = &targetQuery
	}

	response, err := client.NodeDeleteWithResponse(ctx, slug, params)
	if err != nil {
		return err
	}

	if response.StatusCode() != http.StatusOK {
		return nodeDeleteError(response)
	}

	return nil
}

func nodeDeleteError(response *openapi.NodeDeleteResponse) error {
	if response.StatusCode() == http.StatusNotFound {
		return fmt.Errorf("node not found")
	}

	if response.StatusCode() == http.StatusUnauthorized {
		return fmt.Errorf("node delete request was not authorised; run sd auth login again")
	}

	body := strings.TrimSpace(string(response.Body))
	if body != "" {
		return fmt.Errorf("node delete request failed: %s: %s", response.Status(), body)
	}

	return fmt.Errorf("node delete request failed: %s", response.Status())
}
