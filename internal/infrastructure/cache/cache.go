package cache

import (
	"context"
	"time"

	"github.com/redis/rueidis"
	"go.uber.org/fx"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/cache/local"
	"github.com/Southclaws/storyden/internal/infrastructure/cache/redis"
)

type Store interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, object string) error
	Delete(ctx context.Context, key string) error

	HIncrBy(ctx context.Context, key string, field string, incr int64) (int, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, field string) error
	Expire(ctx context.Context, key string, expiration time.Duration) error
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(func(cfg config.Config) (Store, error) {
			switch cfg.CacheProvider {
			case "":
				return local.New(), nil

			case "redis":
				password, _ := cfg.RedisURL.User.Password()

				client, err := rueidis.NewClient(rueidis.ClientOption{
					InitAddress:      []string{cfg.RedisURL.Host},
					Username:         cfg.RedisURL.User.Username(),
					Password:         password,
					DisableCache:     true,
					ConnWriteTimeout: 5 * time.Second,
				})
				if err != nil {
					return nil, fault.Wrap(err, fmsg.With("failed to connect to redis"))
				}

				return redis.New(client), nil
			}

			panic("unknown cache provider: " + cfg.CacheProvider)
		}),
	)
}
