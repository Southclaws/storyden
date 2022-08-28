package thread

import (
	"context"
	"time"

	"github.com/Southclaws/storyden/pkg/resources/thread"
)

func (s *service) ListAll(
	ctx context.Context,
	before time.Time,
	max int,
) ([]*thread.Thread, error) {
	return s.thread_repo.List(ctx, before, max)
}
