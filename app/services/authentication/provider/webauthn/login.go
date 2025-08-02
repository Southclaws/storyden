package webauthn

import (
	"context"
	"encoding/json"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
)

func (p *Provider) BeginLogin(ctx context.Context, handle string) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	acc, exists, err := p.accountQuery.LookupByHandle(ctx, handle)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !exists {
		return nil, nil, fault.Wrap(ErrNotFound,
			fctx.With(ctx),
			ftag.With(ftag.NotFound),
			fmsg.WithDesc("not found", "No account was found with the provided handle."))
	}

	if err := acc.RejectSuspended(); err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	ams, err := p.auth_repo.GetAuthMethods(ctx, acc.ID)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	ams = dt.Filter(ams, func(a *authentication.Authentication) bool { return a.Service == service })
	if len(ams) == 0 {
		return nil, nil, fault.Wrap(ErrNoAuthRecord,
			fctx.With(ctx),
			fmsg.WithDesc("no auth method", "This account does not have a Passkey (WebAuthn) credential."),
		)
	}

	credentials, err := dt.MapErr(ams, func(a *authentication.Authentication) (webauthn.Credential, error) {
		var wac webauthn.Credential
		if err := json.Unmarshal([]byte(a.Token), &wac); err != nil {
			return wac, fault.Wrap(err, fmsg.With("malformed credential from auth storage"))
		}
		return wac, nil
	})
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	t := temporary{handle, credentials}

	credential, sd, err := p.wa.BeginLogin(&t)
	if err != nil {
		ctx = fctx.WithMeta(ctx, waErrMetadata(err)...)
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	return credential, sd, nil
}

func (p *Provider) FinishLogin(ctx context.Context,
	handle string,
	session webauthn.SessionData,
	parsedResponse *protocol.ParsedCredentialAssertionData,
) (*webauthn.Credential, *account.Account, error) {
	acc, exists, err := p.accountQuery.LookupByHandle(ctx, handle)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !exists {
		return nil, nil, fault.Wrap(ErrNotFound,
			fctx.With(ctx),
			ftag.With(ftag.NotFound),
			fmsg.WithDesc("not found", "No account was found with the provided handle."))
	}

	if err := acc.RejectSuspended(); err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	ams, err := p.auth_repo.GetAuthMethods(ctx, acc.ID)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	ams = dt.Filter(ams, func(a *authentication.Authentication) bool { return a.Service == service })
	if len(ams) == 0 {
		return nil, nil, fault.Wrap(ErrNoAuthRecord,
			fctx.With(ctx),
			fmsg.WithDesc("no auth method", "This account does not have a Passkey (WebAuthn) credential."),
		)
	}

	credentials, err := dt.MapErr(ams, func(a *authentication.Authentication) (webauthn.Credential, error) {
		var wac webauthn.Credential
		if err := json.Unmarshal([]byte(a.Token), &wac); err != nil {
			return wac, fault.Wrap(err, fmsg.With("malformed credential from auth storage"))
		}
		return wac, nil
	})
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	t := temporary{handle: handle, credentials: credentials}

	cred, err := p.wa.ValidateLogin(&t, session, parsedResponse)
	if err != nil {
		ctx = fctx.WithMeta(ctx, waErrMetadata(err)...)
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	return cred, &acc.Account, nil
}
