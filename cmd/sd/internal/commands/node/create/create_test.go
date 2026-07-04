package create

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateVisibility(t *testing.T) {
	r := require.New(t)

	r.NoError(validateVisibility(""))
	r.NoError(validateVisibility("draft"))
	r.NoError(validateVisibility("review"))
	r.NoError(validateVisibility("published"))
	r.NoError(validateVisibility("unlisted"))
	r.ErrorContains(validateVisibility("private"), "invalid --visibility: private")
}

func TestContentToHTML(t *testing.T) {
	r := require.New(t)

	// Without the flag, HTML passes through untouched.
	html := "<h1>Title</h1><p>Body</p>"
	out, err := contentToHTML(html, false)
	r.NoError(err)
	r.Equal(html, out)

	// Empty content is a no-op even with the flag set.
	out, err = contentToHTML("", true)
	r.NoError(err)
	r.Empty(out)

	// Markdown is converted to HTML when the flag is set.
	out, err = contentToHTML("# Title\n\nA paragraph with **bold**.", true)
	r.NoError(err)
	r.Contains(out, `<h1>Title</h1>`)
	r.Contains(out, "Title")
	r.Contains(out, "<strong>bold</strong>")
	r.NotContains(out, "# Title")
}
