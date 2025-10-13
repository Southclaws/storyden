package thread_querier

import (
	"context"
	"math"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/alitto/pond/v2"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/collection/collection_item_status"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/internal/ent"
	ent_asset "github.com/Southclaws/storyden/internal/ent/asset"
	"github.com/Southclaws/storyden/internal/ent/collection"
	"github.com/Southclaws/storyden/internal/ent/link"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
)

func (d *Querier) List(
	ctx context.Context,
	page int,
	size int,
	accountID opt.Optional[account.AccountID],
	opts ...Query,
) (*Result, error) {
	if size < 1 {
		size = 1
	}

	if size > 100 {
		size = 100
	}

	query := d.db.Post.Query().Where(ent_post.RootPostIDIsNil())

	for _, fn := range opts {
		fn(query)
	}

	query.
		WithCategory().
		WithAuthor().
		WithAssets(func(aq *ent.AssetQuery) {
			aq.Order(ent_asset.ByUpdatedAt(), ent_asset.ByCreatedAt())
		}).
		WithCollections(func(cq *ent.CollectionQuery) {
			cq.WithOwner().Order(collection.ByUpdatedAt(), collection.ByCreatedAt())
		}).
		WithLink(func(lq *ent.LinkQuery) {
			lq.WithFaviconImage().WithPrimaryImage()
			lq.WithAssets().Order(link.ByCreatedAt(sql.OrderDesc()))
		}).
		Order(ent.Desc(ent_post.FieldLastReplyAt))

	total, err := query.Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	query.
		Limit(size + 1).
		Offset(page * size)

	result, err := query.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	isNextPage := len(result) > size
	nextPage := opt.NewSafe(page+1, isNextPage)
	totalPages := int(math.Ceil(float64(total) / float64(size)))

	if len(result) == 0 {
		return &Result{
			PageSize:    size,
			Results:     0,
			TotalPages:  totalPages,
			CurrentPage: page,
			NextPage:    nextPage,
			Threads:     []*thread.Thread{},
		}, nil
	}

	if isNextPage {
		result = result[:len(result)-1]
	}

	ids := dt.Map(result, func(p *ent.Post) xid.ID { return p.ID })

	pool := pond.NewGroup()

	var readStates post.ReadStateMap
	pool.SubmitErr(func() error {
		r, err := d.getReadStatus(ctx, ids, accountID.String())
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
		readStates = r
		return nil
	})

	var repliesMap post.PostRepliesMap
	pool.SubmitErr(func() error {
		r, err := d.getRepliesStatus(ctx, ids, accountID.String())
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
		repliesMap = r
		return nil
	})

	var likesMap post.PostLikesMap
	pool.SubmitErr(func() error {
		r, err := d.getLikesStatus(ctx, ids, accountID.String())
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
		likesMap = r
		return nil
	})

	var collectionsMap collection_item_status.CollectionStatusMap
	pool.SubmitErr(func() error {
		r, err := d.getCollectionsStatus(ctx, ids, accountID.String())
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
		collectionsMap = r
		return nil
	})

	err = pool.Wait()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	mapper := thread.Mapper(nil, readStates, likesMap, collectionsMap, repliesMap, nil)
	threads, err := dt.MapErr(result, mapper)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &Result{
		PageSize:    size,
		Results:     len(threads),
		TotalPages:  totalPages,
		CurrentPage: page,
		NextPage:    nextPage,
		Threads:     threads,
	}, nil
}
