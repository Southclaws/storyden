package keycloak

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/services/account/register"
	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/endec"
)

var (
	service   = authentication.ServiceOAuthKeycloak
	tokenType = authentication.TokenTypeOAuth
)

// Provider implements OAuthProvider using Keycloak OIDC discovery.
type Provider struct {
	config      oauth.Configuration
	register    *register.Registrar
	ed          endec.EncrypterDecrypter
	issuer      *oidc.Provider
	verifier    *oidc.IDTokenVerifier
	callbackURL string
}

// New constructs a Keycloak OAuth provider, running OIDC discovery.
func New(
	cfg config.Config,
	register *register.Registrar,
	ed endec.EncrypterDecrypter,
) (*Provider, error) {
	if !cfg.KeycloakEnabled {
		return &Provider{
			config: oauth.Configuration{
				Enabled:      cfg.KeycloakEnabled,
				ClientID:     cfg.KeycloakClientID,
				ClientSecret: cfg.KeycloakClientSecret,
			},
		}, nil
	}

	if ed == nil {
		return nil, fault.New("JWT provider must be enabled by setting JWT_SECRET for Keycloak OAuth provider")
	}
	ctx := context.Background()

	issuer, err := oidc.NewProvider(ctx, cfg.KeycloakIssuerURL.String())
	if err != nil {
		return nil, fault.Wrap(err)
	}
	verifier := issuer.Verifier(&oidc.Config{ClientID: cfg.KeycloakClientID})

	// Build the OAuth callback URL using the same pattern as other providers
	callbackURL := oauth.Redirect(cfg.PublicWebAddress, service)

	return &Provider{
		config: oauth.Configuration{
			Enabled:      cfg.KeycloakEnabled,
			ClientID:     cfg.KeycloakClientID,
			ClientSecret: cfg.KeycloakClientSecret,
		},
		register:    register,
		ed:          ed,
		issuer:      issuer,
		verifier:    verifier,
		callbackURL: callbackURL.String(), // Call .String() here instead
	}, nil
}

// Service returns the authentication.Service this provider implements.
func (p *Provider) Service() authentication.Service { return service }

// Token returns the token type used by this provider (always OAuth).
func (p *Provider) Token() authentication.TokenType { return tokenType }

// Enabled reports whether the provider is enabled via config.
func (p *Provider) Enabled(ctx context.Context) (bool, error) {
	return p.config.Enabled, nil
}

// oauthConfig builds the OAuth2 configuration using OIDC discovery.
func (p *Provider) oauthConfig(redirect string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     p.config.ClientID,
		ClientSecret: p.config.ClientSecret,
		Endpoint:     p.issuer.Endpoint(),
		RedirectURL:  redirect,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
}

// Link returns the URL to redirect a user to Keycloak for authentication.
// The redirectPath parameter is the OAuth callback URL (not the post-login destination).
func (p *Provider) Link(redirectPath string) (string, error) {
	state, err := p.ed.Encrypt(map[string]any{"redirect": redirectPath}, time.Minute*10)
	if err != nil {
		return "", fault.Wrap(err)
	}
	oac := p.oauthConfig(redirectPath)
	return oac.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

// Login completes the OAuth2 flow: exchanges code, verifies ID token, and returns the Account.
func (p *Provider) Login(ctx context.Context, state, code string) (*account.Account, error) {
	c, err := p.ed.Decrypt(state)
	if err != nil {
		return nil, fault.Wrap(err,
			fctx.With(ctx),
			fmsg.WithDesc("failed to decrypt state value", "This link has expired, please try again."),
		)
	}

	redirect := c["redirect"].(string)
	oac := p.oauthConfig(redirect)
	tok, err := oac.Exchange(ctx, code)
	if err != nil {
		return nil, fault.Wrap(err,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("failed to exchange code for token", "This login token may have expired, please try again."),
		)
	}

	rawID, ok := tok.Extra("id_token").(string)
	if !ok {
		return nil, fault.New("no id_token field in oauth2 token", fctx.With(ctx))
	}
	idToken, err := p.verifier.Verify(ctx, rawID)
	if err != nil {
		return nil, fault.Wrap(err,
			fctx.With(ctx),
			fmsg.WithDesc("failed to verify Keycloak ID token", "Authentication failed. The login token may be invalid or expired. Please try again."))
	}

	var claims struct {
		Email             string `json:"email"`
		EmailVerified     bool   `json:"email_verified"`
		PreferredUsername string `json:"preferred_username"`
		Name              string `json:"name"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return nil, fault.Wrap(err,
			fctx.With(ctx),
			fmsg.WithDesc("failed to parse Keycloak token claims", "Unable to read authentication information. Please try again."))
	}

	handle := strings.ToLower(claims.PreferredUsername)
	if handle == "" {
		parts := strings.Split(claims.Email, "@")
		handle = parts[0]
	}
	name := claims.Name
	emailAddr, err := mail.ParseAddress(claims.Email)
	if err != nil {
		return nil, fault.Wrap(err,
			fctx.With(ctx),
			fmsg.WithDesc("failed to parse Keycloak email address", "The email address from Keycloak is invalid. Please check your account settings."))
	}
	authName := fmt.Sprintf("Keycloak (%s)", emailAddr.Address)

	return p.register.GetOrCreateViaEmail(
		ctx,
		service,
		authName,
		idToken.Subject,
		tok.AccessToken,
		handle,
		name,
		*emailAddr,
		claims.EmailVerified,
	)
}

// Bootstrap generates OIDC bootstrap values for machine-consumable login flows.
// The redirectPath parameter is the post-login destination (e.g., "/dashboard"),
// NOT the OAuth redirect_uri. The OAuth redirect_uri is always the registered callback URL.
func (p *Provider) Bootstrap(redirectPath string) (map[string]string, error) {
	// Store the post-login destination in state, NOT the OAuth callback URL
	state, err := p.ed.Encrypt(map[string]any{
		"redirect":          p.callbackURL, // OAuth callback URL for token exchange
		"final_destination": redirectPath,  // Where to send user after login
	}, time.Minute*10)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	// Use the registered OAuth callback URL, not the user's destination
	oac := p.oauthConfig(p.callbackURL)
	authorizeURL := oac.AuthCodeURL(state, oauth2.AccessTypeOffline)

	expiresAt := time.Now().Add(time.Minute * 10).Format(time.RFC3339)

	return map[string]string{
		"authorize_url": authorizeURL,
		"expires_at":    expiresAt,
	}, nil
}
