package bindings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/reply"
	reply_service "github.com/Southclaws/storyden/app/services/reply"
	"github.com/Southclaws/storyden/app/services/thread"
	"github.com/Southclaws/storyden/app/services/thread_mark"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Replies struct {
	replyMutator    *reply.Mutator
	thread_mark_svc thread_mark.Service
	threadService   thread.Service
}

func NewReplies(
	replyMutator *reply.Mutator,
	thread_mark_svc thread_mark.Service,
	threadService thread.Service,
) Replies {
	return Replies{
		replyMutator:    replyMutator,
		thread_mark_svc: thread_mark_svc,
		threadService:   threadService,
	}
}

func (p *Replies) ReplyGet(ctx context.Context, request openapi.ReplyGetRequestObject) (openapi.ReplyGetResponseObject, error) {
	threadID, err := p.thread_mark_svc.Lookup(ctx, string(request.ThreadMark))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	replyID, err := xid.FromString(request.PostId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	r, err := p.threadService.GetReply(ctx, threadID, post.ID(replyID))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ReplyGet200JSONResponse{
		ReplyGetOKJSONResponse: openapi.ReplyGetOKJSONResponse(serialiseReplyPtr(r)),
	}, nil
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

	post, err := p.replyMutator.Create(ctx,
		accountID,
		postID,
		partial,
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ReplyCreate200JSONResponse{
		ReplyCreateOKJSONResponse: openapi.ReplyCreateOKJSONResponse(serialiseReplyPtr(post)),
	}, nil
}
