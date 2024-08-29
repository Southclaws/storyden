package thread

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/el-mike/restrict"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/visibility"
)

func (s *service) Create(ctx context.Context,
	title string,
	authorID account.AccountID,
	categoryID category.CategoryID,
	status visibility.Visibility,
	tags []string,
	meta map[string]any,
	partial Partial,
) (*thread.Thread, error) {
	acc, err := s.accountQuery.GetByID(ctx, authorID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to get account"))
	}

	if err := s.rbac.Authorize(&restrict.AccessRequest{
		Subject:  acc,
		Resource: &thread.Thread{},
		Actions:  []string{rbac.ActionCreate},
	}); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to authorize"))
	}

	opts := partial.Opts()
	opts = append(opts,
		thread.WithVisibility(status),
		thread.WithMeta(meta),
	)

	if u, ok := partial.URL.Get(); ok {
		ln, err := s.fetcher.Fetch(ctx, u)
		if err == nil {
			opts = append(opts, thread.WithLink(xid.ID(ln.ID)))
		}
	}

	thr, err := s.thread_repo.Create(ctx,
		title,
		authorID,
		categoryID,
		tags,
		opts...,
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create thread"))
	}

	if err := s.indexQueue.Publish(ctx, mq.IndexPost{
		ID: thr.ID,
	}); err != nil {
		s.l.Error("failed to publish index post message", zap.Error(err))
	}

	s.fetcher.HydrateContentURLs(ctx, thr)

	return thr, nil
}
