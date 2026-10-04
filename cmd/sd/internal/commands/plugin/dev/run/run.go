package run

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	plugindev "github.com/Southclaws/storyden/lib/plugin/dev"
)

func New(store *config.Store) cligen.PluginDevRunHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.PluginDevRunParams) error {
		mf, err := plugindev.ReadManifest(p.Manifest)
		if err != nil {
			return err
		}

		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		plugin, err := plugindev.EnsureExternalPlugin(ctx, client.OpenAPI, mf.Manifest, p.InstanceId, p.NoUpdate)
		if err != nil {
			return err
		}

		rpcURL, err := plugindev.ExternalRPCURL(client.Endpoint, plugin.Token)
		if err != nil {
			return err
		}

		runCommand, runArgs, err := plugindev.CommandFromManifest(mf.Manifest, p.Command)
		if err != nil {
			return err
		}

		fmt.Fprintf(io.Err, "Running %s with %s for plugin %s\n", runCommand, plugindev.RPCURLEnvName, plugin.ID)
		return runPluginCommand(ctx, runCommand, runArgs, rpcURL, filepath.Dir(mf.Path))
	}
}

func runPluginCommand(ctx context.Context, command string, args []string, rpcURL string, dir string) error {
	cmd := exec.CommandContext(ctx, command, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), plugindev.RPCURLEnvName+"="+rpcURL)
	return cmd.Run()
}
