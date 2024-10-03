package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	session1 "github.com/Southclaws/storyden/app/transports/http/middleware/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/script"
)

func main() {
	script.Run(
		fx.Provide(session1.New),
		fx.Invoke(func(
			ctx context.Context,
			ec *ent.Client,
			cj *session1.Jar,
		) (*struct{}, error) {
			if len(os.Args) < 2 {
				return nil, fault.New("no input file specified")
			}

			server := "http://localhost:8000/api"

			client, err := openapi.NewClientWithResponses(server)
			if err != nil {
				return nil, fault.Wrap(err)
			}

			fmt.Printf("Importing Hacker News dataset via Storyden API on %s...\n", server)

			f, err := os.Open(os.Args[1])
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			c, err := csv.NewReader(f).ReadAll()
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			session := func(ctx context.Context, username string) (openapi.RequestEditorFn, error) {
				signin, err := client.AuthPasswordSigninWithResponse(ctx, openapi.AuthPair{
					Identifier: username,
					Token:      "password",
				})
				if err != nil {
					return nil, fault.Wrap(err, fctx.With(ctx))
				}

				if signin.JSON200 != nil {
					ctx = session.WithAccountID(ctx, account.AccountID(openapi.ParseID(signin.JSON200.Id)))
					return e2e.WithSession(ctx, cj), nil
				}

				signup, err := client.AuthPasswordSignupWithResponse(ctx, nil, openapi.AuthPair{
					Identifier: username,
					Token:      "password",
				})
				if err != nil {
					return nil, fault.Wrap(err, fctx.With(ctx))
				}

				if signup.JSON200 != nil {
					ctx = session.WithAccountID(ctx, account.AccountID(openapi.ParseID(signup.JSON200.Id)))
					return e2e.WithSession(ctx, cj), nil
				}

				fmt.Println(signup.Status())

				return nil, fault.New("signup failed")
			}

			create := func(ctx context.Context, by string, title string, body string, url *string) error {
				s, err := session(ctx, by)
				if err != nil {
					return fault.Wrap(err, fctx.With(ctx))
				}

				thread, err := client.ThreadCreateWithResponse(ctx, openapi.ThreadInitialProps{
					Category:   "00000000000000000010",
					Title:      title,
					Body:       body,
					Url:        url,
					Visibility: openapi.Published,
				}, s)
				if err != nil {
					return fault.Wrap(err, fctx.With(ctx))
				}

				if thread.JSON200 == nil {
					return fault.Newf("failed to create thread: %s", thread.Status())
				}

				fmt.Println("Created thread with ID", thread.JSON200.Id)

				return nil
			}

			for _, v := range c[1:] {
				title := v[0]
				url := opt.NewIf(v[1], func(s string) bool { return s != "" })
				text := v[2]
				by := v[4]

				err = create(ctx, by, title, text, url.Ptr())
				if err != nil {
					return nil, fault.Wrap(err, fctx.With(ctx))
				}
			}

			return nil, nil
		}))
}
