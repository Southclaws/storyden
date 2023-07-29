package bindings

import (
	"context"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	account_resource "github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/app/resources/thread"
	"github.com/Southclaws/storyden/app/services/authentication"
	thread_service "github.com/Southclaws/storyden/app/services/thread"
	"github.com/Southclaws/storyden/app/services/thread_mark"
	"github.com/Southclaws/storyden/internal/openapi"
)

type Threads struct {
	thread_svc      thread_service.Service
	thread_mark_svc thread_mark.Service
	account_repo    account_resource.Repository
}

func NewThreads(
	thread_svc thread_service.Service,
	thread_mark_svc thread_mark.Service,
	account_repo account_resource.Repository,
) Threads {
	return Threads{thread_svc, thread_mark_svc, account_repo}
}

func (i *Threads) ThreadCreate(ctx context.Context, request openapi.ThreadCreateRequestObject) (openapi.ThreadCreateResponseObject, error) {
	accountID, err := authentication.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	status, err := deserialiseThreadStatus(request.Body.Status)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var meta map[string]any
	if request.Body.Meta != nil {
		meta = *request.Body.Meta
	}

	opts := []thread.Option{}

	tags := opt.NewPtr(request.Body.Tags)

	thread, err := i.thread_svc.Create(ctx,
		request.Body.Title,
		request.Body.Body,
		accountID,
		category.CategoryID(openapi.ParseID(request.Body.Category)),
		status,
		tags.OrZero(),
		meta,
		opts...,
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sp, err := serialiseThread(thread)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ThreadCreate200JSONResponse{
		ThreadCreateOKJSONResponse: openapi.ThreadCreateOKJSONResponse(*sp),
	}, nil
}

func (i *Threads) ThreadUpdate(ctx context.Context, request openapi.ThreadUpdateRequestObject) (openapi.ThreadUpdateResponseObject, error) {
	postID, err := i.thread_mark_svc.Lookup(ctx, string(request.ThreadMark))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	status, err := opt.MapErr(opt.NewPtr(request.Body.Status), deserialiseThreadStatus)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	thread, err := i.thread_svc.Update(ctx, postID, thread_service.Partial{
		Title:    opt.NewPtr(request.Body.Title),
		Body:     opt.NewPtr(request.Body.Body),
		Tags:     opt.NewPtrMap(request.Body.Tags, tagsIDs),
		Category: opt.NewPtrMap(request.Body.Category, deserialiseID),
		Status:   status,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sp, err := serialiseThread(thread)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ThreadUpdate200JSONResponse{
		ThreadUpdateOKJSONResponse: openapi.ThreadUpdateOKJSONResponse(*sp),
	}, nil
}

func (i *Threads) ThreadDelete(ctx context.Context, request openapi.ThreadDeleteRequestObject) (openapi.ThreadDeleteResponseObject, error) {
	postID, err := i.thread_mark_svc.Lookup(ctx, string(request.ThreadMark))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = i.thread_svc.Delete(ctx, postID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ThreadDelete200Response{}, nil
}

func reacts(reacts []*react.React) []openapi.React {
	return (dt.Map(reacts, serialiseReact))
}

func (i *Threads) ThreadList(ctx context.Context, request openapi.ThreadListRequestObject) (openapi.ThreadListResponseObject, error) {
	// optionally map from OpenAPI account handle type to AccountID type.
	author, err := openapi.OptionalID(ctx, i.account_repo, request.Params.Author)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// optionally map from OpenAPI identifier type to xid.ID type.
	tags := opt.NewPtrMap(request.Params.Tags, func(t []openapi.Identifier) []xid.ID {
		return dt.Map(t, func(i openapi.Identifier) xid.ID {
			return openapi.ParseID(i)
		})
	})

	cats := opt.NewPtr(request.Params.Categories)

	threads, err := i.thread_svc.ListAll(ctx, time.Now(), 10000, thread_service.Params{
		AccountID:  author,
		Tags:       tags,
		Categories: cats,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ThreadList200JSONResponse{
		ThreadListOKJSONResponse: openapi.ThreadListOKJSONResponse{
			Threads: dt.Map(threads, serialiseThreadReference),
		},
	}, nil
}

func (i *Threads) ThreadGet(ctx context.Context, request openapi.ThreadGetRequestObject) (openapi.ThreadGetResponseObject, error) {
	postID, err := i.thread_mark_svc.Lookup(ctx, string(request.ThreadMark))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	thread, err := i.thread_svc.Get(ctx, postID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sp, err := serialiseThread(thread)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ThreadGet200JSONResponse{
		ThreadGetJSONResponse: openapi.ThreadGetJSONResponse(*sp),
	}, nil
}

func deserialiseThreadStatus(in openapi.ThreadStatus) (thread.Status, error) {
	s, err := thread.NewStatus(string(in))
	if err != nil {
		return thread.Status{}, fault.Wrap(err, ftag.With(ftag.InvalidArgument))
	}

	return s, nil
}
