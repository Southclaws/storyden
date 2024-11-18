package google

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	ga "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/services/account/register"
	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth/all"
	"github.com/Southclaws/storyden/app/services/avatar"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/endec"
)

var (
	ErrAccessToken  = fault.New("failed to get access token")
	ErrMissingToken = fault.New("no access token in response")
)

var (
	service   = authentication.ServiceOAuthGoogle
	tokenType = authentication.TokenTypeOAuth
)

type Provider struct {
	auth_repo  authentication.Repository
	register   *register.Registrar
	avatar_svc avatar.Service

	ed       endec.EncrypterDecrypter
	callback string
	config   *all.Configuration
}

func New(
	cfg config.Config,
	auth_repo authentication.Repository,
	register *register.Registrar,
	avatar_svc avatar.Service,
	ed endec.EncrypterDecrypter,
) (*Provider, error) {
	config, err := all.LoadProvider(service)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	callback := all.Redirect(cfg, service)

	return &Provider{
		auth_repo:  auth_repo,
		register:   register,
		avatar_svc: avatar_svc,

		ed:       ed,
		config:   config,
		callback: callback,
	}, nil
}

func (p *Provider) Service() authentication.Service { return service }
func (p *Provider) Token() authentication.TokenType { return tokenType }

func (p *Provider) Enabled(ctx context.Context) (bool, error) {
	return p.config != nil, nil
}

func (p *Provider) Link(redirectPath string) (string, error) {
	state, err := p.ed.Encrypt(map[string]any{
		"redirect": redirectPath,
	}, time.Minute*10)
	if err != nil {
		return "", fault.Wrap(err)
	}

	oac := oauth2.Config{
		ClientID:     p.config.ClientID,
		ClientSecret: p.config.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  p.callback,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}

	return oac.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
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

	oac := oauth2.Config{
		ClientID:     p.config.ClientID,
		ClientSecret: p.config.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  p.callback,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}

	token, err := oac.Exchange(ctx, code, oauth2.AccessTypeOffline)
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
		service,
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

func (p *Provider) getOrCreateAccount(ctx context.Context, service authentication.Service, identifier, token, handle, name string) (*account.Account, error) {
	authmethod, exists, err := p.auth_repo.LookupByIdentifier(ctx, service, identifier)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to lookup existing account"), fctx.With(ctx))
	}

	if exists {
		return &authmethod.Account, nil
	}

	acc, err := p.register.Create(ctx, opt.New(handle),
		account_writer.WithName(name))
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to create new account"), fctx.With(ctx))
	}

	_, err = p.auth_repo.Create(ctx, acc.ID, service, authentication.TokenTypeOAuth, identifier, token, nil)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to create new auth method for account"), fctx.With(ctx))
	}

	return acc, nil
}
