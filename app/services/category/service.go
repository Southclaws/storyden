package category

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication"
)

type Service interface {
	Create(ctx context.Context, name string, description string, colour string, admin bool) (*category.Category, error)
	Reorder(ctx context.Context, ids []category.CategoryID) ([]*category.Category, error)
}

type Partial struct {
	Name        opt.Optional[string]
	Description opt.Optional[string]
	Colour      opt.Optional[string]
	Admin       opt.Optional[bool]
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l    *zap.Logger
	rbac rbac.AccessManager

	account_repo  account.Repository
	category_repo category.Repository
}

func New(
	l *zap.Logger,
	rbac rbac.AccessManager,

	account_repo account.Repository,
	category_repo category.Repository,
) Service {
	return &service{
		l:             l.With(zap.String("service", "collection")),
		rbac:          rbac,
		account_repo:  account_repo,
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

func (s *service) authorise(ctx context.Context) error {
	aid, err := authentication.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := s.account_repo.GetByID(ctx, aid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if !acc.Admin {
		return fault.New("")
	}

	return nil
}
