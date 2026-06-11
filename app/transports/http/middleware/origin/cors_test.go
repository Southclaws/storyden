package origin

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsOriginAllowed_SameHost(t *testing.T) {
	r := require.New(t)

	webAddr := u("http://localhost:3000")
	apiAddr := u("http://localhost:8000")

	r.True(isOriginAllowed("localhost", webAddr, apiAddr))
}

func TestIsOriginAllowed_SameSubdomain(t *testing.T) {
	r := require.New(t)

	webAddr := u("https://www.example.com")
	apiAddr := u("https://api.example.com")

	r.True(isOriginAllowed("www.example.com", webAddr, apiAddr))
	r.True(isOriginAllowed("api.example.com", webAddr, apiAddr))
	r.True(isOriginAllowed("example.com", webAddr, apiAddr))
	r.True(isOriginAllowed("mail.example.com", webAddr, apiAddr))
	r.False(isOriginAllowed("www.example.com.evil.com", webAddr, apiAddr))
}

func TestIsOriginAllowed_APIAndWebSameHost(t *testing.T) {
	r := require.New(t)

	webAddr := u("https://example.com")
	apiAddr := u("https://example.com")

	r.True(isOriginAllowed("example.com", webAddr, apiAddr))
	r.True(isOriginAllowed("www.example.com", webAddr, apiAddr))
	r.True(isOriginAllowed("api.example.com", webAddr, apiAddr))
	r.True(isOriginAllowed("sub.domain.example.com", webAddr, apiAddr))
	r.False(isOriginAllowed("example.org", webAddr, apiAddr))
	r.False(isOriginAllowed("notexample.com", webAddr, apiAddr))
}

func TestIsOriginAllowed_DifferentRootDomains(t *testing.T) {
	r := require.New(t)

	webAddr := u("https://api.example.com")
	apiAddr := u("https://other.com")

	r.True(isOriginAllowed("api.example.com", webAddr, apiAddr))
	r.True(isOriginAllowed("other.com", webAddr, apiAddr))
	r.False(isOriginAllowed("www.example.com", webAddr, apiAddr))
	r.False(isOriginAllowed("www.other.com", webAddr, apiAddr))
	r.False(isOriginAllowed("example.com", webAddr, apiAddr))
}

func TestIsOriginAllowed_InvalidOrigin(t *testing.T) {
	r := require.New(t)

	webAddr := u("https://example.com")
	apiAddr := u("https://example.com")

	r.False(isOriginAllowed("", webAddr, apiAddr))
	r.False(isOriginAllowed("invalid domain", webAddr, apiAddr))
	r.False(isOriginAllowed("192.168.1.1", webAddr, apiAddr))
}

func TestIsOriginAllowed_WithPorts(t *testing.T) {
	r := require.New(t)

	webAddr := u("https://example.com:3000")
	apiAddr := u("https://example.com:8000")

	r.True(isOriginAllowed("example.com", webAddr, apiAddr))
	r.True(isOriginAllowed("www.example.com", webAddr, apiAddr))
	r.True(isOriginAllowed("api.example.com", webAddr, apiAddr))
}

func TestIsOriginAllowed_LocalhostVariations(t *testing.T) {
	r := require.New(t)

	webAddr := u("http://localhost:3000")
	apiAddr := u("http://localhost:8000")

	r.True(isOriginAllowed("localhost", webAddr, apiAddr))
	r.False(isOriginAllowed("localhost.com", webAddr, apiAddr))
	r.False(isOriginAllowed("127.0.0.1", webAddr, apiAddr))
}

func TestIsSubdomainOfRoot(t *testing.T) {
	r := require.New(t)

	r.True(isSubdomainOfRoot("example.com", "example.com"))
	r.True(isSubdomainOfRoot("www.example.com", "example.com"))
	r.True(isSubdomainOfRoot("api.example.com", "example.com"))
	r.True(isSubdomainOfRoot("sub.domain.example.com", "example.com"))
	r.False(isSubdomainOfRoot("notexample.com", "example.com"))
	r.False(isSubdomainOfRoot("example.com.evil.com", "example.com"))
	r.False(isSubdomainOfRoot("localhost:3000", "localhost"))
	r.False(isSubdomainOfRoot("localhost", "localhost:3000"))
}

func u(s string) url.URL {
	u, _ := url.Parse(s)
	return *u
}
