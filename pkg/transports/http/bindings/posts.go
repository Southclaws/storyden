package bindings

import (
	"context"

	"4d63.com/optional"

	"github.com/Southclaws/storyden/pkg/resources/post"
	"github.com/Southclaws/storyden/pkg/services/authentication"
	post_service "github.com/Southclaws/storyden/pkg/services/post"
	"github.com/Southclaws/storyden/pkg/transports/http/openapi"
)

type Posts struct {
	post_svc post_service.Service
}

func NewPosts(post_svc post_service.Service) Posts { return Posts{post_svc} }

func (p *Posts) PostsCreate(ctx context.Context, request openapi.PostsCreateRequestObject) any {
	accountID, err := authentication.GetAccountID(ctx)
	if err != nil {
		return err
	}

	params := func() openapi.PostsCreate {
		if request.FormdataBody != nil {
			return *request.FormdataBody
		} else {
			return *request.JSONBody
		}
	}()

	var reply optional.Optional[post.PostID]

	if params.ReplyTo != nil {
		reply = optional.Of(post.PostID(params.ReplyTo.XID()))
	}

	post, err := p.post_svc.Create(ctx,
		params.Body,
		accountID,
		post.PostID(request.Id.XID()),
		reply,
	)
	if err != nil {
		return err
	}

	return openapi.PostsCreate200JSONResponse(serialisePost(post))
}
