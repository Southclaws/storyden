package collection_read

import (
	"context"
	"fmt"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/resources/collection/collection_querier"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

type Hydrator struct {
	querier *collection_querier.Querier
}

func New(
	querier *collection_querier.Querier,
) *Hydrator {
	return &Hydrator{
		querier: querier,
	}
}

func (r *Hydrator) GetCollection(ctx context.Context, qk collection.QueryKey) (*collection.CollectionWithItems, error) {
	acc := session.GetOptAccount(ctx).OrZero()

	col, err := r.querier.Get(ctx, qk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// The owner and admins can always read the unlisted collection items
	canReadUnlisted := session.GetRoles(ctx).Permissions().HasAny(rbac.PermissionAdministrator) || acc.ID == col.Owner.ID

	col.Items = dt.Filter(col.Items, func(i *collection.CollectionItem) bool {
		if canReadUnlisted {
			return true
		}

		var ownerID account.AccountID
		var vis visibility.Visibility

		switch n := i.Item.(type) {
		case *library.Node:
			vis = n.Visibility
			ownerID = n.Owner.ID

		case *post.Post:
			// TODO: Add visibility to post structure
			// vis = n.Visibility
			vis = visibility.VisibilityPublished
			ownerID = n.Author.ID

		default:
			panic(fmt.Sprintf("unsupported item type: %T", i.Item))
		}

		accountOwnsItem := ownerID == acc.ID

		// TODO: Apply to posts as well, but this needs some more work on post
		// data structure sharing and exposing visibility of the thread properly

		if vis == visibility.VisibilityDraft || i.MembershipType == collection.MembershipTypeSubmissionReview {
			// Don't reveal draft collection items unless the requesting account
			// is either the owner of the collection or the owner of the item.
			return accountOwnsItem
		}

		return true
	})

	return col, nil
}
