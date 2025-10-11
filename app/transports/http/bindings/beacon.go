package bindings

import (
	"context"
	"encoding/json"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Beacon struct {
	bus *pubsub.Bus
}

func NewBeacon(bus *pubsub.Bus) Beacon {
	return Beacon{
		bus: bus,
	}
}

type Message struct {
	Kind datagraph.Kind `json:"k"`
	ID   xid.ID         `json:"id"`
}

// NOTE: Does not handle errors or return anything other than 202.
func (b *Beacon) SendBeacon(ctx context.Context, request openapi.SendBeaconRequestObject) (openapi.SendBeaconResponseObject, error) {
	if request.Body == nil {
		return nil, nil
	}

	session := session.GetOptAccountID(ctx)

	var m Message
	err := json.Unmarshal([]byte(*request.Body), &m)
	if err != nil {
		return nil, nil
	}

	b.bus.SendCommand(ctx, &message.CommandSendBeacon{
		Item: datagraph.Ref{
			Kind: m.Kind,
			ID:   m.ID,
		},
		Subject: session,
	})

	return openapi.SendBeacon202Response{}, nil
}
