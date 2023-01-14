package webauthn

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
)

// TODO: Move to actual accounts model.
type temporary struct {
	handle      string
	credentials []webauthn.Credential
}

func (t *temporary) WebAuthnID() []byte                         { return []byte(t.handle) }
func (t *temporary) WebAuthnName() string                       { return t.handle }
func (t *temporary) WebAuthnDisplayName() string                { return t.handle }
func (t *temporary) WebAuthnIcon() string                       { return "" }
func (t *temporary) WebAuthnCredentials() []webauthn.Credential { return t.credentials }

func (p *Provider) BeginRegistration(ctx context.Context, handle string) (*protocol.CredentialCreation, *webauthn.SessionData, error) {
	t := temporary{handle: handle}

	// TODO: Check if handle already exists
	// if it exists, maybe we can short circuit the flow and switch to login?

	credentialOptions, sessionData, err := p.wa.BeginRegistration(&t,
		webauthn.WithAuthenticatorSelection(
			protocol.AuthenticatorSelection{
				// AuthenticatorAttachment: protocol.AuthenticatorAttachment(authType),
				// RequireResidentKey:      residentKeyRequirement,
				// UserVerification:        protocol.UserVerificationRequirement(userVer),
			}),
		// webauthn.WithConveyancePreference(protocol.ConveyancePreference(attType)),
	)
	if err != nil {
		return nil, nil, fault.Wrap(err, fmsg.With("failed to start registration"))
	}

	return credentialOptions, sessionData, nil
}

func (p *Provider) FinishRegistration(ctx context.Context,
	handle string,
	session webauthn.SessionData,
	parsedResponse *protocol.ParsedCredentialCreationData,
) (*webauthn.Credential, account.AccountID, error) {
	t := temporary{handle: handle}

	credential, err := p.wa.CreateCredential(&t, session, parsedResponse)
	if err != nil {
		ctx = fctx.WithMeta(ctx, waErrMetadata(err)...)
		return nil, account.AccountID(xid.NilID()), fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := p.register(ctx, handle, credential)
	if err != nil {
		return nil, account.AccountID(xid.NilID()), fault.Wrap(err, fctx.With(ctx))
	}

	return credential, acc.ID, nil
}

func waErrMetadata(in error) []string {
	switch err := in.(type) {
	case *protocol.Error:
		return []string{
			"wa_details", err.Details,
			"wa_info", err.DevInfo,
			"wa_type", err.Type,
		}

	default:
		return []string{}
	}
}
