package reply

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
)

func (s *Mutator) Create(
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

	p, err := s.replyWriter.Create(ctx, authorID, parentID, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create reply post in thread"))
	}

	pref, err := s.replyQuerier.Probe(ctx, p.ID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.cache.Invalidate(ctx, xid.ID(pref.RootPostID)); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &message.EventThreadReplyCreated{
		ThreadID:       p.RootPostID,
		ReplyID:        p.ID,
		ThreadAuthorID: p.RootAuthor.ID,
		ReplyAuthorID:  authorID,
	})

	return p, nil
}
