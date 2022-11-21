package github

import (
	"context"

	"github.com/Southclaws/fault"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/authentication"
	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth/all"
	"github.com/Southclaws/storyden/internal/config"
)

const (
	id   = "github"
	name = "GitHub"
	logo = "https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png"
)

type GitHubProvider struct {
	auth_repo authentication.Repository

	callback string
	config   all.Configuration
}

func New(cfg config.Config, auth_repo authentication.Repository) (*GitHubProvider, error) {
	config, err := all.LoadProvider(id)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &GitHubProvider{
		auth_repo: auth_repo,
		config:    config,
		callback:  all.Redirect(cfg, id),
	}, nil
}

func (p *GitHubProvider) Enabled() bool   { return p.config.Enabled }
func (p *GitHubProvider) ID() string      { return id }
func (p *GitHubProvider) Name() string    { return name }
func (p *GitHubProvider) LogoURL() string { return logo }

func (p *GitHubProvider) Link() string {
	c := oauth2.Config{
		ClientID:     p.config.ClientID,
		ClientSecret: p.config.ClientSecret,
		Endpoint:     github.Endpoint,
		RedirectURL:  p.callback,
		Scopes:       []string{},
	}

	return c.AuthCodeURL("", oauth2.AccessTypeOffline)
}

func (p *GitHubProvider) Login(ctx context.Context, state, code string) (*account.Account, error) {
	return nil, nil
}
