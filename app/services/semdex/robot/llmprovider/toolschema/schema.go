package toolschema

import (
	"encoding/json"
	"strings"

	"google.golang.org/genai"
)

// FromFunctionDeclaration projects either of genai's mutually exclusive
// function parameter schema representations into standard JSON Schema.
func FromFunctionDeclaration(fn *genai.FunctionDeclaration) map[string]any {
	var schema map[string]any

	if fn != nil && fn.Parameters != nil {
		schema = schemaMap(fn.Parameters)
	}
	if len(schema) == 0 && fn != nil && fn.ParametersJsonSchema != nil {
		schema = schemaMap(fn.ParametersJsonSchema)
	}
	if len(schema) == 0 {
		schema = map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		}
	}

	normalizeTypeNames(schema)
	return schema
}

func schemaMap(value any) map[string]any {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}

	var schema map[string]any
	if err := json.Unmarshal(raw, &schema); err != nil {
		return nil
	}
	return schema
}

// normalizeTypeNames converts genai's uppercase OpenAPI type enum values into
// the lowercase names required by JSON Schema consumers such as OpenAI and
// Anthropic.
func normalizeTypeNames(schema any) {
	switch value := schema.(type) {
	case map[string]any:
		if schemaType, ok := value["type"].(string); ok {
			value["type"] = strings.ToLower(schemaType)
		}
		if schemaTypes, ok := value["type"].([]any); ok {
			for i, schemaType := range schemaTypes {
				if name, ok := schemaType.(string); ok {
					schemaTypes[i] = strings.ToLower(name)
				}
			}
		}
		for _, child := range value {
			normalizeTypeNames(child)
		}
	case []any:
		for _, child := range value {
			normalizeTypeNames(child)
		}
	}
}
