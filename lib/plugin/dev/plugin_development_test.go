package dev

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func TestExternalRPCURL(t *testing.T) {
	got, err := ExternalRPCURL("https://example.com", "plugin_ext_secret")
	if err != nil {
		t.Fatal(err)
	}
	want := "wss://example.com/rpc?token=plugin_ext_secret"
	if got != want {
		t.Fatalf("ExternalRPCURL() = %q, want %q", got, want)
	}
}

func TestExternalPluginFromAPI(t *testing.T) {
	var connection openapi.PluginModeUnion
	require.NoError(t, connection.FromPluginExternalProps(openapi.PluginExternalProps{
		Mode:  openapi.External,
		Token: "plugin_ext_secret",
	}))

	plugin, err := ExternalPluginFromAPI(openapi.Plugin{
		Id:         "plugin_123",
		Connection: connection,
	})

	require.NoError(t, err)
	require.Equal(t, "plugin_123", plugin.ID)
	require.Equal(t, "plugin_ext_secret", plugin.Token)
}

func TestCommandFromManifestOverride(t *testing.T) {
	command, args, err := CommandFromManifest(rpc.Manifest{Command: "go", Args: []string{"run", "."}}, []string{"node", "index.js"})
	if err != nil {
		t.Fatal(err)
	}
	if command != "node" || len(args) != 1 || args[0] != "index.js" {
		t.Fatalf("CommandFromManifest override = %q %#v", command, args)
	}
}

func TestWriteNewManifest(t *testing.T) {
	dir := t.TempDir()
	manifest := rpc.Manifest{
		ID:          "example-plugin",
		Name:        "Example Plugin",
		Author:      "tester",
		Description: "An example plugin.",
		Version:     "0.1.0",
		Command:     "./example-plugin",
	}

	if err := WriteNewManifest(os.Stdout, dir, manifest, false); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(dir, ManifestFilename)); err != nil {
		t.Fatal(err)
	}

	if _, err := ReadManifest(filepath.Join(dir, ManifestFilename)); err != nil {
		t.Fatal(err)
	}
}

func TestReadManifestWithAccessBlock(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ManifestFilename)
	err := os.WriteFile(path, []byte(`id: reactbot
name: React Bot
author: Southclaws
description: React to every new thread reply with a fire emoji.
version: 1.0.0
command: "./reactbot"
events_consumed:
  - EventThreadReplyCreated
access:
  handle: reactbot
  name: React Bot
  permissions:
    - CREATE_REACTION
`), 0o644)
	require.NoError(t, err)

	manifest, err := ReadManifest(path)
	require.NoError(t, err)
	require.Equal(t, "reactbot", manifest.Manifest.ID)

	access, ok := manifest.Manifest.Access.Get()
	require.True(t, ok)
	require.Equal(t, "reactbot", access.Handle)
	require.Equal(t, []string{"CREATE_REACTION"}, access.Permissions)
}

func TestBuildPackageCreatesValidatedArchive(t *testing.T) {
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, ManifestFilename), []byte(`id: example-plugin
name: Example Plugin
author: tester
description: An example plugin.
version: 0.1.0
command: "./example-plugin"
`), 0o644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(dir, "example-plugin"), []byte("#!/bin/sh\n"), 0o755)
	require.NoError(t, err)

	pkg, err := BuildPackage(context.Background(), dir, ManifestFilename)
	require.NoError(t, err)
	require.Equal(t, "example-plugin", pkg.Manifest.ID)

	names := archiveNames(t, pkg.Bytes)
	require.Contains(t, names, "manifest.json")
	require.Contains(t, names, "example-plugin")
	require.NotContains(t, names, ManifestFilename)

	rc, err := names["manifest.json"].Open()
	require.NoError(t, err)
	defer rc.Close()

	var manifest map[string]any
	require.NoError(t, json.NewDecoder(rc).Decode(&manifest))
	require.Equal(t, "example-plugin", manifest["id"])
}

func TestBuildPackageSkipsSymlinks(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation requires elevated privileges on Windows")
	}

	dir := t.TempDir()
	writeExampleManifest(t, dir)
	outside := filepath.Join(t.TempDir(), "outside.txt")
	require.NoError(t, os.WriteFile(outside, []byte("secret"), 0o644))
	require.NoError(t, os.Symlink(outside, filepath.Join(dir, "linked-outside.txt")))

	pkg, err := BuildPackage(context.Background(), dir, ManifestFilename)
	require.NoError(t, err)

	names := archiveNames(t, pkg.Bytes)
	require.NotContains(t, names, "linked-outside.txt")
}

func TestBuildPackageExcludesRelativePathsAndDirectories(t *testing.T) {
	dir := t.TempDir()
	writeExampleManifest(t, dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "example-plugin.zip"), []byte("stale zip"), 0o644))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "dist"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "dist", "bundle.txt"), []byte("bundle"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "kept.txt"), []byte("kept"), 0o644))

	pkg, err := BuildPackage(context.Background(), dir, ManifestFilename, "example-plugin.zip", "dist")
	require.NoError(t, err)

	names := archiveNames(t, pkg.Bytes)
	require.Contains(t, names, "kept.txt")
	require.NotContains(t, names, "example-plugin.zip")
	require.NotContains(t, names, "dist/bundle.txt")
}

func writeExampleManifest(t *testing.T, dir string) {
	t.Helper()
	err := os.WriteFile(filepath.Join(dir, ManifestFilename), []byte(`id: example-plugin
name: Example Plugin
author: tester
description: An example plugin.
version: 0.1.0
command: "./example-plugin"
`), 0o644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(dir, "example-plugin"), []byte("#!/bin/sh\n"), 0o755)
	require.NoError(t, err)
}

func archiveNames(t *testing.T, data []byte) map[string]*zip.File {
	t.Helper()

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	require.NoError(t, err)

	names := map[string]*zip.File{}
	for _, file := range zr.File {
		names[file.Name] = file
	}
	return names
}
