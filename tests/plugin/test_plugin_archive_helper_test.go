package plugin_test

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
)

var archiveBuildLocks sync.Map

// packageTestPluginArchive builds a plugin fixture and packages it as a zip
// archive containing a manifest and executable, returning the absolute path.
func packageTestPluginArchive(t *testing.T, fixtureDir, pluginName string) string {
	t.Helper()

	absFixture, err := filepath.Abs(fixtureDir)
	if err != nil {
		t.Fatalf("failed to resolve fixture path %q: %v", fixtureDir, err)
	}

	if _, err := os.Stat(filepath.Join(absFixture, "manifest.json")); err != nil {
		t.Fatalf("fixture %q missing manifest.json: %v", absFixture, err)
	}

	archivePath := cachedArchivePath(t, absFixture, pluginName)
	lock := archiveLock(archivePath)
	lock.Lock()
	defer lock.Unlock()

	if info, err := os.Stat(archivePath); err == nil && info.Size() > 0 {
		return archivePath
	}

	tempDir := t.TempDir()
	binaryPath := filepath.Join(tempDir, pluginName)

	cmd := exec.Command("go", "build", "-trimpath", "-ldflags", "-s -w", "-o", binaryPath, ".")
	cmd.Dir = absFixture
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")

	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build test plugin %q in %q: %v\n%s", pluginName, absFixture, err, string(output))
	}

	if _, err := os.Stat(binaryPath); err != nil {
		binaryPathWithExe := binaryPath + ".exe"
		if _, altErr := os.Stat(binaryPathWithExe); altErr != nil {
			t.Fatalf("built plugin binary missing (%q / %q): %v", binaryPath, binaryPathWithExe, err)
		}
		binaryPath = binaryPathWithExe
	}

	tmpArchivePath := filepath.Join(tempDir, pluginName+".zip")
	archiveFile, err := os.Create(tmpArchivePath)
	if err != nil {
		t.Fatalf("failed to create plugin archive %q: %v", tmpArchivePath, err)
	}

	zipWriter := zip.NewWriter(archiveFile)

	addFileToZip(t, zipWriter, filepath.Join(absFixture, "manifest.json"), "manifest.json", 0o644)
	addFileToZip(t, zipWriter, binaryPath, pluginName, 0o755)

	if err := zipWriter.Close(); err != nil {
		t.Fatalf("failed to finalize plugin archive %q: %v", tmpArchivePath, err)
	}
	if err := archiveFile.Close(); err != nil {
		t.Fatalf("failed to close plugin archive %q: %v", tmpArchivePath, err)
	}

	if err := os.MkdirAll(filepath.Dir(archivePath), 0o755); err != nil {
		t.Fatalf("failed to create plugin archive cache directory: %v", err)
	}
	if err := os.Rename(tmpArchivePath, archivePath); err != nil {
		t.Fatalf("failed to move plugin archive into cache (%q -> %q): %v", tmpArchivePath, archivePath, err)
	}

	return archivePath
}

func addFileToZip(t *testing.T, zw *zip.Writer, srcPath, destName string, mode fs.FileMode) {
	t.Helper()

	sourceFile, err := os.Open(srcPath)
	if err != nil {
		t.Fatalf("failed to open source file %q: %v", srcPath, err)
	}
	defer sourceFile.Close()

	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		t.Fatalf("failed to stat source file %q: %v", srcPath, err)
	}

	header, err := zip.FileInfoHeader(sourceInfo)
	if err != nil {
		t.Fatalf("failed to create zip header for %q: %v", srcPath, err)
	}
	header.Name = destName
	header.Method = zip.Deflate
	header.SetMode(mode)

	writer, err := zw.CreateHeader(header)
	if err != nil {
		t.Fatalf("failed to create zip entry %q: %v", destName, err)
	}

	if _, err := io.Copy(writer, sourceFile); err != nil {
		t.Fatalf("failed to write zip entry %q: %v", destName, err)
	}
}

func cachedArchivePath(t *testing.T, fixturePath, pluginName string) string {
	t.Helper()

	h := sha1.Sum([]byte(fixturePath))
	hashPrefix := hex.EncodeToString(h[:])[:10]
	filename := pluginName + "-" + hashPrefix + ".zip"

	return filepath.Join(
		os.TempDir(),
		"storyden-plugin-test-archives",
		runtime.GOOS+"-"+runtime.GOARCH,
		filename,
	)
}

func archiveLock(path string) *sync.Mutex {
	actual, _ := archiveBuildLocks.LoadOrStore(path, &sync.Mutex{})
	lock, ok := actual.(*sync.Mutex)
	if !ok {
		return &sync.Mutex{}
	}
	return lock
}
