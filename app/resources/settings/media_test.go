package settings

import (
	"encoding/json"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/ent"
)

func TestImageLimits(t *testing.T) {
	t.Run("missing_values_are_errors", func(t *testing.T) {
		_, err := (&Settings{}).ImageLimits()
		require.Error(t, err)
	})

	t.Run("partial_persisted_policy_is_usable_without_replacing_saved_limits", func(t *testing.T) {
		stored := Settings{
			Services: opt.New(ServiceSettings{
				Media: opt.New(MediaServiceSettings{
					ImageMaxPixels: opt.New[int64](8000),
				}),
			}),
		}
		data, err := json.Marshal(stored)
		require.NoError(t, err)

		repository := &SettingsRepository{}
		configured, err := repository.hydrateConfigDefaults(&ent.Setting{
			ID:    StorydenPrimarySettingsKey,
			Value: string(data),
		})
		require.NoError(t, err)

		limits, err := configured.ImageLimits()
		require.NoError(t, err)
		require.EqualValues(t, 8000, limits.MaxPixels)
		require.Positive(t, limits.MaxSize)
	})

	for _, value := range []int64{0, -1} {
		t.Run("invalid_size", func(t *testing.T) {
			err := (MediaServiceSettings{ImageMaxSize: opt.New(value)}).Validate()
			require.Error(t, err)
		})

		t.Run("invalid_pixels", func(t *testing.T) {
			err := (MediaServiceSettings{ImageMaxPixels: opt.New(value)}).Validate()
			require.Error(t, err)
		})
	}

	t.Run("partial_patch", func(t *testing.T) {
		err := (MediaServiceSettings{ImageMaxPixels: opt.New[int64](6000000)}).Validate()
		require.NoError(t, err)
	})
}
