package password

import (
	"context"
	"net/mail"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/pkg/errors"

	"github.com/Southclaws/storyden/backend/pkg/resources/account"
	"github.com/Southclaws/storyden/backend/pkg/resources/authentication"
)

const AuthServiceName = `password`

var ErrPasswordMismatch = errors.New("password mismatch")

type Password struct {
	auth    authentication.Repository
	account account.Repository
}

var ErrExists = errors.New("already exists")

func errExists(id string) error {
	return errors.Wrapf(ErrExists, "with email '%s'", id)
}

func NewBasicAuth(auth authentication.Repository, account account.Repository) *Password {
	return &Password{auth, account}
}

func (b *Password) Register(ctx context.Context, identifier string, password string) (*account.Account, error) {
	addr, err := mail.ParseAddress(identifier)
	if err != nil {
		return nil, err
	}

	username := strings.Split(addr.Address, "@")[0]

	u, exists, err := b.account.LookupByEmail(ctx, identifier)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errExists(identifier)
	}

	u, err = b.account.Create(ctx, identifier, username)
	if err != nil {
		return nil, err
	}

	hashed, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	_, err = b.auth.Create(ctx, u.ID, AuthServiceName, identifier, string(hashed), nil)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (b *Password) Login(ctx context.Context, identifier string, password string) (*account.Account, error) {
	a, err := b.auth.GetByIdentifier(ctx, AuthServiceName, identifier)
	if err != nil {
		return nil, err
	}

	match, _, err := argon2id.CheckHash(password, a.Token)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, ErrPasswordMismatch
	}

	return &a.Account, nil
}
