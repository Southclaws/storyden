package bindings

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"

	"github.com/Southclaws/storyden/app/resources/account"
	waprovider "github.com/Southclaws/storyden/app/services/authentication/provider/webauthn"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/openapi"
)

type WebAuthn struct {
	sm     *Session
	ar     account.Repository
	wa     *waprovider.Provider
	domain string
}

func NewWebAuthn(
	cfg config.Config,
	ar account.Repository,
	sm *Session,
	wa *waprovider.Provider,
	router *echo.Echo,
) WebAuthn {
	// in order to retain context across the credential request and creation,
	// a session cookie is used which stores the webauthn session information.
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

	return WebAuthn{sm, ar, wa, cfg.CookieDomain}
}

func (a *WebAuthn) WebAuthnRequestCredential(ctx context.Context, request openapi.WebAuthnRequestCredentialRequestObject) (openapi.WebAuthnRequestCredentialResponseObject, error) {
	cred, sessionData, err := a.wa.BeginRegistration(ctx, string(request.AccountHandle))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	cookie, err := json.Marshal(sessionData)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to encode session data"))
	}

	return openapi.WebAuthnPublicKeyCreationOptionsJSONResponse{
		Headers: openapi.WebAuthnPublicKeyCreationOptionsResponseHeaders{
			SetCookie: string(cookie),
		},
		Body: cred,
	}, nil
}

func (a *WebAuthn) WebAuthnMakeCredential(ctx context.Context, request openapi.WebAuthnMakeCredentialRequestObject) (openapi.WebAuthnMakeCredentialResponseObject, error) {
	session, ok := ctx.Value("webauthn").(*webauthn.SessionData)
	if !ok {
		return nil, nil
	}

	// NOTE: This is a hack due to oapi-codegen not giving us raw JSON.

	b, err := json.Marshal(request.Body)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	reader := bytes.NewReader(b)

	cr, err := protocol.ParseCredentialCreationResponseBody(reader)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	_, accountID, err := a.wa.FinishRegistration(ctx, string(session.UserID), *session, cr)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	cookie, err := a.sm.encodeSession(accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.WebAuthnMakeCredential200JSONResponse{
		Body: openapi.AuthSuccess{},
		Headers: openapi.AuthSuccessResponseHeaders{
			SetCookie: string(cookie),
		},
	}, nil
}

func (a *WebAuthn) WebAuthnGetAssertion(ctx context.Context, request openapi.WebAuthnGetAssertionRequestObject) (openapi.WebAuthnGetAssertionResponseObject, error) {
	return nil, nil
}

func (a *WebAuthn) WebAuthnMakeAssertion(ctx context.Context, request openapi.WebAuthnMakeAssertionRequestObject) (openapi.WebAuthnMakeAssertionResponseObject, error) {
	return nil, nil
}
