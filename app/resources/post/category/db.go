package category

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/category"
	"github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/ent/predicate"
)

type Repository struct {
	db  *ent.Client
	raw *sqlx.DB
}

func New(db *ent.Client, raw *sqlx.DB) *Repository {
	return &Repository{db, raw}
}

type Option func(*ent.CategoryMutation)

func WithID(id CategoryID) Option {
	return func(cm *ent.CategoryMutation) {
		cm.SetID(xid.ID(id))
	}
}

func WithName(v string) Option {
	return func(cm *ent.CategoryMutation) {
		cm.SetName(v)
	}
}

func WithSlug(v string) Option {
	return func(cm *ent.CategoryMutation) {
		cm.SetSlug(v)
	}
}

func WithDescription(v string) Option {
	return func(cm *ent.CategoryMutation) {
		cm.SetDescription(v)
	}
}

func WithColour(v string) Option {
	return func(cm *ent.CategoryMutation) {
		cm.SetColour(v)
	}
}

func WithAdmin(v bool) Option {
	return func(cm *ent.CategoryMutation) {
		cm.SetAdmin(v)
	}
}

func WithMeta(v map[string]any) Option {
	return func(cm *ent.CategoryMutation) {
		cm.SetMetadata(v)
	}
}

func WithCoverImageAssetID(id *xid.ID) Option {
	return func(cm *ent.CategoryMutation) {
		if id == nil {
			cm.ClearCoverImage()
			return
		}
		cm.SetCoverImageAssetID(*id)
	}
}

func WithParent(id *CategoryID) Option {
	return func(cm *ent.CategoryMutation) {
		if id == nil {
			cm.ClearParent()
			return
		}

		cm.SetParentID(xid.ID(*id))
	}
}

type MoveOptions struct {
	ParentProvided bool
	ParentID       *CategoryID
	Before         *CategoryID
	After          *CategoryID
}

func (d *Repository) CreateCategory(ctx context.Context, name, desc, colour string, sort int, admin bool, opts ...Option) (*Category, error) {
	create := d.db.Category.Create()
	mutation := create.Mutation()

	mutation.SetName(name)
	mutation.SetSlug(mark.Slugify(name))
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

func (d *Repository) GetCategories(ctx context.Context, admin bool) ([]*Category, error) {
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
					post.RootPostIDIsNil(),
					post.DeletedAtIsNil(),
					post.VisibilityEQ(post.VisibilityPublished),
				).
				WithAuthor().
				Limit(5).
				Order(ent.Desc(post.FieldUpdatedAt))
		}).
		WithCoverImage(func(aq *ent.AssetQuery) {
			aq.WithParent()
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

func (d *Repository) Get(ctx context.Context, slug string) (*Category, error) {
	c, err := d.db.Category.
		Query().
		Where(category.SlugEQ(slug)).
		WithChildren().
		WithCoverImage(func(aq *ent.AssetQuery) {
			aq.WithParent()
		}).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var replies CategoryThreadsResults
	err = d.raw.SelectContext(ctx, &replies, postsCountManyQuery)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	categoryPosts := replies.Map()

	category := FromModel(c)

	// Set post count for this category.
	category.PostCount = categoryPosts[c.ID].PostCount

	// Set post count for child categories.
	category.Children = dt.Map(category.Children, func(child *Category) *Category {
		child.PostCount = categoryPosts[xid.ID(child.ID)].PostCount
		return child
	})

	return category, nil
}

func (d *Repository) UpdateCategory(ctx context.Context, slug string, opts ...Option) (*Category, error) {
	cat, err := d.db.Category.Query().Where(category.SlugEQ(slug)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	update := d.db.Category.UpdateOneID(cat.ID)
	mutation := update.Mutation()

	for _, fn := range opts {
		fn(mutation)
	}

	updated, err := update.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.AlreadyExists),
				fmsg.WithDesc("category slug already exists", "A category with this slug already exists."))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return d.Get(ctx, updated.Slug)
}

func (d *Repository) DeleteCategory(ctx context.Context, slug string, moveto CategoryID) (*Category, error) {
	tx, err := d.db.Tx(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	defer func() {
		err = tx.Rollback()
	}()

	c, err := tx.Category.Query().Where(category.SlugEQ(slug)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	_, err = tx.Post.Update().
		Where(post.CategoryID(c.ID)).
		SetCategoryID(xid.ID(moveto)).
		Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = tx.Category.DeleteOneID(c.ID).Exec(ctx)
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

func (d *Repository) MoveCategory(ctx context.Context, slug string, opts MoveOptions) ([]*Category, error) {
	if opts.Before != nil && opts.After != nil {
		return nil, fault.New("category move cannot specify both before and after", fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	tx, err := d.db.Tx(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	defer func() {
		err = tx.Rollback()
	}()

	cat, err := tx.Category.Query().Where(category.SlugEQ(slug)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sourceParentID := cat.ParentCategoryID
	targetParentID := sourceParentID

	if opts.ParentProvided {
		if opts.ParentID == nil {
			targetParentID = xid.ID{}
		} else {
			targetParentID = xid.ID(*opts.ParentID)
			if targetParentID == cat.ID {
				return nil, fault.New("category cannot be its own parent", fctx.With(ctx), ftag.With(ftag.InvalidArgument))
			}
			parentCat, err := tx.Category.Get(ctx, targetParentID)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
			}
			if err := ensureNoCycle(ctx, tx, cat.ID, parentCat); err != nil {
				return nil, err
			}
		}
	}

	if opts.Before != nil && *opts.Before == CategoryID(cat.ID) {
		return nil, fault.New("cannot move category before itself", fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}
	if opts.After != nil && *opts.After == CategoryID(cat.ID) {
		return nil, fault.New("cannot move category after itself", fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	if opts.ParentProvided && targetParentID != sourceParentID {
		upd := tx.Category.UpdateOneID(cat.ID)
		if opts.ParentID == nil {
			upd.ClearParent()
		} else {
			upd.SetParentID(xid.ID(*opts.ParentID))
		}
		if _, err := upd.Save(ctx); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if err := resequenceCategorySiblings(ctx, tx, sourceParentID, cat.Admin); err != nil {
			return nil, err
		}
	}

	siblings, err := listCategorySiblings(ctx, tx, targetParentID, cat.Admin)
	if err != nil {
		return nil, err
	}

	order := make([]CategoryID, 0, len(siblings))
	for _, sibling := range siblings {
		if sibling.ID == xid.ID(cat.ID) {
			continue
		}
		order = append(order, CategoryID(sibling.ID))
	}

	insertIndex := len(order)

	if opts.Before != nil {
		idx := indexOfCategory(order, *opts.Before)
		if idx == -1 {
			return nil, fault.New("before category not found in target parent", fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		insertIndex = idx
	}

	if opts.After != nil {
		idx := indexOfCategory(order, *opts.After)
		if idx == -1 {
			return nil, fault.New("after category not found in target parent", fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
		insertIndex = idx + 1
	}

	order = insertCategory(order, insertIndex, CategoryID(cat.ID))

	for idx, siblingID := range order {
		if _, err := tx.Category.UpdateOneID(xid.ID(siblingID)).SetSort(idx).Save(ctx); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return d.GetCategories(ctx, false)
}

func listCategorySiblings(ctx context.Context, tx *ent.Tx, parentID xid.ID, admin bool) ([]*ent.Category, error) {
	query := tx.Category.Query().Where(category.AdminEQ(admin))
	if parentID.IsNil() {
		query = query.Where(category.ParentCategoryIDIsNil())
	} else {
		query = query.Where(category.ParentCategoryID(parentID))
	}

	siblings, err := query.Order(ent.Asc(category.FieldSort), ent.Asc(category.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return siblings, nil
}

func resequenceCategorySiblings(ctx context.Context, tx *ent.Tx, parentID xid.ID, admin bool) error {
	siblings, err := listCategorySiblings(ctx, tx, parentID, admin)
	if err != nil {
		return err
	}

	for idx, sibling := range siblings {
		if sibling.Sort == idx {
			continue
		}
		if _, err := tx.Category.UpdateOneID(sibling.ID).SetSort(idx).Save(ctx); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	return nil
}

func ensureNoCycle(ctx context.Context, tx *ent.Tx, originalID xid.ID, parent *ent.Category) error {
	current := parent
	for {
		if current.ID == originalID {
			return fault.New("cannot move category into its descendant", fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}

		if current.ParentCategoryID.IsNil() {
			return nil
		}

		next, err := tx.Category.Query().Where(category.IDEQ(current.ParentCategoryID)).Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil
			}
			return fault.Wrap(err, fctx.With(ctx))
		}

		current = next
	}
}

func indexOfCategory(ids []CategoryID, target CategoryID) int {
	for idx, id := range ids {
		if id == target {
			return idx
		}
	}
	return -1
}

func insertCategory(ids []CategoryID, index int, id CategoryID) []CategoryID {
	if index < 0 {
		index = 0
	}
	if index > len(ids) {
		index = len(ids)
	}
	ids = append(ids, CategoryID{})
	copy(ids[index+1:], ids[index:])
	ids[index] = id
	return ids
}
