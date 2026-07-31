package cligen

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

type InfoMetadataParams struct {
}

type InfoMetadataHandler func(ctx context.Context, cmd *cobra.Command, io IO, p InfoMetadataParams) (Metadata, error)

func newInfoMetadataCommand(infoMetadata InfoMetadataHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metadata",
		Short: "Show raw instance metadata as JSON.",
		Args:  rangeArgs(0, 0),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		p := InfoMetadataParams{}

		result, err := infoMetadata(cmd.Context(), cmd, newIO(cmd), p)
		if err != nil {
			return err
		}
		return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
	}

	return cmd
}

type InfoFormat string

const (
	InfoFormatPlain InfoFormat = "plain"
	InfoFormatJson  InfoFormat = "json"
)

type InfoParams struct {
	Format InfoFormat
}

type InfoHandler func(ctx context.Context, cmd *cobra.Command, io IO, p InfoParams) (InstanceInfo, error)

func NewInfoCommand(
	info InfoHandler,
	infoMetadata InfoMetadataHandler,
) InfoCommand {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Show basic information about the current Storyden instance.",
		Args:  rangeArgs(0, 0),
	}

	var rawFormat string
	cmd.Flags().StringVar(&rawFormat, "format", "plain", "Output format.")

	cmd.AddCommand(
		newInfoMetadataCommand(infoMetadata),
	)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var format InfoFormat
		if rawFormat != "" {
			switch rawFormat {
			case "plain":
				format = InfoFormatPlain
			case "json":
				format = InfoFormatJson
			default:
				return fmt.Errorf("invalid --format %q: must be one of plain, json", rawFormat)
			}
		}

		p := InfoParams{
			Format: format,
		}

		result, err := info(cmd.Context(), cmd, newIO(cmd), p)
		if err != nil {
			return err
		}
		switch format {
		case InfoFormatJson:
			return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
		}
		return nil
	}

	return InfoCommand(cmd)
}

type InfoCommand *cobra.Command
