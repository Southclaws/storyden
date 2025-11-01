package node_writer

import (
	"context"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/lexorank"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_children"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/node"
	"github.com/Southclaws/storyden/internal/ent/propertyschema"
)

type Writer struct {
	db          *ent.Client
	querier     *node_querier.Querier
	childWriter *node_children.Writer
}

func New(db *ent.Client, querier *node_querier.Querier, childWriter *node_children.Writer) *Writer {
	return &Writer{
		db:          db,
		querier:     querier,
		childWriter: childWriter,
	}
}

type Option func(*ent.NodeMutation)

func WithID(id library.NodeID) Option {
	return func(c *ent.NodeMutation) {
		c.SetID(xid.ID(id))
	}
}

func WithIndexed() Option {
	return func(nm *ent.NodeMutation) {
		nm.SetIndexedAt(time.Now())
	}
}

func WithName(v string) Option {
	return func(c *ent.NodeMutation) {
		c.SetName(v)
	}
}

func WithSlug(v string) Option {
	return func(c *ent.NodeMutation) {
		c.SetSlug(v)
	}
}

func WithAssets(a []asset.AssetID) Option {
	return func(m *ent.NodeMutation) {
		m.AddAssetIDs(a...)
	}
}

func WithAssetsRemoved(a []asset.AssetID) Option {
	return func(m *ent.NodeMutation) {
		m.RemoveAssetIDs(a...)
	}
}

func WithLink(id xid.ID) Option {
	return func(pm *ent.NodeMutation) {
		pm.SetLinkID(id)
	}
}

func WithLinkRemove() Option {
	return func(pm *ent.NodeMutation) {
		pm.ClearLink()
	}
}

func WithContentLinks(ids ...xid.ID) Option {
	return func(pm *ent.NodeMutation) {
		pm.AddContentLinkIDs(ids...)
	}
}

func WithPrimaryImage(id asset.AssetID) Option {
	return func(nm *ent.NodeMutation) {
		nm.SetPrimaryAssetID(id)
	}
}

func WithPrimaryImageRemoved() Option {
	return func(nm *ent.NodeMutation) {
		nm.ClearPrimaryAssetID()
	}
}

func WithContent(v datagraph.Content) Option {
	return func(c *ent.NodeMutation) {
		c.SetContent(v.HTML())
		c.SetDescription(v.Short())
	}
}

func WithDescription(v string) Option {
	return func(c *ent.NodeMutation) {
		c.SetDescription(v)
	}
}

func WithParent(v library.NodeID) Option {
	return func(c *ent.NodeMutation) {
		c.SetParentID(xid.ID(v))
	}
}

func WithHideChildren(v bool) Option {
	return func(c *ent.NodeMutation) {
		c.SetHideChildTree(v)
	}
}

func WithVisibility(v visibility.Visibility) Option {
	return func(c *ent.NodeMutation) {
		c.SetVisibility(node.Visibility(v.String()))
	}
}

func WithMetadata(v map[string]any) Option {
	return func(c *ent.NodeMutation) {
		c.SetMetadata(v)
	}
}

func WithChildNodeAdd(id xid.ID) Option {
	return func(c *ent.NodeMutation) {
		c.AddNodeIDs(id)
	}
}

func WithChildNodeRemove(id xid.ID) Option {
	return func(c *ent.NodeMutation) {
		c.RemoveNodeIDs(id)
	}
}

func WithTagsAdd(refs ...tag_ref.ID) Option {
	ids := dt.Map(refs, func(i tag_ref.ID) xid.ID { return xid.ID(i) })
	return func(c *ent.NodeMutation) {
		c.AddTagIDs(ids...)
	}
}

func WithTagsRemove(refs ...tag_ref.ID) Option {
	ids := dt.Map(refs, func(i tag_ref.ID) xid.ID { return xid.ID(i) })
	return func(c *ent.NodeMutation) {
		c.RemoveTagIDs(ids...)
	}
}

func (w *Writer) Create(
	ctx context.Context,
	owner account.AccountID,
	name string,
	slug mark.Slug,
	opts ...Option,
) (*library.Node, error) {
	// TODO: Use a Node Mark for this.
	if slug.String() == "" {
		return nil, fault.New("slug cannot be empty", fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	create := w.db.Node.Create()
	mutate := create.Mutation()

	mutate.SetOwnerID(xid.ID(owner))
	mutate.SetName(name)
	mutate.SetSlug(slug.String())

	for _, fn := range opts {
		fn(mutate)
	}

	parent := opt.NewSafe(mutate.ParentID())
	sortkey, err := w.getNextSortKey(ctx, parent)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	mutate.SetSort(*sortkey)

	col, err := create.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err,
				fctx.With(ctx),
				ftag.With(ftag.AlreadyExists),
				fmsg.WithDesc("already exists", "The node URL slug must be unique and the specified slug is already in use."),
			)
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return w.querier.Get(ctx, library.NewID(col.ID))
}

func (w *Writer) getNextSortKey(ctx context.Context, parent opt.Optional[xid.ID]) (*lexorank.Key, error) {
	siblingQuery := w.db.Node.Query().
		Select(node.FieldSort).
		Limit(1).
		Order(ent.Desc(node.FieldSort))

	// If the parent is not nil, we need to filter by the parent ID. Otherwise
	// the target node is at the root level and its siblings are too.
	if parentID, ok := parent.Get(); ok {
		siblingQuery.Where(node.ParentNodeID(parentID))
	} else {
		siblingQuery.Where(node.ParentNodeIDIsNil())
	}

	var lerr error

	for range 2 {
		siblings, err := siblingQuery.All(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if len(siblings) == 1 {
			sibling := siblings[0]

			sortkey, ok := sibling.Sort.After(100)
			if !ok {
				err := w.childWriter.Normalise(ctx, parent.Ptr())
				if err != nil {
					return nil, fault.Wrap(err, fctx.With(ctx))
				}

				lerr = fault.Newf("failed to get next sort key between %s and %s", lexorank.Top, sibling.Sort)
				continue
			}

			return sortkey, nil
		} else {
			return &lexorank.Middle, nil
		}
	}
	if lerr == nil {
		lerr = fault.New("failed to get next sort key: unknown")
	}

	// TODO: Explore if returning a random sortkey instead of an error is good.
	return nil, fault.Wrap(lerr, fctx.With(ctx))
}

func (w *Writer) Update(ctx context.Context, qk library.QueryKey, opts ...Option) (*library.Node, error) {
	// NOTE: Should be a probe not a full read. Query is necessary because of
	// the Mark being used (id or slug) for updates. Cannot use UpdateOneID.
	pre, err := w.querier.Get(ctx, qk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	update := w.db.Node.Update()
	update.Where(qk.Predicate())

	mutate := update.Mutation()
	for _, fn := range opts {
		fn(mutate)
	}

	err = update.Exec(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	qk = library.NewID(pre.Mark.ID())

	return w.querier.Get(ctx, qk)
}

func (w *Writer) Delete(ctx context.Context, qk library.QueryKey) error {
	delete := w.db.Node.Delete()

	delete.Where(qk.Predicate())

	_, err := delete.Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	// NOTE: This should probably be run separately either as a background job
	// or in parallel. However, running in a goroutine here for some reason does
	// not delete anything, presumably because SQLite has not committed deletion
	// performed above to disk and is running in a mode that prevents this. It
	// may work fine on Postgres or CockroachDB though but for now this is fine.
	w.CleanupOrphanedSchemas(ctx)

	return nil
}

func (w *Writer) CleanupOrphanedSchemas(ctx context.Context) {
	// error handling doesn't matter this is run in parallel and doesn't matter.
	w.db.PropertySchema.
		Delete().
		Where(
			propertyschema.Not(
				propertyschema.HasNode(),
			),
		).
		Exec(ctx)
}
