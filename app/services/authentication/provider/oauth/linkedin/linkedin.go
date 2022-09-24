package linkedin

import (
	"context"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/authentication"
	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth/all"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/errctx"
)

const (
	id   = "linkedin"
	name = "LinkedIn"
	logo = "https://brand.linkedin.com/content/dam/me/business/en-us/amp/brand-site/v2/bg/LI-Bug.svg.original.svg"
)

type LinkedInProvider struct {
	auth_repo authentication.Repository

	callback string
	config   all.Configuration
}

func New(cfg config.Config, auth_repo authentication.Repository) (*LinkedInProvider, error) {
	config, err := all.LoadProvider(id)
	if err != nil {
		return nil, err
	}

	return &LinkedInProvider{
		auth_repo: auth_repo,
		config:    config,
		callback:  all.Redirect(cfg, id),
	}, nil
}

func (p *LinkedInProvider) Enabled() bool   { return p.config.Enabled }
func (p *LinkedInProvider) ID() string      { return id }
func (p *LinkedInProvider) Name() string    { return name }
func (p *LinkedInProvider) LogoURL() string { return logo }

func (p *LinkedInProvider) Link() string {
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

	return c.AuthCodeURL("state", oauth2.AccessTypeOffline)
}

func (p *LinkedInProvider) Login(ctx context.Context, state, code string) (*account.Account, error) {
	var auth struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   string `json:"expires_in"`
	}

	_, err := resty.New().R().SetFormData(map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"client_id":     p.config.ClientID,
		"client_secret": p.config.ClientSecret,
		"redirect_uri":  p.callback,
	}).SetResult(auth).Post("https://www.linkedin.com/oauth/v2/accessToken")
	if err != nil {
		return nil, errctx.Wrap(err, ctx)
	}

	// TODO:
	// - get linkedin account ID
	// - p.auth_repo.GetByIdentifier(linkedin account ID)
	// - create if doesn't exist
	// - return account

	return nil, nil
}
