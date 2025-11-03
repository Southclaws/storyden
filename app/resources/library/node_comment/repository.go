package node_comment

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/postnode"
)

type Repository struct {
	db *ent.Client
}

func New(db *ent.Client) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetThreadIDs(
	ctx context.Context,
	qk library.QueryKey,
	pp pagination.Parameters,
) (*pagination.Result[xid.ID], error) {
	n, err := r.db.Node.Query().
		Where(qk.Predicate()).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	query := r.db.PostNode.Query().
		Where(postnode.NodeIDEQ(n.ID)).
		Order(ent.Desc(postnode.FieldCreatedAt))

	total, err := query.Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	query.
		Limit(pp.Limit()).
		Offset(pp.Offset())

	results, err := query.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	postIDs := make([]xid.ID, 0, len(results))
	for _, pn := range results {
		postIDs = append(postIDs, pn.PostID)
	}

	result := pagination.NewPageResult(pp, total, postIDs)

	return &result, nil
}

func (r *Repository) Create(
	ctx context.Context,
	id library.NodeID,
	threadID post.ID,
) error {
	err := r.db.PostNode.Create().
		SetNodeID(xid.ID(id)).
		SetPostID(xid.ID(threadID)).
		Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		if ent.IsConstraintError(err) {
			return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.AlreadyExists))
		}
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return nil
}
