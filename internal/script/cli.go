package script

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure"
	"github.com/Southclaws/storyden/pkg/resources"
	"github.com/Southclaws/storyden/pkg/services"
)

// Run is a quick helper for writing scripts that use services.
func Run(opts ...fx.Option) {
	ctx, cf := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cf()

	all := append(opts, []fx.Option{
		fx.NopLogger,

		fx.Invoke(func(cfg config.Config) error {
			if strings.Contains(cfg.DatabaseURL, "development") || strings.Contains(cfg.DatabaseURL, "production") {
				return errors.New("refusing to run script on live database")
			}
			return nil
		}),

		config.Build(),
		infrastructure.Build(),
		services.Build(),
		resources.Build(),
	}...)

	app := fx.New(all...)

	if err := app.Start(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
