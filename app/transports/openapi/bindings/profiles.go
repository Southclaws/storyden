package bindings

import (
	"context"
	"time"

	"github.com/Southclaws/dt"

	account_resource "github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/app/services/account"
	"github.com/Southclaws/storyden/app/transports/openapi/openapi"
	"github.com/Southclaws/storyden/internal/errctx"
	"github.com/Southclaws/storyden/internal/utils"
)

type Profiles struct {
	as account.Service
}

func NewProfiles(as account.Service) Profiles {
	return Profiles{as}
}

func (p *Profiles) ProfilesGet(ctx context.Context, request openapi.ProfilesGetRequestObject) (openapi.ProfilesGetResponseObject, error) {
	acc, err := p.as.Get(ctx, account_resource.AccountID(request.AccountId.XID()))
	if err != nil {
		return nil, errctx.Wrap(err, ctx)
	}

	interests := dt.Map(acc.Interests, func(t tag.Tag) string {
		return t.Name
	})

	return openapi.ProfilesGet200JSONResponse{
		Id:        openapi.Identifier(acc.ID.String()),
		Bio:       utils.Ref(acc.Bio.ElseZero()),
		Handle:    &acc.Handle,
		Name:      &acc.Name,
		Interests: &interests,
		CreatedAt: acc.CreatedAt.Format(time.RFC3339),
	}, nil
}
