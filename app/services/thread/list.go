package thread

import (
	"context"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

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
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to list threads"))
	}

	return thr, nil
}
