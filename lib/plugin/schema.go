package plugin

import (
	_ "embed"
	"encoding/json"

	"github.com/google/jsonschema-go/jsonschema"
)

//go:embed plugin.json
var ManifestSchema []byte

func GetManifestSchema() *jsonschema.Schema {
	var v *jsonschema.Schema

	err := json.Unmarshal(ManifestSchema, &v)
	if err != nil {
		panic(err)
	}

	return v
}
