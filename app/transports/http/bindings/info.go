package bindings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/icon"
	"github.com/Southclaws/storyden/app/services/onboarding"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Info struct {
	sr settings.Repository
	os onboarding.Service
	is icon.Service
}

func NewInfo(sr settings.Repository, os onboarding.Service, is icon.Service) Info {
	return Info{
		sr: sr,
		os: os,
		is: is,
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
		GetInfoOKJSONResponse: openapi.GetInfoOKJSONResponse(serialiseInfo(settings, *status)),
	}, nil
}

func (i Info) IconGet(ctx context.Context, request openapi.IconGetRequestObject) (openapi.IconGetResponseObject, error) {
	a, r, err := i.is.Get(ctx, string(request.IconSize))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.IconGet200AsteriskResponse{
		AssetGetOKAsteriskResponse: openapi.AssetGetOKAsteriskResponse{
			Body:          r,
			ContentType:   a.Metadata.GetMIMEType(),
			ContentLength: int64(a.Size),
		},
	}, nil
}

func (i Info) IconUpload(ctx context.Context, request openapi.IconUploadRequestObject) (openapi.IconUploadResponseObject, error) {
	err := i.is.Upload(ctx, request.Body)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.IconUpload200Response{}, nil
}

func serialiseInfo(s *settings.Settings, status onboarding.Status) openapi.Info {
	return openapi.Info{
		Title:            s.Title.Get(),
		Description:      s.Description.Get(),
		Content:          s.Content.Get().HTML(),
		AccentColour:     s.AccentColour.Get(),
		OnboardingStatus: openapi.OnboardingStatus(status.String()),
	}
}
