package node_querier

import (
	"context"
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/pagination"
)

type ChildSortRule struct {
	Dir   string
	Field string
	Page  pagination.Parameters
	Fixed bool // true means the field is a schema field not a property
	raw   string
}

func (s ChildSortRule) OrderClause() sql.OrderTermOption {
	if s.Dir == "asc" {
		return sql.OrderAsc()
	}
	return sql.OrderDesc()
}

func NewChildSortRule(raw string, pp pagination.Parameters) ChildSortRule {
	field := raw
	dir := "asc"
	if strings.HasPrefix(raw, "-") {
		dir = "desc"
		field = raw[1:]
	}

	csr := ChildSortRule{
		Dir:   dir,
		Field: field,
		raw:   raw,
		Page:  pp,
	}

	switch field {
	// NOTE: The same as `MappableNodeField` in web codebase.
	case "name", "link", "description":
		csr.Fixed = true
	}

	return csr
}

const querySortedByPropertyValue_sqlite = `
select
  n.id id
from
  nodes n
  left join properties p on n.id = p.node_id
  inner join property_schema_fields f on p.field_id = f.id and f.name = $1
where
  n.id in (%s)
order by
  case f.type
    when 'text'      then p.value
    when 'number'    then cast(p.value as real)
    when 'timestamp' then cast(p.value as datetime)
    when 'boolean'   then cast(p.value as integer)
    else p.value

  end %s

limit  %d
offset %d
`

const querySortedByPropertyValue_postgres = `
select
  n.id id
from
  nodes n
  left join properties p on n.id = p.node_id
  inner join property_schema_fields f on p.field_id = f.id and f.name = $1
where
  n.id in (%s)
order by
  case f.type when 'text'      then p.value                            end %s,
  case f.type when 'number'    then cast(p.value as numeric)           end %s,
  case f.type when 'timestamp' then cast(p.value as timestamp)         end %s,
  case f.type when 'boolean'   then cast(p.value as boolean)           end %s,
  p.value %s
limit  %d
offset %d
`

func (q *Querier) sortedByPropertyValue(ctx context.Context, ids []string, csr ChildSortRule) (map[xid.ID]int, error) {
	var rows []struct {
		ID xid.ID `db:"id"`
	}

	quotedIDs := dt.Map(ids, func(id string) string { return fmt.Sprintf("'%s'", id) })
	idList := strings.Join(quotedIDs, ",")

	// you can never convince me sql has a standard that anyone has ever read...
	var queryTemplate string
	switch q.raw.DriverName() {
	// NOTE: Safe injection here as csr.Dir is statically assigned to either
	// "asc" or "desc" in NewChildSortRule. Unfortunately both Go and SQL don't
	// allow parameterizing the ORDER BY direction because nothing about SQL has
	// improved since 1973...
	case "sqlite":
		queryTemplate = fmt.Sprintf(querySortedByPropertyValue_sqlite,
			idList,
			csr.Dir,
			csr.Page.Limit(),
			csr.Page.Offset(),
		)

	case "pgx":
		queryTemplate = fmt.Sprintf(querySortedByPropertyValue_postgres,
			idList,
			csr.Dir, // this
			csr.Dir, // is
			csr.Dir, // so
			csr.Dir, // fuckin
			csr.Dir, // dumb
			csr.Page.Limit(),
			csr.Page.Offset(),
		)
	default:
		return nil, fault.New("unexpected failure in database driver switch")
	}

	err := q.raw.SelectContext(ctx, &rows, queryTemplate, csr.Field)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := make(map[xid.ID]int, len(rows))
	for i, row := range rows {
		result[row.ID] = i
	}

	return result, nil
}
