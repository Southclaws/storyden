package get

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
)

// TestYAMLMatchesJSONShape locks in the format-parity fix: --format yaml
// describes the exact same fields as --format json (including Fid), rather
// than the hand-shaped subset the command used to emit.
func TestYAMLMatchesJSONShape(t *testing.T) {
	r := require.New(t)

	schema := []openapi.PropertySchema{
		{
			Fid:  "field_1",
			Name: "status: #1",
			Type: openapi.PropertyTypeText,
			Sort: "asc",
		},
	}

	var out bytes.Buffer

	r.NoError(output.YAML(&out, schema))

	var decoded []map[string]string
	r.NoError(yaml.Unmarshal(out.Bytes(), &decoded))
	r.Len(decoded, 1)
	r.Equal("field_1", decoded[0]["fid"])
	r.Equal("status: #1", decoded[0]["name"])
}
