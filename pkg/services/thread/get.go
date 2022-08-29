package thread

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Southclaws/storyden/pkg/resources/post"
	"github.com/Southclaws/storyden/pkg/resources/thread"
)

func (s *service) Get(
	ctx context.Context,
	threadID post.PostID,
) (*thread.Thread, error) {
	thr, err := s.thread_repo.Get(ctx, threadID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get thread")
	}

	return thr, nil
}
