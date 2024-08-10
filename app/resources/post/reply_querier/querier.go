package reply_querier

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
)

type Querier interface {
	Get(ctx context.Context, id post.ID) (*reply.Reply, error)
}
