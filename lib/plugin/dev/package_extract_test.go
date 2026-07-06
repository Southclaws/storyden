package dev

import (
	"archive/zip"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractPackageArchive(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	mustZipFile(t, zw, "manifest.json", []byte(`{"id":"example-plugin","name":"Example Plugin","author":"test","description":"test","version":"1.0.0","command":"./example"}`))
	mustZipFile(t, zw, "nested/file.txt", []byte("hello"))
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	result, err := ExtractPackageArchive(buf.Bytes(), dir, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Files) != 2 {
		t.Fatalf("expected 2 extracted files, got %d", len(result.Files))
	}

	data, err := os.ReadFile(filepath.Join(dir, "nested", "file.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatalf("expected extracted content, got %q", string(data))
	}

	if _, err := ExtractPackageArchive(buf.Bytes(), dir, false); err == nil {
		t.Fatal("expected existing file error")
	}
}

func TestExtractPackageArchiveRejectsTraversal(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	mustZipFile(t, zw, "../escape.txt", []byte("nope"))
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}

	if _, err := ExtractPackageArchive(buf.Bytes(), t.TempDir(), false); err == nil {
		t.Fatal("expected traversal error")
	}
}

func TestPackageFilenameUsesManifestID(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	mustZipFile(t, zw, "manifest.json", []byte(`{"id":"example-plugin","name":"Example Plugin","author":"test","description":"test","version":"1.0.0","command":"./example"}`))
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}

	if got := PackageFilename(context.Background(), buf.Bytes(), "fallback"); got != "example-plugin.zip" {
		t.Fatalf("expected filename from manifest id, got %q", got)
	}
}

func mustZipFile(t *testing.T, zw *zip.Writer, name string, data []byte) {
	t.Helper()

	w, err := zw.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write(data); err != nil {
		t.Fatal(err)
	}
}
