package react_manager

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reaction"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Reactor struct {
	accountQuerier *account_querier.Querier
	reactWriter    *reaction.Writer
	reactReader    *reaction.Querier
	bus            *pubsub.Bus
}

func New(
	accountQuerier *account_querier.Querier,
	reactWriter *reaction.Writer,
	reactReader *reaction.Querier,
	bus *pubsub.Bus,
) *Reactor {
	return &Reactor{
		accountQuerier: accountQuerier,
		reactWriter:    reactWriter,
		reactReader:    reactReader,
		bus:            bus,
	}
}

func (s *Reactor) Add(ctx context.Context, postID post.ID, emoji string) (*reaction.React, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := s.reactWriter.Add(ctx, accountID, xid.ID(postID), emoji)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &message.EventPostReacted{PostID: postID})

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

	err = s.reactWriter.Remove(ctx, accountID, reactID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &message.EventPostUnreacted{
		PostID: targetID,
	})

	return nil
}
