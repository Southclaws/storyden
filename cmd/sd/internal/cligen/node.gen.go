package cligen

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

type NodeListLinkScheme string

const (
	NodeListLinkSchemeHttp  NodeListLinkScheme = "http"
	NodeListLinkSchemeHttps NodeListLinkScheme = "https"
)

type NodeListNodeFormat string

const (
	NodeListNodeFormatFlat NodeListNodeFormat = "flat"
	NodeListNodeFormatTree NodeListNodeFormat = "tree"
)

type NodeListFormat string

const (
	NodeListFormatAuto  NodeListFormat = "auto"
	NodeListFormatPlain NodeListFormat = "plain"
	NodeListFormatJson  NodeListFormat = "json"
	NodeListFormatJsonl NodeListFormat = "jsonl"
)

type NodeListOutput string

const (
	NodeListOutputDefault NodeListOutput = "default"
	NodeListOutputWide    NodeListOutput = "wide"
)

type NodeListParams struct {
	All             bool
	Page            int
	Limit           int
	Depth           int
	Parent          string
	NodeId          string
	RootOnly        bool
	HasLink         bool
	NoLink          bool
	LinkDomain      []string
	LinkScheme      NodeListLinkScheme
	LinkUrlContains string
	NameContains    string
	OwnerHandle     string
	Visibility      []string
	Search          string
	NodeFormat      NodeListNodeFormat
	Format          NodeListFormat
	Output          NodeListOutput
}

type NodeListHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodeListParams) (NodeListResult, error)

func newNodeListCommand(nodeList NodeListHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Storyden nodes.",
		Example: `  sd node list
  sd node list --all --format jsonl
  sd node list --link-domain youtube.com,youtu.be
  sd node list --format json > nodes.json`,
		Args: rangeArgs(0, 0),
	}

	var rawAll bool
	cmd.Flags().BoolVar(&rawAll, "all", false, "Fetch every page, streaming output as it goes.")
	var rawPage int
	cmd.Flags().IntVar(&rawPage, "page", 1, "Page to request.")
	var rawLimit int
	cmd.Flags().IntVar(&rawLimit, "limit", 0, "Stop after N matches (0 = no limit).")
	var rawDepth int
	cmd.Flags().IntVar(&rawDepth, "depth", 0, "Maximum child depth to return (0 = root nodes only).")
	var rawParent string
	cmd.Flags().StringVar(&rawParent, "parent", "", "Limit to children of this node (accepts slug; resolves to node id automatically).")
	var rawNodeId string
	cmd.Flags().StringVar(&rawNodeId, "node-id", "", "Limit to descendants of this node id.")
	var rawRootOnly bool
	cmd.Flags().BoolVar(&rawRootOnly, "root-only", false, "Filter to nodes with no parent.")
	var rawHasLink bool
	cmd.Flags().BoolVar(&rawHasLink, "has-link", false, "Filter to nodes that have a link attached.")
	var rawNoLink bool
	cmd.Flags().BoolVar(&rawNoLink, "no-link", false, "Filter to nodes that have no link attached.")
	var rawLinkDomain []string
	cmd.Flags().StringSliceVar(&rawLinkDomain, "link-domain", nil, "Filter to nodes whose link.domain matches (repeatable, comma-separated).")
	var rawLinkScheme string
	cmd.Flags().StringVar(&rawLinkScheme, "link-scheme", "", "Filter to nodes whose link.url uses this scheme.")
	var rawLinkUrlContains string
	cmd.Flags().StringVar(&rawLinkUrlContains, "link-url-contains", "", "Filter to nodes whose link.url contains this substring.")
	var rawNameContains string
	cmd.Flags().StringVar(&rawNameContains, "name-contains", "", "Filter to nodes whose name contains this substring (case-insensitive).")
	var rawOwnerHandle string
	cmd.Flags().StringVar(&rawOwnerHandle, "owner-handle", "", "Filter to nodes owned by this account handle.")
	var rawVisibility []string
	cmd.Flags().StringSliceVar(&rawVisibility, "visibility", nil, "Filter by visibility (repeatable, comma-separated).")
	var rawSearch string
	cmd.Flags().StringVarP(&rawSearch, "search", "q", "", "Server-side full-text search query.")
	var rawNodeFormat string
	cmd.Flags().StringVar(&rawNodeFormat, "node-format", "flat", "Server response shape: flat (every match) or tree (roots with nested children).")
	var rawFormat string
	cmd.Flags().StringVar(&rawFormat, "format", "auto", "Output format: auto, plain, json, jsonl.")
	var rawOutput string
	cmd.Flags().StringVarP(&rawOutput, "output", "o", "default", "Column profile: default, wide.")
	cmd.MarkFlagsMutuallyExclusive("has-link", "no-link")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var linkScheme NodeListLinkScheme
		if rawLinkScheme != "" {
			switch rawLinkScheme {
			case "http":
				linkScheme = NodeListLinkSchemeHttp
			case "https":
				linkScheme = NodeListLinkSchemeHttps
			default:
				return fmt.Errorf("invalid --link-scheme %q: must be one of http, https", rawLinkScheme)
			}
		}
		var nodeFormat NodeListNodeFormat
		if rawNodeFormat != "" {
			switch rawNodeFormat {
			case "flat":
				nodeFormat = NodeListNodeFormatFlat
			case "tree":
				nodeFormat = NodeListNodeFormatTree
			default:
				return fmt.Errorf("invalid --node-format %q: must be one of flat, tree", rawNodeFormat)
			}
		}
		var format NodeListFormat
		if rawFormat != "" {
			switch rawFormat {
			case "auto":
				format = NodeListFormatAuto
			case "plain":
				format = NodeListFormatPlain
			case "json":
				format = NodeListFormatJson
			case "jsonl":
				format = NodeListFormatJsonl
			default:
				return fmt.Errorf("invalid --format %q: must be one of auto, plain, json, jsonl", rawFormat)
			}
		}
		var output NodeListOutput
		if rawOutput != "" {
			switch rawOutput {
			case "default":
				output = NodeListOutputDefault
			case "wide":
				output = NodeListOutputWide
			default:
				return fmt.Errorf("invalid --output %q: must be one of default, wide", rawOutput)
			}
		}

		p := NodeListParams{
			All:             rawAll,
			Page:            rawPage,
			Limit:           rawLimit,
			Depth:           rawDepth,
			Parent:          rawParent,
			NodeId:          rawNodeId,
			RootOnly:        rawRootOnly,
			HasLink:         rawHasLink,
			NoLink:          rawNoLink,
			LinkDomain:      rawLinkDomain,
			LinkScheme:      linkScheme,
			LinkUrlContains: rawLinkUrlContains,
			NameContains:    rawNameContains,
			OwnerHandle:     rawOwnerHandle,
			Visibility:      rawVisibility,
			Search:          rawSearch,
			NodeFormat:      nodeFormat,
			Format:          format,
			Output:          output,
		}

		result, err := nodeList(cmd.Context(), cmd, newIO(cmd), p)
		if err != nil {
			return err
		}
		switch format {
		case NodeListFormatJson:
			return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
		}
		return nil
	}

	return cmd
}

type NodeGetFormat string

const (
	NodeGetFormatJson     NodeGetFormat = "json"
	NodeGetFormatMarkdown NodeGetFormat = "markdown"
	NodeGetFormatYaml     NodeGetFormat = "yaml"
)

type NodeGetParams struct {
	Format NodeGetFormat
	Slug   string
}

type NodeGetHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodeGetParams) (NodeWithChildren, error)

func newNodeGetCommand(nodeGet NodeGetHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get <slug>",
		Short:   "Get a node by its slug.",
		Aliases: []string{"show"},
		Args:    rangeArgs(1, 1),
	}

	var rawFormat string
	cmd.Flags().StringVarP(&rawFormat, "format", "f", "markdown", "Output format.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var format NodeGetFormat
		if rawFormat != "" {
			switch rawFormat {
			case "json":
				format = NodeGetFormatJson
			case "markdown":
				format = NodeGetFormatMarkdown
			case "yaml":
				format = NodeGetFormatYaml
			default:
				return fmt.Errorf("invalid --format %q: must be one of json, markdown, yaml", rawFormat)
			}
		}
		rawSlug := args[0]

		p := NodeGetParams{
			Format: format,
			Slug:   rawSlug,
		}

		result, err := nodeGet(cmd.Context(), cmd, newIO(cmd), p)
		if err != nil {
			return err
		}
		switch format {
		case NodeGetFormatJson:
			return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
		case NodeGetFormatYaml:
			enc := yaml.NewEncoder(cmd.OutOrStdout())
			if err := enc.Encode(result); err != nil {
				return err
			}
			return enc.Close()
		}
		return nil
	}

	return cmd
}

type NodeCreateVisibility string

const (
	NodeCreateVisibilityDraft     NodeCreateVisibility = "draft"
	NodeCreateVisibilityReview    NodeCreateVisibility = "review"
	NodeCreateVisibilityPublished NodeCreateVisibility = "published"
	NodeCreateVisibilityUnlisted  NodeCreateVisibility = "unlisted"
)

type NodeCreateParams struct {
	Name          string
	Slug          string
	Description   string
	Content       string
	ContentFile   string
	Markdown      bool
	Url           string
	Parent        string
	Tags          []string
	Visibility    NodeCreateVisibility
	HideChildTree bool
}

type NodeCreateHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodeCreateParams) error

func newNodeCreateCommand(nodeCreate NodeCreateHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new node.",
		Long: `Create a node from inline content or a file. Content may be HTML, or
Markdown that is converted to HTML when --markdown is set.
`,
		Example: `  sd node create --name "My Page"
  sd node create --name "Notes" --markdown --content "# Title"
  sd node create --name "Getting Started" --parent docs
  sd node create --name "Docs" --content-file page.html`,
		Args: rangeArgs(0, 0),
	}

	var rawName string
	cmd.Flags().StringVar(&rawName, "name", "", "Node name.")
	var rawSlug string
	cmd.Flags().StringVar(&rawSlug, "slug", "", "Node slug (auto-generated if not provided).")
	var rawDescription string
	cmd.Flags().StringVar(&rawDescription, "description", "", "Node description.")
	var rawContent string
	cmd.Flags().StringVar(&rawContent, "content", "", "Node content as HTML (or Markdown with --markdown).")
	var rawContentFile string
	cmd.Flags().StringVar(&rawContentFile, "content-file", "", "Read content from a file (use - for stdin).")
	var rawMarkdown bool
	cmd.Flags().BoolVar(&rawMarkdown, "markdown", false, "Treat content as Markdown and convert it to HTML.")
	var rawUrl string
	cmd.Flags().StringVar(&rawUrl, "url", "", "External URL to attach as the node's link.")
	var rawParent string
	cmd.Flags().StringVar(&rawParent, "parent", "", "Parent node slug.")
	var rawTags []string
	cmd.Flags().StringSliceVar(&rawTags, "tags", nil, "Tags to attach (repeatable, comma-separated).")
	var rawVisibility string
	cmd.Flags().StringVar(&rawVisibility, "visibility", "", "Node visibility.")
	var rawHideChildTree bool
	cmd.Flags().BoolVar(&rawHideChildTree, "hide-child-tree", false, "Hide this node's children from tree views.")
	_ = cmd.MarkFlagRequired("name")
	cmd.MarkFlagsMutuallyExclusive("content", "content-file")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var visibility NodeCreateVisibility
		if rawVisibility != "" {
			switch rawVisibility {
			case "draft":
				visibility = NodeCreateVisibilityDraft
			case "review":
				visibility = NodeCreateVisibilityReview
			case "published":
				visibility = NodeCreateVisibilityPublished
			case "unlisted":
				visibility = NodeCreateVisibilityUnlisted
			default:
				return fmt.Errorf("invalid --visibility %q: must be one of draft, review, published, unlisted", rawVisibility)
			}
		}

		p := NodeCreateParams{
			Name:          rawName,
			Slug:          rawSlug,
			Description:   rawDescription,
			Content:       rawContent,
			ContentFile:   rawContentFile,
			Markdown:      rawMarkdown,
			Url:           rawUrl,
			Parent:        rawParent,
			Tags:          rawTags,
			Visibility:    visibility,
			HideChildTree: rawHideChildTree,
		}

		return nodeCreate(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type NodeUpdateParams struct {
	Name             string
	NewSlug          string
	Description      string
	Content          string
	ContentFile      string
	Markdown         bool
	Url              string
	ClearUrl         bool
	HideChildTree    bool
	HideChildTreeSet bool
	Tags             []string
	Json             string
	Slug             string
}

type NodeUpdateHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodeUpdateParams) error

func newNodeUpdateCommand(nodeUpdate NodeUpdateHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <slug>",
		Short: "Update an existing node.",
		Long: `Every flag is a partial update: only the fields whose flags are set
are changed. --json accepts a full or partial NodeMutableProps
object instead (mutually exclusive with the discrete field flags),
which can also reach fields with no dedicated flag (meta, asset
ids, primary image asset id).
`,
		Args: rangeArgs(1, 1),
	}

	var rawName string
	cmd.Flags().StringVar(&rawName, "name", "", "New node name.")
	var rawNewSlug string
	cmd.Flags().StringVar(&rawNewSlug, "new-slug", "", "New slug to rename the node to.")
	var rawDescription string
	cmd.Flags().StringVar(&rawDescription, "description", "", "New node description.")
	var rawContent string
	cmd.Flags().StringVar(&rawContent, "content", "", "New node content as HTML (or Markdown with --markdown).")
	var rawContentFile string
	cmd.Flags().StringVar(&rawContentFile, "content-file", "", "Read new content from a file (use - for stdin).")
	var rawMarkdown bool
	cmd.Flags().BoolVar(&rawMarkdown, "markdown", false, "Treat --content/--content-file as Markdown and convert it to HTML.")
	var rawUrl string
	cmd.Flags().StringVar(&rawUrl, "url", "", "External URL to attach as the node's link.")
	var rawClearUrl bool
	cmd.Flags().BoolVar(&rawClearUrl, "clear-url", false, "Remove the node's attached link.")
	var rawHideChildTree bool
	cmd.Flags().BoolVar(&rawHideChildTree, "hide-child-tree", false, "Hide this node's children from tree views.")
	var rawTags []string
	cmd.Flags().StringSliceVar(&rawTags, "tags", nil, "Replace the node's tags (repeatable, comma-separated).")
	var rawJson string
	cmd.Flags().StringVar(&rawJson, "json", "", "A full or partial NodeMutableProps JSON body (use - for stdin), instead of discrete field flags.")
	cmd.MarkFlagsMutuallyExclusive("content", "content-file")
	cmd.MarkFlagsMutuallyExclusive("url", "clear-url")
	cmd.MarkFlagsMutuallyExclusive("json", "name")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		hideChildTreeSet := cmd.Flags().Changed("hide-child-tree")
		rawSlug := args[0]

		p := NodeUpdateParams{
			Name:             rawName,
			NewSlug:          rawNewSlug,
			Description:      rawDescription,
			Content:          rawContent,
			ContentFile:      rawContentFile,
			Markdown:         rawMarkdown,
			Url:              rawUrl,
			ClearUrl:         rawClearUrl,
			HideChildTree:    rawHideChildTree,
			HideChildTreeSet: hideChildTreeSet,
			Tags:             rawTags,
			Json:             rawJson,
			Slug:             rawSlug,
		}

		return nodeUpdate(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type NodeDeleteFormat string

const (
	NodeDeleteFormatPlain NodeDeleteFormat = "plain"
	NodeDeleteFormatJsonl NodeDeleteFormat = "jsonl"
)

type NodeDeleteParams struct {
	Target    string
	FromStdin bool
	DryRun    bool
	Format    NodeDeleteFormat
	Slug      []string
}

type NodeDeleteHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodeDeleteParams) error

func newNodeDeleteCommand(nodeDelete NodeDeleteHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <slug>...",
		Short: "Delete one or more nodes.",
		Long: `Deletes each named node. Any children of a deleted node are
re-parented to --target (or become root nodes if --target is
unset).
`,
		Aliases: []string{"rm"},
		Args:    rangeArgs(1, -1),
	}

	var rawTarget string
	cmd.Flags().StringVar(&rawTarget, "target", "", "Target node slug to re-parent orphaned children.")
	var rawFromStdin bool
	cmd.Flags().BoolVar(&rawFromStdin, "from-stdin", false, "Read identifiers from stdin (one per line), appended after any positional arguments.")
	var rawDryRun bool
	cmd.Flags().BoolVarP(&rawDryRun, "dry-run", "n", false, "Show the plan without making any changes.")
	var rawFormat string
	cmd.Flags().StringVar(&rawFormat, "format", "plain", "Per-item result format.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var format NodeDeleteFormat
		if rawFormat != "" {
			switch rawFormat {
			case "plain":
				format = NodeDeleteFormatPlain
			case "jsonl":
				format = NodeDeleteFormatJsonl
			default:
				return fmt.Errorf("invalid --format %q: must be one of plain, jsonl", rawFormat)
			}
		}
		slug := args[0:]

		p := NodeDeleteParams{
			Target:    rawTarget,
			FromStdin: rawFromStdin,
			DryRun:    rawDryRun,
			Format:    format,
			Slug:      slug,
		}

		return nodeDelete(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type NodeMoveFormat string

const (
	NodeMoveFormatPlain NodeMoveFormat = "plain"
	NodeMoveFormatJsonl NodeMoveFormat = "jsonl"
)

type NodeMoveParams struct {
	Parent    string
	ToRoot    bool
	Before    string
	After     string
	FromStdin bool
	DryRun    bool
	Format    NodeMoveFormat
	Slug      []string
}

type NodeMoveHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodeMoveParams) error

func newNodeMoveCommand(nodeMove NodeMoveHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "move <slug>...",
		Short: "Move one or more nodes in the hierarchy.",
		Args:  rangeArgs(1, -1),
	}

	var rawParent string
	cmd.Flags().StringVar(&rawParent, "parent", "", "New parent node slug.")
	var rawToRoot bool
	cmd.Flags().BoolVar(&rawToRoot, "to-root", false, "Move the node(s) to become root nodes.")
	var rawBefore string
	cmd.Flags().StringVar(&rawBefore, "before", "", "Position before this sibling slug (only valid with a single target node).")
	var rawAfter string
	cmd.Flags().StringVar(&rawAfter, "after", "", "Position after this sibling slug (only valid with a single target node).")
	var rawFromStdin bool
	cmd.Flags().BoolVar(&rawFromStdin, "from-stdin", false, "Read identifiers from stdin (one per line), appended after any positional arguments.")
	var rawDryRun bool
	cmd.Flags().BoolVarP(&rawDryRun, "dry-run", "n", false, "Show the plan without making any changes.")
	var rawFormat string
	cmd.Flags().StringVar(&rawFormat, "format", "plain", "Per-item result format.")
	cmd.MarkFlagsMutuallyExclusive("parent", "to-root")
	cmd.MarkFlagsMutuallyExclusive("before", "after")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var format NodeMoveFormat
		if rawFormat != "" {
			switch rawFormat {
			case "plain":
				format = NodeMoveFormatPlain
			case "jsonl":
				format = NodeMoveFormatJsonl
			default:
				return fmt.Errorf("invalid --format %q: must be one of plain, jsonl", rawFormat)
			}
		}
		slug := args[0:]

		p := NodeMoveParams{
			Parent:    rawParent,
			ToRoot:    rawToRoot,
			Before:    rawBefore,
			After:     rawAfter,
			FromStdin: rawFromStdin,
			DryRun:    rawDryRun,
			Format:    format,
			Slug:      slug,
		}

		return nodeMove(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type NodeVisibilityVisibility string

const (
	NodeVisibilityVisibilityDraft     NodeVisibilityVisibility = "draft"
	NodeVisibilityVisibilityReview    NodeVisibilityVisibility = "review"
	NodeVisibilityVisibilityPublished NodeVisibilityVisibility = "published"
	NodeVisibilityVisibilityUnlisted  NodeVisibilityVisibility = "unlisted"
)

type NodeVisibilityFormat string

const (
	NodeVisibilityFormatPlain NodeVisibilityFormat = "plain"
	NodeVisibilityFormatJsonl NodeVisibilityFormat = "jsonl"
)

type NodeVisibilityParams struct {
	Visibility NodeVisibilityVisibility
	FromStdin  bool
	DryRun     bool
	Format     NodeVisibilityFormat
	Slug       []string
}

type NodeVisibilityHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodeVisibilityParams) error

func newNodeVisibilityCommand(nodeVisibility NodeVisibilityHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "visibility <slug>...",
		Short: "Update node visibility.",
		Args:  rangeArgs(1, -1),
	}

	var rawVisibility string
	cmd.Flags().StringVar(&rawVisibility, "visibility", "", "The visibility to set.")
	var rawFromStdin bool
	cmd.Flags().BoolVar(&rawFromStdin, "from-stdin", false, "Read identifiers from stdin (one per line), appended after any positional arguments.")
	var rawDryRun bool
	cmd.Flags().BoolVarP(&rawDryRun, "dry-run", "n", false, "Show the plan without making any changes.")
	var rawFormat string
	cmd.Flags().StringVar(&rawFormat, "format", "plain", "Per-item result format.")
	_ = cmd.MarkFlagRequired("visibility")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var visibility NodeVisibilityVisibility
		if rawVisibility != "" {
			switch rawVisibility {
			case "draft":
				visibility = NodeVisibilityVisibilityDraft
			case "review":
				visibility = NodeVisibilityVisibilityReview
			case "published":
				visibility = NodeVisibilityVisibilityPublished
			case "unlisted":
				visibility = NodeVisibilityVisibilityUnlisted
			default:
				return fmt.Errorf("invalid --visibility %q: must be one of draft, review, published, unlisted", rawVisibility)
			}
		}
		var format NodeVisibilityFormat
		if rawFormat != "" {
			switch rawFormat {
			case "plain":
				format = NodeVisibilityFormatPlain
			case "jsonl":
				format = NodeVisibilityFormatJsonl
			default:
				return fmt.Errorf("invalid --format %q: must be one of plain, jsonl", rawFormat)
			}
		}
		slug := args[0:]

		p := NodeVisibilityParams{
			Visibility: visibility,
			FromStdin:  rawFromStdin,
			DryRun:     rawDryRun,
			Format:     format,
			Slug:       slug,
		}

		return nodeVisibility(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type NodeChildrenFormat string

const (
	NodeChildrenFormatAuto  NodeChildrenFormat = "auto"
	NodeChildrenFormatPlain NodeChildrenFormat = "plain"
	NodeChildrenFormatJson  NodeChildrenFormat = "json"
	NodeChildrenFormatJsonl NodeChildrenFormat = "jsonl"
)

type NodeChildrenOutput string

const (
	NodeChildrenOutputDefault NodeChildrenOutput = "default"
	NodeChildrenOutputWide    NodeChildrenOutput = "wide"
)

type NodeChildrenLinkScheme string

const (
	NodeChildrenLinkSchemeHttp  NodeChildrenLinkScheme = "http"
	NodeChildrenLinkSchemeHttps NodeChildrenLinkScheme = "https"
)

type NodeChildrenParams struct {
	All             bool
	Page            int
	Limit           int
	Format          NodeChildrenFormat
	Output          NodeChildrenOutput
	Sort            string
	LinkDomain      []string
	LinkUrlContains string
	LinkScheme      NodeChildrenLinkScheme
	HasLink         bool
	NoLink          bool
	RootOnly        bool
	OwnerHandle     string
	NameContains    string
	Slug            string
}

type NodeChildrenHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodeChildrenParams) (NodeListResult, error)

func newNodeChildrenCommand(nodeChildren NodeChildrenHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "children <slug>",
		Short: "List children of a node.",
		Args:  rangeArgs(1, 1),
	}

	var rawAll bool
	cmd.Flags().BoolVar(&rawAll, "all", false, "Fetch every page, streaming output as it goes.")
	var rawPage int
	cmd.Flags().IntVar(&rawPage, "page", 1, "Page to request.")
	var rawLimit int
	cmd.Flags().IntVar(&rawLimit, "limit", 0, "Stop after N matches (0 = no limit).")
	var rawFormat string
	cmd.Flags().StringVar(&rawFormat, "format", "auto", "Output format: auto, plain, json, jsonl.")
	var rawOutput string
	cmd.Flags().StringVarP(&rawOutput, "output", "o", "default", "Column profile: default, wide.")
	var rawSort string
	cmd.Flags().StringVar(&rawSort, "sort", "", "Server-side sort order for children.")
	var rawLinkDomain []string
	cmd.Flags().StringSliceVar(&rawLinkDomain, "link-domain", nil, "Filter to nodes whose link.domain matches (repeatable, comma-separated; applied client-side after fetching).")
	var rawLinkUrlContains string
	cmd.Flags().StringVar(&rawLinkUrlContains, "link-url-contains", "", "Filter to nodes whose link.url contains this substring (client-side).")
	var rawLinkScheme string
	cmd.Flags().StringVar(&rawLinkScheme, "link-scheme", "", "Filter to nodes whose link.url uses this scheme (client-side).")
	var rawHasLink bool
	cmd.Flags().BoolVar(&rawHasLink, "has-link", false, "Filter to nodes that have a link attached (client-side).")
	var rawNoLink bool
	cmd.Flags().BoolVar(&rawNoLink, "no-link", false, "Filter to nodes that have no link attached (client-side).")
	var rawRootOnly bool
	cmd.Flags().BoolVar(&rawRootOnly, "root-only", false, "Filter to nodes with no parent (client-side).")
	var rawOwnerHandle string
	cmd.Flags().StringVar(&rawOwnerHandle, "owner-handle", "", "Filter to nodes owned by this account handle (client-side).")
	var rawNameContains string
	cmd.Flags().StringVar(&rawNameContains, "name-contains", "", "Filter to nodes whose name contains this substring, case-insensitive (client-side).")
	cmd.MarkFlagsMutuallyExclusive("has-link", "no-link")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var format NodeChildrenFormat
		if rawFormat != "" {
			switch rawFormat {
			case "auto":
				format = NodeChildrenFormatAuto
			case "plain":
				format = NodeChildrenFormatPlain
			case "json":
				format = NodeChildrenFormatJson
			case "jsonl":
				format = NodeChildrenFormatJsonl
			default:
				return fmt.Errorf("invalid --format %q: must be one of auto, plain, json, jsonl", rawFormat)
			}
		}
		var output NodeChildrenOutput
		if rawOutput != "" {
			switch rawOutput {
			case "default":
				output = NodeChildrenOutputDefault
			case "wide":
				output = NodeChildrenOutputWide
			default:
				return fmt.Errorf("invalid --output %q: must be one of default, wide", rawOutput)
			}
		}
		var linkScheme NodeChildrenLinkScheme
		if rawLinkScheme != "" {
			switch rawLinkScheme {
			case "http":
				linkScheme = NodeChildrenLinkSchemeHttp
			case "https":
				linkScheme = NodeChildrenLinkSchemeHttps
			default:
				return fmt.Errorf("invalid --link-scheme %q: must be one of http, https", rawLinkScheme)
			}
		}
		rawSlug := args[0]

		p := NodeChildrenParams{
			All:             rawAll,
			Page:            rawPage,
			Limit:           rawLimit,
			Format:          format,
			Output:          output,
			Sort:            rawSort,
			LinkDomain:      rawLinkDomain,
			LinkUrlContains: rawLinkUrlContains,
			LinkScheme:      linkScheme,
			HasLink:         rawHasLink,
			NoLink:          rawNoLink,
			RootOnly:        rawRootOnly,
			OwnerHandle:     rawOwnerHandle,
			NameContains:    rawNameContains,
			Slug:            rawSlug,
		}

		result, err := nodeChildren(cmd.Context(), cmd, newIO(cmd), p)
		if err != nil {
			return err
		}
		switch format {
		case NodeChildrenFormatJson:
			return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
		}
		return nil
	}

	return cmd
}

type NodeOpenParams struct {
	Launch bool
	Force  bool
	Slug   string
}

type NodeOpenHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodeOpenParams) error

func newNodeOpenCommand(nodeOpen NodeOpenHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open <slug>",
		Short: "Print or launch a node's attached link URL.",
		Args:  rangeArgs(1, 1),
	}

	var rawLaunch bool
	cmd.Flags().BoolVar(&rawLaunch, "launch", false, "Open the URL in the system's default browser.")
	var rawForce bool
	cmd.Flags().BoolVar(&rawForce, "force", false, "Launch even when stdout is not a terminal.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawSlug := args[0]

		p := NodeOpenParams{
			Launch: rawLaunch,
			Force:  rawForce,
			Slug:   rawSlug,
		}

		return nodeOpen(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type NodeSearchFormat string

const (
	NodeSearchFormatAuto  NodeSearchFormat = "auto"
	NodeSearchFormatPlain NodeSearchFormat = "plain"
	NodeSearchFormatJson  NodeSearchFormat = "json"
	NodeSearchFormatJsonl NodeSearchFormat = "jsonl"
)

type NodeSearchOutput string

const (
	NodeSearchOutputDefault NodeSearchOutput = "default"
	NodeSearchOutputWide    NodeSearchOutput = "wide"
)

type NodeSearchParams struct {
	All    bool
	Page   int
	Limit  int
	Format NodeSearchFormat
	Output NodeSearchOutput
	Query  string
}

type NodeSearchHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodeSearchParams) (SearchNodeListResult, error)

func newNodeSearchCommand(nodeSearch NodeSearchHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Full-text search for nodes.",
		Args:  rangeArgs(1, 1),
	}

	var rawAll bool
	cmd.Flags().BoolVar(&rawAll, "all", false, "Fetch every page, streaming output as it goes.")
	var rawPage int
	cmd.Flags().IntVar(&rawPage, "page", 1, "Page to request.")
	var rawLimit int
	cmd.Flags().IntVar(&rawLimit, "limit", 0, "Stop after N matches (0 = no limit).")
	var rawFormat string
	cmd.Flags().StringVar(&rawFormat, "format", "auto", "Output format: auto, plain, json, jsonl.")
	var rawOutput string
	cmd.Flags().StringVarP(&rawOutput, "output", "o", "default", "Column profile: default, wide.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var format NodeSearchFormat
		if rawFormat != "" {
			switch rawFormat {
			case "auto":
				format = NodeSearchFormatAuto
			case "plain":
				format = NodeSearchFormatPlain
			case "json":
				format = NodeSearchFormatJson
			case "jsonl":
				format = NodeSearchFormatJsonl
			default:
				return fmt.Errorf("invalid --format %q: must be one of auto, plain, json, jsonl", rawFormat)
			}
		}
		var output NodeSearchOutput
		if rawOutput != "" {
			switch rawOutput {
			case "default":
				output = NodeSearchOutputDefault
			case "wide":
				output = NodeSearchOutputWide
			default:
				return fmt.Errorf("invalid --output %q: must be one of default, wide", rawOutput)
			}
		}
		rawQuery := args[0]

		p := NodeSearchParams{
			All:    rawAll,
			Page:   rawPage,
			Limit:  rawLimit,
			Format: format,
			Output: output,
			Query:  rawQuery,
		}

		result, err := nodeSearch(cmd.Context(), cmd, newIO(cmd), p)
		if err != nil {
			return err
		}
		switch format {
		case NodeSearchFormatJson:
			return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
		}
		return nil
	}

	return cmd
}

type NodeTreeLinkScheme string

const (
	NodeTreeLinkSchemeHttp  NodeTreeLinkScheme = "http"
	NodeTreeLinkSchemeHttps NodeTreeLinkScheme = "https"
)

type NodeTreeFormat string

const (
	NodeTreeFormatAuto  NodeTreeFormat = "auto"
	NodeTreeFormatPlain NodeTreeFormat = "plain"
	NodeTreeFormatJson  NodeTreeFormat = "json"
)

type NodeTreeParams struct {
	Depth           int
	Visibility      []string
	LinkDomain      []string
	LinkUrlContains string
	LinkScheme      NodeTreeLinkScheme
	HasLink         bool
	NoLink          bool
	RootOnly        bool
	OwnerHandle     string
	NameContains    string
	Format          NodeTreeFormat
	Slug            string
}

type NodeTreeHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodeTreeParams) (NodeListResult, error)

func newNodeTreeCommand(nodeTree NodeTreeHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tree [slug]",
		Short: "Show nodes as a tree.",
		Args:  rangeArgs(0, 1),
	}

	var rawDepth int
	cmd.Flags().IntVar(&rawDepth, "depth", 10, "Maximum depth to display.")
	var rawVisibility []string
	cmd.Flags().StringSliceVar(&rawVisibility, "visibility", nil, "Server-side filter by visibility (repeatable, comma-separated).")
	var rawLinkDomain []string
	cmd.Flags().StringSliceVar(&rawLinkDomain, "link-domain", nil, "Prune to nodes whose link.domain matches, keeping ancestors of any match (repeatable, comma-separated; client-side).")
	var rawLinkUrlContains string
	cmd.Flags().StringVar(&rawLinkUrlContains, "link-url-contains", "", "Prune to nodes whose link.url contains this substring (client-side).")
	var rawLinkScheme string
	cmd.Flags().StringVar(&rawLinkScheme, "link-scheme", "", "Prune to nodes whose link.url uses this scheme (client-side).")
	var rawHasLink bool
	cmd.Flags().BoolVar(&rawHasLink, "has-link", false, "Prune to nodes that have a link attached (client-side).")
	var rawNoLink bool
	cmd.Flags().BoolVar(&rawNoLink, "no-link", false, "Prune to nodes that have no link attached (client-side).")
	var rawRootOnly bool
	cmd.Flags().BoolVar(&rawRootOnly, "root-only", false, "Prune to nodes with no parent (client-side).")
	var rawOwnerHandle string
	cmd.Flags().StringVar(&rawOwnerHandle, "owner-handle", "", "Prune to nodes owned by this account handle (client-side).")
	var rawNameContains string
	cmd.Flags().StringVar(&rawNameContains, "name-contains", "", "Prune to nodes whose name contains this substring, case-insensitive (client-side).")
	var rawFormat string
	cmd.Flags().StringVar(&rawFormat, "format", "auto", "Output format.")
	cmd.MarkFlagsMutuallyExclusive("has-link", "no-link")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var linkScheme NodeTreeLinkScheme
		if rawLinkScheme != "" {
			switch rawLinkScheme {
			case "http":
				linkScheme = NodeTreeLinkSchemeHttp
			case "https":
				linkScheme = NodeTreeLinkSchemeHttps
			default:
				return fmt.Errorf("invalid --link-scheme %q: must be one of http, https", rawLinkScheme)
			}
		}
		var format NodeTreeFormat
		if rawFormat != "" {
			switch rawFormat {
			case "auto":
				format = NodeTreeFormatAuto
			case "plain":
				format = NodeTreeFormatPlain
			case "json":
				format = NodeTreeFormatJson
			default:
				return fmt.Errorf("invalid --format %q: must be one of auto, plain, json", rawFormat)
			}
		}
		var rawSlug string
		if len(args) > 0 {
			rawSlug = args[0]
		}

		p := NodeTreeParams{
			Depth:           rawDepth,
			Visibility:      rawVisibility,
			LinkDomain:      rawLinkDomain,
			LinkUrlContains: rawLinkUrlContains,
			LinkScheme:      linkScheme,
			HasLink:         rawHasLink,
			NoLink:          rawNoLink,
			RootOnly:        rawRootOnly,
			OwnerHandle:     rawOwnerHandle,
			NameContains:    rawNameContains,
			Format:          format,
			Slug:            rawSlug,
		}

		result, err := nodeTree(cmd.Context(), cmd, newIO(cmd), p)
		if err != nil {
			return err
		}
		switch format {
		case NodeTreeFormatJson:
			return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
		}
		return nil
	}

	return cmd
}

type NodeMetaGetParams struct {
	Slug string
}

type NodeMetaGetHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodeMetaGetParams) (Metadata, error)

func newNodeMetaGetCommand(nodeMetaGet NodeMetaGetHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <slug>",
		Short: "Get node metadata as JSON.",
		Args:  rangeArgs(1, 1),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawSlug := args[0]

		p := NodeMetaGetParams{
			Slug: rawSlug,
		}

		result, err := nodeMetaGet(cmd.Context(), cmd, newIO(cmd), p)
		if err != nil {
			return err
		}
		return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
	}

	return cmd
}

type NodeMetaSetParams struct {
	File string
	Slug string
	Json string
}

type NodeMetaSetHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodeMetaSetParams) error

func newNodeMetaSetCommand(nodeMetaSet NodeMetaSetHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <slug> [json]",
		Short: "Replace node metadata with a JSON object.",
		Long: `This is a full replace, not a merge. The object may be given
inline, via --file, or on stdin — supply exactly one of the
JSON argument or --file/stdin.
`,
		Args: rangeArgs(1, 2),
	}

	var rawFile string
	cmd.Flags().StringVar(&rawFile, "file", "", "Path to a JSON file containing the metadata object (use - for stdin).")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawSlug := args[0]
		var rawJson string
		if len(args) > 1 {
			rawJson = args[1]
		}

		p := NodeMetaSetParams{
			File: rawFile,
			Slug: rawSlug,
			Json: rawJson,
		}

		return nodeMetaSet(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

func newNodeMetaCommand(
	nodeMetaGet NodeMetaGetHandler,
	nodeMetaSet NodeMetaSetHandler,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "meta",
		Short: "Get or set raw node metadata.",
	}

	cmd.AddCommand(
		newNodeMetaGetCommand(nodeMetaGet),
		newNodeMetaSetCommand(nodeMetaSet),
	)

	return cmd
}

type NodeAssetsUploadParams struct {
	Name    string
	Primary bool
	Slug    string
	File    string
}

type NodeAssetsUploadHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodeAssetsUploadParams) error

func newNodeAssetsUploadCommand(nodeAssetsUpload NodeAssetsUploadHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload <slug> <file>",
		Short: "Upload a local file and attach it to a node.",
		Args:  rangeArgs(2, 2),
	}

	var rawName string
	cmd.Flags().StringVar(&rawName, "name", "", "Filename to store the asset as (defaults to the local file's basename).")
	var rawPrimary bool
	cmd.Flags().BoolVar(&rawPrimary, "primary", false, "Set the uploaded asset as the node's primary image.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawSlug := args[0]
		rawFile := args[1]

		p := NodeAssetsUploadParams{
			Name:    rawName,
			Primary: rawPrimary,
			Slug:    rawSlug,
			File:    rawFile,
		}

		return nodeAssetsUpload(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type NodeAssetsAddParams struct {
	Primary bool
	Slug    string
	AssetId string
}

type NodeAssetsAddHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodeAssetsAddParams) error

func newNodeAssetsAddCommand(nodeAssetsAdd NodeAssetsAddHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <slug> <asset-id>",
		Short: "Attach an already-uploaded asset to a node.",
		Args:  rangeArgs(2, 2),
	}

	var rawPrimary bool
	cmd.Flags().BoolVar(&rawPrimary, "primary", false, "Set the asset as the node's primary image.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawSlug := args[0]
		rawAssetId := args[1]

		p := NodeAssetsAddParams{
			Primary: rawPrimary,
			Slug:    rawSlug,
			AssetId: rawAssetId,
		}

		return nodeAssetsAdd(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type NodeAssetsRemoveParams struct {
	Slug    string
	AssetId string
}

type NodeAssetsRemoveHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodeAssetsRemoveParams) error

func newNodeAssetsRemoveCommand(nodeAssetsRemove NodeAssetsRemoveHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <slug> <asset-id>",
		Short: "Detach an asset from a node.",
		Args:  rangeArgs(2, 2),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawSlug := args[0]
		rawAssetId := args[1]

		p := NodeAssetsRemoveParams{
			Slug:    rawSlug,
			AssetId: rawAssetId,
		}

		return nodeAssetsRemove(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type NodeAssetsDownloadParams struct {
	OutputFile string
	Force      bool
	Slug       string
	Asset      string
}

type NodeAssetsDownloadHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodeAssetsDownloadParams) error

func newNodeAssetsDownloadCommand(nodeAssetsDownload NodeAssetsDownloadHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download <slug> <asset>",
		Short: "Download an asset attached to a node.",
		Args:  rangeArgs(2, 2),
	}

	var rawOutputFile string
	cmd.Flags().StringVar(&rawOutputFile, "output-file", "", "Output file path (defaults to a name derived from the source being written).")
	var rawForce bool
	cmd.Flags().BoolVar(&rawForce, "force", false, "Overwrite/proceed without confirmation.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawSlug := args[0]
		rawAsset := args[1]

		p := NodeAssetsDownloadParams{
			OutputFile: rawOutputFile,
			Force:      rawForce,
			Slug:       rawSlug,
			Asset:      rawAsset,
		}

		return nodeAssetsDownload(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type NodeAssetsPrimarySetParams struct {
	Slug    string
	AssetId string
}

type NodeAssetsPrimarySetHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodeAssetsPrimarySetParams) error

func newNodeAssetsPrimarySetCommand(nodeAssetsPrimarySet NodeAssetsPrimarySetHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <slug> <asset-id>",
		Short: "Set the node's primary image to an already-attached asset.",
		Args:  rangeArgs(2, 2),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawSlug := args[0]
		rawAssetId := args[1]

		p := NodeAssetsPrimarySetParams{
			Slug:    rawSlug,
			AssetId: rawAssetId,
		}

		return nodeAssetsPrimarySet(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type NodeAssetsPrimaryClearParams struct {
	Slug string
}

type NodeAssetsPrimaryClearHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodeAssetsPrimaryClearParams) error

func newNodeAssetsPrimaryClearCommand(nodeAssetsPrimaryClear NodeAssetsPrimaryClearHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear <slug>",
		Short: "Clear the node's primary image.",
		Args:  rangeArgs(1, 1),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawSlug := args[0]

		p := NodeAssetsPrimaryClearParams{
			Slug: rawSlug,
		}

		return nodeAssetsPrimaryClear(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type NodeAssetsPrimaryDownloadParams struct {
	OutputFile string
	Force      bool
	Slug       string
}

type NodeAssetsPrimaryDownloadHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodeAssetsPrimaryDownloadParams) error

func newNodeAssetsPrimaryDownloadCommand(nodeAssetsPrimaryDownload NodeAssetsPrimaryDownloadHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download <slug>",
		Short: "Download the node's primary image.",
		Args:  rangeArgs(1, 1),
	}

	var rawOutputFile string
	cmd.Flags().StringVar(&rawOutputFile, "output-file", "", "Output file path (defaults to a name derived from the source being written).")
	var rawForce bool
	cmd.Flags().BoolVar(&rawForce, "force", false, "Overwrite/proceed without confirmation.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawSlug := args[0]

		p := NodeAssetsPrimaryDownloadParams{
			OutputFile: rawOutputFile,
			Force:      rawForce,
			Slug:       rawSlug,
		}

		return nodeAssetsPrimaryDownload(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

func newNodeAssetsPrimaryCommand(
	nodeAssetsPrimarySet NodeAssetsPrimarySetHandler,
	nodeAssetsPrimaryClear NodeAssetsPrimaryClearHandler,
	nodeAssetsPrimaryDownload NodeAssetsPrimaryDownloadHandler,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "primary",
		Short: "Set, clear, or download the node's primary image.",
	}

	cmd.AddCommand(
		newNodeAssetsPrimarySetCommand(nodeAssetsPrimarySet),
		newNodeAssetsPrimaryClearCommand(nodeAssetsPrimaryClear),
		newNodeAssetsPrimaryDownloadCommand(nodeAssetsPrimaryDownload),
	)

	return cmd
}

func newNodeAssetsCommand(
	nodeAssetsUpload NodeAssetsUploadHandler,
	nodeAssetsAdd NodeAssetsAddHandler,
	nodeAssetsRemove NodeAssetsRemoveHandler,
	nodeAssetsDownload NodeAssetsDownloadHandler,
	nodeAssetsPrimarySet NodeAssetsPrimarySetHandler,
	nodeAssetsPrimaryClear NodeAssetsPrimaryClearHandler,
	nodeAssetsPrimaryDownload NodeAssetsPrimaryDownloadHandler,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assets",
		Short: "Upload, download, attach, and remove node assets.",
	}

	cmd.AddCommand(
		newNodeAssetsUploadCommand(nodeAssetsUpload),
		newNodeAssetsAddCommand(nodeAssetsAdd),
		newNodeAssetsRemoveCommand(nodeAssetsRemove),
		newNodeAssetsDownloadCommand(nodeAssetsDownload),
		newNodeAssetsPrimaryCommand(
			nodeAssetsPrimarySet,
			nodeAssetsPrimaryClear,
			nodeAssetsPrimaryDownload,
		),
	)

	return cmd
}

type NodePropertiesGetFormat string

const (
	NodePropertiesGetFormatJson NodePropertiesGetFormat = "json"
	NodePropertiesGetFormatYaml NodePropertiesGetFormat = "yaml"
)

type NodePropertiesGetParams struct {
	Format NodePropertiesGetFormat
	Slug   string
}

type NodePropertiesGetHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodePropertiesGetParams) (PropertyList, error)

func newNodePropertiesGetCommand(nodePropertiesGet NodePropertiesGetHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <slug>",
		Short: "Get property values for a node.",
		Args:  rangeArgs(1, 1),
	}

	var rawFormat string
	cmd.Flags().StringVarP(&rawFormat, "format", "f", "json", "Output format.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var format NodePropertiesGetFormat
		if rawFormat != "" {
			switch rawFormat {
			case "json":
				format = NodePropertiesGetFormatJson
			case "yaml":
				format = NodePropertiesGetFormatYaml
			default:
				return fmt.Errorf("invalid --format %q: must be one of json, yaml", rawFormat)
			}
		}
		rawSlug := args[0]

		p := NodePropertiesGetParams{
			Format: format,
			Slug:   rawSlug,
		}

		result, err := nodePropertiesGet(cmd.Context(), cmd, newIO(cmd), p)
		if err != nil {
			return err
		}
		switch format {
		case NodePropertiesGetFormatJson:
			return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
		case NodePropertiesGetFormatYaml:
			enc := yaml.NewEncoder(cmd.OutOrStdout())
			if err := enc.Encode(result); err != nil {
				return err
			}
			return enc.Close()
		}
		return nil
	}

	return cmd
}

type NodePropertiesSetParams struct {
	Slug  string
	Token []string
}

type NodePropertiesSetHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodePropertiesSetParams) error

func newNodePropertiesSetCommand(nodePropertiesSet NodePropertiesSetHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <slug> <token>...",
		Short: "Set property values for a node.",
		Long:  "Each TOKEN is `name=value` to set an existing property (its\ntype is looked up automatically), or `name:type=value` to\ncreate a new property (type is one of text, number, boolean,\ntimestamp).\n",
		Args:  rangeArgs(2, -1),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawSlug := args[0]
		token := args[1:]

		p := NodePropertiesSetParams{
			Slug:  rawSlug,
			Token: token,
		}

		return nodePropertiesSet(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type NodePropertiesSchemaGetFormat string

const (
	NodePropertiesSchemaGetFormatJson NodePropertiesSchemaGetFormat = "json"
	NodePropertiesSchemaGetFormatYaml NodePropertiesSchemaGetFormat = "yaml"
)

type NodePropertiesSchemaGetParams struct {
	Format NodePropertiesSchemaGetFormat
	Slug   string
}

type NodePropertiesSchemaGetHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodePropertiesSchemaGetParams) (PropertySchemaList, error)

func newNodePropertiesSchemaGetCommand(nodePropertiesSchemaGet NodePropertiesSchemaGetHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <slug>",
		Short: "Get the property schema shared by a node and its siblings.",
		Long:  "Reads `<slug>`'s parent's child-property-schema — i.e. the\nschema that applies to `<slug>` and every other child of\nthe same parent.\n",
		Args:  rangeArgs(1, 1),
	}

	var rawFormat string
	cmd.Flags().StringVarP(&rawFormat, "format", "f", "json", "Output format.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var format NodePropertiesSchemaGetFormat
		if rawFormat != "" {
			switch rawFormat {
			case "json":
				format = NodePropertiesSchemaGetFormatJson
			case "yaml":
				format = NodePropertiesSchemaGetFormatYaml
			default:
				return fmt.Errorf("invalid --format %q: must be one of json, yaml", rawFormat)
			}
		}
		rawSlug := args[0]

		p := NodePropertiesSchemaGetParams{
			Format: format,
			Slug:   rawSlug,
		}

		result, err := nodePropertiesSchemaGet(cmd.Context(), cmd, newIO(cmd), p)
		if err != nil {
			return err
		}
		switch format {
		case NodePropertiesSchemaGetFormatJson:
			return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
		case NodePropertiesSchemaGetFormatYaml:
			enc := yaml.NewEncoder(cmd.OutOrStdout())
			if err := enc.Encode(result); err != nil {
				return err
			}
			return enc.Close()
		}
		return nil
	}

	return cmd
}

type NodePropertiesSchemaSetParams struct {
	Slug  string
	Token []string
}

type NodePropertiesSchemaSetHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodePropertiesSchemaSetParams) error

func newNodePropertiesSchemaSetCommand(nodePropertiesSchemaSet NodePropertiesSchemaSetHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <slug> <token>...",
		Short: "Set the property schema shared by a node and its siblings.",
		Long:  "Sets `<slug>`'s parent's child-property-schema — i.e. the\nschema that will then apply to `<slug>` and every other\nchild of the same parent. Each TOKEN is\n`field:type:sort`, where type is one of text, number,\nboolean, timestamp and sort is asc or desc.\n",
		Args:  rangeArgs(2, -1),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawSlug := args[0]
		token := args[1:]

		p := NodePropertiesSchemaSetParams{
			Slug:  rawSlug,
			Token: token,
		}

		return nodePropertiesSchemaSet(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type NodePropertiesSchemaChildrenParams struct {
	Slug  string
	Token []string
}

type NodePropertiesSchemaChildrenHandler func(ctx context.Context, cmd *cobra.Command, io IO, p NodePropertiesSchemaChildrenParams) error

func newNodePropertiesSchemaChildrenCommand(nodePropertiesSchemaChildren NodePropertiesSchemaChildrenHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "children <slug> <token>...",
		Short: "Set the property schema for a node's children.",
		Long:  "The complement of `schema set`: sets the schema that will\napply to `<slug>`'s own children (rather than its\nsiblings). Each TOKEN is `field:type:sort`, same syntax as\n`schema set`.\n",
		Args:  rangeArgs(2, -1),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawSlug := args[0]
		token := args[1:]

		p := NodePropertiesSchemaChildrenParams{
			Slug:  rawSlug,
			Token: token,
		}

		return nodePropertiesSchemaChildren(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

func newNodePropertiesSchemaCommand(
	nodePropertiesSchemaGet NodePropertiesSchemaGetHandler,
	nodePropertiesSchemaSet NodePropertiesSchemaSetHandler,
	nodePropertiesSchemaChildren NodePropertiesSchemaChildrenHandler,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Manage node property schemas.",
	}

	cmd.AddCommand(
		newNodePropertiesSchemaGetCommand(nodePropertiesSchemaGet),
		newNodePropertiesSchemaSetCommand(nodePropertiesSchemaSet),
		newNodePropertiesSchemaChildrenCommand(nodePropertiesSchemaChildren),
	)

	return cmd
}

func newNodePropertiesCommand(
	nodePropertiesGet NodePropertiesGetHandler,
	nodePropertiesSet NodePropertiesSetHandler,
	nodePropertiesSchemaGet NodePropertiesSchemaGetHandler,
	nodePropertiesSchemaSet NodePropertiesSchemaSetHandler,
	nodePropertiesSchemaChildren NodePropertiesSchemaChildrenHandler,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "properties",
		Short: "Manage node properties.",
	}

	cmd.AddCommand(
		newNodePropertiesGetCommand(nodePropertiesGet),
		newNodePropertiesSetCommand(nodePropertiesSet),
		newNodePropertiesSchemaCommand(
			nodePropertiesSchemaGet,
			nodePropertiesSchemaSet,
			nodePropertiesSchemaChildren,
		),
	)

	return cmd
}

func NewNodeCommand(
	nodeList NodeListHandler,
	nodeGet NodeGetHandler,
	nodeCreate NodeCreateHandler,
	nodeUpdate NodeUpdateHandler,
	nodeDelete NodeDeleteHandler,
	nodeMove NodeMoveHandler,
	nodeVisibility NodeVisibilityHandler,
	nodeChildren NodeChildrenHandler,
	nodeOpen NodeOpenHandler,
	nodeSearch NodeSearchHandler,
	nodeTree NodeTreeHandler,
	nodeMetaGet NodeMetaGetHandler,
	nodeMetaSet NodeMetaSetHandler,
	nodeAssetsUpload NodeAssetsUploadHandler,
	nodeAssetsAdd NodeAssetsAddHandler,
	nodeAssetsRemove NodeAssetsRemoveHandler,
	nodeAssetsDownload NodeAssetsDownloadHandler,
	nodeAssetsPrimarySet NodeAssetsPrimarySetHandler,
	nodeAssetsPrimaryClear NodeAssetsPrimaryClearHandler,
	nodeAssetsPrimaryDownload NodeAssetsPrimaryDownloadHandler,
	nodePropertiesGet NodePropertiesGetHandler,
	nodePropertiesSet NodePropertiesSetHandler,
	nodePropertiesSchemaGet NodePropertiesSchemaGetHandler,
	nodePropertiesSchemaSet NodePropertiesSchemaSetHandler,
	nodePropertiesSchemaChildren NodePropertiesSchemaChildrenHandler,
) NodeCommand {
	cmd := &cobra.Command{
		Use:     "node",
		Short:   "Work with Storyden nodes (pages).",
		Aliases: []string{"nodes", "n"},
	}

	cmd.AddCommand(
		newNodeListCommand(nodeList),
		newNodeGetCommand(nodeGet),
		newNodeCreateCommand(nodeCreate),
		newNodeUpdateCommand(nodeUpdate),
		newNodeDeleteCommand(nodeDelete),
		newNodeMoveCommand(nodeMove),
		newNodeVisibilityCommand(nodeVisibility),
		newNodeChildrenCommand(nodeChildren),
		newNodeOpenCommand(nodeOpen),
		newNodeSearchCommand(nodeSearch),
		newNodeTreeCommand(nodeTree),
		newNodeMetaCommand(
			nodeMetaGet,
			nodeMetaSet,
		),
		newNodeAssetsCommand(
			nodeAssetsUpload,
			nodeAssetsAdd,
			nodeAssetsRemove,
			nodeAssetsDownload,
			nodeAssetsPrimarySet,
			nodeAssetsPrimaryClear,
			nodeAssetsPrimaryDownload,
		),
		newNodePropertiesCommand(
			nodePropertiesGet,
			nodePropertiesSet,
			nodePropertiesSchemaGet,
			nodePropertiesSchemaSet,
			nodePropertiesSchemaChildren,
		),
	)

	return NodeCommand(cmd)
}

type NodeCommand *cobra.Command
