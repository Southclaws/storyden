package post

import (
	"context"

	"github.com/Southclaws/storyden/api/src/infra/db"
	"github.com/Southclaws/storyden/api/src/infra/db/model"
)

func CanUserMutatePost(ctx context.Context, d *model.Client, authorID, id string) error {
	// First, check if this user is the author of the post.
	post, err := d.Post.
		FindUnique(db.Post.ID.Equals(id)).
		With(db.Post.Author.Fetch()).
		Exec(ctx)
	if err != nil {
		return err
	}
	if post.Author().ID == authorID {
		return nil
	}

	// If they are not the author, check if they are an admin.
	user, err := d.User.
		FindUnique(db.User.ID.Equals(authorID)).
		Exec(ctx)
	if err != nil {
		return err
	}
	if user.Admin {
		return nil
	}

	// Not either? Not authorised to edit.
	return ErrUnauthorised
}
