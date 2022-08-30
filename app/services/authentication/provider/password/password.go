package password

import (
	"context"
	"net/mail"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/pkg/errors"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/authentication"
	"github.com/Southclaws/storyden/internal/errtag"
)

const AuthServiceName = `password`

var ErrPasswordMismatch = errors.New("password mismatch")

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

	username := strings.Split(addr.Address, "@")[0]

	u, exists, err := b.account.LookupByEmail(ctx, identifier)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}

	if exists {
		return nil, errtag.Wrap(errors.New("exists"), errtag.AlreadyExists{})
	}

	u, err = b.account.Create(ctx, identifier, username)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create account")
	}

	hashed, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create secure password hash")
	}

	_, err = b.auth.Create(ctx, u.ID, AuthServiceName, identifier, string(hashed), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create account authentication instance")
	}

	return u, nil
}

func (b *Password) Login(ctx context.Context, identifier string, password string) (*account.Account, error) {
	a, err := b.auth.GetByIdentifier(ctx, AuthServiceName, identifier)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account")
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
