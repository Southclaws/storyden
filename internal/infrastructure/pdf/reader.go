package pdf

import (
	"context"
	_ "embed"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/webassembly"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Extractor struct {
	l        *zap.Logger
	pool     pdfium.Pool
	instance pdfium.Pdfium
}

func New(lc fx.Lifecycle, l *zap.Logger) (*Extractor, error) {
	pool, err := webassembly.Init(webassembly.Config{
		MinIdle:  1,
		MaxIdle:  1,
		MaxTotal: 1,
	})
	if err != nil {
		return nil, err
	}

	instance, err := pool.GetInstance(time.Second * 30)
	if err != nil {
		return nil, err
	}

	return &Extractor{
		l:        l.With(zap.String("package", "pdf")),
		pool:     pool,
		instance: instance,
	}, nil
}

type ExtractionResult struct {
	// TODO: split by pages etc.
	Text string
}

func (e *Extractor) Extract(ctx context.Context, buf []byte) (*ExtractionResult, error) {
	doc, err := e.instance.OpenDocument(&requests.OpenDocument{
		File: &buf,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	defer func() {
		_, err := e.instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
			Document: doc.Document,
		})
		if err != nil {
			e.l.Error("failed to close document", zap.Error(err))
		}
	}()

	pages, err := e.instance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
		Document: doc.Document,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	blocks := []string{}

	for p := range pages.PageCount {
		pageText, err := e.instance.GetPageText(&requests.GetPageText{
			Page: requests.Page{
				ByIndex: &requests.PageByIndex{
					Document: doc.Document,
					Index:    p,
				},
			},
		})
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, pageText.Text)
	}

	return &ExtractionResult{
		Text: strings.Join(blocks, "\n"),
	}, nil
}
