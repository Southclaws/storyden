package session_cookie

import (
	"net/url"
	"slices"
	"strings"

	"github.com/Southclaws/fault"
	"golang.org/x/net/publicsuffix"
)

type Domain []string

func DomainFromURL(u url.URL) (Domain, error) {
	return DomainFromString(u.Hostname())
}

func DomainFromString(s string) (Domain, error) {
	if s == "" {
		return nil, fault.New("empty domain")
	}

	if s != "localhost" {
		_, err := publicsuffix.EffectiveTLDPlusOne(s)
		if err != nil {
			return nil, fault.Wrap(err)
		}
	}

	parts := strings.Split(s, ".")
	slices.Reverse(parts)

	return parts, nil
}

func (d Domain) String() string {
	parts := d
	slices.Reverse(parts)
	return strings.Join(parts, ".")
}

func (d Domain) IsSubdomainOf(other Domain) bool {
	if len(d) < len(other) {
		return false
	}

	for i := 0; i < len(other); i++ {
		if d[i] != other[i] {
			return false
		}
	}

	return true
}

func (d Domain) IsSiblingOf(other Domain) bool {
	if len(d) != len(other) {
		return false
	}

	length := len(d)

	for i := 0; i < length-1; i++ {
		if d[i] != other[i] {
			return false
		}
	}

	return true
}

func (d Domain) IsEqual(other Domain) bool {
	if len(d) != len(other) {
		return false
	}

	for i := 0; i < len(other); i++ {
		if d[i] != other[i] {
			return false
		}
	}

	return true
}

func (d Domain) IsLocalhost() bool {
	return len(d) == 1 && d[0] == "localhost"
}

func (d Domain) IsTopLevel() bool {
	tldp1, err := publicsuffix.EffectiveTLDPlusOne(strings.Join(d, "."))
	if err != nil {
		return false
	}

	return strings.Join(d, ".") == tldp1
}

func (d Domain) GetETLDp1() Domain {
	s := d.String()

	if s == "localhost" {
		return d
	}

	// error already handled during d's construction
	etldp1, _ := publicsuffix.EffectiveTLDPlusOne(s)
	// only error from this is from the above call, so won't happen.
	d2, _ := DomainFromString(etldp1)

	return d2
}

func getCookieDomain(backend, frontend url.URL) (string, error) {
	// If both frontend and backend are hosted on the same domain just use that.
	if backend.Hostname() == frontend.Hostname() {
		return backend.Hostname(), nil
	}

	backendDomain, err := DomainFromURL(backend)
	if err != nil {
		return "", err
	}

	frontendDomain, err := DomainFromURL(frontend)
	if err != nil {
		return "", err
	}

	// If frontend and backend are on different subdomains use backend's domain.
	// example: api.cats.com and site.cats.com
	if frontendDomain.IsSiblingOf(backendDomain) {
		return backendDomain.String(), nil
	}

	// If frontend is on a lower level domain than the backend use TLD+1.
	// example: api.cats.com and cats.com
	if backendDomain.IsSubdomainOf(frontendDomain) {
		return frontendDomain.String(), nil
	}

	return backendDomain.String(), nil
}
