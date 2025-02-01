package node_traversal

import (
	"context"
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/link"
	"github.com/Southclaws/storyden/internal/ent/node"
)

type database struct {
	db  *ent.Client
	raw *sqlx.DB
}

func New(db *ent.Client, raw *sqlx.DB) Repository {
	return &database{db, raw}
}

func (d *database) Root(ctx context.Context, fs ...Filter) ([]*library.Node, error) {
	query := d.db.Node.Query().
		Where(node.ParentNodeIDIsNil()).
		WithOwner(func(aq *ent.AccountQuery) {
			aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
		}).
		WithAssets()

	f := filters{}
	for _, fn := range fs {
		fn(&f)
	}

	if f.rootAccountHandleFilter != nil {
		query.Where(node.HasOwnerWith(account.Handle(*f.rootAccountHandleFilter)))
	}

	if len(f.visibility) > 0 {
		visibilityTypes := dt.Map(f.visibility, func(v visibility.Visibility) node.Visibility {
			return node.Visibility(v.String())
		})

		query.Where(node.VisibilityIn(visibilityTypes...))
	} else {
		query.Where(node.VisibilityIn(node.VisibilityPublished))
	}

	cs, err := query.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	nodes, err := dt.MapErr(cs, library.MapNode(true, nil))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return nodes, nil
}

const ddl = `with recursive children (parent, id, depth) as (
    select
        parent_node_id,
        id,
        0
    from
        nodes
    where %s
union
    select
        d.parent,
        s.id,
        d.depth + 1
    from
        children d
        join nodes s on d.id = s.parent_node_id
)
select
    distinct n.id       node_id,
    n.account_id        node_account_id,
    n.visibility        node_visibility,
	depth
from
    children
    inner join nodes n on n.id = children.id
    inner join accounts a on a.id = n.account_id

-- optional where clause
%s

order by
    depth
`

type subtreeRow struct {
	NodeId         xid.ID                `db:"node_id"`
	NodeAccountId  xid.ID                `db:"node_account_id"`
	NodeVisibility visibility.Visibility `db:"node_visibility"`
	Depth          int                   `db:"depth"`
}

func (d *database) Subtree(ctx context.Context, id opt.Optional[library.NodeID], flatten bool, fs ...Filter) ([]*library.Node, error) {
	f := filters{}
	for _, fn := range fs {
		fn(&f)
	}

	// NOTE: i fucking hate writing raw sql into source code...

	var rootPredicate string
	predicates := []string{}
	args := []interface{}{}
	argOffset := 0

	getPlaceholder := func() string {
		argOffset += 1
		return fmt.Sprintf("$%d", argOffset)
	}

	if parentNodeID, ok := id.Get(); ok {
		args = append(args, parentNodeID.String())
		rootPredicate = fmt.Sprintf("id = cast(%s as text)", getPlaceholder())
	} else {
		rootPredicate = "parent_node_id is null"
	}

	if f.rootAccountHandleFilter != nil {
		predicates = append(predicates, fmt.Sprintf(
			"a.handle = %s",
			getPlaceholder()))

		args = append(args, *f.rootAccountHandleFilter)
	}

	if f.depth != nil {
		predicates = append(predicates, fmt.Sprintf(
			"depth <= %s",
			getPlaceholder()))

		args = append(args, *f.depth)
	}

	additional := ""
	if len(predicates) > 0 {
		additional = "where " + strings.Join(predicates, " AND ")
	}
	q := fmt.Sprintf(ddl, rootPredicate, additional)

	r, err := d.raw.QueryxContext(ctx, q, args...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	allRows := []subtreeRow{}
	for r.Next() {
		c := subtreeRow{}
		err = r.StructScan(&c)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		allRows = append(allRows, c)
	}

	filtered := dt.Filter(allRows, applyFilterRules(f))

	// Now query every row returned from the recursive query hydrating all data.
	ids := dt.Map(filtered, func(n subtreeRow) xid.ID { return n.NodeId })
	nodeRecords, err := d.db.Node.Query().
		Where(node.IDIn(ids...)).
		WithOwner(func(aq *ent.AccountQuery) {
			aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
		}).
		WithPrimaryImage(func(aq *ent.AssetQuery) {
			aq.WithParent()
		}).
		WithAssets().
		WithLink(func(lq *ent.LinkQuery) {
			lq.WithAssets().Order(link.ByCreatedAt(sql.OrderDesc()))
		}).
		WithParent(func(cq *ent.NodeQuery) {
			cq.
				WithAssets().
				WithOwner(func(aq *ent.AccountQuery) {
					aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
				})
		}).All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	hydratedNodeMap := lo.KeyBy(nodeRecords, func(n *ent.Node) xid.ID { return n.ID })
	flat, err := dt.MapErr(filtered, func(n subtreeRow) (*library.Node, error) {
		hydratedNode, exists := hydratedNodeMap[n.NodeId]
		if !exists {
			panic("recursive query result was not present in hydrated node map")
		}

		return library.MapNode(true, nil)(hydratedNode)
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to hydrate nodes"))
	}

	// Early valid return: if we're flattening the tree, no need to build it.
	if flatten {
		return flat, nil
	}

	// Rebuild the flat list into the tree
	tree := buildTree(flat, id)

	return tree, nil
}

func buildTree(hydrated []*library.Node, id opt.Optional[library.NodeID]) []*library.Node {
	var linkChildrenForParent func(library.Node) []*library.Node

	linkChildrenForParent = func(parent library.Node) []*library.Node {
		filteredParent, isFilteringParent := id.Get()

		return dt.Reduce(hydrated, func(prev []*library.Node, curr *library.Node) []*library.Node {
			if p, ok := curr.Parent.Get(); ok && p.Mark.ID() == parent.Mark.ID() {
				// Take a copy because our mutations cannot apply to `flat`.
				copy := *curr

				copy.Nodes = linkChildrenForParent(copy)

				// If the current iteration is not the root node of a parent
				// node (a subtree query) then blank out the parent field since
				// it's a waste to store this information in tree children.
				if isFilteringParent && filteredParent != library.NodeID(copy.Mark.ID()) {
					copy.Parent = opt.NewEmpty[library.Node]()
				}

				return append(prev, &copy)
			}

			return prev
		}, []*library.Node{})
	}

	return dt.Reduce(hydrated, func(prev []*library.Node, curr *library.Node) []*library.Node {
		// If we're filtering for a specific node and the current iteration is
		// that node, the children are aggregated for this node regardless.
		filteredParent, ok := id.Get()
		if ok && library.NodeID(curr.Mark.ID()) == filteredParent {
			curr.Nodes = linkChildrenForParent(*curr)
			return append(prev, curr)
		}

		// If the current iteration has no parent, it's a root node. When there
		// is no filtered parent the query may contain multiple root nodes.
		_, hasParent := curr.Parent.Get()
		if !hasParent {
			curr.Nodes = linkChildrenForParent(*curr)
			return append(prev, curr)
		}

		return prev
	}, []*library.Node{})
}

// applyFilterRules applies the rather complex filtering logic for nodes in the
// tree while they are still flattened. This is because implementing this logic
// directly into the recursive query is a huge pain (especially because of Go.)
//
// This may cause a bit of over-querying as the query will, in most cases, pull
// every node (the full tree) but this can be addressed if it becomes a problem.
func applyFilterRules(f filters) func(n subtreeRow) bool {
	return func(n subtreeRow) bool {
		// If there are no visibility filters, the default is just published.
		if len(f.visibility) == 0 {
			return n.NodeVisibility == visibility.VisibilityPublished
		}

		includedInVisibilityFilter := lo.Contains(f.visibility, n.NodeVisibility)
		if !includedInVisibilityFilter {
			// The request is not interested in this node, regardless of rules.
			return false
		}

		// The default yield for this filter is to only show published nodes.
		// This state is returned after other more complex checks are done.
		isPublished := n.NodeVisibility == visibility.VisibilityPublished

		// If published and filters include publish, yield this node.
		if isPublished {
			return true
		}

		session, ok := f.requestingAccount.Get()
		if !ok {
			// If a guest is making this request, then only filter on published.
			// If the requesting guest used only other filters they see nothing.
			return isPublished
		}

		isOwner := n.NodeAccountId == xid.ID(session.ID)
		if isOwner {
			// If the requesting account owns this node, and it's within the
			// visibility filter constraint, return it in the list.
			return includedInVisibilityFilter
		}

		// The account is not the owner of the node, so we need to check if
		// they have the manage library permissions.

		isLibraryManager := session.Roles.Permissions().HasAny(rbac.PermissionManageLibrary, rbac.PermissionAdministrator)
		if !isLibraryManager {
			// If the requesting account is not the owner, and not a manager,
			// only yield the node if it's published and the filters include it.

			return n.NodeVisibility == visibility.VisibilityPublished
		}

		// the account is a library manager, but that still doesn't mean they
		// can see everything. Ensure that the only nodes not published or not
		// owned by the requesting account are in-review.

		if n.NodeVisibility == visibility.VisibilityReview {
			return true
		}

		// by this point, all logic is applied and the node is either not owned,
		// the requesting account doesn't have permission, or not in filters.
		return false
	}
}
