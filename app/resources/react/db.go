package react

import (
	"context"

	"github.com/Southclaws/fault/errctx"
	"github.com/Southclaws/fault/errtag"
	"github.com/pkg/errors"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/infrastructure/db/model"
)

var ErrInvalidEmoji = errors.New("invalid emoji codepoint")

type database struct {
	db *model.Client
}

func New(db *model.Client) Repository {
	return &database{db}
}

func (d *database) Add(ctx context.Context, accountID xid.ID, postID xid.ID, emojiID string) (*React, error) {
	e, ok := IsValidEmoji(emojiID)
	if !ok {
		return nil, errtag.Wrap(errctx.Wrap(ErrInvalidEmoji, ctx), errtag.InvalidArgument{})
	}

	react, err := d.db.React.
		Create().
		SetEmoji(e).
		SetAccountID(xid.ID(accountID)).
		SetPostID(xid.ID(postID)).
		Save(ctx)
	if err != nil {
		if model.IsConstraintError(err) {
			return nil, errtag.Wrap(
				errctx.Wrap(errors.Wrap(err, "already reacted to post"), ctx),
				errtag.AlreadyExists{})
		}

		return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.Internal{})
	}

	return FromModel(react), nil
}

func (d *database) Remove(ctx context.Context, accountID xid.ID, reactID ReactID) (*React, error) {
	// First, look up the react to check if this account has permissions to remove.
	p, err := d.db.React.Get(ctx, xid.ID(reactID))
	if err != nil {
		return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.Internal{})
	}

	if !p.Edges.Account.Admin && p.Edges.Account.ID != xid.ID(accountID) {
		return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.PermissionDenied{})
	}

	// the account has permission, remove it.
	if err := d.db.React.DeleteOneID(xid.ID(reactID)).Exec(ctx); err != nil {
		return nil, errtag.Wrap(errctx.Wrap(err, ctx), errtag.Internal{})
	}

	return FromModel(p), nil
}
