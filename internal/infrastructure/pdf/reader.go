package pdf

import (
	"context"
	_ "embed"
	"strings"
	"sync"
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

	initlock  sync.Mutex
	firstUse  bool
	readyChan chan bool
}

func New(lc fx.Lifecycle, l *zap.Logger) (*Extractor, error) {
	e := Extractor{
		l:         l.With(zap.String("package", "pdf")),
		readyChan: make(chan bool, 1),
	}

	lc.Append(fx.StartHook(func(ctx context.Context) error {
		init := func() error {
			start := time.Now()

			pool, err := webassembly.Init(webassembly.Config{
				MinIdle:  1,
				MaxIdle:  1,
				MaxTotal: 1,
			})
			if err != nil {
				return err
			}

			instance, err := pool.GetInstance(time.Second * 30)
			if err != nil {
				return err
			}

			e.initlock.Lock()
			defer e.initlock.Unlock()

			e.pool = pool
			e.instance = instance

			l.Info("pdf worker pool initialised", zap.Duration("time_taken", time.Since(start)))

			e.readyChan <- true

			return nil
		}

		go func() {
			if err := init(); err != nil {
				l.Fatal("failed to initialize PDFium worker pool", zap.Error(err))
			}
		}()

		return nil
	}))

	return &e, nil
}

func (e *Extractor) waitForReady() error {
	if e.firstUse {
		return nil
	}

	e.initlock.Lock()
	defer e.initlock.Unlock()

	select {
	case <-e.readyChan:
		e.firstUse = true

	case <-time.After(time.Second * 30):
		if e.firstUse {
			return nil
		}

		return fault.New("PDFium worker pool did not initialize within 30 seconds")
	}

	return nil
}

type ExtractionResult struct {
	Text string
	HTML *html.Node
}

func (e *Extractor) Extract(ctx context.Context, buf []byte) (*ExtractionResult, error) {
	if err := e.waitForReady(); err != nil {
		return nil, err
	}

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
