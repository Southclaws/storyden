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
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/settings"
)

type DynamicRegistrationAccessTokenCreate struct {
	Label            string
	CreatorAccountID account.AccountID
	ExpiresAt        time.Time
	MaxRegistrations int
}

type DynamicRegistrationAccessTokenIssued struct {
	IAT   *oauthresource.DynamicRegistrationAccessToken
	Token string
}

func (s *Service) CreateDynamicRegistrationAccessToken(ctx context.Context, input DynamicRegistrationAccessTokenCreate) (*DynamicRegistrationAccessTokenIssued, error) {
	label := strings.TrimSpace(input.Label)
	if label == "" || len(label) > 200 || !input.ExpiresAt.After(time.Now()) || input.MaxRegistrations < 1 || input.MaxRegistrations > 2147483647 {
		return nil, fault.New("label, future expiry, and a positive registration limit are required", ftag.With(ftag.InvalidArgument))
	}

	secret, err := randomToken(32)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	token := oauthresource.DynamicRegistrationAccessTokenPrefix + secret
	iat, err := s.tokens.CreateDynamicRegistrationAccessToken(ctx, oauth_writer.DynamicRegistrationAccessTokenCreate{
		Label:            label,
		TokenHash:        hashString(token),
		CreatorAccountID: input.CreatorAccountID,
		ExpiresAt:        input.ExpiresAt,
		MaxRegistrations: input.MaxRegistrations,
	})
	if err != nil {
		return nil, err
	}

	return &DynamicRegistrationAccessTokenIssued{IAT: iat, Token: token}, nil
}

func (s *Service) ListDynamicRegistrationAccessTokens(ctx context.Context, params pagination.Parameters) (*pagination.Result[*oauthresource.DynamicRegistrationAccessToken], error) {
	return s.clients.ListDynamicRegistrationAccessTokens(ctx, params)
}

func (s *Service) GetDynamicRegistrationAccessToken(ctx context.Context, id oauthresource.DynamicRegistrationAccessTokenID) (*oauthresource.DynamicRegistrationAccessToken, error) {
	return s.clients.GetDynamicRegistrationAccessToken(ctx, id)
}

func (s *Service) RevokeDynamicRegistrationAccessToken(ctx context.Context, id oauthresource.DynamicRegistrationAccessTokenID) error {
	updated, err := s.tokens.RevokeDynamicRegistrationAccessToken(ctx, id, time.Now())
	if err != nil || updated {
		return err
	}

	_, err = s.clients.GetDynamicRegistrationAccessToken(ctx, id)
	return err
}

func (s *Service) authoriseAgentRegistration(ctx context.Context, mode settings.OAuthAutonomousRegistrationMode, autonomous bool, token opt.Optional[string]) (opt.Optional[oauthresource.DynamicRegistrationAccessTokenID], *Error, error) {
	var id opt.Optional[oauthresource.DynamicRegistrationAccessTokenID]
	if autonomous {
		switch mode {
		case settings.OAuthAutonomousRegistrationModeOpen, settings.OAuthAutonomousRegistrationModeApproval:
			break

		case settings.OAuthAutonomousRegistrationModeProtected:
			if !token.Ok() {
				return id, oauthError("invalid_token", "Initial Access Token is required"), nil
			}

		default:
			return id, oauthError("invalid_client_metadata", "Autonomous agent registration is not enabled"), nil
		}
	}

	raw, supplied := token.Get()
	if !supplied {
		return id, nil, nil
	}

	if !strings.HasPrefix(raw, oauthresource.DynamicRegistrationAccessTokenPrefix) {
		return id, oauthError("invalid_token", "Invalid Initial Access Token"), nil
	}

	iat, err := s.clients.GetDynamicRegistrationAccessTokenByHash(ctx, hashString(raw))
	if err != nil {
		if ftag.Get(err) == ftag.NotFound {
			return id, oauthError("invalid_token", "Invalid Initial Access Token"), nil
		}

		return id, nil, err
	}

	if iat.RevokedAt.Ok() || !iat.ExpiresAt.After(time.Now()) || iat.RegistrationCount >= iat.MaxRegistrations {
		return id, oauthError("invalid_token", "Initial Access Token is expired, revoked, or exhausted"), nil
	}

	if !autonomous {
		return id, oauthError("insufficient_scope", "Initial Access Token only permits autonomous agent registration"), nil
	}

	return opt.New(iat.ID), nil, nil
}
