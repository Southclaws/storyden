package mcp

import (
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/google/jsonschema-go/jsonschema"
)

//go:embed robots.json
var schemaJSON []byte

// Schema represents the entire schema for storyden agent tools.
var Schema *jsonschema.Schema

type ToolAnnotations struct {
	ReadOnlyHint    bool
	DestructiveHint bool
	IdempotentHint  bool
	OpenWorldHint   bool
}

type ToolDefinition struct {
	Name                 string
	Title                string
	Description          string
	InputSchema          *jsonschema.Schema
	OutputSchema         *jsonschema.Schema
	RequiredPermission   opt.Optional[rbac.Permission]
	RequiresConfirmation bool
	Annotations          ToolAnnotations
}

func init() {
	Schema = &jsonschema.Schema{}
	err := json.Unmarshal(schemaJSON, Schema)
	if err != nil {
		panic("failed to unmarshal schema from JSON: " + err.Error())
	}

	initAllTools()
}

func deepCloneSchema(schema *jsonschema.Schema) *jsonschema.Schema {
	if schema == nil {
		return nil
	}

	b, err := json.Marshal(schema)
	if err != nil {
		panic("failed to clone schema: " + err.Error())
	}

	cloned := &jsonschema.Schema{}
	if err := json.Unmarshal(b, cloned); err != nil {
		panic("failed to clone schema: " + err.Error())
	}

	return cloned
}

func referencedDefinitionNames(schema *jsonschema.Schema) map[string]struct{} {
	names := map[string]struct{}{}
	if schema == nil {
		return names
	}

	b, err := json.Marshal(schema)
	if err != nil {
		panic("failed to read schema refs: " + err.Error())
	}

	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		panic("failed to read schema refs: " + err.Error())
	}

	var walk func(any)
	walk = func(node any) {
		switch n := node.(type) {
		case map[string]any:
			if ref, ok := n["$ref"].(string); ok && strings.HasPrefix(ref, "#/definitions/") {
				if name := strings.TrimPrefix(ref, "#/definitions/"); name != "" && !strings.Contains(name, "/") {
					names[name] = struct{}{}
				}
			}
			for _, child := range n {
				walk(child)
			}
		case []any:
			for _, child := range n {
				walk(child)
			}
		}
	}
	walk(v)

	return names
}

func cloneSchemaWithDefinitions(schema *jsonschema.Schema) *jsonschema.Schema {
	cloned := deepCloneSchema(schema)
	if cloned == nil {
		return nil
	}

	definitions := map[string]*jsonschema.Schema{}
	visited := map[string]struct{}{}
	var addDefinition func(string)
	addDefinition = func(name string) {
		if _, ok := visited[name]; ok {
			return
		}
		visited[name] = struct{}{}

		def, ok := Schema.Definitions[name]
		if !ok {
			return
		}
		definitions[name] = deepCloneSchema(def)

		for nested := range referencedDefinitionNames(def) {
			addDefinition(nested)
		}
	}

	for name := range referencedDefinitionNames(cloned) {
		addDefinition(name)
	}

	if len(definitions) > 0 {
		cloned.Definitions = definitions
	}

	return cloned
}

// initTool initializes a tool definition from the schema following the naming convention:
// - ToolDefinition: "Tool{name}" (e.g., "ToolRobotSwitch")
// - InputSchema: "Tool{name}Input" (e.g., "ToolRobotSwitchInput")
// - OutputSchema: "Tool{name}Output" (e.g., "ToolRobotSwitchOutput")
func initTool(name string) *ToolDefinition {
	toolDefName := "Tool" + name

	toolDef, ok := Schema.Definitions[toolDefName]
	if !ok {
		panic(toolDefName + " not found in schema definitions")
	}

	toolName := toolDef.Title
	description := toolDef.Description

	var inputSchema *jsonschema.Schema
	if inputProp, ok := toolDef.Properties["input"]; ok {
		if inputProp.Ref != "" {

			parts := strings.Split(inputProp.Ref, "/")
			if len(parts) > 0 {
				defName := parts[len(parts)-1]
				if def, ok := Schema.Definitions[defName]; ok {
					inputSchema = cloneSchemaWithDefinitions(def)
				}
			}
		} else if inputProp.Type != "" {
			inputSchema = cloneSchemaWithDefinitions(inputProp)
		}
	}

	var outputSchema *jsonschema.Schema
	if outputProp, ok := toolDef.Properties["output"]; ok {
		if outputProp.Ref != "" {

			parts := strings.Split(outputProp.Ref, "/")
			if len(parts) > 0 {
				defName := parts[len(parts)-1]
				if def, ok := Schema.Definitions[defName]; ok {
					outputSchema = cloneSchemaWithDefinitions(def)
				}
			}
		} else if outputProp.Type != "" {
			outputSchema = cloneSchemaWithDefinitions(outputProp)
		}
	}

	var requiredPermission opt.Optional[rbac.Permission]
	var displayTitle string
	var requiresConfirmation bool
	var annotations ToolAnnotations

	if toolDef.Extra != nil {
		if ext, ok := toolDef.Extra["x-storyden-tool"].(map[string]any); ok {
			if roleStr, ok := ext["role"].(string); ok {
				r, err := rbac.NewPermission(roleStr)
				if err != nil {
					panic(err)
				}
				requiredPermission = opt.New(r)
			}
			if t, ok := ext["title"].(string); ok {
				displayTitle = t
			}
			requiresConfirmation, _ = ext["requires_confirmation"].(bool)
			if ann, ok := ext["annotations"].(map[string]any); ok {
				annotations.ReadOnlyHint, _ = ann["readOnlyHint"].(bool)
				annotations.DestructiveHint, _ = ann["destructiveHint"].(bool)
				annotations.IdempotentHint, _ = ann["idempotentHint"].(bool)
				annotations.OpenWorldHint, _ = ann["openWorldHint"].(bool)
			}
		}
	}

	return &ToolDefinition{
		Name:                 toolName,
		Title:                displayTitle,
		Description:          description,
		InputSchema:          inputSchema,
		OutputSchema:         outputSchema,
		RequiredPermission:   requiredPermission,
		RequiresConfirmation: requiresConfirmation,
		Annotations:          annotations,
	}
}

func GetToolNamesEnum() []any {
	return toolNamesEnum
}

func InjectToolNamesEnum(schema *jsonschema.Schema, propertyName string) {
	if schema.Properties == nil {
		return
	}

	prop, ok := schema.Properties[propertyName]
	if !ok || prop.Items == nil {
		return
	}

	prop.Items.Enum = toolNamesEnum
}
