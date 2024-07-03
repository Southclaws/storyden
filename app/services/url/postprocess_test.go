package url

import (
	"bytes"
	"context"
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed testdata/*.html
var testContent embed.FS

const dir = "testdata"

func Test_webScraper_postprocess(t *testing.T) {
	r := require.New(t)
	// a := assert.New(t)
	ctx := context.Background()

	es, err := testContent.ReadDir(dir)
	r.NoError(err)

	w := webScraper{}

	for _, e := range es {
		filename := filepath.Join(dir, e.Name())
		// title := strings.TrimSuffix(e.Name(), ".html")

		b, err := os.ReadFile(filename)
		r.NoError(err)

		wc, err := w.postprocess(ctx, "https://storyden.org", bytes.NewReader(b))
		r.NoError(err)
		r.NotNil(wc)

		// NOTE: quite exploratory at the moment so we're not asserting anything
		os.WriteFile(filepath.Join("data", e.Name()), []byte(wc.Text), fs.ModePerm)
	}
}
