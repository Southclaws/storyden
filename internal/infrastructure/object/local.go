package object

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/internal/config"
)

type localStorer struct {
	s    fs.FS
	path string
}

func NewLocalStorer(cfg config.Config) Storer {
	var path string
	if cfg.AssetStorageLocalPath != "" {
		path = cfg.AssetStorageLocalPath
	} else {
		path = "./data"
	}

	if err := os.MkdirAll(path, 0o755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(filepath.Join(path, "assets"), 0o755); err != nil {
		panic(err)
	}

	s := os.DirFS(path)

	return &localStorer{s, path}
}

func (s *localStorer) Exists(ctx context.Context, path string) (bool, error) {
	_, err := fs.Stat(s.s, path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fault.Wrap(err, fctx.With(ctx))
	}

	return true, nil
}

func (s *localStorer) Read(ctx context.Context, path string) (io.Reader, int64, error) {
	f, err := s.s.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, 0, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return nil, 0, fault.Wrap(err, fctx.With(ctx))
	}

	info, err := f.Stat()
	if err != nil {
		return nil, 0, fault.Wrap(err, fctx.With(ctx))
	}

	return f, info.Size(), nil
}

func (s *localStorer) Write(ctx context.Context, path string, r io.Reader, size int64) error {
	fullpath := filepath.Join(s.path, path)

	dir := filepath.Dir(fullpath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	f, err := os.OpenFile(fullpath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0o755,
	)
	if err != nil {
		if os.IsNotExist(err) {
			return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return fault.Wrap(err, fctx.With(ctx))
	}

	_, err = io.Copy(f, r)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := f.Close(); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (s *localStorer) Delete(ctx context.Context, path string) error {
	fullpath := filepath.Join(s.path, path)

	if err := os.Remove(fullpath); err != nil {
		if os.IsNotExist(err) {
			return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (s *localStorer) List(ctx context.Context, prefix string) ([]string, error) {
	entries, err := fs.ReadDir(s.s, prefix)
	if err != nil {
		// NOTE: The Storer interface is an abstraction over both S3-style and
		// local filesystems. In the case of S3, listing a non-existent prefix
		// returns an empty list, whereas in the case of a local filesystem, it
		// returns an error. To maintain a consistent interface, we don't error.
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}

	return names, nil
}
