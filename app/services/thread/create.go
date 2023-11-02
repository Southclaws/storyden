package thread

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/el-mike/restrict"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/thread"
)

func (s *service) Create(ctx context.Context,
	title string,
	body string,
	authorID account.AccountID,
	categoryID category.CategoryID,
	status post.Status,
	tags []string,
	meta map[string]any,
	partial Partial,
) (*thread.Thread, error) {
	acc, err := s.account_repo.GetByID(ctx, authorID)
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
		thread.WithStatus(status),
		thread.WithMeta(meta),
	)

	opts = append(opts, s.hydrateLink(ctx, partial)...)

	thr, err := s.thread_repo.Create(ctx,
		title,
		body,
		authorID,
		categoryID,
		tags,
		opts...,
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create thread"))
	}

	return thr, nil
}

func (s *service) hydrateLink(ctx context.Context, partial Partial) (opts []thread.Option) {
	v, ok := partial.URL.Get()
	if !ok {
		return
	}

	opts, err := s.hydrator.HydrateThread(ctx, v)
	if err != nil {
		s.l.Warn("failed to hydrate URL",
			zap.String("url", v),
			zap.Error(err))
	}

	return
}
