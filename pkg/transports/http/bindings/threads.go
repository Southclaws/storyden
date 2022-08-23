package bindings

import (
	"context"

	"github.com/Southclaws/storyden/pkg/transports/http/openapi"
)

type Threads struct{}

func NewThreads() Threads { return Threads{} }

func (i *Threads) CreateThread(ctx context.Context, request openapi.CreateThreadRequestObject) any {
	return nil
}
