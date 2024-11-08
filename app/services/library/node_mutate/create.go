package node_mutate

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
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

	opts, err := s.applyOpts(ctx, p)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if v, ok := p.AssetSources.Get(); ok {
		for _, source := range v {
			a, err := s.fetcher.CopyAsset(ctx, source)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			opts = append(opts, node_writer.WithAssets([]asset.AssetID{a.ID}))
		}
	}

	nodeSlug := p.Slug.Or(mark.NewSlugFromName(name))

	if u, ok := p.URL.Get(); ok {
		ln, err := s.fetcher.Fetch(ctx, u)
		if err == nil {
			opts = append(opts, node_writer.WithLink(xid.ID(ln.ID)))
		}
	}

	if tags, ok := p.Tags.Get(); ok {
		newTags, err := s.tagWriter.Add(ctx, tags...)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		tagIDs := dt.Map(newTags, func(t *tag_ref.Tag) tag_ref.ID { return t.ID })

		opts = append(opts, node_writer.WithTagsAdd(tagIDs...))
	}

	n, err := s.nodeWriter.Create(ctx, owner, name, nodeSlug, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.indexQueue.Publish(ctx, mq.IndexNode{ID: library.NodeID(n.Mark.ID())}); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.fetcher.HydrateContentURLs(ctx, n)

	return n, nil
}
