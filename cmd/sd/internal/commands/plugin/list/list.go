package list

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
	"github.com/Southclaws/storyden/cmd/sd/internal/render"
	plugindev "github.com/Southclaws/storyden/lib/plugin/dev"
)

type ListCommand *cobra.Command

func New(store *config.Store) ListCommand {
	var format string
	var wide bool

	command := &cobra.Command{
		Use:   "list",
		Short: "List plugins installed on the current instance",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateFormat(format); err != nil {
				return err
			}

			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}
			plugins, err := plugindev.ListPlugins(cmd.Context(), client.OpenAPI)
			if err != nil {
				return err
			}
			if format == "json" {
				return output.JSON(cmd.OutOrStdout(), plugins)
			}
			return render.Render(cmd.OutOrStdout(), plugins, profile(), wide, render.PageInfo{})
		},
	}
	command.Flags().StringVar(&format, "format", "plain", "Output format: plain or json")
	command.Flags().BoolVar(&wide, "wide", false, "Show additional columns")
	help.SetupMarkdownHelp(command)
	return ListCommand(command)
}

func validateFormat(format string) error {
	if format != "" && format != "plain" && format != "json" {
		return fmt.Errorf("--format must be plain or json")
	}
	return nil
}

func profile() render.Profile[openapi.Plugin] {
	return render.Profile[openapi.Plugin]{Columns: []render.Column[openapi.Plugin]{
		{Header: "ID", Render: func(p openapi.Plugin) string { return string(p.Id) }},
		{Header: "NAME", Render: func(p openapi.Plugin) string { return p.Name }},
		{Header: "MODE", Render: plugindev.PluginMode},
		{Header: "STATUS", Render: plugindev.PluginStatus},
		{Header: "VERSION", Render: func(p openapi.Plugin) string {
			if p.Version == nil {
				return ""
			}
			return *p.Version
		}, Wide: true},
		{Header: "DESCRIPTION", Render: func(p openapi.Plugin) string {
			if p.Description == nil {
				return ""
			}
			return *p.Description
		}, Wide: true},
	}}
}
