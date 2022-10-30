package thread

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/el-mike/restrict"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/resources/thread"
)

func (s *service) Create(ctx context.Context,
	title string,
	body string,
	authorID account.AccountID,
	categoryID category.CategoryID,
	tags []string,
	meta map[string]any,
) (*thread.Thread, error) {
	acc, err := s.account_repo.GetByID(ctx, authorID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to get account"))
	}

	if err := s.rbac.Authorize(&restrict.AccessRequest{
		Subject:  acc,
		Resource: &thread.Thread{},
		Actions:  []string{"create"},
	}); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to authorize"))
	}

	thr, err := s.thread_repo.Create(ctx, title, body, authorID, categoryID, tags, thread.WithMeta(meta))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create thread"))
	}

	return thr, nil
}
