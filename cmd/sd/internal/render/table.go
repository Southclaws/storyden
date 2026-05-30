package render

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/Southclaws/storyden/cmd/sd/internal/output"
)

// Column describes a single column in a Profile. Wide columns are only
// included when the user passes --output wide.
type Column[T any] struct {
	Header string
	Render func(T) string
	Wide   bool
}

// Profile is an ordered list of columns describing how to render a slice of T
// as a plain-text table. Profiles are composed once per command and reused.
type Profile[T any] struct {
	Columns []Column[T]
}

// PageInfo carries optional pagination metadata that lets the renderer print a
// footer beneath the table. Zero values suppress the footer so the renderer
// stays quiet when the backend doesn't populate pagination.
type PageInfo struct {
	CurrentPage int
	TotalPages  int
	PageSize    int
	Results     int
}

// Render writes a plain-text table built from the profile and items to out.
// Wide-only columns are skipped unless wide is true. The footer is printed only
// when PageInfo.TotalPages > 0 so it silently no-ops against backends that
// don't yet populate pagination metadata.
func Render[T any](out io.Writer, items []T, p Profile[T], wide bool, page PageInfo) error {
	writer := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)

	headers := make([]string, 0, len(p.Columns))
	cols := make([]Column[T], 0, len(p.Columns))
	for _, c := range p.Columns {
		if c.Wide && !wide {
			continue
		}
		cols = append(cols, c)
		headers = append(headers, c.Header)
	}

	if _, err := fmt.Fprintln(writer, strings.Join(headers, "\t")); err != nil {
		return err
	}

	width := output.TerminalWidth(out, 0)
	cellLimit := 0
	if width > 0 {
		// Reserve space for separators (2 spaces per column).
		cellLimit = width / max(len(cols), 1)
	}

	for _, item := range items {
		cells := make([]string, len(cols))
		for i, c := range cols {
			cells[i] = clampCell(c.Render(item), cellLimit)
		}
		if _, err := fmt.Fprintln(writer, strings.Join(cells, "\t")); err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	if page.TotalPages > 0 {
		fmt.Fprintf(out, "\nPage %d of %d (showing %d of %d)\n",
			page.CurrentPage, page.TotalPages, page.PageSize, page.Results)
	}

	return nil
}

// clampCell collapses newlines and clips overlong cells so a single huge
// description doesn't blow out the table. limit <= 0 disables clipping.
func clampCell(value string, limit int) string {
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "\r", " ")
	if limit <= 0 || len(value) <= limit {
		return value
	}
	if limit <= 1 {
		return value[:limit]
	}
	return value[:limit-1] + "…"
}
