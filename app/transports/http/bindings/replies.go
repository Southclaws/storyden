package bindings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	reply_service "github.com/Southclaws/storyden/app/services/reply"
	"github.com/Southclaws/storyden/app/services/thread_mark"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Replies struct {
	reply_svc       reply_service.Service
	thread_mark_svc thread_mark.Service
}

func NewReplies(
	reply_svc reply_service.Service,
	thread_mark_svc thread_mark.Service,
) Replies {
	return Replies{
		reply_svc:       reply_svc,
		thread_mark_svc: thread_mark_svc,
	}
}

func (p *Replies) ReplyCreate(ctx context.Context, request openapi.ReplyCreateRequestObject) (openapi.ReplyCreateResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	postID, err := p.thread_mark_svc.Lookup(ctx, string(request.ThreadMark))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	richContent, err := datagraph.NewRichText(request.Body.Body)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	partial := reply_service.Partial{
		Content: opt.New(richContent),
		ReplyTo: opt.Map(opt.NewPtr(request.Body.ReplyTo), deserialisePostID),
		Meta:    opt.NewPtr((*map[string]any)(request.Body.Meta)),
	}

	post, err := p.reply_svc.Create(ctx,
		accountID,
		postID,
		partial,
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ReplyCreate200JSONResponse{
		ReplyCreateOKJSONResponse: openapi.ReplyCreateOKJSONResponse(serialiseReply(post)),
	}, nil
}
