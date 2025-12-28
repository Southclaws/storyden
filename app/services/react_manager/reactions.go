package react_manager

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/post_querier"
	"github.com/Southclaws/storyden/app/resources/post/reaction"
	"github.com/Southclaws/storyden/app/resources/post/thread_cache"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Reactor struct {
	accountQuerier *account_querier.Querier
	reactWriter    *reaction.Writer
	reactReader    *reaction.Querier
	postQuerier    *post_querier.Querier
	bus            *pubsub.Bus
	cache          *thread_cache.Cache
}

func New(
	accountQuerier *account_querier.Querier,
	reactWriter *reaction.Writer,
	reactReader *reaction.Querier,
	postQuerier *post_querier.Querier,
	bus *pubsub.Bus,
	cache *thread_cache.Cache,
) *Reactor {
	return &Reactor{
		accountQuerier: accountQuerier,
		reactWriter:    reactWriter,
		reactReader:    reactReader,
		postQuerier:    postQuerier,
		bus:            bus,
		cache:          cache,
	}
}

func (s *Reactor) Add(ctx context.Context, postID post.ID, emoji string) (*reaction.React, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	pref, err := s.postQuerier.Probe(ctx, postID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.cache.Invalidate(ctx, xid.ID(pref.Root)); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := s.reactWriter.Add(ctx, accountID, xid.ID(postID), emoji)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &message.EventPostReacted{
		PostID:     postID,
		RootPostID: pref.Root,
	})

	return r, nil
}

func (s *Reactor) Remove(ctx context.Context, reactID reaction.ReactID) error {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := s.accountQuerier.GetByID(ctx, accountID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	reac, err := s.reactReader.Get(ctx, reactID)
	if err != nil {
		if ftag.Get(err) == ftag.NotFound {
			return nil
		}
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := acc.Roles.Permissions().Authorise(ctx, func() error {
		if reac.Author.ID != accountID {
			return fault.New("not owner of reaction")
		}
		return nil
	}); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	targetID := post.ID(reac.Target())

	pref, err := s.postQuerier.Probe(ctx, targetID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.cache.Invalidate(ctx, xid.ID(pref.Root)); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	err = s.reactWriter.Remove(ctx, accountID, reactID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &message.EventPostUnreacted{
		PostID:     targetID,
		RootPostID: pref.Root,
	})

	return nil
}
