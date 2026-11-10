package httpsafe

import (
	"compress/gzip"
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSurfClientRejectsUnsafeDestinations(t *testing.T) {
	t.Parallel()

	for _, scheme := range []string{"http", "https"} {
		for _, host := range []string{"localhost", "127.0.0.1", "169.254.169.254", "[::1]", "private.example.com"} {
			t.Run(scheme+"/"+host, func(t *testing.T) {
				client, err := NewSurfClient(Config{
					Timeout: time.Second,
					Resolver: func(context.Context, string, string) ([]netip.Addr, error) {
						return []netip.Addr{
							netip.MustParseAddr("93.184.216.34"),
							netip.MustParseAddr("10.0.0.1"),
						}, nil
					},
				})
				require.NoError(t, err)
				t.Cleanup(client.CloseIdleConnections)

				_, err = client.Get(scheme + "://" + host)
				if host == "private.example.com" {
					require.ErrorIs(t, err, ErrDisallowedAddress)
				} else {
					require.ErrorIs(t, err, ErrDisallowedHost)
				}
			})
		}
	}
}

func TestSurfClientRedirectsAndDecompresses(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/start":
			http.Redirect(w, r, "/final", http.StatusFound)
		case "/final":
			w.Header().Set("Content-Encoding", "gzip")
			writer := gzip.NewWriter(w)
			defer writer.Close()
			io.WriteString(writer, "page content")
		case "/private":
			http.Redirect(w, r, "https://private.example.com/", http.StatusFound)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := surfTestClient(t, server)
	baseURL := "http://public.example.com"

	response, err := client.Get(baseURL + "/start")
	require.NoError(t, err)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, "page content", string(body))
	assert.Equal(t, baseURL+"/final", response.Request.URL.String())

	_, err = client.Get(baseURL + "/private")
	require.ErrorIs(t, err, ErrDisallowedAddress)
}

func TestSurfClientRejectsUntrustedTLS(t *testing.T) {
	t.Parallel()

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("request must not reach a server with an untrusted certificate")
	}))
	defer server.Close()

	client := surfTestClient(t, server)
	_, err := client.Get("https://public.example.com/")
	require.ErrorContains(t, err, "certificate signed by unknown authority")
}

func TestSurfClientTimesOutReadingBody(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}))
	defer server.Close()

	client := surfTestClient(t, server)
	client.Timeout = 100 * time.Millisecond
	response, err := client.Get("http://public.example.com/")
	require.NoError(t, err)
	defer response.Body.Close()

	_, err = io.ReadAll(response.Body)
	var timeoutError net.Error
	require.ErrorAs(t, err, &timeoutError)
	assert.True(t, timeoutError.Timeout())
}

func surfTestClient(t *testing.T, server *httptest.Server) *http.Client {
	t.Helper()

	resolve := func(_ context.Context, _ string, host string) ([]netip.Addr, error) {
		if host == "private.example.com" {
			return []netip.Addr{netip.MustParseAddr("10.0.0.1")}, nil
		}

		return []netip.Addr{netip.MustParseAddr("93.184.216.34")}, nil
	}
	client := NewClient(Config{Timeout: 5 * time.Second, Resolver: resolve})
	transport := client.Transport.(*http.Transport)
	transport.DialContext = Guard(resolve, func(ctx context.Context, network, address string) (net.Conn, error) {
		host, _, err := net.SplitHostPort(address)
		if err != nil {
			return nil, err
		}

		assert.Equal(t, "93.184.216.34", host, "dial the validated address without another DNS lookup")
		return (&net.Dialer{}).DialContext(ctx, network, server.Listener.Addr().String())
	})
	client, err := withSurf(client)
	require.NoError(t, err)
	t.Cleanup(client.CloseIdleConnections)
	return client
}
