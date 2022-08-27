package thread

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/el-mike/restrict"

	"github.com/Southclaws/storyden/pkg/resources/account"
	"github.com/Southclaws/storyden/pkg/resources/category"
	"github.com/Southclaws/storyden/pkg/resources/thread"
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
		return nil, fault.WithValue(err, "failed to get account", "authorID", authorID.String())
	}

	if err := s.rbac.Authorize(&restrict.AccessRequest{
		Subject:  acc,
		Resource: &thread.Thread{},
		Actions:  []string{"create"},
	}); err != nil {
		return nil, err
	}

	thr, err := s.thread_repo.Create(ctx, title, body, authorID, categoryID, tags)
	if err != nil {
		return nil, err
	}

	return thr, nil
}
