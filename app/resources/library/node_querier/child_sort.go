package node_querier

import (
	"context"
	"fmt"
	"strings"

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

const querySortedByPropertyValue_sqlite = `
select
  n.id id
from
  nodes n
  left join properties p on n.id = p.node_id
  inner join property_schema_fields f on p.field_id = f.id and f.name = '%s'
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
  inner join property_schema_fields f on p.field_id = f.id and f.name = '%s'
where
  n.id in (%s)
order by
  case f.type
    when 'text'      then p.value
    when 'number'    then cast(p.value as numeric)
    when 'timestamp' then cast(p.value as timestamp)
    when 'boolean'   then cast(p.value as boolean)
    else p.value

  end %s

limit  %d
offset %d
`

func (q *Querier) sortedByPropertyValue(ctx context.Context, ids []string, csr ChildSortRule) (map[xid.ID]int, error) {
	var rows []struct {
		ID xid.ID `db:"id"`
	}

	// Because of the shortcomings of the driver, can't use a prepared statement
	// below (sqlite throws a cryptic error for the where in (?) clause and pg
	// throws a different error, sqlx.In did not work for another reason...)
	// so we have to manually escape the field name to prevent injection.
	safeFieldName := strings.Replace(csr.Field, "'", "''", -1)

	quotedIDs := dt.Map(ids, func(id string) string { return fmt.Sprintf("'%s'", id) })
	idList := strings.Join(quotedIDs, ",")

	// you can never convince me sql has a standard that anyone has ever read...
	var queryTemplate string
	switch q.raw.DriverName() {
	case "sqlite":
		queryTemplate = querySortedByPropertyValue_sqlite
	case "postgres":
		queryTemplate = querySortedByPropertyValue_postgres
	default:
		return nil, fault.New("unexpected failure in database driver switch")
	}

	// NOTE: Safe injection here as csr.Dir is statically assigned to either
	// "asc" or "desc" in NewChildSortRule. Unfortunately both Go and SQL don't
	// allow parameterizing the ORDER BY direction because nothing about SQL has
	// improved since 1973...
	withParams := fmt.Sprintf(queryTemplate,
		safeFieldName,
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
