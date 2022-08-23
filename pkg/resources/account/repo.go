package account

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, email string, username string) (*Account, error)

	GetByID(ctx context.Context, id AccountID) (*Account, error)
	LookupByEmail(ctx context.Context, email string) (*Account, bool, error)
	List(ctx context.Context, sort string, max, skip int) ([]Account, error)
}
