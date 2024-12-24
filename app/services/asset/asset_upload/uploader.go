package asset_upload

import (
	"context"
	"io"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/asset/asset_writer"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/mime"
)

type Uploader struct {
	l *zap.Logger

	nodewriter *node_writer.Writer
	assets     *asset_writer.Writer
	objects    object.Storer
	queue      pubsub.Topic[mq.AnalyseAsset]
}

func New(
	l *zap.Logger,

	nodewriter *node_writer.Writer,
	assets *asset_writer.Writer,
	objects object.Storer,
	queue pubsub.Topic[mq.AnalyseAsset],
) *Uploader {
	return &Uploader{
		l: l.With(zap.String("service", "asset")),

		nodewriter: nodewriter,
		assets:     assets,
		objects:    objects,
		queue:      queue,
	}
}

type Options struct {
	ContentFill opt.Optional[asset.ContentFillCommand]
	ParentID    opt.Optional[asset.AssetID]
}

func (s *Uploader) Upload(ctx context.Context, or io.Reader, size int64, name asset.Filename, opts Options) (*asset.Asset, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	mt, r, err := mime.Detect(or)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	a, err := func() (asset *asset.Asset, err error) {
		if pid, ok := opts.ParentID.Get(); ok {
			return s.assets.AddVersion(ctx, xid.ID(accountID), name, int(size), *mt, pid)
		} else {
			return s.assets.Add(ctx, xid.ID(accountID), name, int(size), *mt)
		}
	}()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if cfr, ok := opts.ContentFill.Get(); ok {
		targetNode, ok := cfr.TargetNodeID.Get()
		if !ok {
			return nil, fault.New("target node ID not set", fctx.With(ctx))
		}

		nodeID := library.QueryKey{mark.NewQueryKeyID(targetNode)}

		_, err := s.nodewriter.Update(ctx, nodeID, node_writer.WithAssets([]asset.AssetID{a.ID}))
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	path := asset.BuildAssetPath(a.Name)

	if err := s.objects.Write(ctx, path, r, size); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.queue.Publish(ctx, mq.AnalyseAsset{
		AssetID:         a.ID,
		ContentFillRule: opts.ContentFill,
	}); err != nil {
		s.l.Error("failed to publish analyse asset message", zap.Error(err))
	}

	return a, nil
}
