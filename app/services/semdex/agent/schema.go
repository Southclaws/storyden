package agent

import (
	"encoding/json"
	"reflect"

	"github.com/google/jsonschema-go/jsonschema"
)

func convertSchema(schema any, raw json.RawMessage) (*jsonschema.Schema, error) {
	if len(raw) > 0 {
		return decodeSchema(raw)
	}

	// Use reflection to detect zero values for typed schemas (primarily outputs).
	if schema == nil || reflect.ValueOf(schema).IsZero() {
		return nil, nil
	}

	data, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}
	return decodeSchema(data)
}

func decodeSchema(raw json.RawMessage) (*jsonschema.Schema, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	var schema jsonschema.Schema
	if err := json.Unmarshal(raw, &schema); err != nil {
		return nil, err
	}
	return &schema, nil
}
