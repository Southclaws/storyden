package thread

import (
	"context"

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
	return nil, nil
}
