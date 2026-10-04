package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/alexedwards/argon2id"
	"github.com/go-jose/go-jose/v4"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/settings"
)

const (
	TokenEndpointAuthMethodNone              = "none"
	TokenEndpointAuthMethodClientSecretBasic = "client_secret_basic"
	TokenEndpointAuthMethodClientSecretPost  = "client_secret_post"
	TokenEndpointAuthMethodPrivateKeyJWT     = "private_key_jwt"
)

var dcrDefaultScopes = []string{
	"openid",
	"profile",
	"email",
	"offline_access",
}

// DynamicClientRegistration carries RFC 7591 client metadata supplied by a
// dynamically registering client.
type DynamicClientRegistration struct {
	RegistrationModes       []string
	InitialAccessToken      opt.Optional[string]
	ClientName              string
	RedirectURIs            []string
	GrantTypes              []string
	ResponseTypes           []string
	Scope                   string
	TokenEndpointAuthMethod string
	ApplicationType         string
	LogoURI                 string
	ClientURI               string
	TOSURI                  string
	PolicyURI               string
	JWKs                    map[string]any
}

// DynamicClientRegistrationResult is the resolved RFC 7591 client information
// response, including the issued client and (for confidential clients) the
// one-time client secret.
//
// ClientSecretExpiresAt is always 0 per RFC 7591, indicating the secret does
// not expire. Storyden does not implement client secret rotation; secrets
// remain valid until the client is deleted or manually rotated via the
// management API.
type DynamicClientRegistrationResult struct {
	Approval                *RegistrationApprovalChallenge
	RetryAfter              int
	Client                  *oauthresource.Client
	ClientSecret            opt.Optional[string]
	ClientIDIssuedAt        int64
	ClientSecretExpiresAt   int64 // Always 0 (secrets never expire)
	ClientName              string
	RedirectURIs            []string
	GrantTypes              []string
	ResponseTypes           []string
	Scope                   string
	TokenEndpointAuthMethod string
	ApplicationType         string
	LogoURI                 string
	ClientURI               string
	TOSURI                  string
	PolicyURI               string
}

// RegisterClient implements RFC 7591 Dynamic Client Registration.
//
// Dynamically registered clients are tenant-owned: they have no account owner.
// This is the least invasive safe model on top of the existing client table -
// the account_id column is already optional, so a NULL owner cleanly marks a
// client as dynamically registered (admin and member clients always have an
// owner). They use the explicit scope policy so they can only ever obtain the
// scopes they were registered with, intersected with the authorising user's
// permissions at authorize time.
//
// PKCE is required for authorization-code clients but needs no extra storage
// here: the authorize and token endpoints already enforce PKCE (S256) for every
// client unconditionally.
func (s *Service) RegisterClient(ctx context.Context, input DynamicClientRegistration) (*DynamicClientRegistrationResult, *Error, error) {
	if !s.Enabled() {
		return nil, oauthError("temporarily_unavailable", "Dynamic client registration is not enabled"), nil
	}

	configuration, err := s.registrationSettings(ctx)
	if err != nil {
		return nil, nil, err
	}

	if !configuration.DynamicRegistrationEnabled.OrZero() {
		return nil, oauthError("temporarily_unavailable", "Dynamic client registration is not enabled"), nil
	}

	prepared, oauthErr := prepareClientRegistration(ctx, input)
	if oauthErr != nil {
		return nil, oauthErr, nil
	}

	autonomousAgent := contains(prepared.GrantTypes, GrantTypeClientCredentials)
	iatID, oauthErr, err := s.authoriseAgentRegistration(ctx, configuration.AutonomousRegistrationMode.OrZero(), autonomousAgent, input.InitialAccessToken)
	if oauthErr != nil || err != nil {
		return nil, oauthErr, err
	}

	if autonomousAgent && configuration.AutonomousRegistrationMode.OrZero() == settings.OAuthAutonomousRegistrationModeApproval && !iatID.Ok() {
		if !contains(input.RegistrationModes, "approval") {
			return nil, oauthError("invalid_client_metadata", "Autonomous registration requires approval; include approval in registration_mode"), nil
		}

		return s.startRegistrationApproval(ctx, configuration, *prepared)
	}

	return s.provisionClientRegistration(ctx, configuration, *prepared, iatID, nil)
}

func prepareClientRegistration(ctx context.Context, input DynamicClientRegistration) (*DynamicClientRegistration, *Error) {
	authMethod, _, oauthErr := resolveTokenEndpointAuthMethod(input.TokenEndpointAuthMethod)
	if oauthErr != nil {
		return nil, oauthErr
	}

	input.TokenEndpointAuthMethod = authMethod
	grantTypes, oauthErr := resolveDCRGrantTypes(input.GrantTypes)
	if oauthErr != nil {
		return nil, oauthErr
	}

	autonomousAgent := contains(grantTypes, GrantTypeClientCredentials)
	if autonomousAgent && (len(grantTypes) != 1 || authMethod != TokenEndpointAuthMethodPrivateKeyJWT) {
		return nil, oauthError("invalid_client_metadata", "Autonomous registration requires only client_credentials and private_key_jwt")
	}

	if oauthErr := validateDCRClientMetadata(ctx, input, autonomousAgent); oauthErr != nil {
		return nil, oauthErr
	}

	// refresh_token grant should only be allowed with authorization_code
	if contains(grantTypes, GrantTypeRefreshToken) && !contains(grantTypes, GrantTypeAuthorizationCode) {
		return nil, oauthError("invalid_client_metadata", "refresh_token grant requires authorization_code grant")
	}

	// Resolve response types (with defaults)
	responseTypes, oauthErr := resolveDCRResponseTypes(input.ResponseTypes, grantTypes)
	if oauthErr != nil {
		return nil, oauthErr
	}

	// RFC 7591 Section 2.1: Validate grant_types and response_types consistency
	// authorization_code grant requires "code" response type
	if contains(grantTypes, GrantTypeAuthorizationCode) && !contains(responseTypes, "code") {
		return nil, oauthError("invalid_client_metadata", "authorization_code grant requires 'code' response type")
	}

	// "code" response type requires authorization_code grant
	if contains(responseTypes, "code") && !contains(grantTypes, GrantTypeAuthorizationCode) {
		return nil, oauthError("invalid_client_metadata", "'code' response type requires authorization_code grant")
	}

	redirectURIs, oauthErr := validateDCRRedirectURIs(input.RedirectURIs, grantTypes)
	if oauthErr != nil {
		return nil, oauthErr
	}

	scopes, oauthErr := resolveDCRScopes(input.Scope)
	if oauthErr != nil {
		return nil, oauthErr
	}

	input.GrantTypes = grantTypes
	input.ResponseTypes = responseTypes
	input.RedirectURIs = redirectURIs
	if input.Scope == "" && contains(grantTypes, GrantTypeClientCredentials) {
		scopes = nil
	}

	input.Scope = strings.Join(scopes, " ")
	input.ClientName = strings.TrimSpace(input.ClientName)
	input.InitialAccessToken = opt.NewEmpty[string]()
	return &input, nil
}

func (s *Service) provisionClientRegistration(ctx context.Context, configuration settings.OAuthServiceSettings, input DynamicClientRegistration, iatID opt.Optional[oauthresource.DynamicRegistrationAccessTokenID], approval *oauthresource.RegistrationApproval) (*DynamicClientRegistrationResult, *Error, error) {
	authMethod, clientType, _ := resolveTokenEndpointAuthMethod(input.TokenEndpointAuthMethod)
	grantTypes, responseTypes, redirectURIs := input.GrantTypes, input.ResponseTypes, input.RedirectURIs
	scopes := strings.Fields(input.Scope)
	autonomousAgent := contains(grantTypes, GrantTypeClientCredentials)
	clientIDToken, err := randomToken(18)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	clientID := oauthresource.OAuthAccessKeyPrefix + clientIDToken

	clientSecret := opt.NewEmpty[string]()
	clientSecretHash := opt.NewEmpty[string]()
	if clientType == oauthresource.ClientTypeConfidential && authMethod != TokenEndpointAuthMethodPrivateKeyJWT {
		secretToken, err := randomToken(32)
		if err != nil {
			return nil, nil, fault.Wrap(err, fctx.With(ctx))
		}

		secret := oauthresource.OAuthAccessSecretPrefix + secretToken

		hash, err := argon2id.CreateHash(secret, argon2id.DefaultParams)
		if err != nil {
			return nil, nil, fault.Wrap(err, fctx.With(ctx))
		}

		clientSecret = opt.New(secret)
		clientSecretHash = opt.New(hash)
	}

	name := strings.TrimSpace(input.ClientName)
	if name == "" {
		name = clientID
	}

	// Authorization-code clients registered dynamically must use PKCE.
	// This applies to both public and confidential clients for defense-in-depth.
	// OAuth 2.1 and RFC 7636 recommend PKCE for all authorization_code flows.
	pkceRequired := contains(grantTypes, GrantTypeAuthorizationCode)

	createInput := oauth_writer.ClientCreate{
		RegistrationAccessTokenID: iatID,
		RegistrationApproval:      approval,
		AccountID:                 opt.NewEmpty[account.AccountID](),
		ClientID:                  clientID,
		ClientSecretHash:          clientSecretHash,
		Name:                      name,
		Type:                      clientType,
		ScopePolicy:               opt.New(oauthresource.ScopePolicyExplicit),
		TokenEndpointAuthMethod:   opt.New(authMethod),
		JWKs:                      opt.New(input.JWKs),
		PKCERequired:              opt.New(pkceRequired),
		RedirectURIs:              redirectURIs,
		AllowedScopes:             scopes,
		AllowedGrants:             grantTypes,
	}

	var client *oauthresource.Client
	if autonomousAgent {
		createInput.ScopePolicy = opt.New(oauthresource.ScopePolicyInheritUserPermissions)
		createInput.RegistrationRoleID = configuration.AutonomousRegistrationRoleID
		client, err = s.provisionAgentClient(ctx, name, createInput)
	} else {
		client, err = s.tokens.CreateClient(ctx, createInput)
	}

	if err != nil {
		if errors.Is(err, oauth_writer.ErrRegistrationApprovalUnavailable) {
			return nil, invalidRegistrationError(), nil
		}

		if errors.Is(err, oauth_writer.ErrDynamicRegistrationAccessTokenUnavailable) {
			return nil, oauthError("invalid_token", "Initial Access Token is expired, revoked, or exhausted"), nil
		}

		if autonomousAgent && ftag.Get(err) == ftag.AlreadyExists {
			return nil, oauthError("invalid_client_metadata", "client_name is already in use"), nil
		}

		return nil, nil, err
	}

	return &DynamicClientRegistrationResult{
		Client:                  client,
		ClientSecret:            clientSecret,
		ClientIDIssuedAt:        time.Now().Unix(),
		ClientSecretExpiresAt:   0,
		ClientName:              name,
		RedirectURIs:            redirectURIs,
		GrantTypes:              grantTypes,
		ResponseTypes:           responseTypes,
		Scope:                   strings.Join(scopes, " "),
		TokenEndpointAuthMethod: authMethod,
		ApplicationType:         input.ApplicationType,
		LogoURI:                 input.LogoURI,
		ClientURI:               input.ClientURI,
		TOSURI:                  input.TOSURI,
		PolicyURI:               input.PolicyURI,
	}, nil, nil
}

func (s *Service) provisionAgentClient(ctx context.Context, handle string, input oauth_writer.ClientCreate) (*oauthresource.Client, error) {
	id := account.AccountID(xid.New())
	if err := s.avatars.SetRobot(ctx, id); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	input.AccountID = opt.New(id)
	client, err := s.tokens.CreateAgentClient(ctx, handle, input)
	if err != nil {
		cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
		defer cancel()
		if cleanupErr := s.avatars.Delete(cleanupCtx, id); cleanupErr != nil {
			return nil, fault.Wrap(errors.Join(err, cleanupErr), fctx.With(ctx))
		}

		return nil, err
	}

	return client, nil
}

func validateDCRClientMetadata(ctx context.Context, input DynamicClientRegistration, autonomousAgent bool) *Error {
	if autonomousAgent {
		if err := account.ValidateHandle(ctx, strings.TrimSpace(input.ClientName)); err != nil {
			return oauthError("invalid_client_metadata", "client_name must be a valid, available Storyden handle")
		}
	}

	if input.TokenEndpointAuthMethod == TokenEndpointAuthMethodPrivateKeyJWT && !validPublicJWKSet(input.JWKs) {
		return oauthError("invalid_client_metadata", "jwks must contain valid public signing keys")
	}

	for _, metadata := range []struct{ name, value string }{
		{"logo_uri", input.LogoURI},
		{"client_uri", input.ClientURI},
		{"tos_uri", input.TOSURI},
		{"policy_uri", input.PolicyURI},
	} {
		if metadata.value == "" {
			continue
		}

		if err := validateMetadataURI(metadata.value); err != nil {
			return oauthError("invalid_client_metadata", metadata.name+" must be a valid HTTPS URL")
		}
	}

	return nil
}

func resolveTokenEndpointAuthMethod(method string) (string, oauthresource.ClientType, *Error) {
	method = strings.TrimSpace(method)
	if method == "" {
		method = TokenEndpointAuthMethodClientSecretBasic
	}

	switch method {
	case TokenEndpointAuthMethodNone:
		return method, oauthresource.ClientTypePublic, nil
	case TokenEndpointAuthMethodClientSecretPost:
		return method, oauthresource.ClientTypeConfidential, nil
	case TokenEndpointAuthMethodClientSecretBasic:
		return method, oauthresource.ClientTypeConfidential, nil
	case TokenEndpointAuthMethodPrivateKeyJWT:
		return method, oauthresource.ClientTypeConfidential, nil
	default:
		return "", oauthresource.ClientType{}, oauthError("invalid_client_metadata", "Unsupported token_endpoint_auth_method")
	}
}

func resolveDCRGrantTypes(requested []string) ([]string, *Error) {
	if len(requested) == 0 {
		return []string{GrantTypeAuthorizationCode, GrantTypeRefreshToken}, nil
	}

	allowed := map[string]struct{}{
		GrantTypeAuthorizationCode: {},
		GrantTypeRefreshToken:      {},
		GrantTypeClientCredentials: {},
	}

	seen := map[string]struct{}{}
	out := []string{}
	for _, grant := range requested {
		grant = strings.TrimSpace(grant)
		if grant == "" {
			continue
		}

		if _, ok := allowed[grant]; !ok {
			return nil, oauthError("invalid_client_metadata", "Unsupported grant type")
		}

		if _, ok := seen[grant]; ok {
			continue
		}

		seen[grant] = struct{}{}
		out = append(out, grant)
	}

	if len(out) == 0 {
		return []string{GrantTypeAuthorizationCode, GrantTypeRefreshToken}, nil
	}

	return out, nil
}

func validPublicJWKSet(jwks map[string]any) bool {
	raw, err := json.Marshal(jwks)
	if err != nil {
		return false
	}

	var keys jose.JSONWebKeySet
	if err := json.Unmarshal(raw, &keys); err != nil || len(keys.Keys) == 0 {
		return false
	}

	seen := map[string]bool{}
	for _, key := range keys.Keys {
		if !key.Valid() || !key.IsPublic() || (key.Use != "" && key.Use != "sig") || seen[key.KeyID] {
			return false
		}

		seen[key.KeyID] = true
	}

	return true
}

func resolveDCRResponseTypes(requested []string, grantTypes []string) ([]string, *Error) {
	usesAuthCode := contains(grantTypes, GrantTypeAuthorizationCode)

	if len(requested) == 0 {
		if usesAuthCode {
			return []string{"code"}, nil
		}

		return []string{}, nil
	}

	out := []string{}
	for _, responseType := range requested {
		responseType = strings.TrimSpace(responseType)
		if responseType == "" {
			continue
		}

		if responseType != "code" {
			return nil, oauthError("invalid_client_metadata", "Only 'code' response type is supported")
		}

		out = append(out, responseType)
	}

	if contains(out, "code") && !usesAuthCode {
		return nil, oauthError("invalid_client_metadata", "'code' response type requires authorization_code grant")
	}

	return out, nil
}

func validateDCRRedirectURIs(redirectURIs []string, grantTypes []string) ([]string, *Error) {
	usesAuthCode := contains(grantTypes, GrantTypeAuthorizationCode)

	if len(redirectURIs) == 0 {
		if usesAuthCode {
			return nil, oauthError("invalid_redirect_uri", "At least one redirect_uri is required for authorization_code clients")
		}

		return []string{}, nil
	}

	// Deduplicate and validate redirect URIs
	seen := make(map[string]struct{}, len(redirectURIs))
	out := make([]string, 0, len(redirectURIs))
	for _, raw := range redirectURIs {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}

		if err := validateDCRRedirectURI(raw); err != nil {
			return nil, oauthError("invalid_redirect_uri", "Invalid redirect_uri: must be an absolute HTTPS URI (or HTTP for loopback)")
		}

		// Deduplicate
		if _, exists := seen[raw]; exists {
			continue
		}

		seen[raw] = struct{}{}
		out = append(out, raw)
	}

	// After deduplication, check if we still have URIs when required
	if len(out) == 0 && usesAuthCode {
		return nil, oauthError("invalid_redirect_uri", "At least one valid redirect_uri is required for authorization_code clients")
	}

	return out, nil
}

func validateDCRRedirectURI(raw string) error {
	if raw == "" {
		return fault.New("empty redirect uri")
	}

	if strings.Contains(raw, "*") {
		return fault.New("wildcard redirect uri")
	}

	u, err := url.Parse(raw)
	if err != nil {
		return err
	}

	if !u.IsAbs() {
		return fault.New("redirect uri must be absolute")
	}

	// Reject opaque URIs like "https:callback" that have no authority/hostname
	if u.Opaque != "" || u.Hostname() == "" {
		return fault.New("redirect uri must have a valid hostname")
	}

	if u.Fragment != "" || strings.Contains(raw, "#") {
		return fault.New("redirect uri must not contain a fragment")
	}

	switch u.Scheme {
	case "https":
		return nil
	case "http":
		// Loopback redirect URIs are permitted for native/dev clients per the
		// OAuth security best current practice.
		if isLoopbackHost(u.Hostname()) {
			return nil
		}

		return fault.New("http redirect uri only allowed for loopback hosts")
	default:
		return fault.New("redirect uri must use https")
	}
}

func isLoopbackHost(host string) bool {
	// Handle "localhost" as a special case
	if host == "localhost" {
		return true
	}

	// Parse as IP and check if it's a loopback address
	// Covers 127.0.0.0/8 for IPv4 and ::1/128 for IPv6
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}

	return ip.IsLoopback()
}

func resolveDCRScopes(scope string) ([]string, *Error) {
	requested := splitScope(scope)
	if len(requested) == 0 {
		return dcrDefaultScopes, nil
	}

	seen := map[string]struct{}{}
	out := []string{}
	for _, sc := range requested {
		if _, ok := standardScopes[sc]; !ok {
			// TODO: Make allowed scopes configurable. It seems we can't change
			// the allowed scopes based on registration method (manual vs DCR)
			// and most clients appear to just request every scope in the well-
			// known set. For now, just allow all scopes but this should change.
			// Or, alternatively, we can silently strip out ADMINISTRATOR just
			// for DCR clients instead of returning a 400 Bad Request error.
			_, err := rbac.NewPermission(sc)
			if err != nil {
				return nil, oauthError("invalid_client_metadata", "Scope is not permitted for dynamic client registration")
			}

		}

		if _, ok := seen[sc]; ok {
			continue
		}

		seen[sc] = struct{}{}
		out = append(out, sc)
	}

	return out, nil
}

// validateMetadataURI validates URIs used for client metadata (logo_uri, client_uri, etc.)
// These must be absolute HTTPS URIs to prevent phishing and drive-by downloads.
func validateMetadataURI(raw string) error {
	if raw == "" {
		return nil
	}

	u, err := url.Parse(raw)
	if err != nil {
		return fault.Wrap(err)
	}

	if !u.IsAbs() {
		return fault.New("metadata uri must be absolute")
	}

	if u.Scheme != "https" {
		return fault.New("metadata uri must use https")
	}

	return nil
}
