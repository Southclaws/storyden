package rpc_handler

import (
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/plugin"
)

func mustPluginID(t *testing.T, raw string) plugin.InstallationID {
	t.Helper()
	id, err := xid.FromString(raw)
	require.NoError(t, err)
	return plugin.InstallationID(id)
}

func TestPluginAccessHandleSuffixDeterministic(t *testing.T) {
	id := mustPluginID(t, "d6k8gido2dtljn5ct60g")

	got1 := pluginAccessHandle("sharedbot", id)
	got2 := pluginAccessHandle("sharedbot", id)

	assert.Equal(t, "sharedbot-t60g", got1)
	assert.Equal(t, got1, got2)
}

func TestPluginAccessHandleLengthBounded(t *testing.T) {
	id := mustPluginID(t, "d6k8gido2dtljn5ct60g")

	got := pluginAccessHandle("this-handle-name-is-way-too-long-for-account-rules", id)

	assert.LessOrEqual(t, len(got), 30)
	assert.Equal(t, "this-handle-name-is-way-t-t60g", got)
}

func TestPluginAccessHandleSanitisesInvalidInput(t *testing.T) {
	id := mustPluginID(t, "d6k8gido2dtljn5ct60g")

	got := pluginAccessHandle("  !!My Plugin!!  ", id)

	assert.Equal(t, "my-plugin-t60g", got)
}

func TestPluginAccessRoleNameFormat(t *testing.T) {
	got := pluginAccessRoleName("Discord Connector", "discord-connector-a1b2")

	assert.Equal(t, "Discord Connector (Bot a1b2)", got)
}

func TestPluginAccessRoleShortIDFallback(t *testing.T) {
	assert.Equal(t, "xy00", pluginAccessRoleShortID("xy"))
	assert.Equal(t, "0000", pluginAccessRoleShortID("___"))
}

func TestFindManagedAccessRoleMatchesInstallationMetadata(t *testing.T) {
	id := mustPluginID(t, "d6k8gido2dtljn5ct60g")

	got := findManagedAccessRole(role.Roles{
		&role.Role{
			Name: "Some Other Role",
			Metadata: map[string]any{
				pluginAccessRoleMetaInstallationID: id.String(),
			},
		},
	}, id)

	require.NotNil(t, got)
	assert.Equal(t, "Some Other Role", got.Name)
}

func TestFindManagedAccessRoleIgnoresNameOnlyMatch(t *testing.T) {
	id := mustPluginID(t, "d6k8gido2dtljn5ct60g")

	got := findManagedAccessRole(role.Roles{
		&role.Role{
			Name: "Discord Connector (Bot t60g)",
		},
	}, id)

	assert.Nil(t, got)
}
