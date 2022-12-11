package bindings

import (
	"context"
	"encoding/json"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gorilla/securecookie"
	"github.com/kr/pretty"
	"github.com/labstack/echo/v4"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/openapi"
)

type WebAuthn struct {
	sc     *securecookie.SecureCookie
	ar     account.Repository
	wa     *webauthn.WebAuthn
	domain string
}

func NewWebAuthn(
	cfg config.Config,
	ar account.Repository,
	sc *securecookie.SecureCookie,
	wa *webauthn.WebAuthn,
	router *echo.Echo,
) WebAuthn {
	router.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if s, err := c.Cookie("storyden-webauthn-session"); err == nil {
				session := webauthn.SessionData{}
				if err := json.Unmarshal([]byte(s.Value), &session); err == nil {
					r := c.Request()
					ctx := r.Context()
					ctx = context.WithValue(ctx, "webauthn", session)
					c.SetRequest(r.WithContext(ctx))
				}
			}

			return next(c)
		}
	})

	return WebAuthn{sc, ar, wa, cfg.CookieDomain}
}

// TODO: Move to actual accounts model.
type temporary struct{ handle string }

func (t *temporary) WebAuthnID() []byte                         { return []byte(t.handle) }
func (t *temporary) WebAuthnName() string                       { return t.handle }
func (t *temporary) WebAuthnDisplayName() string                { return t.handle }
func (t *temporary) WebAuthnIcon() string                       { return "" }
func (t *temporary) WebAuthnCredentials() []webauthn.Credential { return nil }

func (a *WebAuthn) WebAuthnRequestCredential(ctx context.Context, request openapi.WebAuthnRequestCredentialRequestObject) (openapi.WebAuthnRequestCredentialResponseObject, error) {
	t := temporary{string(request.AccountHandle)}

	credentialOptions, sessionData, err := a.wa.BeginRegistration(&t,
		webauthn.WithAuthenticatorSelection(
			protocol.AuthenticatorSelection{
				// AuthenticatorAttachment: protocol.AuthenticatorAttachment(authType),
				// RequireResidentKey:      residentKeyRequirement,
				// UserVerification:        protocol.UserVerificationRequirement(userVer),
			}),
		// webauthn.WithConveyancePreference(protocol.ConveyancePreference(attType)),
	)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to start registration"))
	}

	cookie, err := json.Marshal(sessionData)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to encode session data"))
	}

	return openapi.WebAuthnPublicKeyCreationOptionsJSONResponse{
		Headers: openapi.WebAuthnPublicKeyCreationOptionsResponseHeaders{
			SetCookie: string(cookie),
		},
		Body: credentialOptions,
	}, nil
}

func (a *WebAuthn) WebAuthnMakeCredential(ctx context.Context, request openapi.WebAuthnMakeCredentialRequestObject) (openapi.WebAuthnMakeCredentialResponseObject, error) {
	session, ok := ctx.Value("webauthn").(*webauthn.SessionData)
	if !ok {
		return nil, nil
	}

	pretty.Println(request.Body)
	pretty.Println(session)

	// TODO: Lookup user
	// t :=temporary{string(session.UserID)}

	// protocol.ParseCredentialCreationResponseBody(request.Body)

	// a.wa.ParseCredentialCreationResponseBody(&t, *session, request)

	return nil, nil
}

func (a *WebAuthn) WebAuthnGetAssertion(ctx context.Context, request openapi.WebAuthnGetAssertionRequestObject) (openapi.WebAuthnGetAssertionResponseObject, error) {
	return nil, nil
}

func (a *WebAuthn) WebAuthnMakeAssertion(ctx context.Context, request openapi.WebAuthnMakeAssertionRequestObject) (openapi.WebAuthnMakeAssertionResponseObject, error) {
	return nil, nil
}
