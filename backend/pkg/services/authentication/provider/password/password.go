package password

import (
	"context"
	"net/mail"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/pkg/errors"

	"github.com/Southclaws/storyden/backend/pkg/resources/authentication"
	"github.com/Southclaws/storyden/backend/pkg/resources/user"
)

const AuthServiceName = `password`

var ErrPasswordMismatch = errors.New("password mismatch")

type Password struct {
	auth authentication.Repository
	user user.Repository
}

func ErrExists(id string) error {
	return errors.Errorf("an account with the email '%s' already exists", id)
}

func NewBasicAuth(auth authentication.Repository, user user.Repository) *Password {
	return &Password{auth, user}
}

func (b *Password) Register(ctx context.Context, identifier string, password string) (*user.User, error) {
	addr, err := mail.ParseAddress(identifier)
	if err != nil {
		return nil, err
	}

	username := strings.Split(addr.Address, "@")[0]

	u, err := b.user.GetUserByEmail(ctx, identifier, false)
	if err != nil {
		return nil, err
	}
	if u != nil {
		return nil, ErrExists(identifier)
	}

	u, err = b.user.CreateUser(ctx, identifier, username)
	if err != nil {
		return nil, err
	}

	hashed, err := argon2id.CreateHash(identifier, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	_, err = b.auth.Create(ctx, u.ID, AuthServiceName, identifier, string(hashed), nil)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (b *Password) Login(ctx context.Context, identifier string, password string) (*user.User, error) {
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

	return &a.User, nil
}
