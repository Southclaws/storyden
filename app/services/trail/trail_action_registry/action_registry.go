package trail_action_registry

import (
	"context"
	"fmt"
	"sync"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/trail"
)

var (
	errUnsupportedAction = fault.New("unsupported trail action", ftag.With(ftag.InvalidArgument))
	errActionNotActive   = fault.New("trail action run is not active", ftag.With(ftag.InvalidArgument))
)

type ActionResult struct {
	Status    trail.ActionRunStatus
	Output    []byte
	ErrorText string
}

type ActionAdapter interface {
	Kind() trail.ActionKind
	Validate(context.Context, *account.Account, trail.ActionSpec) error
	Start(context.Context, *trail.ActionRun) error
	Reconcile(context.Context, *trail.ActionRun) (*ActionResult, error)
	Cancel(context.Context, *trail.ActionRun) error
}

type Registry struct {
	repository *trail.Repository
	mu         sync.RWMutex
	adapters   map[trail.ActionKind]ActionAdapter
}

func New(repository *trail.Repository, robot *RobotActionAdapter) *Registry {
	registry := &Registry{repository: repository, adapters: map[trail.ActionKind]ActionAdapter{}}
	registry.adapters[robot.Kind()] = robot

	return registry
}

func (r *Registry) adapter(kind trail.ActionKind) (ActionAdapter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	adapter, ok := r.adapters[kind]
	if !ok {
		return nil, fmt.Errorf("%w: action type %q is not registered", errUnsupportedAction, kind)
	}

	return adapter, nil
}

func (r *Registry) Validate(ctx context.Context, owner *account.Account, spec trail.ActionSpec) error {
	adapter, err := r.adapter(spec.Kind)
	if err != nil {
		return err
	}

	return adapter.Validate(ctx, owner, spec)
}

func (r *Registry) Start(ctx context.Context, run *trail.ActionRun) error {
	adapter, err := r.adapter(run.Kind)
	if err != nil {
		return err
	}

	return adapter.Start(ctx, run)
}

func (r *Registry) Reconcile(ctx context.Context, run *trail.ActionRun) (*ActionResult, error) {
	adapter, err := r.adapter(run.Kind)
	if err != nil {
		return nil, err
	}

	return adapter.Reconcile(ctx, run)
}

func (r *Registry) Cancel(ctx context.Context, run *trail.ActionRun) error {
	if run.Status == trail.ActionRunStatusQueued {
		_, err := r.repository.CancelQueuedActionRun(ctx, run.ID)
		return err
	}

	if run.Status != trail.ActionRunStatusRunning {
		return errActionNotActive
	}

	adapter, err := r.adapter(run.Kind)
	if err != nil {
		return err
	}

	return adapter.Cancel(ctx, run)
}
