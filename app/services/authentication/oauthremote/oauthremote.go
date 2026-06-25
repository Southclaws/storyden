package oauthremote

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/storyden/app/resources/account"
	oauth_remote "github.com/Southclaws/storyden/app/resources/oauth/remote"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/authentication/oauthremote/oauth_http_client"
	"github.com/Southclaws/storyden/app/services/semdex/robot/mcpclient"
	"github.com/Southclaws/storyden/internal/config"
	"go.uber.org/fx"
	"golang.org/x/oauth2"
)

const (
	discoveryTimeout = 10 * time.Second
	flowTTL          = 10 * time.Minute
	maxMetadataBytes = 256 * 1024
)

type Mode = oauth_remote.Mode

const (
	ModeCIMD   = oauth_remote.ModeCIMD
	ModeDCR    = oauth_remote.ModeDCR
	ModeManual = oauth_remote.ModeManual
)

type Service struct {
	repo     *oauth_remote.Repository
	mcp      *mcpclient.Manager
	settings *settings.SettingsRepository
	config   config.Config
	client   *http.Client
}

func Build() fx.Option {
	return fx.Provide(New)
}

func New(repo *oauth_remote.Repository, mcp *mcpclient.Manager, settings *settings.SettingsRepository, cfg config.Config) *Service {
	return &Service{
		repo:     repo,
		mcp:      mcp,
		settings: settings,
		config:   cfg,
		client:   oauth_http_client.NewHTTPClient(discoveryTimeout),
	}
}

type ProtectedResourceMetadata struct {
	Resource               string   `json:"resource"`
	ResourceName           string   `json:"resource_name,omitempty"`
	AuthorizationServers   []string `json:"authorization_servers"`
	BearerMethodsSupported []string `json:"bearer_methods_supported,omitempty"`
}

type AuthorizationServerMetadata struct {
	Issuer                            string   `json:"issuer"`
	AuthorizationEndpoint             string   `json:"authorization_endpoint"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	RegistrationEndpoint              string   `json:"registration_endpoint,omitempty"`
	ResponseTypesSupported            []string `json:"response_types_supported,omitempty"`
	GrantTypesSupported               []string `json:"grant_types_supported,omitempty"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported,omitempty"`
	CodeChallengeMethodsSupported     []string `json:"code_challenge_methods_supported,omitempty"`
	ClientIDMetadataDocumentSupported bool     `json:"client_id_metadata_document_supported,omitempty"`
}

type DiscoveryResult struct {
	ResourceURL                 string
	ProtectedResourceMetadata   ProtectedResourceMetadata
	AuthorizationServer         string
	AuthorizationServerMetadata AuthorizationServerMetadata
	Mode                        Mode
	ClientID                    string
	RedirectURI                 string
}

type ManualConfig struct {
	ClientID              string
	ClientSecret          string
	AuthorizationEndpoint string
	TokenEndpoint         string
	RedirectURI           string
	AuthorizationServer   string
	Scope                 string
}

type CreateConnectionInput struct {
	ResourceURL string
	Mode        Mode
	Manual      ManualConfig
	AddedBy     account.AccountID
	Scope       string
}

type AuthorizeResult struct {
	Connection oauth_remote.Connection
	AuthURL    string
	State      string
}

type TokenResult struct {
	Connection oauth_remote.Connection
}

type ClientMetadataDocument struct {
	ClientID                string   `json:"client_id"`
	ClientName              string   `json:"client_name"`
	RedirectURIs            []string `json:"redirect_uris"`
	GrantTypes              []string `json:"grant_types"`
	ResponseTypes           []string `json:"response_types"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method"`
	ClientURI               string   `json:"client_uri,omitempty"`
}

func (s *Service) Discover(ctx context.Context, rawResourceURL string) (DiscoveryResult, error) {
	resourceURL, err := normalizeResourceURL(rawResourceURL)
	if err != nil {
		return DiscoveryResult{}, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
	if err := oauth_http_client.ValidateRemoteOAuthURL(resourceURL.String(), "resource URL"); err != nil {
		return DiscoveryResult{}, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	protected, err := fetchJSON[ProtectedResourceMetadata](ctx, s.client, protectedResourceMetadataURL(resourceURL))
	if err != nil {
		return DiscoveryResult{}, fault.Wrap(err, fctx.With(ctx))
	}
	if len(protected.AuthorizationServers) == 0 {
		return DiscoveryResult{}, fault.Wrap(fault.New("protected resource metadata has no authorization_servers"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	authorizationServer := strings.TrimRight(protected.AuthorizationServers[0], "/")
	if err := oauth_http_client.ValidateRemoteOAuthURL(authorizationServer, "authorization server"); err != nil {
		return DiscoveryResult{}, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
	metadata, err := fetchJSON[AuthorizationServerMetadata](ctx, s.client, authorizationServerMetadataURL(authorizationServer))
	if err != nil {
		return DiscoveryResult{}, fault.Wrap(err, fctx.With(ctx))
	}
	if err := validateAuthorizationServerMetadata(metadata); err != nil {
		return DiscoveryResult{}, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	mode := selectMode(metadata)
	return DiscoveryResult{
		ResourceURL:                 resourceURL.String(),
		ProtectedResourceMetadata:   protected,
		AuthorizationServer:         authorizationServer,
		AuthorizationServerMetadata: metadata,
		Mode:                        mode,
		ClientID:                    s.ClientMetadataDocumentURL(),
		RedirectURI:                 s.DefaultRedirectURI(),
	}, nil
}

func (s *Service) ListConnections(ctx context.Context) ([]oauth_remote.Connection, error) {
	return s.repo.ListConnections(ctx)
}

func (s *Service) CreateConnection(ctx context.Context, input CreateConnectionInput) (oauth_remote.Connection, error) {
	scope := strings.TrimSpace(input.Scope)
	if scope == "" {
		scope = strings.TrimSpace(input.Manual.Scope)
	}

	if input.Mode == ModeManual {
		resourceURL, err := normalizeResourceURL(input.ResourceURL)
		if err != nil {
			return oauth_remote.Connection{}, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		if err := oauth_http_client.ValidateRemoteOAuthURL(resourceURL.String(), "resource URL"); err != nil {
			return oauth_remote.Connection{}, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		if strings.TrimSpace(input.Manual.ClientID) == "" {
			return oauth_remote.Connection{}, fault.Wrap(fault.New("manual OAuth configuration requires client_id"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		if strings.TrimSpace(input.Manual.AuthorizationEndpoint) == "" {
			return oauth_remote.Connection{}, fault.Wrap(fault.New("manual OAuth configuration requires authorization_endpoint"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		if strings.TrimSpace(input.Manual.TokenEndpoint) == "" {
			return oauth_remote.Connection{}, fault.Wrap(fault.New("manual OAuth configuration requires token_endpoint"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		redirectURI := input.Manual.RedirectURI
		if redirectURI == "" {
			redirectURI = s.DefaultRedirectURI()
		}
		if err := oauth_http_client.ValidateRemoteOAuthURL(input.Manual.AuthorizationEndpoint, "authorization endpoint"); err != nil {
			return oauth_remote.Connection{}, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		if err := oauth_http_client.ValidateRemoteOAuthURL(input.Manual.TokenEndpoint, "token endpoint"); err != nil {
			return oauth_remote.Connection{}, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		if err := oauth_http_client.ValidateRemoteOAuthURL(redirectURI, "redirect URI"); err != nil {
			return oauth_remote.Connection{}, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		tokenAuthMethod := "none"
		if input.Manual.ClientSecret != "" {
			tokenAuthMethod = "client_secret_basic"
		}
		if err := s.rejectExistingConnectedConnection(ctx, resourceURL.String(), input.Manual.AuthorizationServer, input.AddedBy); err != nil {
			return oauth_remote.Connection{}, err
		}
		return s.createReplacingUnconnected(ctx, oauth_remote.ConnectionCreate{
			ResourceURL:             resourceURL.String(),
			AuthorizationServer:     input.Manual.AuthorizationServer,
			Mode:                    ModeManual,
			Status:                  oauth_remote.StatusPending,
			ClientID:                input.Manual.ClientID,
			ClientSecret:            input.Manual.ClientSecret,
			AuthorizationEndpoint:   input.Manual.AuthorizationEndpoint,
			TokenEndpoint:           input.Manual.TokenEndpoint,
			TokenEndpointAuthMethod: tokenAuthMethod,
			RedirectURI:             redirectURI,
			RedirectURIs:            []string{redirectURI},
			Scope:                   scope,
			AddedBy:                 input.AddedBy,
		})
	}

	discovered, err := s.Discover(ctx, input.ResourceURL)
	if err != nil {
		return oauth_remote.Connection{}, err
	}
	mode := discovered.Mode
	if input.Mode != "" {
		mode = input.Mode
	}
	if mode == ModeDCR && strings.TrimSpace(discovered.AuthorizationServerMetadata.RegistrationEndpoint) == "" {
		return oauth_remote.Connection{}, fault.Wrap(fault.New("authorization server does not expose a registration_endpoint"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
	if mode == ModeCIMD && !oauth_http_client.SupportsTokenEndpointAuth(discovered.AuthorizationServerMetadata.TokenEndpointAuthMethodsSupported, "none") {
		return oauth_remote.Connection{}, fault.Wrap(fault.New("authorization server does not support public-client token authentication for CIMD"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
	if err := s.rejectExistingConnectedConnection(ctx, discovered.ResourceURL, discovered.AuthorizationServer, input.AddedBy); err != nil {
		return oauth_remote.Connection{}, err
	}

	clientID := discovered.ClientID
	clientSecret := ""
	tokenAuthMethod := "none"
	redirectURIs := []string{discovered.RedirectURI}

	if mode == ModeDCR {
		if err := oauth_http_client.ValidateRemoteOAuthURL(discovered.AuthorizationServerMetadata.RegistrationEndpoint, "registration endpoint"); err != nil {
			return oauth_remote.Connection{}, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		set, err := s.settings.Get(ctx)
		if err != nil {
			return oauth_remote.Connection{}, err
		}
		registered, err := registerDynamicClient(ctx, s.client, discovered.AuthorizationServerMetadata.RegistrationEndpoint, dynamicClientMetadata{
			RedirectURIs:            redirectURIs,
			TokenEndpointAuthMethod: "none",
			GrantTypes:              []string{"authorization_code", "refresh_token"},
			ResponseTypes:           []string{"code"},
			ClientName:              set.Title.Or(settings.DefaultTitle),
			ClientURI:               strings.TrimRight(s.config.PublicWebAddress.String(), "/"),
			Scope:                   scope,
		})
		if err != nil {
			return oauth_remote.Connection{}, fault.Wrap(err, fctx.With(ctx))
		}
		clientID = registered.ClientID
		clientSecret = registered.ClientSecret
		tokenAuthMethod = registered.TokenEndpointAuthMethod
		if tokenAuthMethod == "" {
			if clientSecret != "" {
				tokenAuthMethod = "client_secret_basic"
			} else {
				tokenAuthMethod = "none"
			}
		}
		if !oauth_http_client.SupportedTokenEndpointAuthMethod(tokenAuthMethod) {
			return oauth_remote.Connection{}, fault.Wrap(fault.New("dynamic client registration returned unsupported token_endpoint_auth_method"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		if len(registered.RedirectURIs) > 0 {
			if !oauth_http_client.Contains(registered.RedirectURIs, redirectURIs[0]) {
				return oauth_remote.Connection{}, fault.Wrap(fault.New("dynamic client registration response did not preserve requested redirect_uri"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
			}
			redirectURIs = registered.RedirectURIs
		}
	}
	if len(redirectURIs) == 0 {
		return oauth_remote.Connection{}, fault.Wrap(fault.New("OAuth connection requires at least one redirect URI"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
	for _, redirectURI := range redirectURIs {
		if err := oauth_http_client.ValidateRemoteOAuthURL(redirectURI, "redirect URI"); err != nil {
			return oauth_remote.Connection{}, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
	}

	return s.createReplacingUnconnected(ctx, oauth_remote.ConnectionCreate{
		ResourceURL:                 discovered.ResourceURL,
		Resource:                    discovered.ProtectedResourceMetadata.Resource,
		ResourceName:                discovered.ProtectedResourceMetadata.ResourceName,
		ProtectedResourceMetadata:   structToMap(discovered.ProtectedResourceMetadata),
		AuthorizationServer:         discovered.AuthorizationServer,
		AuthorizationServerMetadata: structToMap(discovered.AuthorizationServerMetadata),
		Mode:                        mode,
		Status:                      oauth_remote.StatusPending,
		ClientID:                    clientID,
		ClientSecret:                clientSecret,
		AuthorizationEndpoint:       discovered.AuthorizationServerMetadata.AuthorizationEndpoint,
		TokenEndpoint:               discovered.AuthorizationServerMetadata.TokenEndpoint,
		RegistrationEndpoint:        discovered.AuthorizationServerMetadata.RegistrationEndpoint,
		TokenEndpointAuthMethod:     tokenAuthMethod,
		RedirectURIs:                redirectURIs,
		RedirectURI:                 redirectURIs[0],
		Scope:                       scope,
		AddedBy:                     input.AddedBy,
	})
}

func (s *Service) rejectExistingConnectedConnection(ctx context.Context, resourceURL string, authorizationServer string, addedBy account.AccountID) error {
	exists, err := s.repo.HasConnectedConnectionByIdentity(ctx, resourceURL, authorizationServer, addedBy)
	if err != nil {
		return err
	}
	if exists {
		return fault.Wrap(fault.New("remote OAuth connection already exists"), fctx.With(ctx), ftag.With(ftag.AlreadyExists))
	}
	return nil
}

func (s *Service) createReplacingUnconnected(ctx context.Context, in oauth_remote.ConnectionCreate) (oauth_remote.Connection, error) {
	if _, err := s.repo.DeleteUnconnectedConnectionByIdentity(ctx, in.ResourceURL, in.AuthorizationServer, in.AddedBy); err != nil {
		return oauth_remote.Connection{}, err
	}

	return s.repo.CreateConnection(ctx, in)
}

func (s *Service) StartAuthorization(ctx context.Context, id oauth_remote.ConnectionID) (AuthorizeResult, error) {
	connection, err := s.repo.GetConnection(ctx, id)
	if err != nil {
		return AuthorizeResult{}, err
	}
	if connection.AuthorizationEndpoint == "" || connection.TokenEndpoint == "" || connection.ClientID == "" {
		return AuthorizeResult{}, fault.Wrap(fault.New("remote OAuth connection is missing authorization configuration"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
	if err := oauth_http_client.ValidateRemoteOAuthURL(connection.AuthorizationEndpoint, "authorization endpoint"); err != nil {
		return AuthorizeResult{}, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
	if err := oauth_http_client.ValidateRemoteOAuthURL(connection.TokenEndpoint, "token endpoint"); err != nil {
		return AuthorizeResult{}, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	state, err := randomURLToken(32)
	if err != nil {
		return AuthorizeResult{}, err
	}
	verifier, err := randomURLToken(48)
	if err != nil {
		return AuthorizeResult{}, err
	}

	redirectURI := connection.RedirectURI
	if redirectURI == "" {
		redirectURI = s.DefaultRedirectURI()
	}
	if err := oauth_http_client.ValidateRemoteOAuthURL(redirectURI, "redirect URI"); err != nil {
		return AuthorizeResult{}, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	if _, err := s.repo.CreateFlow(ctx, connection.ID, stateHash(state), verifier, redirectURI, time.Now().Add(flowTTL)); err != nil {
		return AuthorizeResult{}, err
	}

	cfg := oauth2.Config{
		ClientID:    connection.ClientID,
		RedirectURL: redirectURI,
		Scopes:      oauth_http_client.SplitScope(connection.Scope),
		Endpoint:    oauth_http_client.Endpoint(connection),
	}
	authURL := cfg.AuthCodeURL(state,
		oauth2.SetAuthURLParam("code_challenge", pkceChallenge(verifier)),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	return AuthorizeResult{Connection: connection, AuthURL: authURL, State: state}, nil
}

func (s *Service) HandleCallback(ctx context.Context, state string, code string) (TokenResult, error) {
	if strings.TrimSpace(state) == "" || strings.TrimSpace(code) == "" {
		return TokenResult{}, fault.Wrap(fault.New("missing OAuth state or code"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	flow, err := s.repo.GetFlowByStateHash(ctx, stateHash(state))
	if err != nil {
		return TokenResult{}, err
	}
	if flow.Connection == nil {
		return TokenResult{}, fault.Wrap(fault.New("OAuth flow has no connection"), fctx.With(ctx))
	}
	now := time.Now()
	if flow.ConsumedAt != nil || now.After(flow.ExpiresAt) {
		return TokenResult{}, fault.Wrap(fault.New("OAuth flow expired or already used"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
	claimed, err := s.repo.ClaimFlow(ctx, flow.ID, now)
	if err != nil {
		return TokenResult{}, err
	}
	if !claimed {
		return TokenResult{}, fault.Wrap(fault.New("OAuth flow expired or already used"), fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	connection := *flow.Connection
	cfg := oauth2.Config{
		ClientID:     connection.ClientID,
		ClientSecret: connection.ClientSecret,
		RedirectURL:  flow.RedirectURI,
		Scopes:       oauth_http_client.SplitScope(connection.Scope),
		Endpoint:     oauth_http_client.Endpoint(connection),
	}

	exchangeCtx := oauth_http_client.ContextWithHTTPClient(ctx, s.client)
	token, err := cfg.Exchange(exchangeCtx, code, oauth2.SetAuthURLParam("code_verifier", flow.PKCEVerifier))
	if err != nil {
		_ = s.repo.MarkError(context.WithoutCancel(ctx), connection.ID, err.Error())
		return TokenResult{}, fault.Wrap(err, fctx.With(ctx))
	}
	if token.AccessToken == "" {
		err := fault.New("OAuth token response missing access_token")
		_ = s.repo.MarkError(context.WithoutCancel(ctx), connection.ID, err.Error())
		return TokenResult{}, fault.Wrap(err, fctx.With(ctx))
	}
	if token.TokenType == "" {
		err := fault.New("OAuth token response missing token_type")
		_ = s.repo.MarkError(context.WithoutCancel(ctx), connection.ID, err.Error())
		return TokenResult{}, fault.Wrap(err, fctx.With(ctx))
	}

	var expiry *time.Time
	if !token.Expiry.IsZero() {
		expiry = &token.Expiry
	}
	updated, err := s.repo.StoreTokens(ctx, connection.ID, oauth_remote.TokenUpdate{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		TokenExpiry:  expiry,
		Scope:        oauth_http_client.StringExtra(token, "scope"),
	})
	if err != nil {
		return TokenResult{}, err
	}
	refreshCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 30*time.Second)
	defer cancel()
	_ = s.mcp.RefreshOAuthConnectionServers(refreshCtx, updated.ID)

	return TokenResult{Connection: updated}, nil
}

func (s *Service) ClientMetadataDocument(ctx context.Context) (ClientMetadataDocument, error) {
	set, err := s.settings.Get(ctx)
	if err != nil {
		return ClientMetadataDocument{}, err
	}
	return ClientMetadataDocument{
		ClientID:                s.ClientMetadataDocumentURL(),
		ClientName:              set.Title.Or(settings.DefaultTitle),
		RedirectURIs:            []string{s.DefaultRedirectURI()},
		GrantTypes:              []string{"authorization_code", "refresh_token"},
		ResponseTypes:           []string{"code"},
		TokenEndpointAuthMethod: "none",
		ClientURI:               strings.TrimRight(s.config.PublicWebAddress.String(), "/"),
	}, nil
}

func (s *Service) ClientMetadataDocumentURL() string {
	base := s.config.PublicAPIAddress
	base.Path = "/.well-known/oauth-client-metadata"
	base.RawQuery = ""
	base.Fragment = ""
	return base.String()
}

func (s *Service) DefaultRedirectURI() string {
	base := s.config.PublicWebAddress
	base.Path = "/oauth/remote/callback"
	base.RawQuery = ""
	base.Fragment = ""
	return base.String()
}

func selectMode(metadata AuthorizationServerMetadata) Mode {
	if metadata.ClientIDMetadataDocumentSupported {
		return ModeCIMD
	}
	if strings.TrimSpace(metadata.RegistrationEndpoint) != "" {
		return ModeDCR
	}
	return ModeManual
}

func validateAuthorizationServerMetadata(metadata AuthorizationServerMetadata) error {
	if strings.TrimSpace(metadata.Issuer) == "" {
		return fault.New("authorization server metadata missing issuer")
	}
	if err := oauth_http_client.ValidateRemoteOAuthURL(metadata.Issuer, "issuer"); err != nil {
		return err
	}
	if err := oauth_http_client.ValidateRemoteOAuthURL(metadata.AuthorizationEndpoint, "authorization endpoint"); err != nil {
		return err
	}
	if err := oauth_http_client.ValidateRemoteOAuthURL(metadata.TokenEndpoint, "token endpoint"); err != nil {
		return err
	}
	if metadata.RegistrationEndpoint != "" {
		if err := oauth_http_client.ValidateRemoteOAuthURL(metadata.RegistrationEndpoint, "registration endpoint"); err != nil {
			return err
		}
	}
	if len(metadata.ResponseTypesSupported) > 0 && !oauth_http_client.Contains(metadata.ResponseTypesSupported, "code") {
		return fault.New("authorization server does not support authorization code response type")
	}
	if len(metadata.GrantTypesSupported) > 0 && !oauth_http_client.Contains(metadata.GrantTypesSupported, "authorization_code") {
		return fault.New("authorization server does not support authorization_code grant")
	}
	if len(metadata.CodeChallengeMethodsSupported) > 0 && !oauth_http_client.Contains(metadata.CodeChallengeMethodsSupported, "S256") {
		return fault.New("authorization server does not support S256 PKCE")
	}
	return nil
}

func normalizeResourceURL(raw string) (*url.URL, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return nil, err
	}
	if u.Scheme != "https" && u.Scheme != "http" {
		return nil, fault.New("resource URL must use http or https")
	}
	if u.Host == "" {
		return nil, fault.New("resource URL host is required")
	}
	if u.User != nil || u.Fragment != "" {
		return nil, fault.New("resource URL must not include user info or fragment")
	}
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)
	if u.Path == "" {
		u.Path = "/"
	}
	return u, nil
}

func protectedResourceMetadataURL(resource *url.URL) string {
	return wellKnownURL(resource, "oauth-protected-resource")
}

func authorizationServerMetadataURL(authorizationServer string) string {
	u, err := url.Parse(strings.TrimSpace(authorizationServer))
	if err != nil || u.Scheme == "" || u.Host == "" {
		return strings.TrimRight(authorizationServer, "/") + "/.well-known/oauth-authorization-server"
	}
	return wellKnownURL(u, "oauth-authorization-server")
}

func wellKnownURL(base *url.URL, name string) string {
	u := *base
	basePath := strings.Trim(u.Path, "/")
	u.Path = "/.well-known/" + name
	if basePath != "" {
		u.Path += "/" + basePath
	}
	u.RawQuery = ""
	u.Fragment = ""
	return u.String()
}

func fetchJSON[T any](ctx context.Context, client *http.Client, url string) (T, error) {
	var out T
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return out, err
	}
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return out, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return out, fmt.Errorf("metadata discovery failed: %s", res.Status)
	}
	return out, json.NewDecoder(io.LimitReader(res.Body, maxMetadataBytes)).Decode(&out)
}

func randomURLToken(size int) (string, error) {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func stateHash(state string) string {
	sum := sha256.Sum256([]byte(state))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func pkceChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func structToMap(in any) map[string]any {
	body, err := json.Marshal(in)
	if err != nil {
		return map[string]any{}
	}
	var out map[string]any
	if err := json.Unmarshal(body, &out); err != nil {
		return map[string]any{}
	}
	return out
}
