package credentials

import (
	"context"
	"fmt"
	"time"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
	"github.com/spf13/cobra"
)

type TokenCommand *cobra.Command
type HeadersCommand *cobra.Command

func NewToken(store *config.Store) TokenCommand {
	cmd := &cobra.Command{Use: "token", Short: "Print a short-lived OAuth access token", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, current, err := store.Current()
			if err != nil {
				return err
			}
			if current.Auth == nil || current.Auth.MethodOrDefault() == config.AuthMethodAccessKey {
				return fmt.Errorf("selected context does not use OAuth")
			}
			auth, err := api.Credentials(cmd.Context(), store, false)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), auth.AccessToken)
			return err
		}}
	help.SetupMarkdownHelp(cmd)
	return TokenCommand(cmd)
}

func NewHeaders(store *config.Store) HeadersCommand {
	cmd := &cobra.Command{Use: "headers", Short: "Print authenticated HTTP headers as JSON for an MCP host", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if store.SelectedContext == "" {
				return fmt.Errorf("auth headers requires an explicit --context")
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), 8*time.Second)
			defer cancel()
			auth, err := api.Credentials(ctx, store, true)
			if err != nil {
				return err
			}
			return output.JSON(cmd.OutOrStdout(), map[string]string{"Authorization": "Bearer " + auth.AccessToken})
		}}
	help.SetupMarkdownHelp(cmd)
	return HeadersCommand(cmd)
}

type StatusCommand *cobra.Command

func NewStatus(store *config.Store) StatusCommand {
	var format string
	cmd := &cobra.Command{Use: "status", Short: "Show local authentication state without contacting the server", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if format != "plain" && format != "json" {
				return fmt.Errorf("--format must be plain or json")
			}
			name, current, err := store.Current()
			if err != nil {
				return err
			}
			if current.Auth == nil {
				return fmt.Errorf("credentials for context %q are unavailable", name)
			}
			auth := current.Auth
			result := struct {
				Context   string            `json:"context"`
				Endpoint  string            `json:"endpoint"`
				Method    config.AuthMethod `json:"method"`
				ClientID  string            `json:"client_id,omitempty"`
				Scope     string            `json:"scope,omitempty"`
				ExpiresAt time.Time         `json:"token_expires_at,omitzero"`
				*config.Registration
			}{name, current.APIURL, auth.MethodOrDefault(), auth.ClientID, auth.Scope, auth.ExpiresAt, auth.Registration}
			if format == "json" {
				return output.JSON(cmd.OutOrStdout(), result)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Context: %s\nEndpoint: %s\nAuthentication: %s\n", name, current.APIURL, auth.MethodOrDefault())
			if auth.Registration != nil {
				fmt.Fprintf(cmd.OutOrStdout(), "Handle: %s\nState: %s\n", auth.Registration.Handle, auth.Registration.State)
			}
			return nil
		}}
	cmd.Flags().StringVar(&format, "format", "plain", "Output format: plain, json")
	help.SetupMarkdownHelp(cmd)
	return StatusCommand(cmd)
}
