package settings_test

import (
	"context"
	"testing"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSettings(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	integration.Test(t, nil, fx.Invoke(func(
		sr settings.Repository,
	) {
		r := require.New(t)
		a := assert.New(t)

		s, err := sr.Get(ctx)
		r.NoError(err)
		r.NotNil(s)

		// The app sets a default during init.
		err = s.Title.Set(ctx, sr, "Storyden")
		r.NoError(err)

		// The value is assigned to the locally held settings struct
		a.Equal("Storyden", s.Title.Get())

		// The value is also persisted
		s, err = sr.Get(ctx)
		r.NoError(err)
		r.NotNil(s)

		a.Equal("Storyden", s.Title.Get())

		err = s.Public.Set(ctx, sr, true)
		r.NoError(err)

		// The value is assigned to the locally held settings struct
		a.Equal(true, s.Public.Get())

		// The value is also persisted
		s, err = sr.Get(ctx)
		r.NoError(err)
		r.NotNil(s)

		a.Equal(true, s.Public.Get())
	}))
}
