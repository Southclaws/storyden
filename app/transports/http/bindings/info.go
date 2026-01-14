package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/branding/banner"
	"github.com/Southclaws/storyden/app/services/branding/icon"
	"github.com/Southclaws/storyden/app/services/system/instance_info"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Info struct {
	systemInfo *instance_info.Provider
	is         icon.Service
	os         banner.Service
}

func NewInfo(systemInfo *instance_info.Provider, is icon.Service, os banner.Service) Info {
	return Info{
		systemInfo: systemInfo,
		is:         is,
		os:         os,
	}
}

func (i Info) GetInfo(ctx context.Context, request openapi.GetInfoRequestObject) (openapi.GetInfoResponseObject, error) {
	info, err := i.systemInfo.Get(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.GetInfo200JSONResponse{
		GetInfoOKJSONResponse: openapi.GetInfoOKJSONResponse(serialiseInfo(info)),
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
			ContentType:   a.MIME.String(),
			ContentLength: int64(a.Size),
			Headers: openapi.AssetGetOKResponseHeaders{
				CacheControl: "public, max-age=31536000",
			},
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

func (i Info) BannerGet(ctx context.Context, request openapi.BannerGetRequestObject) (openapi.BannerGetResponseObject, error) {
	a, r, err := i.os.Get(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.BannerGet200AsteriskResponse{
		AssetGetOKAsteriskResponse: openapi.AssetGetOKAsteriskResponse{
			Body:          r,
			ContentType:   a.MIME.String(),
			ContentLength: int64(a.Size),
			Headers: openapi.AssetGetOKResponseHeaders{
				CacheControl: "public, max-age=3600",
			},
		},
	}, nil
}

func (i Info) BannerUpload(ctx context.Context, request openapi.BannerUploadRequestObject) (openapi.BannerUploadResponseObject, error) {
	err := i.os.Upload(ctx, request.Body)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.BannerUpload200Response{}, nil
}

func serialiseInfo(info *instance_info.Info) openapi.Info {
	quickReactions := info.Settings.QuickReactions.Or(settings.DefaultQuickReactions)

	return openapi.Info{
		Title:              info.Settings.Title.OrZero(),
		Description:        info.Settings.Description.OrZero(),
		Content:            info.Settings.Content.OrZero().HTML(),
		AccentColour:       info.Settings.AccentColour.OrZero(),
		OnboardingStatus:   openapi.OnboardingStatus(info.OnboardingStatus.String()),
		AuthenticationMode: openapi.AuthMode(info.Settings.AuthenticationMode.Or(authentication.ModeHandle).String()),
		Capabilities:       serialiseCapabilitiesList(info.Capabilities),
		QuickReactions:     quickReactions,
		Metadata:           (*openapi.Metadata)(info.Settings.Metadata.Ptr()),
	}
}

func serialiseCapabilitiesList(cs instance_info.Capabilities) openapi.InstanceCapabilityList {
	return dt.Map(cs, func(c instance_info.Capability) openapi.InstanceCapability {
		return openapi.InstanceCapability(c.String())
	})
}
