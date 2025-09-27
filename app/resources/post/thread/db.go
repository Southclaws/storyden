package thread

import (
	"context"
	"fmt"
	"math"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/alitto/pond/v2"
	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/collection/collection_item_status"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reaction"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/internal/ent"
	ent_account "github.com/Southclaws/storyden/internal/ent/account"
	ent_asset "github.com/Southclaws/storyden/internal/ent/asset"
	"github.com/Southclaws/storyden/internal/ent/category"
	"github.com/Southclaws/storyden/internal/ent/collection"
	"github.com/Southclaws/storyden/internal/ent/link"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
	ent_react "github.com/Southclaws/storyden/internal/ent/react"
	ent_tag "github.com/Southclaws/storyden/internal/ent/tag"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/kv"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/spanner"
)

type database struct {
	ins spanner.Instrumentation
	db  *ent.Client
	raw *sqlx.DB
}

func New(ins spanner.Builder, db *ent.Client, raw *sqlx.DB) Repository {
	return &database{
		ins: ins.Build(),
		db:  db,
		raw: raw,
	}
}

func (d *database) Create(
	ctx context.Context,
	title string,
	authorID account.AccountID,
	opts ...Option,
) (*Thread, error) {
	create := d.db.Post.Create()
	mutate := create.Mutation()
	mutate.SetUpdatedAt(time.Now())

	for _, fn := range opts {
		fn(mutate)
	}

	// If a category was specified, check if it exists first.
	if categoryID, ok := mutate.CategoryID(); ok {
		exists, err := d.db.Category.Query().Where(category.ID(xid.ID(categoryID))).Exist(ctx)
		if err != nil || !exists {
			if ent.IsNotFound(err) {
				return nil, fault.Wrap(err,
					fctx.With(ctx),
					ftag.With(ftag.NotFound),
					fmsg.WithDesc("category not found",
						"The specified category was not found while creating the thread."))
			}
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
		}
	}

	mutate.SetTitle(title)
	mutate.SetFirst(true)
	mutate.SetAuthorID(xid.ID(authorID))
	mutate.SetTitle(title)

	p, err := create.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.AlreadyExists))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	// Update the slug so it includes the ID for uniqueness.

	_, err = d.db.Post.
		UpdateOneID(p.ID).
		SetSlug(fmt.Sprintf("%s-%s", p.ID, slug.Make(title))).
		Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	// Finally, query the created thread with related entities.

	p, err = d.db.Post.
		Query().
		Where(ent_post.IDEQ(p.ID)).
		WithAuthor().
		WithCategory().
		WithTags().
		WithAssets().
		WithLink(func(lq *ent.LinkQuery) {
			lq.WithFaviconImage().WithPrimaryImage()
		}).
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return Map(p)
}

func (d *database) Update(ctx context.Context, id post.ID, opts ...Option) (*Thread, error) {
	update := d.db.Post.UpdateOneID(xid.ID(id))
	mutate := update.Mutation()

	for _, fn := range opts {
		fn(mutate)
	}

	// Only set the updated_at field if not changing the indexed_at field.
	if _, set := mutate.IndexedAt(); !set {
		mutate.SetUpdatedAt(time.Now())
	}

	err := update.Exec(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	p, err := d.db.Post.
		Query().
		Where(ent_post.IDEQ(xid.ID(id))).
		WithAuthor().
		WithCategory().
		WithTags().
		WithAssets().
		WithLink(func(lq *ent.LinkQuery) {
			lq.WithFaviconImage().WithPrimaryImage()
		}).
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return Map(p)
}

const repliesCountManyQuery = `select
  p.id        post_id, -- post ID
  count(r.id) replies, -- number of replies,
  count(a.id) replied  -- has this account replied
from
  posts p
  inner join posts r on r.root_post_id = p.id and r.deleted_at is null
  left join accounts a on a.id = p.account_posts and a.id = $1
group by p.id
`

const likesCountManyQuery = `select
  p.id        post_id, -- the post (thread or reply) ID
  count(*)    likes,   -- number of likes
  count(a.id) liked    -- has the account making the query liked this post?
from
  like_posts lp
  inner join posts p on p.id = lp.post_id
  left join accounts a on lp.account_id = a.id
  and a.id = $1
group by p.id
`

const collectionsCountManyQuery = `select
  p.id        item_id,          -- the post (thread or reply) ID
  count(*)    collections,      -- number of likes
  count(a.id) has_in_collection -- has the account making the query liked this post?
from
  collection_posts cp
  inner join posts p on p.id = cp.post_id
  inner join collections c on c.id = cp.collection_id
  left join accounts a on c.account_collections = a.id
  and a.id = $1
group by p.id
`

func (d *database) List(
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

	query := d.db.Post.Query().Where(ent_post.First(true))

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
		Order(func(s *sql.Selector) {
			s.OrderBy(fmt.Sprintf("COALESCE(%s, %s) DESC",
				s.C(ent_post.FieldLastReplyAt),
				s.C(ent_post.FieldCreatedAt)),
			)
		})

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

	isNextPage := len(result) >= size
	nextPage := opt.NewSafe(page+1, isNextPage)
	totalPages := int(math.Ceil(float64(total) / float64(size)))

	if isNextPage {
		result = result[:len(result)-1]
	}

	var replies post.PostRepliesResults
	err = d.raw.SelectContext(ctx, &replies, repliesCountManyQuery, accountID.String())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var likes post.PostLikesResults
	err = d.raw.SelectContext(ctx, &likes, likesCountManyQuery, accountID.String())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var collections collection_item_status.CollectionStatusResults
	err = d.raw.SelectContext(ctx, &collections, collectionsCountManyQuery, accountID.String())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	mapper := Mapper(nil, likes.Map(), collections.Map(), replies.Map(), nil)
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

const repliesCountQuery = `select
  p.id        post_id, -- post ID
  count(r.id) replies, -- number of replies,
  count(a.id) replied  -- has this account replied
from
  posts p
  inner join posts r on r.root_post_id = p.id and r.deleted_at is null
  left join accounts a on a.id = r.account_posts and a.id = $2
where
  p.id = $1 or p.root_post_id = $1
group by p.id
`

const likesCountQuery = `select
  p.id        post_id, -- the post (thread or reply) ID
  count(*)    likes,   -- number of likes
  count(a.id) liked    -- has the account making the query liked this post?
from
  like_posts lp
  inner join posts p on p.id = lp.post_id
  left join accounts a on lp.account_id = a.id and a.id = $2
where
  p.id = $1 or p.root_post_id = $1
group by p.id
`

const collectionsCountQuery = `select
  p.id        item_id,          -- the post (thread or reply) ID
  count(*)    collections,      -- number of likes
  count(a.id) has_in_collection -- has the account making the query liked this post?
from
  collection_posts cp
  inner join posts p on p.id = cp.post_id
  inner join collections c on c.id = cp.collection_id
  left join accounts a on c.account_collections = a.id
  and a.id = $2
where
  p.id = $1 or p.root_post_id = $1
group by p.id
`

func (d *database) Get(ctx context.Context, threadID post.ID, pageParams pagination.Parameters, accountID opt.Optional[account.AccountID]) (*Thread, error) {
	ctx, span := d.ins.Instrument(ctx,
		kv.String("thread_id", threadID.String()),
		kv.String("account_id", accountID.String()),
	)
	defer span.End()

	pool1 := pond.NewGroup()

	var replyStatsMap post.PostRepliesMap
	pool1.SubmitErr(func() error {
		ctx, span := d.ins.InstrumentNamed(ctx, "replies_status")
		defer span.End()

		var replyStats post.PostRepliesResults
		err := d.raw.SelectContext(ctx, &replyStats, repliesCountQuery, threadID.String(), accountID.String())
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
		replyStatsMap = replyStats.Map()
		return nil
	})

	var likesMap post.PostLikesMap
	pool1.SubmitErr(func() error {
		ctx, span := d.ins.InstrumentNamed(ctx, "likes_status")
		defer span.End()

		var likes post.PostLikesResults
		err := d.raw.SelectContext(ctx, &likes, likesCountQuery, threadID.String(), accountID.String())
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		likesMap = likes.Map()

		return nil
	})

	var collectionsMap collection_item_status.CollectionStatusMap
	pool1.SubmitErr(func() error {
		ctx, span := d.ins.InstrumentNamed(ctx, "collections_status")
		defer span.End()

		var collections collection_item_status.CollectionStatusResults
		err := d.raw.SelectContext(ctx, &collections, collectionsCountQuery, threadID.String(), accountID.String())
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		collectionsMap = collections.Map()

		return nil
	})

	var tags tag_ref.Tags
	pool1.SubmitErr(func() error {
		ctx, span := d.ins.InstrumentNamed(ctx, "thread_tags")
		defer span.End()

		tagsResult, err := d.db.Tag.Query().Where(ent_tag.HasPostsWith(ent_post.ID(xid.ID(threadID)))).All(ctx)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		tags = dt.Map(tagsResult, tag_ref.Map(nil))

		return nil
	})

	var assets []*asset.Asset
	pool1.SubmitErr(func() error {
		ctx, span := d.ins.InstrumentNamed(ctx, "thread_assets")
		defer span.End()

		r, err := d.db.Asset.Query().Where(ent_asset.HasPostsWith(ent_post.ID(xid.ID(threadID)))).All(ctx)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		assets = dt.Map(r, asset.Map)

		return nil
	})

	var repliesResult []*ent.Post
	pool1.SubmitErr(func() error {
		ctx, span := d.ins.InstrumentNamed(ctx, "thread_replies")
		defer span.End()

		r, err := d.db.Post.Query().
			Where(
				ent_post.DeletedAtIsNil(),
				ent_post.First(false),
				ent_post.RootPostID(xid.ID(threadID)),
			).
			Limit(pageParams.Limit()).
			Offset(pageParams.Offset()).
			Order(ent.Asc(ent_post.FieldCreatedAt)).
			All(ctx)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		repliesResult = r

		return nil
	})

	var threadResult *ent.Post
	pool1.SubmitErr(func() error {
		ctx, span := d.ins.InstrumentNamed(ctx, "thread_root")
		defer span.End()

		r, err := d.db.Post.Query().
			Where(
				ent_post.DeletedAtIsNil(),
				ent_post.First(true),
				ent_post.ID(xid.ID(threadID)),
			).
			WithCategory().
			WithLink(func(lq *ent.LinkQuery) {
				lq.WithFaviconImage().WithPrimaryImage()
			}).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
			}

			return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
		}

		threadResult = r

		return nil
	})

	// Wait for first stage to complete.
	err := pool1.Wait()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	allPosts := append(repliesResult, threadResult)
	postIDs := dt.Map(allPosts, func(p *ent.Post) xid.ID { return p.ID })

	accountIDs := dt.Map(allPosts, func(p *ent.Post) xid.ID { return p.AccountPosts })

	// Fetch dependent edges.

	reactResult, err := d.db.React.Query().
		Where(ent_react.PostIDIn(postIDs...)).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// React lookup contributes to the account query.
	reacters := dt.Map(reactResult, func(r *ent.React) xid.ID { return r.AccountID })
	accountIDs = append(accountIDs, reacters...)

	accountIDs = lo.Uniq(accountIDs)

	// Lookup all accounts relevant to this thread.
	var accountLookup account.Lookup
	accountEdges, err := d.db.Account.Query().
		Where(ent_account.IDIn(accountIDs...)).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	accountLookup = account.NewAccountLookup(accountEdges)

	// Join all data together

	reacts, err := dt.MapErr(reactResult, reaction.Mapper(accountLookup))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	reactLookup := reaction.Reacts(reacts).Map()

	replyMapper := reply.Mapper(accountLookup, likesMap, reactLookup)
	threadMapper := Mapper(accountLookup, likesMap, collectionsMap, replyStatsMap, reactLookup)

	replies, err := dt.MapErr(repliesResult, replyMapper)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	p, err := threadMapper(threadResult)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	totalReplies := replyStatsMap[threadResult.ID].Count
	repliesPage := pagination.NewPageResult(pageParams, totalReplies, replies)

	p.Replies = repliesPage
	p.Tags = tags
	p.Assets = assets

	return p, nil
}

func (d *database) Delete(ctx context.Context, id post.ID) error {
	err := d.db.Post.
		UpdateOneID(xid.ID(id)).
		SetDeletedAt(time.Now()).
		Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to archive thread root post"))
	}

	err = d.db.Post.
		Update().
		Where(ent_post.RootPostID(xid.ID(id))).
		SetDeletedAt(time.Now()).
		Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to archive thread posts"))
	}

	return nil
}
