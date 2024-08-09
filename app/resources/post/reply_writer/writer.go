package reply_writer

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/internal/ent"
)

type Writer interface {
	Create(
		ctx context.Context,
		authorID account.AccountID,
		parentID post.ID,
		opts ...Option,
	) (*reply.Reply, error)

	Update(ctx context.Context, id post.ID, opts ...Option) (*reply.Reply, error)

	Delete(ctx context.Context, id post.ID) error
}

type Option func(*ent.PostMutation)
