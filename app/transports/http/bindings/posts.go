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
	"github.com/Southclaws/storyden/app/resources/post/post_search"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/reply"
	"github.com/Southclaws/storyden/app/services/thread_mark"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Posts struct {
	post_repo       post_search.Repository
	replyMutator    *reply.Mutator
	thread_mark_svc thread_mark.Service
}

func NewPosts(
	post_repo post_search.Repository,
	replyMutator *reply.Mutator,
	thread_mark_svc thread_mark.Service,
) Posts {
	return Posts{
		post_repo:       post_repo,
		replyMutator:    replyMutator,
		thread_mark_svc: thread_mark_svc,
	}
}

func (p *Posts) PostUpdate(ctx context.Context, request openapi.PostUpdateRequestObject) (openapi.PostUpdateResponseObject, error) {
	postID, err := p.thread_mark_svc.Lookup(ctx, string(request.PostId))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	richContent, err := opt.MapErr(opt.NewPtr(request.Body.Body), datagraph.NewRichText)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	vis, err := opt.MapErr(opt.NewPtr(request.Body.Visibility), func(v openapi.Visibility) (visibility.Visibility, error) {
		return visibility.NewVisibility(string(v))
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	partial := reply.Partial{
		Content:    richContent,
		Visibility: vis,
		Meta:       opt.NewPtr((*map[string]any)(request.Body.Meta)),
	}

	post, err := p.replyMutator.Update(ctx, postID, partial)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PostUpdate200JSONResponse{
		PostUpdateOKJSONResponse: openapi.PostUpdateOKJSONResponse(serialisePost(&post.Post)),
	}, nil
}

func (p *Posts) PostDelete(ctx context.Context, request openapi.PostDeleteRequestObject) (openapi.PostDeleteResponseObject, error) {
	postID, err := p.thread_mark_svc.Lookup(ctx, string(request.PostId))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = p.replyMutator.Delete(ctx, postID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PostDelete200Response{}, nil
}

func (p *Posts) PostLocationGet(ctx context.Context, request openapi.PostLocationGetRequestObject) (openapi.PostLocationGetResponseObject, error) {
	id, err := xid.FromString(string(request.Params.Id))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	location, err := p.post_repo.Locate(ctx, post.ID(id))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PostLocationGet200JSONResponse{
		PostLocationGetOKJSONResponse: openapi.PostLocationGetOKJSONResponse{
			Slug:     location.Slug,
			Kind:     openapi.PostLocationKind(location.Kind.String()),
			Index:    location.Index.Ptr(),
			Page:     location.Page.Ptr(),
			Position: location.Position.Ptr(),
		},
	}, nil
}
