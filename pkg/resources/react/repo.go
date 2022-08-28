package react

import (
	"context"

	"github.com/rs/xid"
)

type Repository interface {
	Add(ctx context.Context, userID xid.ID, postID xid.ID, emojiID string) (*React, error)
	Remove(ctx context.Context, userID xid.ID, reactID ReactID) (*React, error)
}
