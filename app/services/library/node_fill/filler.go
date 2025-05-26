package node_fill

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/gosimple/slug"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/link/scrape"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

var errFillRuleNotAvailale = fault.New("fill rule not available")

type Filler struct {
	nodeWriter *node_writer.Writer
	indexQueue pubsub.Topic[mq.IndexNode]
	assetQueue pubsub.Topic[mq.DownloadAsset]
}

func New(
	nodeWriter *node_writer.Writer,
	indexQueue pubsub.Topic[mq.IndexNode],
	assetQueue pubsub.Topic[mq.DownloadAsset],
) *Filler {
	return &Filler{
		nodeWriter: nodeWriter,
		indexQueue: indexQueue,
		assetQueue: assetQueue,
	}
}

func (f *Filler) FillContentFromLink(ctx context.Context, link *link_ref.LinkRef, wc *scrape.WebContent, cfr asset.ContentFillCommand) error {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	vis := visibility.VisibilityPublished

	switch cfr.FillRule {
	case asset.ContentFillRuleCreate:
		title := wc.Title
		if title == "" {
			title = "Untitled"
		}
		slug, _ := mark.NewSlug(slug.Make(title + "-" + xid.New().String()))

		opts := []node_writer.Option{
			node_writer.WithLink(link.ID),
			node_writer.WithContent(wc.Content),
			node_writer.WithDescription(wc.Description),
			node_writer.WithVisibility(vis),
		}

		if v, ok := cfr.TargetNodeID.Get(); ok {
			opts = append(opts, node_writer.WithParent(library.NodeID(v)))
		}

		if v, ok := link.PrimaryImage.Get(); ok {
			opts = append(opts, node_writer.WithPrimaryImage(v.ID))
		}

		n, err := f.nodeWriter.Create(ctx, accountID, wc.Title, *slug, opts...)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		err = f.assetQueue.Publish(ctx, dt.Map(wc.Content.Media(), func(u string) mq.DownloadAsset {
			return mq.DownloadAsset{
				URL: u,
				ContentFillRule: opt.New(asset.ContentFillCommand{
					TargetNodeID: opt.New(n.GetID()),
					FillRule:     asset.ContentFillRuleAppend,
				}),
			}
		})...)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		// TODO: Decide what to do here, code is still useful.
		// err = f.autoFillQueue.Publish(ctx, mq.AutoFillNode{
		// 	ID:      library.NodeID(n.Mark.ID()),
		// 	AutoTag: true,
		// })
		// if err != nil {
		// 	return fault.Wrap(err, fctx.With(ctx))
		// }

		if vis == visibility.VisibilityPublished {
			if err := f.indexQueue.Publish(ctx, mq.IndexNode{
				ID: library.NodeID(n.Mark.ID()),
			}); err != nil {
				return fault.Wrap(err, fctx.With(ctx))
			}
		}

	case asset.ContentFillRuleReplace:
		// TODO: Decide what to do here, code is still useful.

		// targetNode, ok := cfr.TargetNodeID.Get()
		// if !ok {
		// 	return fault.New("target node ID not set", fctx.With(ctx))
		// }

		// err = f.autoFillQueue.Publish(ctx, mq.AutoFillNode{
		// 	ID:      library.NodeID(library.NodeID(targetNode)),
		// 	AutoTag: true,
		// })
		// if err != nil {
		// 	return fault.Wrap(err, fctx.With(ctx))
		// }

	default:
		return fault.Wrap(errFillRuleNotAvailale, fctx.With(ctx))
	}
	return nil
}
