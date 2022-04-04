package category

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Southclaws/storyden/api/src/infra/db"
	"github.com/Southclaws/storyden/api/src/infra/db/model"
)

type DB struct {
	db *model.Client
}

func New(db *model.Client) Repository {
	return &DB{db}
}

func (d *DB) CreateCategory(ctx context.Context, name, desc, colour string, sort int, admin bool) (*Category, error) {
	c, err := d.db.Category.
		UpsertOne(db.Category.Name.Equals(name)).
		Create(
			db.Category.Name.Set(name),
			db.Category.Description.Set(desc),
			db.Category.Colour.Set(colour),
			db.Category.Sort.Set(sort),
			db.Category.Admin.Set(admin),
		).
		Update(
			db.Category.Name.Set(name),
			db.Category.Description.Set(desc),
			db.Category.Colour.Set(colour),
			db.Category.Sort.Set(sort),
			db.Category.Admin.Set(admin),
		).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return FromModel(c), nil
}

func (d *DB) GetCategories(ctx context.Context, admin bool) ([]Category, error) {
	filters := []db.CategoryWhereParam{}

	if !admin {
		filters = append(filters, db.Category.Admin.Equals(false))
	}

	categories, err := d.db.Category.
		FindMany(filters...).
		OrderBy(
			db.Category.Sort.Order(db.SortOrderAsc),
		).
		With(
			db.Category.Posts.
				Fetch(
					db.Post.First.Equals(true),
					db.Post.DeletedAt.IsNull(),
				).
				With(
					db.Post.Author.Fetch(),
				).
				Take(5).
				OrderBy(db.Post.UpdatedAt.Order(db.SortOrderDesc)),
		).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return nil, nil
	}

	// NOTE:
	// Lazy two queries because Prisma doesn't yet support Count aggregations.
	// I could write the above query as raw SQL too but... screw that. Joins are
	// super annoying to get nested data out of because SQL is awful. So for now
	// there are two separate queries and the data is joined below. Besides,
	// there won't be many categories anyway so it's not going to affect
	// performance much.
	type CategoryPostCount struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Sort  string `json:"sort"`
		Posts int    `json:"posts"`
	}
	var categoryPostsList []CategoryPostCount
	d.db.Prisma.QueryRaw(`
		select c.id, c.name, c.sort, count(*) as posts
		from "Category" c
		inner join "Post" p on c.id = p."categoryId"
		group by c.id
		order by c.sort asc
	`).Exec(ctx, &categoryPostsList)
	categoryPosts := make(map[string]int)
	for _, c := range categoryPostsList {
		categoryPosts[c.ID] = c.Posts
	}

	result := []Category{}
	for _, c := range categories {
		category := FromModel(&c)
		category.PostCount = categoryPosts[c.ID]
		result = append(result, *category)
	}

	return result, nil
}

func (d *DB) UpdateCategory(ctx context.Context, id string, name, desc, colour *string, sort *int, admin *bool) (*Category, error) {
	c, err := d.db.Category.
		FindUnique(db.Category.ID.Equals(id)).
		Update(
			db.Category.Name.SetIfPresent(name),
			db.Category.Description.SetIfPresent(desc),
			db.Category.Colour.SetIfPresent(colour),
			db.Category.Sort.SetIfPresent(sort),
			db.Category.Admin.SetIfPresent(admin),
		).
		Exec(ctx)

	return FromModel(c), err
}

func (d *DB) DeleteCategory(ctx context.Context, id string, moveto string) (*Category, error) {
	movePosts := d.db.Post.
		FindMany(
			db.Post.Category.Where(
				db.Category.ID.Equals(id),
			),
		).
		Update(
			db.Post.CategoryID.Set(moveto),
		).
		Tx()

	deleteCategory := d.db.Category.
		FindUnique(
			db.Category.ID.Equals(id),
		).
		Delete().
		Tx()

	err := d.db.Prisma.Transaction(
		movePosts,
		deleteCategory,
	).Exec(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform move+delete transaction")
	}

	return FromModel(deleteCategory.Result()), nil
}
