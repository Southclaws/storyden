package e2e

import (
	"context"
	"net/http"
	"testing"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/token"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func WithAccount(ctx context.Context, aw *account_writer.Writer, template account.Account, opts ...account_writer.Option) (context.Context, *account.Account) {
	unique := xid.New()
	template.ID = account.AccountID(unique)
	template.Handle = template.Handle + "-" + unique.String()

	opts = append(opts,
		account_writer.WithID(template.ID),
		account_writer.WithName(template.Name),
		account_writer.WithAdmin(template.Admin),
	)

	acc, err := aw.Create(ctx, template.Handle, opts...)
	if err != nil {
		panic(err)
	}

	ctx = session.WithAccountID(ctx, acc.ID)
	return ctx, &acc.Account
}

func WithSessionFromHeader(t *testing.T, ctx context.Context, header http.Header) openapi.RequestEditorFn {
	cookieHeader := header.Get("Set-Cookie")
	cookie, err := http.ParseSetCookie(cookieHeader)
	if err != nil {
		t.Error("failed to parse cookies from header", err)
	}

	return func(ctx context.Context, req *http.Request) error {
		req.AddCookie(cookie)
		return nil
	}
}

type SessionHelper struct {
	cj *session_cookie.Jar
	tr token.Repository
}

func newSessionHelper(
	cj *session_cookie.Jar,
	tr token.Repository,
) *SessionHelper {
	return &SessionHelper{
		cj: cj,
		tr: tr,
	}
}

func (h *SessionHelper) WithSession(ctx context.Context) openapi.RequestEditorFn {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		panic(err)
	}

	t, err := h.tr.Issue(ctx, accountID)
	if err != nil {
		panic(err)
	}

	return func(ctx context.Context, req *http.Request) error {
		req.AddCookie(h.cj.Create(t.Token))
		return nil
	}
}
