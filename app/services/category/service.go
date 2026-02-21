package category

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/post/category_cache"
	"github.com/Southclaws/storyden/internal/deletable"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

var errInvalidCategoryCreate = fault.New("invalid create args", ftag.With(ftag.InvalidArgument))

type Service interface {
	Create(ctx context.Context, partial Partial) (*category.Category, error)
	Update(ctx context.Context, slug string, partial Partial) (*category.Category, error)
	Move(ctx context.Context, slug string, move Move) ([]*category.Category, error)
	Delete(ctx context.Context, slug string, moveToID category.CategoryID) (*category.Category, error)
}

type Partial struct {
	Name              opt.Optional[string]
	Slug              opt.Optional[string]
	Description       opt.Optional[string]
	Colour            opt.Optional[string]
	Parent            opt.Optional[category.CategoryID]
	CoverImageAssetID deletable.Value[*xid.ID]
	Meta              opt.Optional[map[string]any]
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
	cache         *category_cache.Cache
	bus           *pubsub.Bus
}

func New(
	accountQuery *account_querier.Querier,
	category_repo *category.Repository,
	cache *category_cache.Cache,
	bus *pubsub.Bus,
) Service {
	return &service{
		accountQuery:  accountQuery,
		category_repo: category_repo,
		cache:         cache,
		bus:           bus,
	}
}

func (s *service) Create(ctx context.Context, partial Partial) (*category.Category, error) {
	opts := []category.Option{}

	if v, ok := partial.Parent.Get(); ok {
		pid := v
		opts = append(opts, category.WithParent(&pid))
	}

	if v, ok := partial.Slug.Get(); ok {
		opts = append(opts, category.WithSlug(v))
	}

	coverImage, _ := partial.CoverImageAssetID.Get()
	if v, ok := coverImage.Get(); ok {
		opts = append(opts, category.WithCoverImageAssetID(v))
	}

	if v, ok := partial.Meta.Get(); ok {
		opts = append(opts, category.WithMeta(v))
	}

	name, ok := partial.Name.Get()
	if !ok {
		return nil, fault.Wrap(errInvalidCategoryCreate, fctx.With(ctx), fmsg.WithDesc("missing name", "Category name is required."))
	}

	description, ok := partial.Description.Get()
	if !ok {
		return nil, fault.Wrap(errInvalidCategoryCreate, fctx.With(ctx), fmsg.WithDesc("missing description", "Category description is required."))
	}

	colour, ok := partial.Colour.Get()
	if !ok {
		return nil, fault.Wrap(errInvalidCategoryCreate, fctx.With(ctx), fmsg.WithDesc("missing colour", "Category colour is required."))
	}

	cat, err := s.category_repo.CreateCategory(ctx, name, description, colour, 0, false, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &rpc.EventCategoryUpdated{Slug: cat.Slug})

	return cat, nil
}

func (s *service) Update(ctx context.Context, slug string, partial Partial) (*category.Category, error) {
	if err := s.cache.Invalidate(ctx, slug); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

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
	coverImageOpt, shouldDelete := partial.CoverImageAssetID.Get()
	if shouldDelete {
		opts = append(opts, category.WithCoverImageAssetID(nil))
	} else if v, ok := coverImageOpt.Get(); ok {
		opts = append(opts, category.WithCoverImageAssetID(v))
	}
	if v, ok := partial.Meta.Get(); ok {
		opts = append(opts, category.WithMeta(v))
	}

	cat, err := s.category_repo.UpdateCategory(ctx, slug, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &rpc.EventCategoryUpdated{Slug: cat.Slug})

	return cat, nil
}

func (s *service) Move(ctx context.Context, slug string, move Move) ([]*category.Category, error) {
	if err := s.cache.Invalidate(ctx, slug); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

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

	s.bus.Publish(ctx, &rpc.EventCategoryUpdated{Slug: slug})

	return cats, nil
}

func (s *service) Delete(ctx context.Context, slug string, moveToID category.CategoryID) (*category.Category, error) {
	cat, err := s.category_repo.DeleteCategory(ctx, slug, moveToID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &rpc.EventCategoryDeleted{Slug: slug})

	return cat, nil
}
