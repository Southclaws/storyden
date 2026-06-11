package oauth

import (
	"context"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/opt"
	"github.com/golang-jwt/jwt/v5"

	"github.com/Southclaws/storyden/app/resources/account"
	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_writer"
)

func (s *Service) issueTokens(ctx context.Context, client *oauthresource.Client, accountID account.AccountID, scope string, nonce opt.Optional[string]) (*Token, error) {
	resp, _, err := s.issueTokensWithRefresh(ctx, client, accountID, scope, nonce)
	return resp, err
}

func (s *Service) issueTokensWithRefresh(ctx context.Context, client *oauthresource.Client, accountID account.AccountID, scope string, nonce opt.Optional[string]) (*Token, oauthresource.RefreshTokenID, error) {
	now := time.Now()
	exp := now.Add(s.cfg.OAuthAccessTokenTTL)

	jti, err := randomToken(12)
	if err != nil {
		return nil, oauthresource.RefreshTokenID{}, err
	}

	accessToken, err := s.sign(jwt.MapClaims{
		"iss":   s.issuer,
		"sub":   accountID.String(),
		"aud":   "storyden",
		"exp":   exp.Unix(),
		"iat":   now.Unix(),
		"jti":   jti,
		"scope": scope,
	})
	if err != nil {
		return nil, oauthresource.RefreshTokenID{}, err
	}

	resp := &Token{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(s.cfg.OAuthAccessTokenTTL.Seconds()),
		Scope:       scope,
	}

	scopes := splitScope(scope)
	if contains(scopes, "openid") {
		claims, err := s.idTokenClaims(ctx, accountID, client.ClientID, scopes, nonce, exp, now)
		if err != nil {
			return nil, oauthresource.RefreshTokenID{}, err
		}

		idToken, err := s.sign(claims)
		if err != nil {
			return nil, oauthresource.RefreshTokenID{}, err
		}
		resp.IDToken = opt.New(idToken)
	}

	var refreshTokenID oauthresource.RefreshTokenID
	if contains(scopes, "offline_access") && contains(client.AllowedGrants, GrantTypeRefreshToken) {
		raw, err := randomToken(32)
		if err != nil {
			return nil, oauthresource.RefreshTokenID{}, err
		}
		rt, err := s.tokens.CreateRefreshToken(ctx, oauth_writer.RefreshTokenCreate{
			ClientID:  client.ID,
			AccountID: accountID,
			TokenHash: hashString(raw),
			Scope:     scope,
			ExpiresAt: now.Add(s.cfg.OAuthRefreshTokenTTL),
		})
		if err != nil {
			return nil, oauthresource.RefreshTokenID{}, err
		}
		resp.RefreshToken = opt.New(raw)
		refreshTokenID = rt.ID
	}

	return resp, refreshTokenID, nil
}

func (s *Service) idTokenClaims(ctx context.Context, accountID account.AccountID, audience string, scopes []string, nonce opt.Optional[string], exp, now time.Time) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{
		"iss": s.issuer,
		"sub": accountID.String(),
		"aud": audience,
		"exp": exp.Unix(),
		"iat": now.Unix(),
	}

	nonce.Call(func(value string) {
		claims["nonce"] = value
	})

	if !contains(scopes, "profile") && !contains(scopes, "email") {
		return claims, nil
	}

	acc, err := s.account.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if contains(scopes, "profile") {
		claims["name"] = acc.Name
	}

	if contains(scopes, "email") {
		email := ""
		emailVerified := false
		for _, address := range acc.EmailAddresses {
			email = address.Email.Address
			emailVerified = address.Verified
			if address.Verified {
				break
			}
		}
		claims["email"] = email
		claims["email_verified"] = emailVerified
	}

	return claims, nil
}

func (s *Service) sign(claims jwt.MapClaims) (string, error) {
	if s.signer == nil {
		return "", fault.New("oauth signing key is not configured")
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	t.Header["kid"] = s.kid
	return t.SignedString(s.signer)
}
