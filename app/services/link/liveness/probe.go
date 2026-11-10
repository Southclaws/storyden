package liveness

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/link/link_check"
	"github.com/Southclaws/storyden/app/resources/link/link_url"
	"github.com/Southclaws/storyden/internal/infrastructure/httpsafe"
)

type Prober interface {
	Check(context.Context, link_url.URL) (link_check.Status, error)
}

type httpProber struct {
	client *http.Client
}

func NewProber() Prober {
	client := httpsafe.NewClient(httpsafe.Config{Timeout: 10 * time.Second})
	redirectPolicy := client.CheckRedirect
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if _, err := link_url.Parse(req.URL.String()); err != nil {
			return err
		}

		return redirectPolicy(req, via)
	}

	return &httpProber{client: client}
}

func (p *httpProber) Check(ctx context.Context, u link_url.URL) (link_check.Status, error) {
	if u.IsZero() {
		return link_check.Unknown, link_url.ErrEmpty
	}

	status, err := p.request(ctx, http.MethodHead, u)
	if err != nil {
		return link_check.Unknown, err
	}

	if status == http.StatusMethodNotAllowed || status == http.StatusNotImplemented {
		status, err = p.request(ctx, http.MethodGet, u)
		if err != nil {
			return link_check.Unknown, err
		}
	}

	if status >= 200 && status < 300 {
		return link_check.Alive, nil
	}

	if status == http.StatusUnauthorized || status == http.StatusForbidden || status == http.StatusTooManyRequests {
		return link_check.Unknown, nil
	}

	return link_check.Offline, nil
}

func (p *httpProber) request(ctx context.Context, method string, u link_url.URL) (int, error) {
	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("User-Agent", "Storyden Link Check")
	if method == http.MethodGet {
		req.Header.Set("Range", "bytes=0-0")
	}

	resp, err := p.client.Do(req)
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}

	if errors.Is(err, httpsafe.ErrDisallowedHost) || errors.Is(err, httpsafe.ErrDisallowedAddress) {
		return 0, fault.Wrap(err, ftag.With(ftag.InvalidArgument))
	}

	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}
