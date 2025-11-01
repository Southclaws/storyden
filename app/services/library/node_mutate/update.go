package node_mutate

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/library/node_auth"
)

func (s *Manager) Update(ctx context.Context, qk library.QueryKey, p Partial) (*library.Node, error) {
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
	s.bus.Publish(ctx, &message.EventNodeUpdated{
		ID:   library.NodeID(n.Mark.ID()),
		Mark: n.Mark.String(),
	})

	// Emit visibility transition events
	if oldVisibility != n.Visibility {
		switch n.Visibility {
		case visibility.VisibilityPublished:
			s.bus.Publish(ctx, &message.EventNodePublished{
				ID:   library.NodeID(n.Mark.ID()),
				Mark: n.Mark.String(),
			})

		case visibility.VisibilityReview:
			s.bus.Publish(ctx, &message.EventNodeSubmittedForReview{
				ID:   library.NodeID(n.Mark.ID()),
				Mark: n.Mark.String(),
			})

		case visibility.VisibilityUnlisted, visibility.VisibilityDraft, visibility.VisibilityReview:
			if oldVisibility == visibility.VisibilityPublished {
				s.bus.Publish(ctx, &message.EventNodeUnpublished{
					ID:   library.NodeID(n.Mark.ID()),
					Mark: n.Mark.String(),
				})
			}
		}
	}

	return n, nil
}
