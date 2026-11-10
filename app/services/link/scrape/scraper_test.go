package scrape

import (
	"context"
	"testing"

	"github.com/Southclaws/storyden/app/resources/link/link_url"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_scraper_Scrape(t *testing.T) {
	ctx := context.Background()
	sc := New()

	u, err := link_url.Parse("https://ogp.me/")
	require.NoError(t, err)

	wc, err := sc.Scrape(ctx, u)

	require.NoError(t, err)
	assert.Equal(t, "Open Graph protocol", wc.Title)
	assert.Equal(t, "The Open Graph protocol enables any web page to become a rich object in a social graph.", wc.Description)
	assert.Equal(t, "https://ogp.me/logo.png", wc.Image)
	assert.Equal(t, "https://ogp.me/favicon.ico", wc.Favicon)
}

func TestScrapeRejectsZeroURL(t *testing.T) {
	t.Parallel()
	_, err := New().Scrape(context.Background(), link_url.URL{})
	require.ErrorIs(t, err, link_url.ErrEmpty)
}
