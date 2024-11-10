package reply

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/datagraph"
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

	if err := s.indexQueue.Publish(ctx, mq.IndexReply{
		ID: p.ID,
	}); err != nil {
		s.l.Error("failed to publish index post message", zap.Error(err))
	}

	s.fetcher.HydrateContentURLs(ctx, p)

	s.notifier.Send(ctx, p.RootAuthor.ID, notification.EventThreadReply, &datagraph.Ref{
		ID:   xid.ID(p.RootPostID),
		Kind: datagraph.KindPost,
	})

	return p, nil
}
