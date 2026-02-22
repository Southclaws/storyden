package reply

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/reply_writer"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/moderation/checker"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func (s *Mutator) Create(
	ctx context.Context,
	authorID account.AccountID,
	parentID post.ID,
	partial Partial,
) (*reply.Reply, error) {
	opts := partial.Opts()
	opts = append(opts, reply_writer.WithVisibility(visibility.VisibilityPublished))

	p, err := s.replyWriter.Create(ctx, authorID, parentID, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create reply post in thread"))
	}

	wasMovedToReview := false
	if content, ok := partial.Content.Get(); ok {
		result, err := s.cpm.CheckContent(ctx, xid.ID(p.ID), datagraph.KindReply, "", content)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if result.Action == checker.ActionReport {
			updatedReply, err := s.replyWriter.Update(ctx, p.ID, reply_writer.WithVisibility(visibility.VisibilityReview))
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
			p = updatedReply
			wasMovedToReview = true
		}
	}

	pref, err := s.replyQuerier.Probe(ctx, p.ID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.cache.Invalidate(ctx, xid.ID(pref.RootPostID)); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	replyToAuthorID := opt.Map(p.ReplyTo, func(r reply.Reply) account.AccountID {
		return r.Author.ID
	})
	replyToReplyID := opt.Map(p.ReplyTo, func(r reply.Reply) post.ID {
		return r.ID
	})

	// Only emit created event (which triggers indexing) if reply is published
	if !wasMovedToReview {
		s.bus.Publish(ctx, &rpc.EventThreadReplyCreated{
			ThreadID:        p.RootPostID,
			ReplyID:         p.ID,
			ThreadAuthorID:  p.RootAuthor.ID,
			ReplyAuthorID:   authorID,
			ReplyToAuthorID: replyToAuthorID,
			ReplyToTargetID: replyToReplyID,
		})
	}

	return p, nil
}
