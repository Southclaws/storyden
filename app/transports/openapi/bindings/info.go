package bindings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/internal/openapi"
)

type Info struct {
	sr settings.Repository
}

func NewInfo(sr settings.Repository) Info {
	return Info{
		sr: sr,
	}
}

func (i Info) GetInfo(ctx context.Context, request openapi.GetInfoRequestObject) (openapi.GetInfoResponseObject, error) {
	settings, err := i.sr.Get(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.GetInfo200JSONResponse{
		GetInfoOKJSONResponse: openapi.GetInfoOKJSONResponse{
			Title:       settings.Title.Get(),
			Description: settings.Description.Get(),
		},
	}, nil
}
