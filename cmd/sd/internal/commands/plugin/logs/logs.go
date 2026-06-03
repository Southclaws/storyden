package logs

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
)

type LogsCommand *cobra.Command

func New(store *config.Store) LogsCommand {
	command := &cobra.Command{
		Use:   "logs <plugin-instance-id>",
		Short: "Stream supervised plugin logs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.NewAuthenticatedClient(cmd.Context(), store, api.WithRequestTimeout(0))
			if err != nil {
				return err
			}
			response, err := client.OpenAPI.PluginGetLogs(cmd.Context(), openapi.PluginIDParam(args[0]))
			if err != nil {
				return err
			}
			defer response.Body.Close()
			if response.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(response.Body)
				return output.RequestErrorWithMessages("plugin logs request", statusAdapter{response}, body, output.UnauthorizedMessage("plugin logs request"))
			}
			return printEventStream(cmd.OutOrStdout(), response.Body)
		},
	}
	help.SetupMarkdownHelp(command)
	return LogsCommand(command)
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
