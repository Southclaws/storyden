package get

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func TestRenderYAMLUsesStructuredEncoder(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	properties := []openapi.Property{{
		Name:  "property: with # characters",
		Type:  openapi.PropertyTypeText,
		Value: "value: with # characters",
	}}
	var out bytes.Buffer

	r.NoError(renderYAML(&out, properties))

	var decoded map[string][]map[string]string
	r.NoError(yaml.Unmarshal(out.Bytes(), &decoded))
	r.Len(decoded["properties"], 1)
	a.Equal("property: with # characters", decoded["properties"][0]["name"])
	a.Equal("value: with # characters", decoded["properties"][0]["value"])
}
