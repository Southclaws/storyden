package user

import (
	"context"
)

//go:generate mockery --inpackage --name=Repository --case=underscore

type Repository interface {
	CreateUser(ctx context.Context, email string, username string) (*User, error)

	GetUser(ctx context.Context, userId UserID, public bool) (*User, error)
	GetUserByEmail(ctx context.Context, email string, public bool) (*User, error)
	GetUsers(ctx context.Context, sort string, max, skip int, public bool) ([]User, error)

	UpdateUser(ctx context.Context, userId UserID, email, name, bio *string) (*User, error)
	SetAdmin(ctx context.Context, userId UserID, status bool) error

	Ban(ctx context.Context, userId UserID) (*User, error)
	Unban(ctx context.Context, userId UserID) (*User, error)
}
