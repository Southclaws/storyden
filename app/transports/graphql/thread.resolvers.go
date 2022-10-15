package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/Southclaws/fault/errctx"

	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/transports/graphql/models"
	"github.com/Southclaws/storyden/app/transports/graphql/server"
)

// CreateThread is the resolver for the createThread field.
func (r *mutationResolver) CreateThread(ctx context.Context, input models.NewThread) (*models.Thread, error) {
	acc, err := authentication.GetAccountID(ctx)
	if err != nil {
		return nil, err
	}

	thread, err := r.thread_service.Create(ctx, input.Title, input.Body, acc, category.CategoryID{}, nil)
	if err != nil {
		return nil, errctx.Wrap(err, ctx)
	}

	return &models.Thread{
		ID: thread.ID.String(),
	}, nil
}

// Mutation returns server.MutationResolver implementation.
func (r *Resolver) Mutation() server.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
