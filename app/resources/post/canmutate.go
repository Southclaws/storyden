package post

import (
	"context"

	"github.com/rs/xid"

	"github.com/Southclaws/fault/errctx"
	"github.com/Southclaws/fault/errtag"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/post"
)

func CanUserMutatePost(ctx context.Context, d *model.Client, authorID, postID PostID) error {
	// First, check if this user is the author of the post.
	post, err := d.Post.Query().Where(post.IDEQ(xid.ID(postID))).WithAuthor().Only(ctx)
	if err != nil {
		return errtag.Wrap(errctx.Wrap(err, ctx), errtag.Internal{})
	}

	user := post.Edges.Author
	if PostID(user.ID) == authorID {
		return nil
	}

	if user.Admin {
		return nil
	}

	// Not either? Not authorised to edit.
	return errtag.Wrap(errctx.Wrap(ErrUnauthorised, ctx), errtag.PermissionDenied{})
}
