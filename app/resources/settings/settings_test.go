package settings_test

import (
	"context"
	"testing"

	"go.uber.org/fx"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSettings(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	integration.Test(t, nil, fx.Invoke(func(
		_ integration.Migrated,
		sr settings.Repository,
	) {
		t.Run("set_many_partial", func(t *testing.T) {
			r := require.New(t)
			a := assert.New(t)

			content, err := datagraph.NewRichText("<body><p>Hello, Makeroom!</p></body>")
			r.NoError(err)

			set, err := sr.Set(ctx, settings.Partial{
				Title:   opt.New("Makeroom"),
				Content: opt.New(content),
			})
			r.NoError(err)
			r.NotNil(set)

			got, err := sr.Get(ctx)
			r.NoError(err)
			r.NotNil(got)

			a.Equal("Makeroom", got.Title.Get())
			a.Equal(content.HTML(), got.Content.Get().HTML())
		})

		t.Run("set_individual", func(t *testing.T) {
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

			content, err := datagraph.NewRichText("<body><p>Hello, Storyden!</p></body>")
			r.NoError(err)

			err = s.Content.Set(ctx, sr, content)
			r.NoError(err)

			// The value is assigned to the locally held settings struct
			a.Equal(content.HTML(), s.Content.Get().HTML())

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
		})
	}))
}
