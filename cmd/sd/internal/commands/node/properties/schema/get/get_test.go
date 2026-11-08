package get

import (
	"bytes"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func TestRenderYAMLUsesStructuredEncoding(t *testing.T) {
	r := require.New(t)

	schema := []openapi.PropertySchema{
		{
			Name: "status: #1",
			Type: openapi.PropertyTypeText,
			Sort: "asc",
		},
	}

	var out bytes.Buffer

	r.NoError(renderYAML(&out, schema))

	var decoded struct {
		Schema []struct {
			Name string `yaml:"name"`
			Type string `yaml:"type"`
			Sort string `yaml:"sort"`
		} `yaml:"schema"`
	}
	r.NoError(yaml.Unmarshal(out.Bytes(), &decoded))
	r.Len(decoded.Schema, 1)
	r.Equal("status: #1", decoded.Schema[0].Name)
}
