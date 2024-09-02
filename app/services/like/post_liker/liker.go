package post_liker

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/like/like_writer"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type PostLiker struct {
	likeWriter  *like_writer.LikeWriter
	likeQueue   pubsub.Topic[mq.LikePost]
	unlikeQueue pubsub.Topic[mq.UnlikePost]
}

func New(
	likeWriter *like_writer.LikeWriter,
	likeQueue pubsub.Topic[mq.LikePost],
	unlikeQueue pubsub.Topic[mq.UnlikePost],
) *PostLiker {
	return &PostLiker{
		likeWriter:  likeWriter,
		likeQueue:   likeQueue,
		unlikeQueue: unlikeQueue,
	}
}

func (l *PostLiker) AddPostLike(ctx context.Context, accountID account.AccountID, postID post.ID) error {
	err := l.likeWriter.AddPostLike(ctx, accountID, postID)
	if err != nil {
		return err
	}

	if err = l.likeQueue.Publish(ctx, mq.LikePost{PostID: postID}); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (l *PostLiker) RemovePostLike(ctx context.Context, accountID account.AccountID, postID post.ID) error {
	err := l.likeWriter.RemovePostLike(ctx, accountID, postID)
	if err != nil {
		return err
	}

	if err = l.unlikeQueue.Publish(ctx, mq.UnlikePost{PostID: postID}); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
