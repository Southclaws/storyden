package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStreamJSONLEmitsOneLinePerItem(t *testing.T) {
	r := require.New(t)

	type item struct {
		Slug string `json:"slug"`
	}

	var buf bytes.Buffer
	r.NoError(StreamJSONL(&buf, []item{{Slug: "a"}, {Slug: "b"}}))

	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	r.Len(lines, 2)
	r.JSONEq(`{"slug":"a"}`, lines[0])
	r.JSONEq(`{"slug":"b"}`, lines[1])
}

func TestStreamJSONLEmptySlice(t *testing.T) {
	r := require.New(t)

	type item struct{}

	var buf bytes.Buffer
	r.NoError(StreamJSONL(&buf, []item{}))
	r.Empty(buf.String())
}
