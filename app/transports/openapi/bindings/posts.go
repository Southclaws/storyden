package bindings

import (
	"context"

	"4d63.com/optional"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/services/authentication"
	post_service "github.com/Southclaws/storyden/app/services/post"
	"github.com/Southclaws/storyden/internal/openapi"
)

type Posts struct {
	post_svc post_service.Service
}

func NewPosts(post_svc post_service.Service) Posts { return Posts{post_svc} }

func (p *Posts) PostsCreate(ctx context.Context, request openapi.PostsCreateRequestObject) (openapi.PostsCreateResponseObject, error) {
	accountID, err := authentication.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var reply optional.Optional[post.PostID]

	if request.Body.ReplyTo != nil {
		reply = optional.Of(post.PostID(request.Body.ReplyTo.XID()))
	}

	var meta map[string]any
	if request.Body.Meta != nil {
		meta = *request.Body.Meta
	}

	post, err := p.post_svc.Create(ctx,
		request.Body.Body,
		accountID,
		post.PostID(request.ThreadId.XID()),
		reply,
		meta,
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PostsCreate200JSONResponse(serialisePost(post)), nil
}
