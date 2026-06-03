package run

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/pluginapi"
)

type RunCommand *cobra.Command

func New(store *config.Store) RunCommand {
	var manifestPath string
	var instanceID string
	var noUpdate bool

	command := &cobra.Command{
		Use:   "run [-- command args...]",
		Short: "Run a local plugin with STORYDEN_RPC_URL from the current instance",
		Long: `# Run Plugin

Read ` + "`manifest.yaml`" + `, register or update it as an external plugin on the current Storyden instance, build a WebSocket RPC URL from the plugin's static external token, and run the manifest command with ` + "`STORYDEN_RPC_URL`" + ` in the environment.

Pass ` + "`--`" + ` followed by a command to override the manifest command for this run.
`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			mf, err := pluginapi.ReadManifest(manifestPath)
			if err != nil {
				return err
			}

			client, err := api.NewAuthenticatedClient(cmd.Context(), store)
			if err != nil {
				return err
			}

			plugin, err := pluginapi.EnsureExternalPlugin(cmd.Context(), client.OpenAPI, mf.Manifest, instanceID, noUpdate)
			if err != nil {
				return err
			}

			rpcURL, err := pluginapi.ExternalRPCURL(client.Endpoint, plugin.Token)
			if err != nil {
				return err
			}

			runCommand, runArgs, err := pluginapi.CommandFromManifest(mf.Manifest, args)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.ErrOrStderr(), "Running %s with %s for plugin %s\n", runCommand, pluginapi.RPCURLEnvName, plugin.ID)
			return runPluginCommand(cmd.Context(), runCommand, runArgs, rpcURL, filepath.Dir(mf.Path))
		},
	}

	command.Flags().StringVarP(&manifestPath, "manifest", "m", pluginapi.ManifestFilename, "Path to plugin manifest YAML")
	command.Flags().StringVar(&instanceID, "instance-id", "", "Existing plugin installation ID to update instead of matching by manifest id")
	command.Flags().BoolVar(&noUpdate, "no-update", false, "Do not update an existing external plugin manifest before running")

	help.SetupMarkdownHelp(command)

	return RunCommand(command)
}

func runPluginCommand(ctx context.Context, command string, args []string, rpcURL string, dir string) error {
	cmd := exec.CommandContext(ctx, command, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), pluginapi.RPCURLEnvName+"="+rpcURL)
	return cmd.Run()
}
