package account

import (
	"context"

	"4d63.com/optional"
)

type option func(*Account)

type Repository interface {
	Create(ctx context.Context, handle string, opts ...option) (*Account, error)

	GetByID(ctx context.Context, id AccountID) (*Account, error)
	GetByHandle(ctx context.Context, handle string) (*Account, error)

	List(ctx context.Context, sort string, max, skip int) ([]Account, error)
}

func WithID(id AccountID) option {
	return func(a *Account) {
		a.ID = AccountID(id)
	}
}

func WithName(name string) option {
	return func(a *Account) {
		a.Name = name
	}
}

func WithBio(bio string) option {
	return func(a *Account) {
		a.Bio = optional.Of(bio)
	}
}
