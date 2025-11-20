package agent

import (
	"errors"
	"fmt"
	"strings"

	storydentools "github.com/Southclaws/storyden/app/services/semdex/agent/tools"
	markmcp "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	agentpkg "google.golang.org/adk/agent"
	"google.golang.org/adk/model"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/genai"
)

type adapterToolset struct {
	name  string
	tools []adktool.Tool
}

func newAdapterToolset(all storydentools.All) (adktool.Toolset, error) {
	adapted := make([]adktool.Tool, 0, len(all))
	for _, srvTool := range all {
		adapter, err := newServerToolAdapter(srvTool)
		if err != nil {
			return nil, err
		}
		adapted = append(adapted, adapter)
	}

	return &adapterToolset{name: defaultToolsetName, tools: adapted}, nil
}

func (t *adapterToolset) Name() string {
	return t.name
}

func (t *adapterToolset) Tools(ctx agentpkg.ReadonlyContext) ([]adktool.Tool, error) {
	return t.tools, nil
}

type serverToolAdapter struct {
	base        server.ServerTool
	declaration *genai.FunctionDeclaration
}

func newServerToolAdapter(tool server.ServerTool) (*serverToolAdapter, error) {
	params, err := convertSchema(tool.Tool.InputSchema, tool.Tool.RawInputSchema)
	if err != nil {
		return nil, fmt.Errorf("convert input schema for %s: %w", tool.Tool.Name, err)
	}
	output, err := convertSchema(tool.Tool.OutputSchema, tool.Tool.RawOutputSchema)
	if err != nil {
		return nil, fmt.Errorf("convert output schema for %s: %w", tool.Tool.Name, err)
	}

	decl := &genai.FunctionDeclaration{
		Name:                 tool.Tool.Name,
		Description:          tool.Tool.Description,
		ParametersJsonSchema: params,
		ResponseJsonSchema:   output,
	}

	return &serverToolAdapter{base: tool, declaration: decl}, nil
}

func (s *serverToolAdapter) Name() string {
	return s.base.Tool.Name
}

func (s *serverToolAdapter) Description() string {
	return s.base.Tool.Description
}

func (*serverToolAdapter) IsLongRunning() bool { return false }

func (s *serverToolAdapter) Declaration() *genai.FunctionDeclaration {
	return s.declaration
}

func (s *serverToolAdapter) ProcessRequest(ctx adktool.Context, req *model.LLMRequest) error {
	return packTool(req, s)
}

func (s *serverToolAdapter) Run(ctx adktool.Context, args any) (map[string]any, error) {
	request := markmcp.CallToolRequest{
		Params: markmcp.CallToolParams{
			Name:      s.base.Tool.Name,
			Arguments: args,
		},
	}

	res, err := s.base.Handler(ctx, request)
	if err != nil {
		return nil, err
	}

	return adaptToolResult(res)
}

func adaptToolResult(res *markmcp.CallToolResult) (map[string]any, error) {
	if res == nil {
		return nil, errors.New("tool returned no result")
	}

	if res.IsError {
		message := extractText(res.Content)
		if message == "" {
			message = "tool execution failed"
		}
		return nil, errors.New(message)
	}

	if res.StructuredContent != nil {
		return map[string]any{"output": res.StructuredContent}, nil
	}

	text := extractText(res.Content)
	if text == "" {
		return nil, errors.New("tool response did not include text content")
	}

	return map[string]any{"output": text}, nil
}

func extractText(contents []markmcp.Content) string {
	var builder strings.Builder
	for _, c := range contents {
		switch val := c.(type) {
		case markmcp.TextContent:
			builder.WriteString(val.Text)
		case *markmcp.TextContent:
			builder.WriteString(val.Text)
		}
	}
	return builder.String()
}

func packTool(req *model.LLMRequest, tool toolWithDeclaration) error {
	if req.Tools == nil {
		req.Tools = make(map[string]any)
	}

	name := tool.Name()
	if _, exists := req.Tools[name]; exists {
		return fmt.Errorf("duplicate tool: %q", name)
	}

	req.Tools[name] = tool

	if req.Config == nil {
		req.Config = &genai.GenerateContentConfig{}
	}

	decl := tool.Declaration()
	if decl == nil {
		return nil
	}

	var fnTool *genai.Tool
	for _, cfgTool := range req.Config.Tools {
		if cfgTool != nil && cfgTool.FunctionDeclarations != nil {
			fnTool = cfgTool
			break
		}
	}

	if fnTool == nil {
		req.Config.Tools = append(req.Config.Tools, &genai.Tool{
			FunctionDeclarations: []*genai.FunctionDeclaration{decl},
		})
		return nil
	}

	fnTool.FunctionDeclarations = append(fnTool.FunctionDeclarations, decl)
	return nil
}

type toolWithDeclaration interface {
	adktool.Tool
	Declaration() *genai.FunctionDeclaration
}
