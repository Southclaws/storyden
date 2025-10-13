package bindings

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Beacon struct {
	logger *slog.Logger
	bus    *pubsub.Bus
}

func NewBeacon(logger *slog.Logger, bus *pubsub.Bus) Beacon {
	return Beacon{
		logger: logger,
		bus:    bus,
	}
}

type Message struct {
	Kind datagraph.Kind `json:"k"`
	ID   xid.ID         `json:"id"`
}

// NOTE: Does not handle errors or return anything other than 202.
func (b *Beacon) SendBeacon(ctx context.Context, request openapi.SendBeaconRequestObject) (openapi.SendBeaconResponseObject, error) {
	accountID := session.GetOptAccountID(ctx)

	log := b.logger.With(slog.String("handler", "beacon.SendBeacon"))

	if v, ok := accountID.Get(); ok {
		log = log.With(slog.String("account_id", v.String()))
	}

	if request.Body == nil {
		log.Warn("beacon request with empty body")
		return openapi.SendBeacon202Response{}, nil
	}

	bodyStr := *request.Body

	log = log.With(slog.String("body", bodyStr))

	var m Message
	err := json.Unmarshal([]byte(bodyStr), &m)
	if err != nil {
		log.Warn("failed to unmarshal beacon body", slog.String("error", err.Error()))
		return openapi.SendBeacon202Response{}, nil
	}

	if err := b.bus.SendCommand(ctx, &message.CommandSendBeacon{
		Item: datagraph.Ref{
			Kind: m.Kind,
			ID:   m.ID,
		},
		Subject: accountID,
	}); err != nil {
		log.Error("failed to send beacon command", slog.String("error", err.Error()))
	}

	return openapi.SendBeacon202Response{}, nil
}
