package bindings

import (
	"context"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	oauthservice "github.com/Southclaws/storyden/app/services/authentication/oauth"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func (i *Admin) AdminOAuthRegistrationApprovalGet(ctx context.Context, request openapi.AdminOAuthRegistrationApprovalGetRequestObject) (openapi.AdminOAuthRegistrationApprovalGetResponseObject, error) {
	review, oauthErr, err := i.oauth.GetRegistrationApproval(ctx, request.Params.VerificationCode)
	if err != nil {
		return nil, err
	}

	if oauthErr != nil {
		return openapi.AdminOAuthRegistrationApprovalGet400JSONResponse{
			Error:            oauthErr.Code,
			ErrorDescription: &oauthErr.Description,
		}, nil
	}

	return openapi.AdminOAuthRegistrationApprovalGet200JSONResponse(serialiseOAuthRegistrationApproval(review)), nil
}

func serialiseOAuthRegistrationApproval(review *oauthservice.RegistrationApprovalReview) openapi.OAuthRegistrationApprovalReview {
	in := review.Metadata
	return openapi.OAuthRegistrationApprovalReview{
		VerificationCode: review.VerificationCode,
		CreatedAt:        review.CreatedAt,
		ExpiresAt:        review.ExpiresAt,
		Metadata: openapi.OAuthClientRegisterProps{
			ClientName:              &in.ClientName,
			GrantTypes:              &in.GrantTypes,
			ResponseTypes:           &in.ResponseTypes,
			RedirectUris:            &in.RedirectURIs,
			Scope:                   &in.Scope,
			TokenEndpointAuthMethod: &in.TokenEndpointAuthMethod,
			Jwks:                    &in.JWKs,
			ApplicationType:         optionalString(in.ApplicationType),
			LogoUri:                 optionalString(in.LogoURI),
			ClientUri:               optionalString(in.ClientURI),
			TosUri:                  optionalString(in.TOSURI),
			PolicyUri:               optionalString(in.PolicyURI),
		},
	}
}

func (i *Admin) AdminOAuthRegistrationApprovalSubmit(ctx context.Context, request openapi.AdminOAuthRegistrationApprovalSubmitRequestObject) (openapi.AdminOAuthRegistrationApprovalSubmitResponseObject, error) {
	if request.Body == nil {
		return openapi.AdminOAuthRegistrationApprovalSubmit400JSONResponse{Error: "invalid_request"}, nil
	}

	approver, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	oauthErr, err := i.oauth.DecideRegistrationApproval(ctx, approver, request.Body.VerificationCode, request.Body.Approved)
	if err != nil {
		return nil, err
	}

	if oauthErr != nil {
		return openapi.AdminOAuthRegistrationApprovalSubmit400JSONResponse{
			Error:            oauthErr.Code,
			ErrorDescription: &oauthErr.Description,
		}, nil
	}

	return openapi.AdminOAuthRegistrationApprovalSubmit204Response{}, nil
}

func (i *Admin) AdminOAuthRegistrationApprovalList(ctx context.Context, request openapi.AdminOAuthRegistrationApprovalListRequestObject) (openapi.AdminOAuthRegistrationApprovalListResponseObject, error) {
	snapshot := time.Now()
	result, err := i.oauth.ListRegistrationApprovals(ctx, deserialisePageParams(request.Params.Page, 50))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AdminOAuthRegistrationApprovalList200JSONResponse{
		SnapshotAt:    snapshot,
		CurrentPage:   result.CurrentPage,
		NextPage:      result.NextPage.Ptr(),
		PageSize:      result.Size,
		Results:       result.Results,
		TotalPages:    result.TotalPages,
		Registrations: dt.Map(result.Items, serialiseOAuthRegistrationApproval),
	}, nil
}

func (i *Admin) AdminOAuthRegistrationApprovalBulkSubmit(ctx context.Context, request openapi.AdminOAuthRegistrationApprovalBulkSubmitRequestObject) (openapi.AdminOAuthRegistrationApprovalBulkSubmitResponseObject, error) {
	if request.Body == nil {
		return openapi.AdminOAuthRegistrationApprovalBulkSubmit400Response{}, nil
	}

	approver, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	updated, err := i.oauth.DecideAllRegistrationApprovals(ctx, approver, request.Body.Approved, request.Body.CreatedBefore)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AdminOAuthRegistrationApprovalBulkSubmit200JSONResponse{Updated: updated}, nil
}
