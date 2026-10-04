package credentials

import (
	"context"
	"fmt"
	"time"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
	"github.com/spf13/cobra"
)

func NewToken(store *config.Store) cligen.AuthTokenHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.AuthTokenParams) error {

		_, current, err := store.Current()
		if err != nil {
			return err
		}
		if current.Auth == nil || current.Auth.MethodOrDefault() == config.AuthMethodAccessKey {
			return fmt.Errorf("selected context does not use OAuth")
		}
		auth, err := api.Credentials(ctx, store, false)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(cio.Out, auth.AccessToken)
		return err
	}
}
func NewHeaders(store *config.Store) cligen.AuthHeadersHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.AuthHeadersParams) error {

		if store.SelectedContext == "" {
			return fmt.Errorf("auth headers requires an explicit --context")
		}
		ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
		defer cancel()
		auth, err := api.Credentials(ctx, store, true)
		if err != nil {
			return err
		}
		return output.JSON(cio.Out, map[string]string{"Authorization": "Bearer " + auth.AccessToken})
	}
}
func NewStatus(store *config.Store) cligen.AuthStatusHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.AuthStatusParams) error {
		format := string(p.Output)

		if format != "plain" && format != "json" {
			return fmt.Errorf("--output must be plain or json")
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
			return output.JSON(cio.Out, result)
		}
		fmt.Fprintf(cio.Out, "Context: %s\nEndpoint: %s\nAuthentication: %s\n", name, current.APIURL, auth.MethodOrDefault())
		if auth.Registration != nil {
			fmt.Fprintf(cio.Out, "Handle: %s\nState: %s\n", auth.Registration.Handle, auth.Registration.State)
		}
		return nil
	}
}
