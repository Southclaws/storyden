package agent_registry

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestClientToolContextUsesDurableJSONFieldNames(t *testing.T) {
	context, err := NewClientToolContext("browser-1", []ClientToolDefinition{{
		Name:        "library_page_layout_get",
		Description: "Read the current library layout.",
		InputSchema: map[string]any{"type": "object"},
		Annotations: ClientToolAnnotations{ReadOnlyHint: true},
	}})
	if err != nil {
		t.Fatalf("NewClientToolContext() error = %v", err)
	}

	encoded, err := json.Marshal(context)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	const expected = `{"client_id":"browser-1","tools":[{"name":"library_page_layout_get","description":"Read the current library layout.","inputSchema":{"type":"object"},"annotations":{"readOnlyHint":true}}]}`
	if string(encoded) != expected {
		t.Fatalf("JSON = %s, want %s", encoded, expected)
	}
}

func TestNewClientToolContextRejectsDuplicateNames(t *testing.T) {
	_, err := NewClientToolContext("browser-1", []ClientToolDefinition{
		{Name: "duplicate", Description: "First.", InputSchema: map[string]any{"type": "object"}},
		{Name: "duplicate", Description: "Second.", InputSchema: map[string]any{"type": "object"}},
	})
	if err == nil {
		t.Fatal("NewClientToolContext() error = nil, want duplicate-name error")
	}
}

func TestNewClientToolContextRejectsNonObjectInputSchema(t *testing.T) {
	_, err := NewClientToolContext("browser-1", []ClientToolDefinition{{
		Name:        "invalid_schema",
		Description: "Has invalid input.",
		InputSchema: map[string]any{"type": "string"},
	}})
	if err == nil {
		t.Fatal("NewClientToolContext() error = nil, want object-schema error")
	}
}

func TestNewClientToolContextRejectsOversizedDescription(t *testing.T) {
	_, err := NewClientToolContext("browser-1", []ClientToolDefinition{{
		Name:        "oversized_description",
		Description: strings.Repeat("x", maxClientToolTextSize+1),
		InputSchema: map[string]any{"type": "object"},
	}})
	if err == nil {
		t.Fatal("NewClientToolContext() error = nil, want description-size error")
	}

	const expected = `client tool "oversized_description" description exceeds 4096 bytes`
	if err.Error() != expected {
		t.Fatalf("NewClientToolContext() error = %q, want %q", err, expected)
	}
}
