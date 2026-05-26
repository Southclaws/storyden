package oauth

import (
	"context"
	"crypto/sha256"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/alexedwards/argon2id"

	"github.com/Southclaws/storyden/app/resources/account"
	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/rbac"
)

type TokenRequest struct {
	GrantType    string
	ClientID     string
	ClientSecret opt.Optional[string]
	Scope        opt.Optional[string]
	DeviceCode   opt.Optional[string]
	Code         opt.Optional[string]
	RedirectURI  opt.Optional[string]
	CodeVerifier opt.Optional[string]
	RefreshToken opt.Optional[string]
}

type Token struct {
	AccessToken  string
	TokenType    string
	ExpiresIn    int
	Scope        string
	IDToken      opt.Optional[string]
	RefreshToken opt.Optional[string]
}

func (s *Service) ExchangeToken(ctx context.Context, input TokenRequest) (*Token, *Error, error) {
	if !s.Enabled() {
		return nil, oauthError("temporarily_unavailable"), nil
	}

	switch input.GrantType {
	case GrantTypeDeviceCode:
		return s.exchangeDeviceCode(ctx, input)
	case GrantTypeAuthorizationCode:
		return s.exchangeAuthorizationCode(ctx, input)
	case GrantTypeRefreshToken:
		return s.exchangeRefreshToken(ctx, input)
	case GrantTypeClientCredentials:
		return s.exchangeClientCredentials(ctx, input)
	default:
		return nil, oauthError("unsupported_grant_type"), nil
	}
}

func (s *Service) exchangeClientCredentials(ctx context.Context, input TokenRequest) (*Token, *Error, error) {
	cl, err := s.clients.GetClientByClientID(ctx, input.ClientID)
	if err != nil {
		return nil, oauthError("invalid_client"), nil
	}
	if cl.Type != oauthresource.ClientTypeConfidential {
		return nil, oauthError("unauthorized_client"), nil
	}
	if !contains(cl.AllowedGrants, GrantTypeClientCredentials) {
		return nil, oauthError("unauthorized_client"), nil
	}
	if oauthErr, err := s.authenticateConfidentialClient(ctx, cl, input.ClientSecret); oauthErr != nil || err != nil {
		return nil, oauthErr, err
	}

	accountID, ok := cl.AccountID.Get()
	if !ok {
		return nil, oauthError("unauthorized_client"), nil
	}

	accountPermissions, err := s.accountPermissions(ctx, accountID)
	if err != nil {
		return nil, nil, err
	}

	grantedScope, err := grantClientCredentialsScope(input.Scope.OrZero(), cl, accountPermissions)
	if err != nil {
		return nil, oauthError("invalid_scope"), nil
	}

	token, err := s.issueTokens(ctx, cl, accountID, grantedScope)
	if err != nil {
		return nil, nil, err
	}

	return token, nil, nil
}

func (s *Service) exchangeAuthorizationCode(ctx context.Context, input TokenRequest) (*Token, *Error, error) {
	code, ok := input.Code.Get()
	if !ok {
		return nil, oauthError("invalid_request"), nil
	}
	redirectURI, ok := input.RedirectURI.Get()
	if !ok {
		return nil, oauthError("invalid_request"), nil
	}
	codeVerifier, ok := input.CodeVerifier.Get()
	if !ok {
		return nil, oauthError("invalid_request"), nil
	}

	cl, err := s.clients.GetClientByClientID(ctx, input.ClientID)
	if err != nil {
		return nil, oauthError("invalid_client"), nil
	}
	if !contains(cl.AllowedGrants, GrantTypeAuthorizationCode) {
		return nil, oauthError("unauthorized_client"), nil
	}
	if oauthErr, err := s.authenticateConfidentialClient(ctx, cl, input.ClientSecret); oauthErr != nil || err != nil {
		return nil, oauthErr, err
	}
	if !validCodeVerifier(codeVerifier) {
		return nil, oauthError("invalid_request"), nil
	}

	rec, err := s.clients.GetAuthorisationCodeByCodeHash(ctx, hashString(code))
	if err != nil || rec.ClientID != cl.ID || rec.RedirectURI != redirectURI || rec.ExpiresAt.Before(time.Now()) || rec.ConsumedAt.Ok() {
		return nil, oauthError("invalid_grant"), nil
	}

	sum := sha256.Sum256([]byte(codeVerifier))
	sumb64 := b64url(sum[:])
	if sumb64 != rec.CodeChallenge {
		return nil, oauthError("invalid_grant"), nil
	}

	consumed, err := s.tokens.ConsumeAuthorisationCode(ctx, rec.ID, time.Now())
	if err != nil {
		return nil, nil, err
	}
	if !consumed {
		return nil, oauthError("invalid_grant"), nil
	}

	token, err := s.issueTokens(ctx, cl, rec.AccountID, rec.Scope)
	if err != nil {
		return nil, nil, err
	}

	return token, nil, nil
}

func (s *Service) exchangeRefreshToken(ctx context.Context, input TokenRequest) (*Token, *Error, error) {
	refreshToken, ok := input.RefreshToken.Get()
	if !ok {
		return nil, oauthError("invalid_request"), nil
	}

	cl, err := s.clients.GetClientByClientID(ctx, input.ClientID)
	if err != nil {
		return nil, oauthError("invalid_client"), nil
	}
	if !contains(cl.AllowedGrants, GrantTypeRefreshToken) {
		return nil, oauthError("unauthorized_client"), nil
	}
	if oauthErr, err := s.authenticateConfidentialClient(ctx, cl, input.ClientSecret); oauthErr != nil || err != nil {
		return nil, oauthErr, err
	}

	rec, err := s.clients.GetRefreshTokenByTokenHash(ctx, hashString(refreshToken))
	if err != nil || rec.ClientID != cl.ID || rec.ExpiresAt.Before(time.Now()) || rec.RevokedAt.Ok() {
		return nil, oauthError("invalid_grant"), nil
	}

	accountPermissions, err := s.accountPermissions(ctx, rec.AccountID)
	if err != nil {
		return nil, nil, err
	}

	grantedScope, err := refreshScope(rec.Scope, cl, accountPermissions)
	if err != nil {
		return nil, oauthError("invalid_grant"), nil
	}

	consumedAt := time.Now()
	consumed, err := s.tokens.ConsumeRefreshToken(ctx, rec.ID, consumedAt)
	if err != nil {
		return nil, nil, err
	}
	if !consumed {
		return nil, oauthError("invalid_grant"), nil
	}

	token, newID, err := s.issueTokensWithRefresh(ctx, cl, rec.AccountID, grantedScope)
	if err != nil {
		return nil, nil, err
	}
	if newID != (oauthresource.RefreshTokenID{}) {
		if err := s.tokens.SetRefreshTokenReplacement(ctx, rec.ID, newID); err != nil {
			return nil, nil, err
		}
	}

	return token, nil, nil
}

func (s *Service) accountPermissions(ctx context.Context, accountID account.AccountID) (rbac.Permissions, error) {
	acc, err := s.account.GetRefByID(ctx, accountID)
	if err != nil {
		return rbac.Permissions{}, err
	}

	return acc.Roles.Permissions(), nil
}

func (s *Service) authenticateConfidentialClient(ctx context.Context, client *oauthresource.Client, secret opt.Optional[string]) (*Error, error) {
	if client.Type != oauthresource.ClientTypeConfidential {
		return nil, nil
	}

	hash, ok := client.ClientSecretHash.Get()
	if !ok {
		return oauthError("invalid_client"), nil
	}

	raw, ok := secret.Get()
	if !ok {
		return oauthError("invalid_client"), nil
	}

	match, _, err := argon2id.CheckHash(raw, hash)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if !match {
		return oauthError("invalid_client"), nil
	}

	return nil, nil
}

func validCodeVerifier(v string) bool {
	if len(v) < 43 || len(v) > 128 {
		return false
	}

	for _, c := range v {
		if c >= 'A' && c <= 'Z' {
			continue
		}
		if c >= 'a' && c <= 'z' {
			continue
		}
		if c >= '0' && c <= '9' {
			continue
		}
		switch c {
		case '-', '.', '_', '~':
			continue
		default:
			return false
		}
	}

	return true
}
