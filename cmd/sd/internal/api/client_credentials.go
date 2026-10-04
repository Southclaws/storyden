package api

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/go-jose/go-jose/v4"
	"github.com/golang-jwt/jwt/v5"
)

func Credentials(ctx context.Context, store *config.Store, fresh bool) (*config.Auth, error) {
	client, err := NewAuthenticatedClient(ctx, store)
	if err != nil {
		return nil, err
	}

	if fresh && client.session.context.Auth.Method == config.AuthMethodOAuthClient {
		return client.session.forceRefresh(ctx)
	}

	return client.session.auth(ctx)
}

func generateClientKey() (string, string, map[string]interface{}, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", nil, err
	}

	kid, err := assertionID()
	if err != nil {
		return "", "", nil, err
	}

	jwk := jose.JSONWebKey{Key: &key.PublicKey, KeyID: kid, Algorithm: "RS256", Use: "sig"}
	data, err := json.Marshal(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{jwk}})
	if err != nil {
		return "", "", nil, err
	}

	var public map[string]interface{}
	if err := json.Unmarshal(data, &public); err != nil {
		return "", "", nil, err
	}

	private := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	return string(private), kid, public, nil
}

func assertionID() (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b[:]), nil
}

func clientAssertion(auth *config.Auth) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(auth.PrivateKey))
	if err != nil {
		return "", fmt.Errorf("cannot read client private key: %w", err)
	}

	id, err := assertionID()
	if err != nil {
		return "", err
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Issuer: auth.ClientID, Subject: auth.ClientID, Audience: jwt.ClaimStrings{auth.TokenEndpoint},
		IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)), ID: id,
	})
	token.Header["kid"] = auth.KeyID
	return token.SignedString(key)
}

func (s *AuthSession) clientCredentials(ctx context.Context, force bool) (*config.Auth, error) {
	var result *config.Auth
	err := s.store.WithContextLock(ctx, s.contextName, func() error {
		cfg, err := s.store.Load()
		if err != nil {
			return err
		}

		current, ok := cfg.Contexts[s.contextName]
		if !ok || current.Auth == nil {
			return fmt.Errorf("credentials for context %q are unavailable", s.contextName)
		}

		auth := current.Auth
		if auth.Method != config.AuthMethodOAuthClient || auth.ClientID == "" {
			return fmt.Errorf("context %q has no registered autonomous client; use sd auth register check", s.contextName)
		}

		if !force && auth.AccessToken != "" && time.Now().Before(auth.ExpiresAt.Add(-refreshLeeway)) {
			result = auth
			return nil
		}

		next, err := exchangeClientCredentials(ctx, current.APIURL, auth)
		if err != nil {
			return err
		}

		if err := s.save(*next); err != nil {
			return err
		}

		result = next
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.context.Auth = result
	return result, nil
}

func oauthRequest(ctx context.Context, origin, endpoint, contentType string, body io.Reader, bearer string) (*http.Response, error) {
	if err := validateOAuthEndpoint(origin, endpoint); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}

	client := &http.Client{Timeout: discoveryTimeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	return client.Do(req)
}

func validateOAuthEndpoint(origin, endpoint string) error {
	base, err := url.Parse(origin)
	if err != nil {
		return err
	}

	target, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	if target.User != nil || target.Fragment != "" || target.Host == "" || target.Host != base.Host || target.Scheme != base.Scheme {
		return fmt.Errorf("OAuth endpoint must belong to the selected Storyden instance")
	}

	return nil
}

func oauthResponseError(response *http.Response) error {
	var body openapi.OAuthError
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&body); err == nil && body.Error != "" {
		return fmt.Errorf("OAuth request failed: %s (HTTP %d)", body.Error, response.StatusCode)
	}

	return fmt.Errorf("OAuth request failed (HTTP %d)", response.StatusCode)
}

func exchangeClientCredentials(ctx context.Context, origin string, auth *config.Auth) (*config.Auth, error) {
	assertion, err := clientAssertion(auth)
	if err != nil {
		return nil, err
	}

	form := url.Values{"grant_type": {"client_credentials"}, "client_id": {auth.ClientID},
		"client_assertion_type": {"urn:ietf:params:oauth:client-assertion-type:jwt-bearer"}, "client_assertion": {assertion}}
	if auth.RequestedScope != "" {
		form.Set("scope", auth.RequestedScope)
	}

	response, err := oauthRequest(ctx, origin, auth.TokenEndpoint, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()), "")
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, oauthResponseError(response)
	}

	var token openapi.OAuthToken
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&token); err != nil {
		return nil, err
	}

	if token.AccessToken == nil || *token.AccessToken == "" || token.ExpiresIn == nil || *token.ExpiresIn <= 0 || token.TokenType == nil || !strings.EqualFold(*token.TokenType, "Bearer") {
		return nil, fmt.Errorf("token endpoint returned an invalid token response")
	}

	next := *auth
	next.AccessToken = *token.AccessToken
	next.TokenType = "Bearer"
	next.ExpiresAt = time.Now().Add(time.Duration(*token.ExpiresIn) * time.Second)
	if token.Scope != nil {
		next.Scope = *token.Scope
	}

	if next.Registration != nil {
		r := *next.Registration
		r.State = "ready"
		next.Registration = &r
	}

	return &next, nil
}
