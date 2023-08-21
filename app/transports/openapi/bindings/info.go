package bindings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/onboarding"
	"github.com/Southclaws/storyden/internal/openapi"
)

type Info struct {
	sr settings.Repository
	os onboarding.Service
}

func NewInfo(sr settings.Repository, os onboarding.Service) Info {
	return Info{
		sr: sr,
		os: os,
	}
}

func (i Info) GetInfo(ctx context.Context, request openapi.GetInfoRequestObject) (openapi.GetInfoResponseObject, error) {
	settings, err := i.sr.Get(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	status, err := i.os.GetOnboardingStatus(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.GetInfo200JSONResponse{
		GetInfoOKJSONResponse: openapi.GetInfoOKJSONResponse{
			Title:            settings.Title.Get(),
			Description:      settings.Description.Get(),
			AccentColour:     settings.AccentColour.Get(),
			OnboardingStatus: openapi.OnboardingStatus(status.String()),
		},
	}, nil
}
