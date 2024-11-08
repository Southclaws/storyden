package analyse

import (
	"context"
	"io"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/asset/asset_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/services/asset/asset_download"
	"github.com/Southclaws/storyden/app/services/library/node_mutate"
	"github.com/Southclaws/storyden/internal/infrastructure/pdf"
)

type Analyser struct {
	assetQuerier *asset_querier.Querier
	downloader   *asset_download.Downloader
	nodereader   *node_querier.Querier
	nodewriter   *node_mutate.Manager
	pdfextractor *pdf.Extractor
}

func New(
	assetQuerier *asset_querier.Querier,
	downloader *asset_download.Downloader,
	nodereader *node_querier.Querier,
	nodewriter *node_mutate.Manager,
	pdfextractor *pdf.Extractor,
) *Analyser {
	return &Analyser{
		assetQuerier: assetQuerier,
		downloader:   downloader,
		nodereader:   nodereader,
		nodewriter:   nodewriter,
		pdfextractor: pdfextractor,
	}
}

func (a *Analyser) Analyse(ctx context.Context, id asset.AssetID, fillrule opt.Optional[asset.ContentFillCommand]) error {
	ast, err := a.assetQuerier.GetByID(ctx, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	switch ast.MIME.String() {
	case "application/pdf":
		_, reader, err := a.downloader.GetByID(ctx, id)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		buf, err := io.ReadAll(reader)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		return a.analysePDF(ctx, buf, fillrule)
	}

	return nil
}
