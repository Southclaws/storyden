package keycloak

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/config"
)

func TestProviderDisplayNameDefault(t *testing.T) {
	t.Parallel()

	p, err := New(config.Config{KeycloakEnabled: false}, nil, nil)
	require.NoError(t, err)

	assert.Equal(t, "Keycloak", p.DisplayName())
}

func TestProviderDisplayNameCustom(t *testing.T) {
	t.Parallel()

	p, err := New(config.Config{
		KeycloakEnabled:     false,
		KeycloakDisplayName: "Corporate SSO",
	}, nil, nil)
	require.NoError(t, err)

	assert.Equal(t, "Corporate SSO", p.DisplayName())
}
