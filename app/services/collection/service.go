package collection

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	"github.com/el-mike/restrict"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

type Service interface {
	Update(ctx context.Context, cid collection.CollectionID, partial Partial) (*collection.Collection, error)
	Delete(ctx context.Context, cid collection.CollectionID) error

	PostAdd(ctx context.Context, cid collection.CollectionID, pid post.ID) (*collection.Collection, error)
	PostRemove(ctx context.Context, cid collection.CollectionID, pid post.ID) (*collection.Collection, error)

	NodeAdd(ctx context.Context, cid collection.CollectionID, pid datagraph.NodeID) (*collection.Collection, error)
	NodeRemove(ctx context.Context, cid collection.CollectionID, pid datagraph.NodeID) (*collection.Collection, error)
}

type Partial struct {
	Name        opt.Optional[string]
	Description opt.Optional[string]
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l    *zap.Logger
	rbac rbac.AccessManager

	account_repo    account.Repository
	collection_repo collection.Repository
}

func New(
	l *zap.Logger,
	rbac rbac.AccessManager,

	account_repo account.Repository,
	collection_repo collection.Repository,
) Service {
	return &service{
		l:               l.With(zap.String("service", "collection")),
		rbac:            rbac,
		account_repo:    account_repo,
		collection_repo: collection_repo,
	}
}

func (s *service) Update(ctx context.Context, cid collection.CollectionID, partial Partial) (*collection.Collection, error) {
	if err := s.authorise(ctx, cid); err != nil {
		return nil, err
	}

	opts := []collection.Option{}

	partial.Name.Call(func(v string) { opts = append(opts, collection.WithName(v)) })
	partial.Description.Call(func(v string) { opts = append(opts, collection.WithDescription(v)) })

	col, err := s.collection_repo.Update(ctx, cid, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (s *service) Delete(ctx context.Context, cid collection.CollectionID) error {
	if err := s.authorise(ctx, cid); err != nil {
		return err
	}

	err := s.collection_repo.Delete(ctx, cid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (s *service) PostAdd(ctx context.Context, cid collection.CollectionID, pid post.ID) (*collection.Collection, error) {
	if err := s.authorise(ctx, cid); err != nil {
		return nil, err
	}

	col, err := s.collection_repo.Update(ctx, cid, collection.WithPostAdd(pid))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (s *service) PostRemove(ctx context.Context, cid collection.CollectionID, pid post.ID) (*collection.Collection, error) {
	if err := s.authorise(ctx, cid); err != nil {
		return nil, err
	}

	col, err := s.collection_repo.Update(ctx, cid, collection.WithPostRemove(pid))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (s *service) NodeAdd(ctx context.Context, cid collection.CollectionID, id datagraph.NodeID) (*collection.Collection, error) {
	if err := s.authorise(ctx, cid); err != nil {
		return nil, err
	}

	col, err := s.collection_repo.Update(ctx, cid, collection.WithNodeAdd(id))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (s *service) NodeRemove(ctx context.Context, cid collection.CollectionID, id datagraph.NodeID) (*collection.Collection, error) {
	if err := s.authorise(ctx, cid); err != nil {
		return nil, err
	}

	col, err := s.collection_repo.Update(ctx, cid, collection.WithNodeRemove(id))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (s *service) authorise(ctx context.Context, cid collection.CollectionID) error {
	aid, err := session.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := s.account_repo.GetByID(ctx, aid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	col, err := s.collection_repo.Get(ctx, cid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.rbac.Authorize(&restrict.AccessRequest{
		Subject:  acc,
		Resource: col,
		Actions:  []string{rbac.ActionUpdate},
	}); err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to authorize"))
	}

	return nil
}
