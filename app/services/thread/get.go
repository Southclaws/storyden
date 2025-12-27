package thread

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

var ErrNoPermission = fault.New("unauthenticated user cannot view unpublished threads", ftag.With(ftag.PermissionDenied))

func (s *service) Get(
	ctx context.Context,
	threadID post.ID,
	pageParams pagination.Parameters,
) (*thread.Thread, error) {
	ctx, span := s.ins.Instrument(ctx)
	defer span.End()

	accountID := session.GetOptAccountID(ctx)

	thr, err := s.threadQuerier.Get(ctx, threadID, pageParams, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to get thread"))
	}

	if thr.Visibility != visibility.VisibilityPublished {
		aid, ok := accountID.Get()
		if !ok {
			return nil, fault.Wrap(ErrNoPermission, fctx.With(ctx))
		}

		if err := session.Authorise(ctx, func() error {
			if thr.Author.ID == aid {
				return nil
			}

			return ErrNoPermission
		}, rbac.PermissionManagePosts); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	thr.Replies.Items = s.filterRepliesByVisibility(ctx, thr.Replies.Items, accountID)

	// recommendations, err := s.recommender.Recommend(ctx, thr)
	// if err != nil {
	// 	s.l.Warn("failed to aggregate recommendations", slog.String("error", err.Error()))
	// } else {
	// 	thr.Related = append(thr.Related, recommendations...)
	// }

	return thr, nil
}

func (s *service) filterRepliesByVisibility(ctx context.Context, replies []*reply.Reply, accountID opt.Optional[account.AccountID]) []*reply.Reply {
	if accountID.Ok() && session.GetRoles(ctx).Permissions().HasAny(rbac.PermissionManagePosts, rbac.PermissionAdministrator) {
		return replies
	}

	return dt.Filter(replies, func(r *reply.Reply) bool {
		if r.Visibility == visibility.VisibilityPublished {
			return true
		}

		if r.Visibility == visibility.VisibilityReview {
			if aid, ok := accountID.Get(); ok {
				return r.Author.ID == aid
			}
		}

		return false
	})
}
