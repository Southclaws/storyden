package react

import (
	"context"
	"errors"

	"github.com/Southclaws/storyden/internal/infrastructure/db/model"
	"github.com/rs/xid"
)

var (
	ErrInvalidEmoji   = errors.New("invalid emoji codepoint")
	ErrAlreadyReacted = errors.New("account already reacted emoji to post")
	ErrUnauthorised   = errors.New("not allowed to remove another account's react")
)

type database struct {
	db *model.Client
}

func New(db *model.Client) Repository {
	return &database{db}
}

func (d *database) Add(ctx context.Context, accountID xid.ID, postID xid.ID, emojiID string) (*React, error) {
	e, ok := IsValidEmoji(emojiID)
	if !ok {
		return nil, ErrInvalidEmoji
	}

	react, err := d.db.React.
		Create().
		SetEmoji(e).
		SetAccountID(xid.ID(accountID)).
		SetPostID(xid.ID(postID)).
		Save(ctx)
	if err != nil {
		if model.IsConstraintError(err) {
			return nil, ErrAlreadyReacted
		}

		return nil, err
	}

	return FromModel(react), nil
}

func (d *database) Remove(ctx context.Context, accountID xid.ID, reactID ReactID) (*React, error) {
	// First, look up the react to check if this account has permissions to remove.
	p, err := d.db.React.Get(ctx, xid.ID(reactID))
	if err != nil {
		return nil, err
	}

	if !p.Edges.Account.Admin && p.Edges.Account.ID != xid.ID(accountID) {
		return nil, ErrUnauthorised
	}

	// the account has permission, remove it.
	if err := d.db.React.DeleteOneID(xid.ID(reactID)).Exec(ctx); err != nil {
		return nil, err
	}

	return FromModel(p), nil
}
