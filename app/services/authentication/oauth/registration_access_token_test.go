package oauth

import (
	"context"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/settings"
)

func TestAgentRegistrationRejectsUnknownMode(t *testing.T) {
	service := &Service{}
	result, oauthErr, err := service.authoriseAgentRegistration(context.Background(), settings.OAuthAutonomousRegistrationMode{}, true, opt.NewEmpty[string]())

	require.NoError(t, err)
	require.False(t, result.Ok())
	require.NotNil(t, oauthErr)
	require.Equal(t, "invalid_client_metadata", oauthErr.Code)
}

func TestAgentRegistrationRequiresOAuth(t *testing.T) {
	service := &Service{}
	result, oauthErr, err := service.RegisterClient(context.Background(), DynamicClientRegistration{TokenEndpointAuthMethod: "private_key_jwt"})
	require.NoError(t, err)
	require.Nil(t, result)
	require.NotNil(t, oauthErr)
	require.Equal(t, "temporarily_unavailable", oauthErr.Code)
}
