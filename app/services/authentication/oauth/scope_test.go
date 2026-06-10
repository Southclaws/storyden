package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
