package bindings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/oapi-codegen/echo-middleware"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/bindings/openapi_rbac"
)

type Authorisation struct {
	accountQuery *account_querier.Querier
}

func newAuthorisation(aq *account_querier.Querier) *Authorisation {
	return &Authorisation{accountQuery: aq}
}

func (i *Authorisation) validator(ctx context.Context, ai *openapi3filter.AuthenticationInput) error {
	sessionRequired, perm := GetPermissionForOperation(ai.RequestValidationInput.Route.Operation.OperationID)
	if perm == nil {
		// No specific permission required, just need a session.
		return nil
	}

	// security scheme name from openapi.yaml
	if ai.SecuritySchemeName != "browser" {
		if sessionRequired {
			// TODO: Handle more gracefully
			panic("unexpected security scheme for session-required operation: " + ai.SecuritySchemeName)
		}

		return nil
	}

	c := ctx.Value(echomiddleware.EchoContextKey).(echo.Context)

	// first check if the middleware injected an account ID, if not, fail.
	aid, err := session.GetAccountID(c.Request().Context())
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Unauthenticated))
	}

	// Then look up the account.
	// TODO: Cache this.
	a, err := i.accountQuery.GetByID(ctx, aid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	// Reject any requests from suspended accounts.
	if err := a.RejectSuspended(); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	isAllowed := a.Roles.Permissions().HasAny(*perm, rbac.PermissionAdministrator)
	if !isAllowed {
		return fault.New("required role not held", fctx.With(ctx), ftag.With(ftag.PermissionDenied))
	}

	return nil
}

// TODO: Use Scopes field of OpenAPI security spec.
//
// GetPermissionForOperation maps an operation ID to a permission requirement.
// Most operations are quite simple and just require a permission to operate but
// other operations have special requirements that are implemented elsewhere.
//
// There are three kinds of return state from this check function:
// 1. Publicly accessible with no account (false, nil)
// 2. Requires a session but no specific permission (true, nil)
// 3. Requires a specific permission (true, rbac.Permission{...})
//
// NOTE: Some operations are checked in the service layer instead of here due to
// additional logic required to determine attributes such as resource ownership.
func GetPermissionForOperation(operationID string) (requiresSession bool, requiresPermission *rbac.Permission) {
	return openapi_rbac.GetOperationPermission(&openapi_rbac.Mapping{}, operationID)
}
