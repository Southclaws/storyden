package symbols

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	plugindev "github.com/Southclaws/storyden/lib/plugin/dev"
)

func NewPackages() cligen.PluginDevSymbolsPackagesHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.PluginDevSymbolsPackagesParams) error {
		result, err := plugindev.ListGoPackages(ctx, p.Dir, plugindev.PackageListOptions{
			Pattern:     p.Pattern,
			IncludeDeps: p.IncludeDeps,
			MaxPackages: p.Limit,
		})
		if err != nil {
			return err
		}
		return writeJSON(cio.Out, result)
	}
}

func NewPackage() cligen.PluginDevSymbolsPackageHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.PluginDevSymbolsPackageParams) error {
		result, err := plugindev.GoPackageSymbols(ctx, p.Dir, plugindev.PackageSymbolsOptions{
			ImportPath:        p.ImportPath,
			IncludeUnexported: p.IncludeUnexported,
			MaxSymbols:        p.Limit,
		})
		if err != nil {
			return err
		}
		return writeJSON(cio.Out, result)
	}
}

func NewDetail() cligen.PluginDevSymbolsDetailHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.PluginDevSymbolsDetailParams) error {
		result, err := plugindev.GoSymbolDetail(ctx, p.Dir, plugindev.SymbolDetailOptions{
			ImportPath: p.ImportPath,
			Symbol:     p.Symbol,
		})
		if err != nil {
			return err
		}
		return writeJSON(cio.Out, result)
	}
}

func NewSearch() cligen.PluginDevSymbolsSearchHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.PluginDevSymbolsSearchParams) error {
		result, err := plugindev.GoSymbolSearch(ctx, p.Dir, plugindev.SymbolSearchOptions{
			Query:             p.Query,
			Pattern:           p.Pattern,
			IncludeDeps:       p.IncludeDeps,
			IncludeUnexported: p.IncludeUnexported,
			MaxResults:        p.Limit,
		})
		if err != nil {
			return err
		}
		return writeJSON(cio.Out, result)
	}
}

func writeJSON(out io.Writer, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, string(data))
	return err
}
