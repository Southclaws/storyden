package scrape

import (
	"context"
	"io"
	"net/http"
	"net/netip"
	"strings"
	"testing"

	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/storyden/app/resources/link/link_url"
	"github.com/Southclaws/storyden/internal/infrastructure/httpsafe"
	"github.com/stretchr/testify/require"
)

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestScrapeFollowsRedirectsAndResolvesMetadataAtDestination(t *testing.T) {
	t.Parallel()
	s := webScraper{client: &http.Client{Transport: transportFunc(func(req *http.Request) (*http.Response, error) {
		response := &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Request: req}
		body := `<meta property="og:title" content="Final title"><meta property="og:image" content="image.png">`
		switch req.URL.String() {
		case "https://www.example.com/start":
			response.StatusCode = http.StatusMovedPermanently
			response.Header.Set("Location", "https://example.com/intermediate")
		case "https://example.com/intermediate":
			response.StatusCode = http.StatusFound
			response.Header.Set("Location", "/articles/final")
		case "https://example.com/articles/final":
		default:
			t.Fatalf("unexpected request: %s", req.URL)
		}

		response.Body = io.NopCloser(strings.NewReader(body))
		return response, nil
	})}}
	u, err := link_url.Parse("https://www.example.com/start")
	require.NoError(t, err)
	wc, err := s.Scrape(context.Background(), u)
	require.NoError(t, err)
	require.Equal(t, "https://example.com/articles/final", wc.URL.String())
	require.Equal(t, "Final title", wc.Title)
	require.Equal(t, "https://example.com/articles/image.png", wc.Image)
}

func TestScrapeClassifiesBlockedDNSAsInvalidInput(t *testing.T) {
	t.Parallel()
	s := webScraper{client: httpsafe.NewClient(httpsafe.Config{
		Resolver: func(context.Context, string, string) ([]netip.Addr, error) {
			return []netip.Addr{netip.MustParseAddr("10.0.0.1")}, nil
		},
	})}
	u, err := link_url.Parse("https://public-looking.example.com/")
	require.NoError(t, err)
	_, err = s.Scrape(context.Background(), u)
	require.Error(t, err)
	require.Equal(t, ftag.InvalidArgument, ftag.Get(err))
}

func TestScrapeRejectsRedirectToBlockedDNS(t *testing.T) {
	t.Parallel()
	client := httpsafe.NewClient(httpsafe.Config{
		Resolver: func(context.Context, string, string) ([]netip.Addr, error) {
			return []netip.Addr{netip.MustParseAddr("169.254.169.254")}, nil
		},
	})
	guarded := client.Transport
	client.Transport = transportFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "public.example.com" {
			return &http.Response{
				StatusCode: http.StatusFound,
				Header:     http.Header{"Location": []string{"https://redirect.example.com/"}},
				Body:       io.NopCloser(strings.NewReader("")),
				Request:    req,
			}, nil
		}

		return guarded.RoundTrip(req)
	})
	s := webScraper{client: client}
	u, err := link_url.Parse("https://public.example.com/")
	require.NoError(t, err)
	_, err = s.Scrape(context.Background(), u)
	require.Error(t, err)
	require.Equal(t, ftag.InvalidArgument, ftag.Get(err))
}

func TestScrapeRejectsRedirectsOutsideLinkURLPolicyBeforeRequest(t *testing.T) {
	for _, target := range []string{
		"https://example.com:443/",
		"http://example.com:80/",
		"https://example.com:8443/",
		"https://8.8.8.8/",
		"https://[2606:4700:4700::1111]/",
		"https://intranet/",
		"https://bad_host.example.com/",
	} {
		t.Run(target, func(t *testing.T) {
			s := New().(*webScraper)
			requests := 0
			s.client.Transport = transportFunc(func(req *http.Request) (*http.Response, error) {
				requests++
				return &http.Response{
					StatusCode: http.StatusFound,
					Header:     http.Header{"Location": []string{target}},
					Body:       io.NopCloser(strings.NewReader("")),
					Request:    req,
				}, nil
			})
			u, err := link_url.Parse("https://example.com/start")
			require.NoError(t, err)
			_, err = s.Scrape(context.Background(), u)
			require.Error(t, err)
			require.Equal(t, ftag.InvalidArgument, ftag.Get(err))
			require.Equal(t, 1, requests, "only the original URL may be requested")
		})
	}
}
