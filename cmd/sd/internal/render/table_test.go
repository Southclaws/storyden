package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type row struct {
	Name string
	Slug string
}

func TestRenderSkipsWideColumnsWhenNotWide(t *testing.T) {
	r := require.New(t)

	profile := Profile[row]{
		Columns: []Column[row]{
			{Header: "NAME", Render: func(x row) string { return x.Name }},
			{Header: "SLUG", Render: func(x row) string { return x.Slug }, Wide: true},
		},
	}

	var buf bytes.Buffer
	r.NoError(Render(&buf, []row{{Name: "a", Slug: "s"}}, profile, false, PageInfo{}))

	out := buf.String()
	r.Contains(out, "NAME")
	r.NotContains(out, "SLUG")
	r.NotContains(out, "Page")
}

func TestRenderIncludesWideColumns(t *testing.T) {
	r := require.New(t)

	profile := Profile[row]{
		Columns: []Column[row]{
			{Header: "NAME", Render: func(x row) string { return x.Name }},
			{Header: "SLUG", Render: func(x row) string { return x.Slug }, Wide: true},
		},
	}

	var buf bytes.Buffer
	r.NoError(Render(&buf, []row{{Name: "a", Slug: "s"}}, profile, true, PageInfo{}))

	out := buf.String()
	r.Contains(out, "NAME")
	r.Contains(out, "SLUG")
}

func TestRenderEmitsFooterWhenPaginationKnown(t *testing.T) {
	r := require.New(t)

	profile := Profile[row]{
		Columns: []Column[row]{{Header: "NAME", Render: func(x row) string { return x.Name }}},
	}

	var buf bytes.Buffer
	r.NoError(Render(&buf, []row{{Name: "a"}}, profile, false, PageInfo{
		CurrentPage: 1,
		TotalPages:  5,
		PageSize:    50,
		Results:     247,
	}))

	r.Contains(buf.String(), "Page 1 of 5 (showing 50 of 247)")
}

func TestRenderSuppressesFooterWhenNoPagination(t *testing.T) {
	r := require.New(t)

	profile := Profile[row]{
		Columns: []Column[row]{{Header: "NAME", Render: func(x row) string { return x.Name }}},
	}

	var buf bytes.Buffer
	r.NoError(Render(&buf, []row{{Name: "a"}}, profile, false, PageInfo{}))

	r.NotContains(buf.String(), "Page")
}

func TestClampCellCollapsesNewlines(t *testing.T) {
	r := require.New(t)
	r.Equal("a b c", clampCell("a\nb\rc", 0))
}

func TestClampCellTrimsToLimit(t *testing.T) {
	r := require.New(t)
	out := clampCell(strings.Repeat("x", 50), 10)
	r.Len([]rune(out), 10)
	r.True(strings.HasSuffix(out, "…"))
}
