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

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/profile/follow_querier"
	"github.com/Southclaws/storyden/app/resources/profile/profile_search"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/profile/following"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
)

type Profiles struct {
	apiAddress    url.URL
	accountQuery  *account_querier.Querier
	ps            profile_search.Repository
	followQuerier *follow_querier.Querier
	followManager *following.FollowManager
}

func NewProfiles(
	cfg config.Config,
	accountQuery *account_querier.Querier,
	ps profile_search.Repository,
	followQuerier *follow_querier.Querier,
	followManager *following.FollowManager,
) Profiles {
	return Profiles{
		apiAddress:    cfg.PublicWebAddress,
		accountQuery:  accountQuery,
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
	id, err := openapi.ResolveHandle(ctx, p.accountQuery, request.AccountHandle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := p.accountQuery.GetByID(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	pro := profile.ProfileFromAccount(acc)

	return openapi.ProfileGet200JSONResponse{
		ProfileGetOKJSONResponse: openapi.ProfileGetOKJSONResponse(serialiseProfile(pro)),
	}, nil
}

func (p *Profiles) ProfileFollowersGet(ctx context.Context, request openapi.ProfileFollowersGetRequestObject) (openapi.ProfileFollowersGetResponseObject, error) {
	targetID, err := openapi.ResolveHandle(ctx, p.accountQuery, request.AccountHandle)
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
	targetID, err := openapi.ResolveHandle(ctx, p.accountQuery, request.AccountHandle)
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

	targetID, err := openapi.ResolveHandle(ctx, p.accountQuery, request.AccountHandle)
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

	targetID, err := openapi.ResolveHandle(ctx, p.accountQuery, request.AccountHandle)
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
	invitedBy := opt.Map(in.InvitedBy, func(ib profile.Public) openapi.ProfileReference {
		return serialiseProfileReference(ib)
	})

	return openapi.PublicProfile{
		Id:        openapi.Identifier(in.ID.String()),
		CreatedAt: in.Created.Format(time.RFC3339),
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
