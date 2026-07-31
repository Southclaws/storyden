package cligen

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

type ThreadGetFormat string

const (
	ThreadGetFormatJson     ThreadGetFormat = "json"
	ThreadGetFormatMarkdown ThreadGetFormat = "markdown"
	ThreadGetFormatYaml     ThreadGetFormat = "yaml"
)

type ThreadGetParams struct {
	Format     ThreadGetFormat
	ThreadMark string
}

type ThreadGetHandler func(ctx context.Context, cmd *cobra.Command, io IO, p ThreadGetParams) (Thread, error)

func newThreadGetCommand(threadGet ThreadGetHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <thread-mark>",
		Short: "Get a thread by its mark.",
		Example: `  sd thread get my-thread --format json
  sd thread get my-thread --format yaml`,
		Args: rangeArgs(1, 1),
	}

	var rawFormat string
	cmd.Flags().StringVarP(&rawFormat, "format", "f", "markdown", "Output format.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var format ThreadGetFormat
		if rawFormat != "" {
			switch rawFormat {
			case "json":
				format = ThreadGetFormatJson
			case "markdown":
				format = ThreadGetFormatMarkdown
			case "yaml":
				format = ThreadGetFormatYaml
			default:
				return fmt.Errorf("invalid --format %q: must be one of json, markdown, yaml", rawFormat)
			}
		}
		rawThreadMark := args[0]

		p := ThreadGetParams{
			Format:     format,
			ThreadMark: rawThreadMark,
		}

		result, err := threadGet(cmd.Context(), cmd, newIO(cmd), p)
		if err != nil {
			return err
		}
		switch format {
		case ThreadGetFormatJson:
			return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
		case ThreadGetFormatYaml:
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

type ThreadListFormat string

const (
	ThreadListFormatAuto  ThreadListFormat = "auto"
	ThreadListFormatPlain ThreadListFormat = "plain"
	ThreadListFormatJson  ThreadListFormat = "json"
	ThreadListFormatJsonl ThreadListFormat = "jsonl"
)

type ThreadListOutput string

const (
	ThreadListOutputDefault ThreadListOutput = "default"
	ThreadListOutputWide    ThreadListOutput = "wide"
)

type ThreadListParams struct {
	All    bool
	Page   int
	Limit  int
	Format ThreadListFormat
	Output ThreadListOutput
}

type ThreadListHandler func(ctx context.Context, cmd *cobra.Command, io IO, p ThreadListParams) (ThreadListResult, error)

func newThreadListCommand(threadList ThreadListHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List recent Storyden threads.",
		Example: `  sd thread list
  sd thread list --all --format jsonl`,
		Args: rangeArgs(0, 0),
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
		var format ThreadListFormat
		if rawFormat != "" {
			switch rawFormat {
			case "auto":
				format = ThreadListFormatAuto
			case "plain":
				format = ThreadListFormatPlain
			case "json":
				format = ThreadListFormatJson
			case "jsonl":
				format = ThreadListFormatJsonl
			default:
				return fmt.Errorf("invalid --format %q: must be one of auto, plain, json, jsonl", rawFormat)
			}
		}
		var output ThreadListOutput
		if rawOutput != "" {
			switch rawOutput {
			case "default":
				output = ThreadListOutputDefault
			case "wide":
				output = ThreadListOutputWide
			default:
				return fmt.Errorf("invalid --output %q: must be one of default, wide", rawOutput)
			}
		}

		p := ThreadListParams{
			All:    rawAll,
			Page:   rawPage,
			Limit:  rawLimit,
			Format: format,
			Output: output,
		}

		result, err := threadList(cmd.Context(), cmd, newIO(cmd), p)
		if err != nil {
			return err
		}
		switch format {
		case ThreadListFormatJson:
			return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
		}
		return nil
	}

	return cmd
}

func NewThreadCommand(
	threadGet ThreadGetHandler,
	threadList ThreadListHandler,
) ThreadCommand {
	cmd := &cobra.Command{
		Use:   "thread",
		Short: "Work with Storyden threads.",
	}

	cmd.AddCommand(
		newThreadGetCommand(threadGet),
		newThreadListCommand(threadList),
	)

	return ThreadCommand(cmd)
}

type ThreadCommand *cobra.Command
