package dev

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	pluginresource "github.com/Southclaws/storyden/app/resources/plugin"
)

type ExtractResult struct {
	Files []string
}

func PackageFilename(ctx context.Context, data []byte, fallbackID string) string {
	validated, err := pluginresource.Binary(data).Validate(ctx)
	if err == nil && strings.TrimSpace(validated.Metadata.ID) != "" {
		return Slugify(validated.Metadata.ID) + ".zip"
	}
	if strings.TrimSpace(fallbackID) != "" {
		return Slugify(fallbackID) + ".zip"
	}
	return "plugin.zip"
}

func ExtractPackageArchive(data []byte, dir string, force bool) (*ExtractResult, error) {
	if dir == "" {
		dir = "."
	}

	destination, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(destination, 0o755); err != nil {
		return nil, err
	}

	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("open package archive: %w", err)
	}

	written := []string{}
	for _, file := range reader.File {
		target, rel, err := extractTargetPath(destination, file.Name)
		if err != nil {
			return nil, err
		}

		mode := file.FileInfo().Mode()
		if mode&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("refusing to extract symlink from package: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return nil, err
			}
			continue
		}

		if !force {
			if _, err := os.Stat(target); err == nil {
				return nil, fmt.Errorf("output file already exists: %s (use --force to overwrite)", target)
			} else if !errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("inspect output file: %w", err)
			}
		}

		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return nil, err
		}

		rc, err := file.Open()
		if err != nil {
			return nil, err
		}
		if err := writeExtractedFile(target, rc, mode.Perm()); err != nil {
			_ = rc.Close()
			return nil, err
		}
		if err := rc.Close(); err != nil {
			return nil, err
		}

		written = append(written, rel)
	}

	return &ExtractResult{Files: written}, nil
}

func extractTargetPath(destination string, name string) (string, string, error) {
	if name == "" || filepath.IsAbs(name) || strings.Contains(name, `\`) {
		return "", "", fmt.Errorf("invalid package path: %s", name)
	}

	clean := filepath.Clean(filepath.FromSlash(name))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", "", fmt.Errorf("invalid package path: %s", name)
	}

	target := filepath.Join(destination, clean)
	rel, err := filepath.Rel(destination, target)
	if err != nil {
		return "", "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return "", "", fmt.Errorf("invalid package path: %s", name)
	}

	return target, filepath.ToSlash(rel), nil
}

func writeExtractedFile(target string, reader io.Reader, mode os.FileMode) error {
	if mode == 0 {
		mode = 0o644
	}

	file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(file, reader); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}
