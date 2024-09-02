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

	acc, err := s.accountQuery.GetByID(ctx, aid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	// TODO: Minimal reader interface for thread.
	thr, err := s.thread_repo.Get(ctx, id, nil)
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

	// if err := s.indexQueue.Publish(ctx, mq.DeindexPost{
	// 	ID: thr.ID,
	// }); err != nil {
	// 	s.l.Error("failed to publish index post message", zap.Error(err))
	// }

	return nil
}
