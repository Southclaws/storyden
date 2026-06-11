package visibility

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

type VisibilityCommand *cobra.Command

func New(store *config.Store) VisibilityCommand {
	var fromStdin bool
	var dryRun bool
	var format string

	command := &cobra.Command{
		Use:   "visibility <slug>... <visibility>",
		Short: "Update node visibility",
		Long: `# Update Node Visibility

Control who can see your node and where it appears. The last positional argument is always the new visibility; everything before it is treated as a node identifier (slug or xid). Pass ` + "`--from-stdin`" + ` to read identifiers from stdin and supply only the visibility positionally.

## Visibility Levels

- **draft** - Only you can see it (work in progress)
- **review** - Moderators can see it (pending approval)
- **published** - Everyone can see and find it (public content)
- **unlisted** - Accessible by direct link but not searchable (shared but not promoted)

## Examples

Publish a draft:
~~~bash
sd node visibility my-page published
~~~

Publish several pages at once:
~~~bash
sd node visibility page-1 page-2 page-3 published
~~~

Preview without making changes:
~~~bash
sd node visibility article review --dry-run
~~~

Pipe identifiers from a list command:
~~~bash
sd node list --visibility review --link-domain tenor.com --format jsonl \
  | sd node visibility published --from-stdin
~~~
`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			visibility := args[len(args)-1]
			slugArgs := args[:len(args)-1]

			if err := validateVisibility(visibility); err != nil {
				return err
			}

			ids, err := resolveIDs(slugArgs, fromStdin)
			if err != nil {
				return err
			}
			if len(ids) == 0 {
				return fmt.Errorf("no nodes specified; pass slugs as positional args or use --from-stdin")
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
					node, err := updateVisibility(ctx, client.OpenAPI, id, visibility)
					if err != nil {
						return "", err
					}
					return fmt.Sprintf("%s → %s", node.Name, node.Visibility), nil
				},
				batch.Options{DryRun: dryRun},
				emit,
				summary,
			)
			if !ok {
				return fmt.Errorf("one or more visibility updates failed")
			}
			return nil
		},
	}

	command.Flags().BoolVar(&fromStdin, "from-stdin", false, "Read identifiers from stdin (one per line, plain or JSONL)")
	command.Flags().BoolVar(&dryRun, "dry-run", false, "Show the plan without making any changes")
	command.Flags().StringVar(&format, "format", "", "Per-item result format: plain (default), jsonl")

	help.SetupMarkdownHelp(command)

	return VisibilityCommand(command)
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

func updateVisibility(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	slug string,
	visibility string,
) (*openapi.NodeWithChildren, error) {
	vis := openapi.Visibility(visibility)

	response, err := client.NodeUpdateVisibilityWithResponse(ctx, slug, openapi.VisibilityMutationProps{
		Visibility: vis,
	})
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, visibilityUpdateError(response)
	}

	return response.JSON200, nil
}

func visibilityUpdateError(response *openapi.NodeUpdateVisibilityResponse) error {
	if response.StatusCode() == http.StatusNotFound {
		return fmt.Errorf("node not found")
	}

	if response.StatusCode() == http.StatusUnauthorized {
		return fmt.Errorf("visibility update request was not authorised; run sd auth login again")
	}

	body := strings.TrimSpace(string(response.Body))
	if body != "" {
		return fmt.Errorf("visibility update request failed: %s: %s", response.Status(), body)
	}

	return fmt.Errorf("visibility update request failed: %s", response.Status())
}

func validateVisibility(visibility string) error {
	switch openapi.Visibility(visibility) {
	case openapi.VisibilityDraft, openapi.VisibilityReview, openapi.VisibilityPublished, openapi.VisibilityUnlisted:
		return nil
	default:
		return fmt.Errorf("invalid visibility %q; must be one of: draft, review, published, unlisted", visibility)
	}
}
