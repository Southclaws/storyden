package settings_test

import (
	"testing"

	"github.com/Southclaws/opt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/settings"
)

func TestSettingsMerge(t *testing.T) {
	t.Parallel()
	r := require.New(t)
	a := assert.New(t)

	old := settings.Settings{
		Title:       opt.New("Old Title"),
		Description: opt.New("untouched description"),
	}

	updated := settings.Settings{
		Title: opt.New("New Title"),
	}

	err := old.Merge(updated)
	r.NoError(err)

	a.Equal("New Title", old.Title.OrZero())
	a.Equal("untouched description", old.Description.OrZero())
}
