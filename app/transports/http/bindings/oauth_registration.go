package bindings

import (
	"context"
	"net/http"
	"strings"

	"github.com/Southclaws/opt"
	"github.com/labstack/echo/v4"

	oauthservice "github.com/Southclaws/storyden/app/services/authentication/oauth"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type registrationAuthKey struct{}

func oauthRegistrationAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Path() == "/api/oauth/register" || (c.Path() == "/api/admin/oauth/registration-approvals" || c.Path() == "/api/admin/oauth/registration-approvals/verification") {
			c.Response().Header().Set("Cache-Control", "no-store")
		}

		if c.Request().Method == http.MethodPost && c.Path() == "/api/oauth/register" {
			ctx := context.WithValue(c.Request().Context(), registrationAuthKey{}, c.Request().Header.Values(echo.HeaderAuthorization))
			c.SetRequest(c.Request().WithContext(ctx))
		}

		return next(c)
	}
}

func registrationAccessToken(ctx context.Context) opt.Optional[string] {
	headers, _ := ctx.Value(registrationAuthKey{}).([]string)
	if len(headers) == 0 {
		return opt.NewEmpty[string]()
	}

	if len(headers) != 1 {
		return opt.New("")
	}

	parts := strings.Fields(headers[0])
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return opt.New("")
	}

	return opt.New(parts[1])
}

func registrationAccessError(err *oauthservice.Error, supplied bool) openapi.OAuthClientRegisterResponseObject {
	body := openapi.OAuthError{
		Error:            err.Code,
		ErrorDescription: &err.Description,
	}

	cache := "no-store"
	switch err.Code {
	case "invalid_token":
		challenge := `Bearer realm="oauth/register"`
		if supplied {
			challenge += `, error="invalid_token"`
		}

		return openapi.OAuthClientRegister401JSONResponse{
			Body: body,
			Headers: openapi.OAuthClientRegister401ResponseHeaders{
				CacheControl:    &cache,
				WWWAuthenticate: &challenge,
			},
		}
	case "insufficient_scope":
		challenge := `Bearer realm="oauth/register", error="insufficient_scope"`
		return openapi.OAuthClientRegister403JSONResponse{
			Body: body,
			Headers: openapi.OAuthClientRegister403ResponseHeaders{
				CacheControl:    &cache,
				WWWAuthenticate: &challenge,
			},
		}
	default:
		return nil
	}
}

func (o OAuth) pollOAuthRegistration(ctx context.Context, code string, cancel bool) (openapi.OAuthClientRegisterResponseObject, error) {
	if cancel {
		oauthErr, err := o.oauth.CancelRegistrationApproval(ctx, code)
		if err != nil {
			return nil, err
		}

		if oauthErr != nil {
			return registrationResponse(nil, oauthErr, false), nil
		}

		return openapi.OAuthClientRegister204Response{}, nil
	}

	result, oauthErr, err := o.oauth.PollRegistrationApproval(ctx, code)
	if err != nil {
		return nil, err
	}

	return registrationResponse(result, oauthErr, false), nil
}

func registrationResponse(result *oauthservice.DynamicClientRegistrationResult, oauthErr *oauthservice.Error, supplied bool) openapi.OAuthClientRegisterResponseObject {
	if oauthErr != nil {
		if response := registrationAccessError(oauthErr, supplied); response != nil {
			return response
		}

		return openapi.OAuthClientRegister400JSONResponse{OAuthClientRegisterErrorJSONResponse: openapi.OAuthClientRegisterErrorJSONResponse(openapi.OAuthError{
			Error:            oauthErr.Code,
			ErrorDescription: &oauthErr.Description,
		})}
	}

	if result.RetryAfter > 0 {
		return openapi.OAuthClientRegister429Response{Headers: openapi.OAuthClientRegister429ResponseHeaders{RetryAfter: &result.RetryAfter}}
	}

	if pending := result.Approval; pending != nil {
		body := openapi.OAuthClientRegister202JSONResponse{}
		if pending.RegistrationCode != "" {
			body.RegistrationCode = &pending.RegistrationCode
			body.VerificationCode = &pending.VerificationCode
			body.VerificationUri = &pending.VerificationURI
			body.VerificationUriComplete = &pending.VerificationURIComplete
			body.ExpiresIn = &pending.ExpiresIn
			body.Interval = &pending.Interval
		}

		return body
	}

	return openapi.OAuthClientRegister201JSONResponse{OAuthClientRegisterOKJSONResponse: openapi.OAuthClientRegisterOKJSONResponse(serialiseOAuthClientRegistration(result))}
}
