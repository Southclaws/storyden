package move

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

type MoveCommand *cobra.Command

func New(store *config.Store) MoveCommand {
	var parent string
	var before string
	var after string
	var toRoot bool
	var fromStdin bool
	var dryRun bool
	var format string

	command := &cobra.Command{
		Use:   "move <slug>...",
		Short: "Move one or more nodes in the hierarchy",
		Long: `# Move a Node

Reorganize your content by moving nodes to different parents or reordering siblings. Accepts one or more identifiers as positional arguments, or reads them from stdin with ` + "`--from-stdin`" + `.

## Examples

Move to a different parent:
~~~bash
sd node move my-page --parent docs
~~~

Move several pages to the same parent:
~~~bash
sd node move page-1 page-2 page-3 --parent docs
~~~

Move to root level:
~~~bash
sd node move my-page --to-root
~~~

Pipe identifiers from a list command:
~~~bash
sd node list --visibility review --link-domain tenor.com --format jsonl \
  | sd node move --parent random-fun-stuff --from-stdin
~~~

Reorder siblings (one at a time):
~~~bash
sd node move chapter-2 --before chapter-1
sd node move intro --after preface
~~~

Move and position in one go:
~~~bash
sd node move guide --parent tutorials --after basics
~~~

Preview without making any changes:
~~~bash
sd node move my-page --parent docs --dry-run
~~~

Note: You cannot use both ` + "`--before`" + ` and ` + "`--after`" + ` together. Reordering flags only make sense with a single target node.
`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if before != "" && after != "" {
				return fmt.Errorf("cannot specify both --before and --after")
			}
			if parent != "" && toRoot {
				return fmt.Errorf("cannot specify both --parent and --to-root")
			}
			ids, err := resolveIDs(args, fromStdin)
			if err != nil {
				return err
			}
			if len(ids) == 0 {
				return fmt.Errorf("no nodes specified; pass slugs as arguments or use --from-stdin")
			}
			if (before != "" || after != "") && len(ids) > 1 {
				return fmt.Errorf("--before/--after only make sense with a single node, not %d", len(ids))
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
					props := buildProps(parent, before, after, toRoot)
					node, err := moveNode(ctx, client.OpenAPI, id, props)
					if err != nil {
						return "", err
					}
					return string(node.Name), nil
				},
				batch.Options{DryRun: dryRun},
				emit,
				summary,
			)
			if !ok {
				return fmt.Errorf("one or more moves failed")
			}
			return nil
		},
	}

	command.Flags().StringVar(&parent, "parent", "", "New parent node slug or xid")
	command.Flags().BoolVar(&toRoot, "to-root", false, "Move node to root (remove parent)")
	command.Flags().StringVar(&before, "before", "", "Move before this sibling node")
	command.Flags().StringVar(&after, "after", "", "Move after this sibling node")
	command.Flags().BoolVar(&fromStdin, "from-stdin", false, "Read identifiers from stdin (one per line, plain or JSONL)")
	command.Flags().BoolVar(&dryRun, "dry-run", false, "Show the plan without making any changes")
	command.Flags().StringVar(&format, "format", "", "Per-item result format: plain (default), jsonl")

	help.SetupMarkdownHelp(command)

	return MoveCommand(command)
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

func emitters(cmd *cobra.Command, format string) (func(batch.Result), func(int, int)) {
	if format == listflags.FormatJSONL {
		return batch.JSONLEmitter(cmd.OutOrStdout()), nil
	}
	return batch.PlainEmitter(cmd.ErrOrStderr()), batch.PlainSummary(cmd.ErrOrStderr())
}

func buildProps(parent, before, after string, toRoot bool) openapi.NodePositionMutableProps {
	props := openapi.NodePositionMutableProps{}
	if parent != "" {
		props.Parent.Set(parent)
	} else if toRoot {
		props.Parent.Set("")
		props.Parent.SetNull()
	}
	if before != "" {
		props.Before = &before
	} else if after != "" {
		props.After = &after
	}
	return props
}

func moveNode(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	slug string,
	props openapi.NodePositionMutableProps,
) (*openapi.NodeWithChildren, error) {
	response, err := client.NodeUpdatePositionWithResponse(ctx, slug, props)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, nodeMoveError(response)
	}

	return response.JSON200, nil
}

func nodeMoveError(response *openapi.NodeUpdatePositionResponse) error {
	if response.StatusCode() == http.StatusNotFound {
		return fmt.Errorf("node not found")
	}

	if response.StatusCode() == http.StatusUnauthorized {
		return fmt.Errorf("node move request was not authorised; run sd auth login again")
	}

	if response.StatusCode() == http.StatusBadRequest {
		body := strings.TrimSpace(string(response.Body))
		if body != "" {
			return fmt.Errorf("invalid move request: %s", body)
		}

		return fmt.Errorf("invalid move request")
	}

	body := strings.TrimSpace(string(response.Body))
	if body != "" {
		return fmt.Errorf("node move request failed: %s: %s", response.Status(), body)
	}

	return fmt.Errorf("node move request failed: %s", response.Status())
}
