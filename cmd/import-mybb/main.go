package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/openapi/bindings"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/openapi"
	"github.com/Southclaws/storyden/internal/script"
)

type UserRow struct {
	Username string `db:"username"`
	Email    string `db:"email"`
}

const UsersQuery = `
	select username, email from mybb_users
`

type PostRow struct {
	Dateline int    `db:"dateline"`
	PostID   int    `db:"pid"`
	ReplyTo  int    `db:"replyto"`
	Subject  string `db:"subject"`
	Username string `db:"username"`
	Message  string `db:"message"`
}

const PostsQuery = `
	select dateline, pid, replyto, subject, username, message from mybb_posts
`

func main() {
	script.Run(
		fx.Provide(bindings.NewCookieJar),
		fx.Invoke(func(
			ctx context.Context,
			ec *ent.Client,
			cj *bindings.CookieJar,
		) (*struct{}, error) {
			if len(os.Args) < 2 {
				return nil, fault.New("no database specified")
			}

			server := "http://localhost:8000/api"

			client, err := openapi.NewClientWithResponses(server)
			if err != nil {
				return nil, fault.Wrap(err)
			}

			fmt.Printf("Importing from MyBB via Storyden API on %s...\n", server)

			databaseURL := os.Args[1]

			db, err := sql.Open("pgx", databaseURL)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			postResults, err := db.QueryContext(ctx, PostsQuery)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}

			rows := []PostRow{}
			for postResults.Next() {
				row := PostRow{}
				err := postResults.Scan(&row.Dateline, &row.PostID, &row.ReplyTo, &row.Subject, &row.Username, &row.Message)
				if err != nil {
					return nil, fault.Wrap(err, fctx.With(ctx))
				}

				rows = append(rows, row)
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
					fmt.Println("SIGNING IN AS", username)
					ctx = session.WithAccountID(ctx, account.AccountID(openapi.ParseID(signin.JSON200.Id)))
					return e2e.WithSession(ctx, cj), nil
				}

				fmt.Println("USERNAME", username)

				signup, err := client.AuthPasswordSignupWithResponse(ctx, openapi.AuthPair{
					Identifier: username,
					Token:      "password",
				})
				if err != nil {
					return nil, fault.Wrap(err, fctx.With(ctx))
				}

				if signup.JSON200 != nil {
					fmt.Println("SIGNING UP AS", signup.JSON200.Id, username)
					ctx = session.WithAccountID(ctx, account.AccountID(openapi.ParseID(signup.JSON200.Id)))
					return e2e.WithSession(ctx, cj), nil
				}

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

			// reply := func(ctx context.Context, by string, threadID string, body string) error {
			// 	return nil
			// }

			for _, v := range rows {
				if v.ReplyTo == 0 {
					err = create(ctx, v.Username, v.Subject, v.Message, nil)
					if err != nil {
						fmt.Println("ERROR:", err)
						return nil, err
					}
				} else {
					// err = reply(ctx, v.Username, v.ReplyTo, v.Message)
				}
			}

			return nil, nil
		}))
}
