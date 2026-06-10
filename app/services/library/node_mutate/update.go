package node_mutate

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func (s *Manager) Update(ctx context.Context, qk library.QueryKey, p Partial) (*library.Node, error) {
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

	if err := s.cache.Invalidate(ctx, n.GetSlug()); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	oldVisibility := n.Visibility

	pre, err := s.preMutation(ctx, p, opt.NewPtr(n))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	n, err = s.nodeWriter.Update(ctx, qk, pre.opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if props, ok := p.Properties.Get(); ok {
		updatedProperties, err := s.applyPropertyMutations(ctx, n, props)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		n.Properties = opt.New(*updatedProperties)
	}

	// Emit update event
	s.bus.Publish(ctx, &rpc.EventNodeUpdated{
		ID:   library.NodeID(n.Mark.ID()),
		Slug: n.GetSlug(),
	})

	// Emit visibility transition events
	if oldVisibility != n.Visibility {
		switch n.Visibility {
		case visibility.VisibilityPublished:
			s.bus.Publish(ctx, &rpc.EventNodePublished{
				ID:   library.NodeID(n.Mark.ID()),
				Slug: n.GetSlug(),
			})

		case visibility.VisibilityReview:
			s.bus.Publish(ctx, &rpc.EventNodeSubmittedForReview{
				ID:   library.NodeID(n.Mark.ID()),
				Slug: n.GetSlug(),
			})

		case visibility.VisibilityUnlisted, visibility.VisibilityDraft, visibility.VisibilityReview:
			if oldVisibility == visibility.VisibilityPublished {
				s.bus.Publish(ctx, &rpc.EventNodeUnpublished{
					ID:   library.NodeID(n.Mark.ID()),
					Slug: n.GetSlug(),
				})
			}
		}
	}

	return n, nil
}
