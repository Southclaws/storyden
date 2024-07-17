package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/content"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/post_search"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	reply_service "github.com/Southclaws/storyden/app/services/reply"
	"github.com/Southclaws/storyden/app/services/search"
	"github.com/Southclaws/storyden/app/services/thread_mark"
	"github.com/Southclaws/storyden/app/transports/openapi"
)

type Posts struct {
	reply_svc       reply_service.Service
	thread_mark_svc thread_mark.Service
	search_svc      search.Service
}

func NewPosts(
	reply_svc reply_service.Service,
	thread_mark_svc thread_mark.Service,
	search_svc search.Service,
) Posts {
	return Posts{
		reply_svc:       reply_svc,
		thread_mark_svc: thread_mark_svc,
		search_svc:      search_svc,
	}
}

func (p *Posts) PostCreate(ctx context.Context, request openapi.PostCreateRequestObject) (openapi.PostCreateResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	postID, err := p.thread_mark_svc.Lookup(ctx, string(request.ThreadMark))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	richContent, err := content.NewRichText(request.Body.Body)
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

	return openapi.PostCreate200JSONResponse{
		PostCreateOKJSONResponse: openapi.PostCreateOKJSONResponse(serialisePost(post)),
	}, nil
}

func (p *Posts) PostUpdate(ctx context.Context, request openapi.PostUpdateRequestObject) (openapi.PostUpdateResponseObject, error) {
	postID, err := p.thread_mark_svc.Lookup(ctx, string(request.PostId))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	richContent, err := opt.MapErr(opt.NewPtr(request.Body.Body), content.NewRichText)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	partial := reply_service.Partial{
		Content: richContent,
		Meta:    opt.NewPtr((*map[string]any)(request.Body.Meta)),
	}

	post, err := p.reply_svc.Update(ctx, postID, partial)
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

	err = p.reply_svc.Delete(ctx, postID)
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

func deserialisePostID(s string) post.ID {
	return post.ID(openapi.ParseID(s))
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
