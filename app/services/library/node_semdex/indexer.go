package node_semdex

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
)

func (i *semdexer) index(ctx context.Context, id library.NodeID) error {
	qk := library.NewID(xid.ID(id))

	p, err := i.nodeQuerier.Get(ctx, qk)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	err = i.indexer.Index(ctx, p)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	_, err = i.nodeWriter.Update(ctx, qk, node_writer.WithIndexed())
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (i *semdexer) deindex(ctx context.Context, id library.NodeID) error {
	qk := library.NewID(xid.ID(id))

	err := i.deleter.Delete(ctx, xid.ID(id))
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	_, err = i.nodeWriter.Update(ctx, qk, node_writer.WithIndexed())
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
