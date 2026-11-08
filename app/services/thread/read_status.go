package thread

import (
	"context"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

func (s *service) UpdateReadStatus(ctx context.Context, threadID post.ID, unread bool, at time.Time) (*post.ReadStatus, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if _, err := s.Get(ctx, threadID, pagination.Parameters{}); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if unread {
		err = s.readStateWriter.MarkUnread(ctx, accountID, threadID, at)
	} else {
		err = s.readStateWriter.MarkRead(ctx, accountID, threadID, at)
	}
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return s.threadQuerier.GetReadStatus(ctx, threadID, accountID)
}
