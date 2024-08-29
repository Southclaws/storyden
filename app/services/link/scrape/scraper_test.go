package scrape

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_scraper_Scrape(t *testing.T) {
	ctx := context.Background()
	sc := New()

	wc, err := sc.Scrape(ctx, "https://ogp.me")

	require.NoError(t, err)
	assert.Equal(t, "Open Graph protocol", wc.Title)
	assert.Equal(t, "The Open Graph protocol enables any web page to become a rich object in a social graph.", wc.Description)
	assert.Equal(t, "https://ogp.me/logo.png", wc.Image)
}
