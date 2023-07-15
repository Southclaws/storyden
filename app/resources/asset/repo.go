package asset

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/account"
)

type Repository interface {
	Add(ctx context.Context,
		owner account.AccountID,
		id, url, mt string,
		width, height int,
	) (*Asset, error)

	Remove(ctx context.Context, owner account.AccountID, id AssetID) error
}
