package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/pkg/resources/category"
	"github.com/Southclaws/storyden/pkg/resources/post"
	"github.com/Southclaws/storyden/pkg/resources/react"
	"github.com/Southclaws/storyden/pkg/services/authentication"
	"github.com/Southclaws/storyden/pkg/services/thread"
	"github.com/Southclaws/storyden/pkg/transports/http/openapi"
)

type Threads struct {
	thread_svc thread.Service
}

func NewThreads(thread_svc thread.Service) Threads { return Threads{thread_svc} }

func (i *Threads) CreateThread(ctx context.Context, request openapi.CreateThreadRequestObject) any {
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

	thread, err := i.thread_svc.Create(ctx, params.Title, params.Body, accountID, category.CategoryID(params.Category), params.Tags)
	if err != nil {
		return err
	}

	return openapi.CreateThread200JSONResponse{
		Id:        openapi.Identifier(thread.ID),
		CreatedAt: thread.CreatedAt,
		UpdatedAt: thread.UpdatedAt,
		DeletedAt: utils.OptionalToPointer(thread.DeletedAt),

		Title: thread.Title,
		Author: &openapi.ProfileReference{
			Id:   (*openapi.Identifier)(&thread.Author.ID),
			Name: &thread.Author.Name,
		},
		Slug:  &thread.Slug,
		Short: &thread.Short,

		Category: &openapi.Category{
			// thread.Category.ID
		},
		Pinned: &thread.Pinned,
		Posts: utils.Ref(dt.Reduce(thread.Posts, func(s int, p *post.Post) int {
			return s + 1
		}, 0)),
		Reacts: utils.Ref(dt.Map(thread.Reacts, func(r *react.React) openapi.React {
			return openapi.React{
				Id:    (*openapi.Identifier)(&r.ID),
				Emoji: &r.Emoji,
			}
		})),
		Tags: &thread.Tags,
	}
}
