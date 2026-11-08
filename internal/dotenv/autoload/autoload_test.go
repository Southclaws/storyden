package autoload

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestAutoload(t *testing.T) {
	const child = "STORYDEN_DOTENV_AUTOLOAD_CHILD"
	const key = "STORYDEN_DOTENV_AUTOLOAD_VALUE"
	if os.Getenv(child) == "1" {
		got := os.Getenv(key)
		want := os.Getenv("STORYDEN_DOTENV_AUTOLOAD_EXPECTED")
		if got != want {
			t.Fatalf("autoload value = %q, want %q", got, want)
		}

		return
	}

	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name     string
		file     bool
		existing string
		want     string
	}{
		{name: "loads from working directory", file: true, want: "from file"},
		{name: "preserves environment", file: true, existing: "from environment", want: "from environment"},
		{name: "missing file is optional"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			if tc.file {
				path := filepath.Join(dir, ".env")
				if err := os.WriteFile(path, []byte(key+"=from file\n"), 0o600); err != nil {
					t.Fatal(err)
				}
			}

			command := exec.Command(executable, "-test.run=^TestAutoload$")
			command.Dir = dir
			for _, entry := range os.Environ() {
				name, _, _ := strings.Cut(entry, "=")
				if strings.EqualFold(name, key) {
					continue
				}

				command.Env = append(command.Env, entry)
			}

			command.Env = append(command.Env, child+"=1", "STORYDEN_DOTENV_AUTOLOAD_EXPECTED="+tc.want)
			if tc.existing != "" {
				command.Env = append(command.Env, key+"="+tc.existing)
			}

			if output, err := command.CombinedOutput(); err != nil {
				t.Fatalf("autoload subprocess: %v\n%s", err, output)
			}
		})
	}
}
