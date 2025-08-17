package post_liker

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/like/like_writer"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type PostLiker struct {
	likeWriter *like_writer.LikeWriter
	bus        *pubsub.Bus
}

func New(
	likeWriter *like_writer.LikeWriter,
	bus *pubsub.Bus,
) *PostLiker {
	return &PostLiker{
		likeWriter: likeWriter,
		bus:        bus,
	}
}

func (l *PostLiker) AddPostLike(ctx context.Context, accountID account.AccountID, postID post.ID) error {
	err := l.likeWriter.AddPostLike(ctx, accountID, postID)
	if err != nil {
		return err
	}

	l.bus.Publish(ctx, &message.EventPostLiked{
		PostID: postID,
	})

	return nil
}

func (l *PostLiker) RemovePostLike(ctx context.Context, accountID account.AccountID, postID post.ID) error {
	err := l.likeWriter.RemovePostLike(ctx, accountID, postID)
	if err != nil {
		return err
	}

	l.bus.Publish(ctx, &message.EventPostUnliked{
		PostID: postID,
	})

	return nil
}
