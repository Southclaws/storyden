package node_semdex

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
)

func (i *semdexer) index(ctx context.Context, id library.NodeID, summarise bool, autotag bool) error {
	qk := library.NewID(xid.ID(id))

	node, err := i.nodeQuerier.Get(ctx, qk)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	err = i.indexer.Index(ctx, node)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	opts := []node_writer.Option{
		node_writer.WithIndexed(),
	}

	content := node.GetContent()

	if summarise {
		summarisedContent, err := i.getSummary(ctx, node)
		if err != nil {
			i.logger.Warn("failed to summarise node", zap.Error(err), zap.String("node_id", node.GetID().String()))
		} else {
			opts = append(opts, node_writer.WithContent(*summarisedContent))
		}

		content = *summarisedContent
	}

	if autotag {
		tagOpts, err := i.generateTags(ctx, node, content)
		if err != nil {
			i.logger.Warn("failed to autotag node", zap.Error(err), zap.String("node_id", node.GetID().String()))
		}

		opts = append(opts, tagOpts...)
	}

	_, err = i.nodeWriter.Update(ctx, qk, opts...)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (i *semdexer) getSummary(ctx context.Context, p datagraph.Item) (*datagraph.Content, error) {
	summary, err := i.summariser.Summarise(ctx, p)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	content, err := datagraph.NewRichText(summary)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &content, nil
}

func (i *semdexer) deindex(ctx context.Context, id library.NodeID) error {
	qk := library.NewID(xid.ID(id))

	err := i.deleter.Delete(ctx, xid.ID(id))
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	_, err = i.nodeWriter.Update(ctx, qk, node_writer.WithIndexed())
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (i *semdexer) generateTags(ctx context.Context, n *library.Node, content datagraph.Content) ([]node_writer.Option, error) {
	gathered, err := i.tagger.Gather(ctx, tag.TagFillRuleReplace, content)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	tagOpts, err := i.applyTags(ctx, n, gathered)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return tagOpts, nil
}

func (i *semdexer) applyTags(ctx context.Context, n *library.Node, tags []tag_ref.Name) (opts []node_writer.Option, err error) {
	currentTagNames := n.Tags.Names()

	toCreate, toRemove := lo.Difference(tags, currentTagNames)

	newTags, err := i.tagWriter.Add(ctx, toCreate...)
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
