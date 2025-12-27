package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/account/authentication/access_key"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/profile/profile_querier"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/account/account_suspension"
	"github.com/Southclaws/storyden/app/services/admin/settings_manager"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

var errNotAuthorised = fault.Wrap(fault.New("not authorised"), ftag.With(ftag.PermissionDenied))

type Admin struct {
	accountQuery    *account_querier.Querier
	profileQuery    *profile_querier.Querier
	as              account_suspension.Service
	settingsManager *settings_manager.Manager
	akr             *access_key.Repository
}

func NewAdmin(
	accountQuery *account_querier.Querier,
	profileQuery *profile_querier.Querier,
	as account_suspension.Service,
	settingsManager *settings_manager.Manager,
	akr *access_key.Repository,
) Admin {
	return Admin{
		accountQuery:    accountQuery,
		profileQuery:    profileQuery,
		as:              as,
		settingsManager: settingsManager,
		akr:             akr,
	}
}

func (a *Admin) AdminSettingsGet(ctx context.Context, request openapi.AdminSettingsGetRequestObject) (openapi.AdminSettingsGetResponseObject, error) {
	settings, err := a.settingsManager.Get(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AdminSettingsGet200JSONResponse{
		AdminSettingsGetOKJSONResponse: openapi.AdminSettingsGetOKJSONResponse(serialiseSettings(settings)),
	}, nil
}

func (a *Admin) AdminSettingsUpdate(ctx context.Context, request openapi.AdminSettingsUpdateRequestObject) (openapi.AdminSettingsUpdateResponseObject, error) {
	content, err := opt.MapErr(opt.NewPtr(request.Body.Content), datagraph.NewRichText)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	authMode, err := opt.MapErr(opt.NewPtr(request.Body.AuthenticationMode), deserialiseAuthMode)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var services opt.Optional[settings.ServiceSettings]
	if request.Body.Services != nil && request.Body.Services.Moderation != nil {
		moderation := request.Body.Services.Moderation
		services = opt.New(settings.ServiceSettings{
			Moderation: opt.New(settings.ModerationServiceSettings{
				MaxThreadBodyLength: opt.NewPtr(moderation.ThreadBodySizeLimit),
				MaxReplyBodyLength:  opt.NewPtr(moderation.ReplyBodySizeLimit),
				WordBlocklist:       opt.NewPtr(moderation.WordBlockList),
			}),
		})
	}

	settings, err := a.settingsManager.Set(ctx, settings.Settings{
		Title:              opt.NewPtr(request.Body.Title),
		Description:        opt.NewPtr(request.Body.Description),
		Content:            content,
		AccentColour:       opt.NewPtr(request.Body.AccentColour),
		AuthenticationMode: authMode,
		Services:           services,
		Metadata:           opt.NewPtr((*map[string]any)(request.Body.Metadata)),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AdminSettingsUpdate200JSONResponse{
		AdminSettingsUpdateOKJSONResponse: openapi.AdminSettingsUpdateOKJSONResponse(serialiseSettings(settings)),
	}, nil
}

func (i *Admin) AdminAccountBanCreate(ctx context.Context, request openapi.AdminAccountBanCreateRequestObject) (openapi.AdminAccountBanCreateResponseObject, error) {
	id, err := openapi.ResolveHandle(ctx, i.profileQuery, request.AccountHandle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := i.accountQuery.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !acc.Admin {
		return nil, fault.Wrap(errNotAuthorised, fctx.With(ctx), ftag.With(ftag.PermissionDenied))
	}

	acc, err = i.as.Suspend(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AdminAccountBanCreate200JSONResponse{
		AccountGetOKJSONResponse: openapi.AccountGetOKJSONResponse{
			Body: serialiseAccount(acc),
		},
	}, nil
}

func (i *Admin) AdminAccountBanRemove(ctx context.Context, request openapi.AdminAccountBanRemoveRequestObject) (openapi.AdminAccountBanRemoveResponseObject, error) {
	id, err := openapi.ResolveHandle(ctx, i.profileQuery, request.AccountHandle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := i.accountQuery.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !acc.Admin {
		return nil, fault.Wrap(errNotAuthorised, fctx.With(ctx), ftag.With(ftag.PermissionDenied))
	}

	acc, err = i.as.Reinstate(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AdminAccountBanRemove200JSONResponse{
		AccountGetOKJSONResponse: openapi.AccountGetOKJSONResponse{
			Body: serialiseAccount(acc),
		},
	}, nil
}

func (i *Admin) AdminAccessKeyList(ctx context.Context, request openapi.AdminAccessKeyListRequestObject) (openapi.AdminAccessKeyListResponseObject, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	list, err := i.akr.ListAllAsAdmin(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AdminAccessKeyList200JSONResponse{
		AdminAccessKeyListOKJSONResponse: openapi.AdminAccessKeyListOKJSONResponse{
			Keys: serialiseOwnedAccessKeyList(list),
		},
	}, nil
}

func (i *Admin) AdminAccessKeyDelete(ctx context.Context, request openapi.AdminAccessKeyDeleteRequestObject) (openapi.AdminAccessKeyDeleteResponseObject, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionAdministrator); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	_, err := i.akr.RevokeAsAdmin(ctx, deserialiseID(request.AccessKeyId))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NoContentResponse{}, nil
}

func serialiseSettings(in *settings.Settings) openapi.AdminSettingsProps {
	return openapi.AdminSettingsProps{
		AccentColour:       in.AccentColour.OrZero(),
		Description:        in.Description.OrZero(),
		Content:            in.Content.OrZero().HTML(),
		Title:              in.Title.OrZero(),
		AuthenticationMode: openapi.AuthMode(in.AuthenticationMode.Or(authentication.ModeHandle).String()),
		Services:           opt.Map(in.Services, serialiseServiceSettings).Ptr(),
		Metadata:           (*openapi.Metadata)(in.Metadata.Ptr()),
	}
}

func serialiseServiceSettings(in settings.ServiceSettings) openapi.AdminSettingsServiceProps {
	return openapi.AdminSettingsServiceProps{
		Moderation: opt.Map(in.Moderation, serialiseModerationSettings).Ptr(),
	}
}

func serialiseModerationSettings(in settings.ModerationServiceSettings) openapi.ModerationServiceSettings {
	return openapi.ModerationServiceSettings{
		ThreadBodySizeLimit: in.MaxThreadBodyLength.Ptr(),
		ReplyBodySizeLimit:  in.MaxReplyBodyLength.Ptr(),
		WordBlockList:       in.WordBlocklist.Ptr(),
	}
}

func serialiseOwnedAccessKey(in *authentication.Authentication) openapi.OwnedAccessKey {
	return openapi.OwnedAccessKey{
		Id:        in.ID.String(),
		CreatedAt: in.Account.CreatedAt,
		ExpiresAt: in.Expires.Ptr(),
		Enabled:   !in.Disabled,
		Name:      in.Name.Or("Unnamed"),
		CreatedBy: serialiseProfileReferenceFromAccount(in.Account),
	}
}

func serialiseOwnedAccessKeyList(in []*authentication.Authentication) openapi.OwnedAccessKeyList {
	return dt.Map(in, serialiseOwnedAccessKey)
}
