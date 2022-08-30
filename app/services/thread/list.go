package thread

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/Southclaws/storyden/app/resources/thread"
)

func (s *service) ListAll(
	ctx context.Context,
	before time.Time,
	max int,
) ([]*thread.Thread, error) {
	thr, err := s.thread_repo.List(ctx, before, max)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list threads")
	}

	return thr, nil
}
