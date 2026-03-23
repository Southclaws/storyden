package headers

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/config"
)

func TestParseTrustedSSRSourceRangesWithLocalhostFrontendProxy(t *testing.T) {
	t.Parallel()

	frontendProxyURL, err := url.Parse("http://localhost:3000")
	require.NoError(t, err)

	cfg := config.Config{
		FrontendProxy:         *frontendProxyURL,
		SSRTrustedSourceCIDRs: "203.0.113.9,\n2001:db8::7/64,not-a-cidr",
	}

	ranges, invalid := parseTrustedSSRSourceRanges(cfg)
	got := make([]string, 0, len(ranges))
	for _, p := range ranges {
		got = append(got, p.String())
	}

	assert.ElementsMatch(t, []string{
		"127.0.0.1/32",
		"::1/128",
		"203.0.113.9/32",
		"2001:db8::/64",
	}, got)
	assert.Equal(t, []string{"not-a-cidr"}, invalid)
}

func TestParseTrustedSSRSourceRangesWithLiteralFrontendProxyIP(t *testing.T) {
	t.Parallel()

	frontendProxyURL, err := url.Parse("http://[2001:db8::9]:3000")
	require.NoError(t, err)

	cfg := config.Config{
		FrontendProxy: *frontendProxyURL,
	}

	ranges, invalid := parseTrustedSSRSourceRanges(cfg)
	require.Empty(t, invalid)
	require.Len(t, ranges, 1)
	assert.Equal(t, "2001:db8::9/128", ranges[0].String())
}

func TestSplitTrustedSourceCIDRList(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"2001:db8::/32",
	}, splitTrustedSourceCIDRList("10.0.0.0/8,\n 172.16.0.0/12 , 2001:db8::/32"))
}
