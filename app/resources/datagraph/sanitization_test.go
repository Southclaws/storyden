package datagraph

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRichTextSanitizesScriptableMarkup(t *testing.T) {
	r := require.New(t)

	c, err := NewRichText(`<body><p>safe</p><img src=x onerror=alert(1)><script>alert(1)</script><iframe src="https://evil.example"></iframe></body>`)
	r.NoError(err)

	html := c.HTML()
	r.Contains(html, "<body>")
	r.Contains(html, "<p>safe</p>")
	r.NotContains(strings.ToLower(html), "onerror")
	r.NotContains(strings.ToLower(html), "<script")
	r.NotContains(strings.ToLower(html), "<iframe")
}

func TestNewRichTextFromMarkdownDoesNotPreserveArbitraryClasses(t *testing.T) {
	r := require.New(t)

	c, err := NewRichTextFromMarkdown("```js&#x20;xss\nalert(1)\n```")
	r.NoError(err)

	html := c.HTML()
	r.Contains(html, "<code>")
	r.NotContains(html, `class="language-js xss"`)
	r.NotContains(html, `class="xss"`)
}
