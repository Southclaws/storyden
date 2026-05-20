package bindings

import (
	"context"
	"net/mail"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func (h *Accounts) AccountEmailAdd(ctx context.Context, request openapi.AccountEmailAddRequestObject) (openapi.AccountEmailAddResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	emailAddress, err := mail.ParseAddress(strings.ToLower(request.Body.EmailAddress))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	ae, err := h.accountEmail.Add(ctx, accountID, *emailAddress)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountEmailAdd200JSONResponse{
		AccountEmailUpdateOKJSONResponse: openapi.AccountEmailUpdateOKJSONResponse(serialiseEmailAddressPtr(ae)),
	}, nil
}

func (h *Accounts) AccountEmailRemove(ctx context.Context, request openapi.AccountEmailRemoveRequestObject) (openapi.AccountEmailRemoveResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	id, err := xid.FromString(request.EmailAddressId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	err = h.accountEmail.Remove(ctx, accountID, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountEmailRemove200Response{}, nil
}

func (h *Accounts) AccountManageUpdateEmailVerifiedStatus(ctx context.Context, request openapi.AccountManageUpdateEmailVerifiedStatusRequestObject) (openapi.AccountManageUpdateEmailVerifiedStatusResponseObject, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionManageAccounts); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.PermissionDenied))
	}

	accountID, err := xid.FromString(request.AccountId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	emailAddressID, err := xid.FromString(request.EmailAddressId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	ae, err := h.accountEmail.SetVerifiedStatus(ctx, account.AccountID(accountID), emailAddressID, request.Body.Verified)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountManageUpdateEmailVerifiedStatus200JSONResponse{
		AccountEmailUpdateOKJSONResponse: openapi.AccountEmailUpdateOKJSONResponse(serialiseEmailAddressPtr(ae)),
	}, nil
}
