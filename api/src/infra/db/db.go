package db

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert --feature sql/modifier --feature sql/upsert ./schema --target ./model

import (
	"context"
	"database/sql"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/api/src/config"
	"github.com/Southclaws/storyden/api/src/infra/db/model"
)

func Build() fx.Option {
	return fx.Provide(func(lc fx.Lifecycle, cfg config.Config) (*model.Client, error) {
		var client *model.Client

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) (err error) {
				client, _, err = connect(cfg.DatabaseURL)
				if err != nil {
					return err
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				return client.Close()
			},
		})

		return client, nil
	})
}

func connect(url string) (*model.Client, *sql.DB, error) {
	driver, err := sql.Open("pgx", url)
	if err != nil {
		return nil, nil, err
	}

	client := model.NewClient(model.Driver(entsql.OpenDB(dialect.Postgres, driver)))

	opts := []schema.MigrateOption{
		schema.WithDropColumn(true),
		schema.WithDropIndex(true),
	}

	// Run only additive migrations
	if err := client.Schema.Create(context.Background(), opts...); err != nil {
		return nil, nil, err
	}

	return client, driver, nil
}
