package bindings

import (
	"context"
	"time"

	"github.com/Southclaws/dt"

	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/app/services/authentication"
	thread_service "github.com/Southclaws/storyden/app/services/thread"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/errctx"
)

type Threads struct {
	thread_svc thread_service.Service
}

func NewThreads(thread_svc thread_service.Service) Threads { return Threads{thread_svc} }

func (i *Threads) ThreadsCreate(ctx context.Context, request openapi.ThreadsCreateRequestObject) any {
	params := func() openapi.ThreadSubmission {
		if request.FormdataBody != nil {
			return *request.FormdataBody
		} else {
			return *request.JSONBody
		}
	}()

	accountID, err := authentication.GetAccountID(ctx)
	if err != nil {
		return err
	}

	thread, err := i.thread_svc.Create(ctx, params.Title, params.Body, accountID, category.CategoryID(params.Category.XID()), params.Tags)
	if err != nil {
		return errctx.Wrap(err, ctx)
	}

	return openapi.ThreadsCreate200JSONResponse(serialiseThread(thread))
}

func reacts(reacts []*react.React) []openapi.React {
	return (dt.Map(reacts, serialiseReact))
}

func (i *Threads) ThreadsList(ctx context.Context, request openapi.ThreadsListRequestObject) any {
	threads, err := i.thread_svc.ListAll(ctx, time.Now(), 10000)
	if err != nil {
		return errctx.Wrap(err, ctx)
	}

	return openapi.ThreadsList200JSONResponse(dt.Map(threads, serialiseThreadReference))
}

func (i *Threads) ThreadsGet(ctx context.Context, request openapi.ThreadsGetRequestObject) any {
	thread, err := i.thread_svc.Get(ctx, post.PostID(request.ThreadId.XID()))
	if err != nil {
		return errctx.Wrap(err, ctx)
	}

	return openapi.ThreadsGet200JSONResponse(serialiseThread(thread))
}
