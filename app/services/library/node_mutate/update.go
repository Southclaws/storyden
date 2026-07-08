package node_mutate

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_version"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func (s *Manager) Update(ctx context.Context, qk library.QueryKey, p Partial) (*library.Node, error) {
	return s.update(ctx, qk, p, opt.NewEmpty[xid.ID]())
}

func (s *Manager) UpdateFromVersion(ctx context.Context, qk library.QueryKey, p Partial, versionID node_version.VersionID) (*library.Node, error) {
	return s.update(ctx, qk, p, opt.New(xid.ID(versionID)))
}

func (s *Manager) update(ctx context.Context, qk library.QueryKey, p Partial, appliedVersion opt.Optional[xid.ID]) (*library.Node, error) {
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

	versionedMutation := p.HasVersionedFields()
	if versionedMutation && !appliedVersion.Ok() {
		hasDraft, err := s.versionQuerier.HasNodeDraft(ctx, qk)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if hasDraft {
			return nil, fault.New("node has a draft version",
				fctx.With(ctx),
				ftag.With(ftag.AlreadyExists),
				fmsg.WithDesc("draft exists", "This node has a working draft. Apply or delete the draft before editing versioned page fields directly."),
			)
		}
	}

	oldVisibility := n.Visibility
	previousSlug := n.GetSlug()

	pre, err := s.preMutation(ctx, p, opt.NewPtr(n))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if versionedMutation {
		if id, ok := appliedVersion.Get(); ok {
			pre.opts = append(pre.opts, node_writer.WithCurrentVersion(id))
		} else {
			pre.opts = append(pre.opts, node_writer.WithCurrentVersionCleared())
		}
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

	if err := s.cache.Invalidate(ctx, previousSlug); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
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
