package openapi_rbac

import (
	"testing"

	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOAuthRemoteCallbackRequiresManageRobots(t *testing.T) {
	t.Parallel()

	sessionRequired, permission := (&Mapping{}).OAuthRemoteCallback()

	require.True(t, sessionRequired)
	require.NotNil(t, permission)
	assert.Equal(t, rbac.PermissionManageRobots, *permission)
}
