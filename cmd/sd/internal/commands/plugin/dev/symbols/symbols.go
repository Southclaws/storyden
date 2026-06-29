package symbols

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	plugindev "github.com/Southclaws/storyden/lib/plugin/dev"
)

type SymbolsCommand *cobra.Command

func New() SymbolsCommand {
	var dir string

	command := &cobra.Command{
		Use:   "symbols",
		Short: "Discover Go packages and symbols in a plugin project",
		Long: `# Plugin Symbol Discovery

Inspect Go packages, symbols, fields, methods, and docs using the same discovery engine used by the plugin builder Robot tools.
`,
	}

	command.PersistentFlags().StringVar(&dir, "dir", ".", "Plugin project directory")
	command.AddCommand(newPackagesCommand(&dir))
	command.AddCommand(newPackageCommand(&dir))
	command.AddCommand(newDetailCommand(&dir))
	command.AddCommand(newSearchCommand(&dir))

	help.SetupMarkdownHelp(command)

	return SymbolsCommand(command)
}

func newPackagesCommand(dir *string) *cobra.Command {
	var pattern string
	var includeDeps bool
	var maxPackages int

	command := &cobra.Command{
		Use:   "packages",
		Short: "List Go packages",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := plugindev.ListGoPackages(cmd.Context(), *dir, plugindev.PackageListOptions{
				Pattern:     pattern,
				IncludeDeps: includeDeps,
				MaxPackages: maxPackages,
			})
			if err != nil {
				return err
			}
			return writeJSON(cmd, result)
		},
	}

	command.Flags().StringVar(&pattern, "pattern", "./...", "Go package pattern")
	command.Flags().BoolVar(&includeDeps, "include-deps", false, "Include transitive dependency packages")
	command.Flags().IntVar(&maxPackages, "max", 100, "Maximum packages to return")

	return command
}

func newPackageCommand(dir *string) *cobra.Command {
	var includeUnexported bool
	var maxSymbols int

	command := &cobra.Command{
		Use:   "package <import-path>",
		Short: "List symbols in a Go package",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := plugindev.GoPackageSymbols(cmd.Context(), *dir, plugindev.PackageSymbolsOptions{
				ImportPath:        args[0],
				IncludeUnexported: includeUnexported,
				MaxSymbols:        maxSymbols,
			})
			if err != nil {
				return err
			}
			return writeJSON(cmd, result)
		},
	}

	command.Flags().BoolVar(&includeUnexported, "include-unexported", false, "Include unexported package symbols")
	command.Flags().IntVar(&maxSymbols, "max", 100, "Maximum symbols to return")

	return command
}

func newDetailCommand(dir *string) *cobra.Command {
	command := &cobra.Command{
		Use:   "detail <import-path> <symbol>",
		Short: "Inspect one Go symbol",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := plugindev.GoSymbolDetail(cmd.Context(), *dir, plugindev.SymbolDetailOptions{
				ImportPath: args[0],
				Symbol:     args[1],
			})
			if err != nil {
				return err
			}
			return writeJSON(cmd, result)
		},
	}

	return command
}

func newSearchCommand(dir *string) *cobra.Command {
	var pattern string
	var includeDeps bool
	var includeUnexported bool
	var maxResults int

	command := &cobra.Command{
		Use:   "search <query>",
		Short: "Search Go symbols by name, signature, or docs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := plugindev.GoSymbolSearch(cmd.Context(), *dir, plugindev.SymbolSearchOptions{
				Query:             args[0],
				Pattern:           pattern,
				IncludeDeps:       includeDeps,
				IncludeUnexported: includeUnexported,
				MaxResults:        maxResults,
			})
			if err != nil {
				return err
			}
			return writeJSON(cmd, result)
		},
	}

	command.Flags().StringVar(&pattern, "pattern", "./...", "Go package pattern")
	command.Flags().BoolVar(&includeDeps, "include-deps", false, "Include transitive dependency packages")
	command.Flags().BoolVar(&includeUnexported, "include-unexported", false, "Include unexported package symbols")
	command.Flags().IntVar(&maxResults, "max", 100, "Maximum symbol matches to return")

	return command
}

func writeJSON(cmd *cobra.Command, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(cmd.OutOrStdout(), string(data))
	return err
}
