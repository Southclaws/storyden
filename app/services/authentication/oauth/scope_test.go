package oauth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/rbac"
)

func TestValidateScopeNames(t *testing.T) {
	tests := []struct {
		name           string
		requestedScope string
		expectError    bool
	}{
		{
			name:           "valid_permission_scopes",
			requestedScope: "CREATE_POST READ_PUBLISHED_THREADS",
			expectError:    false,
		},
		{
			name:           "valid_standard_scopes",
			requestedScope: "openid profile email",
			expectError:    false,
		},
		{
			name:           "mixed_standard_and_permission_scopes",
			requestedScope: "openid profile CREATE_POST",
			expectError:    false,
		},
		{
			name:           "invalid_scope_name",
			requestedScope: "CREATE_POST TOTALLY_INVALID_SCOPE",
			expectError:    true,
		},
		{
			name:           "administrator_is_valid",
			requestedScope: "CREATE_POST ADMINISTRATOR",
			expectError:    false,
		},
		{
			name:           "empty_scope",
			requestedScope: "",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateScopeNames(tt.requestedScope)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClientCredentialsScopePolicy(t *testing.T) {
	permissions := rbac.NewList(rbac.PermissionCreatePost, rbac.PermissionReadPublishedLibrary)
	for _, tc := range []struct {
		name      string
		policy    oauthresource.ScopePolicy
		requested string
		want      []string
		wantError bool
	}{
		{"inherit_all", oauthresource.ScopePolicyInheritUserPermissions, "", []string{"CREATE_POST", "READ_PUBLISHED_LIBRARY"}, false},
		{"inherit_restricted", oauthresource.ScopePolicyInheritUserPermissions, "READ_PUBLISHED_LIBRARY", []string{"READ_PUBLISHED_LIBRARY"}, false},
		{"inherit_cannot_escalate", oauthresource.ScopePolicyInheritUserPermissions, "ADMINISTRATOR", nil, false},
		{"inherit_invalid_scope", oauthresource.ScopePolicyInheritUserPermissions, "UNKNOWN", nil, true},
		{"explicit_ceiling", oauthresource.ScopePolicyExplicit, "", []string{"CREATE_POST"}, false},
		{"explicit_rejects_outside_allowance", oauthresource.ScopePolicyExplicit, "READ_PUBLISHED_LIBRARY", nil, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := &oauthresource.Client{ScopePolicy: tc.policy, AllowedScopes: []string{"CREATE_POST"}}
			granted, err := grantClientCredentialsScope(tc.requested, client, permissions)
			if tc.wantError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.ElementsMatch(t, tc.want, strings.Fields(granted))
		})
	}
}
