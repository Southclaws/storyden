package url

import (
	"context"
	"net/http"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/content"
)

var errFailedToScrape = fault.New("failed to scrape")

type Scraper interface {
	Scrape(ctx context.Context, url string) (*WebContent, error)
}

type WebContent struct {
	Title       string
	Description string
	Text        string
	Image       string
	Content     content.Rich
}

func Build() fx.Option {
	return fx.Provide(New)
}

type webScraper struct{}

func New() Scraper {
	return &webScraper{}
}

func (s *webScraper) Scrape(ctx context.Context, addr string) (*WebContent, error) {
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

	if resp.StatusCode != http.StatusOK {
		return nil, fault.Wrap(errFailedToScrape, fctx.With(ctx))
	}

	return s.postprocess(ctx, addr, resp.Body)
}
