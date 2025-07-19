package bindings

import (
	"context"
	"strconv"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/like"
	"github.com/Southclaws/storyden/app/resources/like/item_like"
	"github.com/Southclaws/storyden/app/resources/like/like_querier"
	"github.com/Southclaws/storyden/app/resources/like/profile_like"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/like/post_liker"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Likes struct {
	likeQuerier *like_querier.LikeQuerier
	postLiker   *post_liker.PostLiker
}

func NewLikes(
	likeQuerier *like_querier.LikeQuerier,
	postLiker *post_liker.PostLiker,
) Likes {
	return Likes{
		likeQuerier: likeQuerier,
		postLiker:   postLiker,
	}
}

func (h *Likes) LikePostGet(ctx context.Context, request openapi.LikePostGetRequestObject) (openapi.LikePostGetResponseObject, error) {
	postID := deserialisePostID(request.PostId)

	likes, err := h.likeQuerier.GetPostLikes(ctx, post.ID(postID))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	mapped := dt.Map(likes, serialiseItemLike)

	return openapi.LikePostGet200JSONResponse{
		LikePostGetOKJSONResponse: openapi.LikePostGetOKJSONResponse{
			Likes: mapped,
		},
	}, nil
}

func (h *Likes) LikePostAdd(ctx context.Context, request openapi.LikePostAddRequestObject) (openapi.LikePostAddResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	postID := deserialisePostID(request.PostId)

	err = h.postLiker.AddPostLike(ctx, accountID, postID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.LikePostAdd200Response{}, nil
}

func (h *Likes) LikePostRemove(ctx context.Context, request openapi.LikePostRemoveRequestObject) (openapi.LikePostRemoveResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	postID := deserialisePostID(request.PostId)

	err = h.postLiker.RemovePostLike(ctx, accountID, postID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.LikePostRemove200Response{}, nil
}

func (h *Likes) LikeProfileGet(ctx context.Context, request openapi.LikeProfileGetRequestObject) (openapi.LikeProfileGetResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	pageSize := 50

	page := opt.NewPtrMap(request.Params.Page, func(s string) int {
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return 0
		}

		return max(1, int(v))
	}).Or(1)

	// API is 1-indexed, internally it's 0-indexed.
	page = max(0, page-1)

	result, err := h.likeQuerier.GetProfileLikes(ctx, accountID, page, pageSize)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// API is 1-indexed, internally it's 0-indexed.
	page = result.CurrentPage + 1

	return openapi.LikeProfileGet200JSONResponse{
		LikeProfileGetOKJSONResponse: openapi.LikeProfileGetOKJSONResponse{
			PageSize:    pageSize,
			Results:     result.Results,
			TotalPages:  result.TotalPages,
			CurrentPage: page,
			NextPage:    result.NextPage.Ptr(),
			Likes:       dt.Map(result.Likes, serialiseProfileLike),
		},
	}, nil
}

func serialiseItemLike(like *item_like.Like) openapi.ItemLike {
	return openapi.ItemLike{
		Id:        like.ID.String(),
		CreatedAt: like.Created,
		Owner:     serialiseProfileReference(like.Owner),
	}
}

func serialiseProfileLike(like *profile_like.Like) openapi.ProfileLike {
	return openapi.ProfileLike{
		Id:        like.ID.String(),
		CreatedAt: like.Created,
		Item:      serialiseDatagraphItem(like.Item),
	}
}

func serialiseLikeStatus(like *like.Status) openapi.LikeData {
	return openapi.LikeData{
		Likes: like.Count,
		Liked: like.Status,
	}
}
