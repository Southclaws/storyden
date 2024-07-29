package asset_upload

import (
	"context"
	"io"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Uploader struct {
	l    *zap.Logger
	rbac rbac.AccessManager

	assets  asset.Repository
	objects object.Storer
	queue   pubsub.Topic[mq.AnalyseAsset]
}

func New(
	l *zap.Logger,
	rbac rbac.AccessManager,

	assets asset.Repository,
	objects object.Storer,
	queue pubsub.Topic[mq.AnalyseAsset],
) *Uploader {
	return &Uploader{
		l:    l.With(zap.String("service", "asset")),
		rbac: rbac,

		assets:  assets,
		objects: objects,
		queue:   queue,
	}
}

type Options struct {
	ContentFill opt.Optional[asset.ContentFillCommand]
}

func (s *Uploader) Upload(ctx context.Context, r io.Reader, size int64, name asset.Filename, url string, opts Options) (*asset.Asset, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	a, err := s.assets.Add(ctx, accountID, name, url)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// TODO Use the content fill rule to also link the asset to the target node.

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
