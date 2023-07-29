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
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/thread"
)

type Params struct {
	AccountID  opt.Optional[account.AccountID]
	Tags       opt.Optional[[]xid.ID]
	Categories opt.Optional[[]string]
}

func (s *service) ListAll(
	ctx context.Context,
	before time.Time,
	max int,
	opts Params,
) ([]*thread.Thread, error) {
	q := []thread.Query{
		// User's drafts are always private so we always filter published only.
		thread.HasStatus(post.StatusPublished),
	}

	opts.AccountID.Call(func(a account.AccountID) { q = append(q, thread.HasAuthor(a)) })
	opts.Tags.Call(func(a []xid.ID) { q = append(q, thread.HasTags(a)) })
	opts.Categories.Call(func(a []string) { q = append(q, thread.HasCategories(a)) })

	thr, err := s.thread_repo.List(ctx, before, max, q...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to list threads"))
	}

	return thr, nil
}
