package webauthn

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/services/authentication/session"
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
	accountID := session.GetOptAccountID(ctx)

	t := temporary{handle: handle}

	acc, exists, err := p.accountQuery.LookupByHandle(ctx, handle)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}
	if exists {
		id, ok := accountID.Get()
		if !ok {
			// TODO: maybe we can short circuit the flow and switch to login?
			return nil, nil, fault.Wrap(ErrAccountExists,
				fctx.With(ctx),
				ftag.With(ftag.AlreadyExists),
				fmsg.WithDesc(
					"already exists",
					"An account with this handle has already been registered.",
				),
			)
		}

		if id != acc.ID {
			return nil, nil, fault.Wrap(ErrExistsOnAnotherAccount,
				fctx.With(ctx),
				ftag.With(ftag.PermissionDenied),
				fmsg.WithDesc(
					"account id mismatch",
					"The specified handle does not match the authenticated account so the credential cannot be linked.",
				),
			)
		}
	}

	credentialOptions, sessionData, err := p.wa.BeginRegistration(&t,
		webauthn.WithAuthenticatorSelection(
			protocol.AuthenticatorSelection{
				AuthenticatorAttachment: protocol.AuthenticatorAttachment(protocol.Platform),
				RequireResidentKey:      protocol.ResidentKeyRequired(),
				ResidentKey:             protocol.ResidentKeyRequirementRequired,
				UserVerification:        protocol.VerificationPreferred,
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
	inviteCode opt.Optional[xid.ID],
) (*webauthn.Credential, account.AccountID, error) {
	t := temporary{handle: handle}

	credential, err := p.wa.CreateCredential(&t, session, parsedResponse)
	if err != nil {
		ctx = fctx.WithMeta(ctx, waErrMetadata(err)...)
		return nil, account.AccountID(xid.NilID()), fault.Wrap(err, fctx.With(ctx))
	}

	a, err := p.createOrUpdateAccount(ctx, handle, credential, inviteCode)
	if err != nil {
		return nil, account.AccountID(xid.NilID()), fault.Wrap(err, fctx.With(ctx))
	}

	return credential, a, err
}

func (p *Provider) createOrUpdateAccount(ctx context.Context,
	handle string,
	credential *webauthn.Credential,
	inviteCode opt.Optional[xid.ID],
) (account.AccountID, error) {
	accountID := session.GetOptAccountID(ctx)

	// If the request is from an account that's already logged in, instead of
	// treating this as a registration for a new account, add the credential to
	// the existing account's authentications keychain.
	if id, ok := accountID.Get(); ok {
		acc, err := p.add(ctx, id, credential)
		if err != nil {
			return account.AccountID(xid.NilID()), fault.Wrap(err, fctx.With(ctx))
		}

		return acc.ID, nil
	}

	opts := []account_writer.Option{}
	inviteCode.Call(func(id xid.ID) { opts = append(opts, account_writer.WithInvitedBy(id)) })

	acc, err := p.register(ctx, handle, credential, inviteCode)
	if err != nil {
		return account.AccountID(xid.NilID()), fault.Wrap(err, fctx.With(ctx))
	}

	return acc.ID, nil
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
