package oauth

import (
	"context"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
)

type ClientUpdate struct {
	Name             opt.Optional[string]
	ClientSecretHash opt.Optional[string]
	ScopePolicy      opt.Optional[oauthresource.ScopePolicy]
	RedirectURIs     opt.Optional[[]string]
	AllowedScopes    opt.Optional[[]string]
	AllowedGrants    opt.Optional[[]string]
}

type ClientSelfUpdate struct {
	AccountPermissions rbac.Permissions
	Name               opt.Optional[string]
	RedirectURIs       opt.Optional[[]string]
	AllowedScopes      opt.Optional[[]string]
}

func (s *Service) ListClients(ctx context.Context) ([]*oauthresource.Client, error) {
	return s.clients.ListClients(ctx)
}

func (s *Service) ListClientsByAccount(ctx context.Context, accountID account.AccountID) ([]*oauthresource.Client, error) {
	return s.clients.ListClientsByAccount(ctx, accountID)
}

func (s *Service) GetClient(ctx context.Context, id oauthresource.ClientID) (*oauthresource.Client, error) {
	return s.clients.GetClient(ctx, id)
}

func (s *Service) GetClientByAccount(ctx context.Context, accountID account.AccountID, id oauthresource.ClientID) (*oauthresource.Client, *Error, error) {
	client, err := s.clients.GetClient(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if owner, ok := client.AccountID.Get(); !ok || owner != accountID {
		return nil, oauthError("invalid_request", "Client does not belong to the requesting account"), nil
	}

	return client, nil, nil
}

func (s *Service) UpdateClient(ctx context.Context, id oauthresource.ClientID, input ClientUpdate) (*oauthresource.Client, error) {
	return s.tokens.UpdateClient(ctx, id, oauth_writer.ClientUpdate{
		Name:             input.Name,
		ClientSecretHash: input.ClientSecretHash,
		ScopePolicy:      input.ScopePolicy,
		RedirectURIs:     input.RedirectURIs,
		AllowedScopes:    input.AllowedScopes,
		AllowedGrants:    input.AllowedGrants,
	})
}

func (s *Service) UpdateClientByAccount(ctx context.Context, accountID account.AccountID, id oauthresource.ClientID, input ClientSelfUpdate) (*oauthresource.Client, *Error, error) {
	_, oauthErr, err := s.GetClientByAccount(ctx, accountID, id)
	if err != nil || oauthErr != nil {
		return nil, oauthErr, err
	}

	if scopes, ok := input.AllowedScopes.Get(); ok {
		if err := validatePermissionOnlyScopes(scopes); err != nil {
			return nil, nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		if err := validatePermissionScopes(strings.Join(scopes, " "), input.AccountPermissions); err != nil {
			return nil, nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.PermissionDenied))
		}
	}

	client, err := s.tokens.UpdateClient(ctx, id, oauth_writer.ClientUpdate{
		Name:          input.Name,
		RedirectURIs:  input.RedirectURIs,
		AllowedScopes: input.AllowedScopes,
	})
	if err != nil {
		return nil, nil, err
	}

	return client, nil, nil
}

func (s *Service) DeleteClient(ctx context.Context, id oauthresource.ClientID) error {
	return s.tokens.DeleteClient(ctx, id)
}

func (s *Service) DeleteClientByAccount(ctx context.Context, accountID account.AccountID, id oauthresource.ClientID) *Error {
	_, oauthErr, err := s.GetClientByAccount(ctx, accountID, id)
	if err != nil || oauthErr != nil {
		return oauthError("invalid_request", "Client not found or does not belong to account")
	}

	if err := s.tokens.DeleteClient(ctx, id); err != nil {
		return oauthError("invalid_request", "Failed to delete client")
	}

	return nil
}

func (s *Service) ListDeviceAuthorisations(ctx context.Context) ([]*oauthresource.DeviceAuthorisation, error) {
	return s.clients.ListDeviceAuthorisations(ctx)
}

func (s *Service) ListRefreshTokens(ctx context.Context) ([]*oauthresource.RefreshToken, error) {
	return s.clients.ListRefreshTokens(ctx)
}

func (s *Service) ListRefreshTokensByAccount(ctx context.Context, accountID account.AccountID) ([]*oauthresource.RefreshToken, error) {
	return s.clients.ListRefreshTokensByAccount(ctx, accountID)
}

func (s *Service) RevokeRefreshToken(ctx context.Context, id oauthresource.RefreshTokenID) error {
	_, err := s.tokens.RevokeRefreshToken(ctx, id, time.Now(), opt.NewEmpty[oauthresource.RefreshTokenID]())
	return err
}

func (s *Service) RevokeRefreshTokenByAccount(ctx context.Context, accountID account.AccountID, id oauthresource.RefreshTokenID) *Error {
	token, err := s.clients.GetRefreshToken(ctx, id)
	if err != nil || token.AccountID != accountID {
		return oauthError("invalid_request", "Refresh token not found or does not belong to account")
	}

	if _, err := s.tokens.RevokeRefreshToken(ctx, id, time.Now(), opt.NewEmpty[oauthresource.RefreshTokenID]()); err != nil {
		return oauthError("invalid_request", "Failed to revoke refresh token")
	}

	return nil
}
