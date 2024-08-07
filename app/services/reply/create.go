package reply

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
)

func (s *service) Create(
	ctx context.Context,
	authorID account.AccountID,
	parentID post.ID,
	partial Partial,
) (*reply.Reply, error) {
	opts := partial.Opts()

	opts = append(opts, s.hydrate(ctx, partial)...)

	p, err := s.post_repo.Create(ctx, authorID, parentID, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create reply post in thread"))
	}

	if err := s.indexQueue.Publish(ctx, mq.IndexPost{
		ID: p.ID,
	}); err != nil {
		s.l.Error("failed to publish index post message", zap.Error(err))
	}

	return p, nil
}
