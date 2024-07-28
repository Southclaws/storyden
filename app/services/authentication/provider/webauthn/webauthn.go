package webauthn

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/services/account/register"
	"github.com/Southclaws/storyden/app/transports/openapi/glue"
)

var (
	ErrNoAuthRecord           = fault.New("webauthn does not match account")
	ErrExistsOnAnotherAccount = fault.New("webauthn id already bound to another account")
	ErrNotFound               = fault.New("account not found")
	ErrAccountExists          = fault.New("requester already has an account")
)

const (
	id   = "webauthn"
	name = "WebAuthn"
)

type Provider struct {
	auth_repo    authentication.Repository
	accountQuery account_querier.Querier
	reg          *register.Registrar

	wa *webauthn.WebAuthn
}

func New(
	auth_repo authentication.Repository,
	accountQuery account_querier.Querier,
	reg *register.Registrar,

	wa *webauthn.WebAuthn,
) (*Provider, error) {
	return &Provider{
		auth_repo:    auth_repo,
		accountQuery: accountQuery,
		reg:          reg,
		wa:           wa,
	}, nil
}

func (p *Provider) Enabled() bool { return true }
func (p *Provider) ID() string    { return id }
func (p *Provider) Name() string  { return name }

func (b *Provider) Link(_ string) (string, error) {
	return "", nil
}

func (p *Provider) Login(ctx context.Context, handle, pubkey string) (*account.Account, error) {
	return nil, nil
}

func (p *Provider) register(ctx context.Context, handle string, credential *webauthn.Credential) (*account.Account, error) {
	acc, exists, err := p.accountQuery.LookupByHandle(ctx, handle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if exists {
		return nil, fault.Wrap(ErrAccountExists,
			fctx.With(ctx),
			ftag.With(ftag.AlreadyExists),
			fmsg.WithDesc(
				"already exists",
				"An account with this handle has already been registered without a Passkey (WebAuthn) credential.",
			),
		)
	}

	acc, err = p.reg.Create(ctx, handle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	encoded, err := json.Marshal(credential)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	_, err = p.auth_repo.Create(ctx,
		acc.ID,
		id,
		base64.RawURLEncoding.EncodeToString(credential.ID),
		string(encoded),
		nil,
		authentication.WithName(glue.GetDeviceName(ctx)),
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, nil
}

func (p *Provider) add(ctx context.Context, accountID account.AccountID, credential *webauthn.Credential) (*account.Account, error) {
	acc, err := p.accountQuery.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	encoded, err := json.Marshal(credential)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	_, err = p.auth_repo.Create(ctx,
		acc.ID,
		id,
		base64.RawURLEncoding.EncodeToString(credential.ID),
		string(encoded),
		nil,
		authentication.WithName(glue.GetDeviceName(ctx)),
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, nil
}
