package category

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/Southclaws/storyden/api/src/infra/db/model"
	"github.com/Southclaws/storyden/api/src/infra/db/model/category"
	"github.com/Southclaws/storyden/api/src/infra/db/model/post"
	"github.com/Southclaws/storyden/api/src/infra/db/model/predicate"
)

type database struct {
	db *model.Client
}

func New(db *model.Client) Repository {
	return &database{db}
}

func (d *database) CreateCategory(ctx context.Context, name, desc, colour string, sort int, admin bool) (*Category, error) {
	id, err := d.db.Category.
		Create().
		SetName(name).
		SetDescription(desc).
		SetColour(colour).
		SetSort(sort).
		SetAdmin(admin).
		OnConflictColumns(category.FieldID).
		UpdateNewValues().
		ID(ctx)
	if err != nil {
		return nil, err
	}

	c, err := d.db.Category.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return FromModel(c), nil
}

func (d *database) GetCategories(ctx context.Context, admin bool) ([]Category, error) {
	filters := []predicate.Category{}

	if !admin {
		filters = append(filters, category.AdminEQ(false))
	}

	categories, err := d.db.Category.
		Query().
		Where(filters...).
		WithPosts(func(pq *model.PostQuery) {
			pq.
				Where(
					post.FirstEQ(true),
					post.DeletedAtIsNil(),
				).
				WithAuthor().
				Limit(5).
				Order(model.Desc(post.FieldUpdatedAt))
		}).
		Order(model.Asc(category.FieldSort)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return []Category{}, nil
	}

	// NOTE:
	// Lazy two queries because Ent doesn't yet support Count aggregations.
	// I could write the above query as raw SQL too but... screw that. Joins are
	// super annoying to get nested data out of because SQL is awful. So for now
	// there are two separate queries and the data is joined below. Besides,
	// there won't be many categories anyway so it's not going to affect
	// performance much.
	type CategoryPostCount struct {
		ID    uuid.UUID `json:"id"`
		Posts int       `json:"posts"`
	}

	var categoryPostsList []CategoryPostCount

	err = d.db.Category.Query().Modify(func(s *sql.Selector) {
		s.
			Select(
				sql.As(s.C("id"), "id"),
				sql.As(sql.Count("*"), "posts"),
			).
			Join(sql.Table(post.Table).As("p")).On(s.C(post.FieldID), "category_id").
			GroupBy(s.C("id")).
			OrderBy(sql.Desc("posts"))
	}).Scan(ctx, &categoryPostsList)
	if err != nil {
		return nil, err
	}

	categoryPosts := make(map[uuid.UUID]int)
	for _, c := range categoryPostsList {
		categoryPosts[c.ID] = c.Posts
	}

	result := []Category{}

	for _, c := range categories {
		category := FromModel(c)
		category.PostCount = categoryPosts[c.ID]
		result = append(result, *category)
	}

	return result, nil
}

func (d *database) UpdateCategory(ctx context.Context, id CategoryID, name, desc, colour *string, sort *int, admin *bool) (*Category, error) {
	u := d.db.Category.UpdateOneID(uuid.UUID(id))

	// TODO: Write a less explicit, more ergonomic way to do this:

	//nocheck:wsl
	if name != nil {
		u.SetName(*name)
	}
	if desc != nil {
		u.SetDescription(*desc)
	}
	if colour != nil {
		u.SetColour(*colour)
	}
	if sort != nil {
		u.SetSort(*sort)
	}
	if admin != nil {
		u.SetAdmin(*admin)
	}

	c, err := u.Save(ctx)
	if err != nil {
		return nil, err
	}

	return FromModel(c), nil
}

func (d *database) DeleteCategory(ctx context.Context, id CategoryID, moveto CategoryID) (*Category, error) {
	tx, err := d.db.Tx(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	c, err := tx.Category.Get(ctx, uuid.UUID(id))
	if err != nil {
		return nil, err
	}

	_, err = tx.Post.Update().Where(post.CategoryID(uuid.UUID(id))).SetCategoryID(uuid.UUID(moveto)).Save(ctx)
	if err != nil {
		return nil, err
	}

	err = tx.Category.DeleteOneID(uuid.UUID(id)).Exec(ctx)
	if err != nil {
		return nil, err
	}

	tx.Commit()
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform move+delete transaction")
	}

	return FromModel(c), nil
}
