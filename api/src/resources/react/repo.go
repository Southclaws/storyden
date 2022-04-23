package react

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Add(ctx context.Context, userID uuid.UUID, postID uuid.UUID, emojiID string) (*React, error)
	Remove(ctx context.Context, userID uuid.UUID, reactID ReactID) (*React, error)
}
