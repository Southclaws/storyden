// package httpsafe provides http clients hardened against server-side request forgery
package httpsafe

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"time"
)

var (
	ErrDisallowedHost    = errors.New("host is not allowed")
	ErrDisallowedAddress = errors.New("host resolves to a disallowed address")
	ErrNoResolution      = errors.New("host did not resolve to any addresses")
)

var disallowedPrefixes = []netip.Prefix{
	netip.MustParsePrefix("10.0.0.0/8"),
	netip.MustParsePrefix("100.64.0.0/10"),
	netip.MustParsePrefix("127.0.0.0/8"),
	netip.MustParsePrefix("169.254.0.0/16"),
	netip.MustParsePrefix("172.16.0.0/12"),
	netip.MustParsePrefix("192.168.0.0/16"),
	netip.MustParsePrefix("::1/128"),
	netip.MustParsePrefix("64:ff9b::/96"),
	netip.MustParsePrefix("fc00::/7"),
	netip.MustParsePrefix("fe80::/10"),
}

// resolveFunc matches net.Resolver.LookupNetIP so a stub can be injected in tests
type ResolveFunc func(ctx context.Context, network, host string) ([]netip.Addr, error)

// dialFunc matches net.Dialer.DialContext
type DialFunc func(ctx context.Context, network, address string) (net.Conn, error)

// IsDisallowedAddr reports whether addr is unsafe to reach from the server
func IsDisallowedAddr(addr netip.Addr) bool {
	addr = addr.Unmap()
	if addr.IsLoopback() ||
		addr.IsLinkLocalUnicast() ||
		addr.IsLinkLocalMulticast() ||
		addr.IsMulticast() ||
		addr.IsUnspecified() {
		return true
	}
	for _, prefix := range disallowedPrefixes {
		if prefix.Contains(addr) {
			return true
		}
	}
	return false
}

// IsDisallowedHost reports whether a host literal is unsafe, hostnames are validated after resolution
func IsDisallowedHost(host string) bool {
	host = strings.TrimSpace(strings.ToLower(host))
	if host == "" {
		return true
	}
	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return true
	}

	addr, err := netip.ParseAddr(host)
	if err != nil {
		return false
	}
	return IsDisallowedAddr(addr)
}

// Guard resolves once, refuses any disallowed address, then dials a validated literal to close the dns-rebinding window
func Guard(resolve ResolveFunc, dial DialFunc) DialFunc {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			host, port = address, ""
		}

		if IsDisallowedHost(host) {
			return nil, ErrDisallowedHost
		}

		if addr, err := netip.ParseAddr(host); err == nil {
			if IsDisallowedAddr(addr) {
				return nil, ErrDisallowedAddress
			}
			return dial(ctx, network, address)
		}

		addrs, err := resolve(ctx, "ip", host)
		if err != nil {
			return nil, err
		}
		if len(addrs) == 0 {
			return nil, ErrNoResolution
		}
		for _, addr := range addrs {
			if IsDisallowedAddr(addr) {
				return nil, ErrDisallowedAddress
			}
		}

		var lastErr error
		for _, addr := range addrs {
			target := addr.String()
			if port != "" {
				target = net.JoinHostPort(addr.String(), port)
			}
			conn, err := dial(ctx, network, target)
			if err == nil {
				return conn, nil
			}
			lastErr = err
		}
		return nil, lastErr
	}
}

// Config tunes a guarded client
type Config struct {
	Timeout      time.Duration
	DialTimeout  time.Duration
	MaxRedirects int
	UseEnvProxy  bool
	Resolver     ResolveFunc
}

const defaultMaxRedirects = 10

// NewClient builds an http.Client whose connections are restricted to public addresses
func NewClient(cfg Config) *http.Client {
	resolve := cfg.Resolver
	if resolve == nil {
		resolve = net.DefaultResolver.LookupNetIP
	}

	// dial timeout falls back to the overall timeout so streaming callers can set only DialTimeout
	dialTimeout := cfg.DialTimeout
	if dialTimeout == 0 {
		dialTimeout = cfg.Timeout
	}
	dialer := &net.Dialer{Timeout: dialTimeout}

	var proxy func(*http.Request) (*url.URL, error)
	if cfg.UseEnvProxy {
		proxy = http.ProxyFromEnvironment
	}

	maxRedirects := cfg.MaxRedirects
	if maxRedirects == 0 {
		maxRedirects = defaultMaxRedirects
	}

	return &http.Client{
		Timeout: cfg.Timeout,
		Transport: &http.Transport{
			Proxy:       proxy,
			DialContext: Guard(resolve, dialer.DialContext),
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if maxRedirects < 0 {
				return http.ErrUseLastResponse
			}
			if len(via) >= maxRedirects {
				return errors.New("too many redirects")
			}
			if len(via) > 0 &&
				strings.EqualFold(via[len(via)-1].URL.Scheme, "https") &&
				!strings.EqualFold(req.URL.Scheme, "https") {
				return errors.New("refusing insecure redirect from https to http")
			}
			return nil
		},
	}
}
