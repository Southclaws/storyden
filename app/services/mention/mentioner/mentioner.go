package mentioner

import (
	"context"
	"log/slog"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type Mentioner struct {
	logger *slog.Logger
	bus    *pubsub.Bus
}

func New(logger *slog.Logger, bus *pubsub.Bus) *Mentioner {
	return &Mentioner{logger: logger, bus: bus}
}

func (n *Mentioner) Send(ctx context.Context, by account.AccountID, source datagraph.Ref, items ...*datagraph.Ref) {
	sender, err := session.GetAccountID(ctx)
	if err != nil {
		n.logger.Warn("cannot send notification without source session", slog.String("error", err.Error()))
		return
	}

	for _, i := range items {
		if i.Kind == datagraph.KindProfile && sender == account.AccountID(i.ID) {
			// Skip self-mentions
			continue
		}

		n.bus.Publish(ctx, &rpc.EventMemberMentioned{
			By:     by,
			Source: rpc.DatagraphRefToRPC(source),
			Item:   rpc.DatagraphRefToRPC(*i),
		})
	}
}
