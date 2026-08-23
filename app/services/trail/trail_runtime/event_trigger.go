package trail_runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Southclaws/storyden/app/resources/trail"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type eventSubscription struct {
	event        string
	subscription *pubsub.Subscription
}

func (r *Runtime) eventTriggerLoop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	defer r.closeEventSubscriptions()

	for {
		if err := r.syncEventSubscriptions(r.ctx); err != nil {
			r.logger.Error("failed to sync Trail event subscriptions", slog.String("error", err.Error()))
		}

		select {
		case <-r.ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (r *Runtime) syncEventSubscriptions(ctx context.Context) error {
	definitions, err := r.repository.ActiveEventTriggers(ctx)
	if err != nil {
		return err
	}

	desired := make(map[trail.ID]string)
	for _, definition := range definitions {
		var config trail.EventTriggerConfig
		if err := json.Unmarshal(definition.Config, &config); err != nil {
			r.logger.Error(
				"invalid persisted Trail event trigger",
				slog.String("trail_id", definition.ID.String()),
				slog.String("error", err.Error()),
			)
			continue
		}

		desired[definition.ID] = config.Event
	}

	r.eventMu.Lock()
	defer r.eventMu.Unlock()

	for id, current := range r.eventSubscriptions {
		event, keep := desired[id]
		if keep && event == current.event {
			continue
		}

		if err := current.subscription.Close(); err != nil {
			r.logger.Error("failed to close Trail event subscription", slog.String("trail_id", id.String()), slog.String("error", err.Error()))
		}
		delete(r.eventSubscriptions, id)
	}

	for id, event := range desired {
		if _, exists := r.eventSubscriptions[id]; exists {
			continue
		}

		trailID := id
		eventName := event
		handlerName := fmt.Sprintf("trail_event_%s", id.String())
		topicName := "rpc." + eventName
		subscription, err := pubsub.SubscribeNamed(ctx, r.bus, topicName, handlerName, func(ctx context.Context, payload json.RawMessage) error {
			_, _, err := r.repository.MaterialiseEvent(ctx, trailID, eventName, payload, time.Now().UTC())
			return err
		})
		if err != nil {
			r.logger.Error(
				"failed to attach Trail event subscription",
				slog.String("trail_id", id.String()),
				slog.String("event", eventName),
				slog.String("error", err.Error()),
			)
			continue
		}

		r.eventSubscriptions[id] = eventSubscription{event: eventName, subscription: subscription}
	}

	return nil
}

func (r *Runtime) closeEventSubscriptions() {
	r.eventMu.Lock()
	defer r.eventMu.Unlock()

	for id, current := range r.eventSubscriptions {
		if err := current.subscription.Close(); err != nil {
			r.logger.Error("failed to close Trail event subscription", slog.String("trail_id", id.String()), slog.String("error", err.Error()))
		}
		delete(r.eventSubscriptions, id)
	}
}
