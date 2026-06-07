package node_version_writer

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_version"
	"github.com/Southclaws/storyden/app/resources/library/node_version/node_version_querier"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/nodeversion"
)

type Writer struct {
	db      *ent.Client
	querier *node_version_querier.Querier
}

func New(db *ent.Client, querier *node_version_querier.Querier) *Writer {
	return &Writer{db: db, querier: querier}
}

type Option func(*ent.NodeVersionMutation)

func WithName(v string) Option {
	return func(m *ent.NodeVersionMutation) { m.SetName(v) }
}

func WithSlug(v string) Option {
	return func(m *ent.NodeVersionMutation) { m.SetSlug(v) }
}

func WithDescription(v opt.Optional[string]) Option {
	return func(m *ent.NodeVersionMutation) {
		if s, ok := v.Get(); ok {
			m.SetDescription(s)
		} else {
			m.ClearDescription()
		}
	}
}

func WithContent(v opt.Optional[datagraph.Content]) Option {
	return func(m *ent.NodeVersionMutation) {
		if s, ok := v.Get(); ok {
			m.SetContent(s.HTML())
		} else {
			m.ClearContent()
		}
	}
}

func WithMetadata(v map[string]any) Option {
	return func(m *ent.NodeVersionMutation) { m.SetMetadata(v) }
}

func WithPropertiesSnapshot(v []node_version.PropertySnapshot) Option {
	return func(m *ent.NodeVersionMutation) {
		m.SetPropertiesSnapshot(node_version.UnmapPropertySnapshots(v))
	}
}

func WithStatus(v node_version.VersionStatus) Option {
	return func(m *ent.NodeVersionMutation) {
		m.SetStatus(nodeversion.Status(v.String()))
	}
}

func (w *Writer) Create(
	ctx context.Context,
	qk library.QueryKey,
	authorID account.AccountID,
	filter node_version_querier.NodeFilter,
	opts ...Option,
) (*node_version.NodeVersion, error) {
	nodeID, err := w.db.Node.Query().
		Where(filter.Predicates(qk)...).
		OnlyID(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	create := w.db.NodeVersion.Create()
	m := create.Mutation()

	m.SetNodeID(nodeID)
	m.SetAuthorID(xid.ID(authorID))

	for _, fn := range opts {
		fn(m)
	}

	v, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return w.querier.Get(ctx, node_version.VersionID(v.ID))
}

func (w *Writer) Update(
	ctx context.Context,
	id node_version.VersionID,
	opts ...Option,
) (*node_version.NodeVersion, error) {
	update := w.db.NodeVersion.UpdateOneID(xid.ID(id))
	m := update.Mutation()

	for _, fn := range opts {
		fn(m)
	}

	err := update.Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return w.querier.Get(ctx, id)
}

func (w *Writer) Delete(ctx context.Context, id node_version.VersionID) error {
	err := w.db.NodeVersion.DeleteOneID(xid.ID(id)).Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	return nil
}
