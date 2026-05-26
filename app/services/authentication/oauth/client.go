package oauth

import (
	"context"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
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
	RedirectURIs       []string
	AllowedScopes      []string
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

	clientIDToken, err := randomToken(18)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	clientSecret := opt.NewEmpty[string]()
	clientSecretHash := opt.NewEmpty[string]()

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

	client, err := s.tokens.CreateClient(ctx, oauth_writer.ClientCreate{
		AccountID:        opt.New(input.AccountID),
		ClientID:         oauthresource.OAuthAccessKeyPrefix + clientIDToken,
		ClientSecretHash: clientSecretHash,
		Name:             input.Name,
		Type:             oauthresource.ClientTypeConfidential,
		ScopePolicy:      opt.New(oauthresource.ScopePolicyExplicit),
		RedirectURIs:     input.RedirectURIs,
		AllowedScopes:    input.AllowedScopes,
		AllowedGrants: []string{
			GrantTypeClientCredentials,
		},
	})
	if err != nil {
		return nil, err
	}

	return &ClientSelfCreateResult{
		Client:       client,
		ClientSecret: clientSecret,
	}, nil
}
