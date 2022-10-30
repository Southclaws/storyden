package post

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/post"
)

func CanUserMutatePost(ctx context.Context, d *model.Client, authorID, postID PostID) error {
	// First, check if this user is the author of the post.
	post, err := d.Post.Query().Where(post.IDEQ(xid.ID(postID))).WithAuthor().Only(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	user := post.Edges.Author
	if PostID(user.ID) == authorID {
		return nil
	}

	if user.Admin {
		return nil
	}

	// Not either? Not authorised to edit.
	return fault.Wrap(ErrUnauthorised, fctx.With(ctx), ftag.With(ftag.PermissionDenied))
}
