package node_mutate

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/library/node_auth"
)

type Option func(*updateOptions)

type updateOptions struct {
	tagFillRule opt.Optional[tag.TagFillRule]
}

func WithTagFillRule(fr tag.TagFillRule) Option {
	return func(uo *updateOptions) {
		uo.tagFillRule = opt.New(fr)
	}
}

type Updated struct {
	library.Node
	TagSuggestions opt.Optional[tag_ref.Names]
}

func (s *Manager) Update(ctx context.Context, qk library.QueryKey, p Partial, options ...Option) (*Updated, error) {
	updateOpts := updateOptions{}
	for _, fn := range options {
		fn(&updateOpts)
	}

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

	opts, err := s.applyOpts(ctx, p)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// TODO: Queue this for background processing
	if v, ok := p.AssetSources.Get(); ok {
		for _, source := range v {
			a, err := s.fetcher.CopyAsset(ctx, source)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			opts = append(opts, node_writer.WithAssets([]asset.AssetID{a.ID}))
		}
	}

	assetsAdd, assetsAddSet := p.AssetsAdd.Get()
	if assetsAddSet && p.ContentFill.Ok() {

		messages := dt.Map(assetsAdd, func(a asset.AssetID) mq.AnalyseAsset {
			return mq.AnalyseAsset{
				AssetID:         a,
				ContentFillRule: p.ContentFill,
			}
		})

		if err := s.assetAnalyseQueue.Publish(ctx, messages...); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	if u, ok := p.URL.Get(); ok {
		ln, err := s.fetcher.Fetch(ctx, u)
		if err == nil {
			opts = append(opts, node_writer.WithLink(xid.ID(ln.ID)))
		}
	}

	suggestedTags := opt.NewEmpty[tag_ref.Names]()
	if tfr, ok := updateOpts.tagFillRule.Get(); ok {
		// If the update query contains new content, use that, otherwise, fall
		// back to the current content in the node.
		content := p.Content.Or(n.Content.OrZero())

		// Only bother if there's any actual content to work with!
		if !content.IsEmpty() {
			gathered, err := s.tagger.Gather(ctx, tfr, content)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			switch tfr {
			case tag.TagFillRuleQuery:
				suggestedTags = opt.New(gathered)

			case tag.TagFillRuleReplace:
				if t, ok := p.Tags.Get(); ok {
					p.Tags = opt.New(append(t, gathered...))
				} else {
					p.Tags = opt.New(gathered)
				}
			default:
			}
		}
	}

	if tags, ok := p.Tags.Get(); ok {
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
	}

	n, err = s.nodeWriter.Update(ctx, qk, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if n.Visibility == visibility.VisibilityPublished {
		if err := s.indexQueue.Publish(ctx, mq.IndexNode{
			ID: library.NodeID(n.Mark.ID()),
		}); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	} else {
		if err := s.deleteQueue.Publish(ctx, mq.DeleteNode{
			ID: library.NodeID(n.GetID()),
		}); err != nil {
			s.logger.Error("failed to publish index post message", zap.Error(err))
		}
	}

	s.fetcher.HydrateContentURLs(ctx, n)

	u := Updated{
		Node:           *n,
		TagSuggestions: suggestedTags,
	}

	return &u, nil
}
