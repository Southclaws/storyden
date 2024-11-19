package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/rueidis"
)

var errNotFound = fmt.Errorf("not found")

type RedisCache struct {
	client  rueidis.Client
	options Options
}

type Options struct {
	Expiration                time.Duration
	ClientSideCacheExpiration time.Duration
}

func New(client rueidis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	cmd := c.client.B().Get().Key(key).Cache()
	res := c.client.DoCache(ctx, cmd, c.options.ClientSideCacheExpiration)

	str, err := res.ToString()
	if rueidis.IsRedisNil(err) {
		return "", errNotFound
	}

	return str, err
}

func (c *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	cmd := c.client.B().
		Set().
		Key(key).
		Value(value).
		Ex(ttl).
		Build()

	err := c.client.Do(ctx, cmd).Error()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	cmd := c.client.B().
		Del().
		Key(key).
		Build()

	err := c.client.Do(ctx, cmd).Error()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisCache) HIncrBy(ctx context.Context, key string, field string, incr int64) (int, error) {
	cmd := c.client.B().
		Hincrby().
		Key(key).
		Field(field).
		Increment(incr).
		Build()

	r, err := c.client.Do(ctx, cmd).ToInt64()
	if err != nil {
		return 0, err
	}

	return int(r), nil
}

func (c *RedisCache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	cmd := c.client.B().
		Hgetall().
		Key(key).
		Build()

	r, err := c.client.Do(ctx, cmd).ToMap()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(r))
	for k, v := range r {
		result[k], _ = v.ToString()
	}

	return result, nil
}

func (c *RedisCache) HDel(ctx context.Context, key string, field string) error {
	cmd := c.client.B().
		Hdel().
		Key(key).
		Field(field).
		Build()

	err := c.client.Do(ctx, cmd).Error()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	cmd := c.client.B().
		Expire().
		Key(key).
		Seconds(int64(expiration.Seconds())).
		Build()

	err := c.client.Do(ctx, cmd).Error()
	if err != nil {
		return err
	}

	return nil
}
