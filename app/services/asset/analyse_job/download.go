package analyse_job

import (
	"context"
	"net/http"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/services/asset/asset_upload"
)

func (c *analyseConsumer) downloadAsset(ctx context.Context, src string, fillrule opt.Optional[asset.ContentFillCommand]) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, src, nil)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ctx = fctx.WithMeta(ctx, "status", resp.Status)
		return fault.Wrap(fault.New("failed to get asset"), fctx.With(ctx))
	}

	// TODO: Better naming???
	name := mark.Slugify(src)

	a, err := c.uploader.Upload(ctx, resp.Body, resp.ContentLength, asset.NewFilename(name), asset_upload.Options{})
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if fr, ok := fillrule.Get(); ok {
		targetNodeID, targetSet := fr.TargetNodeID.Get()
		if !targetSet {
			return fault.New("target node ID not set", fctx.With(ctx))
		}

		_, err = c.nodeWriter.Update(ctx, library.NewID(targetNodeID), node_writer.WithAssets([]asset.AssetID{a.ID}))
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	return nil
}
