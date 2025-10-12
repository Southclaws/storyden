package thread

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
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

	session := session.GetOptAccountID(ctx)

	thr, err := s.threadQuerier.Get(ctx, threadID, pageParams, session)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to get thread"))
	}

	if thr.Visibility != visibility.VisibilityPublished {
		accountID, ok := session.Get()
		if !ok {
			return nil, fault.Wrap(ErrNoPermission, fctx.With(ctx))
		}

		acc, err := s.accountQuery.GetByID(ctx, accountID)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if err := acc.Roles.Permissions().Authorise(ctx, func() error {
			if thr.Author.ID == accountID {
				return nil
			}

			return ErrNoPermission
		}, rbac.PermissionManagePosts); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	// recommendations, err := s.recommender.Recommend(ctx, thr)
	// if err != nil {
	// 	s.l.Warn("failed to aggregate recommendations", slog.String("error", err.Error()))
	// } else {
	// 	thr.Related = append(thr.Related, recommendations...)
	// }

	return thr, nil
}
