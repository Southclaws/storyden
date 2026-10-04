package list

import (
	"context"
	"encoding/json"

	"github.com/Southclaws/storyden/cmd/sd/internal/output"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/render"
	plugindev "github.com/Southclaws/storyden/lib/plugin/dev"
)

func New(store *config.Store) cligen.PluginListHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.PluginListParams) error {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}
		plugins, err := plugindev.ListPlugins(ctx, client.OpenAPI)
		if err != nil {
			return err
		}

		wide := p.Wide

		switch p.Output {
		case cligen.PluginListOutputJson:
			return output.JSON(io.Out, plugins)
		case cligen.PluginListOutputJsonl:
			encoder := json.NewEncoder(io.Out)
			for _, plugin := range plugins {
				if err := encoder.Encode(plugin); err != nil {
					return err
				}
			}
			return nil
		default:
			return render.Render(io.Out, plugins, profile(), wide, render.PageInfo{})
		}
	}
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
