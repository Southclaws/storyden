package post

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/el-mike/restrict"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication"
)

func (s *service) Delete(ctx context.Context, postID post.PostID) error {
	aid, err := authentication.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := s.account_repo.GetByID(ctx, aid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	p, err := s.post_repo.Get(ctx, postID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.rbac.Authorize(&restrict.AccessRequest{
		Subject:  acc,
		Resource: p,
		Actions:  []string{rbac.ActionDelete},
	}); err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to authorize"))
	}

	err = s.post_repo.Delete(ctx, postID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
