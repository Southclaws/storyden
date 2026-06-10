package node_mutate

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type DeleteOptions struct {
	NewParent opt.Optional[library.QueryKey]
}

func (s *Manager) Delete(ctx context.Context, qk library.QueryKey, d DeleteOptions) (*library.Node, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	n, err := s.nodeQuerier.Get(ctx, qk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := session.Authorise(ctx, func() error {
		if n.Owner.ID != accountID {
			return fault.Wrap(rbac.ErrPermissions,
				fctx.With(ctx),
				fmsg.WithDesc("not owner", "You are not the owner of the page and do not have the Manage Library permission."))
		}
		return nil
	}, rbac.PermissionManageLibrary); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	destination, err := opt.MapErr(d.NewParent, func(target library.QueryKey) (library.Node, error) {
		destination, err := s.nc.Move(ctx, qk, target)
		if err != nil {
			return library.Node{}, fault.Wrap(err, fctx.With(ctx))
		}

		return *destination, fault.Wrap(err, fctx.With(ctx))
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = s.nodeWriter.Delete(ctx, qk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &rpc.EventNodeDeleted{
		ID:   library.NodeID(n.GetID()),
		Slug: n.GetSlug(),
	})

	return destination.Ptr(), nil
}
