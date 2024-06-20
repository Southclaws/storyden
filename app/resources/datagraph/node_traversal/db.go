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

	account_repo "github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/node"
)

type database struct {
	db  *ent.Client
	raw *sqlx.DB
}

func New(db *ent.Client, raw *sqlx.DB) Repository {
	return &database{db, raw}
}

func (d *database) Root(ctx context.Context, fs ...Filter) ([]*datagraph.Node, error) {
	query := d.db.Node.Query().
		Where(node.ParentNodeIDIsNil()).
		WithOwner().
		WithAssets()

	f := filters{}
	for _, fn := range fs {
		fn(&f)
	}

	if f.accountSlug != nil {
		query.Where(node.HasOwnerWith(account.Handle(*f.accountSlug)))
	}

	if len(f.visibility) > 0 {
		visibilityTypes := dt.Map(f.visibility, func(v post.Visibility) node.Visibility {
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

	nodes, err := dt.MapErr(cs, datagraph.NodeFromModel)
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
    n.created_at        node_created_at,
    n.updated_at        node_updated_at,
    n.deleted_at        node_deleted_at,
    n.name              node_name,
    n.slug              node_slug,
    n.description       node_description,
    n.parent_node_id    node_parent_node_id,
    n.account_id        node_account_id,
    n.properties        node_properties,
    a.id                owner_id,
    a.created_at        owner_created_at,
    a.updated_at        owner_updated_at,
    a.deleted_at        owner_deleted_at,
    a.handle            owner_handle,
    a.name              owner_name,
    a.bio               owner_bio,
    a.admin             owner_admin
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
	NodeId           xid.ID     `db:"node_id"`
	NodeCreatedAt    time.Time  `db:"node_created_at"`
	NodeUpdatedAt    time.Time  `db:"node_updated_at"`
	NodeDeletedAt    *time.Time `db:"node_deleted_at"`
	NodeName         string     `db:"node_name"`
	NodeSlug         string     `db:"node_slug"`
	NodeDescription  string     `db:"node_description"`
	NodeParentNodeId xid.ID     `db:"node_parent_node_id"`
	NodeAccountId    xid.ID     `db:"node_account_id"`
	NodeProperties   any        `db:"node_properties"`
	OwnerId          xid.ID     `db:"owner_id"`
	OwnerCreatedAt   time.Time  `db:"owner_created_at"`
	OwnerUpdatedAt   time.Time  `db:"owner_updated_at"`
	OwnerDeletedAt   *time.Time `db:"owner_deleted_at"`
	OwnerHandle      string     `db:"owner_handle"`
	OwnerName        string     `db:"owner_name"`
	OwnerBio         *string    `db:"owner_bio"`
	OwnerAdmin       bool       `db:"owner_admin"`
}

func fromRow(r subtreeRow) (*datagraph.Node, error) {
	return &datagraph.Node{
		ID:          datagraph.NodeID(r.NodeId),
		CreatedAt:   r.NodeCreatedAt,
		UpdatedAt:   r.NodeUpdatedAt,
		Name:        r.NodeName,
		Slug:        r.NodeSlug,
		Description: r.NodeDescription,
		Parent: opt.NewSafe(datagraph.Node{
			ID: datagraph.NodeID(r.NodeParentNodeId),
		}, !r.NodeParentNodeId.IsNil()),
		Owner: profile.Profile{
			ID:      account_repo.AccountID(r.OwnerId),
			Created: r.OwnerCreatedAt,
			Handle:  r.OwnerHandle,
			Name:    r.OwnerName,
			Bio:     opt.NewPtr(r.OwnerBio).OrZero(),
			Admin:   r.OwnerAdmin,
		},
		Properties: r.NodeProperties,
	}, nil
}

func (d *database) Subtree(ctx context.Context, id opt.Optional[datagraph.NodeID], fs ...Filter) ([]*datagraph.Node, error) {
	f := filters{}
	for _, fn := range fs {
		fn(&f)
	}

	// NOTE: i fucking hate writing raw sql into source code...

	var rootPredicate string
	predicates := []string{}
	args := []interface{}{}
	argOffset := 1

	if parentNodeID, ok := id.Get(); ok {
		args = append(args, parentNodeID.String())
		rootPredicate = "id is $1"
	} else {
		rootPredicate = "parent_node_id is null"
	}

	if f.accountSlug != nil {
		aidx := len(args) + argOffset
		predicates = append(predicates, fmt.Sprintf(
			"a.handle = $%d",
			aidx))

		args = append(args, *f.accountSlug)
	}

	if f.depth != nil {
		aidx := len(args) + argOffset
		predicates = append(predicates, fmt.Sprintf(
			"depth <= $%d",
			aidx))

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

	flat := []*datagraph.Node{}

	for r.Next() {
		c := subtreeRow{}
		err = r.StructScan(&c)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		n, err := fromRow(c)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		flat = append(flat, n)
	}

	var linkChildrenForParent func(datagraph.Node) []*datagraph.Node

	linkChildrenForParent = func(parent datagraph.Node) []*datagraph.Node {
		filteredParent, isFilteringParent := id.Get()

		return dt.Reduce(filtered, func(prev []*datagraph.Node, curr *datagraph.Node) []*datagraph.Node {
			if p, ok := curr.Parent.Get(); ok && p.ID == parent.ID {
				// Take a copy because our mutations cannot apply to `flat`.
				copy := *curr

				copy.Nodes = linkChildrenForParent(copy)

				// If the current iteration is not the root node of a parent
				// node (a subtree query) then blank out the parent field since
				// it's a waste to store this information in tree children.
				if isFilteringParent && filteredParent != copy.ID {
					copy.Parent = opt.NewEmpty[datagraph.Node]()
				}

				return append(prev, &copy)
			}

			return prev
		}, []*datagraph.Node{})

	}

	// Rebuild the flat list into the tree
	nodes := dt.Reduce(filtered, func(prev []*datagraph.Node, curr *datagraph.Node) []*datagraph.Node {

		// If we're filtering for a specific node and the current iteration is
		// that node, the children are aggregated for this node regardless.
		filteredParent, ok := id.Get()
		if ok && curr.ID == filteredParent {
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
	}, []*datagraph.Node{})

	return nodes, nil
}

func (d *database) FilterSubtree(ctx context.Context, id datagraph.NodeID, filter string) ([]*datagraph.Node, error) {
	return nil, nil
}
