package node_version_querier

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/role/role_hydrate"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_version"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/node"
	"github.com/Southclaws/storyden/internal/ent/nodeversion"
	"github.com/Southclaws/storyden/internal/ent/predicate"
)

type Querier struct {
	db          *ent.Client
	roleQuerier *role_hydrate.Hydrator
}

func New(db *ent.Client, roleQuerier *role_hydrate.Hydrator) *Querier {
	return &Querier{db: db, roleQuerier: roleQuerier}
}

type NodeFilter struct {
	AccountID opt.Optional[account.AccountID]
	CanManage bool
}

func (f NodeFilter) Predicates(qk library.QueryKey) []predicate.Node {
	predicates := []predicate.Node{qk.Predicate()}

	if accountID, ok := f.AccountID.Get(); ok {
		if f.CanManage {
			predicates = append(predicates, node.Or(
				node.AccountID(xid.ID(accountID)),
				node.VisibilityIn(node.VisibilityPublished, node.VisibilityReview),
			))
		} else {
			predicates = append(predicates, node.Or(
				node.AccountID(xid.ID(accountID)),
				node.VisibilityEQ(node.VisibilityPublished),
			))
		}
	} else {
		predicates = append(predicates, node.VisibilityEQ(node.VisibilityPublished))
	}

	return predicates
}

func (q *Querier) Get(ctx context.Context, id node_version.VersionID) (*node_version.NodeVersion, error) {
	v, err := q.db.NodeVersion.Query().
		Where(nodeversion.IDEQ(xid.ID(id))).
		WithAuthor().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := q.hydrateAuthorRoles(ctx, v); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return node_version.Map(v)
}

func (q *Querier) GetForNode(
	ctx context.Context,
	qk library.QueryKey,
	id node_version.VersionID,
	filter NodeFilter,
) (*node_version.NodeVersion, error) {
	v, err := q.db.NodeVersion.Query().
		Where(
			nodeversion.IDEQ(xid.ID(id)),
			nodeversion.HasNodeWith(filter.Predicates(qk)...),
		).
		WithAuthor().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := q.hydrateAuthorRoles(ctx, v); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return node_version.Map(v)
}

func (q *Querier) GetPreviousReference(
	ctx context.Context,
	v *node_version.NodeVersion,
) (opt.Optional[node_version.VersionReference], error) {
	if v.Status != node_version.VersionStatusApplied {
		return opt.NewEmpty[node_version.VersionReference](), nil
	}

	previous, err := q.db.NodeVersion.Query().
		Where(
			nodeversion.NodeIDEQ(xid.ID(v.NodeID)),
			nodeversion.StatusEQ(nodeversion.StatusApplied),
			nodeversion.Or(
				nodeversion.UpdatedAtLT(v.UpdatedAt),
				nodeversion.And(
					nodeversion.UpdatedAtEQ(v.UpdatedAt),
					nodeversion.IDLT(xid.ID(v.ID)),
				),
			),
		).
		Select(
			nodeversion.FieldID,
			nodeversion.FieldCreatedAt,
			nodeversion.FieldUpdatedAt,
			nodeversion.FieldAuthorID,
			nodeversion.FieldStatus,
		).
		Order(ent.Desc(nodeversion.FieldUpdatedAt), ent.Desc(nodeversion.FieldID)).
		WithAuthor().
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return opt.NewEmpty[node_version.VersionReference](), nil
		}
		return opt.NewEmpty[node_version.VersionReference](), fault.Wrap(err, fctx.With(ctx))
	}

	if err := q.hydrateAuthorRoles(ctx, previous); err != nil {
		return opt.NewEmpty[node_version.VersionReference](), fault.Wrap(err, fctx.With(ctx))
	}

	mapped, err := node_version.MapReference(previous)
	if err != nil {
		return opt.NewEmpty[node_version.VersionReference](), fault.Wrap(err, fctx.With(ctx))
	}

	return opt.New(*mapped), nil
}

func (q *Querier) List(
	ctx context.Context,
	qk library.QueryKey,
	statuses []node_version.VersionStatus,
	filter NodeFilter,
	page pagination.Parameters,
) (pagination.Result[*node_version.NodeVersion], error) {
	statusVals := dt.Map(statuses, func(s node_version.VersionStatus) nodeversion.Status {
		return nodeversion.Status(s.String())
	})

	query := q.db.NodeVersion.Query().
		Where(nodeversion.HasNodeWith(filter.Predicates(qk)...)).
		Order(ent.Desc(nodeversion.FieldUpdatedAt), ent.Desc(nodeversion.FieldID)).
		WithAuthor()

	if len(statusVals) > 0 {
		query.Where(nodeversion.StatusIn(statusVals...))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return pagination.Result[*node_version.NodeVersion]{}, fault.Wrap(err, fctx.With(ctx))
	}

	rows, err := query.
		Limit(page.Limit()).
		Offset(page.Offset()).
		All(ctx)
	if err != nil {
		return pagination.Result[*node_version.NodeVersion]{}, fault.Wrap(err, fctx.With(ctx))
	}

	if err := q.hydrateAuthorsRoles(ctx, rows...); err != nil {
		return pagination.Result[*node_version.NodeVersion]{}, fault.Wrap(err, fctx.With(ctx))
	}

	versions, err := dt.MapErr(rows, node_version.Map)
	if err != nil {
		return pagination.Result[*node_version.NodeVersion]{}, fault.Wrap(err, fctx.With(ctx))
	}

	return pagination.NewPageResult(page, total, versions), nil
}

func (q *Querier) ListVisible(
	ctx context.Context,
	qk library.QueryKey,
	filter NodeFilter,
	page pagination.Parameters,
) (pagination.Result[*node_version.NodeVersion], error) {
	query := q.db.NodeVersion.Query().
		Where(nodeversion.HasNodeWith(filter.Predicates(qk)...)).
		Order(ent.Desc(nodeversion.FieldUpdatedAt), ent.Desc(nodeversion.FieldID)).
		WithAuthor()

	if filter.CanManage {
		query.Where(nodeversion.StatusIn(nodeversion.StatusApplied, nodeversion.StatusDraft))
	} else if accountID, ok := filter.AccountID.Get(); ok {
		query.Where(nodeversion.Or(
			nodeversion.StatusEQ(nodeversion.StatusApplied),
			nodeversion.And(
				nodeversion.StatusEQ(nodeversion.StatusDraft),
				nodeversion.AuthorIDEQ(xid.ID(accountID)),
			),
		))
	} else {
		query.Where(nodeversion.StatusEQ(nodeversion.StatusApplied))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return pagination.Result[*node_version.NodeVersion]{}, fault.Wrap(err, fctx.With(ctx))
	}

	rows, err := query.
		Limit(page.Limit()).
		Offset(page.Offset()).
		All(ctx)
	if err != nil {
		return pagination.Result[*node_version.NodeVersion]{}, fault.Wrap(err, fctx.With(ctx))
	}

	if err := q.hydrateAuthorsRoles(ctx, rows...); err != nil {
		return pagination.Result[*node_version.NodeVersion]{}, fault.Wrap(err, fctx.With(ctx))
	}

	versions, err := dt.MapErr(rows, node_version.Map)
	if err != nil {
		return pagination.Result[*node_version.NodeVersion]{}, fault.Wrap(err, fctx.With(ctx))
	}

	return pagination.NewPageResult(page, total, versions), nil
}

func (q *Querier) ListAllDrafts(
	ctx context.Context,
	filter NodeFilter,
	page pagination.Parameters,
) (pagination.Result[*node_version.NodeVersionWithNode], error) {
	query := q.db.NodeVersion.Query().
		Where(nodeversion.StatusEQ(nodeversion.StatusDraft)).
		Order(ent.Desc(nodeversion.FieldUpdatedAt), ent.Desc(nodeversion.FieldID)).
		WithAuthor().
		WithNode(func(nq *ent.NodeQuery) {
			nq.WithOwner()
		})

	if filter.CanManage {
		// Managers can see all drafts
	} else if accountID, ok := filter.AccountID.Get(); ok {
		// Authors can only see their own drafts
		query.Where(nodeversion.AuthorIDEQ(xid.ID(accountID)))
	} else {
		// Unauthenticated users see no drafts
		query.Where(nodeversion.IDEQ(xid.NilID()))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return pagination.Result[*node_version.NodeVersionWithNode]{}, fault.Wrap(err, fctx.With(ctx))
	}

	rows, err := query.
		Limit(page.Limit()).
		Offset(page.Offset()).
		All(ctx)
	if err != nil {
		return pagination.Result[*node_version.NodeVersionWithNode]{}, fault.Wrap(err, fctx.With(ctx))
	}

	// Build list of all accounts that need role hydration (version authors + node owners)
	roleTargets := make([]*ent.Account, 0, len(rows)*2)
	for _, v := range rows {
		// Add version author
		if author := v.Edges.Author; author != nil {
			roleTargets = append(roleTargets, author)
		}

		// Add node owner
		if node := v.Edges.Node; node != nil {
			if owner := node.Edges.Owner; owner != nil {
				roleTargets = append(roleTargets, owner)
			}
		}
	}

	// Hydrate all roles in one batch using cached role data
	if err := q.roleQuerier.HydrateRoleEdges(ctx, roleTargets...); err != nil {
		return pagination.Result[*node_version.NodeVersionWithNode]{}, fault.Wrap(err, fctx.With(ctx))
	}

	drafts, err := dt.MapErr(rows, func(v *ent.NodeVersion) (*node_version.NodeVersionWithNode, error) {
		version, err := node_version.Map(v)
		if err != nil {
			return nil, fault.Wrap(err)
		}

		nodeEdge, err := v.Edges.NodeOrErr()
		if err != nil {
			return nil, fault.Wrap(err)
		}

		node, err := library.MapNode(true, nil)(nodeEdge)
		if err != nil {
			return nil, fault.Wrap(err)
		}

		return &node_version.NodeVersionWithNode{
			NodeVersion: *version,
			Node:        node,
		}, nil
	})
	if err != nil {
		return pagination.Result[*node_version.NodeVersionWithNode]{}, fault.Wrap(err, fctx.With(ctx))
	}

	return pagination.NewPageResult(page, total, drafts), nil
}

func (q *Querier) GetDraft(
	ctx context.Context,
	qk library.QueryKey,
	authorID account.AccountID,
	filter NodeFilter,
) (opt.Optional[node_version.NodeVersion], error) {
	v, err := q.db.NodeVersion.Query().
		Where(
			nodeversion.HasNodeWith(filter.Predicates(qk)...),
			nodeversion.AuthorIDEQ(xid.ID(authorID)),
			nodeversion.StatusEQ(nodeversion.StatusDraft),
		).
		Order(ent.Desc(nodeversion.FieldUpdatedAt)).
		WithAuthor().
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return opt.NewEmpty[node_version.NodeVersion](), nil
		}
		return opt.NewEmpty[node_version.NodeVersion](), fault.Wrap(err, fctx.With(ctx))
	}

	if err := q.hydrateAuthorRoles(ctx, v); err != nil {
		return opt.NewEmpty[node_version.NodeVersion](), fault.Wrap(err, fctx.With(ctx))
	}

	mapped, err := node_version.Map(v)
	if err != nil {
		return opt.NewEmpty[node_version.NodeVersion](), fault.Wrap(err, fctx.With(ctx))
	}

	return opt.New(*mapped), nil
}

func (q *Querier) GetNodeDraft(
	ctx context.Context,
	qk library.QueryKey,
	filter NodeFilter,
) (opt.Optional[node_version.NodeVersion], error) {
	v, err := q.db.NodeVersion.Query().
		Where(
			nodeversion.HasNodeWith(filter.Predicates(qk)...),
			nodeversion.StatusEQ(nodeversion.StatusDraft),
		).
		Order(ent.Desc(nodeversion.FieldUpdatedAt)).
		WithAuthor().
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return opt.NewEmpty[node_version.NodeVersion](), nil
		}
		return opt.NewEmpty[node_version.NodeVersion](), fault.Wrap(err, fctx.With(ctx))
	}

	if err := q.hydrateAuthorRoles(ctx, v); err != nil {
		return opt.NewEmpty[node_version.NodeVersion](), fault.Wrap(err, fctx.With(ctx))
	}

	mapped, err := node_version.Map(v)
	if err != nil {
		return opt.NewEmpty[node_version.NodeVersion](), fault.Wrap(err, fctx.With(ctx))
	}

	return opt.New(*mapped), nil
}

func (q *Querier) HasNodeDraft(
	ctx context.Context,
	qk library.QueryKey,
) (bool, error) {
	exists, err := q.db.NodeVersion.Query().
		Where(
			nodeversion.HasNodeWith(qk.Predicate()),
			nodeversion.StatusEQ(nodeversion.StatusDraft),
		).
		Exist(ctx)
	if err != nil {
		return false, fault.Wrap(err, fctx.With(ctx))
	}

	return exists, nil
}

func (q *Querier) hydrateAuthorRoles(ctx context.Context, v *ent.NodeVersion) error {
	author, err := v.Edges.AuthorOrErr()
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return q.roleQuerier.HydrateRoleEdges(ctx, author)
}

func (q *Querier) hydrateAuthorsRoles(ctx context.Context, versions ...*ent.NodeVersion) error {
	authors := make([]*ent.Account, 0, len(versions))

	for _, v := range versions {
		author, err := v.Edges.AuthorOrErr()
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}

		authors = append(authors, author)
	}

	return q.roleQuerier.HydrateRoleEdges(ctx, authors...)
}
