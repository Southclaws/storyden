package react

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/pkg/errors"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
)

var ErrInvalidEmoji = errors.New("invalid emoji codepoint")

type database struct {
	db *ent.Client
}

func New(db *ent.Client) Repository {
	return &database{db}
}

func (d *database) Add(ctx context.Context, accountID account.AccountID, postID xid.ID, emojiID string) (*React, error) {
	e, ok := IsValidEmoji(emojiID)
	if !ok {
		return nil, fault.Wrap(ErrInvalidEmoji, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	react, err := d.db.React.
		Create().
		SetEmoji(e).
		SetAccountID(xid.ID(accountID)).
		SetPostID(xid.ID(postID)).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err,
				fctx.With(ctx),
				ftag.With(ftag.AlreadyExists),
				fmsg.WithDesc("already reacted", "You have already reacted to this post."),
			)
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return FromModel(react), nil
}

func (d *database) Remove(ctx context.Context, accountID account.AccountID, reactID ReactID) (*React, error) {
	// First, look up the react to check if this account has permissions to remove.
	p, err := d.db.React.Get(ctx, xid.ID(reactID))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	if !p.Edges.Account.Admin && p.Edges.Account.ID != xid.ID(accountID) {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.PermissionDenied))
	}

	// the account has permission, remove it.
	if err := d.db.React.DeleteOneID(xid.ID(reactID)).Exec(ctx); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return FromModel(p), nil
}
