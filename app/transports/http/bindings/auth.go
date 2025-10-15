package bindings

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/account/authentication/access_key"
	"github.com/Southclaws/storyden/app/resources/account/email"
	"github.com/Southclaws/storyden/app/resources/account/token"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/settings"
	auth_svc "github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/services/authentication/email_verify"
	"github.com/Southclaws/storyden/app/services/authentication/provider/email_only"
	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth"
	"github.com/Southclaws/storyden/app/services/authentication/provider/password"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
)

type Authentication struct {
	logger                        *slog.Logger
	cj                            *session_cookie.Jar
	si                            *session.Issuer
	tokenRepo                     token.Repository
	settings                      *settings.SettingsRepository
	passwordAuthProvider          *password.Provider
	emailVerificationAuthProvider *email_only.Provider
	accountQuery                  *account_querier.Querier
	emailRepo                     *email.Repository
	authManager                   *auth_svc.Manager
	emailVerifier                 *email_verify.Verifier
	access_key                    *access_key.Repository
	webAddress                    url.URL
}

func NewAuthentication(
	cfg config.Config,
	logger *slog.Logger,
	cj *session_cookie.Jar,
	si *session.Issuer,
	tokenRepo token.Repository,
	settings *settings.SettingsRepository,
	passwordAuthProvider *password.Provider,
	emailVerificationAuthProvider *email_only.Provider,
	accountQuery *account_querier.Querier,
	emailRepo *email.Repository,
	authManager *auth_svc.Manager,
	emailVerifier *email_verify.Verifier,
	access_key *access_key.Repository,
) Authentication {
	return Authentication{
		logger:                        logger,
		cj:                            cj,
		si:                            si,
		tokenRepo:                     tokenRepo,
		settings:                      settings,
		passwordAuthProvider:          passwordAuthProvider,
		emailVerificationAuthProvider: emailVerificationAuthProvider,
		accountQuery:                  accountQuery,
		emailRepo:                     emailRepo,
		authManager:                   authManager,
		emailVerifier:                 emailVerifier,
		access_key:                    access_key,
		webAddress:                    cfg.PublicWebAddress,
	}
}

func (o *Authentication) AuthProviderList(ctx context.Context, request openapi.AuthProviderListRequestObject) (openapi.AuthProviderListResponseObject, error) {
	settings, err := o.settings.Get(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	providers, err := o.authManager.GetProviderList(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	list, err := dt.MapErr(providers, serialiseAuthProvider(buildRedirectURL(o.webAddress)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	mode := settings.AuthenticationMode.Or(authentication.ModeHandle)

	return openapi.AuthProviderList200JSONResponse{
		AuthProviderListOKJSONResponse: openapi.AuthProviderListOKJSONResponse{
			Providers: list,
			Mode:      openapi.AuthMode(mode.String()),
		},
	}, nil
}

func (a *Authentication) AuthProviderLogout(ctx context.Context, request openapi.AuthProviderLogoutRequestObject) (openapi.AuthProviderLogoutResponseObject, error) {
	redirectTo := a.webAddress

	if request.Params.Redirect != nil && *request.Params.Redirect != "" {
		parsed, err := url.Parse(*request.Params.Redirect)
		if err == nil {
			// Only use the path component to prevent open redirects.
			redirectTo.Path = parsed.Path
			redirectTo.RawQuery = parsed.RawQuery
			redirectTo.Fragment = parsed.Fragment
		}
	}

	if sessionToken, ok := session.GetSessionToken(ctx).Get(); ok {
		t, err := token.FromString(sessionToken)
		if err == nil {
			if err := a.tokenRepo.Revoke(ctx, t); err != nil {
				a.logger.Warn("failed to revoke session token on logout",
					slog.String("error", err.Error()))
			}
		}
	}

	return openapi.AuthProviderLogout302Response{
		Headers: openapi.AuthProviderLogout302ResponseHeaders{
			SetCookie:     a.cj.Destroy().String(),
			ClearSiteData: `"cache", "cookies", "storage", "executionContexts"`,
			CacheControl:  "no-cache, no-store, must-revalidate",
			Location:      redirectTo.String(),
		},
	}, nil
}

func (a *Authentication) AccessKeyList(ctx context.Context, request openapi.AccessKeyListRequestObject) (openapi.AccessKeyListResponseObject, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionUsePersonalAccessKeys); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	accID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	list, err := a.access_key.List(ctx, accID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccessKeyList200JSONResponse{
		AccessKeyListOKJSONResponse: openapi.AccessKeyListOKJSONResponse{
			Keys: serialiseAccessKeyList(list),
		},
	}, nil
}

func (a *Authentication) AccessKeyCreate(ctx context.Context, request openapi.AccessKeyCreateRequestObject) (openapi.AccessKeyCreateResponseObject, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionUsePersonalAccessKeys); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	accID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	aks, err := a.access_key.Create(ctx, accID, access_key.AccessKeyKindPersonal, request.Body.Name, opt.NewPtr(request.Body.ExpiresAt))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccessKeyCreate200JSONResponse{
		AccessKeyCreateOKJSONResponse: openapi.AccessKeyCreateOKJSONResponse(openapi.AccessKeyIssued{
			Id:        openapi.Identifier(aks.AuthID.String()),
			CreatedAt: aks.CreatedAt,
			ExpiresAt: aks.Expires.Ptr(),
			Name:      aks.Name,
			Secret:    aks.String(),
		}),
	}, nil
}

func (a *Authentication) AccessKeyDelete(ctx context.Context, request openapi.AccessKeyDeleteRequestObject) (openapi.AccessKeyDeleteResponseObject, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionUsePersonalAccessKeys); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	accID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	_, err = a.access_key.Revoke(ctx, accID, deserialiseID(request.AccessKeyId))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccessKeyDelete204Response{}, nil
}

func buildRedirectURL(webAddress url.URL) func(s authentication.Service) url.URL {
	return func(s authentication.Service) url.URL {
		return oauth.Redirect(webAddress, s)
	}
}

func serialiseAuthProvider(redirectFn func(authentication.Service) url.URL) func(p auth_svc.Provider) (openapi.AuthProvider, error) {
	return func(p auth_svc.Provider) (openapi.AuthProvider, error) {
		if op, ok := p.(auth_svc.OAuthProvider); ok {
			uri := redirectFn(p.Service())

			link, err := op.Link(uri.String())
			if err != nil {
				return openapi.AuthProvider{}, fault.Wrap(err)
			}
			return openapi.AuthProvider{
				Provider: p.Service().String(),
				Name:     fmt.Sprintf("%v", p.Service()),
				Link:     &link,
			}, nil
		}

		return openapi.AuthProvider{
			Provider: p.Service().String(),
			Name:     fmt.Sprintf("%v", p.Service()),
		}, nil
	}
}

func deserialiseAuthMode(in openapi.AuthMode) (authentication.Mode, error) {
	mode, err := authentication.NewMode(string(in))
	if err != nil {
		return authentication.Mode{}, fault.Wrap(err, ftag.With(ftag.InvalidArgument))
	}
	return mode, nil
}

func serialiseAccessKey(k *authentication.Authentication) openapi.AccessKey {
	return openapi.AccessKey{
		Id:        k.ID.String(),
		CreatedAt: k.Created,
		ExpiresAt: k.Expires.Ptr(),
		Enabled:   !k.Disabled,
		Name:      k.Name.Or("Unnamed key"),
	}
}

func serialiseAccessKeyList(list []*authentication.Authentication) []openapi.AccessKey {
	return dt.Map(list, serialiseAccessKey)
}
