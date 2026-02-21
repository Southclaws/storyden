package node_mutate

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/library/node_auth"
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

	acc, err := s.accountQuery.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	n, err := s.nodeQuerier.Get(ctx, qk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := node_auth.AuthoriseNodeMutation(ctx, acc, n); err != nil {
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
