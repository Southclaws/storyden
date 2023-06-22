package bindings

import (
	"context"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	account_repo "github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/services/account"
	"github.com/Southclaws/storyden/internal/openapi"
	"github.com/Southclaws/storyden/internal/utils"
)

type Profiles struct {
	as account.Service
	ar account_repo.Repository
}

func NewProfiles(as account.Service, ar account_repo.Repository) Profiles {
	return Profiles{as, ar}
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

	interests := dt.Map(acc.Interests, serialiseTag)

	return openapi.ProfileGet200JSONResponse{
		ProfileGetOKJSONResponse: openapi.ProfileGetOKJSONResponse{
			Id:        openapi.Identifier(acc.ID.String()),
			Bio:       utils.Ref(acc.Bio.OrZero()),
			Handle:    acc.Handle,
			Name:      acc.Name,
			Interests: interests,
			CreatedAt: acc.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}
