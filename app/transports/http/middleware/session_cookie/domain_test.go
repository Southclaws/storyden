package session_cookie

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDomainType(t *testing.T) {
	r := require.New(t)

	d, err := DomainFromString("localhost")
	r.NoError(err)
	r.Equal(Domain{"localhost"}, d)
	r.Equal("localhost", d.String())

	d, err = DomainFromString("example.com")
	r.NoError(err)
	r.Equal(Domain{"com", "example"}, d)
	r.Equal("example.com", d.String())

	d, err = DomainFromString("sub.example.com")
	r.NoError(err)
	r.Equal(Domain{"com", "example", "sub"}, d)
	r.Equal("sub.example.com", d.String())

	d, err = DomainFromString("sub.sub.example.com")
	r.NoError(err)
	r.Equal(Domain{"com", "example", "sub", "sub"}, d)
	r.Equal("sub.sub.example.com", d.String())

	_, err = DomainFromString("")
	r.Error(err)

	_, err = DomainFromString("example")
	r.Error(err)
}

func TestIsSubdomainOf(t *testing.T) {
	r := require.New(t)

	d1, _ := DomainFromString("sub.example.com")
	d2, _ := DomainFromString("example.com")

	r.True(d1.IsSubdomainOf(d2))
	r.False(d2.IsSubdomainOf(d1))

	d3, _ := DomainFromString("sub.example.com")
	d4, _ := DomainFromString("dom.example.com")

	r.False(d3.IsSubdomainOf(d4))
	r.False(d4.IsSubdomainOf(d3))
}

func TestIsSiblingOf(t *testing.T) {
	r := require.New(t)

	d1, _ := DomainFromString("sub.example.com")
	d2, _ := DomainFromString("example.com")

	r.False(d1.IsSiblingOf(d2))
	r.False(d2.IsSiblingOf(d1))

	d3, _ := DomainFromString("sub.example.com")
	d4, _ := DomainFromString("dom.example.com")

	r.True(d3.IsSiblingOf(d4))
	r.True(d4.IsSiblingOf(d3))
}

func TestIsEqual(t *testing.T) {
	r := require.New(t)

	d1, _ := DomainFromString("sub.example.com")
	d2, _ := DomainFromString("example.com")

	r.False(d1.IsEqual(d2))
	r.False(d2.IsEqual(d1))

	d3, _ := DomainFromString("sub.example.com")
	d4, _ := DomainFromString("sub.example.com")

	r.True(d3.IsEqual(d4))
	r.True(d4.IsEqual(d3))
}

func TestIsLocalhost(t *testing.T) {
	r := require.New(t)

	d1, _ := DomainFromString("localhost")
	r.True(d1.IsLocalhost())

	d2, _ := DomainFromString("example.com")
	r.False(d2.IsLocalhost())
}

func TestIsETLDPlus1(t *testing.T) {
	r := require.New(t)

	d1, _ := DomainFromString("localhost")
	r.False(d1.IsTopLevel())

	d2, _ := DomainFromString("example.com")
	r.True(d2.IsTopLevel())

	d3, _ := DomainFromString("sub.example.com")
	r.False(d3.IsTopLevel())
}

func TestGetETLDp1(t *testing.T) {
	r := require.New(t)

	d1, _ := DomainFromString("localhost")
	r.Equal(Domain{"localhost"}, d1.GetETLDp1())

	d2, _ := DomainFromString("example.com")
	r.Equal(Domain{"com", "example"}, d2.GetETLDp1())

	d3, _ := DomainFromString("sub.example.com")
	r.Equal(Domain{"com", "example"}, d3.GetETLDp1())
}

func TestGetCookieDomain(t *testing.T) {
	r := require.New(t)

	d1, err := getCookieDomain(u("http://localhost:8080"), u("http://localhost:8080"))
	r.NoError(err)
	r.Equal("localhost", d1)

	d2, err := getCookieDomain(u("https://api.makeroom.club"), u("https://makeroom.club"))
	r.NoError(err)
	r.Equal("makeroom.club", d2)

	d3, err := getCookieDomain(u("https://api.makeroom.club"), u("https://www.makeroom.club"))
	r.NoError(err)
	r.Equal("api.makeroom.club", d3)

	d4, err := getCookieDomain(u("https://makeroom.club"), u("https://community.makeroom.club"))
	r.NoError(err)
	r.Equal("makeroom.club", d4)

	d5, err := getCookieDomain(u("https://makeroom.club"), u("https://makeroom.club"))
	r.NoError(err)
	r.Equal("makeroom.club", d5)
}

func u(s string) url.URL {
	u, _ := url.Parse(s)
	return *u
}
