package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
)

type RegistrationOptions struct {
	Endpoint    string
	Name        string
	Handle      string
	Scope       string
	AccessToken string
	Storage     config.AuthStorage
}

type RegistrationResult struct {
	Context  string `json:"context"`
	ClientID string `json:"client_id,omitempty"`
	Scope    string `json:"scope,omitempty"`
	*config.Registration
}

func Register(ctx context.Context, store *config.Store, opts RegistrationOptions) (*RegistrationResult, error) {
	if opts.Name == "" || opts.Handle == "" {
		return nil, fmt.Errorf("--name and --handle are required")
	}

	err := store.WithContextLock(ctx, opts.Name, func() error {
		cfg, err := store.Load()
		if err != nil {
			return err
		}

		if _, exists := cfg.Contexts[opts.Name]; exists {
			return fmt.Errorf("context %q already exists; use sd auth register check to resume", opts.Name)
		}

		client, err := NewClient(ctx, opts.Endpoint)
		if err != nil {
			return err
		}

		endpoint := client.Discovery.RegistrationEndpoint
		if err := validateOAuthEndpoint(client.Endpoint, endpoint); err != nil {
			return err
		}

		if err := validateOAuthEndpoint(client.Endpoint, client.Discovery.TokenEndpoint); err != nil {
			return err
		}

		private, kid, public, err := generateClientKey()
		if err != nil {
			return err
		}

		current := config.Context{APIURL: client.Endpoint, AuthType: opts.Storage, Auth: &config.Auth{
			Method: config.AuthMethodOAuthClient, PrivateKey: private, KeyID: kid, Issuer: client.Discovery.Issuer,
			TokenEndpoint: client.Discovery.TokenEndpoint, RequestedScope: opts.Scope,
			Registration: &config.Registration{State: "outcome_unknown", Handle: opts.Handle, Endpoint: endpoint},
		}}
		if err := store.Update(ctx, func(cfg *config.Config) error {
			if _, exists := cfg.Contexts[opts.Name]; exists {
				return fmt.Errorf("context %q already exists", opts.Name)
			}

			cfg.UpsertContext(opts.Name, current)
			return nil
		}); err != nil {
			return err
		}

		method := "private_key_jwt"
		grants := []string{"client_credentials"}
		modes := []string{"approval"}
		request := openapi.OAuthClientRegisterProps{
			ClientName: &opts.Handle, GrantTypes: &grants,
			TokenEndpointAuthMethod: &method, Jwks: &public, RegistrationMode: &modes,
		}

		if opts.Scope != "" {
			request.Scope = &opts.Scope
		}

		return exchangeRegistration(ctx, store, opts.Name, current, request, opts.AccessToken)
	})
	if err != nil {
		return registrationResult(store, opts.Name, err)
	}

	return finishRegistration(ctx, store, opts.Name)
}

func CheckRegistration(ctx context.Context, store *config.Store, cancel bool) (*RegistrationResult, error) {
	name, _, err := store.Current()
	if err != nil {
		return nil, err
	}

	err = store.WithContextLock(ctx, name, func() error {
		cfg, err := store.Load()
		if err != nil {
			return err
		}

		current, ok := cfg.Contexts[name]
		if !ok || current.Auth == nil || current.Auth.Registration == nil {
			return fmt.Errorf("context %q has no autonomous registration", name)
		}

		return pollRegistration(ctx, store, name, current, cancel)
	})
	if err != nil {
		return registrationResult(store, name, err)
	}

	if cancel {
		return registrationResult(store, name, nil)
	}

	return finishRegistration(ctx, store, name)
}

func exchangeRegistration(ctx context.Context, store *config.Store, name string, current config.Context, request openapi.OAuthClientRegisterProps, bearer string) error {
	data, err := json.Marshal(request)
	if err != nil {
		return err
	}

	response, err := oauthRequest(ctx, current.APIURL, current.Auth.Registration.Endpoint, "application/json", bytes.NewReader(data), bearer)
	if err != nil {
		return fmt.Errorf("registration outcome is unknown; do not register another identity automatically: %w", err)
	}

	defer response.Body.Close()
	resultErr := applyRegistrationResponse(response, current.Auth)
	if err := saveRegistration(context.WithoutCancel(ctx), store, name, current); err != nil {
		return err
	}

	return resultErr
}

func applyRegistrationResponse(response *http.Response, auth *config.Auth) error {
	r := auth.Registration
	decoder := json.NewDecoder(io.LimitReader(response.Body, 1<<20))
	switch response.StatusCode {
	case http.StatusCreated:
		var registered openapi.OAuthClientRegistration
		if err := decoder.Decode(&registered); err != nil {
			return err
		}

		if registered.ClientId == "" || registered.TokenEndpointAuthMethod != "private_key_jwt" {
			return fmt.Errorf("invalid client registration response")
		}

		auth.ClientID = registered.ClientId
		if registered.Scope != nil {
			auth.Scope = *registered.Scope
		}

		*r = config.Registration{State: "registered", Handle: r.Handle, Endpoint: r.Endpoint}
	case http.StatusAccepted:
		return applyPendingRegistration(decoder, r)
	case http.StatusTooManyRequests:
		if r.Code == "" {
			r.State = "rejected"
			return fmt.Errorf("registration rate limited; retry with a new local context after the server's cooldown")
		}

		r.State = "pending"
		r.Interval += 5
		r.NextPollAt = time.Now().Add(time.Duration(r.Interval) * time.Second)
		if retry, ok := parseHeaderTime(response.Header.Get("Retry-After"), time.Now()); ok && retry.After(r.NextPollAt) {
			r.NextPollAt = retry
		}
	case http.StatusNoContent:
		r.State, r.Code = "cancelled", ""
	default:
		var failure openapi.OAuthError
		if err := decoder.Decode(&failure); err != nil {
			return fmt.Errorf("registration failed (HTTP %d); outcome is unknown", response.StatusCode)
		}

		if response.StatusCode >= 500 {
			return fmt.Errorf("registration failed (HTTP %d); outcome is unknown", response.StatusCode)
		}

		if failure.Error == "" {
			failure.Error = "rejected"
		}

		if r.Code != "" && failure.Error != "access_denied" && failure.Error != "expired_token" && failure.Error != "invalid_registration" {
			r.State = "pending"
		} else {
			r.State, r.Code = failure.Error, ""
		}

		return fmt.Errorf("registration failed: %s", failure.Error)
	}

	return nil
}

func saveRegistration(ctx context.Context, store *config.Store, name string, current config.Context) error {
	return store.Update(ctx, func(cfg *config.Config) error {
		if _, ok := cfg.Contexts[name]; !ok {
			return fmt.Errorf("context %q was removed", name)
		}

		cfg.UpsertContext(name, current)
		return nil
	})
}

func finishRegistration(ctx context.Context, store *config.Store, name string) (*RegistrationResult, error) {
	cfg, err := store.Load()
	if err != nil {
		return nil, err
	}

	current := cfg.Contexts[name]
	if current.Auth == nil || current.Auth.Registration == nil {
		return nil, fmt.Errorf("context %q has no registration credentials", name)
	}

	if current.Auth.ClientID != "" {
		session := &AuthSession{store: store, contextName: name, context: current}
		if _, err := session.clientCredentials(ctx, false); err != nil {
			return registrationResult(store, name, err)
		}
	}

	return registrationResult(store, name, nil)
}

func registrationResult(store *config.Store, name string, cause error) (*RegistrationResult, error) {
	cfg, err := store.Load()
	if err != nil {
		return nil, err
	}

	current, ok := cfg.Contexts[name]
	if !ok || current.Auth == nil || current.Auth.Registration == nil {
		return nil, cause
	}

	return &RegistrationResult{Context: name, ClientID: current.Auth.ClientID, Scope: current.Auth.Scope, Registration: current.Auth.Registration}, cause
}

func applyPendingRegistration(decoder *json.Decoder, r *config.Registration) error {
	var pending openapi.OAuthRegistrationPending
	if err := decoder.Decode(&pending); err != nil {
		return err
	}

	if r.Code == "" {
		if pending.RegistrationCode == nil || *pending.RegistrationCode == "" || pending.ExpiresIn == nil || *pending.ExpiresIn <= 0 || pending.Interval == nil || *pending.Interval <= 0 || pending.VerificationCode == nil || pending.VerificationUri == nil {
			return fmt.Errorf("invalid pending registration response")
		}

		r.Code = *pending.RegistrationCode
		r.VerificationCode = *pending.VerificationCode
		r.VerificationURI = *pending.VerificationUri
		if pending.VerificationUriComplete != nil {
			r.VerificationURIComplete = *pending.VerificationUriComplete
		}

		r.ExpiresAt = time.Now().Add(time.Duration(*pending.ExpiresIn) * time.Second)
		r.Interval = *pending.Interval
	}

	r.State = "pending"
	r.NextPollAt = time.Now().Add(time.Duration(r.Interval) * time.Second)
	return nil
}

func pollRegistration(ctx context.Context, store *config.Store, name string, current config.Context, cancel bool) error {
	r := current.Auth.Registration

	// A failed poll can be retried with its existing single-use code. Initial
	// registrations without a code must not be replayed after an unknown outcome.
	if r.State != "pending" && !(r.State == "outcome_unknown" && r.Code != "") {
		if !cancel && (r.State == "ready" || r.State == "registered") {
			return nil
		}

		return fmt.Errorf("registration is %s", r.State)
	}

	if !time.Now().Before(r.ExpiresAt) {
		r.State = "expired_token"
		r.Code = ""
		if err := saveRegistration(ctx, store, name, current); err != nil {
			return err
		}

		return fmt.Errorf("registration expired")
	}

	if !cancel && time.Now().Before(r.NextPollAt) {
		return nil
	}

	request := openapi.OAuthClientRegisterProps{RegistrationCode: &r.Code}
	if cancel {
		request.CancelRegistration = &cancel
	}

	r.State = "outcome_unknown"
	if err := saveRegistration(ctx, store, name, current); err != nil {
		return err
	}

	return exchangeRegistration(ctx, store, name, current, request, "")
}
