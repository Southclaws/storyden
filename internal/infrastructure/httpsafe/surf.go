package httpsafe

import (
	"net/http"

	surfhttp "github.com/enetx/http"
	"github.com/enetx/surf"
)

func NewSurfClient(cfg Config) (*http.Client, error) {
	return withSurf(NewClient(cfg))
}

func withSurf(client *http.Client) (*http.Client, error) {
	guarded := client.Transport.(*http.Transport)
	sc := surf.NewClient()
	transport := sc.GetTransport().(*surfhttp.Transport)
	transport.DialContext = guarded.DialContext
	transport.Proxy = nil
	if guarded.Proxy != nil {
		transport.Proxy = surfhttp.ProxyFromEnvironment
	}

	result := sc.Builder().Impersonate().MacOS().Chrome().SecureTLS().Build()
	if result.IsErr() {
		return nil, result.Err()
	}

	client.Transport = result.Ok().Std().Transport
	return client, nil
}
