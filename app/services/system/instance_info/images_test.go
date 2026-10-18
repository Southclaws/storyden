package instance_info

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/internal/config"
)

func TestDefaultMediaLimitsFitMemoryBudget(t *testing.T) {
	limits, err := settings.DefaultSettings.ImageLimits()
	require.NoError(t, err)

	estimate := estimateImageMemory(limits)
	require.Positive(t, estimate)
	require.LessOrEqual(t, estimate, config.DefaultImageMemoryBudget)
}
