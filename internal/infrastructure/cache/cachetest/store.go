package cachetest

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/Southclaws/storyden/internal/infrastructure/cache"
)

var _ cache.Store = (*Store)(nil)

type Store struct {
	mu     sync.Mutex
	values map[string]string
	hashes map[string]map[string]int
}

func New() *Store {
	return &Store{
		values: map[string]string{},
		hashes: map[string]map[string]int{},
	}
}

func (s *Store) Get(ctx context.Context, key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.values[key]
	if !ok {
		return "", errors.New("not found")
	}

	return value, nil
}

func (s *Store) Set(ctx context.Context, key string, object string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.values[key] = object
	return nil
}

func (s *Store) Delete(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.values, key)
	delete(s.hashes, key)
	return nil
}

func (s *Store) HIncrBy(ctx context.Context, key string, field string, incr int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash, ok := s.hashes[key]
	if !ok {
		hash = map[string]int{}
		s.hashes[key] = hash
	}

	hash[field] += int(incr)
	return hash[field], nil
}

func (s *Store) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := map[string]string{}
	for field, value := range s.hashes[key] {
		out[field] = strconv.Itoa(value)
	}

	return out, nil
}

func (s *Store) HDel(ctx context.Context, key string, field string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash, ok := s.hashes[key]
	if !ok {
		return nil
	}

	delete(hash, field)
	if len(hash) == 0 {
		delete(s.hashes, key)
	}

	return nil
}

func (s *Store) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return nil
}
