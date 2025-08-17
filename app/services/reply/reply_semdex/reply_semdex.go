package reply_semdex

import (
	"context"
	"log/slog"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func Build() fx.Option {
	return fx.Options(
		fx.Invoke(newReplySemdexer),
	)
}

type semdexer struct {
	logger        *slog.Logger
	replyQuerier  reply.Repository
	replyWriter   reply.Repository
	semdexMutator semdex.Mutator
	bus           *pubsub.Bus
}

func newReplySemdexer(
	ctx context.Context,
	lc fx.Lifecycle,
	cfg config.Config,
	logger *slog.Logger,
	replyQuerier reply.Repository,
	replyWriter reply.Repository,
	semdexMutator semdex.Mutator,
	bus *pubsub.Bus,
) {
	if cfg.SemdexProvider == "" {
		return
	}

	re := semdexer{
		logger:        logger,
		replyQuerier:  replyQuerier,
		replyWriter:   replyWriter,
		semdexMutator: semdexMutator,
		bus:           bus,
	}

	lc.Append(fx.StartHook(func(hctx context.Context) error {
		_, err := pubsub.Subscribe(hctx, bus, "reply_semdex.index_created", func(ctx context.Context, evt *message.EventThreadReplyCreated) error {
			return bus.SendCommand(ctx, &message.CommandReplyIndex{ID: evt.ReplyID})
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(hctx, bus, "reply_semdex.index_updated", func(ctx context.Context, evt *message.EventThreadReplyUpdated) error {
			return bus.SendCommand(ctx, &message.CommandReplyIndex{ID: evt.ReplyID})
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(hctx, bus, "reply_semdex.deindex_deleted", func(ctx context.Context, evt *message.EventThreadReplyDeleted) error {
			return bus.SendCommand(ctx, &message.CommandReplyDeindex{ID: evt.ReplyID})
		})
		if err != nil {
			return err
		}

		return nil
	}))

	lc.Append(fx.StartHook(func(hctx context.Context) error {
		_, err := pubsub.SubscribeCommand(hctx, bus, "reply_semdex.index", func(ctx context.Context, cmd *message.CommandReplyIndex) error {
			return re.indexReply(ctx, cmd.ID)
		})
		if err != nil {
			return err
		}

		_, err = pubsub.SubscribeCommand(hctx, bus, "reply_semdex.deindex", func(ctx context.Context, cmd *message.CommandReplyDeindex) error {
			return re.deindexReply(ctx, cmd.ID)
		})
		if err != nil {
			return err
		}

		return nil
	}))
}

func (s *semdexer) indexReply(ctx context.Context, id post.ID) error {
	p, err := s.replyQuerier.Get(ctx, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	updates, err := s.semdexMutator.Index(ctx, p)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if updates > 0 {
		_, err = s.replyWriter.Update(ctx, id)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	return nil
}

func (s *semdexer) deindexReply(ctx context.Context, id post.ID) error {
	_, err := s.semdexMutator.Delete(ctx, xid.ID(id))
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	_, err = s.replyWriter.Update(ctx, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
