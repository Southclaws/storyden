package node_traversal

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"github.com/samber/lo"

	account_repo "github.com/Southclaws/storyden/app/resources/account"
	asset_repo "github.com/Southclaws/storyden/app/resources/asset"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/asset"
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

	nodes, err := dt.MapErr(cs, library.NodeFromModel)
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
    n.parent_node_id    node_parent_node_id,
    n.account_id        node_account_id,
    n.visibility        node_visibility,
    n.metadata        node_metadata,
    a.id                owner_id,
    a.created_at        owner_created_at,
    a.updated_at        owner_updated_at,
    a.deleted_at        owner_deleted_at,
    a.handle            owner_handle,
    a.name              owner_name,
    a.bio               owner_bio,
    a.admin             owner_admin,
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
	NodeId             xid.ID                `db:"node_id"`
	NodeCreatedAt      time.Time             `db:"node_created_at"`
	NodeUpdatedAt      time.Time             `db:"node_updated_at"`
	NodeDeletedAt      *time.Time            `db:"node_deleted_at"`
	NodeName           string                `db:"node_name"`
	NodeSlug           string                `db:"node_slug"`
	NodeParentNodeId   xid.ID                `db:"node_parent_node_id"`
	NodeParentNodeSlug xid.ID                `db:"node_parent_node_slug"`
	NodeAccountId      xid.ID                `db:"node_account_id"`
	NodeVisibility     visibility.Visibility `db:"node_visibility"`
	NodeMetadata       *[]byte               `db:"node_metadata"`
	OwnerId            xid.ID                `db:"owner_id"`
	OwnerCreatedAt     time.Time             `db:"owner_created_at"`
	OwnerUpdatedAt     time.Time             `db:"owner_updated_at"`
	OwnerDeletedAt     *time.Time            `db:"owner_deleted_at"`
	OwnerHandle        string                `db:"owner_handle"`
	OwnerName          string                `db:"owner_name"`
	OwnerBio           *string               `db:"owner_bio"`
	OwnerAdmin         bool                  `db:"owner_admin"`
	Depth              int                   `db:"depth"`
}

func fromRow(r subtreeRow) (*library.Node, error) {
	bio, err := opt.MapErr(opt.NewPtr(r.OwnerBio), datagraph.NewRichText)
	if err != nil {
		return nil, err
	}

	meta := opt.NewPtrMap(r.NodeMetadata, func(b []byte) map[string]any {
		meta := map[string]any{}
		err = json.Unmarshal(b, &meta)
		if err != nil {
			return nil
		}

		return meta
	})

	return &library.Node{
		Mark:       library.NewMark(r.NodeId, r.NodeSlug),
		CreatedAt:  r.NodeCreatedAt,
		UpdatedAt:  r.NodeUpdatedAt,
		Name:       r.NodeName,
		Visibility: r.NodeVisibility,
		Parent: opt.NewSafe(library.Node{
			Mark: library.NewMark(r.NodeParentNodeId, ""),
		}, !r.NodeParentNodeId.IsNil()),
		Owner: profile.Public{
			ID:      account_repo.AccountID(r.OwnerId),
			Created: r.OwnerCreatedAt,
			Handle:  r.OwnerHandle,
			Name:    r.OwnerName,
			Bio:     bio.OrZero(),
			// Roles not inclucded because silly flat query result...
			// to be hydrated elsewhere via second query.
		},
		Metadata: meta.OrZero(),
	}, nil
}

func (d *database) Subtree(ctx context.Context, id opt.Optional[library.NodeID], fs ...Filter) ([]*library.Node, error) {
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

	flat := []*library.Node{}

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

	filtered := dt.Filter(flat, applyFilterRules(f))

	ids := dt.Map(filtered, func(n *library.Node) xid.ID { return xid.ID(n.Mark.ID()) })

	// TODO: Build a table of pointers to look up each asset via node ID

	relatedAssets := d.db.Asset.Query().
		Where(asset.HasNodesWith(node.IDIn(ids...))).
		WithNodes().
		AllX(ctx)

	hydrated := dt.Map(filtered, func(n *library.Node) *library.Node {
		// NOTE: This is slow as fuck (2 nested loops lol) needs the
		// aforementioned hash table lookup for node <> asset relations.
		assets := dt.Filter(relatedAssets, func(a *ent.Asset) bool {
			_, found := lo.Find(a.Edges.Nodes, func(an *ent.Node) bool {
				return n.GetID() == an.ID
			})

			return found
		})

		n.Assets = dt.Map(assets, asset_repo.FromModel)

		return n
	})

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

	// Rebuild the flat list into the tree
	nodes := dt.Reduce(hydrated, func(prev []*library.Node, curr *library.Node) []*library.Node {
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

	return nodes, nil
}

func (d *database) FilterSubtree(ctx context.Context, id library.NodeID, filter string) ([]*library.Node, error) {
	return nil, nil
}

// applyFilterRules applies the rather complex filtering logic for nodes in the
// tree while they are still flattened. This is because implementing this logic
// directly into the recursive query is a huge pain (especially because of Go.)
//
// This may cause a bit of over-querying as the query will, in most cases, pull
// every node (the full tree) but this can be addressed if it becomes a problem.
func applyFilterRules(f filters) func(n *library.Node) bool {
	return func(n *library.Node) bool {
		// If there are no visibility filters, the default is just published.
		if len(f.visibility) == 0 {
			return n.Visibility == visibility.VisibilityPublished
		}

		includedInVisibilityFilter := lo.Contains(f.visibility, n.Visibility)
		if !includedInVisibilityFilter {
			// The request is not interested in this node, regardless of rules.
			return false
		}

		// The default yield for this filter is to only show published nodes.
		// This state is returned after other more complex checks are done.
		isPublished := n.Visibility == visibility.VisibilityPublished

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

		isOwner := n.Owner.ID == session.ID
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

			return n.Visibility == visibility.VisibilityPublished
		}

		// the account is a library manager, but that still doesn't mean they
		// can see everything. Ensure that the only nodes not published or not
		// owned by the requesting account are in-review.

		if n.Visibility == visibility.VisibilityReview {
			return true
		}

		// by this point, all logic is applied and the node is either not owned,
		// the requesting account doesn't have permission, or not in filters.
		return false
	}
}
