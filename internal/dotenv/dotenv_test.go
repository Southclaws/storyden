package dotenv

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"
)

func TestParse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input string
		want  map[string]string
	}{
		{
			name:  "empty file",
			input: "",
			want:  map[string]string{},
		},
		{
			name:  "comments and whitespace",
			input: "\n  # comment\n\t\n NAME = Storyden \nEMPTY=\n",
			want:  map[string]string{"NAME": "Storyden", "EMPTY": ""},
		},
		{
			name:  "export prefix",
			input: "export NAME=value\nexport\t_OTHER2 = second\nexported=third",
			want:  map[string]string{"NAME": "value", "_OTHER2": "second", "exported": "third"},
		},
		{
			name:  "windows line endings and UTF-8 BOM",
			input: "\ufeffNAME=Storyden\r\nOTHER=value\r\n",
			want:  map[string]string{"NAME": "Storyden", "OTHER": "value"},
		},
		{
			name:  "last duplicate wins",
			input: "NAME=first\nNAME=last",
			want:  map[string]string{"NAME": "last"},
		},
		{
			name:  "values containing equals and unicode",
			input: "TOKEN=abc==\nNAME=こんにちは 🌍",
			want:  map[string]string{"TOKEN": "abc==", "NAME": "こんにちは 🌍"},
		},
		{
			name:  "quoted values and trailing comments",
			input: "SINGLE=' spaced # value ' # comment\nDOUBLE=\"Storyden\" # comment\nEMPTY=''",
			want:  map[string]string{"SINGLE": " spaced # value ", "DOUBLE": "Storyden", "EMPTY": ""},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parse(strings.NewReader(tc.input))
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("parse() = %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestParseValue(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input string
		want  string
	}{
		{name: "empty", input: "", want: ""},
		{name: "unquoted", input: "Storyden community", want: "Storyden community"},
		{name: "comment only", input: "# comment", want: ""},
		{name: "inline comment", input: "value # comment", want: "value"},
		{name: "tab before comment", input: "value\t# comment", want: "value"},
		{name: "hash in value", input: "https://example.com/#section", want: "https://example.com/#section"},
		{name: "single quoted literal", input: `'\n $HOME "quoted"'`, want: `\n $HOME "quoted"`},
		{name: "empty double quotes", input: `""`, want: ""},
		{name: "opposite quotes", input: `"it's Storyden"`, want: "it's Storyden"},
		{name: "double quoted escapes", input: `"one\ntwo\r\t\"quoted\"\\end"`, want: "one\ntwo\r\t\"quoted\"\\end"},
		{name: "escaped escape", input: `"\\n"`, want: `\n`},
		{name: "unknown escape preserved", input: `"\q"`, want: `\q`},
		{name: "unquoted windows path", input: `C:\Users\Storyden`, want: `C:\Users\Storyden`},
		{name: "quoted windows path", input: `"C:\Users\Storyden"`, want: `C:\Users\Storyden`},
		{name: "unquoted variables stay literal", input: "$HOME/${NAME}", want: "$HOME/${NAME}"},
		{name: "quoted variables stay literal", input: `"$HOME/${NAME}"`, want: "$HOME/${NAME}"},
		{name: "shell expression stays literal", input: "$(echo Storyden)", want: "$(echo Storyden)"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseValue(tc.input)
			if err != nil {
				t.Fatal(err)
			}

			if got != tc.want {
				t.Fatalf("parseValue() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input string
	}{
		{name: "missing equals", input: "NAME"},
		{name: "empty key", input: "=secret"},
		{name: "digit starts key", input: "1NAME=secret"},
		{name: "space in key", input: "BAD NAME=secret"},
		{name: "nonportable key", input: "BAD-NAME=secret"},
		{name: "unicode key", input: "名=secret"},
		{name: "NUL in key", input: "NA\x00ME=secret"},
		{name: "NUL in value", input: "NAME=secret\x00"},
		{name: "unterminated single quote", input: "NAME='secret"},
		{name: "unterminated double quote", input: `NAME="secret`},
		{name: "escaped closing quote", input: `NAME="secret\"`},
		{name: "trailing text", input: `NAME="value" secret`},
		{name: "multiline quoted value", input: "NAME='secret\nvalue'"},
		{name: "oversized line", input: "NAME=secret" + strings.Repeat("x", 65536)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			values, err := parse(strings.NewReader("# comment\n" + tc.input))
			if err == nil || values != nil {
				t.Fatalf("parse() = %#v, %v; want nil values and an error", values, err)
			}

			if !strings.Contains(err.Error(), "line 2:") {
				t.Fatalf("error lacks line number: %v", err)
			}

			if strings.Contains(err.Error(), "secret") {
				t.Fatalf("error exposes the value: %v", err)
			}
		})
	}
}

func TestParseReadError(t *testing.T) {
	t.Parallel()

	want := errors.New("read failed")
	_, err := parse(iotest.ErrReader(want))
	if !errors.Is(err, want) {
		t.Fatalf("parse() error = %v, want %v", err, want)
	}
}

func TestLoad(t *testing.T) {
	const existing = "STORYDEN_DOTENV_TEST_EXISTING"
	const empty = "STORYDEN_DOTENV_TEST_EMPTY"
	const loaded = "STORYDEN_DOTENV_TEST_LOADED"

	t.Setenv(existing, "from environment")
	t.Setenv(empty, "")
	t.Setenv(loaded, "")
	if err := os.Unsetenv(loaded); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(t.TempDir(), "environment settings.env")
	contents := existing + "=from file\r\n" + empty + "=from file\r\n" + loaded + "=first\r\n" + loaded + "=last\r\n"
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}

	for range 2 {
		if err := Load(path); err != nil {
			t.Fatal(err)
		}

		for key, want := range map[string]string{existing: "from environment", empty: "", loaded: "last"} {
			got, exists := os.LookupEnv(key)
			if !exists || got != want {
				t.Errorf("environment %s = %q (exists %t), want %q", key, got, exists, want)
			}
		}
	}
}

func TestLoadInvalidFileDoesNotSetEnvironment(t *testing.T) {
	const key = "STORYDEN_DOTENV_TEST_INVALID"
	t.Setenv(key, "")
	if err := os.Unsetenv(key); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte(key+"=value\nINVALID"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := Load(path)
	if err == nil || !strings.Contains(err.Error(), path) {
		t.Fatalf("Load() error = %v, want an error containing the path", err)
	}

	if _, exists := os.LookupEnv(key); exists {
		t.Fatal("invalid file partially populated the environment")
	}
}

func TestLoadMissingFile(t *testing.T) {
	err := Load(filepath.Join(t.TempDir(), "missing.env"))
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Load() error = %v, want os.ErrNotExist", err)
	}
}
