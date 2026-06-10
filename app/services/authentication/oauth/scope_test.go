package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Southclaws/storyden/app/resources/rbac"
)

func TestValidateClientRequestedScopes(t *testing.T) {
	tests := []struct {
		name          string
		requestedScope string
		allowedScopes  []string
		expectError    bool
	}{
		{
			name:           "exact_match_allowed",
			requestedScope: "CREATE_POST READ_PUBLISHED_THREADS",
			allowedScopes:  []string{"CREATE_POST", "READ_PUBLISHED_THREADS"},
			expectError:    false,
		},
		{
			name:           "standard_scopes_allowed",
			requestedScope: "openid profile email",
			allowedScopes:  []string{"openid", "profile", "email"},
			expectError:    false,
		},
		{
			name:           "mixed_standard_and_permission_scopes",
			requestedScope: "openid profile CREATE_POST",
			allowedScopes:  []string{"openid", "profile", "email", "CREATE_POST"},
			expectError:    false,
		},
		{
			name:           "requested_scope_not_in_allowed",
			requestedScope: "CREATE_POST DELETE_POST",
			allowedScopes:  []string{"CREATE_POST"},
			expectError:    true,
		},
		{
			name:           "administrator_allows_all_permissions",
			requestedScope: "CREATE_POST READ_PUBLISHED_THREADS MANAGE_POSTS",
			allowedScopes:  []string{rbac.PermissionAdministrator.String()},
			expectError:    false,
		},
		{
			name:           "administrator_with_standard_scopes",
			requestedScope: "openid profile email CREATE_POST READ_PUBLISHED_THREADS",
			allowedScopes:  []string{"openid", "profile", "email", "offline_access", rbac.PermissionAdministrator.String()},
			expectError:    false,
		},
		{
			name:           "standard_scope_still_required_with_administrator",
			requestedScope: "openid profile email CREATE_POST",
			allowedScopes:  []string{"openid", "profile", rbac.PermissionAdministrator.String()},
			expectError:    true,
		},
		{
			name:           "empty_scope",
			requestedScope: "",
			allowedScopes:  []string{"openid"},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateClientRequestedScopes(tt.requestedScope, tt.allowedScopes)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
