package local

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"strconv"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/shirou/gopsutil/v4/mem"
)

var errNotFound = fmt.Errorf("not found")

type LocalCache struct {
	cache *ristretto.Cache[string, []byte]
}

type HSet map[string]int

func init() {
	gob.Register(HSet{})
}

func New() (*LocalCache, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	// cache size is 25% of available memory
	// is this a good idea? who knows...
	maxCost := int64(vm.Available / 4)

	cache, err := ristretto.NewCache(&ristretto.Config[string, []byte]{
		NumCounters:            1e7,
		MaxCost:                maxCost,
		BufferItems:            64,
		TtlTickerDurationInSec: 30,
	})
	if err != nil {
		return nil, err
	}

	return &LocalCache{
		cache: cache,
	}, nil
}

func (c *LocalCache) Get(ctx context.Context, key string) (string, error) {
	v, found := c.cache.Get(key)
	if !found {
		return "", errNotFound
	}

	return string(v), nil
}

func (c *LocalCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	c.cache.SetWithTTL(key, []byte(value), 0, ttl)
	c.cache.Wait()
	return nil
}

func (c *LocalCache) Delete(ctx context.Context, key string) error {
	c.cache.Del(key)
	return nil
}

func (c *LocalCache) HIncrBy(ctx context.Context, key string, field string, incr int64) (int, error) {
	hash, exists, err := c.getHSET(key)
	if err != nil {
		return 0, err
	}

	if exists {
		if i, ok := hash[field]; ok {
			next := i + int(incr)
			hash[field] = next
			err := c.setHSET(key, hash)
			return next, err
		} else {
			hash[field] = int(incr)
			err := c.setHSET(key, hash)
			return int(incr), err
		}
	} else {
		hash := map[string]int{
			field: int(incr),
		}
		err := c.setHSET(key, hash)
		return int(incr), err
	}
}

func (c *LocalCache) getHSET(key string) (HSet, bool, error) {
	v, ok := c.cache.Get(key)
	if !ok {
		return nil, false, nil
	}

	var hset HSet
	err := gob.NewDecoder(bytes.NewBuffer(v)).Decode(&hset)
	if err != nil {
		return nil, false, err
	}

	return hset, true, nil
}

func (c *LocalCache) setHSET(key string, hset HSet) error {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(hset)
	if err != nil {
		return err
	}

	c.cache.Set(key, buf.Bytes(), 0)
	return nil
}

func (c *LocalCache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	v, ok, err := c.getHSET(key)
	if err != nil {
		return nil, err
	}

	if !ok {
		return map[string]string{}, nil
	}

	// in my infinite wisedom i wrote the rate limit store to use strings as
	// values instead of integers, something to do with redis i guess... oh well
	ms := map[string]string{}
	for k, v := range v {
		ms[k] = strconv.Itoa(v)
	}

	return ms, nil
}

func (c *LocalCache) HDel(ctx context.Context, key string, field string) error {
	hash, exists, err := c.getHSET(key)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	_, ok := hash[field]
	if !ok {
		return nil
	}

	delete(hash, field)

	if len(hash) == 0 {
		c.cache.Del(key)
	}

	return c.setHSET(key, hash)
}

func (c *LocalCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	v, exists := c.cache.Get(key)
	if exists {
		c.cache.SetWithTTL(key, v, 0, expiration)
	}

	return nil
}
