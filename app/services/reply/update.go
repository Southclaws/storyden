package reply

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

func (s *Mutator) Update(ctx context.Context, threadID post.ID, partial Partial) (*reply.Reply, error) {
	if content, ok := partial.Content.Get(); ok {
		if err := s.cpm.CheckContent(ctx, content); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	aid, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := s.accountQuery.GetByID(ctx, aid)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	p, err := s.replyQuerier.Get(ctx, threadID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := acc.Roles.Permissions().Authorise(ctx, func() error {
		if p.Author.ID != aid {
			return fault.Wrap(rbac.ErrPermissions,
				fctx.With(ctx),
				fmsg.WithDesc("not owner", "You are not the owner of the post and do not have the Manage Posts permission."))
		}
		return nil
	}, rbac.PermissionManagePosts); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	opts := partial.Opts()

	pref, err := s.replyQuerier.Probe(ctx, threadID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.cache.Invalidate(ctx, xid.ID(pref.RootPostID)); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	p, err = s.replyWriter.Update(ctx, threadID, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &message.EventThreadReplyUpdated{
		ThreadID: p.RootPostID,
		ReplyID:  p.ID,
	})

	return p, nil
}
