package category

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"github.com/samber/lo"
	"go.uber.org/multierr"

	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/category"
	"github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/ent/predicate"
)

type database struct {
	db  *ent.Client
	raw *sqlx.DB
}

func New(db *ent.Client, raw *sqlx.DB) Repository {
	return &database{db, raw}
}

func (d *database) CreateCategory(ctx context.Context, name, desc, colour string, sort int, admin bool, opts ...Option) (*Category, error) {
	create := d.db.Category.Create()
	mutation := create.Mutation()

	mutation.SetName(name)
	mutation.SetSlug(slug.Make(name))
	mutation.SetDescription(desc)
	mutation.SetColour(colour)
	mutation.SetSort(sort)
	mutation.SetAdmin(admin)

	for _, fn := range opts {
		fn(mutation)
	}

	id, err := create.
		OnConflictColumns(category.FieldID).
		UpdateNewValues().
		ID(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.AlreadyExists))
		}

		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	c, err := d.db.Category.Get(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return FromModel(c), nil
}

const postsCountManyQuery = `select
  c.id        cat_id, -- category id
  count(p.id) threads -- thread count
from
  categories c
  inner join posts p on p.category_id = c.id
    and p.deleted_at is null
    and p.visibility = 'published'
group by
  c.id
`

type CategoryThreadsResult struct {
	CategoryID xid.ID `db:"cat_id"`
	PostCount  int    `db:"threads"`
}

type (
	CategoryThreadsResults []CategoryThreadsResult
	CategoryThreadsMap     map[xid.ID]CategoryThreadsResult
)

func (p CategoryThreadsResults) Map() CategoryThreadsMap {
	return lo.KeyBy(p, func(x CategoryThreadsResult) xid.ID { return x.CategoryID })
}

func (d *database) GetCategories(ctx context.Context, admin bool) ([]*Category, error) {
	filters := []predicate.Category{}

	if !admin {
		filters = append(filters, category.AdminEQ(false))
	}

	categories, err := d.db.Category.
		Query().
		Where(filters...).
		WithPosts(func(pq *ent.PostQuery) {
			pq.
				Where(
					post.FirstEQ(true),
					post.DeletedAtIsNil(),
					post.VisibilityEQ(post.VisibilityPublished),
				).
				WithAuthor(func(aq *ent.AccountQuery) {
					aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
				}).
				Limit(5).
				Order(ent.Desc(post.FieldUpdatedAt))
		}).
		Order(ent.Asc(category.FieldSort)).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if len(categories) == 0 {
		return []*Category{}, nil
	}

	var replies CategoryThreadsResults
	err = d.raw.SelectContext(ctx, &replies, postsCountManyQuery)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	categoryPosts := replies.Map()

	return dt.Map(categories, func(in *ent.Category) *Category {
		category := FromModel(in)
		category.PostCount = categoryPosts[in.ID].PostCount
		return category
	}), nil
}

func (d *database) Reorder(ctx context.Context, ids []CategoryID) ([]*Category, error) {
	cats, err := d.db.Category.Query().Where(category.Admin(false)).All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if len(ids) != len(cats) {
		return nil, fault.Wrap(
			fault.Newf("cannot reorder %d categories with %d ids, id list mismatch", len(cats), len(ids)),
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
		)
	}

	tx, err := d.db.Tx(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	newcats := []*Category{}

	for order, id := range ids {
		cat, err := tx.Category.UpdateOneID(xid.ID(id)).SetSort(order).Save(ctx)
		if err != nil {
			if rerr := tx.Rollback(); rerr != nil {
				return nil, fault.Wrap(multierr.Combine(err, rerr))
			}
			return nil, fault.Wrap(err)
		}

		newcats = append(newcats, FromModel(cat))
	}

	err = tx.Commit()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return newcats, nil
}

func (d *database) UpdateCategory(ctx context.Context, id CategoryID, opts ...Option) (*Category, error) {
	update := d.db.Category.UpdateOneID(xid.ID(id))
	mutation := update.Mutation()

	for _, fn := range opts {
		fn(mutation)
	}

	c, err := update.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return FromModel(c), nil
}

func (d *database) DeleteCategory(ctx context.Context, id CategoryID, moveto CategoryID) (*Category, error) {
	tx, err := d.db.Tx(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	defer func() {
		err = tx.Rollback()
	}()

	c, err := tx.Category.Get(ctx, xid.ID(id))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	_, err = tx.Post.Update().
		Where(post.CategoryID(xid.ID(id))).
		SetCategoryID(xid.ID(moveto)).
		Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = tx.Category.DeleteOneID(xid.ID(id)).Exec(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := tx.Commit(); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to perform move+delete transaction"))
	}

	return FromModel(c), nil
}
