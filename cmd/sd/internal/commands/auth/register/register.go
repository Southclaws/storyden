package register

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
	"github.com/spf13/cobra"
)

func New(store *config.Store) cligen.AuthRegisterHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.AuthRegisterParams) error {
		format := string(p.Output)
		opts := api.RegistrationOptions{Name: p.Name, Handle: p.Handle}
		storage := string(p.AuthStorage)

		if err := validateFormat(format); err != nil {
			return err
		}

		if opts.Name == "" || opts.Handle == "" {
			return fmt.Errorf("--name and --handle are required")
		}

		opts.Endpoint = p.StorydenApiUrl
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

		if p.RegistrationTokenStdin {
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

		opts.Scope = strings.Join(p.Scope, " ")
		result, err := api.Register(ctx, store, opts)
		return render(cmd, format, result, err)
	}
}
func NewCheck(store *config.Store) cligen.AuthRegisterCheckHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.AuthRegisterCheckParams) error {
		format := string(p.Output)

		if err := validateFormat(format); err != nil {
			return err
		}

		result, err := api.CheckRegistration(ctx, store, false)
		return render(cmd, format, result, err)
	}
}
func NewWait(store *config.Store) cligen.AuthRegisterWaitHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.AuthRegisterWaitParams) error {
		format := string(p.Output)
		timeout := p.Timeout

		if err := validateFormat(format); err != nil {
			return err
		}

		if timeout < 0 {
			return fmt.Errorf("--timeout must not be negative")
		}

		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		for {
			result, err := api.CheckRegistration(ctx, store, false)
			if err != nil || result == nil || result.State != "pending" {
				return render(cmd, format, result, err)
			}

			delay := time.Until(result.NextPollAt)
			if remaining := time.Until(result.ExpiresAt); remaining < delay {
				delay = remaining
			}

			timer := time.NewTimer(max(delay, time.Millisecond))
			select {
			case <-ctx.Done():
				timer.Stop()
				return render(cmd, format, result, pendingError{ctx.Err()})
			case <-timer.C:
			}
		}
	}
}
func NewCancel(store *config.Store) cligen.AuthRegisterCancelHandler {
	return func(ctx context.Context, cmd *cobra.Command, cio cligen.IO, p cligen.AuthRegisterCancelParams) error {
		format := string(p.Output)

		if err := validateFormat(format); err != nil {
			return err
		}

		result, err := api.CheckRegistration(ctx, store, true)
		return render(cmd, format, result, err)
	}
}
func validateFormat(format string) error {
	if format != "plain" && format != "json" {
		return fmt.Errorf("--output must be plain or json")
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
