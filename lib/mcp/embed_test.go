package mcp

import (
	"strings"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRobotSearchTool(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	toolDef := GetRobotSearchTool()
	r.NotNil(toolDef, "Tool definition should be loaded")

	a.Equal("robot_search", toolDef.Name)
	a.NotEmpty(toolDef.Description, "Description should be loaded from YAML")

	r.NotNil(toolDef.InputSchema, "Input schema should be resolved")
	r.NotNil(toolDef.OutputSchema, "Output schema should be resolved")

	// Verify input schema structure
	a.Equal("object", toolDef.InputSchema.Type)
	a.NotNil(toolDef.InputSchema.Properties)

	queryProp, ok := toolDef.InputSchema.Properties["query"]
	a.True(ok, "Input schema should have query property")
	a.Equal("string", queryProp.Type)

	// Verify output schema structure
	a.Equal("object", toolDef.OutputSchema.Type)
	a.NotNil(toolDef.OutputSchema.Properties)

	robotsProp, ok := toolDef.OutputSchema.Properties["robots"]
	a.True(ok, "Output schema should have robots property")
	a.Equal("array", robotsProp.Type)
}

func TestAllToolSchemasResolve(t *testing.T) {
	tools := []*ToolDefinition{
		toolContentSearch,
		toolDocumentClose,
		toolDocumentGet,
		toolDocumentList,
		toolDocumentSearch,
		toolRobotSearch,
		toolToolSearch,
		toolToolGet,
		toolToolLoad,
		toolRobotCreate,
		toolRobotList,
		toolRobotGet,
		toolRobotUpdate,
		toolRobotDelete,
		toolToolsetCreate,
		toolToolsetDelete,
		toolToolsetGet,
		toolToolsetList,
		toolToolsetSearch,
		toolToolsetLoad,
		toolToolsetUpdate,
		toolLibraryPageTree,
		toolLibraryPageGet,
		toolLibraryPageOpen,
		toolLibraryPageCreate,
		toolLibraryPageUpdate,
		toolLibrarySearchPages,
		toolLibraryPagePropertySchemaGet,
		toolLibraryPagePropertySchemaUpdate,
		toolLibraryPagePropertiesUpdate,
		toolTagList,
		toolLinkCreate,
		toolWebFetch,
		toolWebOpen,
		toolThreadCreate,
		toolThreadGet,
		toolThreadOpen,
		toolThreadList,
		toolThreadUpdate,
		toolThreadReply,
		toolCategoryList,
		toolThreadSearch,
		toolReplySearch,
		toolPostSearch,
		toolMemberSearch,
	}

	for _, toolDef := range tools {
		require.NotNil(t, toolDef)

		t.Run(toolDef.Name, func(t *testing.T) {
			if toolDef.InputSchema != nil {
				_, err := toolDef.InputSchema.Resolve(nil)
				require.NoError(t, err, "input schema should resolve")
			}

			if toolDef.OutputSchema != nil {
				_, err := toolDef.OutputSchema.Resolve(nil)
				require.NoError(t, err, "output schema should resolve")
			}
		})
	}
}

func TestEveryToolHasUnifiedStorydenMetadata(t *testing.T) {
	for definitionName, schema := range Schema.Definitions {
		if !strings.HasPrefix(definitionName, "Tool") {
			continue
		}
		if _, ok := schema.Properties["input"]; !ok {
			continue
		}
		if _, ok := schema.Properties["output"]; !ok {
			continue
		}

		_, hasLegacyExtension := schema.Extra["x-storyden"]
		assert.False(t, hasLegacyExtension, "%s must not define legacy x-storyden metadata", definitionName)

		extension, ok := schema.Extra["x-storyden-tool"].(map[string]any)
		require.True(t, ok, "%s must define x-storyden-tool", definitionName)
		title, ok := extension["title"].(string)
		require.True(t, ok, "%s must define x-storyden-tool.title", definitionName)
		assert.NotEmpty(t, strings.TrimSpace(title), "%s must define a non-empty x-storyden-tool.title", definitionName)
		_, ok = extension["toolsets"].([]any)
		require.True(t, ok, "%s must define x-storyden-tool.toolsets", definitionName)
	}
}

func TestSessionLoadToolsExposeFrontendMetadata(t *testing.T) {
	toolLoad := GetToolLoadTool()
	require.NotNil(t, toolLoad)
	assert.True(t, toolLoad.Internal)
	assert.Equal(t, "Load Tools", toolLoad.Title)
	assert.False(t, toolLoad.Annotations.ReadOnlyHint)
	assert.False(t, toolLoad.Annotations.DestructiveHint)
	assert.True(t, toolLoad.Annotations.IdempotentHint)
	assert.False(t, toolLoad.Annotations.OpenWorldHint)

	toolsetLoad := GetToolsetLoadTool()
	require.NotNil(t, toolsetLoad)
	assert.True(t, toolsetLoad.Internal)
	assert.Equal(t, "Load Toolsets", toolsetLoad.Title)
	assert.False(t, toolsetLoad.Annotations.ReadOnlyHint)
	assert.False(t, toolsetLoad.Annotations.DestructiveHint)
	assert.True(t, toolsetLoad.Annotations.IdempotentHint)
	assert.False(t, toolsetLoad.Annotations.OpenWorldHint)
}

func TestInternalToolMetadata(t *testing.T) {
	assert.True(t, GetLibraryRequestPageTool().Internal)
	assert.True(t, GetDocumentGetTool().Internal)
	assert.True(t, GetDocumentSearchTool().Internal)
	assert.True(t, GetDocumentListTool().Internal)
	assert.True(t, GetDocumentCloseTool().Internal)
	assert.True(t, GetLibraryPageOpenTool().Internal)
	assert.True(t, GetThreadOpenTool().Internal)
	assert.True(t, GetWebOpenTool().Internal)
	assert.False(t, GetLibraryPageGetTool().Internal)
	assert.False(t, GetWebFetchTool().Internal)
}

func TestDocumentNavigationToolsAreToolsetOnly(t *testing.T) {
	for _, tool := range []*ToolDefinition{
		GetDocumentGetTool(),
		GetDocumentSearchTool(),
		GetDocumentListTool(),
		GetDocumentCloseTool(),
	} {
		assert.True(t, tool.ToolsetOnly, tool.Name)
		assert.Equal(t, []string{"system.documents"}, tool.Toolsets, tool.Name)
	}
}

func TestCapabilityDiscoverySchemasAreProgressive(t *testing.T) {
	toolSearch := GetToolSearchTool()
	require.NotNil(t, toolSearch.OutputSchema)
	toolItems := toolSearch.OutputSchema.Properties["tools"].Items
	require.NotNil(t, toolItems)
	assert.ElementsMatch(t, []string{"id", "name", "description"}, schemaPropertyNames(toolItems))
	assert.NotContains(t, toolItems.Properties, "source")
	assert.NotContains(t, toolItems.Properties, "input_schema")

	toolGet := GetToolGetTool()
	require.NotNil(t, toolGet.OutputSchema)
	assert.ElementsMatch(t,
		[]string{"id", "name", "description", "input_schema", "output_schema", "toolsets", "toolset_only", "requires_confirmation", "requires_workspace"},
		schemaPropertyNames(toolGet.OutputSchema),
	)
	assert.NotContains(t, toolGet.OutputSchema.Properties, "source")

	toolsetSearch := GetToolsetSearchTool()
	require.NotNil(t, toolsetSearch.OutputSchema)
	toolsetItems := toolsetSearch.OutputSchema.Properties["toolsets"].Items
	require.NotNil(t, toolsetItems)
	assert.ElementsMatch(t, []string{"id", "name", "description"}, schemaPropertyNames(toolsetItems))
	assert.NotContains(t, toolsetItems.Properties, "source")
	assert.NotContains(t, toolsetItems.Properties, "tools")

	toolsetGet := GetToolsetGetTool()
	require.NotNil(t, toolsetGet.OutputSchema)
	assert.ElementsMatch(t,
		[]string{"id", "name", "description", "tools", "instruction", "editable", "requires_workspace"},
		schemaPropertyNames(toolsetGet.OutputSchema),
	)
	assert.NotContains(t, toolsetGet.OutputSchema.Properties, "source")
}

func schemaPropertyNames(schema *jsonschema.Schema) []string {
	properties := make([]string, 0, len(schema.Properties))
	for name := range schema.Properties {
		properties = append(properties, name)
	}
	return properties
}
