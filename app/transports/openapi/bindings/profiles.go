package bindings

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	account_repo "github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/profile_search"
	"github.com/Southclaws/storyden/app/services/account"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/openapi"
	"github.com/Southclaws/storyden/internal/utils"
)

type Profiles struct {
	apiAddress string
	as         account.Service
	ar         account_repo.Repository
	ps         profile_search.Repository
}

func NewProfiles(cfg config.Config, as account.Service, ar account_repo.Repository, ps profile_search.Repository) Profiles {
	return Profiles{cfg.PublicWebAddress, as, ar, ps}
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
			profile_search.WithDisplayNameContains(*request.Params.Q),
			profile_search.WithHandleContains(*request.Params.Q),
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
	id, err := openapi.ResolveHandle(ctx, p.ar, request.AccountHandle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := p.as.Get(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// TODO: Make this a bit more well designed and less coupled to the hostname
	avatarURL := fmt.Sprintf("%s/api/v1/accounts/%s/avatar", p.apiAddress, acc.Handle)

	return openapi.ProfileGet200JSONResponse{
		ProfileGetOKJSONResponse: openapi.ProfileGetOKJSONResponse{
			Id:        openapi.Identifier(acc.ID.String()),
			Bio:       utils.Ref(acc.Bio.OrZero()),
			Handle:    acc.Handle,
			Image:     &avatarURL,
			Name:      acc.Name,
			CreatedAt: acc.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func serialiseProfile(in *profile.Profile) openapi.PublicProfile {
	return openapi.PublicProfile{
		Id:     openapi.Identifier(in.ID.String()),
		Bio:    &in.Bio,
		Handle: in.Handle,
		// Image:     &avatarURL,
		Name:      in.Name,
		CreatedAt: in.Created.Format(time.RFC3339),
	}
}
