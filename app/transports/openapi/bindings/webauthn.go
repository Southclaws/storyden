package bindings

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	waprovider "github.com/Southclaws/storyden/app/services/authentication/provider/webauthn"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/openapi"
)

const cookieName = "storyden-webauthn-session"

var errNoCookie = fault.New("no webauthn session cookie")

type WebAuthn struct {
	sm     Session
	ar     account.Repository
	wa     *waprovider.Provider
	domain string
}

func NewWebAuthn(
	cfg config.Config,
	ar account.Repository,
	sm Session,
	wa *waprovider.Provider,
	router *echo.Echo,
) WebAuthn {
	// in order to retain context across the credential request and creation,
	// a session cookie is used which stores the webauthn session information.
	router.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if s, err := c.Cookie(cookieName); err == nil {

				r := base64.NewDecoder(base64.URLEncoding, strings.NewReader(s.Value))

				session := &webauthn.SessionData{}
				if err := json.NewDecoder(r).Decode(&session); err == nil {
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

	// Encode the session data as a base64 JSON string

	j, err := json.Marshal(sessionData)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	value := base64.URLEncoding.EncodeToString(j)

	// save the base64 as a cookie for the WebAuthnMakeCredential call

	cookie := http.Cookie{
		Name:  cookieName,
		Value: value,
		// Expire this exchange after 10 minutes
		Expires:  time.Now().Add(time.Minute * 10),
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		Domain:   a.domain,
		Secure:   true,
		HttpOnly: true,
	}

	return openapi.WebAuthnPublicKeyCreationOptionsJSONResponse{
		Headers: openapi.WebAuthnPublicKeyCreationOptionsResponseHeaders{
			SetCookie: cookie.String(),
		},
		Body: cred,
	}, nil
}

func (a *WebAuthn) WebAuthnMakeCredential(ctx context.Context, request openapi.WebAuthnMakeCredentialRequestObject) (openapi.WebAuthnMakeCredentialResponseObject, error) {
	c := ctx.Value("webauthn")
	session, ok := c.(*webauthn.SessionData)
	if !ok {
		return nil, fault.Wrap(errNoCookie,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
		)
	}

	// NOTE: This is a hack due to oapi-codegen not giving us raw JSON.

	b, err := json.Marshal(request.Body)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	reader := bytes.NewReader(b)

	cr, err := protocol.ParseCredentialCreationResponseBody(reader)
	if err != nil {
		pe := err.(*protocol.Error)
		ctx = fctx.WithMeta(ctx,
			"type", pe.Type,
			"details", pe.Details,
			"info", pe.DevInfo,
		)
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With(pe.DevInfo))
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
		Body: openapi.AuthSuccess{
			Id: xid.NilID().String(),
		},
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
