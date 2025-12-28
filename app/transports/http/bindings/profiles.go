package bindings

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/cachecontrol"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/profile/follow_querier"
	"github.com/Southclaws/storyden/app/resources/profile/profile_cache"
	"github.com/Southclaws/storyden/app/resources/profile/profile_querier"
	"github.com/Southclaws/storyden/app/resources/profile/profile_search"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/profile/following"
	"github.com/Southclaws/storyden/app/services/reqinfo"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
)

type Profiles struct {
	apiAddress    url.URL
	profileQuery  *profile_querier.Querier
	profile_cache *profile_cache.Cache
	ps            profile_search.Repository
	followQuerier *follow_querier.Querier
	followManager *following.FollowManager
}

func NewProfiles(
	cfg config.Config,
	profileQuery *profile_querier.Querier,
	profile_cache *profile_cache.Cache,
	ps profile_search.Repository,
	followQuerier *follow_querier.Querier,
	followManager *following.FollowManager,
) Profiles {
	return Profiles{
		apiAddress:    cfg.PublicWebAddress,
		profileQuery:  profileQuery,
		profile_cache: profile_cache,
		ps:            ps,
		followQuerier: followQuerier,
		followManager: followManager,
	}
}

func (p *Profiles) ProfileList(ctx context.Context, request openapi.ProfileListRequestObject) (openapi.ProfileListResponseObject, error) {
	pageSize := 50

	page := opt.NewPtrMap(request.Params.Page, func(s string) int {
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return 0
		}

		return max(1, int(v))
	}).Or(1)

	opts := []profile_search.Filter{}

	if request.Params.Q != nil {
		opts = append(opts,
			profile_search.WithNamesLike(*request.Params.Q),
		)
	}

	// API is 1-indexed, internally it's 0-indexed.
	page = max(0, page-1)

	result, err := p.ps.Search(ctx, page, pageSize, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// API is 1-indexed, internally it's 0-indexed.
	page = result.CurrentPage + 1

	return openapi.ProfileList200JSONResponse{
		ProfileListOKJSONResponse: openapi.ProfileListOKJSONResponse{
			PageSize:    pageSize,
			Results:     result.Results,
			TotalPages:  result.TotalPages,
			CurrentPage: page,
			NextPage:    result.NextPage.Ptr(),
			Profiles:    dt.Map(result.Profiles, serialiseProfile),
		},
	}, nil
}

func (p *Profiles) ProfileGet(ctx context.Context, request openapi.ProfileGetRequestObject) (openapi.ProfileGetResponseObject, error) {
	id, err := openapi.ResolveHandle(ctx, p.profileQuery, request.AccountHandle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	etag, notModified := p.profile_cache.Check(ctx, reqinfo.GetCacheQuery(ctx), xid.ID(id))
	if notModified {
		return openapi.ProfileGet304Response{
			Headers: openapi.NotModifiedResponseHeaders{
				CacheControl: getAuthStateCacheControl(ctx, "no-cache"),
				LastModified: etag.Time.Format(time.RFC1123),
				ETag:         etag.String(),
			},
		}, nil
	}

	pro, err := p.profileQuery.GetByID(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if etag == nil {
		p.profile_cache.Store(ctx, xid.ID(id), pro.Updated)
		etag = cachecontrol.NewETag(pro.Updated)
	}

	return openapi.ProfileGet200JSONResponse{
		ProfileGetOKJSONResponse: openapi.ProfileGetOKJSONResponse{
			Body: serialiseProfile(pro),
			Headers: openapi.ProfileGetOKResponseHeaders{
				CacheControl: getAuthStateCacheControl(ctx, "no-cache"),
				LastModified: etag.Time.Format(time.RFC1123),
				ETag:         etag.String(),
			},
		},
	}, nil
}

func (p *Profiles) ProfileFollowersGet(ctx context.Context, request openapi.ProfileFollowersGetRequestObject) (openapi.ProfileFollowersGetResponseObject, error) {
	targetID, err := openapi.ResolveHandle(ctx, p.profileQuery, request.AccountHandle)
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

	result, err := p.followQuerier.GetFollowers(ctx, targetID, page, pageSize)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// API is 1-indexed, internally it's 0-indexed.
	page = result.CurrentPage + 1

	return openapi.ProfileFollowersGet200JSONResponse{
		ProfileFollowersGetOKJSONResponse: openapi.ProfileFollowersGetOKJSONResponse{
			CurrentPage: page,
			Followers:   dt.Map(result.Profiles, serialiseProfileReferencePtr),
			NextPage:    result.NextPage.Ptr(),
			PageSize:    pageSize,
			Results:     result.Results,
			TotalPages:  result.TotalPages,
		},
	}, nil
}

func (p *Profiles) ProfileFollowingGet(ctx context.Context, request openapi.ProfileFollowingGetRequestObject) (openapi.ProfileFollowingGetResponseObject, error) {
	targetID, err := openapi.ResolveHandle(ctx, p.profileQuery, request.AccountHandle)
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

	result, err := p.followQuerier.GetFollowing(ctx, targetID, page, pageSize)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// API is 1-indexed, internally it's 0-indexed.
	page = result.CurrentPage + 1

	return openapi.ProfileFollowingGet200JSONResponse{
		ProfileFollowingGetOKJSONResponse: openapi.ProfileFollowingGetOKJSONResponse{
			CurrentPage: page,
			Following:   dt.Map(result.Profiles, serialiseProfileReferencePtr),
			NextPage:    result.NextPage.Ptr(),
			PageSize:    pageSize,
			Results:     result.Results,
			TotalPages:  result.TotalPages,
		},
	}, nil
}

func (p *Profiles) ProfileFollowersAdd(ctx context.Context, request openapi.ProfileFollowersAddRequestObject) (openapi.ProfileFollowersAddResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	targetID, err := openapi.ResolveHandle(ctx, p.profileQuery, request.AccountHandle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = p.followManager.Follow(ctx, accountID, targetID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ProfileFollowersAdd200Response{}, nil
}

func (p *Profiles) ProfileFollowersRemove(ctx context.Context, request openapi.ProfileFollowersRemoveRequestObject) (openapi.ProfileFollowersRemoveResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	targetID, err := openapi.ResolveHandle(ctx, p.profileQuery, request.AccountHandle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = p.followManager.Unfollow(ctx, accountID, targetID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.ProfileFollowersRemove200Response{}, nil
}

func serialiseProfile(in *profile.Public) openapi.PublicProfile {
	invitedBy := opt.Map(in.InvitedBy, func(ib profile.Ref) openapi.ProfileReference {
		return serialiseProfileReference(ib)
	})

	return openapi.PublicProfile{
		Id:        openapi.Identifier(in.ID.String()),
		CreatedAt: in.Created.Format(time.RFC3339),
		Joined:    in.Created,
		Suspended: in.Deleted.Ptr(),
		DeletedAt: in.Deleted.Ptr(),
		Bio:       in.Bio.HTML(),
		Handle:    in.Handle,
		Name:      in.Name,
		Roles:     serialiseHeldRoleList(in.Roles),
		Followers: in.Followers,
		Following: in.Following,
		LikeScore: in.LikeScore,
		Links:     serialiseExternalLinks(in.ExternalLinks),
		InvitedBy: invitedBy.Ptr(),
		Meta:      in.Metadata,
	}
}

func serialiseProfileReference(a profile.Ref) openapi.ProfileReference {
	return openapi.ProfileReference{
		Id:        *openapi.IdentifierFrom(xid.ID(a.ID)),
		Joined:    a.Created,
		Suspended: a.Deleted.Ptr(),
		Handle:    (openapi.AccountHandle)(a.Handle),
		Name:      a.Name,
	}
}

func serialiseProfileReferenceFromAccount(a account.Account) openapi.ProfileReference {
	return openapi.ProfileReference{
		Id:        *openapi.IdentifierFrom(xid.ID(a.ID)),
		Joined:    a.CreatedAt,
		Suspended: a.DeletedAt.Ptr(),
		Handle:    (openapi.AccountHandle)(a.Handle),
		Name:      a.Name,
	}
}

func serialiseProfileReferencePtr(a *profile.Ref) openapi.ProfileReference {
	return serialiseProfileReference(*a)
}
