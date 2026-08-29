package robot

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	robot_resource "github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_memory"
)

func TestFormatMemoryRootIsBoundedAndShallow(t *testing.T) {
	items := make([]*robot_memory.Item, 20)
	for i := range items {
		id := robot_resource.MemoryID(xid.New())
		items[i] = &robot_memory.Item{Memory: &robot_resource.Memory{ID: id}, Excerpt: fmt.Sprintf("memory excerpt %02d", i)}
	}
	items[0].Memory.Fact = opt.New(robot_resource.MemoryFact{Subject: "freyja", Predicate: "owned_by", Object: "southclaws"})

	result := formatMemoryRoot(items, true)

	assert.Equal(t, 20, strings.Count(result, "memory excerpt"))
	assert.Contains(t, result, "[freyja, owned_by, southclaws]")
	assert.Contains(t, result, "Active top-level knowledge graph nodes")
	assert.Contains(t, result, "Additional top-level memories were omitted")
	assert.NotContains(t, result, "complete content")
	assert.NotContains(t, result, "memory_list")
	require.LessOrEqual(t, strings.Count(result, "\n"), 21)
}

func TestFormatMemoryRootEmptyRejectsSessionState(t *testing.T) {
	assert.Equal(t,
		"There are no active top-level knowledge graph nodes. Do not create session-specific task state. Save a useful long-term fact directly when one emerges.",
		formatMemoryRoot(nil, false),
	)
}
