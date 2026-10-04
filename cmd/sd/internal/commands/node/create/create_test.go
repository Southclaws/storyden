package create

import (
	"testing"

	commandcontent "github.com/Southclaws/storyden/cmd/sd/internal/content"

	"github.com/stretchr/testify/require"
)

func TestContentToHTML(t *testing.T) {
	r := require.New(t)

	// Without the flag, HTML passes through untouched.
	html := "<h1>Title</h1><p>Body</p>"
	out, err := commandcontent.ToHTML(html, false)
	r.NoError(err)
	r.Equal(html, out)

	// Empty content is a no-op even with the flag set.
	out, err = commandcontent.ToHTML("", true)
	r.NoError(err)
	r.Empty(out)

	// Markdown is converted to HTML when the flag is set.
	out, err = commandcontent.ToHTML("# Title\n\nA paragraph with **bold**.", true)
	r.NoError(err)
	r.Contains(out, `<h1>Title</h1>`)
	r.Contains(out, "Title")
	r.Contains(out, "<strong>bold</strong>")
	r.NotContains(out, "# Title")
}
