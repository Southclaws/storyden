package login

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"runtime"
	"strings"
	"time"
	"unicode"

	"charm.land/huh/v2"
	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/services/authentication/oauth"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/help"
	"github.com/Southclaws/storyden/cmd/sd/internal/tui"
)

type LoginCommand *cobra.Command

const (
	deviceAuthScope     = "openid profile offline_access"
	defaultPollInterval = 5 * time.Second
)

func New(
	store *config.Store,
) LoginCommand {
	authStorage := "auto"
	accessKey := false
	accessKeyStdin := false

	command := &cobra.Command{
		Use:   "login [storyden-api-url]",
		Short: "Log in to a Storyden instance",
		Long:  loginLongHelp(runtime.GOOS),
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			endpoint, err := endpointFromArgs(cmd, args, store)
			if err != nil {
				return err
			}

			storage, err := resolveAuthStorage(authStorage, store)
			if err != nil {
				return err
			}

			if accessKey || accessKeyStdin {
				auth, err := authenticateAccessKey(cmd, accessKeyStdin)
				if err != nil {
					return err
				}

				endpoint, err = api.CanonicalEndpoint(endpoint)
				if err != nil {
					return err
				}

				cfg, err := store.Load()
				if err != nil {
					return err
				}

				name := contextName(cfg, endpoint)
				cfg.UpsertContext(name, config.Context{
					APIURL:   endpoint,
					AuthType: storage,
					Auth:     auth,
				})
				cfg.SetCurrentContext(name)

				if err := store.Save(cfg); err != nil {
					return err
				}

				if storage == config.AuthStorageFile {
					warnFileAuthStorage(cmd, store)
				}

				fmt.Fprintf(cmd.OutOrStdout(), "%s %s as context %q\n", tui.Accent.Render("Authenticated with"), endpoint, name)

				return nil
			}

			client, err := api.NewClient(cmd.Context(), endpoint)
			if err != nil {
				return err
			}

			auth, err := authenticate(cmd.Context(), cmd, client)
			if err != nil {
				return err
			}

			cfg, err := store.Load()
			if err != nil {
				return err
			}

			name := contextName(cfg, client.Endpoint)
			cfg.UpsertContext(name, config.Context{
				APIURL:   client.Endpoint,
				AuthType: storage,
				Auth:     auth,
			})
			cfg.SetCurrentContext(name)

			if err := store.Save(cfg); err != nil {
				return err
			}

			if storage == config.AuthStorageFile {
				warnFileAuthStorage(cmd, store)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%s %s as context %q\n", tui.Accent.Render("Authenticated with"), client.Endpoint, name)

			return nil
		},
	}

	command.Flags().StringVar(&authStorage, "auth-storage", "auto", "Where to store credentials: auto, credential-store, or file")
	command.Flags().BoolVar(&accessKey, "access-key", false, "Authenticate with a Storyden access key instead of OAuth device auth")
	command.Flags().BoolVar(&accessKeyStdin, "access-key-stdin", false, "Read the access key from stdin")

	help.SetupMarkdownHelp(command)

	return command
}

func loginLongHelp(goos string) string {
	shell, stdinExample := accessKeyStdinExample(goos)

	return fmt.Sprintf(`# Authenticate with Storyden

Log in using OAuth2/OIDC device flow - secure browser-based authentication with no password input needed.

Each authenticated instance is saved as a "context" that you can switch between.

## Examples

Login with URL:
~~~bash
sd auth login https://my-community.com
~~~

Login interactively (prompts for URL):
~~~bash
sd auth login
~~~

Login to localhost:
~~~bash
sd auth login http://localhost:8000
~~~

Login with an access key:
~~~bash
sd auth login http://localhost:8000 --access-key --auth-storage file
~~~

Read an access key from stdin:
~~~%s
%s
~~~

Manage multiple instances:
~~~bash
sd auth login https://community1.com
sd auth login https://community2.com
sd auth switch  # Choose which to use
~~~
`, shell, stdinExample)
}

func accessKeyStdinExample(goos string) (string, string) {
	switch goos {
	case "windows":
		return "powershell", `$env:STORYDEN_ACCESS_KEY | sd auth login http://localhost:8000 --access-key-stdin --auth-storage file`
	case "darwin":
		return "zsh", `printf '%s' "$STORYDEN_ACCESS_KEY" | sd auth login http://localhost:8000 --access-key-stdin --auth-storage file`
	default:
		return "bash", `printf '%s' "$STORYDEN_ACCESS_KEY" | sd auth login http://localhost:8000 --access-key-stdin --auth-storage file`
	}
}

func resolveAuthStorage(value string, store *config.Store) (config.AuthStorage, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "auto":
		return store.DefaultAuthStorage(), nil

	case "file", "config", "config-file":
		return config.AuthStorageFile, nil

	case "credential-store", "credential_store", "keyring":
		if !store.CredentialStoreAvailable() {
			return "", fmt.Errorf("credential store is not available; use --auth-storage file to store credentials in the config file")
		}

		return config.AuthStorageCredentialStore, nil

	default:
		return "", fmt.Errorf("unsupported auth storage %q; expected auto, credential-store, or file", value)
	}
}

func warnFileAuthStorage(cmd *cobra.Command, store *config.Store) {
	fmt.Fprintf(
		cmd.ErrOrStderr(),
		"Warning: credentials will be stored in the config file at %s. Use --auth-storage credential-store on a supported desktop OS to store them securely.\n",
		store.Path(),
	)
}

func endpointFromArgs(cmd *cobra.Command, args []string, store *config.Store) (string, error) {
	if len(args) == 1 {
		return args[0], nil
	}

	endpoint := currentContextEndpoint(store)
	if endpoint != "" {
		return endpoint, nil
	}

	form := tui.NewForm(
		cmd.InOrStdin(),
		cmd.ErrOrStderr(),
		huh.NewGroup(
			huh.NewInput().
				Title(tui.Title.Render("Storyden instance")).
				Description("Paste the Storyden URL. /api is optional.").
				Placeholder("http://localhost:8000").
				Value(&endpoint),
		),
	)
	if err := form.RunWithContext(cmd.Context()); err != nil {
		return "", err
	}

	return endpoint, nil
}

func currentContextEndpoint(store *config.Store) string {
	cfg, err := store.Load()
	if err != nil || cfg.CurrentContext == "" {
		return ""
	}

	ctx, ok := cfg.Contexts[cfg.CurrentContext]
	if !ok {
		return ""
	}

	return ctx.APIURL
}

func authenticate(ctx context.Context, cmd *cobra.Command, client *api.Client) (*config.Auth, error) {
	start, err := client.OpenAPI.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(ctx, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
		ClientId: oauth.StorydenCLIClientID,
		Scope:    ptr(deviceAuthScope),
	})
	if err != nil {
		return nil, err
	}
	if start.StatusCode() != 200 || start.JSON200 == nil {
		return nil, oauthDeviceAuthorizationError(start)
	}

	device := start.JSON200
	if device.DeviceCode == nil || !hasDeviceVerification(device) {
		return nil, fmt.Errorf("OAuth device authorization response was missing required fields")
	}

	writeDeviceInstructions(cmd.OutOrStdout(), device)

	token, err := pollToken(ctx, client.OpenAPI, *device.DeviceCode, pollInterval(device.Interval), expiresAt(device.ExpiresIn))
	if err != nil {
		return nil, err
	}

	if token.AccessToken == nil || token.TokenType == nil {
		return nil, fmt.Errorf("OAuth token response was missing required fields")
	}

	auth := &config.Auth{
		Method:      config.AuthMethodOAuthDevice,
		AccessToken: *token.AccessToken,
		TokenType:   *token.TokenType,
		Issuer:      client.Discovery.Issuer,
		ClientID:    oauth.StorydenCLIClientID,
	}

	if token.RefreshToken != nil {
		auth.RefreshToken = *token.RefreshToken
	}
	if token.ExpiresIn != nil {
		auth.ExpiresAt = time.Now().Add(time.Duration(*token.ExpiresIn) * time.Second)
	}
	if token.Scope != nil {
		auth.Scope = *token.Scope
	}

	return auth, nil
}

func authenticateAccessKey(cmd *cobra.Command, stdin bool) (*config.Auth, error) {
	key := ""
	if stdin {
		data, err := io.ReadAll(cmd.InOrStdin())
		if err != nil {
			return nil, err
		}
		key = string(data)
	}

	key = strings.TrimSpace(key)
	if key == "" && !stdin {
		form := tui.NewForm(
			cmd.InOrStdin(),
			cmd.ErrOrStderr(),
			huh.NewGroup(
				huh.NewInput().
					Title(tui.Title.Render("Storyden access key")).
					Description("Paste an access key for this Storyden instance.").
					Password(true).
					Value(&key),
			),
		)
		if err := form.RunWithContext(cmd.Context()); err != nil {
			return nil, err
		}
		key = strings.TrimSpace(key)
	}

	if key == "" {
		return nil, fmt.Errorf("access key is required")
	}

	return &config.Auth{
		Method:      config.AuthMethodAccessKey,
		AccessToken: key,
		TokenType:   "Bearer",
	}, nil
}

func hasDeviceVerification(device *openapi.OAuthDeviceAuthorisation) bool {
	return device.VerificationUriComplete != nil || (device.VerificationUri != nil && device.UserCode != nil)
}

func writeDeviceInstructions(out io.Writer, device *openapi.OAuthDeviceAuthorisation) {
	fmt.Fprintln(out, tui.Title.Render("Open this URL to authenticate Storyden CLI:"))

	if device.VerificationUriComplete != nil {
		fmt.Fprintf(out, "\n%s\n\n", tui.URL.Render(*device.VerificationUriComplete))
		if device.UserCode != nil {
			fmt.Fprintf(out, "%s %s\n\n", tui.Muted.Render("Code:"), tui.Accent.Render(*device.UserCode))
		}
	} else {
		fmt.Fprintf(out, "\n%s\n\n", tui.URL.Render(*device.VerificationUri))
		fmt.Fprintf(out, "%s %s\n\n", tui.Muted.Render("Enter this code when prompted:"), tui.Accent.Render(*device.UserCode))
	}

	fmt.Fprintln(out, tui.Muted.Render("Waiting for browser approval..."))
}

func pollToken(
	ctx context.Context,
	client *openapi.ClientWithResponses,
	deviceCode string,
	interval time.Duration,
	expiresAt time.Time,
) (*openapi.OAuthToken, error) {
	if interval <= 0 {
		interval = defaultPollInterval
	}

	for {
		if !expiresAt.IsZero() && time.Now().After(expiresAt) {
			return nil, fmt.Errorf("OAuth device authorization expired")
		}

		token, err := client.OAuthTokenWithFormdataBodyWithResponse(ctx, openapi.OAuthTokenFormdataRequestBody{
			GrantType:  oauth.GrantTypeDeviceCode,
			ClientId:   oauth.StorydenCLIClientID,
			DeviceCode: &deviceCode,
		})
		if err != nil {
			return nil, err
		}

		switch token.StatusCode() {
		case 200:
			if token.JSON200 == nil {
				return nil, fmt.Errorf("OAuth token response was empty")
			}
			return token.JSON200, nil

		case 400:
			if token.JSON400 == nil {
				return nil, fmt.Errorf("OAuth token request failed: %s", token.Status())
			}

			switch token.JSON400.Error {
			case "authorization_pending":
				if err := wait(ctx, interval); err != nil {
					return nil, err
				}

			case "slow_down":
				interval += 5 * time.Second
				if err := wait(ctx, interval); err != nil {
					return nil, err
				}

			case "access_denied":
				return nil, fmt.Errorf("OAuth authorization was denied")

			case "expired_token":
				return nil, fmt.Errorf("OAuth device authorization expired")

			default:
				return nil, oauthTokenError(token.JSON400)
			}

		default:
			return nil, fmt.Errorf("OAuth token request failed: %s", token.Status())
		}
	}
}

func oauthDeviceAuthorizationError(response *openapi.OAuthDeviceAuthorisationResponse) error {
	if response.JSON400 != nil {
		return oauthTokenError(response.JSON400)
	}

	return fmt.Errorf("OAuth device authorization failed: %s", response.Status())
}

func oauthTokenError(err *openapi.OAuthError) error {
	if err.ErrorDescription != nil && *err.ErrorDescription != "" {
		return fmt.Errorf("OAuth error %q: %s", err.Error, *err.ErrorDescription)
	}

	return fmt.Errorf("OAuth error %q", err.Error)
}

func wait(ctx context.Context, interval time.Duration) error {
	timer := time.NewTimer(interval)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func pollInterval(value *int) time.Duration {
	if value == nil {
		return defaultPollInterval
	}

	return time.Duration(*value) * time.Second
}

func expiresAt(value *int) time.Time {
	if value == nil {
		return time.Time{}
	}

	return time.Now().Add(time.Duration(*value) * time.Second)
}

func ptr[T any](value T) *T {
	return &value
}

func contextName(cfg *config.Config, apiURL string) string {
	parsed, err := url.Parse(apiURL)
	if err != nil {
		return uniqueContextName(cfg, "storyden", apiURL)
	}

	base := slug(parsed.Host)
	if base == "" {
		base = "storyden"
	}

	return uniqueContextName(cfg, base, apiURL)
}

func uniqueContextName(cfg *config.Config, base string, apiURL string) string {
	if ctx, ok := cfg.Contexts[base]; !ok || ctx.APIURL == "" || ctx.APIURL == apiURL {
		return base
	}

	for i := 2; ; i++ {
		name := fmt.Sprintf("%s-%d", base, i)
		if ctx, ok := cfg.Contexts[name]; !ok || ctx.APIURL == apiURL {
			return name
		}
	}
}

func slug(value string) string {
	var builder strings.Builder
	previousDash := false

	for _, r := range strings.ToLower(value) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
			previousDash = false
			continue
		}

		if !previousDash {
			builder.WriteRune('-')
			previousDash = true
		}
	}

	return strings.Trim(builder.String(), "-")
}
