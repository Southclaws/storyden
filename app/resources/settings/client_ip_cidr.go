package settings

import (
	"net/netip"
	"strconv"
	"strings"
)

// ParseTrustedProxyPrefix parses a trusted proxy identifier into a normalised
// prefix. It accepts canonical CIDR, host-form CIDR, and plain IP values.
func ParseTrustedProxyPrefix(value string) (netip.Prefix, bool) {
	v := strings.TrimSpace(value)
	if v == "" {
		return netip.Prefix{}, false
	}

	// Plain IPs are accepted as single-host prefixes.
	if !strings.Contains(v, "/") {
		addr, err := netip.ParseAddr(v)
		if err != nil {
			return netip.Prefix{}, false
		}
		addr = addr.Unmap()
		bits := 128
		if addr.Is4() {
			bits = 32
		}
		return netip.PrefixFrom(addr, bits), true
	}

	// Canonical CIDR form, fast path.
	if prefix, err := netip.ParsePrefix(v); err == nil {
		return prefix.Masked(), true
	}

	// Host-form prefixes like 172.16.38.226/24.
	addrPart, bitsPart, ok := strings.Cut(v, "/")
	if !ok {
		return netip.Prefix{}, false
	}
	addr, err := netip.ParseAddr(strings.TrimSpace(addrPart))
	if err != nil {
		return netip.Prefix{}, false
	}
	addr = addr.Unmap()

	bits, err := strconv.Atoi(strings.TrimSpace(bitsPart))
	if err != nil {
		return netip.Prefix{}, false
	}

	maxBits := 128
	if addr.Is4() {
		maxBits = 32
	}
	if bits < 0 || bits > maxBits {
		return netip.Prefix{}, false
	}

	return netip.PrefixFrom(addr, bits).Masked(), true
}

// NormaliseTrustedProxyCIDRs trims and normalises trusted proxy values while
// returning a separate list of invalid entries for validation errors.
func NormaliseTrustedProxyCIDRs(values []string) ([]string, []string) {
	out := make([]string, 0, len(values))
	invalid := make([]string, 0)
	for _, value := range values {
		v := strings.TrimSpace(value)
		if v == "" {
			continue
		}

		prefix, ok := ParseTrustedProxyPrefix(v)
		if !ok {
			invalid = append(invalid, v)
			continue
		}

		out = append(out, prefix.String())
	}

	return out, invalid
}
