package trail_runtime

import (
	"context"
	"log/slog"
	"sync"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/trail"
	"github.com/Southclaws/storyden/app/services/trail/trail_action_registry"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Runtime struct {
	ctx        context.Context
	cancel     context.CancelFunc
	logger     *slog.Logger
	repository *trail.Repository
	registry   *trail_action_registry.Registry
	bus        *pubsub.Bus

	schedulerMu    sync.Mutex
	schedulerToken string
	schedulerReady bool

	eventMu            sync.Mutex
	eventSubscriptions map[trail.ID]eventSubscription
	wg                 sync.WaitGroup
}

func New(
	ctx context.Context,
	lifecycle fx.Lifecycle,
	logger *slog.Logger,
	repository *trail.Repository,
	registry *trail_action_registry.Registry,
	bus *pubsub.Bus,
) *Runtime {
	runtimeCtx, cancel := context.WithCancel(ctx)
	runtime := &Runtime{
		ctx:                runtimeCtx,
		cancel:             cancel,
		logger:             logger,
		repository:         repository,
		registry:           registry,
		bus:                bus,
		eventSubscriptions: map[trail.ID]eventSubscription{},
	}
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			runtime.wg.Add(3)
			go func() {
				defer runtime.wg.Done()
				runtime.schedulerLoop()
			}()
			go func() {
				defer runtime.wg.Done()
				runtime.eventTriggerLoop()
			}()
			go func() {
				defer runtime.wg.Done()
				runtime.dispatchLoop()
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			runtime.cancel()
			runtime.wg.Wait()
			return nil
		},
	})
	return runtime
}
