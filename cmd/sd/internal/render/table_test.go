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

	r.Contains(buf.String(), "Page 1 of 5 (showing 1 of 247)")
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

func TestPaginationShowsActualRows(t *testing.T) {
	for _, tc := range []struct {
		name                        string
		rows, current, pages, total int
		want                        string
	}{
		{"single short page", 8, 1, 1, 8, "Page 1 of 1 (showing 8 of 8)"},
		{"full page", 50, 1, 5, 247, "Page 1 of 5 (showing 50 of 247)"},
		{"last page", 47, 5, 5, 247, "Page 5 of 5 (showing 47 of 247)"},
		{"filtered page", 3, 1, 5, 247, "Page 1 of 5 (showing 3 of 247)"},
		{"empty page", 0, 1, 1, 0, "Page 1 of 1 (showing 0 of 0)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			profile := Profile[row]{Columns: []Column[row]{{Header: "NAME", Render: func(r row) string { return r.Name }}}}
			require.NoError(t, Render(&out, make([]row, tc.rows), profile, false, PageInfo{CurrentPage: tc.current, TotalPages: tc.pages, PageSize: 50, Results: tc.total}))
			require.Contains(t, out.String(), tc.want)
		})
	}
}
