package cligen

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

type PluginGetParams struct {
	PluginInstanceId string
}

type PluginGetHandler func(ctx context.Context, cmd *cobra.Command, io IO, p PluginGetParams) (Plugin, error)

func newPluginGetCommand(pluginGet PluginGetHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get <plugin-instance-id>",
		Short:   "Get plugin status and manifest information.",
		Aliases: []string{"status"},
		Args:    rangeArgs(1, 1),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawPluginInstanceId := args[0]

		p := PluginGetParams{
			PluginInstanceId: rawPluginInstanceId,
		}

		result, err := pluginGet(cmd.Context(), cmd, newIO(cmd), p)
		if err != nil {
			return err
		}
		return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
	}

	return cmd
}

type PluginListFormat string

const (
	PluginListFormatAuto  PluginListFormat = "auto"
	PluginListFormatPlain PluginListFormat = "plain"
	PluginListFormatJson  PluginListFormat = "json"
	PluginListFormatJsonl PluginListFormat = "jsonl"
)

type PluginListOutput string

const (
	PluginListOutputDefault PluginListOutput = "default"
	PluginListOutputWide    PluginListOutput = "wide"
)

type PluginListParams struct {
	Format PluginListFormat
	Output PluginListOutput
}

type PluginListHandler func(ctx context.Context, cmd *cobra.Command, io IO, p PluginListParams) ([]Plugin, error)

func newPluginListCommand(pluginList PluginListHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List plugins installed on the current instance.",
		Args:  rangeArgs(0, 0),
	}

	var rawFormat string
	cmd.Flags().StringVar(&rawFormat, "format", "auto", "Output format: auto, plain, json, jsonl.")
	var rawOutput string
	cmd.Flags().StringVarP(&rawOutput, "output", "o", "default", "Column profile: default, wide.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var format PluginListFormat
		if rawFormat != "" {
			switch rawFormat {
			case "auto":
				format = PluginListFormatAuto
			case "plain":
				format = PluginListFormatPlain
			case "json":
				format = PluginListFormatJson
			case "jsonl":
				format = PluginListFormatJsonl
			default:
				return fmt.Errorf("invalid --format %q: must be one of auto, plain, json, jsonl", rawFormat)
			}
		}
		var output PluginListOutput
		if rawOutput != "" {
			switch rawOutput {
			case "default":
				output = PluginListOutputDefault
			case "wide":
				output = PluginListOutputWide
			default:
				return fmt.Errorf("invalid --output %q: must be one of default, wide", rawOutput)
			}
		}

		p := PluginListParams{
			Format: format,
			Output: output,
		}

		result, err := pluginList(cmd.Context(), cmd, newIO(cmd), p)
		if err != nil {
			return err
		}
		switch format {
		case PluginListFormatJson:
			return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
		}
		return nil
	}

	return cmd
}

type PluginDeleteParams struct {
	PluginInstanceId string
}

type PluginDeleteHandler func(ctx context.Context, cmd *cobra.Command, io IO, p PluginDeleteParams) error

func newPluginDeleteCommand(pluginDelete PluginDeleteHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <plugin-instance-id>",
		Short: "Delete a plugin from the current instance.",
		Args:  rangeArgs(1, 1),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawPluginInstanceId := args[0]

		p := PluginDeleteParams{
			PluginInstanceId: rawPluginInstanceId,
		}

		return pluginDelete(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type PluginActivateParams struct {
	PluginInstanceId string
}

type PluginActivateHandler func(ctx context.Context, cmd *cobra.Command, io IO, p PluginActivateParams) (Plugin, error)

func newPluginActivateCommand(pluginActivate PluginActivateHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activate <plugin-instance-id>",
		Short: "Start a supervised plugin.",
		Args:  rangeArgs(1, 1),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawPluginInstanceId := args[0]

		p := PluginActivateParams{
			PluginInstanceId: rawPluginInstanceId,
		}

		result, err := pluginActivate(cmd.Context(), cmd, newIO(cmd), p)
		if err != nil {
			return err
		}
		return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
	}

	return cmd
}

type PluginDeactivateParams struct {
	PluginInstanceId string
}

type PluginDeactivateHandler func(ctx context.Context, cmd *cobra.Command, io IO, p PluginDeactivateParams) (Plugin, error)

func newPluginDeactivateCommand(pluginDeactivate PluginDeactivateHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deactivate <plugin-instance-id>",
		Short: "Stop a supervised plugin.",
		Args:  rangeArgs(1, 1),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawPluginInstanceId := args[0]

		p := PluginDeactivateParams{
			PluginInstanceId: rawPluginInstanceId,
		}

		result, err := pluginDeactivate(cmd.Context(), cmd, newIO(cmd), p)
		if err != nil {
			return err
		}
		return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
	}

	return cmd
}

type PluginLogsParams struct {
	PluginInstanceId string
}

type PluginLogsHandler func(ctx context.Context, cmd *cobra.Command, io IO, p PluginLogsParams) error

func newPluginLogsCommand(pluginLogs PluginLogsHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs <plugin-instance-id>",
		Short: "Stream supervised plugin logs.",
		Long: `Streams log lines until the plugin's log stream closes or the
command is interrupted. There is currently no --follow/--tail/
--since — this always streams from the beginning, unbounded.
`,
		Args: rangeArgs(1, 1),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawPluginInstanceId := args[0]

		p := PluginLogsParams{
			PluginInstanceId: rawPluginInstanceId,
		}

		return pluginLogs(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type PluginTokenRotateParams struct {
	PluginInstanceId string
}

type PluginTokenRotateHandler func(ctx context.Context, cmd *cobra.Command, io IO, p PluginTokenRotateParams) error

func newPluginTokenRotateCommand(pluginTokenRotate PluginTokenRotateHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rotate <plugin-instance-id>",
		Short: "Regenerate and print an external plugin RPC token.",
		Long: `The new token is shown once, in plain text, and is not
recoverable afterwards — save it immediately.
`,
		Args: rangeArgs(1, 1),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawPluginInstanceId := args[0]

		p := PluginTokenRotateParams{
			PluginInstanceId: rawPluginInstanceId,
		}

		return pluginTokenRotate(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

func newPluginTokenCommand(
	pluginTokenRotate PluginTokenRotateHandler,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Manage external plugin RPC tokens.",
	}

	cmd.AddCommand(
		newPluginTokenRotateCommand(pluginTokenRotate),
	)

	return cmd
}

type PluginDevNewParams struct {
	Command   string
	Author    string
	Force     bool
	Directory string
}

type PluginDevNewHandler func(ctx context.Context, cmd *cobra.Command, io IO, p PluginDevNewParams) error

func newPluginDevNewCommand(pluginDevNew PluginDevNewHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new <directory>",
		Short: "Create a minimal plugin project.",
		Args:  rangeArgs(1, 1),
	}

	var rawCommand string
	cmd.Flags().StringVar(&rawCommand, "command", "", "Command used to run the plugin during dev (defaults to ./<directory-name>).")
	var rawAuthor string
	cmd.Flags().StringVar(&rawAuthor, "author", "", "Plugin author slug (defaults to an environment/git-derived value).")
	var rawForce bool
	cmd.Flags().BoolVar(&rawForce, "force", false, "Overwrite/proceed without confirmation.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawDirectory := args[0]

		p := PluginDevNewParams{
			Command:   rawCommand,
			Author:    rawAuthor,
			Force:     rawForce,
			Directory: rawDirectory,
		}

		return pluginDevNew(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type PluginDevRunParams struct {
	Manifest   string
	InstanceId string
	NoUpdate   bool
	Command    []string
}

type PluginDevRunHandler func(ctx context.Context, cmd *cobra.Command, io IO, p PluginDevRunParams) error

func newPluginDevRunCommand(pluginDevRun PluginDevRunHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run [command]...",
		Short: "Run a local plugin with STORYDEN_RPC_URL set from the current instance.",
		Long: `With no COMMAND, runs the command from the plugin's manifest.
Pass a command (conventionally after --) to override it.
`,
		Args: rangeArgs(0, -1),
	}

	var rawManifest string
	cmd.Flags().StringVarP(&rawManifest, "manifest", "m", "manifest.yaml", "Path to the plugin manifest file.")
	var rawInstanceId string
	cmd.Flags().StringVar(&rawInstanceId, "instance-id", "", "Existing plugin installation ID to update instead of matching by manifest id.")
	var rawNoUpdate bool
	cmd.Flags().BoolVar(&rawNoUpdate, "no-update", false, "Skip updating the existing external plugin manifest before running.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		command := args[0:]

		p := PluginDevRunParams{
			Manifest:   rawManifest,
			InstanceId: rawInstanceId,
			NoUpdate:   rawNoUpdate,
			Command:    command,
		}

		return pluginDevRun(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type PluginDevPackageParams struct {
	Dir        string
	Manifest   string
	OutputFile string
	Force      bool
}

type PluginDevPackageHandler func(ctx context.Context, cmd *cobra.Command, io IO, p PluginDevPackageParams) error

func newPluginDevPackageCommand(pluginDevPackage PluginDevPackageHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "package",
		Short: "Create a supervised plugin package zip.",
		Args:  rangeArgs(0, 0),
	}

	var rawDir string
	cmd.Flags().StringVar(&rawDir, "dir", ".", "Plugin project directory.")
	var rawManifest string
	cmd.Flags().StringVarP(&rawManifest, "manifest", "m", "manifest.yaml", "Path to the plugin manifest file.")
	var rawOutputFile string
	cmd.Flags().StringVar(&rawOutputFile, "output-file", "", "Output file path (defaults to a name derived from the source being written).")
	var rawForce bool
	cmd.Flags().BoolVar(&rawForce, "force", false, "Overwrite/proceed without confirmation.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		p := PluginDevPackageParams{
			Dir:        rawDir,
			Manifest:   rawManifest,
			OutputFile: rawOutputFile,
			Force:      rawForce,
		}

		return pluginDevPackage(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type PluginDevValidateParams struct {
	Dir      string
	Manifest string
}

type PluginDevValidateHandler func(ctx context.Context, cmd *cobra.Command, io IO, p PluginDevValidateParams) error

func newPluginDevValidateCommand(pluginDevValidate PluginDevValidateHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a plugin manifest and supervised package.",
		Args:  rangeArgs(0, 0),
	}

	var rawDir string
	cmd.Flags().StringVar(&rawDir, "dir", ".", "Plugin project directory.")
	var rawManifest string
	cmd.Flags().StringVarP(&rawManifest, "manifest", "m", "manifest.yaml", "Path to the plugin manifest file.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		p := PluginDevValidateParams{
			Dir:      rawDir,
			Manifest: rawManifest,
		}

		return pluginDevValidate(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type PluginDevInstallParams struct {
	Dir      string
	Manifest string
}

type PluginDevInstallHandler func(ctx context.Context, cmd *cobra.Command, io IO, p PluginDevInstallParams) error

func newPluginDevInstallCommand(pluginDevInstall PluginDevInstallHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install or update a local plugin as supervised.",
		Args:  rangeArgs(0, 0),
	}

	var rawDir string
	cmd.Flags().StringVar(&rawDir, "dir", ".", "Plugin project directory.")
	var rawManifest string
	cmd.Flags().StringVarP(&rawManifest, "manifest", "m", "manifest.yaml", "Path to the plugin manifest file.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		p := PluginDevInstallParams{
			Dir:      rawDir,
			Manifest: rawManifest,
		}

		return pluginDevInstall(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type PluginDevDownloadParams struct {
	Dir              string
	NoUnzip          bool
	Force            bool
	PluginInstanceId string
}

type PluginDevDownloadHandler func(ctx context.Context, cmd *cobra.Command, io IO, p PluginDevDownloadParams) error

func newPluginDevDownloadCommand(pluginDevDownload PluginDevDownloadHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download <plugin-instance-id>",
		Short: "Download an installed plugin package.",
		Args:  rangeArgs(1, 1),
	}

	var rawDir string
	cmd.Flags().StringVar(&rawDir, "dir", ".", "Plugin project directory.")
	var rawNoUnzip bool
	cmd.Flags().BoolVar(&rawNoUnzip, "no-unzip", false, "Save the zip without extracting it.")
	var rawForce bool
	cmd.Flags().BoolVar(&rawForce, "force", false, "Overwrite/proceed without confirmation.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawPluginInstanceId := args[0]

		p := PluginDevDownloadParams{
			Dir:              rawDir,
			NoUnzip:          rawNoUnzip,
			Force:            rawForce,
			PluginInstanceId: rawPluginInstanceId,
		}

		return pluginDevDownload(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type PluginDevSymbolsPackagesParams struct {
	Dir         string
	Pattern     string
	IncludeDeps bool
	Limit       int
}

type PluginDevSymbolsPackagesHandler func(ctx context.Context, cmd *cobra.Command, io IO, p PluginDevSymbolsPackagesParams) error

func newPluginDevSymbolsPackagesCommand(pluginDevSymbolsPackages PluginDevSymbolsPackagesHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "packages",
		Short: "List Go packages in a plugin project.",
		Args:  rangeArgs(0, 0),
	}

	var rawDir string
	cmd.Flags().StringVar(&rawDir, "dir", ".", "Plugin project directory.")
	var rawPattern string
	cmd.Flags().StringVar(&rawPattern, "pattern", "./...", "Go package pattern to scan.")
	var rawIncludeDeps bool
	cmd.Flags().BoolVar(&rawIncludeDeps, "include-deps", false, "Include third-party dependency packages.")
	var rawLimit int
	cmd.Flags().IntVar(&rawLimit, "limit", 0, "Stop after N matches (0 = no limit).")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		p := PluginDevSymbolsPackagesParams{
			Dir:         rawDir,
			Pattern:     rawPattern,
			IncludeDeps: rawIncludeDeps,
			Limit:       rawLimit,
		}

		return pluginDevSymbolsPackages(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type PluginDevSymbolsPackageParams struct {
	Dir               string
	IncludeUnexported bool
	Limit             int
	ImportPath        string
}

type PluginDevSymbolsPackageHandler func(ctx context.Context, cmd *cobra.Command, io IO, p PluginDevSymbolsPackageParams) error

func newPluginDevSymbolsPackageCommand(pluginDevSymbolsPackage PluginDevSymbolsPackageHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "package <import-path>",
		Short: "List exported symbols in a Go package.",
		Args:  rangeArgs(1, 1),
	}

	var rawDir string
	cmd.Flags().StringVar(&rawDir, "dir", ".", "Plugin project directory.")
	var rawIncludeUnexported bool
	cmd.Flags().BoolVar(&rawIncludeUnexported, "include-unexported", false, "Include unexported symbols.")
	var rawLimit int
	cmd.Flags().IntVar(&rawLimit, "limit", 0, "Stop after N matches (0 = no limit).")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawImportPath := args[0]

		p := PluginDevSymbolsPackageParams{
			Dir:               rawDir,
			IncludeUnexported: rawIncludeUnexported,
			Limit:             rawLimit,
			ImportPath:        rawImportPath,
		}

		return pluginDevSymbolsPackage(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type PluginDevSymbolsDetailParams struct {
	Dir        string
	ImportPath string
	Symbol     string
}

type PluginDevSymbolsDetailHandler func(ctx context.Context, cmd *cobra.Command, io IO, p PluginDevSymbolsDetailParams) error

func newPluginDevSymbolsDetailCommand(pluginDevSymbolsDetail PluginDevSymbolsDetailHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detail <import-path> <symbol>",
		Short: "Show detail for a single Go symbol.",
		Args:  rangeArgs(2, 2),
	}

	var rawDir string
	cmd.Flags().StringVar(&rawDir, "dir", ".", "Plugin project directory.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawImportPath := args[0]
		rawSymbol := args[1]

		p := PluginDevSymbolsDetailParams{
			Dir:        rawDir,
			ImportPath: rawImportPath,
			Symbol:     rawSymbol,
		}

		return pluginDevSymbolsDetail(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type PluginDevSymbolsSearchParams struct {
	Dir               string
	Pattern           string
	IncludeDeps       bool
	IncludeUnexported bool
	Limit             int
	Query             string
}

type PluginDevSymbolsSearchHandler func(ctx context.Context, cmd *cobra.Command, io IO, p PluginDevSymbolsSearchParams) error

func newPluginDevSymbolsSearchCommand(pluginDevSymbolsSearch PluginDevSymbolsSearchHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search Go symbols in a plugin project.",
		Long:  "This searches Go source symbols within a local plugin\nproject — unrelated to `sd search` (Storyden content) or\n`sd node search` (nodes only).\n",
		Args:  rangeArgs(1, 1),
	}

	var rawDir string
	cmd.Flags().StringVar(&rawDir, "dir", ".", "Plugin project directory.")
	var rawPattern string
	cmd.Flags().StringVar(&rawPattern, "pattern", "./...", "Go package pattern to scan.")
	var rawIncludeDeps bool
	cmd.Flags().BoolVar(&rawIncludeDeps, "include-deps", false, "Include third-party dependency packages.")
	var rawIncludeUnexported bool
	cmd.Flags().BoolVar(&rawIncludeUnexported, "include-unexported", false, "Include unexported symbols.")
	var rawLimit int
	cmd.Flags().IntVar(&rawLimit, "limit", 0, "Stop after N matches (0 = no limit).")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		rawQuery := args[0]

		p := PluginDevSymbolsSearchParams{
			Dir:               rawDir,
			Pattern:           rawPattern,
			IncludeDeps:       rawIncludeDeps,
			IncludeUnexported: rawIncludeUnexported,
			Limit:             rawLimit,
			Query:             rawQuery,
		}

		return pluginDevSymbolsSearch(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

func newPluginDevSymbolsCommand(
	pluginDevSymbolsPackages PluginDevSymbolsPackagesHandler,
	pluginDevSymbolsPackage PluginDevSymbolsPackageHandler,
	pluginDevSymbolsDetail PluginDevSymbolsDetailHandler,
	pluginDevSymbolsSearch PluginDevSymbolsSearchHandler,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "symbols",
		Short: "Discover Go packages and symbols in a plugin project.",
	}

	cmd.AddCommand(
		newPluginDevSymbolsPackagesCommand(pluginDevSymbolsPackages),
		newPluginDevSymbolsPackageCommand(pluginDevSymbolsPackage),
		newPluginDevSymbolsDetailCommand(pluginDevSymbolsDetail),
		newPluginDevSymbolsSearchCommand(pluginDevSymbolsSearch),
	)

	return cmd
}

func newPluginDevCommand(
	pluginDevNew PluginDevNewHandler,
	pluginDevRun PluginDevRunHandler,
	pluginDevPackage PluginDevPackageHandler,
	pluginDevValidate PluginDevValidateHandler,
	pluginDevInstall PluginDevInstallHandler,
	pluginDevDownload PluginDevDownloadHandler,
	pluginDevSymbolsPackages PluginDevSymbolsPackagesHandler,
	pluginDevSymbolsPackage PluginDevSymbolsPackageHandler,
	pluginDevSymbolsDetail PluginDevSymbolsDetailHandler,
	pluginDevSymbolsSearch PluginDevSymbolsSearchHandler,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dev",
		Short: "Create, run, package, validate, and install plugin projects.",
	}

	cmd.AddCommand(
		newPluginDevNewCommand(pluginDevNew),
		newPluginDevRunCommand(pluginDevRun),
		newPluginDevPackageCommand(pluginDevPackage),
		newPluginDevValidateCommand(pluginDevValidate),
		newPluginDevInstallCommand(pluginDevInstall),
		newPluginDevDownloadCommand(pluginDevDownload),
		newPluginDevSymbolsCommand(
			pluginDevSymbolsPackages,
			pluginDevSymbolsPackage,
			pluginDevSymbolsDetail,
			pluginDevSymbolsSearch,
		),
	)

	return cmd
}

func NewPluginCommand(
	pluginGet PluginGetHandler,
	pluginList PluginListHandler,
	pluginDelete PluginDeleteHandler,
	pluginActivate PluginActivateHandler,
	pluginDeactivate PluginDeactivateHandler,
	pluginLogs PluginLogsHandler,
	pluginTokenRotate PluginTokenRotateHandler,
	pluginDevNew PluginDevNewHandler,
	pluginDevRun PluginDevRunHandler,
	pluginDevPackage PluginDevPackageHandler,
	pluginDevValidate PluginDevValidateHandler,
	pluginDevInstall PluginDevInstallHandler,
	pluginDevDownload PluginDevDownloadHandler,
	pluginDevSymbolsPackages PluginDevSymbolsPackagesHandler,
	pluginDevSymbolsPackage PluginDevSymbolsPackageHandler,
	pluginDevSymbolsDetail PluginDevSymbolsDetailHandler,
	pluginDevSymbolsSearch PluginDevSymbolsSearchHandler,
) PluginCommand {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Build, develop, and manage Storyden plugins.",
	}

	cmd.AddCommand(
		newPluginGetCommand(pluginGet),
		newPluginListCommand(pluginList),
		newPluginDeleteCommand(pluginDelete),
		newPluginActivateCommand(pluginActivate),
		newPluginDeactivateCommand(pluginDeactivate),
		newPluginLogsCommand(pluginLogs),
		newPluginTokenCommand(
			pluginTokenRotate,
		),
		newPluginDevCommand(
			pluginDevNew,
			pluginDevRun,
			pluginDevPackage,
			pluginDevValidate,
			pluginDevInstall,
			pluginDevDownload,
			pluginDevSymbolsPackages,
			pluginDevSymbolsPackage,
			pluginDevSymbolsDetail,
			pluginDevSymbolsSearch,
		),
	)

	return PluginCommand(cmd)
}

type PluginCommand *cobra.Command
