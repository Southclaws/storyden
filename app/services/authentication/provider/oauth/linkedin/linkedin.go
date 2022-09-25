package linkedin

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
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
	auth_repo    authentication.Repository
	account_repo account.Repository

	callback string
	config   all.Configuration
}

func New(cfg config.Config, auth_repo authentication.Repository, account_repo account.Repository) (*LinkedInProvider, error) {
	config, err := all.LoadProvider(id)
	if err != nil {
		return nil, err
	}

	return &LinkedInProvider{
		auth_repo:    auth_repo,
		account_repo: account_repo,
		config:       config,
		callback:     all.Redirect(cfg, id),
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
		ExpiresIn   int    `json:"expires_in"`
	}

	var authError struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}

	rest := resty.New()

	_, err := rest.R().
		SetFormData(map[string]string{
			"grant_type":    "authorization_code",
			"code":          code,
			"client_id":     p.config.ClientID,
			"client_secret": p.config.ClientSecret,
			"redirect_uri":  p.callback,
		}).
		SetResult(&auth).
		SetError(&authError).
		Post("https://www.linkedin.com/oauth/v2/accessToken")
	if err != nil {
		return nil, errctx.Wrap(err, ctx)
	}

	if authError.Error != "" {
		return nil, errctx.Wrap(errors.New(authError.ErrorDescription), ctx)
	}

	if auth.AccessToken == "" {
		return nil, errctx.Wrap(errors.New("no access token in response"), ctx)
	}

	var profile struct {
		ID        string `json:"id"`
		FirstName string `json:"localizedFirstName"`
		LastName  string `json:"localizedLastName"`
	}

	_, err = rest.R().
		SetAuthToken(auth.AccessToken).
		SetResult(&profile).
		Get("https://api.linkedin.com/v2/me")
	if err != nil {
		return nil, errctx.Wrap(err, ctx)
	}

	// TODO: Everything below this can be made generic for all OAuth providers.

	authmethod, exists, err := p.auth_repo.LookupByIdentifier(ctx, "linkedin", profile.ID)
	if err != nil {
		return nil, errctx.Wrap(err, ctx)
	}

	if exists {
		return &authmethod.Account, nil
	}

	acc, err := p.account_repo.Create(ctx, "temp@temp.com", "temp",
		account.WithName(fmt.Sprint(profile.FirstName, " ", profile.LastName)))
	if err != nil {
		return nil, errctx.Wrap(err, ctx)
	}

	_, err = p.auth_repo.Create(ctx, acc.ID, "linkedin", profile.ID, auth.AccessToken, nil)
	if err != nil {
		return nil, errctx.Wrap(errors.Wrap(err, "failed to create new auth method for account"), ctx)
	}

	return acc, nil
}
