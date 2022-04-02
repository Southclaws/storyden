package react

import (
	"context"
	"errors"
	"strings"

	"github.com/Southclaws/storyden/api/src/infra/db"
	"github.com/Southclaws/storyden/ent"
)

var (
	ErrInvalidEmoji   = errors.New("invalid emoji codepoint")
	ErrAlreadyReacted = errors.New("user already reacted emoji to post")
	ErrUnauthorised   = errors.New("not allowed to remove another user's react")
)

type DB struct {
	db *ent.Client
}

func New(db *ent.Client) Repository {
	return &DB{db}
}

func (d *DB) Add(ctx context.Context, userID, postID, emojiID string) (*React, error) {
	e, ok := IsValidEmoji(emojiID)
	if !ok {
		return nil, ErrInvalidEmoji
	}

	react, err := d.db.React.
		CreateOne(
			db.React.Emoji.Set(e),
			db.React.User.Link(
				db.User.ID.Equals(userID),
			),
			db.React.Post.Link(
				db.Post.ID.Equals(postID),
			),
		).
		With(
			db.React.User.Fetch(),
			db.React.Post.Fetch(),
		).
		Exec(ctx)
	if err != nil {
		// NOTE: This depends on internal Prisma error implementation details
		// and may change without notice. This will be reflected by the frontend
		// showing an Internal Server Error when a user attempts to react to a
		// post with an emoji they have already reacted with. This comment
		// signposts this possibility.
		// TODO: Clean this up when Prisma's Go client adds better errors.
		if strings.Contains(err.Error(), "UniqueConstraintViolation") {
			return nil, ErrAlreadyReacted
		}
		return nil, err
	}

	_, err = d.db.Post.
		FindUnique(db.Post.ID.Equals(postID)).
		Update(
			db.Post.Reacts.Link(
				db.React.ID.Equals(react.ID),
			),
		).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return FromModel(react, react.RelationsReact.Post.ID), nil
}

func (d *DB) Remove(ctx context.Context, userID, reactID string) (*React, error) {
	// First, look up the react to check if this user has permissions to remove.
	p, err := d.db.React.
		FindUnique(db.React.ID.Equals(reactID)).
		With(
			db.React.User.Fetch(),
			db.React.Post.Fetch(),
		).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	if !p.RelationsReact.User.Admin && p.UserID != userID {
		return nil, ErrUnauthorised
	}

	// the user has permission, remove it.
	if _, err = d.db.React.
		FindUnique(db.React.ID.Equals(reactID)).
		Delete().
		Exec(ctx); err != nil {
		return nil, err
	}

	return FromModel(p, p.RelationsReact.Post.ID), nil
}
