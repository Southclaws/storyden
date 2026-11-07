package push_subscription

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/config"
)

func TestSubscriptionLimitConfiguration(t *testing.T) {
	for _, tc := range []struct {
		name    string
		enabled bool
		limit   int
		valid   bool
	}{
		{
			name:    "positive",
			enabled: true,
			limit:   5,
			valid:   true,
		},
		{
			name:    "zero",
			enabled: true,
			limit:   0,
			valid:   false,
		},
		{
			name:    "negative",
			enabled: true,
			limit:   -1,
			valid:   false,
		},
		{
			name:    "disabled_ignores_zero",
			enabled: false,
			limit:   0,
			valid:   true,
		},
		{
			name:    "disabled_ignores_negative",
			enabled: false,
			limit:   -1,
			valid:   true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)
			a := assert.New(t)

			repository, err := New(nil, config.Config{
				WebPushEnabled:          tc.enabled,
				WebPushMaxSubscriptions: tc.limit,
			})
			if !tc.valid {
				r.ErrorContains(err, "WEB_PUSH_MAX_SUBSCRIPTIONS")
				return
			}

			r.NoError(err)
			a.Positive(repository.maxSubscriptions)
		})
	}
}
