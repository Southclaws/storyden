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
	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/asset"
	"github.com/Southclaws/storyden/internal/ent/collection"
	"github.com/Southclaws/storyden/internal/ent/link"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/ent/react"
	"github.com/Southclaws/storyden/internal/ent/tag"
)

type database struct {
	db  *ent.Client
	raw *sqlx.DB
}

func New(db *ent.Client, raw *sqlx.DB) Repository {
	return &database{db, raw}
}

func (d *database) Create(
	ctx context.Context,
	title string,
	authorID account.AccountID,
	categoryID category.CategoryID,
	tags []string,
	opts ...Option,
) (*Thread, error) {
	// tagset, err := d.createTags(ctx, tags)
	// if err != nil {
	// 	return nil, fault.Wrap(err, "failed to upsert tags for linking to post")
	// }

	cat, err := d.db.Category.Get(ctx, xid.ID(categoryID))
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err,
				fctx.With(ctx),
				ftag.With(ftag.NotFound),
				fmsg.WithDesc("category not found",
					"The specified category was not found while creating the thread."))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	create := d.db.Post.Create()
	mutate := create.Mutation()

	for _, fn := range opts {
		fn(mutate)
	}

	mutate.SetTitle(title)
	mutate.SetFirst(true)
	mutate.SetAuthorID(xid.ID(authorID))
	mutate.SetTitle(title)
	mutate.SetCategoryID(cat.ID)

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
		WithAuthor(func(aq *ent.AccountQuery) {
			aq.WithRoles()
		}).
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

	return FromModel(nil, nil)(p)
}

// func (d *database) createTags(ctx context.Context, tags []string) ([]db.TagWhereParam, error) {
// 	setters := []db.TagWhereParam{}
// 	for _, tag := range tags {
// 		if len(tag) > 24 {
// 			return nil, post.ErrTagNameTooLong
// 		}
// 		_, err := d.db.Tag.
// 			UpsertOne(db.Tag.Name.Equals(tag)).
// 			Update().
// 			Create(db.Tag.Name.Set(tag)).
// 			Exec(ctx)
// 		if err != nil {
// 			return nil, fault.Wrap(err, "failed to upsert tag")
// 		}
// 		setters = append(setters, db.Tag.Name.Equals(tag))
// 	}
// 	return setters, nil
// }

func (d *database) Update(ctx context.Context, id post.ID, opts ...Option) (*Thread, error) {
	update := d.db.Post.UpdateOneID(xid.ID(id))
	mutate := update.Mutation()

	for _, fn := range opts {
		fn(mutate)
	}

	err := update.Exec(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	p, err := d.db.Post.
		Query().
		Where(ent_post.IDEQ(xid.ID(id))).
		WithAuthor(func(aq *ent.AccountQuery) {
			aq.WithRoles()
		}).
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

	return FromModel(nil, nil)(p)
}

const repliesCountManyQuery = `select
  p.id        post_id, -- post ID
  count(r.id) replies, -- number of replies,
  count(a.id) replied  -- has this account replied
from
  posts p
  inner join posts r on r.root_post_id = p.id
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
		WithAuthor(func(aq *ent.AccountQuery) {
			aq.WithRoles()
		}).
		WithAssets(func(aq *ent.AssetQuery) {
			aq.Order(asset.ByUpdatedAt(), asset.ByCreatedAt())
		}).
		WithCollections(func(cq *ent.CollectionQuery) {
			cq.WithOwner(func(aq *ent.AccountQuery) {
				aq.WithRoles()
			}).Order(collection.ByUpdatedAt(), collection.ByCreatedAt())
		}).
		WithLink(func(lq *ent.LinkQuery) {
			lq.WithFaviconImage().WithPrimaryImage()
			lq.WithAssets().Order(link.ByCreatedAt(sql.OrderDesc()))
		}).
		Order(ent_post.ByUpdatedAt(sql.OrderDesc()), ent_post.ByCreatedAt(sql.OrderDesc()))

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

	threads, err := dt.MapErr(result, FromModel(likes.Map(), replies.Map()))
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

func (d *database) Get(ctx context.Context, threadID post.ID, accountID opt.Optional[account.AccountID]) (*Thread, error) {
	var likes post.PostLikesResults
	err := d.raw.SelectContext(ctx, &likes, likesCountQuery, threadID.String(), accountID.String())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := d.db.Post.
		Query().
		Where(
			ent_post.First(true),
			ent_post.ID(xid.ID(threadID)),
		).
		WithPosts(func(pq *ent.PostQuery) {
			pq.
				Where(
					ent_post.DeletedAtIsNil(),
				).
				WithReplyTo(func(pq *ent.PostQuery) {
					pq.WithAuthor(func(aq *ent.AccountQuery) {
						aq.WithRoles()
					})
				}).
				WithReacts().
				WithAuthor(func(aq *ent.AccountQuery) {
					aq.WithRoles()
				}).
				WithAssets().
				WithLink(func(lq *ent.LinkQuery) {
					lq.WithFaviconImage().WithPrimaryImage()
					lq.WithAssets().Order(link.ByCreatedAt(sql.OrderDesc()))
				}).
				Order(ent.Asc(ent_post.FieldCreatedAt))
		}).
		WithAuthor(func(aq *ent.AccountQuery) {
			aq.WithRoles()
		}).
		WithCategory().
		WithTags(func(tq *ent.TagQuery) {
			tq.Order(tag.ByCreatedAt())
		}).
		WithReacts(func(rq *ent.ReactQuery) {
			rq.Order(react.ByCreatedAt())
		}).
		WithAssets().
		WithLink(func(lq *ent.LinkQuery) {
			lq.WithFaviconImage().WithPrimaryImage()
			lq.WithAssets().Order(link.ByCreatedAt(sql.OrderDesc()))
		}).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	replies := post.PostRepliesMap{
		xid.ID(threadID): post.PostRepliesResult{
			PostID: xid.ID(threadID),
			Count:  len(r.Edges.Posts),
			Replied: opt.Map(accountID, func(a account.AccountID) (replied int) {
				for _, p := range r.Edges.Posts {
					if p.Edges.Author.ID == xid.ID(a) {
						replied++
					}
				}
				return
			}).OrZero(),
		},
	}

	p, err := FromModel(likes.Map(), replies)(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

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
