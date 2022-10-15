package linkedin

import (
	"context"
	"fmt"
	"time"

	"4d63.com/optional"
	"github.com/Southclaws/fault/errctx"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/authentication"
	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth/all"
	"github.com/Southclaws/storyden/app/services/avatar"
	"github.com/Southclaws/storyden/internal/config"
)

const (
	id   = "linkedin"
	name = "LinkedIn"
	logo = "https://brand.linkedin.com/content/dam/me/business/en-us/amp/brand-site/v2/bg/LI-Bug.svg.original.svg"
)

type LinkedInProvider struct {
	auth_repo    authentication.Repository
	account_repo account.Repository
	avatar_svc   avatar.Service

	callback string
	config   all.Configuration
}

func New(cfg config.Config, auth_repo authentication.Repository, account_repo account.Repository, avatar_svc avatar.Service) (*LinkedInProvider, error) {
	config, err := all.LoadProvider(id)
	if err != nil {
		return nil, err
	}

	return &LinkedInProvider{
		auth_repo:    auth_repo,
		account_repo: account_repo,
		avatar_svc:   avatar_svc,
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

	// Use the auth token for all future requests
	rest.SetAuthToken(auth.AccessToken)

	var profile struct {
		ID             string `json:"id"`
		FirstName      string `json:"localizedFirstName"`
		LastName       string `json:"localizedLastName"`
		ProfilePicture struct {
			DisplayImage struct {
				Elements []struct {
					Identifiers []struct {
						Identifier string `json:"identifier"`
					} `json:"identifiers"`
				} `json:"elements"`
			} `json:"displayImage~"`
		} `json:"profilePicture"`
	}

	_, err = rest.R().
		SetResult(&profile).
		Get("https://api.linkedin.com/v2/me?projection=(id,localizedFirstName,localizedLastName,profilePicture(displayImage~digitalmediaAsset:playableStreams))")
	if err != nil {
		return nil, errctx.Wrap(err, ctx)
	}

	// TODO: Iterate "elements" and search for the largest avatar instead of
	// just picking the first item from the list.
	avatarURL := optional.OfPtr(func() *string {
		if len(profile.ProfilePicture.DisplayImage.Elements) > 0 && len(profile.ProfilePicture.DisplayImage.Elements[0].Identifiers) > 0 {
			return &profile.ProfilePicture.DisplayImage.Elements[0].Identifiers[0].Identifier
		}
		return nil
	}())

	// TODO: Invent a better handle generator
	handle := fmt.Sprintf("%s-%s-%d", profile.FirstName, profile.LastName, time.Now().Day())

	name := fmt.Sprint(profile.FirstName, " ", profile.LastName)

	// TODO: Everything below this can be made generic for all OAuth providers.

	acc, err := p.getOrCreateAccount(ctx, "linkedin", profile.ID, auth.AccessToken, handle, name)
	if err != nil {
		return nil, errctx.Wrap(err, ctx)
	}

	if err := p.setAvatar(ctx, rest, acc, avatarURL); err != nil {
		// failing to set the avatar is not a big issue.
		fmt.Println(err)
	}

	return acc, nil
}

func (p *LinkedInProvider) getOrCreateAccount(ctx context.Context, provider authentication.Service, identifier, token, handle, name string) (*account.Account, error) {
	authmethod, exists, err := p.auth_repo.LookupByIdentifier(ctx, provider, identifier)
	if err != nil {
		return nil, errctx.Wrap(errors.Wrap(err, "failed to lookup existing account"), ctx)
	}

	if exists {
		return &authmethod.Account, nil
	}

	acc, err := p.account_repo.Create(ctx, handle,
		account.WithName(name))
	if err != nil {
		return nil, errctx.Wrap(errors.Wrap(err, "failed to create new account"), ctx)
	}

	_, err = p.auth_repo.Create(ctx, acc.ID, provider, identifier, token, nil)
	if err != nil {
		return nil, errctx.Wrap(errors.Wrap(err, "failed to create new auth method for account"), ctx)
	}

	return acc, nil
}

func (p *LinkedInProvider) setAvatar(ctx context.Context, rest *resty.Client, acc *account.Account, avatarURL optional.Optional[string]) error {
	url, ok := avatarURL.Get()
	if !ok {
		return nil
	}

	avatarBinary, err := rest.R().
		SetDoNotParseResponse(true).
		Get(url)
	if err != nil {
		return errctx.Wrap(err, ctx)
	}

	if avatarBinary == nil {
		return nil
	}

	if err := p.avatar_svc.Set(ctx, acc.ID, avatarBinary.RawBody()); err != nil {
		return errctx.Wrap(err, ctx)
	}

	return nil
}
