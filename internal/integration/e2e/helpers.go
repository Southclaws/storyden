package e2e

import (
	"context"
	"net/http"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/transports/openapi/bindings"
	"github.com/Southclaws/storyden/internal/openapi"
)

func WithAccount(ctx context.Context, ar account.Repository, template account.Account, opts ...account.Option) (context.Context, *account.Account) {
	unique := xid.New()
	template.ID = account.AccountID(unique)
	template.Handle = template.Handle + "-" + unique.String()

	opts = append(opts, account.WithID(template.ID), account.WithName(template.Name), account.WithAdmin(template.Admin))

	acc, err := ar.Create(ctx, template.Handle, opts...)
	if err != nil {
		panic(err)
	}

	ctx = authentication.WithAccountID(ctx, acc.ID)
	return ctx, acc
}

func WithSession(ctx context.Context, cj *bindings.CookieJar) openapi.RequestEditorFn {
	accountID, err := authentication.GetAccountID(ctx)
	if err != nil {
		panic(err)
	}

	return func(ctx context.Context, req *http.Request) error {
		req.AddCookie(cj.Create(accountID.String()))
		return nil
	}
}
