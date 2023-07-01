package post

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/post"
)

func (s *service) Create(
	ctx context.Context,
	body post.Content,
	authorID account.AccountID,
	parentID post.PostID,
	replyToID opt.Optional[post.PostID],
	meta map[string]any,
	opts ...post.Option,
) (*post.Post, error) {
	p, err := s.post_repo.Create(ctx, body, authorID, parentID, replyToID, meta, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create post in thread"))
	}

	return p, nil
}
