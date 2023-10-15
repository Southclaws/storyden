package main

import (
	"context"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/script"
)

func main() {
	script.Run(fx.Invoke(func(ctx context.Context, client *ent.Client) {
		if err := client.Schema.Create(ctx); err != nil {
			panic(err)
		}
	}))
}
