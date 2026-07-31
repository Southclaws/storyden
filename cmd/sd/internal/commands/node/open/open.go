// Package open implements `sd node open <slug>`. It fetches a node and emits
// its attached link URL. The default is print-only so agents in a non-TTY
// pipeline don't accidentally launch a browser; pass --launch to actually
// open the URL.
package open

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/nodeapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
)

func New(store *config.Store) cligen.NodeOpenHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.NodeOpenParams) error {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		node, err := nodeapi.Fetch(ctx, client.OpenAPI, p.Slug)
		if err != nil {
			return err
		}

		if node.Link == nil {
			return fmt.Errorf("node %s has no attached link", p.Slug)
		}

		url := string(node.Link.Url)
		fmt.Fprintln(io.Out, url)

		if !p.Launch {
			return nil
		}

		if !output.IsTerminal(io.Out) && !p.Force {
			return fmt.Errorf("refusing to launch browser from non-tty; pass --force to override")
		}

		return launchURL(url)
	}
}

func launchURL(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("launch URL: %w", err)
	}
	return nil
}
