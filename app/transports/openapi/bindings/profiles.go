package bindings

import (
	"context"
	"fmt"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	account_repo "github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/services/account"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/openapi"
	"github.com/Southclaws/storyden/internal/utils"
)

type Profiles struct {
	apiAddress string
	as         account.Service
	ar         account_repo.Repository
}

func NewProfiles(cfg config.Config, as account.Service, ar account_repo.Repository) Profiles {
	return Profiles{cfg.PublicWebAddress, as, ar}
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
