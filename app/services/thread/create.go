package thread

import (
	"context"

	"github.com/el-mike/restrict"
	"github.com/pkg/errors"

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
) (*thread.Thread, error) {
	acc, err := s.account_repo.GetByID(ctx, authorID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}

	if err := s.rbac.Authorize(&restrict.AccessRequest{
		Subject:  acc,
		Resource: &thread.Thread{},
		Actions:  []string{"create"},
	}); err != nil {
		return nil, errors.Wrap(err, "failed to authorize")
	}

	thr, err := s.thread_repo.Create(ctx, title, body, authorID, categoryID, tags)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create thread")
	}

	return thr, nil
}
