// Package filter contains client-side predicates that slice node lists by
// fields the server can't filter on natively (link domain, url contents,
// presence of a link, parentage, etc.). Flag names mirror JSON field paths so
// `--link-domain` corresponds to `.link.domain`.
package filter

import (
	"net/url"
	"strings"
)

// DomainMatches reports whether candidate matches pattern on dot boundaries:
// "tenor.com" matches "tenor.com" and "media.tenor.com" but not "nottenor.com".
// Empty pattern matches everything; empty candidate matches nothing.
func DomainMatches(candidate, pattern string) bool {
	candidate = strings.ToLower(strings.TrimSpace(candidate))
	pattern = strings.ToLower(strings.TrimSpace(pattern))
	if pattern == "" {
		return true
	}
	if candidate == "" {
		return false
	}
	if candidate == pattern {
		return true
	}
	return strings.HasSuffix(candidate, "."+pattern)
}

// URLContains reports whether u contains needle (case-insensitive). Empty
// needle returns true.
func URLContains(u, needle string) bool {
	if needle == "" {
		return true
	}
	return strings.Contains(strings.ToLower(u), strings.ToLower(needle))
}

// URLScheme returns the scheme portion of u (lowercased). Returns "" if u is
// not a parseable URL or has no scheme.
func URLScheme(u string) string {
	parsed, err := url.Parse(u)
	if err != nil {
		return ""
	}
	return strings.ToLower(parsed.Scheme)
}

// NameContains reports whether name contains needle (case-insensitive). Empty
// needle returns true.
func NameContains(name, needle string) bool {
	if needle == "" {
		return true
	}
	return strings.Contains(strings.ToLower(name), strings.ToLower(needle))
}
