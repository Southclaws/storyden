package profile_semdex

import (
	"context"
	"log/slog"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/profile/profile_querier"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

func Build() fx.Option {
	return fx.Options(
		fx.Invoke(newProfileSemdexer),
	)
}

type semdexer struct {
	logger        *slog.Logger
	profileQuery  *profile_querier.Querier
	semdexMutator semdex.Mutator
	bus           *pubsub.Bus
}

func newProfileSemdexer(
	ctx context.Context,
	lc fx.Lifecycle,
	cfg config.Config,
	logger *slog.Logger,
	profileQuery *profile_querier.Querier,
	semdexMutator semdex.Mutator,
	bus *pubsub.Bus,
) {
	if cfg.SemdexProvider == "" {
		return
	}

	s := semdexer{
		logger:        logger,
		profileQuery:  profileQuery,
		semdexMutator: semdexMutator,
		bus:           bus,
	}

	lc.Append(fx.StartHook(func(hctx context.Context) error {
		_, err := pubsub.Subscribe(hctx, bus, "profile_semdex.index_created", func(ctx context.Context, evt *message.EventAccountCreated) error {
			return bus.SendCommand(ctx, &message.CommandProfileIndex{ID: evt.ID})
		})
		if err != nil {
			return err
		}

		_, err = pubsub.Subscribe(hctx, bus, "profile_semdex.index_updated", func(ctx context.Context, evt *message.EventAccountUpdated) error {
			return bus.SendCommand(ctx, &message.CommandProfileIndex{ID: evt.ID})
		})
		if err != nil {
			return err
		}

		return nil
	}))

	lc.Append(fx.StartHook(func(hctx context.Context) error {
		_, err := pubsub.SubscribeCommand(hctx, bus, "profile_semdex.index", func(ctx context.Context, cmd *message.CommandProfileIndex) error {
			return s.indexProfile(ctx, cmd.ID)
		})
		if err != nil {
			return err
		}

		return nil
	}))
}

func (s *semdexer) indexProfile(ctx context.Context, id account.AccountID) error {
	p, err := s.profileQuery.GetByID(ctx, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if p.GetContent().IsEmpty() {
		return nil
	}

	_, err = s.semdexMutator.Index(ctx, p)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
