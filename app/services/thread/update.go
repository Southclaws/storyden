package thread

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/el-mike/restrict"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/thread"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

func (s *service) Update(ctx context.Context, threadID post.ID, partial Partial) (*thread.Thread, error) {
	aid, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := s.account_repo.GetByID(ctx, aid)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	thr, err := s.thread_repo.Get(ctx, threadID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.rbac.Authorize(&restrict.AccessRequest{
		Subject:  acc,
		Resource: thr,
		Actions:  []string{rbac.ActionUpdate},
	}); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to authorize"))
	}

	opts := partial.Opts()

	opts = append(opts, s.hydrate(ctx, partial)...)

	thr, err = s.thread_repo.Update(ctx, threadID, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return thr, nil
}
