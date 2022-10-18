package thread

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/xid"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/thread"
)

type Params struct {
	AccountID opt.Optional[account.AccountID]
	Tags      opt.Optional[[]xid.ID]
}

func (s *service) ListAll(
	ctx context.Context,
	before time.Time,
	max int,
	opts Params,
) ([]*thread.Thread, error) {
	q := []thread.Query{}

	opts.AccountID.If(func(a account.AccountID) { q = append(q, thread.WithAuthor(a)) })
	opts.Tags.If(func(a []xid.ID) { q = append(q, thread.WithTags(a)) })

	thr, err := s.thread_repo.List(ctx, before, max, q...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list threads")
	}

	return thr, nil
}
