package bindings

import (
	"testing"

	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/stretchr/testify/assert"
)

func TestNormaliseCIDRList(t *testing.T) {
	t.Parallel()

	out, invalid := normaliseCIDRList([]string{
		" 10.0.0.0/8 ",
		"",
		"   ",
		"172.16.0.0/12",
		"172.16.38.226/24",
		"203.0.113.7",
	})

	assert.Empty(t, invalid)
	assert.Equal(t, []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"172.16.38.0/24",
		"203.0.113.7/32",
	}, out)
}

func TestNormaliseCIDRListInvalidCIDRs(t *testing.T) {
	t.Parallel()

	out, invalid := normaliseCIDRList([]string{
		"10.0.0.0/8",
		"not-a-cidr",
		"172.16.0.0",
		"2001:db8::/32",
	})

	assert.Equal(t, []string{"10.0.0.0/8", "172.16.0.0/32", "2001:db8::/32"}, out)
	assert.Equal(t, []string{"not-a-cidr"}, invalid)
}

func TestMapOpenAPIClientIPMode(t *testing.T) {
	t.Parallel()

	mode, ok := mapOpenAPIClientIPMode(openapi.RemoteAddr)
	assert.True(t, ok)
	assert.Equal(t, settings.ClientIPModeRemoteAddr, mode)

	mode, ok = mapOpenAPIClientIPMode(openapi.SingleHeader)
	assert.True(t, ok)
	assert.Equal(t, settings.ClientIPModeSingleHeader, mode)

	mode, ok = mapOpenAPIClientIPMode(openapi.XffTrustedProxies)
	assert.True(t, ok)
	assert.Equal(t, settings.ClientIPModeXFFTrustedProxies, mode)
}

func TestMapOpenAPIClientIPModeRejectsUnknown(t *testing.T) {
	t.Parallel()

	_, ok := mapOpenAPIClientIPMode(openapi.ClientIPServiceSettingsClientIpMode("invalid"))
	assert.False(t, ok)
}
