package local

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/puzpuzpuz/xsync/v3"
)

var errNotFound = fmt.Errorf("not found")

type Entry struct {
	Value  any
	Expiry *time.Time
}

type LocalCache struct {
	local *xsync.MapOf[string, Entry]
}

func New() *LocalCache {
	return &LocalCache{
		local: xsync.NewMapOf[string, Entry](),
	}
}

func (c *LocalCache) Get(ctx context.Context, key string) (string, error) {
	v, found := c.local.Load(key)
	if !found {
		return "", errNotFound
	}

	if v.Expiry.Before(time.Now()) {
		c.local.Delete(key)
		return "", errNotFound
	}

	return v.Value.(string), nil
}

func (c *LocalCache) Set(ctx context.Context, key string, value string) error {
	c.local.Store(key, Entry{
		Value: value,
	})
	return nil
}

func (c *LocalCache) Delete(ctx context.Context, key string) error {
	c.local.Delete(key)
	return nil
}

func (c *LocalCache) HIncrBy(ctx context.Context, key string, field string, incr int64) (int, error) {
	ac, _ := c.local.Compute(key, func(old Entry, found bool) (Entry, bool) {
		if found {
			hash := old.Value.(map[string]string)
			if curr, ok := hash[field]; ok {
				i, err := strconv.Atoi(curr)
				if err != nil {
					return old, false
				}

				i += int(incr)
				hash[field] = strconv.Itoa(i)
				old.Value = hash
				return old, false
			} else {
				hash[field] = strconv.Itoa(int(incr))
				old.Value = hash
				return old, false
			}

		} else {
			hash := map[string]string{
				field: strconv.Itoa(int(incr)),
			}
			old.Value = hash
			return old, false
		}
	})

	hash := ac.Value.(map[string]string)
	i, err := strconv.Atoi(hash[field])
	if err != nil {
		return 0, fmt.Errorf("failed to convert hash field to integer")
	}

	return i, nil
}

func (c *LocalCache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	v, ok := c.local.Load(key)
	if !ok {
		return map[string]string{}, nil
	}

	return v.Value.(map[string]string), nil
}

func (c *LocalCache) HDel(ctx context.Context, key string, field string) error {
	_, ok := c.local.Compute(key, func(old Entry, found bool) (Entry, bool) {
		if found {
			hash := old.Value.(map[string]string)
			delete(hash, field)
			old.Value = hash
			return old, false
		}
		return old, false
	})
	if !ok {
		return fmt.Errorf("failed to delete hash field")
	}

	return nil
}

func (c *LocalCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	c.local.Compute(key, func(old Entry, found bool) (Entry, bool) {
		if found {
			if old.Expiry != nil && old.Expiry.Before(time.Now()) {
				return old, true
			}

			expiry := time.Now().Add(expiration)
			old.Expiry = &expiry
			return old, false
		}
		return old, false
	})

	return nil
}
