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
	subscriptions map[string]*pubsub.Subscription
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

	desired := make(map[trail.ID][]string)
	for _, definition := range definitions {
		desired[definition.ID] = definition.Events
	}

	r.eventMu.Lock()
	defer r.eventMu.Unlock()

	for id, current := range r.eventSubscriptions {
		events, keep := desired[id]
		if keep && sameEventSubscriptions(current, events) {
			continue
		}

		r.closeEventSubscription(id, current)
		delete(r.eventSubscriptions, id)
	}

	for id, events := range desired {
		if _, exists := r.eventSubscriptions[id]; exists {
			continue
		}

		current := eventSubscription{subscriptions: make(map[string]*pubsub.Subscription, len(events))}
		for _, event := range events {
			trailID := id
			eventName := event
			handlerName := fmt.Sprintf("trail_event_%s_%s", id.String(), eventName)
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
				r.closeEventSubscription(id, current)
				current.subscriptions = nil
				break
			}

			current.subscriptions[eventName] = subscription
		}

		if current.subscriptions != nil {
			r.eventSubscriptions[id] = current
		}
	}

	return nil
}

func (r *Runtime) closeEventSubscriptions() {
	r.eventMu.Lock()
	defer r.eventMu.Unlock()

	for id, current := range r.eventSubscriptions {
		r.closeEventSubscription(id, current)
		delete(r.eventSubscriptions, id)
	}
}

func sameEventSubscriptions(current eventSubscription, events []string) bool {
	if len(current.subscriptions) != len(events) {
		return false
	}

	for _, event := range events {
		if _, exists := current.subscriptions[event]; !exists {
			return false
		}
	}

	return true
}

func (r *Runtime) closeEventSubscription(id trail.ID, current eventSubscription) {
	for event, subscription := range current.subscriptions {
		if err := subscription.Close(); err != nil {
			r.logger.Error(
				"failed to close Trail event subscription",
				slog.String("trail_id", id.String()),
				slog.String("event", event),
				slog.String("error", err.Error()),
			)
		}
	}
}
