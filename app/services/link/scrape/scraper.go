package scrape

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/link/link_url"
	"github.com/Southclaws/storyden/internal/infrastructure/httpsafe"
)

const scrapeFetchTimeout = 30 * time.Second

var errFailedToScrape = fault.New("failed to scrape")

type Scraper interface {
	Scrape(ctx context.Context, url link_url.URL) (*WebContent, error)
}

type WebContent struct {
	URL         link_url.URL
	Title       string
	Description string
	Text        string
	Favicon     string
	Image       string
	Content     datagraph.Content
}

type webScraper struct {
	client *http.Client
}

func New() (Scraper, error) {
	client, err := httpsafe.NewSurfClient(httpsafe.Config{Timeout: scrapeFetchTimeout})
	if err != nil {
		return nil, err
	}

	redirectPolicy := client.CheckRedirect
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if _, err := link_url.Parse(req.URL.String()); err != nil {
			return err
		}

		return redirectPolicy(req, via)
	}

	return &webScraper{client: client}, nil
}

func (s *webScraper) Scrape(ctx context.Context, addr link_url.URL) (*WebContent, error) {
	if addr.IsZero() {
		return nil, link_url.ErrEmpty
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, addr.String(), nil)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	resp, err := s.client.Do(req)
	if err != nil {
		if errors.Is(err, httpsafe.ErrDisallowedHost) || errors.Is(err, httpsafe.ErrDisallowedAddress) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument),
				fmsg.WithDesc("unsupported URL host", "Only public HTTP and HTTPS URLs are supported."))
		}

		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	defer resp.Body.Close()

	finalURL, err := link_url.Parse(resp.Request.URL.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return &WebContent{URL: finalURL}, fault.Wrap(errFailedToScrape, fctx.With(ctx))
	}

	wc, err := s.postprocess(ctx, finalURL, resp.Body)
	if err != nil {
		return &WebContent{URL: finalURL}, fault.Wrap(err, fctx.With(ctx))
	}

	wc.URL = finalURL
	return wc, nil
}
