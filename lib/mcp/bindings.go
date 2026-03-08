package mcp

import (
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/google/jsonschema-go/jsonschema"
)

//go:embed robots.json
var schemaJSON []byte

// Schema represents the entire schema for storyden agent tools.
var Schema *jsonschema.Schema

type ToolDefinition struct {
	Name               string
	Description        string
	InputSchema        *jsonschema.Schema
	OutputSchema       *jsonschema.Schema
	RequiredPermission rbac.Permission
}

var (
	toolSearch                          *ToolDefinition
	toolRobotSwitch                     *ToolDefinition
	toolRobotGetAllToolNames            *ToolDefinition
	toolRobotCreate                     *ToolDefinition
	toolRobotList                       *ToolDefinition
	toolRobotGet                        *ToolDefinition
	toolRobotUpdate                     *ToolDefinition
	toolRobotDelete                     *ToolDefinition
	toolLibraryPageTree                 *ToolDefinition
	toolLibraryPageGet                  *ToolDefinition
	toolLibraryPageCreate               *ToolDefinition
	toolLibraryPageUpdate               *ToolDefinition
	toolLibraryPageSearch               *ToolDefinition
	toolLibraryPagePropertySchemaGet    *ToolDefinition
	toolLibraryPagePropertySchemaUpdate *ToolDefinition
	toolLibraryPagePropertiesUpdate     *ToolDefinition
	toolTagList                         *ToolDefinition
	toolLinkCreate                      *ToolDefinition
	toolThreadCreate                    *ToolDefinition
	toolThreadGet                       *ToolDefinition
	toolThreadList                      *ToolDefinition
	toolThreadUpdate                    *ToolDefinition
	toolThreadReply                     *ToolDefinition
	toolCategoryList                    *ToolDefinition
	toolNamesEnum                       []any
)

func initAllTools() {
	toolSearch = initTool("Search")
	toolRobotSwitch = initTool("RobotSwitch")
	toolRobotGetAllToolNames = initTool("RobotGetAllToolNames")
	toolRobotCreate = initTool("RobotCreate")
	toolRobotList = initTool("RobotList")
	toolRobotGet = initTool("RobotGet")
	toolRobotUpdate = initTool("RobotUpdate")
	toolRobotDelete = initTool("RobotDelete")
	toolLibraryPageTree = initTool("LibraryPageTree")
	toolLibraryPageGet = initTool("LibraryPageGet")
	toolLibraryPageCreate = initTool("LibraryPageCreate")
	toolLibraryPageUpdate = initTool("LibraryPageUpdate")
	toolLibraryPageSearch = initTool("LibraryPageSearch")
	toolLibraryPagePropertySchemaGet = initTool("LibraryPagePropertySchemaGet")
	toolLibraryPagePropertySchemaUpdate = initTool("LibraryPagePropertySchemaUpdate")
	toolLibraryPagePropertiesUpdate = initTool("LibraryPagePropertiesUpdate")
	toolTagList = initTool("TagList")
	toolLinkCreate = initTool("LinkCreate")
	toolThreadCreate = initTool("ThreadCreate")
	toolThreadGet = initTool("ThreadGet")
	toolThreadList = initTool("ThreadList")
	toolThreadUpdate = initTool("ThreadUpdate")
	toolThreadReply = initTool("ThreadReply")
	toolCategoryList = initTool("CategoryList")

	names := AllToolNames()
	toolNamesEnum = make([]any, len(names))
	for i, name := range names {
		toolNamesEnum[i] = name
	}
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

	// Extract x-storyden-role from Extra field if present
	var requiredPermission rbac.Permission
	if toolDef.Extra != nil {
		if role, ok := toolDef.Extra["x-storyden-role"]; ok {
			if roleStr, ok := role.(string); ok {
				r, err := rbac.NewPermission(roleStr)
				if err != nil {
					// TODO: Actually check this at codegen time, not runtime.
					panic(err)
				}
				requiredPermission = r
			}
		}
	}

	return &ToolDefinition{
		Name:               toolName,
		Description:        description,
		InputSchema:        inputSchema,
		OutputSchema:       outputSchema,
		RequiredPermission: requiredPermission,
	}
}

func GetSearchTool() *ToolDefinition {
	return toolSearch
}

func GetRobotSwitchTool() *ToolDefinition {
	return toolRobotSwitch
}

func GetRobotGetAllToolNamesTool() *ToolDefinition {
	return toolRobotGetAllToolNames
}

func GetRobotCreateTool() *ToolDefinition {
	return toolRobotCreate
}

func GetRobotListTool() *ToolDefinition {
	return toolRobotList
}

func GetRobotGetTool() *ToolDefinition {
	return toolRobotGet
}

func GetRobotUpdateTool() *ToolDefinition {
	return toolRobotUpdate
}

func GetRobotDeleteTool() *ToolDefinition {
	return toolRobotDelete
}

func GetLibraryPageTreeTool() *ToolDefinition {
	return toolLibraryPageTree
}

func GetLibraryPageGetTool() *ToolDefinition {
	return toolLibraryPageGet
}

func GetLibraryPageCreateTool() *ToolDefinition {
	return toolLibraryPageCreate
}

func GetLibraryPageUpdateTool() *ToolDefinition {
	return toolLibraryPageUpdate
}

func GetLibraryPageSearchTool() *ToolDefinition {
	return toolLibraryPageSearch
}

func GetLibraryPagePropertySchemaGetTool() *ToolDefinition {
	return toolLibraryPagePropertySchemaGet
}

func GetLibraryPagePropertySchemaUpdateTool() *ToolDefinition {
	return toolLibraryPagePropertySchemaUpdate
}

func GetLibraryPagePropertiesUpdateTool() *ToolDefinition {
	return toolLibraryPagePropertiesUpdate
}

func GetTagListTool() *ToolDefinition {
	return toolTagList
}

func GetLinkCreateTool() *ToolDefinition {
	return toolLinkCreate
}

func GetThreadCreateTool() *ToolDefinition {
	return toolThreadCreate
}

func GetThreadGetTool() *ToolDefinition {
	return toolThreadGet
}

func GetThreadListTool() *ToolDefinition {
	return toolThreadList
}

func GetThreadUpdateTool() *ToolDefinition {
	return toolThreadUpdate
}

func GetThreadReplyTool() *ToolDefinition {
	return toolThreadReply
}

func GetCategoryListTool() *ToolDefinition {
	return toolCategoryList
}

func AllToolNames() []string {
	return []string{
		toolSearch.Name,
		toolRobotSwitch.Name,
		toolRobotGetAllToolNames.Name,
		toolRobotCreate.Name,
		toolRobotList.Name,
		toolRobotGet.Name,
		toolRobotUpdate.Name,
		toolRobotDelete.Name,
		toolLibraryPageTree.Name,
		toolLibraryPageGet.Name,
		toolLibraryPageCreate.Name,
		toolLibraryPageUpdate.Name,
		toolLibraryPageSearch.Name,
		toolLibraryPagePropertySchemaGet.Name,
		toolLibraryPagePropertySchemaUpdate.Name,
		toolLibraryPagePropertiesUpdate.Name,
		toolTagList.Name,
		toolLinkCreate.Name,
		toolThreadCreate.Name,
		toolThreadGet.Name,
		toolThreadList.Name,
		toolThreadUpdate.Name,
		toolThreadReply.Name,
		toolCategoryList.Name,
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
