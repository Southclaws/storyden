package asset_download

import (
	"context"
	"io"
	"log/slog"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/asset/asset_image"
	"github.com/Southclaws/storyden/app/resources/asset/asset_querier"
	"github.com/Southclaws/storyden/app/services/background_job/job_dispatcher"
	"github.com/Southclaws/storyden/internal/infrastructure/object"
)

type Downloader struct {
	jobs       *asset_image.Repository
	logger     *slog.Logger
	assets     *asset_querier.Querier
	objects    object.Storer
	dispatcher job_dispatcher.Dispatcher
}

func New(
	assets *asset_querier.Querier,
	objects object.Storer,
	jobs *asset_image.Repository,
	dispatcher job_dispatcher.Dispatcher,
	logger *slog.Logger,
) *Downloader {
	return &Downloader{
		assets:     assets,
		jobs:       jobs,
		logger:     logger,
		objects:    objects,
		dispatcher: dispatcher,
	}
}

func (d *Downloader) Get(ctx context.Context, id asset.Filename) (*asset.Asset, io.Reader, error) {
	a, err := d.assets.Get(ctx, id)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	path := asset.BuildAssetPath(a.Name)
	ctx = fctx.WithMeta(ctx, "path", path, "asset_id", id.String())

	r, size, err := d.objects.Read(ctx, path)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	a.Size = int(size)

	return a, r, nil
}

func (d *Downloader) GetByID(ctx context.Context, id asset.AssetID) (*asset.Asset, io.Reader, error) {
	a, err := d.assets.GetByID(ctx, id)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	path := asset.BuildAssetPath(a.Name)
	ctx = fctx.WithMeta(ctx, "path", path, "asset_id", id.String())

	r, size, err := d.objects.Read(ctx, path)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	a.Size = int(size)

	return a, r, nil
}

func (d *Downloader) Variant(
	ctx context.Context,
	name asset.Filename,
	variant asset.Variant,
) (*asset.Asset, io.Reader, bool, error) {
	source, err := d.assets.Get(ctx, name)
	if err != nil {
		return nil, nil, false, err
	}

	return d.VariantByAsset(ctx, source, variant)
}

func (d *Downloader) VariantByID(
	ctx context.Context,
	id asset.AssetID,
	variant asset.Variant,
) (*asset.Asset, io.Reader, bool, error) {
	source, err := d.assets.GetByID(ctx, id)
	if err != nil {
		return nil, nil, false, err
	}

	return d.VariantByAsset(ctx, source, variant)
}

func (d *Downloader) VariantByAsset(
	ctx context.Context,
	source *asset.Asset,
	variant asset.Variant,
) (*asset.Asset, io.Reader, bool, error) {
	selected := source
	pending := false
	preset, transform := asset.Preset(variant)
	canTransform := transform && asset.ImageEligible(source.MIME.String())
	if canTransform {
		metadata, err := source.ImageMetadata()
		if err != nil {
			return nil, nil, false, err
		}
		canTransform = metadata.Preset == ""
	}

	if canTransform {
		variants, err := d.jobs.Variants(ctx, source.ID)
		if err != nil {
			return nil, nil, false, err
		}
		selected, err = selectVariant(source, variants, preset)
		if err != nil {
			return nil, nil, false, err
		}
		selectedMetadata, err := selected.ImageMetadata()
		if err != nil {
			return nil, nil, false, err
		}
		pending = selectedMetadata.Preset != variant
		if pending {
			if err = d.jobs.Request(ctx, source.ID, preset.Bit|asset.EagerImagePresets); err != nil {
				d.logger.Warn("image request failed", "asset_id", source.ID, "error", err)
			} else {
				d.dispatcher.Dispatch(ctx, source.ID)
			}
		}
	}

	r, size, err := d.objects.Read(ctx, asset.BuildAssetPath(selected.Name))
	if err != nil {
		return nil, nil, false, err
	}
	selected.Size = int(size)

	return selected, r, pending, nil
}

func selectVariant(source *asset.Asset, variants []*asset.Asset, target asset.ImagePreset) (*asset.Asset, error) {
	var best *asset.Asset
	bestEdge := 0
	for _, candidate := range variants {
		metadata, err := candidate.ImageMetadata()
		if err != nil {
			return nil, err
		}
		if metadata == nil {
			continue
		}

		preset, ok := asset.Preset(metadata.Preset)
		if !ok {
			continue
		}
		if preset.Name == target.Name {
			return candidate, nil
		}
		if preset.Square != target.Square || preset.PNG != target.PNG {
			continue
		}

		edge := preset.Width
		if !target.Square && edge < target.Width {
			continue
		}

		if best != nil {
			if bestEdge >= target.Width {
				if edge < target.Width || edge >= bestEdge {
					continue
				}
			} else if edge < target.Width && edge <= bestEdge {
				continue
			}
		}

		best = candidate
		bestEdge = edge
	}

	if best != nil {
		return best, nil
	}

	return source, nil
}
