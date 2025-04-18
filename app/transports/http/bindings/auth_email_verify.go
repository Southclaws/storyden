package bindings

import (
	"context"
	"net/mail"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func (i *Authentication) AuthEmailVerify(ctx context.Context, request openapi.AuthEmailVerifyRequestObject) (openapi.AuthEmailVerifyResponseObject, error) {
	email, err := mail.ParseAddress(strings.ToLower(request.Body.Email))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	acc, err := i.emailVerifier.Verify(ctx, *email, request.Body.Code)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	t, err := i.si.Issue(ctx, acc.ID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthEmailVerify200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccessOK{Id: acc.ID.String()},
			Headers: openapi.AuthSuccessOKResponseHeaders{
				SetCookie: i.cj.Create(*t).String(),
			},
		},
	}, nil
}
