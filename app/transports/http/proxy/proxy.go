package proxy

import (
	"context"
	"net/url"
	"os"
	"os/exec"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/utils"
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

			router.GET("/test", func(c echo.Context) error {
				c.JSON(200, "hello")
				return nil
			})

			router.Use(middleware.Proxy(middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{{
				URL: utils.Must(url.Parse("http://localhost:3000/")),
			}})))

			return nil
		},
	})
}

func run(ctx context.Context) {
	// TODO: format logs properly
	cmd := exec.CommandContext(ctx, "yarn", "dev")

	cmd.Dir = "./web"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
