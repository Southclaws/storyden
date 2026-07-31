package cligen

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
)

type AuthLoginAuthStorage string

const (
	AuthLoginAuthStorageAuto            AuthLoginAuthStorage = "auto"
	AuthLoginAuthStorageCredentialStore AuthLoginAuthStorage = "credential-store"
	AuthLoginAuthStorageFile            AuthLoginAuthStorage = "file"
)

type AuthLoginParams struct {
	AuthStorage    AuthLoginAuthStorage
	AccessKey      bool
	AccessKeyStdin bool
	StorydenApiUrl string
}

type AuthLoginHandler func(ctx context.Context, cmd *cobra.Command, io IO, p AuthLoginParams) error

func newAuthLoginCommand(authLogin AuthLoginHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login [storyden-api-url]",
		Short: "Log in to a Storyden instance.",
		Long:  "Authenticate against an instance's public web/API address. With no\nURL, either reuses the current context's endpoint or prompts for one\ninteractively. Defaults to an interactive OAuth device-code flow;\npass `--access-key` (or `--access-key-stdin`) to authenticate with a\nstatic access key instead.\n",
		Example: `  # Log in interactively via OAuth device flow.
  sd auth login https://demo.storyden.org
  # Log in with an access key (prompted).
  sd auth login https://demo.storyden.org --access-key`,
		Args: rangeArgs(0, 1),
	}

	var rawAuthStorage string
	cmd.Flags().StringVar(&rawAuthStorage, "auth-storage", "auto", "Where to store credentials.")
	var rawAccessKey bool
	cmd.Flags().BoolVar(&rawAccessKey, "access-key", false, "Authenticate with a static access key instead of the OAuth device flow.")
	var rawAccessKeyStdin bool
	cmd.Flags().BoolVar(&rawAccessKeyStdin, "access-key-stdin", false, "Like --access-key, but read the key from stdin.")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var authStorage AuthLoginAuthStorage
		if rawAuthStorage != "" {
			switch rawAuthStorage {
			case "auto":
				authStorage = AuthLoginAuthStorageAuto
			case "credential-store":
				authStorage = AuthLoginAuthStorageCredentialStore
			case "file":
				authStorage = AuthLoginAuthStorageFile
			default:
				return fmt.Errorf("invalid --auth-storage %q: must be one of auto, credential-store, file", rawAuthStorage)
			}
		}
		var rawStorydenApiUrl string
		if len(args) > 0 {
			rawStorydenApiUrl = args[0]
		}

		p := AuthLoginParams{
			AuthStorage:    authStorage,
			AccessKey:      rawAccessKey,
			AccessKeyStdin: rawAccessKeyStdin,
			StorydenApiUrl: rawStorydenApiUrl,
		}

		return authLogin(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type AuthRemoveParams struct {
	ContextName string
}

type AuthRemoveHandler func(ctx context.Context, cmd *cobra.Command, io IO, p AuthRemoveParams) error

func newAuthRemoveCommand(authRemove AuthRemoveHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [context]",
		Short: "Remove a Storyden auth context.",
		Args:  rangeArgs(0, 1),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var rawContextName string
		if len(args) > 0 {
			rawContextName = args[0]
		}

		p := AuthRemoveParams{
			ContextName: rawContextName,
		}

		return authRemove(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

type AuthSwitchParams struct {
	ContextName string
}

type AuthSwitchHandler func(ctx context.Context, cmd *cobra.Command, io IO, p AuthSwitchParams) error

func newAuthSwitchCommand(authSwitch AuthSwitchHandler) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "switch [context]",
		Short: "Switch the active Storyden auth context.",
		Args:  rangeArgs(0, 1),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var rawContextName string
		if len(args) > 0 {
			rawContextName = args[0]
		}

		p := AuthSwitchParams{
			ContextName: rawContextName,
		}

		return authSwitch(cmd.Context(), cmd, newIO(cmd), p)
	}

	return cmd
}

func NewAuthCommand(
	authLogin AuthLoginHandler,
	authRemove AuthRemoveHandler,
	authSwitch AuthSwitchHandler,
) AuthCommand {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with Storyden instances.",
	}

	cmd.AddCommand(
		newAuthLoginCommand(authLogin),
		newAuthRemoveCommand(authRemove),
		newAuthSwitchCommand(authSwitch),
	)

	return AuthCommand(cmd)
}

type AuthCommand *cobra.Command
