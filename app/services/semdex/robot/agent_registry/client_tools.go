package agent_registry

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
)

const (
	ClientToolSourceWebMCP = "webmcp"
	ClientToolMetadataKey  = "storyden_client_tool"

	maxClientToolCount      = 32
	maxClientToolTextSize   = 4 * 1024
	maxClientToolSchemaSize = 64 * 1024
)

var clientToolNamePattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_-]{0,63}$`)

type ClientToolAnnotations struct {
	ReadOnlyHint         bool `json:"readOnlyHint,omitempty"`
	UntrustedContentHint bool `json:"untrustedContentHint,omitempty"`
	ConsequentialHint    bool `json:"consequentialHint,omitempty"`
}

type ClientToolDefinition struct {
	Name        string                `json:"name"`
	Title       string                `json:"title,omitempty"`
	Description string                `json:"description"`
	InputSchema map[string]any        `json:"inputSchema"`
	Annotations ClientToolAnnotations `json:"annotations,omitempty"`
}

type ClientToolContext struct {
	ClientID string                 `json:"client_id"`
	Tools    []ClientToolDefinition `json:"tools"`
}

func NewClientToolContext(clientID string, definitions []ClientToolDefinition) (*ClientToolContext, error) {
	clientID = strings.TrimSpace(clientID)
	if clientID == "" {
		return nil, fmt.Errorf("client tool context requires a client ID")
	}
	if len(clientID) > 128 {
		return nil, fmt.Errorf("client tool context client ID must be at most 128 characters")
	}
	if len(definitions) == 0 {
		return nil, nil
	}
	if len(definitions) > maxClientToolCount {
		return nil, fmt.Errorf("client tool context contains %d tools; maximum is %d", len(definitions), maxClientToolCount)
	}

	validated := make([]ClientToolDefinition, 0, len(definitions))
	seen := make(map[string]struct{}, len(definitions))
	for _, definition := range definitions {
		definition.Name = strings.TrimSpace(definition.Name)
		definition.Title = strings.TrimSpace(definition.Title)
		definition.Description = strings.TrimSpace(definition.Description)
		if !clientToolNamePattern.MatchString(definition.Name) {
			return nil, fmt.Errorf("invalid client tool name %q", definition.Name)
		}
		if definition.Description == "" {
			return nil, fmt.Errorf("client tool %q requires a description", definition.Name)
		}
		if len(definition.Description) > maxClientToolTextSize {
			return nil, fmt.Errorf("client tool %q description exceeds %d bytes", definition.Name, maxClientToolTextSize)
		}
		if _, ok := seen[definition.Name]; ok {
			return nil, fmt.Errorf("duplicate client tool name %q", definition.Name)
		}
		seen[definition.Name] = struct{}{}

		if _, err := definition.JSONSchema(); err != nil {
			return nil, err
		}
		validated = append(validated, definition)
	}

	return &ClientToolContext{ClientID: clientID, Tools: validated}, nil
}

func (c *ClientToolContext) Find(name string) (ClientToolDefinition, bool) {
	if c == nil {
		return ClientToolDefinition{}, false
	}
	for _, definition := range c.Tools {
		if definition.Name == name {
			return definition, true
		}
	}
	return ClientToolDefinition{}, false
}

func (d ClientToolDefinition) JSONSchema() (*jsonschema.Schema, error) {
	if d.InputSchema == nil {
		return nil, fmt.Errorf("client tool %q requires an input schema", d.Name)
	}
	if schemaType, _ := d.InputSchema["type"].(string); schemaType != "object" {
		return nil, fmt.Errorf("client tool %q input schema must have type object", d.Name)
	}

	encoded, err := json.Marshal(d.InputSchema)
	if err != nil {
		return nil, fmt.Errorf("encode client tool %q input schema: %w", d.Name, err)
	}
	if len(encoded) > maxClientToolSchemaSize {
		return nil, fmt.Errorf("client tool %q input schema exceeds %d bytes", d.Name, maxClientToolSchemaSize)
	}

	var schema jsonschema.Schema
	if err := json.Unmarshal(encoded, &schema); err != nil {
		return nil, fmt.Errorf("decode client tool %q input schema: %w", d.Name, err)
	}
	if _, err := schema.Resolve(nil); err != nil {
		return nil, fmt.Errorf("resolve client tool %q input schema: %w", d.Name, err)
	}
	return &schema, nil
}
