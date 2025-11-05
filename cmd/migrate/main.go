package main

import (
	"context"

	"entgo.io/ent/dialect/sql/schema"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/script"
)

func main() {
	script.Run(fx.Invoke(func(ctx context.Context, client *ent.Client) {
		if err := client.Schema.Create(
			ctx,
			schema.WithDropIndex(true),
			schema.WithDropColumn(true),
		); err != nil {
			panic(err)
		}
	}))
}
