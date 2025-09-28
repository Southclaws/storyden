package category

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/internal/deletable"
)

type Service interface {
	Create(ctx context.Context, name string, description string, colour string, admin bool) (*category.Category, error)
	Update(ctx context.Context, slug string, partial Partial) (*category.Category, error)
	Move(ctx context.Context, slug string, move Move) ([]*category.Category, error)
}

type Partial struct {
	Name        opt.Optional[string]
	Slug        opt.Optional[string]
	Description opt.Optional[string]
	Colour      opt.Optional[string]
	Admin       opt.Optional[bool]
	Meta        opt.Optional[map[string]any]
}

type Move struct {
	Parent deletable.Value[category.CategoryID]
	Before opt.Optional[category.CategoryID]
	After  opt.Optional[category.CategoryID]
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	accountQuery  *account_querier.Querier
	category_repo *category.Repository
}

func New(
	accountQuery *account_querier.Querier,
	category_repo *category.Repository,
) Service {
	return &service{
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

func (s *service) Update(ctx context.Context, slug string, partial Partial) (*category.Category, error) {
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

	cat, err := s.category_repo.UpdateCategory(ctx, slug, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return cat, nil
}

func (s *service) Move(ctx context.Context, slug string, move Move) ([]*category.Category, error) {
	parentOpt, deleteParent := move.Parent.Get()
	parentProvided := deleteParent
	var parentID *category.CategoryID

	if !deleteParent {
		if v, ok := parentOpt.Get(); ok {
			parentProvided = true
			pv := v
			parentID = &pv
		}
	}

	var beforeID *category.CategoryID
	if v, ok := move.Before.Get(); ok {
		pv := v
		beforeID = &pv
	}

	var afterID *category.CategoryID
	if v, ok := move.After.Get(); ok {
		pv := v
		afterID = &pv
	}

	cats, err := s.category_repo.MoveCategory(ctx, slug, category.MoveOptions{
		ParentProvided: parentProvided,
		ParentID:       parentID,
		Before:         beforeID,
		After:          afterID,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return cats, nil
}
