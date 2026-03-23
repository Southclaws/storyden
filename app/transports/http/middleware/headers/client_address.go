package headers

import (
	"context"
	"net"
	"net/http"
	"net/netip"
	"strings"

	"github.com/Southclaws/storyden/app/resources/settings"
)

const (
	defaultClientIPHeader = "X-Real-IP"
	ssrRequestHeader      = "X-Storyden-SSR"
)

type clientIPConfiguration struct {
	Mode               settings.ClientIPMode
	Header             string
	TrustedProxyCIDRs  []string
	trustedProxyRanges []netip.Prefix
}

func defaultClientIPConfiguration() clientIPConfiguration {
	return clientIPConfiguration{
		Mode:   settings.ClientIPModeRemoteAddr,
		Header: defaultClientIPHeader,
	}
}

func (m *Middleware) reloadClientIPConfiguration(ctx context.Context) {
	cfg, err := m.getClientIPConfiguration(ctx)
	if err != nil {
		m.logger.Warn("failed to fetch settings for client IP configuration",
			"error", err.Error(),
		)
		return
	}

	m.clientIPConfig.Store(cfg)
}

func (m *Middleware) getClientIPConfiguration(ctx context.Context) (clientIPConfiguration, error) {
	cfg := defaultClientIPConfiguration()

	appSettings, err := m.settingsRepo.Get(ctx)
	if err != nil {
		return cfg, err
	}

	services, ok := appSettings.Services.Get()
	if !ok {
		return cfg, nil
	}

	clientIP, ok := services.ClientIP.Get()
	if !ok {
		return cfg, nil
	}

	if v, ok := clientIP.ClientIPMode.Get(); ok {
		cfg.Mode = v
	}
	if v, ok := clientIP.ClientIPHeader.Get(); ok {
		header := strings.TrimSpace(v)
		if header != "" {
			cfg.Header = header
		}
	}
	if v, ok := clientIP.TrustedProxyCIDRs.Get(); ok {
		cfg.TrustedProxyCIDRs = v
		cfg.trustedProxyRanges = parseTrustedProxyCIDRs(v)
	}

	return cfg, nil
}

func (m *Middleware) clientAddress(r *http.Request) string {
	cfg := m.currentClientIPConfiguration()
	return m.clientAddressWithConfig(r, cfg)
}

func (m *Middleware) currentClientIPConfiguration() clientIPConfiguration {
	cfgAny := m.clientIPConfig.Load()
	cfg, ok := cfgAny.(clientIPConfiguration)
	if !ok {
		cfg = defaultClientIPConfiguration()
	}
	return cfg
}

func (m *Middleware) clientAddressWithConfig(r *http.Request, cfg clientIPConfiguration) string {
	if key := getClientIPKey(r, cfg); key != "" {
		return key
	}

	return strings.TrimSpace(r.RemoteAddr)
}

func (m *Middleware) ssrClientAddress(r *http.Request, resolvedClientAddress string) string {
	if strings.TrimSpace(r.Header.Get(ssrRequestHeader)) == "" {
		return ""
	}

	return resolvedClientAddress
}

func parseTrustedProxyCIDRs(cidrs []string) []netip.Prefix {
	prefixes := make([]netip.Prefix, 0, len(cidrs))
	for _, c := range cidrs {
		prefix, ok := settings.ParseTrustedProxyPrefix(c)
		if !ok {
			continue
		}
		prefixes = append(prefixes, prefix)
	}
	return prefixes
}

func getClientIPKey(r *http.Request, cfg clientIPConfiguration) string {
	remote := parseRemoteAddrIP(r.RemoteAddr)
	if remote == "" {
		remote = strings.TrimSpace(r.RemoteAddr)
	}

	switch cfg.Mode {
	case settings.ClientIPModeSingleHeader:
		headerName := strings.TrimSpace(cfg.Header)
		if headerName == "" {
			headerName = defaultClientIPHeader
		}
		if ip := parseIPToken(r.Header.Get(headerName)); ip != "" {
			return ip
		}
		return remote
	case settings.ClientIPModeXFFTrustedProxies:
		return getTrustedProxyXFFIP(r, remote, cfg.trustedProxyRanges, isSSRRequest(r))
	default:
		return remote
	}
}

func getTrustedProxyXFFIP(r *http.Request, remote string, trusted []netip.Prefix, allowLoopback bool) string {
	remoteAddr, ok := parseAddr(remote)
	if !ok {
		return remote
	}
	if !isTrustedProxy(remoteAddr, trusted) && !(allowLoopback && remoteAddr.IsLoopback()) {
		return remote
	}

	chain := parseXFFChain(r.Header.Values("X-Forwarded-For"))
	if len(chain) == 0 {
		return remote
	}

	for i := len(chain) - 1; i >= 0; i-- {
		addr, ok := parseAddr(chain[i])
		if !ok {
			continue
		}
		if isTrustedProxy(addr, trusted) {
			continue
		}
		return addr.String()
	}

	return remote
}

func parseXFFChain(values []string) []string {
	out := make([]string, 0)
	for _, v := range values {
		if strings.TrimSpace(v) == "" {
			continue
		}
		for _, p := range strings.Split(v, ",") {
			if ip := parseIPToken(p); ip != "" {
				out = append(out, ip)
			}
		}
	}
	return out
}

func parseIPToken(value string) string {
	v := strings.Trim(strings.TrimSpace(value), `"`)
	if v == "" || strings.EqualFold(v, "unknown") {
		return ""
	}

	if host, _, err := net.SplitHostPort(v); err == nil {
		v = host
	}

	if addr, ok := parseAddr(v); ok {
		return addr.String()
	}

	return ""
}

func parseRemoteAddrIP(remoteAddr string) string {
	v := strings.TrimSpace(remoteAddr)
	if v == "" {
		return ""
	}

	if host, _, err := net.SplitHostPort(v); err == nil {
		v = host
	}

	if addr, ok := parseAddr(v); ok {
		return addr.String()
	}

	return ""
}

func parseAddr(v string) (netip.Addr, bool) {
	addr, err := netip.ParseAddr(strings.Trim(strings.TrimSpace(v), "[]"))
	if err != nil {
		return netip.Addr{}, false
	}
	return addr.Unmap(), true
}

func isTrustedProxy(addr netip.Addr, trusted []netip.Prefix) bool {
	for _, p := range trusted {
		if p.Contains(addr) {
			return true
		}
	}
	return false
}

func isSSRRequest(r *http.Request) bool {
	return strings.TrimSpace(r.Header.Get(ssrRequestHeader)) != ""
}
