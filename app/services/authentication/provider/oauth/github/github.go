package github

import (
	"context"

	"github.com/Southclaws/fault"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth/all"
	"github.com/Southclaws/storyden/internal/config"
)

var provider = authentication.ServiceOAuthGitHub

type Provider struct {
	auth_repo authentication.Repository

	callback string
	config   *all.Configuration
}

func New(cfg config.Config, auth_repo authentication.Repository) (*Provider, error) {
	config, err := all.LoadProvider(provider)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &Provider{
		auth_repo: auth_repo,
		config:    config,
		callback:  all.Redirect(cfg, provider),
	}, nil
}

func (p *Provider) Provides() authentication.Service { return provider }

func (p *Provider) Enabled(ctx context.Context) (bool, error) {
	return p.config != nil, nil
}

func (p *Provider) Link(_ string) (string, error) {
	c := oauth2.Config{
		ClientID:     p.config.ClientID,
		ClientSecret: p.config.ClientSecret,
		Endpoint:     github.Endpoint,
		RedirectURL:  p.callback,
		Scopes:       []string{},
	}

	return c.AuthCodeURL("", oauth2.AccessTypeOffline), nil
}

func (p *Provider) Login(ctx context.Context, state, code string) (*account.Account, error) {
	return nil, nil
}
