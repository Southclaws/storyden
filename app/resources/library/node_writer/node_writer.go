package node_writer

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/node"
)

type Writer struct {
	db      *ent.Client
	querier *node_querier.Querier
}

func New(db *ent.Client, querier *node_querier.Querier) *Writer {
	return &Writer{db, querier}
}

type Option func(*ent.NodeMutation)

func WithID(id library.NodeID) Option {
	return func(c *ent.NodeMutation) {
		c.SetID(xid.ID(id))
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

	return w.querier.Get(ctx, library.QueryKey{mark.NewQueryKeyID(col.ID)})
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

	qk = library.QueryKey{mark.NewQueryKeyID(pre.Mark.ID())}

	return w.querier.Get(ctx, qk)
}

func (w *Writer) Delete(ctx context.Context, qk library.QueryKey) error {
	delete := w.db.Node.Delete()

	delete.Where(qk.Predicate())

	_, err := delete.Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
