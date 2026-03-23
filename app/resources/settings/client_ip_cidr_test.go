package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTrustedProxyPrefix(t *testing.T) {
	t.Parallel()

	p, ok := ParseTrustedProxyPrefix("172.16.38.226/24")
	require.True(t, ok)
	assert.Equal(t, "172.16.38.0/24", p.String())

	p, ok = ParseTrustedProxyPrefix("2001:db8::7/64")
	require.True(t, ok)
	assert.Equal(t, "2001:db8::/64", p.String())

	p, ok = ParseTrustedProxyPrefix("203.0.113.7")
	require.True(t, ok)
	assert.Equal(t, "203.0.113.7/32", p.String())
}

func TestNormaliseTrustedProxyCIDRs(t *testing.T) {
	t.Parallel()

	out, invalid := NormaliseTrustedProxyCIDRs([]string{
		"",
		" 10.0.0.0/8 ",
		"172.16.38.226/24",
		"2001:db8::7/64",
		"not-a-cidr",
	})

	assert.Equal(t, []string{
		"10.0.0.0/8",
		"172.16.38.0/24",
		"2001:db8::/64",
	}, out)
	assert.Equal(t, []string{"not-a-cidr"}, invalid)
}
