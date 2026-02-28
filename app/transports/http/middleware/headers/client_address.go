package headers

import (
	"net"
	"net/http"
	"strings"
)

func clientAddress(r *http.Request) string {
	if addr := parseForwardedValues(r.Header.Values("Forwarded")); addr != "" {
		return addr
	}
	if addr := parseXForwardedForValues(r.Header.Values("X-Forwarded-For")); addr != "" {
		return addr
	}
	return parseRemoteAddr(r.RemoteAddr)
}

func parseForwardedValues(values []string) string {
	for i := len(values) - 1; i >= 0; i-- {
		header := strings.TrimSpace(values[i])
		if header == "" {
			continue
		}

		elements := strings.Split(header, ",")
		for j := len(elements) - 1; j >= 0; j-- {
			for _, param := range strings.Split(elements[j], ";") {
				key, value, ok := strings.Cut(strings.TrimSpace(param), "=")
				if !ok || !strings.EqualFold(key, "for") {
					continue
				}

				addr := normaliseAddrToken(value)
				if addr != "" && !strings.EqualFold(addr, "unknown") {
					return addr
				}
			}
		}
	}

	return ""
}

func parseXForwardedForValues(values []string) string {
	for i := len(values) - 1; i >= 0; i-- {
		header := strings.TrimSpace(values[i])
		if header == "" {
			continue
		}

		parts := strings.Split(header, ",")
		for j := len(parts) - 1; j >= 0; j-- {
			addr := normaliseAddrToken(parts[j])
			if addr != "" && !strings.EqualFold(addr, "unknown") {
				return addr
			}
		}
	}

	return ""
}

func parseRemoteAddr(remoteAddr string) string {
	addr := strings.TrimSpace(remoteAddr)
	if addr == "" {
		return ""
	}

	host, _, err := net.SplitHostPort(addr)
	if err == nil {
		return host
	}

	return addr
}

func normaliseAddrToken(token string) string {
	v := strings.TrimSpace(token)
	if v == "" {
		return ""
	}

	v = strings.Trim(v, `"`)
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}

	host, _, err := net.SplitHostPort(v)
	if err == nil {
		return host
	}

	return v
}
