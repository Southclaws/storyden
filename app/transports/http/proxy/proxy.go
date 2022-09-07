package proxy

import (
	"context"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/config"
)

func Build() fx.Option {
	return fx.Options(
		fx.Invoke(mount),
	)
}

func mount(lc fx.Lifecycle, cfg config.Config, router *echo.Echo) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go run(ctx)

			return nil
		},
	})
}

func run(ctx context.Context) {
	// TODO: Implement standalone mode which means also setting up a separate
	// Docker image which contains the frontend dependencies too.

	// cmd := exec.CommandContext(ctx, "yarn", "dev")

	// cmd.Dir = "./web"
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	// err := cmd.Run()
	// if err != nil {
	// 	panic(err)
	// }
}
