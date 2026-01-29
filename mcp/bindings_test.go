package mcp

import (
	"testing"
)

func TestSchemaLoading(t *testing.T) {
	if Schema == nil {
		t.Fatal("Schema should not be nil")
	}

	if Schema.Definitions == nil {
		t.Fatal("Schema.Definitions should not be nil")
	}

	if len(Schema.Definitions) == 0 {
		t.Fatal("Schema.Definitions should not be empty")
	}

	toolSearch := GetSearchTool()
	if toolSearch == nil {
		t.Fatal("toolSearch should not be nil")
	}

	if toolSearch.InputSchema == nil {
		t.Fatal("toolSearch.InputSchema should not be nil")
	}

	if toolSearch.InputSchema.Properties == nil {
		t.Fatal("toolSearch.InputSchema.Properties should not be nil")
	}

	kindProp, ok := toolSearch.InputSchema.Properties["kind"]
	if !ok {
		t.Fatal("kind property should exist in toolSearch input schema")
	}

	if kindProp.Type != "array" {
		t.Errorf("kind should be an array, got %s", kindProp.Type)
	}

	if kindProp.Items == nil {
		t.Fatal("kind.items should not be nil")
	}

	if len(kindProp.Items.Enum) == 0 {
		t.Fatal("kind.items.enum should not be empty - this means $ref was not dereferenced")
	}

	t.Logf("Successfully loaded dereferenced schema with %d enum values", len(kindProp.Items.Enum))
}
