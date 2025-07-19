package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/post/reaction"
	"github.com/Southclaws/storyden/app/services/react_manager"
	"github.com/Southclaws/storyden/app/services/thread_mark"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Reacts struct {
	thread_mark_svc thread_mark.Service
	reactor         *react_manager.Reactor
}

func NewReacts(thread_mark_svc thread_mark.Service, reactor *react_manager.Reactor) Reacts {
	return Reacts{thread_mark_svc, reactor}
}

func (p *Reacts) PostReactAdd(ctx context.Context, request openapi.PostReactAddRequestObject) (openapi.PostReactAddResponseObject, error) {
	postID, err := p.thread_mark_svc.Lookup(ctx, string(request.PostId))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	react, err := p.reactor.Add(ctx, postID, request.Body.Emoji)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PostReactAdd200JSONResponse{
		PostReactAddOKJSONResponse: openapi.PostReactAddOKJSONResponse(serialiseReact(react)),
	}, nil
}

func (h *Reacts) PostReactRemove(ctx context.Context, request openapi.PostReactRemoveRequestObject) (openapi.PostReactRemoveResponseObject, error) {
	reactID := reaction.ReactID(openapi.ParseID(request.ReactId))

	err := h.reactor.Remove(ctx, reactID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PostReactRemove200Response{}, nil
}

func serialiseReact(r *reaction.React) openapi.React {
	return openapi.React{
		Id:     xid.ID(r.ID).String(),
		Emoji:  r.Emoji,
		Author: serialiseProfileReference(r.Author),
	}
}

func serialiseReactList(reacts []*reaction.React) []openapi.React {
	return dt.Map(reacts, serialiseReact)
}
