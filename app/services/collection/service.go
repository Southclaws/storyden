package collection

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

var ErrNotResourcesOwner = fault.New("not the owner of either collection or item", ftag.With(ftag.PermissionDenied))

type Service interface {
	Update(ctx context.Context, cid collection.CollectionID, partial Partial) (*collection.CollectionWithItems, error)
	Delete(ctx context.Context, cid collection.CollectionID) error

	PostAdd(ctx context.Context, cid collection.CollectionID, pid post.ID) (*collection.CollectionWithItems, error)
	PostRemove(ctx context.Context, cid collection.CollectionID, pid post.ID) (*collection.CollectionWithItems, error)

	NodeAdd(ctx context.Context, cid collection.CollectionID, pid library.NodeID) (*collection.CollectionWithItems, error)
	NodeRemove(ctx context.Context, cid collection.CollectionID, pid library.NodeID) (*collection.CollectionWithItems, error)
}

type Partial struct {
	Name        opt.Optional[string]
	Description opt.Optional[string]
}

func Build() fx.Option {
	return fx.Provide(New)
}

type service struct {
	l *zap.Logger

	accountQuery *account_querier.Querier
	repo         collection.Repository
}

func New(
	l *zap.Logger,

	accountQuery *account_querier.Querier,
	repo collection.Repository,
) Service {
	return &service{
		l: l.With(zap.String("service", "collection")),

		accountQuery: accountQuery,
		repo:         repo,
	}
}

func (s *service) Update(ctx context.Context, cid collection.CollectionID, partial Partial) (*collection.CollectionWithItems, error) {
	if err := s.authoriseDirectUpdate(ctx, cid); err != nil {
		return nil, err
	}

	opts := []collection.Option{}

	partial.Name.Call(func(v string) { opts = append(opts, collection.WithName(v)) })
	partial.Description.Call(func(v string) { opts = append(opts, collection.WithDescription(v)) })

	col, err := s.repo.Update(ctx, cid, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (s *service) Delete(ctx context.Context, cid collection.CollectionID) error {
	if err := s.authoriseDirectUpdate(ctx, cid); err != nil {
		return err
	}

	err := s.repo.Delete(ctx, cid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (s *service) PostAdd(ctx context.Context, cid collection.CollectionID, pid post.ID) (*collection.CollectionWithItems, error) {
	err, mt := s.authoriseSubmission(ctx, cid, xid.ID(pid))
	if err != nil {
		return nil, err
	}

	col, err := s.repo.UpdateItems(ctx, cid, collection.WithPost(pid, mt))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (s *service) PostRemove(ctx context.Context, cid collection.CollectionID, pid post.ID) (*collection.CollectionWithItems, error) {
	col, err := s.repo.UpdateItems(ctx, cid, collection.WithPostRemove(pid))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (s *service) NodeAdd(ctx context.Context, cid collection.CollectionID, id library.NodeID) (*collection.CollectionWithItems, error) {
	err, mt := s.authoriseSubmission(ctx, cid, xid.ID(id))
	if err != nil {
		return nil, err
	}

	col, err := s.repo.UpdateItems(ctx, cid, collection.WithNode(id, mt))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (s *service) NodeRemove(ctx context.Context, cid collection.CollectionID, id library.NodeID) (*collection.CollectionWithItems, error) {
	if err := s.authoriseDirectUpdate(ctx, cid); err != nil {
		return nil, err
	}

	col, err := s.repo.UpdateItems(ctx, cid, collection.WithNodeRemove(id))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return col, nil
}

func (s *service) authoriseDirectUpdate(ctx context.Context, cid collection.CollectionID) error {
	aid, err := session.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := s.accountQuery.GetByID(ctx, aid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	col, err := s.repo.Get(ctx, cid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := acc.Roles.Permissions().Authorise(ctx, func() error {
		if acc.ID != col.Owner.ID {
			return fault.Wrap(rbac.ErrPermissions,
				fctx.With(ctx),
				fmsg.WithDesc("not owner", "You are not the owner of the collection and do not have the Manage Collections permission."),
			)
		}
		return nil
	}, rbac.PermissionManageCollections); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (s *service) authoriseSubmission(ctx context.Context, cid collection.CollectionID, iid xid.ID) (error, collection.MembershipType) {
	aid, err := session.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx)), collection.MembershipType{}
	}

	acc, err := s.accountQuery.GetByID(ctx, aid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx)), collection.MembershipType{}
	}

	col, err := s.repo.ProbeItem(ctx, cid, iid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx)), collection.MembershipType{}
	}

	if item, ok := col.Item.Get(); ok && item.Author.ID != acc.ID {
		return nil, collection.MembershipTypeSubmissionAccepted
	}

	if col.Collection.Owner.ID != acc.ID {
		return nil, collection.MembershipTypeSubmissionReview
	}

	if acc.ID == col.Collection.Owner.ID {
		return nil, collection.MembershipTypeNormal
	}

	return fault.Wrap(ErrNotResourcesOwner, fctx.With(ctx)), collection.MembershipType{}
}
