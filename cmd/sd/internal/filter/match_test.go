package filter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDomainMatches(t *testing.T) {
	r := require.New(t)

	r.True(DomainMatches("tenor.com", "tenor.com"))
	r.True(DomainMatches("media.tenor.com", "tenor.com"))
	r.True(DomainMatches("a.b.tenor.com", "tenor.com"))
	r.False(DomainMatches("nottenor.com", "tenor.com"))
	r.False(DomainMatches("tenor.computer", "tenor.com"))
	r.False(DomainMatches("", "tenor.com"))
	r.True(DomainMatches("anything", ""))

	// Case-insensitive
	r.True(DomainMatches("MEDIA.Tenor.COM", "tenor.com"))
}

func TestURLContains(t *testing.T) {
	r := require.New(t)

	r.True(URLContains("https://www.youtube.com/watch?v=abc", "/watch"))
	r.False(URLContains("https://example.com", "/watch"))
	r.True(URLContains("anything", ""))
	r.True(URLContains("https://EXAMPLE.com/Foo", "foo"))
}

func TestURLScheme(t *testing.T) {
	r := require.New(t)

	r.Equal("https", URLScheme("https://example.com"))
	r.Equal("http", URLScheme("HTTP://example.com"))
	r.Equal("", URLScheme(""))
}

func TestNameContains(t *testing.T) {
	r := require.New(t)

	r.True(NameContains("Hello World", "world"))
	r.True(NameContains("anything", ""))
	r.False(NameContains("Hello", "xyz"))
}
