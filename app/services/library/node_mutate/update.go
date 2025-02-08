package node_mutate

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/library/node_auth"
)

type Updated struct {
	library.Node
	TitleSuggestion   opt.Optional[string]
	TagSuggestions    opt.Optional[tag_ref.Names]
	ContentSuggestion opt.Optional[datagraph.Content]
}

func (s *Manager) Update(ctx context.Context, qk library.QueryKey, p Partial) (*Updated, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := s.accountQuery.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	n, err := s.nodeQuerier.Get(ctx, qk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := node_auth.AuthoriseNodeMutation(ctx, acc, n); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	pre, err := s.preMutation(ctx, p, opt.NewPtr(n))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	n, err = s.nodeWriter.Update(ctx, qk, pre.opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	post, err := s.postMutation(ctx, n, pre)
	if err != nil {
		// TODO: Does this need to error?
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if post.properties != nil {
		n.Properties = post.properties
	}

	if n.Visibility == visibility.VisibilityPublished {
		if err := s.indexQueue.Publish(ctx, mq.IndexNode{
			ID: library.NodeID(n.Mark.ID()),
		}); err != nil {
			s.logger.Error("failed to publish index post message", zap.Error(err))
		}
	} else {
		if err := s.deleteQueue.Publish(ctx, mq.DeleteNode{
			ID: library.NodeID(n.GetID()),
		}); err != nil {
			// failing to publish the deletion message is worthy of an error.
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	s.fetcher.HydrateContentURLs(ctx, n)

	u := Updated{
		Node:              *n,
		TagSuggestions:    pre.tags,
		TitleSuggestion:   pre.title,
		ContentSuggestion: pre.content,
	}

	return &u, nil
}
