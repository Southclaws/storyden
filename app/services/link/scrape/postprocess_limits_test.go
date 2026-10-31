package scrape

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostprocessBodyLimit(t *testing.T) {
	const limit = 8 * 1024 * 1024
	const page = `<html><head><meta property="og:title" content="Example"></head><body><p>Article text.</p></body></html>`
	addr := url.URL{Scheme: "https", Host: "example.com"}
	scraper := webScraper{}

	for _, size := range []int{len(page), limit} {
		t.Run(strconv.Itoa(size), func(t *testing.T) {
			body := page + strings.Repeat(" ", size-len(page))
			content, err := scraper.postprocess(context.Background(), addr, strings.NewReader(body))

			require.NoError(t, err)
			require.NotNil(t, content)
			assert.Equal(t, "Example", content.Title)
		})
	}

	for _, size := range []int{limit + 1, limit + 1024} {
		t.Run(strconv.Itoa(size), func(t *testing.T) {
			body := strings.NewReader(page + strings.Repeat(" ", size-len(page)))
			content, err := scraper.postprocess(context.Background(), addr, body)

			require.Error(t, err)
			assert.True(t, errors.Is(err, errBodyTooLarge))
			assert.Nil(t, content)
			assert.Equal(t, size-limit-1, body.Len(), "read at most the limit plus one byte")
		})
	}
}
