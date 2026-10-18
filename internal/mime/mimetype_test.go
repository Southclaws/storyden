package mime_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/mime"
)

func TestNew(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		input string
		want  string
	}{
		{input: "text/css", want: "text/css"},
		{input: "text/css; charset=utf-8", want: "text/css"},
		{input: "application/javascript", want: "text/javascript"},
		{input: "image/png", want: "image/png"},
		{input: "application/vnd.storyden.example", want: "application/vnd.storyden.example"},
		{input: "invalid", want: "application/octet-stream"},
		{input: "invalid MIME type", want: "application/octet-stream"},
	} {
		t.Run(test.input, func(t *testing.T) {
			require.Equal(t, test.want, mime.New(test.input).String())
		})
	}
}
