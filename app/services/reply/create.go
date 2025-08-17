package reply

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
)

func (s *service) Create(
	ctx context.Context,
	authorID account.AccountID,
	parentID post.ID,
	partial Partial,
) (*reply.Reply, error) {
	if content, ok := partial.Content.Get(); ok {
		if err := s.cpm.CheckContent(ctx, content); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	opts := partial.Opts()

	p, err := s.post_repo.Create(ctx, authorID, parentID, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create reply post in thread"))
	}

	s.bus.Publish(ctx, &message.EventThreadReplyCreated{
		ThreadID:       p.RootPostID,
		ReplyID:        p.ID,
		ThreadAuthorID: p.RootAuthor.ID,
		ReplyAuthorID:  authorID,
	})


	return p, nil
}
