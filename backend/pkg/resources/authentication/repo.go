package authentication

import (
	"context"

	"github.com/Southclaws/storyden/backend/pkg/resources/user"
)

type Repository interface {
	// Create an auth method for a user.
	Create(ctx context.Context,
		userID user.UserID,
		service Service,
		identifier string,
		token string,
		metadata map[string]any,
	) (*Authentication, error)

	// Gets an auth method based on a service's external account ID.
	GetByIdentifier(ctx context.Context, service Service, identifier string) (*Authentication, error)

	// Gets all auth methods that a user has.
	GetAuthMethods(ctx context.Context, userID user.UserID) ([]Authentication, error)

	// Checks if the given token is equal to the stored auth method's token.
	IsEqual(ctx context.Context, userID user.UserID, identifier string, token string) (bool, error)
}
