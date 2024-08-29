package scrape

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_scraper_Scrape(t *testing.T) {
	ctx := context.Background()
	sc := New()

	u, _ := url.Parse("https://ogp.me/")

	wc, err := sc.Scrape(ctx, *u)

	require.NoError(t, err)
	assert.Equal(t, "Open Graph protocol", wc.Title)
	assert.Equal(t, "The Open Graph protocol enables any web page to become a rich object in a social graph.", wc.Description)
	assert.Equal(t, "https://ogp.me/logo.png", wc.Image)
	assert.Equal(t, "https://ogp.me/favicon.ico", wc.Favicon)
}
