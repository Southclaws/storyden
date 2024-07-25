package saml

import (
	"context"

	"github.com/Southclaws/fault"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth/all"
	"github.com/Southclaws/storyden/app/services/authentication/register"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/saml"
)

var (
	ErrAccessToken  = fault.New("failed to get access token")
	ErrMissingToken = fault.New("no access token in response")
)

const (
	id   = "saml"
	name = "SAML"
)

type SAMLProvider struct {
	auth_repo authentication.Repository
	register  register.Service
	samlsp    *saml.SAML

	callback string
	config   all.Configuration
}

func New(cfg config.Config, auth_repo authentication.Repository, register register.Service, samlsp *saml.SAML) (*SAMLProvider, error) {
	config, err := all.LoadProvider(id)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &SAMLProvider{
		auth_repo: auth_repo,
		register:  register,
		samlsp:    samlsp,
		config:    config,
		callback:  all.Redirect(cfg, id),
	}, nil
}

func (p *SAMLProvider) Enabled() bool { return p.config.Enabled }
func (p *SAMLProvider) ID() string    { return id }
func (p *SAMLProvider) Name() string  { return name }

func (p *SAMLProvider) Link(_ string) (string, error) {
	c := oauth2.Config{
		ClientID:     p.config.ClientID,
		ClientSecret: p.config.ClientSecret,
		Endpoint:     linkedin.Endpoint,
		RedirectURL:  p.callback,
		Scopes: []string{
			"r_emailaddress",
			"r_liteprofile",
		},
	}

	return c.AuthCodeURL("state", oauth2.AccessTypeOffline), nil
}

func (p *SAMLProvider) Login(ctx context.Context, state, code string) (*account.Account, error) {
	return nil, nil
}
