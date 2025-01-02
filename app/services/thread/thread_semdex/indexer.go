package thread_semdex

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/thread"
)

func (i *semdexer) indexThread(ctx context.Context, id post.ID) error {
	p, err := i.threadQuerier.Get(ctx, id, nil)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	updates, err := i.semdexMutator.Index(ctx, p)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if updates > 0 {
		_, err = i.threadWriter.Update(ctx, id, thread.WithIndexed())
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	return nil
}

func (i *semdexer) deindexThread(ctx context.Context, id post.ID) error {
	_, err := i.semdexMutator.Delete(ctx, xid.ID(id))
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	_, err = i.threadWriter.Update(ctx, id, thread.WithIndexed())
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
