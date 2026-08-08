package operation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOAuthDeviceConsentRateLimitCosts(t *testing.T) {
	assert.Equal(t, 20, RateLimitCostOverrides[OperationIDOAuthDeviceConsent])
	assert.Equal(t, 20, RateLimitCostOverrides[OperationIDOAuthDeviceConsentSubmit])
}
