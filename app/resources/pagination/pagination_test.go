package pagination

import (
	"testing"

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

		a.Equal(Size(10), r.Size)
		a.Equal(10, r.Results)
		a.Equal(13, r.TotalPages)
		a.Equal(Page(1), r.CurrentPage)
		a.Equal(Page(2), r.NextPage.OrZero())
		a.Len(r.Items, 10)
	})

	t.Run("page_7", func(t *testing.T) {
		p := NewPageParams(7, 10)

		rows := lo.Slice(database, p.Offset(), p.Offset()+p.Limit())

		r := NewPageResult(p, tableSize, rows)

		a.Equal(Size(10), r.Size)
		a.Equal(10, r.Results)
		a.Equal(13, r.TotalPages)
		a.Equal(Page(7), r.CurrentPage)
		a.Equal(Page(8), r.NextPage.OrZero())
		a.Len(r.Items, 10)
	})

	t.Run("page_13", func(t *testing.T) {
		p := NewPageParams(13, 10)

		rows := lo.Slice(database, p.Offset(), p.Offset()+p.Limit())

		r := NewPageResult(p, tableSize, rows)

		a.Equal(Size(10), r.Size)
		a.Equal(3, r.Results)
		a.Equal(13, r.TotalPages)
		a.Equal(Page(13), r.CurrentPage)
		a.False(r.NextPage.Ok(), "there are no more pages after 13")
		a.Len(r.Items, 3, "the final page contains the final 3 items")
	})
}
