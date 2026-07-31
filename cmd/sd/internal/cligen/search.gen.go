package cligen

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

type SearchFormat string

const (
	SearchFormatAuto  SearchFormat = "auto"
	SearchFormatPlain SearchFormat = "plain"
	SearchFormatJson  SearchFormat = "json"
	SearchFormatJsonl SearchFormat = "jsonl"
)

type SearchOutput string

const (
	SearchOutputDefault SearchOutput = "default"
	SearchOutputWide    SearchOutput = "wide"
)

type SearchParams struct {
	Kind       []string
	Authors    []string
	Categories []string
	Tags       []string
	Page       int
	Limit      int
	All        bool
	Format     SearchFormat
	Output     SearchOutput
	Query      string
}

type SearchHandler func(ctx context.Context, cmd *cobra.Command, io IO, p SearchParams) (SearchResult, error)

func NewSearchCommand(search SearchHandler) SearchCommand {
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search all Storyden content.",
		Example: `  sd search "design system"
  sd search "release notes" --kind thread,post
  sd search widget --authors southclaws --categories hardware
  sd search widget --all --format jsonl`,
		Args: rangeArgs(1, 1),
	}

	var rawKind []string
	cmd.Flags().StringSliceVar(&rawKind, "kind", nil, "Filter by datagraph item kind (repeatable, comma-separated): post, thread, reply, node, collection, profile, event.")
	var rawAuthors []string
	cmd.Flags().StringSliceVar(&rawAuthors, "authors", nil, "Filter by author account IDs or handles (repeatable, comma-separated).")
	var rawCategories []string
	cmd.Flags().StringSliceVar(&rawCategories, "categories", nil, "Filter by category slugs (repeatable, comma-separated).")
	var rawTags []string
	cmd.Flags().StringSliceVar(&rawTags, "tags", nil, "Filter by tag names (repeatable, comma-separated).")
	var rawPage int
	cmd.Flags().IntVar(&rawPage, "page", 1, "Page to request.")
	var rawLimit int
	cmd.Flags().IntVar(&rawLimit, "limit", 0, "Stop after N matches (0 = no limit).")
	var rawAll bool
	cmd.Flags().BoolVar(&rawAll, "all", false, "Fetch every page, streaming output as it goes.")
	var rawFormat string
	cmd.Flags().StringVar(&rawFormat, "format", "auto", "Output format: auto, plain, json, jsonl.")
	var rawOutput string
	cmd.Flags().StringVarP(&rawOutput, "output", "o", "default", "Column profile: default, wide.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var format SearchFormat
		if rawFormat != "" {
			switch rawFormat {
			case "auto":
				format = SearchFormatAuto
			case "plain":
				format = SearchFormatPlain
			case "json":
				format = SearchFormatJson
			case "jsonl":
				format = SearchFormatJsonl
			default:
				return fmt.Errorf("invalid --format %q: must be one of auto, plain, json, jsonl", rawFormat)
			}
		}
		var output SearchOutput
		if rawOutput != "" {
			switch rawOutput {
			case "default":
				output = SearchOutputDefault
			case "wide":
				output = SearchOutputWide
			default:
				return fmt.Errorf("invalid --output %q: must be one of default, wide", rawOutput)
			}
		}
		rawQuery := args[0]

		p := SearchParams{
			Kind:       rawKind,
			Authors:    rawAuthors,
			Categories: rawCategories,
			Tags:       rawTags,
			Page:       rawPage,
			Limit:      rawLimit,
			All:        rawAll,
			Format:     format,
			Output:     output,
			Query:      rawQuery,
		}

		result, err := search(cmd.Context(), cmd, newIO(cmd), p)
		if err != nil {
			return err
		}
		switch format {
		case SearchFormatJson:
			return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
		}
		return nil
	}

	return SearchCommand(cmd)
}

type SearchCommand *cobra.Command
