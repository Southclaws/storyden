package httpsafe

import (
	"context"
	"net"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsDisallowedAddr(t *testing.T) {
	t.Parallel()

	for _, ip := range []string{"127.0.0.1", "10.0.0.1", "192.168.1.1", "172.16.0.1", "169.254.169.254", "::1", "fc00::1", "fe80::1", "0.0.0.0"} {
		assert.True(t, IsDisallowedAddr(netip.MustParseAddr(ip)), ip)
	}
	for _, ip := range []string{"8.8.8.8", "1.1.1.1", "2606:4700:4700::1111"} {
		assert.False(t, IsDisallowedAddr(netip.MustParseAddr(ip)), ip)
	}
	// ipv4-mapped ipv6 form of the metadata address must also be caught
	assert.True(t, IsDisallowedAddr(netip.MustParseAddr("::ffff:169.254.169.254")))
}

func TestGuardRejectsDisallowedResolution(t *testing.T) {
	t.Parallel()

	resolve := func(context.Context, string, string) ([]netip.Addr, error) {
		return []netip.Addr{netip.MustParseAddr("127.0.0.1")}, nil
	}
	dialed := false
	dial := func(context.Context, string, string) (net.Conn, error) {
		dialed = true
		return nil, nil
	}

	_, err := Guard(resolve, dial)(context.Background(), "tcp", "evil.example.com:443")
	require.ErrorIs(t, err, ErrDisallowedAddress)
	assert.False(t, dialed, "must not dial when resolution is disallowed")
}

func TestGuardDialsValidatedLiteralNotHostname(t *testing.T) {
	t.Parallel()

	// resolve reports a public address, and the guard must dial THAT literal so a
	// second lookup (dns rebinding to a private ip) can never happen
	resolve := func(context.Context, string, string) ([]netip.Addr, error) {
		return []netip.Addr{netip.MustParseAddr("93.184.216.34")}, nil
	}
	var dialedAddress string
	dial := func(_ context.Context, _ string, address string) (net.Conn, error) {
		dialedAddress = address
		return nil, nil
	}

	_, err := Guard(resolve, dial)(context.Background(), "tcp", "example.com:443")
	require.NoError(t, err)
	assert.Equal(t, "93.184.216.34:443", dialedAddress)
}

func TestGuardRejectsDisallowedLiteralHost(t *testing.T) {
	t.Parallel()

	resolve := func(context.Context, string, string) ([]netip.Addr, error) {
		t.Fatal("resolve must not be called for a literal ip host")
		return nil, nil
	}
	_, err := Guard(resolve, nil)(context.Background(), "tcp", "169.254.169.254:80")
	require.ErrorIs(t, err, ErrDisallowedHost)
}
