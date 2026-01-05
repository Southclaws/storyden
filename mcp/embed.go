package mcp

import (
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"go.yaml.in/yaml/v4"
)

//go:embed schema.yaml
var SchemaYAML []byte

// Schema represents the entire schema for storyden agent tools.
var Schema *jsonschema.Schema

type ToolDefinition struct {
	Name         string
	Description  string
	InputSchema  *jsonschema.Schema
	OutputSchema *jsonschema.Schema
}

var (
	toolSearch               *ToolDefinition
	toolRobotSwitch          *ToolDefinition
	toolRobotGetAllToolNames *ToolDefinition
	toolRobotCreate          *ToolDefinition
	toolRobotList            *ToolDefinition
	toolRobotGet             *ToolDefinition
	toolRobotUpdate          *ToolDefinition
	toolRobotDelete          *ToolDefinition
	toolLibraryPageTree      *ToolDefinition
	toolLibraryPageGet       *ToolDefinition
	toolLibraryPageCreate    *ToolDefinition
	toolLibraryPageUpdate    *ToolDefinition
	toolLibraryPageSearch    *ToolDefinition
	toolTagList              *ToolDefinition
	toolLinkCreate           *ToolDefinition
	toolThreadCreate         *ToolDefinition
	toolThreadGet            *ToolDefinition
	toolThreadList           *ToolDefinition
	toolThreadUpdate         *ToolDefinition
	toolThreadReply          *ToolDefinition
	toolCategoryList         *ToolDefinition
	toolNamesEnum            []any
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
	var yamlData map[string]any
	err := yaml.Unmarshal(SchemaYAML, &yamlData)
	if err != nil {
		panic("failed to unmarshal MCP schema YAML: " + err.Error())
	}

	// NOTE: bit of a hack, it fixes some oddities with $refs...
	jsonBytes, err := json.Marshal(yamlData)
	if err != nil {
		panic("failed to marshal schema to JSON: " + err.Error())
	}

	Schema = &jsonschema.Schema{}
	err = json.Unmarshal(jsonBytes, Schema)
	if err != nil {
		panic("failed to unmarshal schema from JSON: " + err.Error())
	}

	initAllTools()
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
					inputSchema = def
				}
			}
		} else if inputProp.Type != "" {
			inputSchema = inputProp
		}
	}

	var outputSchema *jsonschema.Schema
	if outputProp, ok := toolDef.Properties["output"]; ok {
		if outputProp.Ref != "" {

			parts := strings.Split(outputProp.Ref, "/")
			if len(parts) > 0 {
				defName := parts[len(parts)-1]
				if def, ok := Schema.Definitions[defName]; ok {
					outputSchema = def
				}
			}
		} else if outputProp.Type != "" {
			outputSchema = outputProp
		}
	}

	return &ToolDefinition{
		Name:         toolName,
		Description:  description,
		InputSchema:  inputSchema,
		OutputSchema: outputSchema,
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
