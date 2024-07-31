package analyse

import (
	"context"
	"io"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/services/asset/asset_download"
	"github.com/Southclaws/storyden/app/services/library/node_mutate"
	"github.com/Southclaws/storyden/internal/infrastructure/pdf"
	"github.com/gabriel-vasile/mimetype"
)

type Analyser struct {
	downloader   *asset_download.Downloader
	nodereader   library.Repository
	nodewriter   node_mutate.Manager
	pdfextractor *pdf.Extractor
}

func New(
	downloader *asset_download.Downloader,
	nodereader library.Repository,
	nodewriter node_mutate.Manager,
	pdfextractor *pdf.Extractor,
) *Analyser {
	return &Analyser{
		downloader:   downloader,
		nodereader:   nodereader,
		nodewriter:   nodewriter,
		pdfextractor: pdfextractor,
	}
}

func (a *Analyser) Analyse(ctx context.Context, id asset.AssetID, fillrule opt.Optional[asset.ContentFillCommand]) error {
	_, reader, err := a.downloader.GetByID(ctx, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	// TODO: Store the size claimed by the client in the asset table.
	// We need this to prevent DOS attacks with spoofed Content-Size headers.
	// buf := make([]byte, metadata.Size)
	// _, err = io.ReadFull(reader, buf)
	// if err != nil {
	// 	return fault.Wrap(err, fctx.With(ctx))
	// }

	buf, err := io.ReadAll(reader)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	mt := mimetype.Detect(buf)
	mime := mt.String()

	// switch mime { case "application/pdf": EXTRACT MUH PDF TEXT BOI }

	switch mime {
	case "application/pdf":
		err = a.analysePDF(ctx, buf, fillrule)
	}
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
