package bindings

import (
	"context"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	account_resource "github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/app/services/authentication"
	thread_service "github.com/Southclaws/storyden/app/services/thread"
	"github.com/Southclaws/storyden/internal/openapi"
)

type Threads struct {
	thread_svc   thread_service.Service
	account_repo account_resource.Repository
}

func NewThreads(
	thread_svc thread_service.Service,
	account_repo account_resource.Repository,
) Threads {
	return Threads{thread_svc, account_repo}
}

func (i *Threads) ThreadsCreate(ctx context.Context, request openapi.ThreadsCreateRequestObject) (openapi.ThreadsCreateResponseObject, error) {
	params := func() openapi.ThreadsCreateBody {
		if request.FormdataBody != nil {
			return *request.FormdataBody
		} else {
			return *request.JSONBody
		}
	}()

	accountID, err := authentication.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var meta map[string]any
	if params.Meta != nil {
		meta = *params.Meta
	}

	thread, err := i.thread_svc.Create(ctx,
		params.Title,
		params.Body,
		accountID,
		category.CategoryID(params.Category.XID()),
		dt.Map(params.Tags, func(t openapi.Tag) string { return string(t.Id) }),
		meta,
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ThreadsCreate200JSONResponse(serialiseThread(thread)), nil
}

func reacts(reacts []*react.React) []openapi.React {
	return (dt.Map(reacts, serialiseReact))
}

func (i *Threads) ThreadsList(ctx context.Context, request openapi.ThreadsListRequestObject) (openapi.ThreadsListResponseObject, error) {
	// optionally map from OpenAPI account handle type to AccountID type.
	author, err := request.Params.Author.OptionalID(ctx, i.account_repo)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// optionally map from OpenAPI identifier type to xid.ID type.
	tags := opt.NewPtrMap(request.Params.Tags, func(t []openapi.Identifier) []xid.ID {
		return dt.Map(t, func(i openapi.Identifier) xid.ID {
			return i.XID()
		})
	})

	threads, err := i.thread_svc.ListAll(ctx, time.Now(), 10000, thread_service.Params{
		AccountID: author,
		Tags:      tags,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ThreadsList200JSONResponse(dt.Map(threads, serialiseThreadReference)), nil
}

func (i *Threads) ThreadsGet(ctx context.Context, request openapi.ThreadsGetRequestObject) (openapi.ThreadsGetResponseObject, error) {
	thread, err := i.thread_svc.Get(ctx, post.PostID(request.ThreadMark.XID()))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ThreadsGet200JSONResponse(serialiseThread(thread)), nil
}
