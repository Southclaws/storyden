package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/resources/like"
	"github.com/Southclaws/storyden/app/resources/like/item_like"
	"github.com/Southclaws/storyden/app/resources/like/like_querier"
	"github.com/Southclaws/storyden/app/resources/like/profile_like"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/profile/profile_querier"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/like/post_liker"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Likes struct {
	likeQuerier  *like_querier.LikeQuerier
	postLiker    *post_liker.PostLiker
	profileQuery *profile_querier.Querier
}

func NewLikes(
	likeQuerier *like_querier.LikeQuerier,
	postLiker *post_liker.PostLiker,
	profileQuery *profile_querier.Querier,
) Likes {
	return Likes{
		likeQuerier:  likeQuerier,
		postLiker:    postLiker,
		profileQuery: profileQuery,
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
	accountID, err := openapi.ResolveHandle(ctx, h.profileQuery, request.AccountHandle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	params := deserialisePageParams(request.Params.Page, 50)

	result, err := h.likeQuerier.GetProfileLikes(ctx, accountID, params)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.LikeProfileGet200JSONResponse{
		LikeProfileGetOKJSONResponse: openapi.LikeProfileGetOKJSONResponse{
			PageSize:    result.Size,
			Results:     result.Results,
			TotalPages:  result.TotalPages,
			CurrentPage: result.CurrentPage,
			NextPage:    result.NextPage.Ptr(),
			Likes:       dt.Map(result.Items, serialiseProfileLike),
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
