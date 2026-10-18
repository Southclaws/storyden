package background_job

import (
	"context"
	"encoding/json"
	"errors"

	jobs "github.com/Southclaws/storyden/app/resources/background_job"
)

type Result struct {
	Terminal bool
	Deferred bool
}

type Handler[T any] struct {
	Run      func(context.Context, *jobs.Job[T]) (Result, error)
	Complete func(T) bool
}

type handler struct {
	Run      func(context.Context, *jobs.Job[json.RawMessage]) (Result, error)
	Complete func(json.RawMessage) (bool, error)
}

type worker struct {
	handler handler
	slots   chan struct{}
}

func (r *Runner) Register[T any](kind string, workers int, typed Handler[T]) error {
	if typed.Run == nil {
		return errors.New("background job handler requires a run function")
	}

	adapted := handler{
		Run: func(ctx context.Context, raw *jobs.Job[json.RawMessage]) (Result, error) {
			payload, err := jobs.Decode[T](raw.Payload)
			if err != nil {
				return Result{Terminal: true}, err
			}
			return typed.Run(ctx, &jobs.Job[T]{State: raw.State, Payload: payload})
		},
	}

	if typed.Complete != nil {
		adapted.Complete = func(raw json.RawMessage) (bool, error) {
			payload, err := jobs.Decode[T](raw)
			if err != nil {
				return false, err
			}
			return typed.Complete(payload), nil
		}
	}

	return r.register(kind, workers, adapted)
}
