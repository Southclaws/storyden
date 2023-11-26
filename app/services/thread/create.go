package thread

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/el-mike/restrict"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/thread"
)

func (s *service) Create(ctx context.Context,
	title string,
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

	// BUG: partial does not have body set
	opts = append(opts, s.hydrate(ctx, partial)...)

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

	return thr, nil
}

func (s *service) hydrate(ctx context.Context, partial Partial) (opts []thread.Option) {
	body, bodyOK := partial.Body.Get()
	url, urlOK := partial.URL.Get()

	if !bodyOK && !urlOK {
		return
	}

	return s.hydrator.HydrateThread(ctx, body, url)
}
