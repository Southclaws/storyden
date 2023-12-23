package bindings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/account_suspension"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/openapi"
)

var errNotAuthorised = fault.Wrap(fault.New("not authorised"), ftag.With(ftag.PermissionDenied))

type Admin struct {
	ar account.Repository
	as account_suspension.Service
	sr settings.Repository
}

func NewAdmin(ar account.Repository, as account_suspension.Service, sr settings.Repository) Admin {
	return Admin{ar, as, sr}
}

func (a *Admin) AdminSettingsUpdate(ctx context.Context, request openapi.AdminSettingsUpdateRequestObject) (openapi.AdminSettingsUpdateResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := a.ar.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !acc.Admin {
		return nil, fault.Wrap(errNotAuthorised, fctx.With(ctx))
	}

	settings, err := a.sr.Set(ctx, settings.Partial{
		Title:        opt.NewPtr(request.Body.Title),
		Description:  opt.NewPtr(request.Body.Description),
		AccentColour: opt.NewPtr(request.Body.AccentColour),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AdminSettingsUpdate200JSONResponse{
		AdminSettingsUpdateOKJSONResponse: openapi.AdminSettingsUpdateOKJSONResponse(serialiseSettings(settings)),
	}, nil
}

func (i *Admin) AdminAccountBanCreate(ctx context.Context, request openapi.AdminAccountBanCreateRequestObject) (openapi.AdminAccountBanCreateResponseObject, error) {
	id, err := openapi.ResolveHandle(ctx, i.ar, request.AccountHandle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := i.ar.GetByID(ctx, accountID)
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
		AccountGetOKJSONResponse: openapi.AccountGetOKJSONResponse(serialiseAccount(acc)),
	}, nil
}

func (i *Admin) AdminAccountBanRemove(ctx context.Context, request openapi.AdminAccountBanRemoveRequestObject) (openapi.AdminAccountBanRemoveResponseObject, error) {
	id, err := openapi.ResolveHandle(ctx, i.ar, request.AccountHandle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := i.ar.GetByID(ctx, accountID)
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
		AccountGetOKJSONResponse: openapi.AccountGetOKJSONResponse(serialiseAccount(acc)),
	}, nil
}

func serialiseSettings(in *settings.Settings) openapi.AdminSettingsProps {
	return openapi.AdminSettingsProps{
		AccentColour: in.AccentColour.Get(),
		Description:  in.Description.Get(),
		Title:        in.Title.Get(),
	}
}
