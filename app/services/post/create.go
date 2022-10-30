package post

import (
	"context"

	"4d63.com/optional"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/post"
)

func (s *service) Create(
	ctx context.Context,
	body string,
	authorID account.AccountID,
	parentID post.PostID,
	replyToID optional.Optional[post.PostID],
	meta map[string]any,
) (*post.Post, error) {
	p, err := s.post_repo.Create(ctx, body, authorID, parentID, replyToID, meta)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create post in thread"))
	}

	return p, nil
}
