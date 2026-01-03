package mcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRobotSwitchTool(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	toolDef := GetRobotSwitchTool()
	r.NotNil(toolDef, "Tool definition should be loaded")

	a.Equal("robot_switch", toolDef.Name)
	a.NotEmpty(toolDef.Description, "Description should be loaded from YAML")

	r.NotNil(toolDef.InputSchema, "Input schema should be resolved")
	r.NotNil(toolDef.OutputSchema, "Output schema should be resolved")

	// Verify input schema structure
	a.Equal("object", toolDef.InputSchema.Type)
	a.NotNil(toolDef.InputSchema.Properties)

	robotIDProp, ok := toolDef.InputSchema.Properties["robot_id"]
	a.True(ok, "Input schema should have robot_id property")
	a.Equal("string", robotIDProp.Type)

	// Verify output schema structure
	a.Equal("object", toolDef.OutputSchema.Type)
	a.NotNil(toolDef.OutputSchema.Properties)

	successProp, ok := toolDef.OutputSchema.Properties["success"]
	a.True(ok, "Output schema should have success property")
	a.Equal("boolean", successProp.Type)
}
