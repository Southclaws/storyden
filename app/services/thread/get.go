package thread

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/thread"
)

func (s *service) Get(
	ctx context.Context,
	threadID post.ID,
) (*thread.Thread, error) {
	thr, err := s.thread_repo.Get(ctx, threadID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to get thread"))
	}

	recommendations, err := s.recommender.Recommend(ctx, thr)
	if err != nil {
		s.l.Warn("failed to aggregate recommendations", zap.Error(err))
	} else {
		thr.Related = append(thr.Related, recommendations...)
	}

	return thr, nil
}
