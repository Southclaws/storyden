package thread

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/el-mike/restrict"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

func (s *service) Delete(ctx context.Context, id post.ID) error {
	aid, err := session.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := s.account_repo.GetByID(ctx, aid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	thr, err := s.thread_repo.Get(ctx, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.rbac.Authorize(&restrict.AccessRequest{
		Subject:  acc,
		Resource: thr,
		Actions:  []string{rbac.ActionDelete},
	}); err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to authorize"))
	}

	err = s.thread_repo.Delete(ctx, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create thread"))
	}

	return nil
}
