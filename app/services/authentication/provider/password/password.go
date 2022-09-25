package password

import (
	"context"
	"net/mail"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/pkg/errors"

	"github.com/Southclaws/fault/errtag"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/authentication"
)

const AuthServiceName = `password`

var (
	ErrPasswordMismatch = errors.New("password mismatch")
	ErrNotFound         = errors.New("account not found")
)

type Password struct {
	auth    authentication.Repository
	account account.Repository
}

func NewBasicAuth(auth authentication.Repository, account account.Repository) *Password {
	return &Password{auth, account}
}

func (b *Password) Register(ctx context.Context, identifier string, password string) (*account.Account, error) {
	addr, err := mail.ParseAddress(identifier)
	if err != nil {
		return nil, errtag.Wrap(err, errtag.InvalidArgument{})
	}

	handle := strings.Split(addr.Address, "@")[0]

	_, exists, err := b.auth.LookupByIdentifier(ctx, AuthServiceName, identifier)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}

	if exists {
		return nil, errtag.Wrap(errors.New("account already exists"), errtag.AlreadyExists{})
	}

	account, err := b.account.Create(ctx, handle)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create account")
	}

	hashed, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create secure password hash")
	}

	_, err = b.auth.Create(ctx, account.ID, AuthServiceName, identifier, string(hashed), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create account authentication instance")
	}

	return account, nil
}

func (b *Password) Login(ctx context.Context, identifier string, password string) (*account.Account, error) {
	a, exists, err := b.auth.LookupByIdentifier(ctx, AuthServiceName, identifier)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}

	if !exists {
		return nil, errtag.Wrap(ErrNotFound, errtag.NotFound{})
	}

	match, _, err := argon2id.CheckHash(password, a.Token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compare secure password hash")
	}

	if !match {
		return nil, errtag.Wrap(ErrPasswordMismatch, errtag.Unauthenticated{})
	}

	return &a.Account, nil
}
