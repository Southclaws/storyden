package cluster_children

import (
	"context"
	"fmt"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/datagraph/cluster"
	"github.com/Southclaws/storyden/internal/ent"
	cluster_model "github.com/Southclaws/storyden/internal/ent/cluster"
)

type database struct {
	db *ent.Client
	cr cluster.Repository
}

func New(db *ent.Client, cr cluster.Repository) Repository {
	return &database{db, cr}
}

type options struct {
	moveClusters bool
	moveItems    bool
}

func (d *database) Move(ctx context.Context, fromSlug datagraph.ClusterSlug, toSlug datagraph.ClusterSlug, opts ...Option) (*datagraph.Cluster, error) {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	fromCluster, err := d.cr.Get(ctx, fromSlug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	toCluster, err := d.cr.Get(ctx, toSlug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	tx, err := d.db.Tx(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = func() (err error) {
		if o.moveClusters {
			clusters, err := d.db.Cluster.Query().Where(cluster_model.ParentClusterID(xid.ID(fromCluster.ID))).All(ctx)
			if err != nil {
				return fault.Wrap(err)
			}
			childClusterIDs := dt.Map(clusters, func(c *ent.Cluster) xid.ID { return c.ID })

			err = d.db.Cluster.Update().
				SetParentID(xid.ID(toCluster.ID)).
				Where(cluster_model.IDIn(childClusterIDs...)).
				Exec(ctx)
			if err != nil {
				return fault.Wrap(err)
			}
		}
		return
	}()
	if err != nil {
		terr := tx.Rollback()
		if terr != nil {
			panic(fmt.Errorf("while handling error: %w, rollback error: %s", err, terr))
		}

		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := tx.Commit(); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return toCluster, nil
}
