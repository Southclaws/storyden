package collection_read

import (
	"context"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/resources/datagraph/semdex"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/account/session"
)

type CollectionQuerier struct {
	fx.In

	Logger  *zap.Logger
	Repo    collection.Repository
	Semdex  semdex.RelevanceScorer
	Session session.SessionProvider
}

func (r *CollectionQuerier) GetCollection(ctx context.Context, id collection.CollectionID) (*collection.CollectionWithItems, error) {
	session := r.Session.AccountOpt(ctx)
	acc := session.OrZero()

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

	if acc, ok := session.Get(); ok && r.Semdex != nil {
		pro := profile.ProfileFromAccount(&acc)
		ids := dt.Map(col.Items, func(i *collection.CollectionItem) xid.ID { return i.Item.GetID() })

		scores, err := r.Semdex.ScoreRelevance(ctx, pro, ids...)
		if err != nil {
			r.Logger.Warn("failed to score relevance", zap.Error(err))
		}

		col.Items = dt.Map(col.Items, func(i *collection.CollectionItem) *collection.CollectionItem {
			score, ok := scores[i.Item.GetID()]
			if ok {
				i.RelevanceScore = opt.New(score)
			}

			return i
		})
	}

	return col, nil
}
