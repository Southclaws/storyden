package authentication

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
)

type Repository interface {
	// Create an auth method for a account.
	Create(ctx context.Context,
		userID account.AccountID,
		service Service,
		tokenType TokenType,
		identifier string,
		token string,
		metadata map[string]any,
		opts ...Option,
	) (*Authentication, error)

	// Gets an auth method based on a service's external account ID.
	LookupByIdentifier(ctx context.Context, service Service, identifier string) (*Authentication, bool, error)

	// Gets an auth method for a specific account based on a token type and identifier.
	LookupByTokenType(ctx context.Context, accountID account.AccountID, tokenType TokenType, identifier string) (*Authentication, bool, error)

	// Gets all auth methods that a account has.
	GetAuthMethods(ctx context.Context, userID account.AccountID) ([]*Authentication, error)

	Update(ctx context.Context, id ID, options ...Option) (*Authentication, error)

	DeleteByID(ctx context.Context, userID account.AccountID, aid ID) (bool, error)
	Delete(ctx context.Context, userID account.AccountID, identifier string, service Service) (bool, error)
}

type Option func(*ent.AuthenticationMutation)

func WithToken(token string) Option {
	return func(am *ent.AuthenticationMutation) {
		am.SetToken(token)
	}
}

func WithName(name string) Option {
	return func(am *ent.AuthenticationMutation) {
		am.SetName(name)
	}
}
