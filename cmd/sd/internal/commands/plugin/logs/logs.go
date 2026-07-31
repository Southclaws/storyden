package logs

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
)

func New(store *config.Store) cligen.PluginLogsHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.PluginLogsParams) error {
		client, err := api.NewAuthenticatedClient(ctx, store, api.WithRequestTimeout(0))
		if err != nil {
			return err
		}
		response, err := client.OpenAPI.PluginGetLogs(ctx, openapi.PluginIDParam(p.PluginInstanceId))
		if err != nil {
			return err
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(response.Body)
			return output.RequestErrorWithMessages("plugin logs request", statusAdapter{response}, body, output.UnauthorizedMessage("plugin logs request"))
		}
		return printEventStream(cio.Out, response.Body)
	}
}

func printEventStream(out io.Writer, body io.Reader) error {
	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			fmt.Fprintln(out, strings.TrimPrefix(line, "data: "))
		}
	}
	return scanner.Err()
}

type statusAdapter struct{ *http.Response }

func (s statusAdapter) Status() string { return s.Response.Status }

func (s statusAdapter) StatusCode() int { return s.Response.StatusCode }
