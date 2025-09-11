package node_mutate

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
)

type preMutationResult struct {
	opts []node_writer.Option
}

// preMutation constructs node_writer options for a create or partial update.
func (s *Manager) preMutation(ctx context.Context, p Partial, current opt.Optional[library.Node]) (*preMutationResult, error) {
	opts := []node_writer.Option{}

	// Apply all primitive options. These are just basic partial updates.
	p.Name.Call(func(value string) { opts = append(opts, node_writer.WithName(value)) })
	p.Slug.Call(func(value mark.Slug) { opts = append(opts, node_writer.WithSlug(value.String())) })
	p.Description.Call(func(value string) { opts = append(opts, node_writer.WithDescription(value)) })
	p.PrimaryImage.Call(func(value xid.ID) {
		opts = append(opts, node_writer.WithPrimaryImage(value))
	}, func() {
		opts = append(opts, node_writer.WithPrimaryImageRemoved())
	})
	p.Content.Call(func(value datagraph.Content) { opts = append(opts, node_writer.WithContent(value)) })
	p.Metadata.Call(func(value map[string]any) { opts = append(opts, node_writer.WithMetadata(value)) })
	p.AssetsAdd.Call(func(value []asset.AssetID) { opts = append(opts, node_writer.WithAssets(value)) })
	p.AssetsRemove.Call(func(value []asset.AssetID) { opts = append(opts, node_writer.WithAssetsRemoved(value)) })
	p.Visibility.Call(func(value visibility.Visibility) { opts = append(opts, node_writer.WithVisibility(value)) })
	p.HideChildren.Call(func(value bool) { opts = append(opts, node_writer.WithHideChildren(value)) })

	// If the mutation includes a parent node, we need to query it because the
	// WithParent API only accepts a node ID, not a node mark (slug or ID).
	if parentSlug, ok := p.Parent.Get(); ok {
		parent, err := s.nodeQuerier.Get(ctx, parentSlug)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		opts = append(opts, node_writer.WithParent(library.NodeID(parent.Mark.ID())))
	}

	// If the mutation includes asset sources (so, URLs to assets to be added)
	// download them and append them to the node's asset list.
	if v, ok := p.AssetSources.Get(); ok {
		o, err := s.buildAssetSourcesOpts(ctx, v)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
		opts = append(opts, o...)
	}

	// If there's a URL being applied, fetch it and if it returns new content,
	// set that as the content to be used for fill rules.
	linkUrl, removeLinkURL := p.URL.Get()
	if u, ok := linkUrl.Get(); ok && !removeLinkURL {
		ln, _, err := s.fetcher.ScrapeAndStore(ctx, u)
		if err == nil {
			opts = append(opts, node_writer.WithLink(xid.ID(ln.ID)))
		}
	} else if removeLinkURL {
		opts = append(opts, node_writer.WithLinkRemove())
	}

	// -
	// Saving tags
	//
	// Happens after any generative options due to how tag-fill rules may yield
	// new tags that need to be saved before being applied to the node. But even
	// if there's no tag-fill rule set this is applied to any tags in the patch.
	// -

	if t, ok := p.Tags.Get(); ok {
		n, ok := current.Get()
		if ok {
			tagOpts, err := s.createDeleteTagsForExistingNode(ctx, &n, t)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
			opts = append(opts, tagOpts...)
		} else {
			tagOpts, err := s.createDeleteTagsForNewNode(ctx, t)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
			opts = append(opts, tagOpts...)
		}
	}

	return &preMutationResult{
		opts: opts,
	}, nil
}

func (s *Manager) buildAssetSourcesOpts(ctx context.Context, sources []string) ([]node_writer.Option, error) {
	opts := []node_writer.Option{}

	for _, source := range sources {
		a, err := s.fetcher.CopyAsset(ctx, source)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		opts = append(opts, node_writer.WithAssets([]asset.AssetID{a.ID}))
	}

	return opts, nil
}

func (s *Manager) createDeleteTagsForNewNode(ctx context.Context, tags tag_ref.Names) ([]node_writer.Option, error) {
	opts := []node_writer.Option{}

	newTags, err := s.tagWriter.Add(ctx, tags...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	addIDs := dt.Map(newTags, func(t *tag_ref.Tag) tag_ref.ID { return t.ID })

	opts = append(opts, node_writer.WithTagsAdd(addIDs...))

	return opts, nil
}

func (s *Manager) createDeleteTagsForExistingNode(ctx context.Context, n *library.Node, tags tag_ref.Names) ([]node_writer.Option, error) {
	opts := []node_writer.Option{}

	currentTagNames := n.Tags.Names()

	toCreate, toRemove := lo.Difference(tags, currentTagNames)

	newTags, err := s.tagWriter.Add(ctx, toCreate...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	addIDs := dt.Map(newTags, func(t *tag_ref.Tag) tag_ref.ID { return t.ID })
	removeIDs := dt.Reduce(n.Tags, func(acc []tag_ref.ID, prev *tag_ref.Tag) []tag_ref.ID {
		if lo.Contains(toRemove, prev.Name) {
			acc = append(acc, prev.ID)
		}
		return acc
	}, []tag_ref.ID{})

	opts = append(opts, node_writer.WithTagsAdd(addIDs...))
	opts = append(opts, node_writer.WithTagsRemove(removeIDs...))

	return opts, nil
}
