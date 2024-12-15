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
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/link/fetcher"
)

type preMutationResult struct {
	opts []node_writer.Option

	// Ideally, this API should only return node writer options, but because of
	// a weird public API design choice I made, the PATCH /nodes endpoint also
	// returns tag suggestions which can be opted out of being applied directly.
	// This may change in future but it would require breaking public API change
	// and it works pretty well at the moment as an API design, so not critical.
	tags opt.Optional[tag_ref.Names]
}

// preMutation constructs node_writer options for a create or partial update.
func (s *Manager) preMutation(ctx context.Context, p Partial, current opt.Optional[library.Node]) (*preMutationResult, error) {
	opts := []node_writer.Option{}

	// Apply all primitive options. These are just basic partial updates.
	p.Name.Call(func(value string) { opts = append(opts, node_writer.WithName(value)) })
	p.Slug.Call(func(value mark.Slug) { opts = append(opts, node_writer.WithSlug(value.String())) })
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

	// If assets have been added to the node and there's a content fill rule,
	// queue the assets for extraction for the node's content etc.
	assetsAdd, assetsAddSet := p.AssetsAdd.Get()
	if assetsAddSet && p.ContentFill.Ok() {
		if err := s.assetAnalyseQueue.Publish(ctx, dt.Map(assetsAdd, func(a asset.AssetID) mq.AnalyseAsset {
			return mq.AnalyseAsset{
				AssetID:         a,
				ContentFillRule: p.ContentFill,
			}
		})...); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	// If there's a URL being applied, fetch its contents.
	if u, ok := p.URL.Get(); ok {
		ln, err := s.fetcher.Fetch(ctx, u, fetcher.Options{})
		if err == nil {
			opts = append(opts, node_writer.WithLink(xid.ID(ln.ID)))
		}
	}

	// The content to use during pre-mutation tasks such as tag suggestion, auto
	// title generation and content summarisation. If it's a new node, this will
	// be the content submitted for the new node, if it's an update, either pick
	// the new content if specified in the partial, or the current node content.
	content := p.Content.Or(current.OrZero().Content.OrZero())

	var tags opt.Optional[tag_ref.Names]
	if tfr, ok := p.TagFill.Get(); ok {
		suggested, err := s.buildTagSuggestionOpts(ctx, content, tfr.FillRule)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if tfr.FillRule == tag.TagFillRuleReplace {
			tags = opt.New(suggested)
		} else {
			tags = p.Tags
		}
	} else {
		tags = p.Tags
	}

	if t, ok := tags.Get(); ok {
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

	if p.ContentSummarise.OrZero() {
		opt, err := s.buildSummaryOpts(ctx, content)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		opts = append(opts, opt)
	}

	return &preMutationResult{
		opts: opts,
		tags: tags,
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

func (s *Manager) buildTagSuggestionOpts(ctx context.Context, content datagraph.Content, tfr tag.TagFillRule) (tag_ref.Names, error) {
	// Only bother if there's any actual content to work with!
	if content.IsEmpty() {
		return nil, nil
	}

	gathered, err := s.tagger.Gather(ctx, tfr, content)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return gathered, nil
}

func (s *Manager) buildSummaryOpts(ctx context.Context, content datagraph.Content) (node_writer.Option, error) {
	summary, err := s.summariser.Summarise(ctx, content)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	newContent, err := datagraph.NewRichText(summary)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return node_writer.WithContent(newContent), nil
}
