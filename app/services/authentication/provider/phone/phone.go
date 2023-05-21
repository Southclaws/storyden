package phone

import (
	"context"
	"crypto/rand"
	"fmt"
	"math"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/authentication"
	"github.com/Southclaws/storyden/internal/sms"
)

var (
	errHandleMismatch      = fault.New("phone already linked to different account")
	errNoPhoneAuth         = fault.New("no phone auth method linked to account")
	errOneTimeCodeMismatch = fault.New("one time code mismatch")
)

const (
	id   = "phone"
	name = "Phone"
	logo = ""
)

const template = `Your unique one-time login code is: %s`

type Provider struct {
	auth    authentication.Repository
	account account.Repository

	sms sms.Sender
}

func New(auth authentication.Repository, account account.Repository, sms sms.Sender) *Provider {
	return &Provider{auth, account, sms}
}

func (p *Provider) Enabled() bool   { return p.sms != nil }
func (p *Provider) ID() string      { return id }
func (p *Provider) Name() string    { return name }
func (p *Provider) LogoURL() string { return logo }

func (p *Provider) Register(ctx context.Context, handle string, phone string) (*account.Account, error) {
	//
	// STEP 1.
	//
	// Using the provided phone number, look up an authentication record which
	// points to an account already registered with the system. We need to do
	// this because there's no separation between registration and login via the
	// phone login system so if there's an account already, we start auth again.
	//

	authrecord, exists, err := p.auth.LookupByIdentifier(ctx, id, phone)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to get account"))
	}

	var acc *account.Account
	if exists {
		acc = &authrecord.Account
		if acc.Handle != handle {
			return nil, fault.Wrap(errHandleMismatch,
				fctx.With(ctx),
				ftag.With(ftag.PermissionDenied),
				fmsg.WithDesc("handle mismatch", "Handle already registered with a different authentication method."),
			)
		}

		//
		// STEP 1.5:
		//
		// If an account already exists, there's a chance the account also has a
		// phone authentication record associated with it. Currently, we only
		// support a single phone associated with an account so if there is one,
		// it needs to be deleted so it can be created again with a new code.
		//

		auths, err := p.auth.GetAuthMethods(ctx, acc.ID)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		// If there's already a phone auth associated with the account, deleted it
		// and start fresh with the new request.
		// NOTE: This could result in a DoS for the account holder...
		if _, exists = lo.Find(auths, func(a authentication.Authentication) bool {
			return a.Service == id
		}); exists {
			_, err = p.auth.Delete(ctx, acc.ID, phone, id)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
		}

	} else {
		//
		// If there isn't an account already with this phone number, we create
		// a new one using the @handle specified in the request.
		//
		acc, err = p.account.Create(ctx, handle)
		if err != nil {
			if ftag.Get(err) == ftag.AlreadyExists {
				return nil, fault.Wrap(err,
					fctx.With(ctx),
					fmsg.With("failed to create account"),
					fmsg.WithDesc("already exists", "Handle already registered with a different authentication method."))
			}
			return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create account"))
		}
	}

	//
	// STEP 2:
	//
	// Generate a one-time-password which is a 6 digit number and send this to
	// the phone number specified in the request.
	//

	code, err := generateCode()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to generate code"))
	}

	_, err = p.auth.Create(ctx, acc.ID, id, phone, code, nil)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create account authentication instance"))
	}

	// TODO: For whitelabling, allow the instance brand name to be specified in
	// the message template. So the message says "Log in to Acme with xyz..."
	message := fmt.Sprintf(template, code)
	err = p.sms.Send(ctx, phone, message)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, nil
}

func (p *Provider) Link() string {
	// Phone provider does not use external links.
	return ""
}

func (p *Provider) Login(ctx context.Context, handle string, onetimecode string) (*account.Account, error) {
	acc, err := p.account.GetByHandle(ctx, handle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	auths, err := p.auth.GetAuthMethods(ctx, acc.ID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	phoneauth, exists := lo.Find(auths, func(a authentication.Authentication) bool {
		return a.Service == id
	})
	if !exists {
		return nil, fault.Wrap(errNoPhoneAuth)
	}

	if phoneauth.Token != onetimecode {
		return nil, fault.Wrap(errOneTimeCodeMismatch,
			fctx.With(ctx),
			ftag.With(ftag.PermissionDenied),
			fmsg.WithDesc("mismatch", "The code did not match."),
		)
	}

	return acc, nil
}

func generateCode() (string, error) {
	sum := make([]byte, 6)
	_, err := rand.Read(sum)
	if err != nil {
		return "", fault.Wrap(err)
	}

	value := int64(((int(sum[0]) & 0x7f) << 24) |
		((int(sum[1] & 0xff)) << 16) |
		((int(sum[2] & 0xff)) << 8) |
		(int(sum[3]) & 0xff))

	mod := int32(value % int64(math.Pow10(6)))

	return fmt.Sprintf("%06d", mod), nil
}
