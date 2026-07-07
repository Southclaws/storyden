package plugin_manager

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/internal/infrastructure/httpsafe"
)

const (
	pluginFetchTimeout      = 30 * time.Second
	pluginFetchMaxRedirects = 5
)

var lookupNetIP = net.DefaultResolver.LookupNetIP

func readAllBounded(r io.Reader, maxBytes int64, label string) ([]byte, error) {
	limited := io.LimitReader(r, maxBytes+1)
	b, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > maxBytes {
		return nil, fault.New(
			"payload exceeds maximum size",
			fmsg.WithDesc(
				label+" too large",
				label+" must be at most "+byteCountString(maxBytes),
			),
		)
	}
	return b, nil
}

func byteCountString(n int64) string {
	if n%(1024*1024) == 0 {
		return fmt.Sprintf("%d MiB", n/(1024*1024))
	}
	if n%1024 == 0 {
		return fmt.Sprintf("%d KiB", n/1024)
	}
	return fmt.Sprintf("%d bytes", n)
}

func validatePluginSourceURL(u url.URL) error {
	scheme := strings.ToLower(strings.TrimSpace(u.Scheme))
	if scheme != "http" && scheme != "https" {
		return fault.New("plugin URL must use http or https")
	}
	if strings.TrimSpace(u.Hostname()) == "" {
		return fault.New("plugin URL host is required")
	}
	if u.User != nil {
		return fault.New("plugin URL must not include user info")
	}
	if isDisallowedHost(u.Hostname()) {
		return fault.New("plugin URL host is not allowed")
	}
	return nil
}

func isDisallowedHost(host string) bool {
	return httpsafe.IsDisallowedHost(host)
}

func isDisallowedAddr(addr netip.Addr) bool {
	return httpsafe.IsDisallowedAddr(addr)
}

func validateResolvedHost(ctx context.Context, host string) error {
	host = strings.TrimSpace(host)
	if host == "" {
		return fault.New("plugin URL host is required")
	}
	if isDisallowedHost(host) {
		return fault.New("plugin URL host is not allowed")
	}

	// Host literals are already validated above.
	if _, err := netip.ParseAddr(host); err == nil {
		return nil
	}

	addrs, err := lookupNetIP(ctx, "ip", host)
	if err != nil {
		return fault.Wrap(err, fmsg.With("failed to resolve plugin URL host"))
	}
	if len(addrs) == 0 {
		return fault.New("plugin URL host did not resolve to any addresses")
	}

	for _, addr := range addrs {
		if isDisallowedAddr(addr) {
			return fault.New("plugin URL host resolves to a disallowed address")
		}
	}

	return nil
}

func fetchPluginArchive(ctx context.Context, u url.URL) ([]byte, error) {
	if err := validatePluginSourceURL(u); err != nil {
		return nil, err
	}

	client := httpsafe.NewClient(httpsafe.Config{
		Timeout:      pluginFetchTimeout,
		MaxRedirects: pluginFetchMaxRedirects,
		UseEnvProxy:  true,
		Resolver:     lookupNetIP,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fault.Newf("failed to fetch plugin from URL: status code %d", resp.StatusCode)
	}
	if resp.ContentLength > plugin.MaxArchiveSizeBytes {
		return nil, fault.Newf("plugin archive too large: %d bytes", resp.ContentLength)
	}

	return readAllBounded(resp.Body, plugin.MaxArchiveSizeBytes, "plugin archive")
}
