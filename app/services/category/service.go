package category

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/post/category"
)

type Service interface {
	Create(ctx context.Context, name string, description string, colour string, admin bool) (*category.Category, error)
	Reorder(ctx context.Context, ids []category.CategoryID) ([]*category.Category, error)
	Update(ctx context.Context, id category.CategoryID, partial Partial) (*category.Category, error)
}

type Partial struct {
	Name        opt.Optional[string]
	Slug        opt.Optional[string]
	Description opt.Optional[string]
	Colour      opt.Optional[string]
	Admin       opt.Optional[bool]
	Meta        opt.Optional[map[string]any]
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l *zap.Logger

	accountQuery  *account_querier.Querier
	category_repo category.Repository
}

func New(
	l *zap.Logger,

	accountQuery *account_querier.Querier,
	category_repo category.Repository,
) Service {
	return &service{
		l: l.With(zap.String("service", "collection")),

		accountQuery:  accountQuery,
		category_repo: category_repo,
	}
}

func (s *service) Create(ctx context.Context, name string, description string, colour string, admin bool) (*category.Category, error) {
	cat, err := s.category_repo.CreateCategory(ctx, name, description, colour, 0, admin)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return cat, nil
}

func (s *service) Reorder(ctx context.Context, ids []category.CategoryID) ([]*category.Category, error) {
	cats, err := s.category_repo.Reorder(ctx, ids)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return cats, nil
}

func (s *service) Update(ctx context.Context, id category.CategoryID, partial Partial) (*category.Category, error) {
	opts := []category.Option{}

	if v, ok := partial.Name.Get(); ok {
		opts = append(opts, category.WithName(v))
	}
	if v, ok := partial.Slug.Get(); ok {
		opts = append(opts, category.WithSlug(v))
	}
	if v, ok := partial.Description.Get(); ok {
		opts = append(opts, category.WithDescription(v))
	}
	if v, ok := partial.Colour.Get(); ok {
		opts = append(opts, category.WithColour(v))
	}
	if v, ok := partial.Admin.Get(); ok {
		opts = append(opts, category.WithAdmin(v))
	}
	if v, ok := partial.Meta.Get(); ok {
		opts = append(opts, category.WithMeta(v))
	}

	cat, err := s.category_repo.UpdateCategory(ctx, id, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return cat, nil
}
