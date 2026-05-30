package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/Southclaws/storyden/app/services/authentication/oauth"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
)

const (
	refreshLeeway  = time.Minute
	requestTimeout = 30 * time.Second
)

type AuthenticatedClientOption func(*authenticatedClientOptions)

type authenticatedClientOptions struct {
	rateLimitWarnings io.Writer
}

func WithRateLimitWarnings(w io.Writer) AuthenticatedClientOption {
	return func(opts *authenticatedClientOptions) {
		opts.rateLimitWarnings = w
	}
}

type AuthSession struct {
	store       *config.Store
	contextName string
	context     config.Context
	tokenClient *openapi.ClientWithResponses

	mu sync.Mutex
}

func NewAuthenticatedClient(ctx context.Context, store *config.Store, options ...AuthenticatedClientOption) (*Client, error) {
	opts := authenticatedClientOptions{
		rateLimitWarnings: os.Stderr,
	}
	for _, option := range options {
		option(&opts)
	}

	contextName, currentContext, err := loadCurrentContext(store)
	if err != nil {
		return nil, err
	}

	var client *Client
	if currentContext.Auth.MethodOrDefault() == config.AuthMethodAccessKey {
		client, err = NewStaticClient(currentContext.APIURL)
	} else {
		client, err = NewClient(ctx, currentContext.APIURL)
	}
	if err != nil {
		return nil, err
	}

	session := &AuthSession{
		store:       store,
		contextName: contextName,
		context:     currentContext,
		tokenClient: client.OpenAPI,
	}

	authenticated, err := openapi.NewClientWithResponses(
		client.BaseURL,
		openapi.WithHTTPClient(authenticatedDoer{
			base:              &http.Client{Timeout: requestTimeout},
			session:           session,
			rateLimitWarnings: opts.rateLimitWarnings,
		}),
		openapi.WithRequestEditorFn(session.RequestEditor),
	)
	if err != nil {
		return nil, err
	}

	client.OpenAPI = authenticated

	return client, nil
}

func loadCurrentContext(store *config.Store) (string, config.Context, error) {
	cfg, err := store.Load()
	if err != nil {
		return "", config.Context{}, err
	}

	if cfg.CurrentContext == "" {
		return "", config.Context{}, fmt.Errorf("no current Storyden context; run sd auth login first")
	}

	currentContext, ok := cfg.Contexts[cfg.CurrentContext]
	if !ok {
		return "", config.Context{}, fmt.Errorf("current Storyden context %q was not found", cfg.CurrentContext)
	}

	if currentContext.Auth == nil {
		if currentContext.AuthType == config.AuthStorageCredentialStore {
			return "", config.Context{}, reauthenticateError("context %q credentials are stored in the credential store, but they could not be loaded", cfg.CurrentContext)
		}

		return "", config.Context{}, reauthenticateError("context %q is not authenticated", cfg.CurrentContext)
	}
	if currentContext.Auth.AccessToken == "" {
		return "", config.Context{}, reauthenticateError("context %q is not authenticated", cfg.CurrentContext)
	}

	return cfg.CurrentContext, currentContext, nil
}

func (s *AuthSession) RequestEditor(ctx context.Context, req *http.Request) error {
	auth, err := s.auth(ctx)
	if err != nil {
		return err
	}

	tokenType := auth.TokenType
	if tokenType == "" {
		tokenType = "Bearer"
	}

	req.Header.Set("Authorization", tokenType+" "+auth.AccessToken)

	return nil
}

func (s *AuthSession) auth(ctx context.Context) (*config.Auth, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	auth := s.context.Auth
	if auth == nil {
		return nil, reauthenticateError("context %q is not authenticated", s.contextName)
	}

	if auth.MethodOrDefault() == config.AuthMethodAccessKey {
		return auth, nil
	}

	if auth.RefreshToken == "" || auth.ExpiresAt.IsZero() || time.Now().Before(auth.ExpiresAt.Add(-refreshLeeway)) {
		return auth, nil
	}

	return s.refresh(ctx)
}

func (s *AuthSession) forceRefresh(ctx context.Context) (*config.Auth, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.refresh(ctx)
}

func (s *AuthSession) refresh(ctx context.Context) (*config.Auth, error) {
	auth := s.context.Auth
	if auth != nil && auth.MethodOrDefault() == config.AuthMethodAccessKey {
		return nil, reauthenticateError("context %q uses an access key and cannot be refreshed", s.contextName)
	}
	if auth.RefreshToken == "" {
		return nil, reauthenticateError("context %q access token expired and no refresh token is available", s.contextName)
	}

	token, err := s.tokenClient.OAuthTokenWithFormdataBodyWithResponse(ctx, openapi.OAuthTokenFormdataRequestBody{
		GrantType:    oauth.GrantTypeRefreshToken,
		ClientId:     oauth.StorydenCLIClientID,
		RefreshToken: &auth.RefreshToken,
	})
	if err != nil {
		return nil, err
	}
	if token.StatusCode() != http.StatusOK || token.JSON200 == nil {
		return nil, refreshError(s.contextName, s.context.APIURL, token)
	}

	next := *auth
	if token.JSON200.AccessToken != nil {
		next.AccessToken = *token.JSON200.AccessToken
	}
	if token.JSON200.RefreshToken != nil {
		next.RefreshToken = *token.JSON200.RefreshToken
	}
	if token.JSON200.TokenType != nil {
		next.TokenType = *token.JSON200.TokenType
	}
	if token.JSON200.ExpiresIn != nil {
		next.ExpiresAt = time.Now().Add(time.Duration(*token.JSON200.ExpiresIn) * time.Second)
	}
	if token.JSON200.Scope != nil {
		next.Scope = *token.JSON200.Scope
	}

	if err := s.save(next); err != nil {
		return nil, err
	}

	s.context.Auth = &next

	return &next, nil
}

func (s *AuthSession) save(auth config.Auth) error {
	cfg, err := s.store.Load()
	if err != nil {
		return err
	}

	currentContext, ok := cfg.Contexts[s.contextName]
	if !ok {
		return fmt.Errorf("current Storyden context %q was not found", s.contextName)
	}

	currentContext.Auth = &auth
	cfg.UpsertContext(s.contextName, currentContext)

	return s.store.Save(cfg)
}

func refreshError(contextName string, apiURL string, token *openapi.OAuthTokenResponse) error {
	if token.JSON400 != nil {
		switch token.JSON400.Error {
		case "invalid_grant":
			return reauthenticateError("context %q session expired or was revoked", contextName)
		case "invalid_client":
			return fmt.Errorf("Storyden CLI OAuth client was rejected by %s", apiURL)
		}

		if token.JSON400.Error != "" {
			return fmt.Errorf("OAuth refresh failed: %s", token.JSON400.Error)
		}
	}

	body := strings.TrimSpace(string(token.Body))
	if body != "" {
		return fmt.Errorf("OAuth refresh failed: %s: %s", token.Status(), body)
	}

	return fmt.Errorf("OAuth refresh failed: %s", token.Status())
}

func reauthenticateError(format string, args ...any) error {
	return fmt.Errorf("%s\n\nrun\n\n  \x1b[1msd auth login\x1b[0m\n\nto re-authenticate", fmt.Sprintf(format, args...))
}

type authenticatedDoer struct {
	base              openapi.HttpRequestDoer
	session           *AuthSession
	rateLimitWarnings io.Writer
}

func (d authenticatedDoer) Do(req *http.Request) (*http.Response, error) {
	response, err := d.base.Do(req)
	if err != nil || response.StatusCode != http.StatusUnauthorized {
		return d.handleResponse(response, err)
	}
	if !d.session.canRefresh() {
		return d.handleResponse(response, nil)
	}

	retry, ok, err := cloneRequest(req)
	if err != nil {
		return nil, err
	}
	if !ok {
		return response, nil
	}

	_, _ = io.Copy(io.Discard, response.Body)
	_ = response.Body.Close()

	if _, err := d.session.forceRefresh(req.Context()); err != nil {
		return nil, err
	}
	if err := d.session.RequestEditor(req.Context(), retry); err != nil {
		return nil, err
	}

	response, err = d.base.Do(retry)

	return d.handleResponse(response, err)
}

func (d authenticatedDoer) handleResponse(response *http.Response, err error) (*http.Response, error) {
	if err != nil || response == nil {
		return response, err
	}

	if response.StatusCode == http.StatusTooManyRequests {
		if response.Body != nil {
			_, _ = io.Copy(io.Discard, response.Body)
			_ = response.Body.Close()
		}

		return nil, rateLimitExceededFromHeaders(response.Header, time.Now())
	}

	info, ok := parseRateLimit(response.Header, time.Now())
	if !ok {
		return response, nil
	}

	if warning := rateLimitWarning(info); warning != "" && d.rateLimitWarnings != nil {
		fmt.Fprintln(d.rateLimitWarnings, warning)
	}

	return response, nil
}

func (s *AuthSession) canRefresh() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	auth := s.context.Auth
	return auth != nil && auth.MethodOrDefault() != config.AuthMethodAccessKey && auth.RefreshToken != ""
}

func cloneRequest(req *http.Request) (*http.Request, bool, error) {
	retry := req.Clone(req.Context())
	if req.Body == nil {
		return retry, true, nil
	}
	if req.GetBody == nil {
		return nil, false, nil
	}

	body, err := req.GetBody()
	if err != nil {
		return nil, false, err
	}

	retry.Body = body

	return retry, true, nil
}

type rateLimitInfo struct {
	Limit     int
	Remaining int
	ResetAt   time.Time
	RetryAt   time.Time
	Now       time.Time
}

func parseRateLimit(header http.Header, now time.Time) (rateLimitInfo, bool) {
	limit, limitOK := parseHeaderInt(header.Get("X-RateLimit-Limit"))
	remaining, remainingOK := parseHeaderInt(header.Get("X-RateLimit-Remaining"))
	resetAt, resetOK := parseHeaderTime(header.Get("X-RateLimit-Reset"), now)
	retryAt, _ := parseHeaderTime(header.Get("Retry-After"), now)

	if !limitOK && !remainingOK && !resetOK {
		return rateLimitInfo{}, false
	}

	return rateLimitInfo{
		Limit:     limit,
		Remaining: remaining,
		ResetAt:   resetAt,
		RetryAt:   retryAt,
		Now:       now,
	}, true
}

func rateLimitExceededFromHeaders(header http.Header, now time.Time) error {
	info, ok := parseRateLimit(header, now)
	if !ok {
		info = rateLimitInfo{Now: now}
	}

	return rateLimitExceededError(info)
}

func parseHeaderInt(value string) (int, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, false
	}

	return parsed, true
}

func parseHeaderTime(value string, now time.Time) (time.Time, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, false
	}

	if seconds, err := strconv.Atoi(value); err == nil {
		return now.Add(time.Duration(seconds) * time.Second), true
	}

	for _, layout := range []string{time.RFC1123, time.RFC1123Z, time.RFC3339} {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, true
		}
	}

	return time.Time{}, false
}

func rateLimitExceededError(info rateLimitInfo) error {
	resetAt := info.RetryAt
	if resetAt.IsZero() {
		resetAt = info.ResetAt
	}

	message := "Rate limit exceeded."
	if info.Limit > 0 {
		message += fmt.Sprintf("\n\nLimit: %d requests", info.Limit)
	}
	if !resetAt.IsZero() {
		message += fmt.Sprintf("\nReset time: %s (%s)", resetAt.Local().Format("2006-01-02 15:04:05 -07"), relativeResetTime(info.Now, resetAt))
	}

	message += "\n\nPlease wait for the reset window before retrying, or use a more narrowly scoped command."

	return fmt.Errorf("%s", message)
}

func rateLimitWarning(info rateLimitInfo) string {
	if info.Limit <= 0 || info.Remaining <= 0 {
		return ""
	}

	threshold := max(5, info.Limit/10)
	if info.Remaining > threshold {
		return ""
	}

	percent := float64(info.Remaining) / float64(info.Limit) * 100
	message := fmt.Sprintf(
		"Warning: Storyden API rate limit is getting low: %d/%d requests remaining (%.0f%%).",
		info.Remaining,
		info.Limit,
		percent,
	)
	if !info.ResetAt.IsZero() {
		message += fmt.Sprintf(" Reset time: %s (%s).", info.ResetAt.Local().Format("15:04:05 -07"), relativeResetTime(info.Now, info.ResetAt))
	}
	message += " Consider waiting before running broad list or tree commands."

	return message
}

func relativeResetTime(now time.Time, resetAt time.Time) string {
	if resetAt.Before(now) {
		return humanize.RelTime(resetAt, now, "ago", "from now")
	}

	relative := strings.TrimSpace(humanize.RelTime(now, resetAt, "", "ago"))
	if relative == "now" {
		return "now"
	}

	return "in about " + relative
}
