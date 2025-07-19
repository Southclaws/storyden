package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/invitation"
	"github.com/Southclaws/storyden/app/resources/account/invitation/invitation_querier"
	"github.com/Southclaws/storyden/app/resources/account/invitation/invitation_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Invitations struct {
	accountQuerier *account_querier.Querier
	invQuerier     *invitation_querier.Querier
	invWriter      *invitation_writer.Writer
}

func NewInvitations(accountQuerier *account_querier.Querier, invQuerier *invitation_querier.Querier, invWriter *invitation_writer.Writer) Invitations {
	return Invitations{
		accountQuerier: accountQuerier,
		invQuerier:     invQuerier,
		invWriter:      invWriter,
	}
}

func (h *Invitations) InvitationList(ctx context.Context, request openapi.InvitationListRequestObject) (openapi.InvitationListResponseObject, error) {
	session, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := h.accountQuerier.GetByID(ctx, session)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	opts := []invitation_querier.Filter{}

	filterByAccountID := opt.Map(opt.NewPtr(request.Params.AccountId), openapi.GetAccountID)

	if id, ok := filterByAccountID.Get(); ok && acc.Roles.Permissions().HasAll(rbac.PermissionAdministrator) {
		// Only administrators may use the account ID filter parameter.
		opts = append(opts, invitation_querier.WithCreator(account.AccountID(id)))
	} else {
		// Otherwise, always force filter by session account ID
		opts = append(opts, invitation_querier.WithCreator(session))

		// TODO: Provide a way for administrators to list *all* invitations.
		// Perhaps via `account_id=*`.
	}

	invs, err := h.invQuerier.List(ctx, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.InvitationList200JSONResponse{
		InvitationListOKJSONResponse: openapi.InvitationListOKJSONResponse{
			// TODO: Pagination.
			CurrentPage: 1,
			Invitations: dt.Map(invs, serialiseInvitationPtr),
			PageSize:    len(invs),
			Results:     len(invs),
			TotalPages:  1,
		},
	}, nil
}

func (h *Invitations) InvitationCreate(ctx context.Context, request openapi.InvitationCreateRequestObject) (openapi.InvitationCreateResponseObject, error) {
	session, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := h.accountQuerier.GetByID(ctx, session)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	inv, err := h.invWriter.Create(ctx, account.AccountID(acc.ID), opt.NewPtr(request.Body.Message))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.InvitationCreate200JSONResponse{
		InvitationCreateOKJSONResponse: openapi.InvitationCreateOKJSONResponse(serialiseInvitationPtr(inv)),
	}, nil
}

func (h *Invitations) InvitationDelete(ctx context.Context, request openapi.InvitationDeleteRequestObject) (openapi.InvitationDeleteResponseObject, error) {
	session, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := h.accountQuerier.GetByID(ctx, session)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	invid, err := xid.FromString(request.InvitationId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	inv, err := h.invQuerier.GetByID(ctx, invid)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = acc.Roles.Permissions().Authorise(ctx, func() error {
		if inv.Creator.ID == acc.ID {
			return nil
		}
		return fault.New("not the creator of this invitation")
	}, rbac.PermissionAdministrator)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = h.invWriter.Delete(ctx, invid)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.InvitationDelete200Response{}, nil
}

func (h *Invitations) InvitationGet(ctx context.Context, request openapi.InvitationGetRequestObject) (openapi.InvitationGetResponseObject, error) {
	invid, err := xid.FromString(request.InvitationId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	inv, err := h.invQuerier.GetByID(ctx, invid)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.InvitationGet200JSONResponse{
		InvitationGetOKJSONResponse: openapi.InvitationGetOKJSONResponse(serialiseInvitationPtr(inv)),
	}, nil
}

func serialiseInvitationPtr(inv *invitation.Invitation) openapi.Invitation {
	return openapi.Invitation{
		Id:        inv.ID.String(),
		CreatedAt: inv.CreatedAt,
		UpdatedAt: inv.UpdatedAt,
		DeletedAt: inv.DeletedAt.Ptr(),
		Creator:   serialiseProfileReferenceFromAccount(inv.Creator),
		Message:   inv.Message.Ptr(),
	}
}

func deserialiseInvitationID(id *string) (opt.Optional[xid.ID], error) {
	inv, err := opt.MapErr(opt.NewPtr(id), xid.FromString)
	if err != nil {
		return nil, err
	}

	return inv, nil
}
