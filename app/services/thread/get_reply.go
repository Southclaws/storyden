package thread

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

func (s *service) GetReply(ctx context.Context, threadID post.ID, replyID post.ID) (*reply.Reply, error) {
	root, err := s.threadQuerier.Probe(ctx, threadID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if root.DeletedAt.Ok() {
		return nil, fault.New("thread not found", ftag.With(ftag.NotFound))
	}

	if err := authoriseThreadRead(ctx, &root.Post); err != nil {
		return nil, err
	}

	accountID := session.GetOptAccountID(ctx)
	r, err := s.threadQuerier.GetReply(ctx, threadID, replyID, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if len(s.filterRepliesByVisibility(ctx, []*reply.Reply{r}, accountID)) == 0 {
		return nil, fault.New("reply is not visible", ftag.With(ftag.PermissionDenied))
	}

	if target, ok := r.ReplyTo.Get(); ok {
		if len(s.filterRepliesByVisibility(ctx, []*reply.Reply{&target}, accountID)) == 0 {
			r.ReplyTo = opt.NewEmpty[reply.Reply]()
		}
	}

	return r, nil
}
