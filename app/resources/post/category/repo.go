package category

import (
	"context"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/ent"
)

type Option func(*ent.CategoryMutation)

type Repository interface {
	CreateCategory(ctx context.Context,
		name string,
		desc string,
		colour string,
		sort int,
		admin bool,
		opts ...Option) (*Category, error)

	GetCategories(ctx context.Context, admin bool) ([]*Category, error)
	Reorder(ctx context.Context, ids []CategoryID) ([]*Category, error)
	UpdateCategory(ctx context.Context, id CategoryID, opts ...Option) (*Category, error)
	DeleteCategory(ctx context.Context, id CategoryID, moveto CategoryID) (*Category, error)
}

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
