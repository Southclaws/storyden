package bindings

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	waprovider "github.com/Southclaws/storyden/app/services/authentication/provider/webauthn"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
)

const cookieName = "storyden-webauthn-session"

var errNoCookie = fault.New("no webauthn session cookie")

type WebAuthn struct {
	cj           *session_cookie.Jar
	si           *session.Issuer
	accountQuery *account_querier.Querier
	wa           *waprovider.Provider
	address      url.URL
}

func NewWebAuthn(
	cfg config.Config,
	si *session.Issuer,
	accountQuery *account_querier.Querier,
	cj *session_cookie.Jar,
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

	return WebAuthn{cj, si, accountQuery, wa, cfg.PublicAPIAddress}
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
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
		Domain:   a.address.Hostname(),
		Secure:   true,
		HttpOnly: true,
	}

	return openapi.WebAuthnRequestCredential200JSONResponse{
		WebAuthnRequestCredentialOKJSONResponse: openapi.WebAuthnRequestCredentialOKJSONResponse{
			Headers: openapi.WebAuthnRequestCredentialOKResponseHeaders{
				SetCookie: cookie.String(),
			},
			Body: serialiseWebAuthnCredentialCreationOptions(*cred),
		},
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
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument), fmsg.With(pe.DevInfo))
	}

	invitedBy, err := deserialiseInvitationID(request.Params.InvitationId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	_, accountID, err := a.wa.FinishRegistration(ctx, string(session.UserID), *session, cr, invitedBy)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	t, err := a.si.Issue(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.WebAuthnMakeCredential200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccess{
				Id: xid.NilID().String(),
			},
			Headers: openapi.AuthSuccessOKResponseHeaders{
				SetCookie: a.cj.Create(*t).String(),
			},
		},
	}, nil
}

func (a *WebAuthn) WebAuthnGetAssertion(ctx context.Context, request openapi.WebAuthnGetAssertionRequestObject) (openapi.WebAuthnGetAssertionResponseObject, error) {
	cred, sessionData, err := a.wa.BeginLogin(ctx, string(request.AccountHandle))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

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
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
		Domain:   a.address.Hostname(),
		Secure:   true,
		HttpOnly: true,
	}

	return openapi.WebAuthnGetAssertion200JSONResponse{
		WebAuthnGetAssertionOKJSONResponse: openapi.WebAuthnGetAssertionOKJSONResponse{
			Body: serialiseWebAuthnCredentialRequestOptions(cred.Response),
			Headers: openapi.WebAuthnGetAssertionOKResponseHeaders{
				SetCookie: cookie.String(),
			},
		},
	}, nil
}

func (a *WebAuthn) WebAuthnMakeAssertion(ctx context.Context, request openapi.WebAuthnMakeAssertionRequestObject) (openapi.WebAuthnMakeAssertionResponseObject, error) {
	c := ctx.Value("webauthn")
	session, ok := c.(*webauthn.SessionData)
	if !ok {
		return nil, fault.Wrap(errNoCookie,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
		)
	}

	// something here is messing up userHandle
	b, err := json.Marshal(request.Body)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	reader := bytes.NewReader(b)

	cr, err := protocol.ParseCredentialRequestResponseBody(reader)
	if err != nil {
		pe := err.(*protocol.Error)
		ctx = fctx.WithMeta(ctx,
			"type", pe.Type,
			"details", pe.Details,
			"info", pe.DevInfo,
		)
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument), fmsg.With(pe.DevInfo))
	}

	_, acc, err := a.wa.FinishLogin(ctx, string(session.UserID), *session, cr)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	t, err := a.si.Issue(ctx, acc.ID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.WebAuthnMakeAssertion200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccess{
				Id: xid.NilID().String(),
			},
			Headers: openapi.AuthSuccessOKResponseHeaders{
				SetCookie: a.cj.Create(*t).String(),
			},
		},
	}, nil
}

func serialiseWebAuthnCredentialCreationOptions(cred protocol.CredentialCreation) openapi.WebAuthnPublicKeyCreationOptions {
	rp := openapi.PublicKeyCredentialRpEntity{
		Id:   cred.Response.RelyingParty.ID,
		Name: cred.Response.RelyingParty.Name,
	}

	user := openapi.PublicKeyCredentialUserEntity{
		DisplayName: cred.Response.User.DisplayName,
		Id:          fmt.Sprint(cred.Response.User.ID),
		Name:        cred.Response.User.Name,
	}

	pubKeyCredParams := dt.Map(cred.Response.Parameters, func(p protocol.CredentialParameter) openapi.PublicKeyCredentialParameters {
		alg := float32(p.Algorithm)
		return openapi.PublicKeyCredentialParameters{
			Type: openapi.PublicKeyCredentialType(p.Type),
			Alg:  alg,
		}
	})

	excludeCredentials := dt.Map(cred.Response.CredentialExcludeList, func(d protocol.CredentialDescriptor) openapi.PublicKeyCredentialDescriptor {
		transports := dt.Map(d.Transport, func(t protocol.AuthenticatorTransport) openapi.PublicKeyCredentialDescriptorTransports {
			return openapi.PublicKeyCredentialDescriptorTransports(t)
		})
		return openapi.PublicKeyCredentialDescriptor{
			Type:       openapi.PublicKeyCredentialType(d.Type),
			Id:         string(d.CredentialID),
			Transports: &transports,
		}
	})

	authenticatorSelection := &openapi.AuthenticatorSelectionCriteria{
		AuthenticatorAttachment: openapi.AuthenticatorAttachment(cred.Response.AuthenticatorSelection.AuthenticatorAttachment),
		RequireResidentKey:      cred.Response.AuthenticatorSelection.RequireResidentKey,
		ResidentKey:             openapi.ResidentKeyRequirement(cred.Response.AuthenticatorSelection.ResidentKey),
		UserVerification:        (*openapi.UserVerificationRequirement)(&cred.Response.AuthenticatorSelection.UserVerification),
	}

	return openapi.WebAuthnPublicKeyCreationOptions{
		PublicKey: openapi.PublicKeyCredentialCreationOptions{
			Rp:   rp,
			User: user,

			Challenge:        cred.Response.Challenge.String(),
			PubKeyCredParams: pubKeyCredParams,

			Timeout:                &cred.Response.Timeout,
			ExcludeCredentials:     excludeCredentials,
			AuthenticatorSelection: authenticatorSelection,
			Attestation:            (*openapi.AttestationConveyancePreference)(&cred.Response.Attestation),
			Extensions:             (*openapi.AuthenticationExtensionsClientInputs)(&cred.Response.Extensions),
		},
	}
}

func serialiseWebAuthnCredentialRequestOptions(cred protocol.PublicKeyCredentialRequestOptions) openapi.CredentialRequestOptions {
	allowedCredentials := dt.Map(cred.AllowedCredentials, func(cd protocol.CredentialDescriptor) openapi.PublicKeyCredentialDescriptor {
		transports := dt.Map(cd.Transport, func(t protocol.AuthenticatorTransport) openapi.PublicKeyCredentialDescriptorTransports {
			return openapi.PublicKeyCredentialDescriptorTransports(t)
		})
		id := make([]byte, base64.RawStdEncoding.EncodedLen(len(cd.CredentialID)))
		base64.RawURLEncoding.Encode(id, cd.CredentialID)
		return openapi.PublicKeyCredentialDescriptor{
			Id:         string(id),
			Transports: &transports,
			Type:       openapi.PublicKeyCredentialType(cd.Type),
		}
	})

	return openapi.CredentialRequestOptions{
		PublicKey: openapi.PublicKeyCredentialRequestOptions{
			AllowCredentials: &allowedCredentials,
			Challenge:        cred.Challenge.String(),
			RpId:             &cred.RelyingPartyID,
			Timeout:          &cred.Timeout,
			UserVerification: (*openapi.PublicKeyCredentialRequestOptionsUserVerification)(&cred.UserVerification),
		},
	}
}
