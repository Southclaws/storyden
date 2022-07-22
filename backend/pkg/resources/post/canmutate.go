package post

import (
	"context"

	"github.com/google/uuid"

	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/post"
)

func CanUserMutatePost(ctx context.Context, d *model.Client, authorID, postID PostID) error {
	// First, check if this user is the author of the post.
	post, err := d.Post.Query().Where(post.IDEQ(uuid.UUID(postID))).WithAuthor().Only(ctx)
	if err != nil {
		return err
	}

	user := post.Edges.Author
	if PostID(user.ID) == authorID {
		return nil
	}

	if user.Admin {
		return nil
	}

	// Not either? Not authorised to edit.
	return ErrUnauthorised
}
