package react

import (
	"context"
)

type Repository interface {
	Add(ctx context.Context, userID, postID, emojiID string) (*React, error)
	Remove(ctx context.Context, userID, reactID string) (*React, error)
}
