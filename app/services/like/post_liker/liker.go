package post_liker

import (
	"context"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/like/like_writer"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/post_querier"
	"github.com/Southclaws/storyden/app/resources/post/thread_cache"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type PostLiker struct {
	likeWriter  *like_writer.LikeWriter
	postQuerier *post_querier.Querier
	bus         *pubsub.Bus
	cache       *thread_cache.Cache
}

func New(
	likeWriter *like_writer.LikeWriter,
	postQuerier *post_querier.Querier,
	bus *pubsub.Bus,
	cache *thread_cache.Cache,
) *PostLiker {
	return &PostLiker{
		likeWriter:  likeWriter,
		postQuerier: postQuerier,
		bus:         bus,
		cache:       cache,
	}
}

func (l *PostLiker) AddPostLike(ctx context.Context, accountID account.AccountID, postID post.ID) error {
	postRef, err := l.postQuerier.Probe(ctx, postID)
	if err != nil {
		return err
	}

	if err := l.cache.Invalidate(ctx, xid.ID(postRef.Root)); err != nil {
		return err
	}

	err = l.likeWriter.AddPostLike(ctx, accountID, postID)
	if err != nil {
		return err
	}

	l.bus.Publish(ctx, &message.EventPostLiked{
		PostID:     postID,
		RootPostID: postRef.Root,
	})

	return nil
}

func (l *PostLiker) RemovePostLike(ctx context.Context, accountID account.AccountID, postID post.ID) error {
	postRef, err := l.postQuerier.Probe(ctx, postID)
	if err != nil {
		return err
	}

	if err := l.cache.Invalidate(ctx, xid.ID(postRef.Root)); err != nil {
		return err
	}

	err = l.likeWriter.RemovePostLike(ctx, accountID, postID)
	if err != nil {
		return err
	}

	l.bus.Publish(ctx, &message.EventPostUnliked{
		PostID:     postID,
		RootPostID: postRef.Root,
	})

	return nil
}
