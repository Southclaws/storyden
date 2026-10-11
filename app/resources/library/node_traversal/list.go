package node_traversal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/visibility"
)

type ListParams struct {
	CreatedSince          opt.Optional[time.Time]
	CreatedBefore         opt.Optional[time.Time]
	NewestFirst           opt.Optional[bool]
	Root                  opt.Optional[library.NodeID]
	IncludeHiddenChildren opt.Optional[bool]
}

type ListResult struct {
	pagination.Result[*library.Node]
	TotalResults int
}

func (d *database) List(ctx context.Context, page pagination.Parameters, params ListParams, opts ...Filter) (*ListResult, error) {
	f := filters{includeHiddenChildren: params.IncludeHiddenChildren.Or(true)}
	for _, fn := range opts {
		fn(&f)
	}

	query, args, order, err := buildListQuery(params, f, d.raw.DriverName())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var total int
	countQuery := fmt.Sprintf("select count(*) from (%s) eligible_nodes", query)
	if err := d.raw.QueryRowxContext(ctx, d.raw.Rebind(countQuery), args...).Scan(&total); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	pageQuery := fmt.Sprintf("%s order by %s limit ? offset ?", query, order)
	pageArgs := append(args, page.Limit(), page.Offset())
	rows := []subtreeRow{}
	if err := d.raw.SelectContext(ctx, &rows, d.raw.Rebind(pageQuery), pageArgs...); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	ids := dt.Map(rows, func(n subtreeRow) xid.ID { return n.NodeId })
	nodes, err := d.hydrateNodes(ctx, ids)
	if err != nil {
		return nil, err
	}

	return &ListResult{
		Result:       pagination.NewPageResult(page, total, nodes),
		TotalResults: total,
	}, nil
}

func buildListQuery(params ListParams, f filters, driver string) (string, []interface{}, string, error) {
	subtree, args := buildSubtreeSelection(params.Root, f)
	predicates := []string{}

	visibilities := []string{visibility.VisibilityPublished.String()}
	if len(f.visibility) > 0 {
		visibilities = dt.Map(f.visibility, func(v visibility.Visibility) string { return v.String() })
	}
	predicates = append(predicates, "n.visibility in (?)")
	args = append(args, visibilities)

	if acc, ok := f.requestingAccount.Get(); ok && len(f.visibility) > 0 {
		access := "n.visibility = 'published' or n.account_id = ?"
		args = append(args, acc.ID.String())
		if acc.Roles.Permissions().HasAny(rbac.PermissionManageLibrary, rbac.PermissionAdministrator) {
			access += " or n.visibility = 'review'"
		}
		predicates = append(predicates, "("+access+")")
	} else {
		predicates = append(predicates, "n.visibility = 'published'")
	}

	createdAt := "n.created_at"
	bound := "?"
	switch driver {
	case "sqlite", "sqlite3", "libsql":
		createdAt = "julianday(n.created_at)"
		bound = "julianday(?)"
	}
	if since, ok := params.CreatedSince.Get(); ok {
		predicates = append(predicates, createdAt+" >= "+bound)
		args = append(args, since.UTC().Format(time.RFC3339Nano))
	}
	if before, ok := params.CreatedBefore.Get(); ok {
		predicates = append(predicates, createdAt+" < "+bound)
		args = append(args, before.UTC().Format(time.RFC3339Nano))
	}

	order := "depth, node_sort_key"
	if newest, ok := params.NewestFirst.Get(); ok {
		direction := "asc"
		if newest {
			direction = "desc"
		}
		order = fmt.Sprintf("%s %s, n.id %s", createdAt, direction, direction)
	}

	query := fmt.Sprintf(`with traversed_nodes as (
%s
)
select t.*
from traversed_nodes t
join nodes n on n.id = t.node_id
where %s`, subtree, strings.Join(predicates, " and "))
	query, args, err := sqlx.In(query, args...)
	return query, args, order, err
}
