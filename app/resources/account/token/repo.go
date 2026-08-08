package token

import (
	"context"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/infrastructure/cache"
)

const cachePrefix = "account:session:"

type Repository interface {
	Issue(context.Context, account.AccountID) (*Session, error)
	Revoke(context.Context, Token) error
	Validate(context.Context, Token) (*Validated, error)
}

type cachedRepo struct {
	repo  Repository
	store cache.Store
}

func NewCachedRepository(repo Repository, store cache.Store) Repository {
	return &cachedRepo{
		repo:  repo,
		store: store,
	}
}

func (r *cachedRepo) Issue(ctx context.Context, accountID account.AccountID) (*Session, error) {
	s, err := r.repo.Issue(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := r.cache(ctx, *s); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return s, nil
}

func (r *cachedRepo) Revoke(ctx context.Context, token Token) error {
	if err := r.delete(ctx, token); err != nil {
		return err
	}

	if err := r.repo.Revoke(ctx, token); err != nil {
		return err
	}

	if err := r.delete(ctx, token); err != nil {
		return err
	}

	return nil
}

func (r *cachedRepo) Validate(ctx context.Context, t Token) (*Validated, error) {
	sess, found, err := r.get(ctx, t)
	if err != nil {
		// an unreadable entry is evicted so the database decides instead
		if delErr := r.delete(ctx, t); delErr != nil {
			return nil, fault.Wrap(delErr, fctx.With(ctx))
		}
	} else if found {
		v, err := sess.Validate()
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return v, nil
	}

	// Fall back to database query.
	v, err := r.repo.Validate(ctx, t)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// Store in cache for future validations.
	if err := r.cache(ctx, Session(*v)); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return v, nil
}

func (r *cachedRepo) get(ctx context.Context, t Token) (*Session, bool, error) {
	raw, err := r.store.Get(ctx, cacheKey(t))
	if err != nil {
		// Cache miss, found=false
		// TODO: Expose a "cache miss" error/return value and distinguish
		// between network/cache errors and cache misses.
		return nil, false, nil
	}

	session, err := Deserialise([]byte(raw))
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	return session, true, nil
}

func (r *cachedRepo) cache(ctx context.Context, s Session) error {
	payload, err := s.Serialise()
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	ttl := time.Until(s.ExpiresAt)
	if ttl <= 0 {
		return nil
	}

	err = r.store.Set(ctx, cacheKey(s.Token), string(payload), ttl)
	if err != nil {
		return err
	}

	return nil
}

func (r *cachedRepo) delete(ctx context.Context, token Token) error {
	err := r.store.Delete(ctx, cacheKey(token))
	if err != nil {
		return err
	}

	return nil
}

func cacheKey(t Token) string {
	return cachePrefix + t.String()
}
