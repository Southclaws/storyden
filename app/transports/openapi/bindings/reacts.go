package bindings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/react"
	"github.com/Southclaws/storyden/app/services/thread_mark"
	"github.com/Southclaws/storyden/internal/openapi"
)

type Reacts struct {
	thread_mark_svc thread_mark.Service
	react_svc       react.Service
}

func NewReacts(thread_mark_svc thread_mark.Service, react_svc react.Service) Reacts {
	return Reacts{thread_mark_svc, react_svc}
}

func (p *Reacts) PostReactAdd(ctx context.Context, request openapi.PostReactAddRequestObject) (openapi.PostReactAddResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	postID, err := p.thread_mark_svc.Lookup(ctx, string(request.PostId))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	react, err := p.react_svc.Add(ctx, accountID, postID, *request.Body.Emoji)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PostReactAdd200JSONResponse{
		PostReactAddOKJSONResponse: openapi.PostReactAddOKJSONResponse(serialiseReact(react)),
	}, nil
}
