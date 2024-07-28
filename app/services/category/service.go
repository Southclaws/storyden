package category

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

var errNotAuthorised = fault.Wrap(fault.New("not authorised"), ftag.With(ftag.PermissionDenied))

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
	l    *zap.Logger
	rbac rbac.AccessManager

	accountQuery  account_querier.Querier
	category_repo category.Repository
}

func New(
	l *zap.Logger,
	rbac rbac.AccessManager,

	accountQuery account_querier.Querier,
	category_repo category.Repository,
) Service {
	return &service{
		l:             l.With(zap.String("service", "collection")),
		rbac:          rbac,
		accountQuery:  accountQuery,
		category_repo: category_repo,
	}
}

func (s *service) Create(ctx context.Context, name string, description string, colour string, admin bool) (*category.Category, error) {
	if err := s.authorise(ctx); err != nil {
		return nil, err
	}

	cat, err := s.category_repo.CreateCategory(ctx, name, description, colour, 0, admin)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return cat, nil
}

func (s *service) Reorder(ctx context.Context, ids []category.CategoryID) ([]*category.Category, error) {
	if err := s.authorise(ctx); err != nil {
		return nil, err
	}

	cats, err := s.category_repo.Reorder(ctx, ids)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return cats, nil
}

func (s *service) Update(ctx context.Context, id category.CategoryID, partial Partial) (*category.Category, error) {
	if err := s.authorise(ctx); err != nil {
		return nil, err
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

func (s *service) authorise(ctx context.Context) error {
	aid, err := session.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := s.accountQuery.GetByID(ctx, aid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if !acc.Admin {
		return fault.Wrap(errNotAuthorised, fctx.With(ctx))
	}

	return nil
}
