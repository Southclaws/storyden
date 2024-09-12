package bindings

import (
	"context"
	"net/url"
	"strconv"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/datagraph"

	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	thread_service "github.com/Southclaws/storyden/app/services/thread"
	"github.com/Southclaws/storyden/app/services/thread_mark"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Threads struct {
	thread_svc      thread_service.Service
	thread_mark_svc thread_mark.Service
	accountQuery    account_querier.Querier
}

func NewThreads(
	thread_svc thread_service.Service,
	thread_mark_svc thread_mark.Service,
	accountQuery account_querier.Querier,
) Threads {
	return Threads{thread_svc, thread_mark_svc, accountQuery}
}

func (i *Threads) ThreadCreate(ctx context.Context, request openapi.ThreadCreateRequestObject) (openapi.ThreadCreateResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	status, err := deserialiseThreadStatus(request.Body.Visibility)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var meta map[string]any
	if request.Body.Meta != nil {
		meta = *request.Body.Meta
	}

	tags := opt.NewPtr(request.Body.Tags)

	richContent, err := datagraph.NewRichText(request.Body.Body)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	url, err := opt.MapErr(opt.NewPtr(request.Body.Url), func(s string) (url.URL, error) {
		u, err := url.Parse(s)
		if err != nil {
			return url.URL{}, err
		}
		return *u, nil
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	thread, err := i.thread_svc.Create(ctx,
		request.Body.Title,
		accountID,
		category.CategoryID(openapi.ParseID(request.Body.Category)),
		status,
		tags.OrZero(),
		meta,
		thread_service.Partial{
			Content: opt.New(richContent),
			URL:     url,
		},
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ThreadCreate200JSONResponse{
		ThreadCreateOKJSONResponse: openapi.ThreadCreateOKJSONResponse(serialiseThread(thread)),
	}, nil
}

func (i *Threads) ThreadUpdate(ctx context.Context, request openapi.ThreadUpdateRequestObject) (openapi.ThreadUpdateResponseObject, error) {
	postID, err := i.thread_mark_svc.Lookup(ctx, string(request.ThreadMark))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	Visibility, err := opt.MapErr(opt.NewPtr(request.Body.Visibility), deserialiseThreadStatus)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	richContent, err := opt.MapErr(opt.NewPtr(request.Body.Body), datagraph.NewRichText)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	thread, err := i.thread_svc.Update(ctx, postID, thread_service.Partial{
		Title:      opt.NewPtr(request.Body.Title),
		Content:    richContent,
		Tags:       opt.NewPtrMap(request.Body.Tags, tagsIDs),
		Category:   opt.NewPtrMap(request.Body.Category, deserialiseID),
		Visibility: Visibility,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ThreadUpdate200JSONResponse{
		ThreadUpdateOKJSONResponse: openapi.ThreadUpdateOKJSONResponse(serialiseThread(thread)),
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
	pageSize := 50

	page := opt.NewPtrMap(request.Params.Page, func(s string) int {
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return 0
		}

		return max(1, int(v))
	}).Or(1)

	query := opt.NewPtr(request.Params.Q)

	author, err := openapi.OptionalID(ctx, i.accountQuery, request.Params.Author)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	tags := opt.NewPtrMap(request.Params.Tags, func(t []openapi.Identifier) []xid.ID {
		return dt.Map(t, func(i openapi.Identifier) xid.ID {
			return openapi.ParseID(i)
		})
	})

	cats := opt.NewPtr(request.Params.Categories)

	page = max(0, page-1)
	result, err := i.thread_svc.List(ctx, page, pageSize, thread_service.Params{
		Query:      query,
		AccountID:  author,
		Tags:       tags,
		Categories: cats,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	page = result.CurrentPage + 1

	return openapi.ThreadList200JSONResponse{
		ThreadListOKJSONResponse: openapi.ThreadListOKJSONResponse{
			CurrentPage: page,
			NextPage:    result.NextPage.Ptr(),
			PageSize:    result.PageSize,
			Results:     result.Results,
			Threads:     dt.Map(result.Threads, serialiseThreadReference),
			TotalPages:  result.TotalPages,
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

	return openapi.ThreadGet200JSONResponse{
		ThreadGetJSONResponse: openapi.ThreadGetJSONResponse(serialiseThread(thread)),
	}, nil
}

func deserialiseThreadStatus(in openapi.Visibility) (visibility.Visibility, error) {
	s, err := visibility.NewVisibility(string(in))
	if err != nil {
		return visibility.Visibility{}, fault.Wrap(err, ftag.With(ftag.InvalidArgument))
	}

	return s, nil
}
