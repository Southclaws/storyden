package local

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/dgraph-io/ristretto/v2/z"
	"github.com/shirou/gopsutil/v4/mem"
)

var (
	errNotFound      = fmt.Errorf("not found")
	errWriteRejected = fmt.Errorf("cache write rejected")
)

const cacheLockCount = 256

type LocalCache struct {
	// stripes for CAS/SIA
	locks [cacheLockCount]sync.Mutex
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
	lock := &c.locks[c.lockIndex(key)]
	lock.Lock()
	defer lock.Unlock()

	if err := c.setWithTTL(key, []byte(value), ttl); err != nil {
		return err
	}
	c.cache.Wait()
	return nil
}

func (c *LocalCache) SetIfAbsent(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	if ttl <= 0 {
		return false, fmt.Errorf("cache TTL must be positive")
	}

	lock := &c.locks[c.lockIndex(key)]
	lock.Lock()
	defer lock.Unlock()

	if err := ctx.Err(); err != nil {
		return false, err
	}

	if _, exists := c.cache.Get(key); exists {
		return false, nil
	}

	if err := c.setWithTTL(key, []byte(value), ttl); err != nil {
		return false, err
	}
	c.cache.Wait()

	if _, exists := c.cache.Get(key); !exists {
		return false, errWriteRejected
	}

	return true, nil
}

func (c *LocalCache) SetMany(ctx context.Context, values map[string]string, ttl time.Duration) error {
	var needed [cacheLockCount]bool
	for key := range values {
		needed[c.lockIndex(key)] = true
	}

	for i, used := range needed {
		if used {
			c.locks[i].Lock()
			defer c.locks[i].Unlock()
		}
	}

	for key, value := range values {
		if err := c.setWithTTL(key, []byte(value), ttl); err != nil {
			c.cache.Wait()
			return err
		}
	}
	c.cache.Wait()
	return nil
}

func (c *LocalCache) setWithTTL(key string, value []byte, ttl time.Duration) error {
	if c.cache.SetWithTTL(key, value, 0, ttl) {
		return nil
	}

	// A full write buffer rejects new entries. Drain accepted writes once and
	// retry before reporting the rejection to the caller.
	c.cache.Wait()
	if c.cache.SetWithTTL(key, value, 0, ttl) {
		return nil
	}

	return fmt.Errorf("%w for key %q", errWriteRejected, key)
}

func (c *LocalCache) Delete(ctx context.Context, key string) error {
	lock := &c.locks[c.lockIndex(key)]
	lock.Lock()
	defer lock.Unlock()

	c.cache.Del(key)
	c.cache.Wait()
	return nil
}

func (c *LocalCache) HIncrBy(ctx context.Context, key string, field string, incr int64) (int, error) {
	lock := &c.locks[c.lockIndex(key)]
	lock.Lock()
	defer lock.Unlock()

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

	if err := c.setWithTTL(key, buf.Bytes(), 0); err != nil {
		return err
	}
	c.cache.Wait()
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
	lock := &c.locks[c.lockIndex(key)]
	lock.Lock()
	defer lock.Unlock()

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
	lock := &c.locks[c.lockIndex(key)]
	lock.Lock()
	defer lock.Unlock()

	v, exists := c.cache.Get(key)
	if exists {
		if err := c.setWithTTL(key, v, expiration); err != nil {
			return err
		}
		c.cache.Wait()
	}

	return nil
}

func (c *LocalCache) lockIndex(key string) uint64 {
	hash, _ := z.KeyToHash(key)
	return hash % cacheLockCount
}
