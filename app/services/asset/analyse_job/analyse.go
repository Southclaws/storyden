package analyse_job

import (
	"context"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/services/asset/analyse"
	"github.com/Southclaws/storyden/app/services/asset/asset_upload"
)

type analyseConsumer struct {
	analyser   *analyse.Analyser
	uploader   *asset_upload.Uploader
	nodeWriter *node_writer.Writer
}

func newAnalyseConsumer(
	analyser *analyse.Analyser,
	uploader *asset_upload.Uploader,
	nodeWriter *node_writer.Writer,
) *analyseConsumer {
	return &analyseConsumer{
		analyser:   analyser,
		uploader:   uploader,
		nodeWriter: nodeWriter,
	}
}

func (i *analyseConsumer) analyseAsset(ctx context.Context, id asset.AssetID, fillrule opt.Optional[asset.ContentFillCommand]) error {
	return i.analyser.Analyse(ctx, id, fillrule)
}
