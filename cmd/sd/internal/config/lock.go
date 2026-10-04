package config

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofrs/flock"
)

func (s *Store) WithContextLock(ctx context.Context, name string, fn func() error) error {
	return s.lock(ctx, fmt.Sprintf(".%x", sha256.Sum256([]byte(name))), fn)
}

func (s *Store) lock(ctx context.Context, suffix string, fn func() error) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}
	l := flock.New(s.path+suffix+".lock", flock.SetPermissions(0o600))
	locked, err := l.TryLockContext(ctx, 25*time.Millisecond)
	if err != nil {
		return err
	}
	if !locked {
		return ctx.Err()
	}
	defer l.Close()
	return fn()
}

func (s *Store) write(data []byte) error {
	f, err := os.CreateTemp(filepath.Dir(s.path), ".config-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), s.path)
}
