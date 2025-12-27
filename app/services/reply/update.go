package reply

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/reply_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

func (s *Mutator) Update(ctx context.Context, replyID post.ID, partial Partial) (*reply.Reply, error) {
	aid, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	p, err := s.replyQuerier.Get(ctx, replyID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// Check if user can update this reply (owner or has ManagePosts permission)
	if err := session.Authorise(ctx, func() error {
		if p.Author.ID != aid {
			return fault.Wrap(rbac.ErrPermissions,
				fctx.With(ctx),
				fmsg.WithDesc("not owner", "You are not the owner of the post and do not have the Manage Posts permission."))
		}
		return nil
	}, rbac.PermissionManagePosts); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// Check if user is trying to change visibility - only post managers can do this
	userSetVisibility := false
	if _, ok := partial.Visibility.Get(); ok {
		roles := session.GetRoles(ctx)
		if !roles.Permissions().HasAny(rbac.PermissionManagePosts, rbac.PermissionAdministrator) {
			return nil, fault.Wrap(rbac.ErrPermissions,
				fctx.With(ctx),
				fmsg.WithDesc("visibility change denied", "Only users with Manage Posts permission can change post visibility."))
		}
		userSetVisibility = true
	}

	oldVisibility := p.Visibility
	opts := partial.Opts()

	if content, ok := partial.Content.Get(); ok && !userSetVisibility {
		result, err := s.cpm.CheckContent(ctx, xid.ID(replyID), datagraph.KindReply, "", content)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if result.RequiresReview {
			opts = append(opts, reply_writer.WithVisibility(visibility.VisibilityReview))
		}
	}

	pref, err := s.replyQuerier.Probe(ctx, replyID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.cache.Invalidate(ctx, xid.ID(pref.RootPostID)); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	p, err = s.replyWriter.Update(ctx, replyID, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &message.EventThreadReplyUpdated{
		ThreadID: p.RootPostID,
		ReplyID:  p.ID,
	})

	// Emit visibility-specific events when visibility changes
	if oldVisibility != p.Visibility {
		if p.Visibility == visibility.VisibilityPublished {
			s.bus.Publish(ctx, &message.EventThreadReplyPublished{
				ThreadID: p.RootPostID,
				ReplyID:  p.ID,
			})
		} else if oldVisibility == visibility.VisibilityPublished {
			s.bus.Publish(ctx, &message.EventThreadReplyUnpublished{
				ThreadID: p.RootPostID,
				ReplyID:  p.ID,
			})
		}
	}

	return p, nil
}
