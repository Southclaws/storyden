package robot

import (
	"testing"

	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry"
)

func TestAnnotateClientToolCallsBindsCallToBrowserAndPage(t *testing.T) {
	clientTools, err := NewClientToolContext("browser-1", []agent_registry.ClientToolDefinition{{
		Name:        "library_page_block_add",
		Title:       "Add library block",
		Description: "Add one block.",
		InputSchema: map[string]any{"type": "object"},
	}})
	if err != nil {
		t.Fatalf("NewClientToolContext() error = %v", err)
	}
	response := &model.LLMResponse{Content: &genai.Content{Parts: []*genai.Part{{
		FunctionCall: &genai.FunctionCall{ID: "call-1", Name: "library_page_block_add"},
	}}}}

	callback := annotateClientToolCalls(clientTools, InvocationContext{
		InvocationContextKeyPageType: "library",
	})
	if _, err := callback(nil, response, nil); err != nil {
		t.Fatalf("annotateClientToolCalls() error = %v", err)
	}

	metadata, ok := response.Content.Parts[0].PartMetadata[agent_registry.ClientToolMetadataKey].(map[string]any)
	if !ok {
		t.Fatalf("client tool metadata = %#v, want map", response.Content.Parts[0].PartMetadata)
	}
	if metadata["source"] != "webmcp" {
		t.Fatalf("source = %#v, want webmcp", metadata["source"])
	}
	if metadata["client_id"] != "browser-1" {
		t.Fatalf("client_id = %#v, want browser-1", metadata["client_id"])
	}
	scope, ok := metadata["scope"].(InvocationContext)
	if !ok || scope[InvocationContextKeyPageType] != "library" {
		t.Fatalf("scope = %#v, want library page context", metadata["scope"])
	}
}

func TestAnnotateClientToolCallsIgnoresServerTool(t *testing.T) {
	clientTools, err := NewClientToolContext("browser-1", []agent_registry.ClientToolDefinition{{
		Name:        "library_page_block_add",
		Description: "Add one block.",
		InputSchema: map[string]any{"type": "object"},
	}})
	if err != nil {
		t.Fatalf("NewClientToolContext() error = %v", err)
	}
	response := &model.LLMResponse{Content: &genai.Content{Parts: []*genai.Part{{
		FunctionCall: &genai.FunctionCall{ID: "call-1", Name: "content_search"},
	}}}}

	callback := annotateClientToolCalls(clientTools, nil)
	if _, err := callback(nil, response, nil); err != nil {
		t.Fatalf("annotateClientToolCalls() error = %v", err)
	}
	if response.Content.Parts[0].PartMetadata != nil {
		t.Fatalf("PartMetadata = %#v, want nil", response.Content.Parts[0].PartMetadata)
	}
}

func TestBuildClientToolsetDeclaresLongRunningBrowserTool(t *testing.T) {
	clientTools, err := NewClientToolContext("browser-1", []agent_registry.ClientToolDefinition{{
		Name:        "library_page_layout_get",
		Description: "Read the current library layout.",
		InputSchema: map[string]any{"type": "object"},
	}})
	if err != nil {
		t.Fatalf("NewClientToolContext() error = %v", err)
	}

	toolset, err := buildClientToolset(clientTools)
	if err != nil {
		t.Fatalf("buildClientToolset() error = %v", err)
	}
	toolList, err := toolset.Tools(nil)
	if err != nil {
		t.Fatalf("Toolset.Tools() error = %v", err)
	}
	if len(toolList) != 1 {
		t.Fatalf("len(tools) = %d, want 1", len(toolList))
	}
	if toolList[0].Name() != "library_page_layout_get" {
		t.Fatalf("tool name = %q, want library_page_layout_get", toolList[0].Name())
	}
	if !toolList[0].IsLongRunning() {
		t.Fatal("IsLongRunning() = false, want true")
	}
}
