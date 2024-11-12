package pagination

import (
	"fmt"
	"testing"

	"github.com/Southclaws/dt"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestNewPageParams(t *testing.T) {
	a := assert.New(t)

	p := NewPageParams(0, 10)
	a.Equal(Page(1), p.page)
	a.Equal(Size(10), p.size)
	a.Equal(11, p.Limit())
	a.Equal(0, p.Offset())

	p = NewPageParams(1, 10)
	a.Equal(Page(1), p.page)
	a.Equal(Size(10), p.size)
	a.Equal(11, p.Limit())
	a.Equal(0, p.Offset())

	p = NewPageParams(2, 10)
	a.Equal(Page(2), p.page)
	a.Equal(Size(10), p.size)
	a.Equal(11, p.Limit())
	a.Equal(10, p.Offset())
}

func TestNewPageResult(t *testing.T) {
	a := assert.New(t)

	tableSize := 123

	// a rudimentary fake database of integers
	// we use lo.Slice to simulate a limit+offset sql query
	database := lo.Range(tableSize)

	t.Run("page_1", func(t *testing.T) {
		p := NewPageParams(1, 10)

		rows := lo.Slice(database, p.Offset(), p.Offset()+p.Limit())

		r := NewPageResult(p, tableSize, rows)

		a.Equal(10, r.Size)
		a.Equal(10, r.Results)
		a.Equal(13, r.TotalPages)
		a.Equal(1, r.CurrentPage)
		a.Equal(2, r.NextPage.OrZero())
		a.Len(r.Items, 10)
	})

	t.Run("page_7", func(t *testing.T) {
		p := NewPageParams(7, 10)

		rows := lo.Slice(database, p.Offset(), p.Offset()+p.Limit())

		r := NewPageResult(p, tableSize, rows)

		a.Equal(10, r.Size)
		a.Equal(10, r.Results)
		a.Equal(13, r.TotalPages)
		a.Equal(7, r.CurrentPage)
		a.Equal(8, r.NextPage.OrZero())
		a.Len(r.Items, 10)
	})

	t.Run("page_13", func(t *testing.T) {
		p := NewPageParams(13, 10)

		rows := lo.Slice(database, p.Offset(), p.Offset()+p.Limit())

		r := NewPageResult(p, tableSize, rows)

		a.Equal(10, r.Size)
		a.Equal(3, r.Results)
		a.Equal(13, r.TotalPages)
		a.Equal(13, r.CurrentPage)
		a.False(r.NextPage.Ok(), "there are no more pages after 13")
		a.Len(r.Items, 3, "the final page contains the final 3 items")
	})

	t.Run("weaviate_bug", func(t *testing.T) {
		p := NewPageParams(1, 10)

		// due to complexities of vector searches, sometimes the total count and
		// the actual query don't agree, the total may be in the hundreds but an
		// actual query yields a single page with fewer items than the page size
		// due to autocut and other things.
		rows := lo.Slice(database, p.Offset(), p.Offset()+3)

		r := NewPageResult(p, tableSize, rows)

		a.Equal(10, r.Size)
		a.Equal(3, r.Results)
		a.Equal(1, r.TotalPages)
		a.Equal(1, r.CurrentPage)
		a.False(r.NextPage.Ok(), "there are no more pages after 1")
		a.Len(r.Items, 3, "only 3 items are returned")
	})
}

func TestNewPageResultOffByOnes(t *testing.T) {
	a := assert.New(t)

	const pageSize = 10
	tableSize := 100

	// a rudimentary fake database of integers
	// we use lo.Slice to simulate a limit+offset sql query
	database := lo.Range(tableSize)

	t.Run("page_1", func(t *testing.T) {
		p := NewPageParams(1, pageSize)

		rows := lo.Slice(database, p.Offset(), p.Offset()+p.Limit())

		r := NewPageResult(p, tableSize, rows)

		a.Equal(pageSize, r.Size)
		a.Equal(10, r.Results)
		a.Equal(10, r.TotalPages)
		a.Equal(1, r.CurrentPage)
		a.Equal(2, r.NextPage.OrZero())
		a.Len(r.Items, 10)
	})

	t.Run("page_10", func(t *testing.T) {
		p := NewPageParams(10, pageSize)

		rows := lo.Slice(database, p.Offset(), p.Offset()+p.Limit())

		r := NewPageResult(p, tableSize, rows)

		a.Equal(pageSize, r.Size)
		a.Equal(10, r.Results)
		a.Equal(10, r.TotalPages)
		a.Equal(10, r.CurrentPage)
		a.False(r.NextPage.Ok())
		a.Len(r.Items, 10)
	})

	t.Run("edge_case_page_zero", func(t *testing.T) {
		p := NewPageParams(0, pageSize)

		rows := lo.Slice(database, p.Offset(), p.Offset()+p.Limit())

		r := NewPageResult(p, tableSize, rows)

		a.Equal(pageSize, r.Size)
		a.Equal(10, r.Results)
		a.Equal(10, r.TotalPages)
		a.Equal(1, r.CurrentPage)
		a.Equal(2, r.NextPage.OrZero())
		a.Len(r.Items, 10)
	})

	t.Run("edge_case_page_size_zero", func(t *testing.T) {
		p := NewPageParams(0, 0)

		rows := lo.Slice(database, p.Offset(), p.Offset()+p.Limit())

		r := NewPageResult(p, tableSize, rows)

		a.Equal(1, r.Size)
		a.Equal(1, r.Results)
		a.Equal(100, r.TotalPages)
		a.Equal(1, r.CurrentPage)
		a.Equal(2, r.NextPage.OrZero())
		a.Len(r.Items, 1)
	})
}

func TestConvertPageResult(t *testing.T) {
	a := assert.New(t)

	tableSize := 123
	database := lo.Range(tableSize)

	p := NewPageParams(1, 10)
	rows := lo.Slice(database, p.Offset(), p.Offset()+p.Limit())

	r := NewPageResult(p, tableSize, rows)

	mapped := dt.Map(r.Items, func(r int) string { return fmt.Sprint(r) })

	converted := ConvertPageResult(r, mapped)
	a.Equal(10, converted.Size)
	a.Equal(10, converted.Results)
	a.Equal(13, converted.TotalPages)
	a.Equal(1, converted.CurrentPage)
	a.Equal(2, converted.NextPage.OrZero())
	a.Len(converted.Items, 10)
}
