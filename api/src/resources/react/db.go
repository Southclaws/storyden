package react

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/Southclaws/storyden/api/src/infra/db/model"
)

var (
	ErrInvalidEmoji   = errors.New("invalid emoji codepoint")
	ErrAlreadyReacted = errors.New("user already reacted emoji to post")
	ErrUnauthorised   = errors.New("not allowed to remove another user's react")
)

type database struct {
	db *model.Client
}

func New(db *model.Client) Repository {
	return &database{db}
}

func (d *database) Add(ctx context.Context, userID uuid.UUID, postID uuid.UUID, emojiID string) (*React, error) {
	e, ok := IsValidEmoji(emojiID)
	if !ok {
		return nil, ErrInvalidEmoji
	}

	react, err := d.db.React.
		Create().
		SetEmoji(e).
		SetUserID(uuid.UUID(userID)).
		SetPostID(uuid.UUID(postID)).
		Save(ctx)
	if err != nil {
		if model.IsConstraintError(err) {
			return nil, ErrAlreadyReacted
		}
		return nil, err
	}

	return FromModel(react), nil
}

func (d *database) Remove(ctx context.Context, userID uuid.UUID, reactID ReactID) (*React, error) {
	// First, look up the react to check if this user has permissions to remove.
	p, err := d.db.React.Get(ctx, uuid.UUID(reactID))
	if err != nil {
		return nil, err
	}

	if !p.Edges.User.Admin && p.Edges.User.ID != uuid.UUID(userID) {
		return nil, ErrUnauthorised
	}

	// the user has permission, remove it.
	if err := d.db.React.DeleteOneID(uuid.UUID(reactID)).Exec(ctx); err != nil {
		return nil, err
	}

	return FromModel(p), nil
}
