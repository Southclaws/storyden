package node_traversal

import (
	"context"
	"fmt"
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

const ddl = `with recursive descendants (parent, descendant, depth) as (
    select
        parent_node_id,
        id,
        1
    from
        nodes
    union
    all
    select
        d.parent,
        s.id,
        d.depth + 1
    from
        descendants d
        join nodes s on d.descendant = s.parent_node_id
)
select
    c.id                node_id,
    c.created_at        node_created_at,
    c.updated_at        node_updated_at,
    c.deleted_at        node_deleted_at,
    c.name              node_name,
    c.slug              node_slug,
    c.description       node_description,
    c.parent_node_id    node_parent_node_id,
    c.account_id        node_account_id,
    c.properties        node_properties,
	a.id                owner_id,
	a.created_at        owner_created_at,
	a.updated_at        owner_updated_at,
	a.deleted_at        owner_deleted_at,
	a.handle            owner_handle,
	a.name              owner_name,
	a.bio               owner_bio,
	a.admin             owner_admin
from
    descendants
    inner join nodes c on c.id = descendants.descendant
    inner join accounts a on a.id = c.account_id
where
    (
        (
            descendant = $1
            and parent is not null
        ) or
        parent = $1
    )
    -- additional filters
    %s
    -- end
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

func (d *database) Subtree(ctx context.Context, id datagraph.NodeID, fs ...Filter) ([]*datagraph.Node, error) {
	f := filters{}
	for _, fn := range fs {
		fn(&f)
	}

	predicates := ""
	predicateN := 2
	if f.accountSlug != nil {
		predicates = fmt.Sprintf("%s AND a.handle = $%d", predicates, predicateN)
		// predicateN++ // Do this when more predicates are added.
	}

	r, err := d.raw.QueryxContext(ctx, fmt.Sprintf(ddl, predicates), id.String())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	nodes := []*datagraph.Node{}

	for r.Next() {
		c := subtreeRow{}
		err = r.StructScan(&c)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		clus, err := fromRow(c)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		nodes = append(nodes, clus)
	}

	return nodes, nil
}

func (d *database) FilterSubtree(ctx context.Context, id datagraph.NodeID, filter string) ([]*datagraph.Node, error) {
	return nil, nil
}
