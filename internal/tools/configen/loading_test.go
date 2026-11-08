package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSchemaRejectsInvalidConfiguration(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		schema string
		error  string
	}{
		{"unknown type", "type: map[string]string", "unsupported configuration type"},
		{"invalid integer default", "type: int\n      default: nope", "invalid default"},
		{"overflowing integer default", "type: int32\n      default: '2147483648'", "invalid default"},
		{"invalid boolean default", "type: bool\n      default: yes", "invalid default"},
		{"invalid duration default", "type: time.Duration\n      default: soon", "invalid default"},
		{"invalid URL default", "type: net/url.URL\n      default: '%zz'", "invalid default"},
		{"invalid level default", "type: log/slog.Level\n      default: verbose", "invalid default"},
		{"nonfinite default", "type: float64\n      default: NaN", "must be finite"},
		{"non-numeric minimum", "type: string\n      minimum: 0", "bounds require a numeric type"},
		{"reversed bounds", "type: int\n      minimum: 2\n      maximum: 1", "minimum exceeds maximum"},
		{"fractional integer bound", "type: int\n      minimum: 0.5", "integer bounds must be whole numbers"},
		{"overflowing integer bound", "type: int32\n      maximum: 2147483648", "outside the range"},
		{"nonfinite bound", "type: float64\n      minimum: .nan", "bounds must be finite"},
		{"default below minimum", "type: int\n      default: '1'\n      minimum: 2", "below minimum"},
		{"default above maximum", "type: float64\n      default: '2'\n      maximum: 1", "exceeds maximum"},
		{"empty required default", "type: string\n      required: true\n      default: ''", "required value cannot be empty"},
		{"zero without default", "type: int\n      zero_is_default: true", "requires an exported default"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			directory := t.TempDir()
			schema := "- section: Test\n  fields:\n    - env: TEST_VALUE\n      name: Value\n      " + tc.schema + "\n"
			input := filepath.Join(directory, "config.yaml")
			require.NoError(t, os.WriteFile(input, []byte(schema), 0o600))
			err := run(input, filepath.Join(directory, "config.go"), filepath.Join(directory, "configuration.mdx"))
			require.ErrorContains(t, err, tc.error)
			require.NoFileExists(t, filepath.Join(directory, "config.go"))
			require.NoFileExists(t, filepath.Join(directory, "configuration.mdx"))
		})
	}
}

func TestSchemaRejectsDuplicateNames(t *testing.T) {
	t.Parallel()

	for _, fields := range [][]ConfigOption{
		{{Env: "ONE", Name: "Value", Type: "string"}, {Env: "TWO", Name: "Value", Type: "string"}},
		{{Env: "ONE", Name: "First", Type: "string"}, {Env: "ONE", Name: "Second", Type: "string"}},
		{{Env: "ONE", Name: "not-valid", Type: "string"}},
		{{Env: "", Name: "Value", Type: "string"}},
	} {
		require.Error(t, validateSchema([]Section{{Fields: fields}}))
	}
}

func TestGeneratedLoader(t *testing.T) {
	directory := t.TempDir()
	input := filepath.Join(directory, "config.yaml")
	output := filepath.Join(directory, "config.go")
	docs := filepath.Join(directory, "configuration.mdx")
	schema, err := os.ReadFile("testdata/config.yaml")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(input, schema, 0o600))
	require.NoError(t, run(input, output, docs))

	generated, err := os.ReadFile(output)
	require.NoError(t, err)
	require.NotContains(t, string(generated), "reflect")
	require.NotContains(t, string(generated), "envconfig")
	require.Contains(t, string(generated), "Token string")
	require.NotContains(t, string(generated), "TokenRequired")

	document, err := os.ReadFile(docs)
	require.NoError(t, err)
	require.Contains(t, string(document), "<td>required</td>")
	require.Contains(t, string(document), "<td>minimum</td>")
	require.Contains(t, string(document), "<td>maximum</td>")

	require.NoError(t, run(input, output, docs))
	repeated, err := os.ReadFile(output)
	require.NoError(t, err)
	require.Equal(t, generated, repeated)

	require.NoError(t, os.WriteFile(filepath.Join(directory, "go.mod"), []byte("module configfixture\n\ngo 1.24\n"), 0o600))
	testSource, err := os.ReadFile("testdata/config_test.go")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(directory, "config_test.go"), testSource, 0o600))
	command := exec.Command("go", "test", "-p=1", "./...")
	command.Dir = directory
	command.Env = append(os.Environ(), "GOWORK=off")
	result, err := command.CombinedOutput()
	require.NoError(t, err, strings.TrimSpace(string(result)))
}
