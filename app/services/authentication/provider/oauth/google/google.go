package google

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	ga "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth/all"
	"github.com/Southclaws/storyden/app/services/authentication/register"
	"github.com/Southclaws/storyden/app/services/avatar"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/endec"
)

var (
	ErrAccessToken  = fault.New("failed to get access token")
	ErrMissingToken = fault.New("no access token in response")
)

const (
	id   = "google"
	name = "Google"
)

type Provider struct {
	auth_repo  authentication.Repository
	register   register.Service
	avatar_svc avatar.Service

	ed       endec.EncrypterDecrypter
	callback string
	config   all.Configuration
	oac      oauth2.Config
}

func New(
	cfg config.Config,
	auth_repo authentication.Repository,
	register register.Service,
	avatar_svc avatar.Service,
	ed endec.EncrypterDecrypter,
) (*Provider, error) {
	config, err := all.LoadProvider(id)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	callback := all.Redirect(cfg, id)

	oac := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  callback,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}

	return &Provider{
		auth_repo:  auth_repo,
		register:   register,
		avatar_svc: avatar_svc,

		ed:       ed,
		config:   config,
		callback: callback,
		oac:      oac,
	}, nil
}

func (p *Provider) Enabled() bool { return p.config.Enabled }
func (p *Provider) Name() string  { return name }
func (p *Provider) ID() string    { return id }

func (p *Provider) Link(redirectPath string) (string, error) {
	state, err := p.ed.Encrypt(map[string]any{
		"redirect": redirectPath,
	}, time.Minute*10)
	if err != nil {
		return "", fault.Wrap(err)
	}

	return p.oac.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (p *Provider) Login(ctx context.Context, state, code string) (*account.Account, error) {
	_, err := p.ed.Decrypt(state)
	if err != nil {
		return nil, fault.Wrap(err,
			fctx.With(ctx),
			fmsg.WithDesc("failed to decrypt state value", "This link has expired, please try again."),
		)
	}
	// TODO: Process claims for redirect etc.

	token, err := p.oac.Exchange(ctx, code, oauth2.AccessTypeOffline)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	gs, err := ga.NewService(ctx, option.WithTokenSource(oauth2.StaticTokenSource(token)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	u, err := ga.NewUserinfoV2MeService(gs).Get().Do()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// TODO: Invent a better handle generator
	handle := fmt.Sprintf("%s-%s-%d", u.GivenName, u.FamilyName, time.Now().Day())

	name := fmt.Sprint(u.GivenName, " ", u.FamilyName)

	// TODO: Everything below this can be made generic for all OAuth providers.

	acc, err := p.getOrCreateAccount(ctx,
		id,
		strings.ToLower(u.Id),
		token.AccessToken,
		handle,
		name,
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, nil
}

func (p *Provider) getOrCreateAccount(ctx context.Context, provider authentication.Service, identifier, token, handle, name string) (*account.Account, error) {
	authmethod, exists, err := p.auth_repo.LookupByIdentifier(ctx, provider, identifier)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to lookup existing account"), fctx.With(ctx))
	}

	if exists {
		return &authmethod.Account, nil
	}

	acc, err := p.register.Create(ctx, handle,
		account.WithName(name))
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to create new account"), fctx.With(ctx))
	}

	_, err = p.auth_repo.Create(ctx, acc.ID, provider, identifier, token, nil)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to create new auth method for account"), fctx.With(ctx))
	}

	return acc, nil
}
