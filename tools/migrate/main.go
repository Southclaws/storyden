package main

import (
	"context"
	"fmt"
	"log"
	"os"

	atlas "ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/lib/pq"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/migrate"
	"github.com/Southclaws/storyden/internal/script"
)

func main() {
	script.Run(fx.Invoke(func(cfg config.Config) {
		ctx := context.Background()

		dir, err := atlas.NewLocalDir("migrations")
		if err != nil {
			log.Fatalf("failed creating atlas migration directory: %v", err)
		}

		if len(os.Args) != 2 {
			log.Fatalln("migration name is required")
		}

		err = migrate.NamedDiff(ctx, cfg.DatabaseURL, os.Args[1],
			schema.WithSumFile(),                        // hash migrations
			schema.WithDir(dir),                         // provide migration directory
			schema.WithMigrationMode(schema.ModeReplay), // provide migration mode
			schema.WithDialect(dialect.Postgres),        // Ent dialect to use
			schema.WithFormatter(atlas.DefaultFormatter),
		)
		if err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}

		fmt.Println("Done!")
	}))
}
