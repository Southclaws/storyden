package pagination

import (
	"math"

	"github.com/Southclaws/opt"
)

type (
	Page = uint
	Size = uint
)

// Parameters is to be used in any paginated repository query method for paging.
type Parameters struct {
	page Page
	size Size
}

// Result holds a list of rows from a paginated query result with page metadata.
type Result[T any] struct {
	Size        int
	Results     int
	TotalPages  int
	CurrentPage int
	NextPage    opt.Optional[int]
	Items       []T
}

// NewPageParams creates a new pagination parameter object that is 1-indexed, so
// 0 as a page number will default to 1. Use pagination parameters to construct
// queries using Limit and Offset, then use the parameter object in the result.
func NewPageParams(oneIndexedPageNumber uint, pageSize uint) Parameters {
	if oneIndexedPageNumber == 0 {
		oneIndexedPageNumber = 1
	}

	if pageSize == 0 {
		pageSize = 1
	}

	return Parameters{
		page: Page(oneIndexedPageNumber),
		size: Size(pageSize),
	}
}

func (p Parameters) PageOneIndexed() int {
	return int(p.page)
}

func (p Parameters) PageZeroIndexed() int {
	return int(p.page - 1)
}

func (p Parameters) Size() int {
	return int(p.size)
}

// Limit returns a query limit clause value. The actual limit is +1 because we
// want to determine if there are more results by checking if the returned rows
// are above the page size, this would mean there's another page available.
func (p Parameters) Limit() int {
	return int(p.size + 1)
}

// Offset returns an offset clause value for a page query.
func (p Parameters) Offset() int {
	return int(p.page-1) * int(p.size)
}

// NewPageResult constructs a paginated results object containing metadata about
// the paged results such as whether there are more pages available to query.
func NewPageResult[T any](p Parameters, total int, r []T) Result[T] {
	totalPages := int(math.Ceil(float64(total) / float64(p.size)))

	moreResults := len(r) > int(p.size)
	nextPage := opt.NewSafe(int(p.page+1), moreResults)

	var trimmed []T
	if moreResults {
		trimmed = r[:len(r)-1]
	} else {
		trimmed = r
	}

	// The number of rows after trimming, not the len of r because if there is
	// another page, the length of r will be +1 from the actual page size.
	results := len(trimmed)

	// A bit of a hack for Weaviate due to paging with vector searches being a
	// bit unpredictable. This happens when Autocut is used and total is wrong.
	if p.page == 1 && !moreResults && total > int(p.size) {
		totalPages = 1
		nextPage = opt.NewEmpty[int]()
	}

	return Result[T]{
		Size:        int(p.size),
		Results:     results,
		TotalPages:  totalPages,
		CurrentPage: int(p.page),
		NextPage:    nextPage,
		Items:       trimmed,
	}
}

func ConvertPageResult[F, T any](f Result[F], t []T) Result[T] {
	return Result[T]{
		Size:        f.Size,
		Results:     f.Results,
		TotalPages:  f.TotalPages,
		CurrentPage: f.CurrentPage,
		NextPage:    f.NextPage,
		Items:       t,
	}
}
