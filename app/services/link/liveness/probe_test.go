package liveness

import (
	"context"
	"io"
	"net/http"
	"net/netip"
	"strings"
	"testing"

	"github.com/Southclaws/fault/ftag"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/link/link_check"
	"github.com/Southclaws/storyden/app/resources/link/link_url"
	"github.com/Southclaws/storyden/internal/infrastructure/httpsafe"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type trackedBody struct {
	read   bool
	closed bool
}

func (b *trackedBody) Read([]byte) (int, error) {
	b.read = true
	return 0, io.EOF
}

func (b *trackedBody) Close() error {
	b.closed = true
	return nil
}

func TestProbeDoesNotReadPageBodies(t *testing.T) {
	for _, scenario := range []struct {
		name string
		head int
		get  int
		want link_check.Status
	}{
		{name: "success", head: 200, want: link_check.Alive},
		{name: "no_content", head: 204, want: link_check.Alive},
		{name: "not_found", head: 404, want: link_check.Offline},
		{name: "gone", head: 410, want: link_check.Offline},
		{name: "server_failure", head: 503, want: link_check.Offline},
		{name: "authentication", head: 401, want: link_check.Unknown},
		{name: "forbidden", head: 403, want: link_check.Unknown},
		{name: "rate_limit", head: 429, want: link_check.Unknown},
		{name: "head_not_allowed", head: 405, get: 206, want: link_check.Alive},
		{name: "head_not_implemented", head: 501, get: 200, want: link_check.Alive},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			var methods []string
			var bodies []*trackedBody
			p := NewProber().(*httpProber)
			p.client.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
				methods = append(methods, req.Method)
				status := scenario.head
				if req.Method == http.MethodGet {
					require.Equal(t, "bytes=0-0", req.Header.Get("Range"))
					status = scenario.get
				}

				body := &trackedBody{}
				bodies = append(bodies, body)
				return &http.Response{StatusCode: status, Header: make(http.Header), Body: body, Request: req}, nil
			})
			u, err := link_url.Parse("https://example.com/page")
			require.NoError(t, err)
			status, err := p.Check(context.Background(), u)
			require.NoError(t, err)
			require.Equal(t, scenario.want, status)
			if scenario.get != 0 {
				require.Equal(t, []string{http.MethodHead, http.MethodGet}, methods)
			} else {
				require.Equal(t, []string{http.MethodHead}, methods)
			}

			for _, body := range bodies {
				require.False(t, body.read)
				require.True(t, body.closed)
			}
		})
	}
}

func TestProbeValidatesEveryRedirect(t *testing.T) {
	for _, target := range []string{"https://example.com:443/", "http://127.0.0.1/", "ftp://example.com/", "https://intranet/"} {
		t.Run(target, func(t *testing.T) {
			requests := 0
			p := NewProber().(*httpProber)
			p.client.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
				requests++
				return &http.Response{
					StatusCode: http.StatusFound,
					Header:     http.Header{"Location": []string{target}},
					Body:       io.NopCloser(strings.NewReader("")),
					Request:    req,
				}, nil
			})
			u, err := link_url.Parse("https://example.com/")
			require.NoError(t, err)
			_, err = p.Check(context.Background(), u)
			require.Error(t, err)
			require.Equal(t, ftag.InvalidArgument, ftag.Get(err))
			require.Equal(t, 1, requests)
		})
	}
}

func TestProbeRejectsBlockedDNS(t *testing.T) {
	p := &httpProber{client: httpsafe.NewClient(httpsafe.Config{
		Resolver: func(context.Context, string, string) ([]netip.Addr, error) {
			return []netip.Addr{netip.MustParseAddr("10.0.0.1")}, nil
		},
	})}
	u, err := link_url.Parse("https://example.com/")
	require.NoError(t, err)
	_, err = p.Check(context.Background(), u)
	require.Error(t, err)
	require.Equal(t, ftag.InvalidArgument, ftag.Get(err))
}
