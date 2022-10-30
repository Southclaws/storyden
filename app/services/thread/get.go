package thread

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/thread"
)

func (s *service) Get(
	ctx context.Context,
	threadID post.PostID,
) (*thread.Thread, error) {
	thr, err := s.thread_repo.Get(ctx, threadID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to get thread"))
	}

	return thr, nil
}
