package post

import (
	"context"

	"4d63.com/optional"

	"github.com/Southclaws/storyden/pkg/resources/account"
	"github.com/Southclaws/storyden/pkg/resources/post"
)

func (s *service) Create(
	ctx context.Context,
	body string,
	authorID account.AccountID,
	parentID post.PostID,
	replyToID optional.Optional[post.PostID],
) (*post.Post, error) {
	// TODO: RBAC
	return s.post_repo.Create(ctx, body, authorID, parentID, replyToID)
}
