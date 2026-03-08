package plugin

import (
	"archive/zip"
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBinaryValidateRejectsOversizedManifest(t *testing.T) {
	t.Parallel()

	oversized := strings.Repeat("a", int(MaxManifestSizeBytes)+32)
	archive := makeArchive(t, `{"id":"plug1","author":"author1","name":"Plugin","command":"run","description":"`+oversized+`","version":"1.0.0"}`)

	_, err := Binary(archive).Validate(context.Background())
	require.Error(t, err)
	require.ErrorContains(t, err, "manifest exceeds maximum size")
}

func TestBinaryValidateAcceptsSmallManifest(t *testing.T) {
	t.Parallel()

	archive := makeArchive(t, `{"id":"plug1","author":"author1","name":"Plugin","command":"run","description":"ok","version":"1.0.0"}`)

	_, err := Binary(archive).Validate(context.Background())
	require.NoError(t, err)
}

func makeArchive(t *testing.T, manifest string) []byte {
	t.Helper()

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create(ArchiveManifestFileName)
	require.NoError(t, err)
	_, err = w.Write([]byte(manifest))
	require.NoError(t, err)
	require.NoError(t, zw.Close())

	return buf.Bytes()
}
