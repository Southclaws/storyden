package plugin_auth

import (
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/plugin"
)

func TestTokenRoundTrip(t *testing.T) {
	secret, err := GenerateSecret()
	require.NoError(t, err)
	require.Len(t, secret, secretLength)

	for _, char := range secret {
		assert.Contains(t, charset, string(char), "secret should only contain charset characters")
	}

	pluginID := plugin.InstallationID(xid.New())

	token, err := SealToken(pluginID, secret)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	decrypted, err := OpenToken(token, secret)
	require.NoError(t, err)
	assert.Equal(t, pluginID, decrypted)
}

func TestTokenWithWrongSecret(t *testing.T) {
	secret1, err := GenerateSecret()
	require.NoError(t, err)

	secret2, err := GenerateSecret()
	require.NoError(t, err)

	pluginID := plugin.InstallationID(xid.New())

	token, err := SealToken(pluginID, secret1)
	require.NoError(t, err)

	_, err = OpenToken(token, secret2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decrypt token")
}

func TestTokenUniqueness(t *testing.T) {
	secret, err := GenerateSecret()
	require.NoError(t, err)

	pluginID := plugin.InstallationID(xid.New())

	token1, err := SealToken(pluginID, secret)
	require.NoError(t, err)

	token2, err := SealToken(pluginID, secret)
	require.NoError(t, err)

	assert.NotEqual(t, token1, token2, "tokens should be unique due to random nonce")

	decrypted1, err := OpenToken(token1, secret)
	require.NoError(t, err)
	assert.Equal(t, pluginID, decrypted1)

	decrypted2, err := OpenToken(token2, secret)
	require.NoError(t, err)
	assert.Equal(t, pluginID, decrypted2)
}

func TestGenerateExternalToken(t *testing.T) {
	token, err := GenerateExternalToken()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	assert.True(t, IsExternalToken(token))
}
