package pluginbuilder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorydenDocsURLAllowsDocsAndLLMSText(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  string
	}{
		{input: "/llms.txt", want: "https://www.storyden.org/llms.txt"},
		{input: "llms.txt", want: "https://www.storyden.org/llms.txt"},
		{input: "/docs/plugins.md", want: "https://www.storyden.org/docs/plugins.md"},
		{input: "https://storyden.org/docs/plugins?ignored=true#section", want: "https://www.storyden.org/docs/plugins"},
	} {
		t.Run(tc.input, func(t *testing.T) {
			u, err := storydenDocsURL(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.want, u.String())
		})
	}
}

func TestStorydenDocsURLRejectsNonDocsPaths(t *testing.T) {
	for _, input := range []string{
		"https://example.com/docs/plugins",
		"http://www.storyden.org/docs/plugins",
		"/admin",
		"/api/secrets",
	} {
		t.Run(input, func(t *testing.T) {
			_, err := storydenDocsURL(input)
			require.Error(t, err)
		})
	}
}
