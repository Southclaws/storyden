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

func (s *scraper) Scrape(ctx context.Context, addr string) (*WebContent, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, addr, nil)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("Accept-Encoding", "*")
	req.Header.Add("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Dnt", "1")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Sec-Ch-Ua", `"Google Chrome";v="119", "Chromium";v="119", "Not?A_Brand";v="24"`)
	req.Header.Add("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Add("Sec-Ch-Ua-Platform", `"macOS"`)
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "none")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36`)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	t := metatable(doc)

	wc := &WebContent{
		Title:       title(t),
		Description: description(t),
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

func title(t map[string]string) string {
	if t["og:title"] != "" {
		return t["og:title"]
	}
	if t["title"] != "" {
		return t["title"]
	}
	if t["og:site_name"] != "" {
		return t["og:site_name"]
	}
	if t["og:url"] != "" {
		return t["og:url"]
	}
	if t["title"] != "" {
		return t["title"]
	}

	return ""
}

func description(t map[string]string) string {
	if t["og:description"] != "" {
		return t["og:description"]
	}
	if t["description"] != "" {
		return t["description"]
	}

	return ""
}
