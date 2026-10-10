package bindings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/idempotency"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply_querier"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/reply"
	reply_service "github.com/Southclaws/storyden/app/services/reply"
	thread_service "github.com/Southclaws/storyden/app/services/thread"
	"github.com/Southclaws/storyden/app/services/thread_mark"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Replies struct {
	replyMutator    *reply.Mutator
	thread_mark_svc thread_mark.Service
	threadSvc       thread_service.Service
	replyQuerier    *reply_querier.Querier
	receipts        *idempotency.Repository
}

func NewReplies(
	replyMutator *reply.Mutator,
	thread_mark_svc thread_mark.Service,
	threadSvc thread_service.Service,
	replyQuerier *reply_querier.Querier,
	receipts *idempotency.Repository,
) Replies {
	return Replies{
		replyMutator:    replyMutator,
		thread_mark_svc: thread_mark_svc,
		threadSvc:       threadSvc,
		replyQuerier:    replyQuerier,
		receipts:        receipts,
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

	if err := session.Authorise(ctx, nil, rbac.PermissionCreatePost); err != nil {
		return nil, err
	}
	response, err := idempotentCreate(ctx, p.receipts, accountID.String(), "ReplyCreate", request.Params.IdempotencyKey,
		struct {
			ThreadMark openapi.ThreadMarkParam
			Body       *openapi.ReplyCreateJSONRequestBody
		}{request.ThreadMark, request.Body},
		func() (openapi.ReplyCreate200JSONResponse, error) {
			post, err := p.replyMutator.Create(ctx, accountID, postID, partial)
			if err != nil {
				return openapi.ReplyCreate200JSONResponse{}, err
			}
			return openapi.ReplyCreate200JSONResponse{ReplyCreateOKJSONResponse: openapi.ReplyCreateOKJSONResponse(serialiseReplyPtr(post))}, nil
		}, func(response openapi.ReplyCreate200JSONResponse) error {
			_, err := p.threadSvc.Get(ctx, postID, pagination.NewPageParams(1, 1))
			if err != nil {
				return err
			}
			existing, err := p.replyQuerier.Get(ctx, post.ID(openapi.ParseID(response.Id)))
			if err != nil {
				return err
			}
			if existing.DeletedAt.Ok() {
				return fault.New("reply no longer exists", ftag.With(ftag.NotFound))
			}
			if existing.Visibility != visibility.VisibilityPublished &&
				!(existing.Visibility == visibility.VisibilityReview && existing.Author.ID == accountID) {
				if err := session.Authorise(ctx, nil, rbac.PermissionManagePosts); err != nil {
					return err
				}
			}
			return nil
		})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return response, nil
}
