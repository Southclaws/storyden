package register

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
	"github.com/spf13/cobra"
)

type RegisterCommand *cobra.Command

func New(store *config.Store) RegisterCommand {
	var opts api.RegistrationOptions
	var format, storage string
	var tokenStdin bool
	command := &cobra.Command{
		Use: "register [storyden-api-url]", Short: "Register an autonomous bot account",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateFormat(format); err != nil {
				return err
			}

			if opts.Name == "" || opts.Handle == "" {
				return fmt.Errorf("--name and --handle are required")
			}

			opts.Endpoint = args[0]
			switch storage {
			case "auto":
				opts.Storage = store.DefaultAuthStorage()
			case "file":
				opts.Storage = config.AuthStorageFile
			case "credential-store":
				if !store.CredentialStoreAvailable() {
					return fmt.Errorf("credential store is unavailable")
				}

				opts.Storage = config.AuthStorageCredentialStore
			default:
				return fmt.Errorf("--auth-storage must be auto, credential-store or file")
			}

			if tokenStdin {
				data, err := io.ReadAll(io.LimitReader(cmd.InOrStdin(), 64*1024+1))
				if err != nil {
					return err
				}

				opts.AccessToken = strings.TrimSpace(string(data))
				if len(data) > 64*1024 || opts.AccessToken == "" {
					return fmt.Errorf("invalid registration token input")
				}

				defer func() { opts.AccessToken = "" }()
			}

			scopes, err := cmd.Flags().GetStringSlice("scope")
			if err != nil {
				return err
			}

			opts.Scope = strings.Join(scopes, " ")
			result, err := api.Register(cmd.Context(), store, opts)
			return render(cmd, format, result, err)
		},
	}

	command.Flags().StringVar(&opts.Name, "name", "", "Local identity context name (required)")
	command.Flags().StringVar(&opts.Handle, "handle", "", "Requested bot handle (required)")
	command.Flags().StringSlice("scope", nil, "Restrict token permissions to these scopes; defaults to inheriting account roles (repeatable or comma-separated)")
	command.Flags().BoolVar(&tokenStdin, "registration-token-stdin", false, "Read an initial registration access token from stdin")
	command.Flags().StringVar(&storage, "auth-storage", "auto", "Credential storage: auto, credential-store, file")
	command.PersistentFlags().StringVar(&format, "format", "plain", "Output format: plain, json")
	command.AddCommand(newCheck(store, &format), newWait(store, &format), newCancel(store, &format))
	help.SetupMarkdownHelp(command)
	return RegisterCommand(command)
}

func newCheck(store *config.Store, format *string) *cobra.Command {
	cmd := &cobra.Command{Use: "check", Short: "Check registration once without waiting", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateFormat(*format); err != nil {
				return err
			}

			result, err := api.CheckRegistration(cmd.Context(), store, false)
			return render(cmd, *format, result, err)
		}}
	help.SetupMarkdownHelp(cmd)
	return cmd
}

func newWait(store *config.Store, format *string) *cobra.Command {
	var timeout time.Duration
	cmd := &cobra.Command{Use: "wait", Short: "Wait for registration approval and acquire a token", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateFormat(*format); err != nil {
				return err
			}

			if timeout < 0 {
				return fmt.Errorf("--timeout must not be negative")
			}

			ctx := cmd.Context()
			if timeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, timeout)
				defer cancel()
			}

			for {
				result, err := api.CheckRegistration(ctx, store, false)
				if err != nil || result == nil || result.State != "pending" {
					return render(cmd, *format, result, err)
				}

				delay := time.Until(result.NextPollAt)
				if remaining := time.Until(result.ExpiresAt); remaining < delay {
					delay = remaining
				}

				timer := time.NewTimer(max(delay, time.Millisecond))
				select {
				case <-ctx.Done():
					timer.Stop()
					return render(cmd, *format, result, pendingError{ctx.Err()})
				case <-timer.C:
				}
			}
		}}
	cmd.Flags().DurationVar(&timeout, "timeout", 0, "Maximum wait (0 waits until approval or registration expiry)")
	help.SetupMarkdownHelp(cmd)
	return cmd
}

func newCancel(store *config.Store, format *string) *cobra.Command {
	cmd := &cobra.Command{Use: "cancel", Short: "Cancel a pending registration", Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateFormat(*format); err != nil {
				return err
			}

			result, err := api.CheckRegistration(cmd.Context(), store, true)
			return render(cmd, *format, result, err)
		}}
	help.SetupMarkdownHelp(cmd)
	return cmd
}

func validateFormat(format string) error {
	if format != "plain" && format != "json" {
		return fmt.Errorf("--format must be plain or json")
	}

	return nil
}

func render(cmd *cobra.Command, format string, result *api.RegistrationResult, cause error) error {
	if result != nil {
		if format == "json" {
			if err := output.JSON(cmd.OutOrStdout(), result); err != nil {
				return err
			}
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "Context: %s\nHandle: %s\nState: %s\n", result.Context, result.Handle, result.State)
			if result.State == "pending" {
				fmt.Fprintf(cmd.OutOrStdout(), "Administrator verification: %s\nVerification code: %s\nExpires: %s\n", result.VerificationURI, result.VerificationCode, result.ExpiresAt.Format(time.RFC3339))
			}
		}
	}

	return cause
}

type pendingError struct{ cause error }

func (e pendingError) Error() string {
	return fmt.Sprintf("registration is still pending: %s", e.cause)
}
func (e pendingError) Unwrap() error { return e.cause }
func (e pendingError) ExitCode() int { return 2 }
