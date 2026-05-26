package oauth

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/rbac"
)

type AccessTokenClaims struct {
	AccountID   account.AccountID
	Permissions rbac.Permissions
	Scopes      []string
}

func (s *Service) ValidateAccessToken(ctx context.Context, raw string) (*AccessTokenClaims, error) {
	if s.signer == nil {
		return nil, fault.New("oauth signing key is not configured", fctx.With(ctx), ftag.With(ftag.Unauthenticated))
	}

	claims := jwt.MapClaims{}
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
		jwt.WithIssuer(s.issuer),
		jwt.WithAudience("storyden"),
	)

	_, err := parser.ParseWithClaims(raw, claims, func(token *jwt.Token) (any, error) {
		return &s.signer.PublicKey, nil
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Unauthenticated))
	}

	subject, err := claims.GetSubject()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Unauthenticated))
	}

	id, err := xid.FromString(subject)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Unauthenticated))
	}

	permissions, err := permissionsFromScopeClaim(claims)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Unauthenticated))
	}
	scopes, err := scopesFromScopeClaim(claims)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Unauthenticated))
	}

	return &AccessTokenClaims{
		AccountID:   account.AccountID(id),
		Permissions: permissions,
		Scopes:      scopes,
	}, nil
}

func permissionsFromScopeClaim(claims jwt.MapClaims) (rbac.Permissions, error) {
	scopes, err := scopesFromScopeClaim(claims)
	if err != nil {
		return rbac.Permissions{}, err
	}

	return permissionsFromScopes(scopes)
}

func scopesFromScopeClaim(claims jwt.MapClaims) ([]string, error) {
	raw, ok := claims["scope"]
	if !ok {
		return nil, nil
	}

	scope, ok := raw.(string)
	if !ok {
		return nil, fault.New("invalid scope claim")
	}

	return splitScope(scope), nil
}
