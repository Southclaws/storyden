package dev

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"

	pluginresource "github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

const ManifestFilename = "manifest.yaml"

type ManifestFile struct {
	Path     string
	Manifest rpc.Manifest
}

type PackageArchive struct {
	Manifest rpc.Manifest
	Bytes    []byte
	Files    []string
}

type excludePath struct {
	path string
	dir  bool
}

func ReadManifest(path string) (*ManifestFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse manifest: %w", err)
	}
	manifest, err := rpc.ManifestFromMap(raw)
	if err != nil {
		return nil, fmt.Errorf("validate manifest: %w", err)
	}

	return &ManifestFile{Path: path, Manifest: *manifest}, nil
}

func ReadProjectManifest(dir string, manifestPath string) (*ManifestFile, error) {
	path := manifestPath
	if path == "" {
		path = ManifestFilename
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(dir, path)
	}
	return ReadManifest(path)
}

func WriteNewManifest(out io.Writer, dir string, manifest rpc.Manifest, force bool) error {
	if err := manifest.Validate(); err != nil {
		return fmt.Errorf("validate manifest: %w", err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	path := filepath.Join(dir, ManifestFilename)
	if !force {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("%s already exists; use --force to overwrite", path)
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	data, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return err
	}

	fmt.Fprintf(out, "Created plugin manifest at %s\n", path)
	return nil
}

func BuildPackage(ctx context.Context, dir string, manifestPath string, excludePaths ...string) (*PackageArchive, error) {
	root, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	mf, err := ReadProjectManifest(root, manifestPath)
	if err != nil {
		return nil, err
	}

	manifestJSON, err := json.MarshalIndent(mf.Manifest, "", "  ")
	if err != nil {
		return nil, err
	}
	manifestJSON = append(manifestJSON, '\n')

	excludes := []excludePath{}
	for _, path := range excludePaths {
		if path == "" {
			continue
		}
		if !filepath.IsAbs(path) {
			path = filepath.Join(root, path)
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		abs = filepath.Clean(abs)
		info, err := os.Stat(abs)
		excludes = append(excludes, excludePath{
			path: abs,
			dir:  err == nil && info.IsDir(),
		})
	}

	type archiveFile struct {
		abs string
		rel string
	}
	files := []archiveFile{}
	if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}

		abs := filepath.Clean(path)
		if excluded := excludedArchivePath(abs, excludes); excluded {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}

		name := d.Name()
		if d.IsDir() {
			switch name {
			case ".git", "node_modules", ".next":
				return filepath.SkipDir
			}
			return nil
		}
		if name == ".DS_Store" || name == pluginresource.ArchiveManifestFileName {
			return nil
		}
		if sameFilePath(abs, mf.Path) {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files = append(files, archiveFile{abs: abs, rel: filepath.ToSlash(rel)})
		return nil
	}); err != nil {
		return nil, err
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].rel < files[j].rel
	})

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	manifestHeader := &zip.FileHeader{
		Name:   pluginresource.ArchiveManifestFileName,
		Method: zip.Deflate,
	}
	manifestHeader.SetMode(0o644)
	manifestWriter, err := zw.CreateHeader(manifestHeader)
	if err != nil {
		return nil, err
	}
	if _, err := manifestWriter.Write(manifestJSON); err != nil {
		return nil, err
	}

	written := []string{pluginresource.ArchiveManifestFileName}
	for _, file := range files {
		linkInfo, err := os.Lstat(file.abs)
		if err != nil {
			return nil, err
		}
		if linkInfo.Mode()&os.ModeSymlink != 0 {
			continue
		}
		info, err := os.Stat(file.abs)
		if err != nil {
			return nil, err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return nil, err
		}
		header.Name = file.rel
		header.Method = zip.Deflate

		writer, err := zw.CreateHeader(header)
		if err != nil {
			return nil, err
		}
		reader, err := os.Open(file.abs)
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(writer, reader); err != nil {
			_ = reader.Close()
			return nil, err
		}
		if err := reader.Close(); err != nil {
			return nil, err
		}
		written = append(written, file.rel)
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}

	archive := pluginresource.Binary(buf.Bytes())
	validated, err := archive.Validate(ctx)
	if err != nil {
		return nil, fmt.Errorf("validate package: %w", err)
	}
	if err := validated.Metadata.Validate(); err != nil {
		return nil, fmt.Errorf("validate package manifest: %w", err)
	}

	return &PackageArchive{
		Manifest: mf.Manifest,
		Bytes:    buf.Bytes(),
		Files:    written,
	}, nil
}

func excludedArchivePath(path string, excludes []excludePath) bool {
	for _, exclude := range excludes {
		if path == exclude.path {
			return true
		}
		if exclude.dir {
			prefix := exclude.path + string(filepath.Separator)
			if strings.HasPrefix(path, prefix) {
				return true
			}
		}
	}
	return false
}

func WritePackageFile(path string, pkg *PackageArchive, force bool) error {
	if !force {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("%s already exists; use --force to overwrite", path)
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, pkg.Bytes, 0o644)
}

func DefaultPackagePath(dir string, manifest rpc.Manifest) string {
	return filepath.Join(dir, Slugify(manifest.ID)+".zip")
}

func sameFilePath(a string, b string) bool {
	aa, err := filepath.Abs(a)
	if err != nil {
		return filepath.Clean(a) == filepath.Clean(b)
	}
	bb, err := filepath.Abs(b)
	if err != nil {
		return filepath.Clean(a) == filepath.Clean(b)
	}
	return filepath.Clean(aa) == filepath.Clean(bb)
}

func DefaultAuthor() string {
	if user := os.Getenv("USER"); user != "" {
		return user
	}
	if user := os.Getenv("USERNAME"); user != "" {
		return user
	}
	return "you"
}

func Slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var b strings.Builder
	lastDash := false
	for _, r := range value {
		valid := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
		if valid {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash && b.Len() > 0 {
			b.WriteRune('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func Titleize(value string) string {
	value = strings.ReplaceAll(value, "-", " ")
	value = strings.ReplaceAll(value, "_", " ")
	parts := strings.Fields(value)
	for i, part := range parts {
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	if len(parts) == 0 {
		return "My Plugin"
	}
	return strings.Join(parts, " ")
}
