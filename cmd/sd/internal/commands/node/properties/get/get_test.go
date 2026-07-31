package get

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
)

// TestYAMLMatchesJSONShape locks in the format-parity fix: --format yaml
// describes the exact same fields as --format json (including Fid), rather
// than the hand-shaped subset the command used to emit.
func TestYAMLMatchesJSONShape(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	properties := []openapi.Property{{
		Fid:   "prop_1",
		Name:  "property: with # characters",
		Type:  openapi.PropertyTypeText,
		Value: "value: with # characters",
	}}
	var out bytes.Buffer

	r.NoError(output.YAML(&out, properties))

	var decoded []map[string]string
	r.NoError(yaml.Unmarshal(out.Bytes(), &decoded))
	r.Len(decoded, 1)
	a.Equal("prop_1", decoded[0]["fid"])
	a.Equal("property: with # characters", decoded[0]["name"])
	a.Equal("value: with # characters", decoded[0]["value"])
}
