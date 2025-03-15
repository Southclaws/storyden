package node_querier

import (
	"context"
	"fmt"
	"strings"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/rs/xid"
)

type ChildSortRule struct {
	Dir   string
	Field string
	Page  pagination.Parameters
	raw   string
}

func NewChildSortRule(raw string, pp pagination.Parameters) ChildSortRule {
	field := raw
	dir := "asc"
	if strings.HasPrefix(raw, "-") {
		dir = "desc"
		field = raw[1:]
	}

	switch raw {
	default:
		return ChildSortRule{
			Dir:   dir,
			Field: field,
			raw:   raw,
			Page:  pp,
		}
	}
}

const querySortedByPropertyValue = `
select
  n.id id
from
  nodes n
  left join properties p on n.id = p.node_id
  inner join property_schema_fields f on p.field_id = f.id and f.name = '%s'
where
  n.id in (%s)
order by
  case
    f.type

	-- TODO: Add the actual types here when we implement stricter property types
    when 'number' then cast(p.value as integer)
    when 'date' then cast(p.value as datetime)
    else p.value

  end %s

limit  %d
offset %d
`

func (q *Querier) sortedByPropertyValue(ctx context.Context, ids []string, csr ChildSortRule) (map[xid.ID]int, error) {
	var rows []struct {
		ID xid.ID `db:"id"`
	}

	quotedIDs := dt.Map(ids, func(id string) string { return fmt.Sprintf("'%s'", id) })
	idList := strings.Join(quotedIDs, ",")

	// NOTE: Safe injection here as csr.Dir is statically assigned to either
	// "asc" or "desc" in NewChildSortRule. Unfortunately both Go and SQL don't
	// allow parameterizing the ORDER BY direction because nothing about SQL has
	// improved since 1973...
	withParams := fmt.Sprintf(querySortedByPropertyValue,
		csr.Field,
		idList,
		csr.Dir,
		csr.Page.Limit(),
		csr.Page.Offset(),
	)

	err := q.raw.SelectContext(ctx, &rows, withParams)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := make(map[xid.ID]int, len(rows))
	for i, row := range rows {
		result[row.ID] = i
	}

	return result, nil
}
