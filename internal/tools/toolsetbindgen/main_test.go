package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseToolsetBindingsGroupsAndSortsToolNames(t *testing.T) {
	bindings, err := parseToolsetBindings([]byte(`{
  "definitions": {
    "ToolThreadList": {
      "title": "thread_list",
	  "x-storyden-tool": {"title": "List discussions"},
      "x-storyden": {"toolsets": ["system.discussions"]}
    },
    "ToolLibrarySearchPages": {
      "title": "library_search_pages",
	  "x-storyden-tool": {"title": "Search library pages"},
      "x-storyden": {"toolsets": ["system.library", "system.content_search"]}
    },
    "ToolThreadGet": {
      "title": "thread_get",
	  "x-storyden-tool": {"title": "Get discussion"},
      "x-storyden": {"toolsets": ["system.discussions", "system.discussions"]}
    },
    "ToolThreadGetInput": {"type": "object"}
  }
}`))
	require.NoError(t, err)
	require.Len(t, bindings, 3)

	assert.Equal(t, "system.content_search", bindings[0].ID)
	assert.Equal(t, "system_content_search", bindings[0].PackageName)
	assert.Equal(t, []string{"library_search_pages"}, bindings[0].ToolNames)
	assert.Equal(t, "system.discussions", bindings[1].ID)
	assert.Equal(t, []string{"thread_get", "thread_list"}, bindings[1].ToolNames)
	assert.Equal(t, "system.library", bindings[2].ID)
	assert.Equal(t, []string{"library_search_pages"}, bindings[2].ToolNames)
}

func TestParseToolsetBindingsRejectsInvalidMembership(t *testing.T) {
	_, err := parseToolsetBindings([]byte(`{
  "definitions": {
    "ToolThreadGet": {
      "title": "thread_get",
	  "x-storyden-tool": {"title": "Get discussion"},
      "x-storyden": {"toolsets": ["custom.user_tools"]}
    }
  }
}`))
	require.Error(t, err)
	assert.Contains(t, err.Error(), `invalid system Toolset ID "custom.user_tools"`)
}

func TestParseToolsetBindingsRequiresExplicitMembershipArray(t *testing.T) {
	_, err := parseToolsetBindings([]byte(`{
  "definitions": {
    "ToolLinkCreate": {
      "title": "link_create",
      "x-storyden-tool": {"title": "Create Link"}
    }
  }
}`))
	require.Error(t, err)
	assert.EqualError(t, err, "ToolLinkCreate must declare x-storyden.toolsets")
}

func TestWriteBindingsTargetsDefinitionPackages(t *testing.T) {
	outputDir := t.TempDir()
	definitionDir := filepath.Join(outputDir, "system_discussions")
	require.NoError(t, os.MkdirAll(definitionDir, 0o755))
	require.NoError(t, os.WriteFile(
		filepath.Join(definitionDir, "definition.go"),
		[]byte("package system_discussions\n"),
		0o644,
	))

	err := writeBindings(outputDir, []toolsetBinding{{
		ID:          "system.discussions",
		PackageName: "system_discussions",
		ToolNames:   []string{"thread_get", "thread_list"},
	}})
	require.NoError(t, err)

	toolNames, err := os.ReadFile(filepath.Join(definitionDir, "tools_gen.go"))
	require.NoError(t, err)
	assert.Contains(t, string(toolNames), "package system_discussions")
	assert.Contains(t, string(toolNames), `"thread_get"`)
	assert.Contains(t, string(toolNames), `"thread_list"`)

	registrations, err := os.ReadFile(filepath.Join(outputDir, "system_bindings_gen.go"))
	require.NoError(t, err)
	assert.Contains(t, string(registrations), `system_discussions.ID`)
	assert.Contains(t, string(registrations), `system_discussions.Instruction`)
}
