package collection_read

import (
	"context"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/reply"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/account/session"
	"github.com/Southclaws/storyden/app/services/semdex"
)

type CollectionQuerier struct {
	fx.In

	Logger  *zap.Logger
	Repo    collection.Repository
	Semdex  semdex.Retriever
	Session session.SessionProvider
}

func (r *CollectionQuerier) GetCollection(ctx context.Context, id collection.CollectionID) (*collection.Collection, error) {
	acc := r.Session.AccountOpt(ctx).OrZero()

	col, err := r.Repo.Get(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// The owner and admins can always read the unlisted collection items
	canReadUnlisted := acc.Admin || acc.ID == col.Owner.ID

	col.Items = dt.Filter(col.Items, func(i *collection.CollectionItem) bool {
		if canReadUnlisted {
			return true
		}

		var ownerID account.AccountID
		var vis visibility.Visibility

		switch n := i.Item.(type) {
		case *datagraph.Node:
			vis = n.Visibility
			ownerID = n.Owner.ID

		case *reply.Reply:
			// TODO: Add visibility to reply structure
			// vis = n.Visibility
			vis = visibility.VisibilityPublished
			ownerID = n.Author.ID

		default:
			panic(fmt.Sprintf("unsupported item type: %T", i.Item))
		}

		accountOwnsItem := ownerID == acc.ID

		// TODO: Apply to posts as well, but this needs some more work on post
		// data structure sharing and exposing visibility of the thread properly

		if vis != visibility.VisibilityPublished || i.MembershipType == collection.MembershipTypeSubmission {
			// Don't reveal unlisted collection items unless the requesting acc
			// is either the owner of the collection or the owner of the item.

			// If the owner of the node is the requesting account, show it.
			return accountOwnsItem
		}

		return true
	})

	return col, nil
}
