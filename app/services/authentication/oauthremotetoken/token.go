package oauthremotetoken

import (
	"context"
	"net/http"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	oauth_remote "github.com/Southclaws/storyden/app/resources/oauth/remote"
	"github.com/Southclaws/storyden/app/services/authentication/oauthremote/oauth_http_client"
	"go.uber.org/fx"
	"golang.org/x/oauth2"
)

const (
	refreshHTTPTimeout = 10 * time.Second
	refreshSkew        = 5 * time.Minute
	refreshLeaseTTL    = 30 * time.Second
	refreshPoll        = 100 * time.Millisecond
	refreshWaitTimeout = 5 * time.Second
)

type Service struct {
	repo   *oauth_remote.Repository
	client *http.Client
}

func Build() fx.Option {
	return fx.Provide(New)
}

func New(repo *oauth_remote.Repository) *Service {
	return &Service{
		repo:   repo,
		client: oauth_http_client.NewHTTPClient(refreshHTTPTimeout),
	}
}

func (s *Service) AccessToken(ctx context.Context, id oauth_remote.ConnectionID) (string, error) {
	waitUntil := time.Now().Add(refreshWaitTimeout)

	for {
		connection, err := s.repo.GetConnection(ctx, id)
		if err != nil {
			return "", err
		}

		token, fresh, err := currentAccessToken(connection, time.Now())
		if err != nil {
			return "", err
		}
		if fresh {
			return token, nil
		}

		now := time.Now()
		claimed, err := s.repo.ClaimTokenRefresh(ctx, id, now, now.Add(-refreshLeaseTTL))
		if err != nil {
			return "", err
		}
		if claimed {
			return s.refreshClaimed(ctx, id)
		}

		if time.Now().After(waitUntil) {
			return "", fault.Wrap(fault.New("timed out waiting for remote OAuth token refresh"), fctx.With(ctx))
		}

		timer := time.NewTimer(refreshPoll)
		select {
		case <-ctx.Done():
			timer.Stop()
			return "", fault.Wrap(ctx.Err(), fctx.With(ctx))
		case <-timer.C:
		}
	}
}

func (s *Service) refreshClaimed(ctx context.Context, id oauth_remote.ConnectionID) (string, error) {
	releaseLease := true
	defer func() {
		if releaseLease {
			_ = s.repo.ReleaseTokenRefresh(context.WithoutCancel(ctx), id)
		}
	}()

	connection, err := s.repo.GetConnection(ctx, id)
	if err != nil {
		return "", err
	}

	token, fresh, err := currentAccessToken(connection, time.Now())
	if err != nil {
		return "", err
	}
	if fresh {
		return token, nil
	}

	if connection.RefreshToken == "" {
		err := fault.New("remote OAuth access token is expired and no refresh token is available")
		_ = s.repo.MarkError(context.WithoutCancel(ctx), id, err.Error())
		releaseLease = false
		return "", fault.Wrap(err, fctx.With(ctx))
	}
	if connection.TokenEndpoint == "" || connection.ClientID == "" {
		err := fault.New("remote OAuth connection is missing token refresh configuration")
		_ = s.repo.MarkError(context.WithoutCancel(ctx), id, err.Error())
		releaseLease = false
		return "", fault.Wrap(err, fctx.With(ctx))
	}
	if err := oauth_http_client.ValidateRemoteOAuthURL(connection.TokenEndpoint, "token endpoint"); err != nil {
		_ = s.repo.MarkError(context.WithoutCancel(ctx), id, err.Error())
		releaseLease = false
		return "", fault.Wrap(err, fctx.With(ctx))
	}

	cfg := oauth2.Config{
		ClientID:     connection.ClientID,
		ClientSecret: connection.ClientSecret,
		Endpoint:     oauth_http_client.Endpoint(connection),
		Scopes:       oauth_http_client.SplitScope(connection.Scope),
	}

	refreshCtx := oauth_http_client.ContextWithHTTPClient(ctx, s.client)
	refreshed, err := cfg.TokenSource(refreshCtx, &oauth2.Token{RefreshToken: connection.RefreshToken}).Token()
	if err != nil {
		_ = s.repo.MarkError(context.WithoutCancel(ctx), id, err.Error())
		releaseLease = false
		return "", fault.Wrap(err, fctx.With(ctx))
	}
	if refreshed.AccessToken == "" {
		err := fault.New("OAuth token refresh response missing access_token")
		_ = s.repo.MarkError(context.WithoutCancel(ctx), id, err.Error())
		releaseLease = false
		return "", fault.Wrap(err, fctx.With(ctx))
	}
	if refreshed.TokenType == "" {
		err := fault.New("OAuth token refresh response missing token_type")
		_ = s.repo.MarkError(context.WithoutCancel(ctx), id, err.Error())
		releaseLease = false
		return "", fault.Wrap(err, fctx.With(ctx))
	}
	if refreshed.RefreshToken == "" {
		refreshed.RefreshToken = connection.RefreshToken
	}

	var expiry *time.Time
	if !refreshed.Expiry.IsZero() {
		expiry = &refreshed.Expiry
	}

	updated, err := s.repo.StoreTokens(ctx, id, oauth_remote.TokenUpdate{
		AccessToken:  refreshed.AccessToken,
		RefreshToken: refreshed.RefreshToken,
		TokenType:    refreshed.TokenType,
		TokenExpiry:  expiry,
		Scope:        oauth_http_client.StringExtra(refreshed, "scope"),
	})
	if err != nil {
		return "", err
	}

	releaseLease = false
	return updated.AccessToken, nil
}

func currentAccessToken(connection oauth_remote.Connection, now time.Time) (string, bool, error) {
	if connection.AccessToken == "" {
		return "", false, fault.New("remote OAuth connection has no access token")
	}
	if connection.TokenExpiry == nil {
		return connection.AccessToken, true, nil
	}
	if connection.TokenExpiry.After(now.Add(refreshSkew)) {
		return connection.AccessToken, true, nil
	}
	return "", false, nil
}
