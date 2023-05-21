package bindings

import (
	"context"

	"4d63.com/optional"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post_search"
	"github.com/Southclaws/storyden/app/services/authentication"
	post_service "github.com/Southclaws/storyden/app/services/post"
	"github.com/Southclaws/storyden/app/services/search"
	"github.com/Southclaws/storyden/app/services/thread_mark"
	"github.com/Southclaws/storyden/internal/openapi"
)

type Posts struct {
	post_svc        post_service.Service
	thread_mark_svc thread_mark.Service
	search_svc      search.Service
}

func NewPosts(
	post_svc post_service.Service,
	thread_mark_svc thread_mark.Service,
	search_svc search.Service,
) Posts {
	return Posts{
		post_svc:        post_svc,
		thread_mark_svc: thread_mark_svc,
		search_svc:      search_svc,
	}
}

func (p *Posts) PostCreate(ctx context.Context, request openapi.PostCreateRequestObject) (openapi.PostCreateResponseObject, error) {
	accountID, err := authentication.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	postID, err := p.thread_mark_svc.Lookup(ctx, string(request.ThreadMark))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var reply optional.Optional[post.PostID]

	if request.Body.ReplyTo != nil {
		tm := openapi.ParseID(*request.Body.ReplyTo)
		reply = optional.Of(post.PostID(tm))
	}

	var meta map[string]any
	if request.Body.Meta != nil {
		meta = *request.Body.Meta
	}

	post, err := p.post_svc.Create(ctx,
		request.Body.Body,
		accountID,
		postID,
		reply,
		meta,
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PostCreate200JSONResponse{
		PostCreateOKJSONResponse: openapi.PostCreateOKJSONResponse(serialisePost(post)),
	}, nil
}

func (p *Posts) PostUpdate(ctx context.Context, request openapi.PostUpdateRequestObject) (openapi.PostUpdateResponseObject, error) {
	postID, err := p.thread_mark_svc.Lookup(ctx, string(request.PostId))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	post, err := p.post_svc.Update(ctx, postID, post_service.Partial{
		Body: opt.NewPtr(request.Body.Body),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PostUpdate200JSONResponse{
		PostUpdateOKJSONResponse: openapi.PostUpdateOKJSONResponse(serialisePost(post)),
	}, nil
}

func (p *Posts) PostDelete(ctx context.Context, request openapi.PostDeleteRequestObject) (openapi.PostDeleteResponseObject, error) {
	postID, err := p.thread_mark_svc.Lookup(ctx, string(request.PostId))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = p.post_svc.Delete(ctx, postID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PostDelete200Response{}, nil
}

func (p *Posts) PostSearch(ctx context.Context, request openapi.PostSearchRequestObject) (openapi.PostSearchResponseObject, error) {
	ks := opt.MapErrC(deserialiseContentKinds)

	kinds, err := ks(opt.NewPtr(request.Params.Kind))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	posts, err := p.search_svc.Search(ctx, search.Query{
		Body:   opt.NewPtr(request.Params.Body),
		Author: opt.NewPtr(request.Params.Author),
		Kinds:  kinds,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	results := dt.Map(posts, serialisePost)

	return openapi.PostSearch200JSONResponse{
		PostSearchOKJSONResponse: openapi.PostSearchOKJSONResponse{
			Count:   float32(len(results)),
			Results: results,
		},
	}, nil
}

func deserialiseContentKinds(in openapi.ContentKinds) ([]post_search.Kind, error) {
	out, err := dt.MapErr(in, deserialiseContentKind)
	if err != nil {
		return nil, fault.Wrap(err)
	}
	return out, nil
}

func deserialiseContentKind(in openapi.ContentKind) (post_search.Kind, error) {
	out, err := post_search.NewKind(string(in))
	if err != nil {
		return post_search.Kind{}, fault.Wrap(err)
	}

	return out, nil
}
