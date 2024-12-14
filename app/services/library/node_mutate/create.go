package node_mutate

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/visibility"
)

func (s *Manager) Create(ctx context.Context,
	owner account.AccountID,
	name string,
	p Partial,
) (*library.Node, error) {
	if v, ok := p.Visibility.Get(); ok {
		if v == visibility.VisibilityPublished {
			acc, err := s.accountQuery.GetByID(ctx, owner)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			if err := acc.Roles.Permissions().Authorise(ctx, nil, rbac.PermissionManageLibrary); err != nil {
				return nil, fault.Wrap(err,
					fctx.With(ctx),
					fmsg.WithDesc("non admin cannot publish nodes", "You do not have permission to publish, please submit as draft, review or unlisted."),
				)
			}
		}
	}

	pre, err := s.preMutation(ctx, p, opt.NewEmpty[library.Node]())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	opts := pre.opts

	nodeSlug := p.Slug.Or(mark.NewSlugFromName(name))

	n, err := s.nodeWriter.Create(ctx, owner, name, nodeSlug, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if p.Visibility.OrZero() == visibility.VisibilityPublished {
		if err := s.indexQueue.Publish(ctx, mq.IndexNode{ID: library.NodeID(n.Mark.ID())}); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	s.fetcher.HydrateContentURLs(ctx, n)

	return n, nil
}
