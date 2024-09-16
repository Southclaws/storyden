package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/email"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/services/authentication/email_verify"
	"github.com/Southclaws/storyden/app/services/authentication/provider/email/email_only"
	"github.com/Southclaws/storyden/app/services/authentication/provider/email/email_password"
	"github.com/Southclaws/storyden/app/services/authentication/provider/password"
	session1 "github.com/Southclaws/storyden/app/transports/http/middleware/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
)

type Authentication struct {
	p            *password.Provider
	ep           *email_only.Provider
	epp          *email_password.Provider
	cj           *session1.Jar
	accountQuery *account_querier.Querier
	er           email.EmailRepo
	am           *authentication.Manager
	ev           email_verify.Verifier
}

func NewAuthentication(
	cfg config.Config,
	p *password.Provider,
	ep *email_only.Provider,
	epp *email_password.Provider,
	accountQuery *account_querier.Querier,
	er email.EmailRepo,
	sm *session1.Jar,
	am *authentication.Manager,
	ev email_verify.Verifier,
) Authentication {
	return Authentication{p, ep, epp, sm, accountQuery, er, am, ev}
}

func (o *Authentication) AuthProviderList(ctx context.Context, request openapi.AuthProviderListRequestObject) (openapi.AuthProviderListResponseObject, error) {
	list, err := dt.MapErr(o.am.Providers(), serialiseAuthProvider)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthProviderList200JSONResponse{
		AuthProviderListOKJSONResponse: openapi.AuthProviderListOKJSONResponse{
			Providers: list,
		},
	}, nil
}

func (a *Authentication) AuthProviderLogout(ctx context.Context, request openapi.AuthProviderLogoutRequestObject) (openapi.AuthProviderLogoutResponseObject, error) {
	return openapi.AuthProviderLogout200Response{
		Headers: openapi.AuthProviderLogout200ResponseHeaders{
			SetCookie: a.cj.Destroy().String(),
		},
	}, nil
}

func serialiseAuthProvider(p authentication.Provider) (openapi.AuthProvider, error) {
	link, err := p.Link("/")
	if err != nil {
		return openapi.AuthProvider{}, fault.Wrap(err)
	}

	return openapi.AuthProvider{
		Provider: p.ID(),
		Name:     p.Name(),
		Link:     link,
	}, nil
}
