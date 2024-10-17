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
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

type Params struct {
	Query         opt.Optional[string]
	CreatedBefore opt.Optional[time.Time]
	UpdatedBefore opt.Optional[time.Time]
	AccountID     opt.Optional[account.AccountID]
	Visibility    opt.Optional[[]visibility.Visibility]
	Tags          opt.Optional[[]xid.ID]
	Categories    opt.Optional[[]string]
}

func (s *service) List(ctx context.Context,
	page int,
	size int,
	opts Params,
) (*thread.Result, error) {
	accountID := session.GetOptAccountID(ctx)

	q := []thread.Query{
		thread.HasNotBeenDeleted(),
	}

	opts.Query.Call(func(value string) { q = append(q, thread.HasKeyword(value)) })
	opts.CreatedBefore.Call(func(value time.Time) { q = append(q, thread.HasCreatedDateBefore(value)) })
	opts.UpdatedBefore.Call(func(value time.Time) { q = append(q, thread.HasUpdatedDateBefore(value)) })
	opts.AccountID.Call(func(a account.AccountID) { q = append(q, thread.HasAuthor(a)) })
	opts.Tags.Call(func(a []xid.ID) { q = append(q, thread.HasTags(a)) })
	opts.Categories.Call(func(a []string) { q = append(q, thread.HasCategories(a)) })

	vq := func() thread.Query {
		v, ok := opts.Visibility.Get()
		if !ok {
			return thread.HasStatus(visibility.VisibilityPublished)
		}

		onlyRequestingPublished := len(v) == 1 && v[0] == visibility.VisibilityPublished
		if onlyRequestingPublished {
			return thread.HasStatus(visibility.VisibilityPublished)
		}

		filterByAccount, ok := opts.AccountID.Get()
		if !ok {
			return thread.HasStatus(visibility.VisibilityPublished)
		}

		requestedByAccount, ok := accountID.Get()
		if !ok {
			return thread.HasStatus(visibility.VisibilityPublished)
		}

		requestingOwnThreads := filterByAccount == requestedByAccount

		if !requestingOwnThreads {
			thread.HasStatus(visibility.VisibilityPublished)
		}

		return thread.HasStatus(v...)
	}()
	q = append(q, vq)

	thr, err := s.thread_repo.List(ctx, page, size, accountID, q...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to list threads"))
	}

	return thr, nil
}
