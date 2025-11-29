package cache

import (
	"context"
	"time"

	"github.com/redis/rueidis"
	"go.uber.org/fx"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/cache/local"
	"github.com/Southclaws/storyden/internal/infrastructure/cache/redis"
)

type Store interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, object string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error

	HIncrBy(ctx context.Context, key string, field string, incr int64) (int, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, field string) error
	Expire(ctx context.Context, key string, expiration time.Duration) error
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(func(cfg config.Config, redisClient rueidis.Client) (Store, error) {
			switch cfg.CacheProvider {
			case "":
				c, err := local.New()
				return c, err

			case "redis":
				if redisClient == nil {
					return nil, fault.New("REDIS_URL is required when CACHE_PROVIDER is set to 'redis'")
				}

				return redis.New(redisClient), nil
			}

			panic("unknown cache provider: " + cfg.CacheProvider)
		}),
	)
}
