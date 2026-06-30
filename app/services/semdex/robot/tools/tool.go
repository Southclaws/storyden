package tools

import (
	"context"
	"encoding/json"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/storyden/lib/mcp"
	adkagent "google.golang.org/adk/v2/agent"
	adktool "google.golang.org/adk/v2/tool"
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

// Handler is a context.Context-based handler for use by non-ADK transports
// (e.g. the MCP SSE transport). It receives raw JSON arguments and returns raw
// JSON output, avoiding a dependency on google.golang.org/adk/tool.Context.
type Handler func(ctx context.Context, args json.RawMessage) (json.RawMessage, error)

// Tool is a wrapper around an actual tool definition, it includes the actual
// executor as well as the definition and a flag for if it's a client-side tool.
type Tool struct {
	Definition *mcp.ToolDefinition
	// Source identifies where the tool came from for catalogues and UI labels.
	// Empty means a built-in Storyden tool.
	Source string
	// CallableName is the name exposed to ADK/model providers. Native tools use
	// Definition.Name; dynamic MCP tools may need a stricter generated name.
	CallableName string
	// Builder actually constructs the executor function dynamically. This is
	// useful for tools that depend on runtime context, like database data.
	Builder      func(context.Context) (adktool.Tool, error)
	Handler      Handler
	IsClientTool bool
}

func makeHandler[T, O any](execute func(context.Context, T) (*O, error)) Handler {
	return func(ctx context.Context, args json.RawMessage) (json.RawMessage, error) {
		var input T
		if err := json.Unmarshal(args, &input); err != nil {
			return nil, err
		}
		output, err := execute(ctx, input)
		if err != nil {
			return nil, err
		}
		return json.Marshal(output)
	}
}

type Tools []*Tool

type confirmationDisabledContextKey struct{}

func ContextWithConfirmationDisabled(ctx context.Context) context.Context {
	return context.WithValue(ctx, confirmationDisabledContextKey{}, true)
}

func confirmationDisabled(ctx context.Context) bool {
	v, _ := ctx.Value(confirmationDisabledContextKey{}).(bool)
	return v
}

func ConfirmationDisabled(ctx context.Context) bool {
	return confirmationDisabled(ctx)
}

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

func (t *Tool) ADKName() string {
	if t.CallableName != "" {
		return t.CallableName
	}
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
