package authentication

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/account"
)

type Repository interface {
	// Create an auth method for a account.
	Create(ctx context.Context,
		userID account.AccountID,
		service Service,
		identifier string,
		token string,
		metadata map[string]any,
	) (*Authentication, error)

	// Gets an auth method based on a service's external account ID.
	LookupByIdentifier(ctx context.Context, service Service, identifier string) (*Authentication, bool, error)

	// Gets an auth method based on a service and the account's handle.
	LookupByHandle(ctx context.Context, service Service, handle string) (*Authentication, bool, error)

	// Gets all auth methods that a account has.
	GetAuthMethods(ctx context.Context, userID account.AccountID) ([]*Authentication, error)

	// Checks if the given token is equal to the stored auth method's token.
	IsEqual(ctx context.Context, userID account.AccountID, identifier string, token string) (bool, error)

	Delete(ctx context.Context, userID account.AccountID, identifier string, service Service) (bool, error)
}
