package thread_querier

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/alitto/pond/v2"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/collection/collection_item_status"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reaction"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/internal/ent"
	ent_account "github.com/Southclaws/storyden/internal/ent/account"
	ent_asset "github.com/Southclaws/storyden/internal/ent/asset"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
	ent_react "github.com/Southclaws/storyden/internal/ent/react"
	ent_tag "github.com/Southclaws/storyden/internal/ent/tag"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/kv"
)

func (d *Querier) Get(ctx context.Context, threadID post.ID, pageParams pagination.Parameters, accountID opt.Optional[account.AccountID]) (*thread.Thread, error) {
	ctx, span := d.ins.Instrument(ctx,
		kv.String("thread_id", threadID.String()),
		kv.String("account_id", accountID.String()),
	)
	defer span.End()

	pool1 := pond.NewGroup()

	idList := []xid.ID{xid.ID(threadID)}

	var replyStatsMap post.PostRepliesMap
	pool1.SubmitErr(func() error {
		ctx, span := d.ins.InstrumentNamed(ctx, "replies_status")
		defer span.End()

		r, err := d.getRepliesStatus(ctx, idList, accountID.String())
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
		replyStatsMap = r
		return nil
	})

	var collectionsMap collection_item_status.CollectionStatusMap
	pool1.SubmitErr(func() error {
		ctx, span := d.ins.InstrumentNamed(ctx, "collections_status")
		defer span.End()

		r, err := d.getCollectionsStatus(ctx, idList, accountID.String())
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
		collectionsMap = r
		return nil
	})

	var readStateMap post.ReadStateMap
	pool1.SubmitErr(func() error {
		ctx, span := d.ins.InstrumentNamed(ctx, "read_status")
		defer span.End()

		r, err := d.getReadStatus(ctx, idList, accountID.String())
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
		readStateMap = r
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
				ent_post.RootPostID(xid.ID(threadID)),
			).
			Limit(pageParams.Limit()).
			Offset(pageParams.Offset()).
			Order(ent.Asc(ent_post.FieldCreatedAt)).
			WithReplyTo(func(rq *ent.PostQuery) {
				rq.Where(ent_post.DeletedAtIsNil())
			}).
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
				ent_post.RootPostIDIsNil(),
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
	for _, p := range repliesResult {
		if p.Edges.ReplyTo != nil {
			allPosts = append(allPosts, p.Edges.ReplyTo)
		}
	}
	postIDs := dt.Map(allPosts, func(p *ent.Post) xid.ID { return p.ID })

	accountIDs := dt.Map(allPosts, func(p *ent.Post) xid.ID { return p.AccountPosts })

	// Fetch dependent edges.

	reactResult, err := d.db.React.Query().
		Where(ent_react.PostIDIn(postIDs...)).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	likesMap, err := d.getLikesStatus(ctx, postIDs, accountID.String())
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
	threadMapper := thread.Mapper(accountLookup, readStateMap, likesMap, collectionsMap, replyStatsMap, reactLookup)

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
