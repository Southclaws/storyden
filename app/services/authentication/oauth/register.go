package oauth

import (
	"context"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/alexedwards/argon2id"

	"github.com/Southclaws/storyden/app/resources/account"
	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
)

const (
	TokenEndpointAuthMethodNone              = "none"
	TokenEndpointAuthMethodClientSecretBasic = "client_secret_basic"
	TokenEndpointAuthMethodClientSecretPost  = "client_secret_post"
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

// DynamicClientRegistrationResult is the resolved RFC 7591 client information
// response, including the issued client and (for confidential clients) the
// one-time client secret.
//
// ClientSecretExpiresAt is always 0 per RFC 7591, indicating the secret does
// not expire. Storyden does not implement client secret rotation; secrets
// remain valid until the client is deleted or manually rotated via the
// management API.
type DynamicClientRegistrationResult struct {
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
	if !s.Enabled() || !s.cfg.OAuthDynamicRegistrationEnabled {
		return nil, oauthError("temporarily_unavailable", "Dynamic client registration is not enabled"), nil
	}

	authMethod, clientType, oauthErr := resolveTokenEndpointAuthMethod(input.TokenEndpointAuthMethod)
	if oauthErr != nil {
		return nil, oauthErr, nil
	}

	grantTypes, oauthErr := s.resolveDCRGrantTypes(input.GrantTypes)
	if oauthErr != nil {
		return nil, oauthErr, nil
	}

	// OAuth security invariant: client_credentials requires a confidential client.
	// A public client has no client authentication, so it must not receive
	// app-level tokens via client_credentials.
	if contains(grantTypes, GrantTypeClientCredentials) && clientType == oauthresource.ClientTypePublic {
		return nil, oauthError("invalid_client_metadata", "client_credentials grant requires a confidential client"), nil
	}

	// refresh_token grant should only be allowed with authorization_code
	if contains(grantTypes, GrantTypeRefreshToken) && !contains(grantTypes, GrantTypeAuthorizationCode) {
		return nil, oauthError("invalid_client_metadata", "refresh_token grant requires authorization_code grant"), nil
	}

	// Resolve response types (with defaults)
	responseTypes, oauthErr := resolveDCRResponseTypes(input.ResponseTypes, grantTypes)
	if oauthErr != nil {
		return nil, oauthErr, nil
	}

	// RFC 7591 Section 2.1: Validate grant_types and response_types consistency
	// authorization_code grant requires "code" response type
	if contains(grantTypes, GrantTypeAuthorizationCode) && !contains(responseTypes, "code") {
		return nil, oauthError("invalid_client_metadata", "authorization_code grant requires 'code' response type"), nil
	}

	// "code" response type requires authorization_code grant
	if contains(responseTypes, "code") && !contains(grantTypes, GrantTypeAuthorizationCode) {
		return nil, oauthError("invalid_client_metadata", "'code' response type requires authorization_code grant"), nil
	}

	redirectURIs, oauthErr := validateDCRRedirectURIs(input.RedirectURIs, grantTypes)
	if oauthErr != nil {
		return nil, oauthErr, nil
	}

	// Validate metadata URIs if provided
	if input.LogoURI != "" {
		if err := validateMetadataURI(input.LogoURI); err != nil {
			return nil, oauthError("invalid_client_metadata", "logo_uri must be a valid HTTPS URL"), nil
		}
	}
	if input.ClientURI != "" {
		if err := validateMetadataURI(input.ClientURI); err != nil {
			return nil, oauthError("invalid_client_metadata", "client_uri must be a valid HTTPS URL"), nil
		}
	}
	if input.TOSURI != "" {
		if err := validateMetadataURI(input.TOSURI); err != nil {
			return nil, oauthError("invalid_client_metadata", "tos_uri must be a valid HTTPS URL"), nil
		}
	}
	if input.PolicyURI != "" {
		if err := validateMetadataURI(input.PolicyURI); err != nil {
			return nil, oauthError("invalid_client_metadata", "policy_uri must be a valid HTTPS URL"), nil
		}
	}

	scopes, oauthErr := resolveDCRScopes(input.Scope)
	if oauthErr != nil {
		return nil, oauthErr, nil
	}

	clientIDToken, err := randomToken(18)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}
	clientID := oauthresource.OAuthAccessKeyPrefix + clientIDToken

	clientSecret := opt.NewEmpty[string]()
	clientSecretHash := opt.NewEmpty[string]()
	if clientType == oauthresource.ClientTypeConfidential {
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

	client, err := s.tokens.CreateClient(ctx, oauth_writer.ClientCreate{
		AccountID:               opt.NewEmpty[account.AccountID](),
		ClientID:                clientID,
		ClientSecretHash:        clientSecretHash,
		Name:                    name,
		Type:                    clientType,
		ScopePolicy:             opt.New(oauthresource.ScopePolicyExplicit),
		TokenEndpointAuthMethod: opt.New(authMethod),
		PKCERequired:            opt.New(pkceRequired),
		RedirectURIs:            redirectURIs,
		AllowedScopes:           scopes,
		AllowedGrants:           grantTypes,
	})
	if err != nil {
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
	default:
		return "", oauthresource.ClientType{}, oauthError("invalid_client_metadata", "Unsupported token_endpoint_auth_method")
	}
}

func (s *Service) resolveDCRGrantTypes(requested []string) ([]string, *Error) {
	if len(requested) == 0 {
		return []string{GrantTypeAuthorizationCode, GrantTypeRefreshToken}, nil
	}

	// Only authorization_code and refresh_token are allowed for DCR.
	// client_credentials doesn't make sense for dynamically registered clients
	// since DCR is for user-facing OAuth flows, not machine-to-machine auth.
	allowed := map[string]struct{}{
		GrantTypeAuthorizationCode: {},
		GrantTypeRefreshToken:      {},
	}

	seen := map[string]struct{}{}
	out := []string{}
	for _, grant := range requested {
		grant = strings.TrimSpace(grant)
		if grant == "" {
			continue
		}
		if _, ok := allowed[grant]; !ok {
			return nil, oauthError("invalid_client_metadata", "Unsupported grant type; only authorization_code and refresh_token are allowed")
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
