package url

import (
	"context"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"go.uber.org/fx"
	"golang.org/x/net/html"
)

type Scraper interface {
	Scrape(ctx context.Context, url string) (*WebContent, error)
}

type WebContent struct {
	Title       string
	Description string
	Image       string
}

func Build() fx.Option {
	return fx.Provide(New)
}

type scraper struct{}

func New() Scraper {
	return &scraper{}
}

func (s *scraper) Scrape(ctx context.Context, url string) (*WebContent, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	t := metatable(doc)

	wc := &WebContent{
		Title:       t["og:title"],
		Description: t["og:description"],
		Image:       t["og:image"],
	}

	return wc, nil
}

func metatable(doc *goquery.Document) map[string]string {
	return dt.Reduce(doc.Find("head > meta").Nodes, func(wc map[string]string, n *html.Node) map[string]string {
		k, v := ogtable(n.Attr)
		if k != "" && v != "" {
			wc[k] = v
		}

		return wc
	}, map[string]string{})
}

func ogtable(attrs []html.Attribute) (k string, v string) {
	for _, a := range attrs {
		switch a.Key {
		case "property":
			k = a.Val
		case "name":
			k = a.Val
		case "content":
			v = a.Val
		}
	}

	return
}
