package category

import (
	"context"
)

type option func(*Category)

type Repository interface {
	CreateCategory(ctx context.Context,
		name string,
		desc string,
		colour string,
		sort int,
		admin bool,
		opts ...option) (*Category, error)

	GetCategories(ctx context.Context, admin bool) ([]*Category, error)
	Reorder(ctx context.Context, ids []CategoryID) ([]*Category, error)
	UpdateCategory(ctx context.Context, id CategoryID, name, desc, colour *string, sort *int, admin *bool) (*Category, error)
	DeleteCategory(ctx context.Context, id CategoryID, moveto CategoryID) (*Category, error)
}

func WithID(id CategoryID) option {
	return func(c *Category) {
		c.ID = id
	}
}
