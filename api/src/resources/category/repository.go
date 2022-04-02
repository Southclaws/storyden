package category

import (
	"context"
)

type Repository interface {
	CreateCategory(ctx context.Context, name, desc, colour string, sort int, admin bool) (*Category, error)
	GetCategories(ctx context.Context, admin bool) ([]Category, error)
	UpdateCategory(ctx context.Context, id string, name, desc, colour *string, sort *int, admin *bool) (*Category, error)
	DeleteCategory(ctx context.Context, id string, moveto string) (*Category, error)
}
