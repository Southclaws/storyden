package ai

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type SuggestTagsResultSchema struct {
	Tags []string `json:"tags" jsonschema:"title=Tags,description=List of suggested tags,items=string"`
}

func Test_schemaFromObjectInstance(t *testing.T) {
	r := require.New(t)

	sc, err := schemaFromObjectInstance(SuggestTagsResultSchema{})
	r.NoError(err)

	j, _ := json.MarshalIndent(sc, "", "  ")

	r.Equal(`{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "properties": {
    "tags": {
      "items": {
        "type": "string"
      },
      "type": "array",
      "title": "Tags",
      "description": "List of suggested tags"
    }
  },
  "additionalProperties": false,
  "type": "object",
  "required": [
    "tags"
  ]
}`, string(j))
}
