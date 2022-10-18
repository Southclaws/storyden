package thread

import (
	"context"
	"fmt"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault/errctx"
	"github.com/Southclaws/fault/errtag"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model"
	post_model "github.com/Southclaws/storyden/internal/infrastructure/db/model/post"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/predicate"
	"github.com/Southclaws/storyden/internal/utils"
)

var (
	ErrNoTitle      = errors.New("missing title")
	ErrNoBody       = errors.New("missing body")
	ErrUnauthorised = errors.New("unauthorised")
)

type database struct {
	db *model.Client
}

func New(db *model.Client) Repository {
	return &database{db}
}

func (d *database) Create(
	ctx context.Context,
	title string,
	body string,
	authorID account.AccountID,
	categoryID category.CategoryID,
	tags []string,
	opts ...option,
) (*Thread, error) {
	insert := Thread{
		Short: post.MakeShortBody(body),
		Title: title,
	}

	for _, v := range opts {
		v(&insert)
	}

	// tagset, err := d.createTags(ctx, tags)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "failed to upsert tags for linking to post")
	// }

	cat, err := d.db.Category.Get(ctx, xid.ID(categoryID))
	if err != nil {
		if model.IsNotFound(err) {
			return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.NotFound{})
		}

		return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.Internal{})
	}

	p, err := d.db.Post.
		Create().
		SetNillableID(utils.OptionalID(xid.ID(insert.ID))).
		SetFirst(true).
		SetShort(insert.Short).
		SetBody(body).
		SetAuthorID(xid.ID(authorID)).
		SetTitle(title).
		SetCategory(cat).
		// AddTagIDs(tagset).
		Save(ctx)
	if err != nil {
		if model.IsConstraintError(err) {
			return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.AlreadyExists{})
		}

		return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.Internal{})
	}

	// Update the slug so it includes the ID for uniqueness.

	_, err = d.db.Post.
		UpdateOneID(p.ID).
		SetSlug(fmt.Sprintf("%s-%s", p.ID, slug.Make(title))).
		Save(ctx)
	if err != nil {
		return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.Internal{})
	}

	// Finally, query the created thread with related entities.

	p, err = d.db.Post.
		Query().
		Where(post_model.IDEQ(p.ID)).
		WithAuthor().
		WithCategory().
		WithTags().
		Only(ctx)
	if err != nil {
		return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.Internal{})
	}

	return FromModel(p), nil
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
// 			return nil, errors.Wrap(err, "failed to upsert tag")
// 		}
// 		setters = append(setters, db.Tag.Name.Equals(tag))
// 	}
// 	return setters, nil
// }

func (d *database) List(
	ctx context.Context,
	before time.Time,
	max int,
	opts ...Query,
) ([]*Thread, error) {
	filters := []predicate.Post{
		post_model.DeletedAtIsNil(),
		post_model.First(true),
	}

	if !before.IsZero() {
		filters = append(filters, post_model.CreatedAtLT(before))
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
		WithAuthor()

	for _, fn := range opts {
		fn(query)
	}

	result, err := query.All(ctx)
	if err != nil {
		return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.Internal{})
	}

	// counts, err := d.GetPostCounts(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	return dt.Map(result, FromModel), nil
}

func (d *database) Get(ctx context.Context, threadID post.PostID) (*Thread, error) {
	post, err := d.db.Post.
		Query().
		Where(
			post_model.First(true),
			post_model.ID(xid.ID(threadID)),
		).
		WithPosts(func(pq *model.PostQuery) {
			pq.
				WithReplyTo(func(pq *model.PostQuery) {
					pq.WithAuthor()
				}).
				WithReacts().
				WithAuthor().
				Order(model.Asc(post_model.FieldCreatedAt))
		}).
		WithAuthor().
		WithCategory().
		WithTags().
		WithReacts().
		Only(ctx)
	if err != nil {
		if model.IsNotFound(err) {
			return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.NotFound{})
		}

		return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.Internal{})
	}

	return FromModel(post), nil
}

// func (d *database) GetPostCounts(ctx context.Context) (map[string]int, error) {
// 	type PostCount struct {
// 		PostID string `json:"rootPostId"`
// 		Count  int    `json:"count"`
// 	}

// 	var counts []PostCount
// 	err := d.db.Prisma.Raw.QueryRaw(`
// 		with recursive counts AS(
// 			select id, "rootPostId"
// 			from public."Post"
// 			where "rootPostId" is not null

// 			union

// 			select s.id, s."rootPostId"
// 			from public."Post" s
// 			inner join counts c on c.id = s."rootPostId"
// 		) select "rootPostId", count(*) from counts
// 		group by "rootPostId"`).
// 		Exec(ctx, &counts)
// 	if err != nil {
// 		return nil, err
// 	}

// 	result := make(map[string]int)
// 	for _, c := range counts {
// 		result[c.PostID] = c.Count
// 	}

// 	return result, nil
// }

// func (d *database) Update(ctx context.Context, userID, id string, title, categoryID *string, pinned *bool) (*post.Post, error) {
// 	updates := []db.PostSetParam{
// 		db.Post.Title.SetIfPresent(title),
// 		db.Post.CategoryID.SetIfPresent(categoryID),
// 		db.Post.Pinned.SetIfPresent(pinned),
// 	}

// 	if err := post.CanUserMutatePost(ctx, d.db, userID, id); err != nil {
// 		return nil, post.ErrUnauthorised
// 	}

// 	p, err := d.db.Post.
// 		FindUnique(db.Post.ID.Equals(id)).
// 		With(
// 			db.Post.Author.Fetch(),
// 			db.Post.Category.Fetch(),
// 			db.Post.Tags.Fetch(),
// 		).
// 		Update(updates...).
// 		Exec(ctx)
// 	if err != nil {
// 		if errors.Is(err, db.ErrNotFound) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	return post.FromModel(p), nil
// }

// func (d *database) Delete(ctx context.Context, id, authorID string) (int, error) {
// 	// NOTE:
// 	// We really want this authorID to eventually be removed from this API.
// 	// Because this API should be user-agnostic, and should be usable by non-
// 	// human users. Therefore, the validation of access rights should happen at
// 	// a different abstraction layer. Lower than the HTTP API but higher than
// 	// the database implementation.
// 	if err := post.CanUserMutatePost(ctx, d.db, authorID, id); err != nil {
// 		return 0, errors.Wrap(err, "failed to check user permissions")
// 	}

// 	result, err := d.db.Post.FindMany(
// 		db.Post.Or(
// 			db.Post.And(
// 				db.Post.First.Equals(true),
// 				db.Post.ID.Equals(id),
// 			),
// 			db.Post.And(
// 				db.Post.First.Equals(false),
// 				db.Post.Root.Where(db.Post.ID.Equals(id)),
// 			),
// 		),
// 	).Update(
// 		db.Post.DeletedAt.Set(time.Now()),
// 	).Exec(ctx)
// 	if err != nil {
// 		if errors.Is(err, db.ErrNotFound) {
// 			return 0, nil
// 		}
// 		return 0, errors.Wrap(err, "failed to set deletedAt for posts")
// 	}

// 	return result.Count, nil
// }
