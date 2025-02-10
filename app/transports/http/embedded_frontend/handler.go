package embedded_frontend

import (
	"net/http"

	"github.com/Southclaws/storyden/app/services/thread"
	"github.com/Southclaws/storyden/app/transports/http/embedded_frontend/components"
)

type Handler http.Handler

func New(
	threadService thread.Service,
) Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		list, err := threadService.List(ctx, 1, 50, thread.Params{})
		if err != nil {
			return
		}

		components.Index{
			Threads: list,
		}.Page().Render(ctx, w)
	}))

	return mux
}
