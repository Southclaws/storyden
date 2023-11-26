package thread

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/gosimple/slug"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/asset"
	"github.com/Southclaws/storyden/internal/ent/collection"
	"github.com/Southclaws/storyden/internal/ent/link"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/ent/predicate"
	"github.com/Southclaws/storyden/internal/ent/react"
	"github.com/Southclaws/storyden/internal/ent/tag"
)

type database struct {
	db *ent.Client
}

func New(db *ent.Client) Repository {
	return &database{db}
}

func (d *database) Create(
	ctx context.Context,
	title string,
	body string,
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
	mutate.SetBody(body)
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
		WithAuthor().
		WithCategory().
		WithTags().
		WithAssets().
		WithLinks().
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return FromModel(p)
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
		WithAuthor().
		WithCategory().
		WithTags().
		WithAssets().
		WithLinks().
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return FromModel(p)
}

func (d *database) List(
	ctx context.Context,
	before time.Time,
	max int,
	opts ...Query,
) ([]*Thread, error) {
	filters := []predicate.Post{
		ent_post.DeletedAtIsNil(),
		ent_post.First(true),
	}

	if !before.IsZero() {
		filters = append(filters, ent_post.CreatedAtLT(before))
	}

	if max < 1 {
		max = 1
	}

	if max > 100 {
		max = 100
	}

	query := d.db.Post.Query().
		Where(filters...).
		Limit(max).
		WithCategory().
		WithAuthor().
		WithAssets(func(aq *ent.AssetQuery) {
			aq.Order(asset.ByUpdatedAt(), asset.ByCreatedAt())
		}).
		WithCollections(func(cq *ent.CollectionQuery) {
			cq.WithOwner().Order(collection.ByUpdatedAt(), collection.ByCreatedAt())
		}).
		WithLinks(func(lq *ent.LinkQuery) {
			lq.WithAssets().Order(link.ByCreatedAt(sql.OrderDesc()))
		}).
		Order(ent_post.ByUpdatedAt(sql.OrderDesc()), ent_post.ByCreatedAt(sql.OrderDesc()))

	for _, fn := range opts {
		fn(query)
	}

	result, err := query.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	// counts, err := d.GetPostCounts(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	threads, err := dt.MapErr(result, FromModel)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return threads, nil
}

func (d *database) Get(ctx context.Context, threadID post.ID) (*Thread, error) {
	post, err := d.db.Post.
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
					pq.WithAuthor()
				}).
				WithReacts().
				WithAuthor().
				WithAssets().
				WithLinks(func(lq *ent.LinkQuery) {
					lq.WithAssets().Order(link.ByCreatedAt(sql.OrderDesc()))
				}).
				Order(ent.Asc(ent_post.FieldCreatedAt))
		}).
		WithAuthor().
		WithCategory().
		WithTags(func(tq *ent.TagQuery) {
			tq.Order(tag.ByCreatedAt())
		}).
		WithReacts(func(rq *ent.ReactQuery) {
			rq.Order(react.ByCreatedAt())
		}).
		WithAssets().
		WithLinks(func(lq *ent.LinkQuery) {
			lq.WithAssets().Order(link.ByCreatedAt(sql.OrderDesc()))
		}).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return FromModel(post)
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
