package post

import "context"

type Repository interface {
	CreatePost(ctx context.Context, body, authorID, parentID, replyToID string) (*Post, error)
	GetPosts(ctx context.Context, slug string, max, skip int, deleted bool, admin bool) ([]Post, error)
	EditPost(ctx context.Context, authorID, postID string, title *string, body *string) (*Post, error)
	DeletePost(ctx context.Context, authorID, postID string, force bool) (*Post, error)
}
