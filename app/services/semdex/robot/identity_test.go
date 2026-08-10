package robot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRobotIdentityInstructionDescribesCoordinator(t *testing.T) {
	instruction := robotIdentityInstruction(robotIdentityContext{Current: robotIdentity{
		Name:         "Denbot",
		Description:  "Coordinates Storyden work.",
		Capabilities: []string{"toolset_search", "robot_search", "delegation"},
	}})

	assert.Contains(t, instruction, "visible coordinator")
	assert.Contains(t, instruction, "toolset_search")
	assert.Contains(t, instruction, "delegated inside this session")
}
