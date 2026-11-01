package oauth

import (
	"net"
	"net/url"
	"strings"

	"github.com/Southclaws/fault"
)

var unsafeRedirectSchemes = map[string]struct{}{
	"javascript": {},
	"data":       {},
	"vbscript":   {},
	"file":       {},
	"blob":       {},
	"about":      {},
	"filesystem": {},
}

func validateClientRedirectURIs(uris []string) error {
	for _, raw := range uris {
		if err := validateRedirectURI(raw); err != nil {
			return err
		}
	}
	return nil
}

// All registration paths and browser redirects share this policy. Non-web
// schemes support native app callbacks; browser-executable schemes are blocked.
func validateRedirectURI(raw string) error {
	if raw == "" {
		return fault.New("empty redirect uri")
	}

	if strings.Contains(raw, "*") {
		return fault.New("wildcard redirect uri")
	}

	u, err := url.Parse(raw)
	if err != nil {
		return err
	}

	if !u.IsAbs() {
		return fault.New("redirect uri must be absolute")
	}

	if u.User != nil {
		return fault.New("redirect uri must not contain credentials")
	}

	if u.Fragment != "" || strings.Contains(raw, "#") {
		return fault.New("redirect uri must not contain a fragment")
	}

	switch u.Scheme {
	case "https", "http":
		if u.Opaque != "" || u.Hostname() == "" {
			return fault.New("redirect uri must have a valid hostname")
		}
		if u.Scheme == "http" && !isLoopbackHost(u.Hostname()) {
			return fault.New("http redirect uri only allowed for loopback hosts")
		}
	default:
		if _, unsafe := unsafeRedirectSchemes[strings.ToLower(u.Scheme)]; unsafe {
			return fault.Newf("redirect uri scheme %q is not permitted", u.Scheme)
		}
	}
	return nil
}

func isLoopbackHost(host string) bool {
	if host == "localhost" {
		return true
	}

	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
