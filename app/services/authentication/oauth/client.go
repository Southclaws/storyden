package oauth

import (
	"context"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/alexedwards/argon2id"

	"github.com/Southclaws/storyden/app/resources/account"
	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
)

type ClientCreate struct {
	AccountID          account.AccountID
	AccountPermissions rbac.Permissions
	ClientID           string
	ClientSecretHash   opt.Optional[string]
	Name               string
	Type               oauthresource.ClientType
	ScopePolicy        opt.Optional[oauthresource.ScopePolicy]
	RedirectURIs       []string
	AllowedScopes      []string
	AllowedGrants      []string
}

type ClientSelfCreate struct {
	AccountID          account.AccountID
	AccountPermissions rbac.Permissions
	Name               string
	Type               oauthresource.ClientType
	RedirectURIs       []string
	AllowedScopes      []string
	AllowedGrants      []string
	PKCERequired       opt.Optional[bool]
}

type ClientSelfCreateResult struct {
	Client       *oauthresource.Client
	ClientSecret opt.Optional[string]
}

func (s *Service) CreateClient(ctx context.Context, input ClientCreate) (*oauthresource.Client, error) {
	if err := validatePermissionScopes(strings.Join(input.AllowedScopes, " "), input.AccountPermissions); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.PermissionDenied))
	}

	return s.tokens.CreateClient(ctx, oauth_writer.ClientCreate{
		AccountID:        opt.New(input.AccountID),
		ClientID:         input.ClientID,
		ClientSecretHash: input.ClientSecretHash,
		Name:             input.Name,
		Type:             input.Type,
		ScopePolicy:      input.ScopePolicy,
		RedirectURIs:     input.RedirectURIs,
		AllowedScopes:    input.AllowedScopes,
		AllowedGrants:    input.AllowedGrants,
	})
}

func (s *Service) CreateClientForAccount(ctx context.Context, input ClientSelfCreate) (*ClientSelfCreateResult, error) {
	if err := validatePermissionOnlyScopes(input.AllowedScopes); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
	if err := validatePermissionScopes(strings.Join(input.AllowedScopes, " "), input.AccountPermissions); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.PermissionDenied))
	}

	// Validate grant types
	if err := validateGrantTypes(input.AllowedGrants); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	// Validate redirect URIs for authorization code grant
	if hasGrant(input.AllowedGrants, GrantTypeAuthorizationCode) && len(input.RedirectURIs) == 0 {
		return nil, fault.Wrap(
			fault.New("redirect_uris required for authorization_code grant"),
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
		)
	}

	// Validate client credentials grant is confidential only
	if hasGrant(input.AllowedGrants, GrantTypeClientCredentials) && input.Type == oauthresource.ClientTypePublic {
		return nil, fault.Wrap(
			fault.New("client_credentials grant requires confidential client type"),
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
		)
	}

	// Validate public clients require PKCE for authorization_code
	if input.Type == oauthresource.ClientTypePublic && hasGrant(input.AllowedGrants, GrantTypeAuthorizationCode) {
		pkceRequired := input.PKCERequired.OrZero()
		if !pkceRequired {
			return nil, fault.Wrap(
				fault.New("public clients with authorization_code must require PKCE"),
				fctx.With(ctx),
				ftag.With(ftag.InvalidArgument),
			)
		}
	}

	// Validate machine clients should not have redirect URIs
	if hasGrant(input.AllowedGrants, GrantTypeClientCredentials) && !hasGrant(input.AllowedGrants, GrantTypeAuthorizationCode) && len(input.RedirectURIs) > 0 {
		return nil, fault.Wrap(
			fault.New("machine clients (client_credentials only) should not have redirect URIs"),
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
		)
	}

	clientIDToken, err := randomToken(18)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	clientSecret := opt.NewEmpty[string]()
	clientSecretHash := opt.NewEmpty[string]()

	// Only generate client secret for confidential clients
	if input.Type == oauthresource.ClientTypeConfidential {
		secretToken, err := randomToken(32)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		secret := oauthresource.OAuthAccessSecretPrefix + secretToken

		hash, err := argon2id.CreateHash(secret, argon2id.DefaultParams)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		clientSecret = opt.New(secret)
		clientSecretHash = opt.New(hash)
	}

	// Determine token endpoint auth method based on client type
	authMethod := "client_secret_post"
	if input.Type == oauthresource.ClientTypePublic {
		authMethod = "none"
	}

	client, err := s.tokens.CreateClient(ctx, oauth_writer.ClientCreate{
		AccountID:               opt.New(input.AccountID),
		ClientID:                oauthresource.OAuthAccessKeyPrefix + clientIDToken,
		ClientSecretHash:        clientSecretHash,
		Name:                    input.Name,
		Type:                    input.Type,
		ScopePolicy:             opt.New(oauthresource.ScopePolicyExplicit),
		TokenEndpointAuthMethod: opt.New(authMethod),
		PKCERequired:            input.PKCERequired,
		RedirectURIs:            input.RedirectURIs,
		AllowedScopes:           input.AllowedScopes,
		AllowedGrants:           input.AllowedGrants,
	})
	if err != nil {
		return nil, err
	}

	return &ClientSelfCreateResult{
		Client:       client,
		ClientSecret: clientSecret,
	}, nil
}

func validateGrantTypes(grants []string) error {
	if len(grants) == 0 {
		return fault.New("invalid grants", fmsg.WithDesc("no grants specified", "Allowed Grants field must contain at least one grant"))
	}

	validGrants := map[string]bool{
		GrantTypeAuthorizationCode: true,
		GrantTypeRefreshToken:      true,
		GrantTypeClientCredentials: true,
		GrantTypeDeviceCode:        true,
	}

	for _, grant := range grants {
		if !validGrants[grant] {
			return fault.Newf("invalid grant type: %s", grant)
		}
	}

	return nil
}

func hasGrant(grants []string, grant string) bool {
	for _, g := range grants {
		if g == grant {
			return true
		}
	}
	return false
}
