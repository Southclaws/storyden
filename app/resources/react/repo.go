package react

import (
	"context"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
)

type Repository interface {
	Add(ctx context.Context, accountID account.AccountID, postID xid.ID, emojiID string) (*React, error)
	Remove(ctx context.Context, accountID account.AccountID, reactID ReactID) (*React, error)
}
