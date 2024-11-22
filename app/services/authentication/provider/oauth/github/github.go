package github

import (
	"context"
	"fmt"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/google/go-github/v66/github"
	"golang.org/x/oauth2"
	oauth2_github "golang.org/x/oauth2/github"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/services/account/register"
	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth/all"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/endec"
)

var (
	service   = authentication.ServiceOAuthGitHub
	tokenType = authentication.TokenTypeOAuth
)

type Provider struct {
	register *register.Registrar
	ed       endec.EncrypterDecrypter

	callback string
	config   *all.Configuration
}

func New(
	cfg config.Config,
	register *register.Registrar,
	ed endec.EncrypterDecrypter,
) (*Provider, error) {
	config, err := all.LoadProvider(service)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	return &Provider{
		register: register,
		ed:       ed,
		config:   config,
		callback: all.Redirect(cfg, service),
	}, nil
}

func (p *Provider) Service() authentication.Service { return service }
func (p *Provider) Token() authentication.TokenType { return tokenType }

func (p *Provider) Enabled(ctx context.Context) (bool, error) {
	return p.config != nil, nil
}

func (p *Provider) oauthConfig() *oauth2.Config {
	if p.config == nil {
		return nil
	}

	return &oauth2.Config{
		ClientID:     p.config.ClientID,
		ClientSecret: p.config.ClientSecret,
		Endpoint:     oauth2_github.Endpoint,
		RedirectURL:  p.callback,
		Scopes: []string{
			"read:user",
			"user:email",
		},
	}
}

func (p *Provider) Link(redirectPath string) (string, error) {
	state, err := p.ed.Encrypt(map[string]any{
		"redirect": redirectPath,
	}, time.Minute*10)
	if err != nil {
		return "", fault.Wrap(err)
	}

	oac := p.oauthConfig()

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

	oac := p.oauthConfig()

	token, err := oac.Exchange(ctx, code, oauth2.AccessTypeOffline)
	if err != nil {
		return nil, fault.Wrap(err,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("failed to exchange code for token", "This login token may have expired, please try again."),
		)
	}

	client := github.NewClient(nil).WithAuthToken(token.AccessToken)

	clientID := token.Extra("client_id").(string)

	u, _, err := client.Users.Get(ctx, clientID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if u.Login == nil {
		return nil, fault.New("missing login",
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("no username", "The GitHub API did not return a username for this account."),
		)
	}

	handle := strings.ToLower(*u.Login)
	name := opt.NewPtr(u.Name).Or(handle)

	authName := fmt.Sprintf("GitHub (@%s)", handle)

	if u.Email == nil {
		return p.register.GetOrCreateViaHandle(ctx,
			service,
			authName,
			strconv.Itoa(int(*u.ID)),
			token.AccessToken,
			handle,
			name,
		)
	}

	email, err := mail.ParseAddress(*u.Email)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return p.register.GetOrCreateViaEmail(ctx,
		service,
		authName,
		strconv.Itoa(int(*u.ID)),
		token.AccessToken,
		handle,
		name,
		*email,
	)
}
