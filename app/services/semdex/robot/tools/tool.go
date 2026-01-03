package tools

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/storyden/mcp"
	adkagent "google.golang.org/adk/agent"
	adktool "google.golang.org/adk/tool"
)

// ToolResult is just because adk doesn't allow returning an error. Since adk
// just serialises the return into json, this lets us tell the LLM about errors.
type ToolResult[T any] struct {
	Result T      `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func NewSuccess[T any](v T) ToolResult[T] {
	return ToolResult[T]{Result: v}
}

func NewError[T any](err error) ToolResult[T] {
	return ToolResult[T]{Error: err.Error()}
}

func NewErrorMsg[T any](msg string) ToolResult[T] {
	return ToolResult[T]{Error: msg}
}

// Tool is a wrapper around an actual tool definition, it includes the actual
// executor as well as the definition and a flag for if it's a client-side tool.
type Tool struct {
	Definition *mcp.ToolDefinition
	// Builder actually constructs the executor function dynamically. This is
	// useful for tools that depend on runtime context, like database data.
	Builder      func(context.Context) (adktool.Tool, error)
	IsClientTool bool
}

type Tools []*Tool

func (t Tools) ToADKTools(ctx context.Context) ([]adktool.Tool, error) {
	var adkTools []adktool.Tool
	for _, tool := range t {
		execTool, err := tool.Builder(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fmsg.Withf("failed to build tool %q", tool.Definition.Name))
		}
		adkTools = append(adkTools, execTool)
	}
	return adkTools, nil
}

var _ adktool.Tool = (*Tool)(nil)

func (t *Tool) Name() string {
	return t.Definition.Name
}

func (t *Tool) Description() string {
	return t.Definition.Description
}

func (t *Tool) IsLongRunning() bool {
	return t.IsClientTool
}

// Toolset implements adktool.Toolset
type Toolset struct {
	name     string
	ToolList []adktool.Tool
}

var _ adktool.Toolset = (*Toolset)(nil)

func (d *Toolset) Name() string {
	return d.name
}

func (d *Toolset) Tools(ctx adkagent.ReadonlyContext) ([]adktool.Tool, error) {
	return d.ToolList, nil
}
