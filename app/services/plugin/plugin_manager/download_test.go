package plugin_manager

import (
	"context"
	"net/netip"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidatePluginSourceURL(t *testing.T) {
	t.Parallel()

	t.Run("allows public https", func(t *testing.T) {
		u := mustParseURL(t, "https://example.com/plugin.zip")
		require.NoError(t, validatePluginSourceURL(*u))
	})

	t.Run("rejects localhost", func(t *testing.T) {
		u := mustParseURL(t, "http://localhost/plugin.zip")
		require.Error(t, validatePluginSourceURL(*u))
	})

	t.Run("rejects private ip", func(t *testing.T) {
		u := mustParseURL(t, "http://192.168.1.10/plugin.zip")
		require.Error(t, validatePluginSourceURL(*u))
	})

	t.Run("rejects unsupported scheme", func(t *testing.T) {
		u := mustParseURL(t, "ftp://example.com/plugin.zip")
		require.Error(t, validatePluginSourceURL(*u))
	})
}

func TestIsDisallowedAddr(t *testing.T) {
	t.Parallel()

	require.True(t, isDisallowedAddr(netip.MustParseAddr("127.0.0.1")))
	require.True(t, isDisallowedAddr(netip.MustParseAddr("10.0.0.1")))
	require.True(t, isDisallowedAddr(netip.MustParseAddr("fc00::1")))

	require.False(t, isDisallowedAddr(netip.MustParseAddr("8.8.8.8")))
	require.False(t, isDisallowedAddr(netip.MustParseAddr("2606:4700:4700::1111")))
}

func TestValidateResolvedHostRejectsPrivateDNSResolution(t *testing.T) {
	origLookup := lookupNetIP
	t.Cleanup(func() { lookupNetIP = origLookup })

	lookupNetIP = func(context.Context, string, string) ([]netip.Addr, error) {
		return []netip.Addr{netip.MustParseAddr("127.0.0.1")}, nil
	}

	err := validateResolvedHost(context.Background(), "example.com")
	require.Error(t, err)
	require.ErrorContains(t, err, "disallowed")
}

func TestValidateResolvedHostAllowsPublicDNSResolution(t *testing.T) {
	origLookup := lookupNetIP
	t.Cleanup(func() { lookupNetIP = origLookup })

	lookupNetIP = func(context.Context, string, string) ([]netip.Addr, error) {
		return []netip.Addr{netip.MustParseAddr("8.8.8.8")}, nil
	}

	err := validateResolvedHost(context.Background(), "example.com")
	require.NoError(t, err)
}

func TestReadAllBounded(t *testing.T) {
	t.Parallel()

	t.Run("within limit", func(t *testing.T) {
		b, err := readAllBounded(strings.NewReader("hello"), 5, "payload")
		require.NoError(t, err)
		require.Equal(t, "hello", string(b))
	})

	t.Run("exceeds limit", func(t *testing.T) {
		_, err := readAllBounded(strings.NewReader("hello!"), 5, "payload")
		require.Error(t, err)
		require.ErrorContains(t, err, "payload too large")
	})
}

func mustParseURL(t *testing.T, raw string) *url.URL {
	t.Helper()
	u, err := url.Parse(raw)
	require.NoError(t, err)
	return u
}
