package cluster_traversal

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

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/cluster"
)

type database struct {
	db  *ent.Client
	raw *sqlx.DB
}

func New(db *ent.Client, raw *sqlx.DB) Repository {
	return &database{db, raw}
}

func (d *database) Root(ctx context.Context) ([]*datagraph.Cluster, error) {
	query := d.db.Cluster.Query().Where(cluster.ParentClusterIDIsNil()).WithOwner()

	cs, err := query.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	clusters, err := dt.MapErr(cs, datagraph.ClusterFromModel)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return clusters, nil
}

const ddl = `with recursive descendants (parent, descendant, depth) as (
    select
        parent_cluster_id,
        id,
        1
    from
        clusters
    union
    all
    select
        d.parent,
        s.id,
        d.depth + 1
    from
        descendants d
        join clusters s on d.descendant = s.parent_cluster_id
)
select
    c.id                cluster_id,
    c.created_at        cluster_created_at,
    c.updated_at        cluster_updated_at,
    c.deleted_at        cluster_deleted_at,
    c.name              cluster_name,
    c.slug              cluster_slug,
    c.image_url         cluster_image_url,
    c.description       cluster_description,
    c.parent_cluster_id cluster_parent_cluster_id,
    c.account_id        cluster_account_id,
    c.properties        cluster_properties,
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
    inner join clusters c on c.id = descendants.descendant
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
	ClusterId              xid.ID     `db:"cluster_id"`
	ClusterCreatedAt       time.Time  `db:"cluster_created_at"`
	ClusterUpdatedAt       time.Time  `db:"cluster_updated_at"`
	ClusterDeletedAt       *time.Time `db:"cluster_deleted_at"`
	ClusterName            string     `db:"cluster_name"`
	ClusterSlug            string     `db:"cluster_slug"`
	ClusterImageUrl        *string    `db:"cluster_image_url"`
	ClusterDescription     string     `db:"cluster_description"`
	ClusterParentClusterId xid.ID     `db:"cluster_parent_cluster_id"`
	ClusterAccountId       xid.ID     `db:"cluster_account_id"`
	ClusterProperties      any        `db:"cluster_properties"`
	OwnerId                xid.ID     `db:"owner_id"`
	OwnerCreatedAt         time.Time  `db:"owner_created_at"`
	OwnerUpdatedAt         time.Time  `db:"owner_updated_at"`
	OwnerDeletedAt         *time.Time `db:"owner_deleted_at"`
	OwnerHandle            string     `db:"owner_handle"`
	OwnerName              string     `db:"owner_name"`
	OwnerBio               *string    `db:"owner_bio"`
	OwnerAdmin             bool       `db:"owner_admin"`
}

func fromRow(r subtreeRow) (*datagraph.Cluster, error) {
	return &datagraph.Cluster{
		ID:          datagraph.ClusterID(r.ClusterId),
		CreatedAt:   r.ClusterCreatedAt,
		UpdatedAt:   r.ClusterUpdatedAt,
		Name:        r.ClusterName,
		Slug:        r.ClusterSlug,
		ImageURL:    opt.NewPtr(r.ClusterImageUrl),
		Description: r.ClusterDescription,
		Owner: profile.Profile{
			ID:     account.AccountID(r.OwnerId),
			Handle: r.OwnerHandle,
			Name:   r.OwnerName,
			Bio:    opt.NewPtr(r.OwnerBio).OrZero(),
			Admin:  r.OwnerAdmin,
		},
		Properties: r.ClusterProperties,
	}, nil
}

func (d *database) Subtree(ctx context.Context, id datagraph.ClusterID) ([]*datagraph.Cluster, error) {
	filters := ""
	r, err := d.raw.QueryxContext(ctx, fmt.Sprintf(ddl, filters), id.String())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	clusters := []*datagraph.Cluster{}

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

		clusters = append(clusters, clus)
	}

	return clusters, nil
}

func (d *database) FilterSubtree(ctx context.Context, id datagraph.ClusterID, filter string) ([]*datagraph.Cluster, error) {
	return nil, nil
}
