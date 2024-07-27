package bindings

import (
	"context"
	"errors"
	"net/mail"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/services/authentication/provider/email_only"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/openapi"
)

func (i *Authentication) AuthEmailSignup(ctx context.Context, request openapi.AuthEmailSignupRequestObject) (openapi.AuthEmailSignupResponseObject, error) {
	address, err := mail.ParseAddress(request.Body.Email)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	handle := opt.NewPtr(request.Body.Handle)

	acc, err := i.ep.Register(ctx, *address, handle)
	if err != nil {
		// SPEC: If the email exists, return a 422 response with no session.
		if errors.Is(err, email_only.ErrAccountAlreadyExists) {
			return openapi.AuthEmailSignup422Response{}, nil
		}

		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthEmailSignup200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccessOK{Id: acc.ID.String()},
			Headers: openapi.AuthSuccessOKResponseHeaders{
				SetCookie: i.sm.Create(acc.ID.String()).String(),
			},
		},
	}, nil
}

func (i *Authentication) AuthEmailSignin(ctx context.Context, request openapi.AuthEmailSigninRequestObject) (openapi.AuthEmailSigninResponseObject, error) {
	// i.ep.Login(ctx, request.Body.Email)
	return nil, nil
}

func (i *Authentication) AuthEmailVerify(ctx context.Context, request openapi.AuthEmailVerifyRequestObject) (openapi.AuthEmailVerifyResponseObject, error) {
	accountID, err := func() (string, error) {
		sessionAccountID := session.GetOptAccountID(ctx)
		requestEmailAddress, err := opt.MapErr(opt.NewPtr(request.Body.Email), func(s string) (mail.Address, error) {
			e, err := mail.ParseAddress(s)
			if err != nil {
				return mail.Address{}, err
			}

			return *e, nil
		})
		if err != nil {
			return "", fault.Wrap(err, fctx.With(ctx))
		}

		accountID, hasAccountID := sessionAccountID.Get()
		emailAddress, hasEmail := requestEmailAddress.Get()

		switch {
		// SPEC: The verification is being made with an active session.
		case hasAccountID && !hasEmail:
			return accountID.String(), nil

		// SPEC: The verification is being made without an active session.
		case !hasAccountID && hasEmail:
			acc, exists, err := i.er.LookupAccount(ctx, emailAddress)
			if err != nil {
				return "", fault.Wrap(err, fctx.With(ctx))
			}
			if !exists {
				return "", fault.New("account not found", fctx.With(ctx), ftag.With(ftag.NotFound))
			}

			return acc.ID.String(), nil

		case hasAccountID && hasEmail:
			return "", fault.New("both account ID and email provided", fctx.With(ctx), ftag.With(ftag.InvalidArgument))

		default:
			return "", fault.New("neither account ID nor email provided", fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}
	}()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := i.ep.Login(ctx, accountID, request.Body.Code)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthEmailVerify200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccessOK{Id: acc.ID.String()},
			Headers: openapi.AuthSuccessOKResponseHeaders{
				SetCookie: i.sm.Create(acc.ID.String()).String(),
			},
		},
	}, nil
}
