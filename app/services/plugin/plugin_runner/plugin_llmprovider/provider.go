package plugin_llmprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"strings"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"google.golang.org/adk/model"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/robot/llm_provider"
	"github.com/Southclaws/storyden/app/resources/robot/model_ref"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type Provider struct {
	provider model_ref.Provider
	session  plugin_runner.Session
}

func New(provider model_ref.Provider, session plugin_runner.Session) *Provider {
	return &Provider{
		provider: provider,
		session:  session,
	}
}

func (p *Provider) Provider() model_ref.Provider { return p.provider }

func (p *Provider) RequiresAPIKey() bool { return false }

func (p *Provider) Configure(llm_provider.Config) {}

func (p *Provider) ListModels(ctx context.Context) ([]model_ref.Info, error) {
	id := xid.New()
	resp, err := p.session.Send(ctx, id, &rpc.RPCRequestRobotModelProviderListModels{
		ID:      id,
		Jsonrpc: "2.0",
		Method:  "robot_model_provider_list_models",
		Params: rpc.RPCRequestRobotModelProviderListModelsParams{
			Provider: p.provider.String(),
		},
	})
	if err != nil {
		return nil, err
	}

	body, ok := resp.HostToPluginResponseUnionUnion.(*rpc.RPCResponseRobotModelProviderListModels)
	if !ok {
		return nil, fmt.Errorf("unexpected robot model provider list models response: %T", resp.HostToPluginResponseUnionUnion)
	}

	now := time.Now()
	models := make([]model_ref.Info, 0, len(body.Models))
	for _, m := range body.Models {
		models = append(models, model_ref.Info{
			Ref: model_ref.ModelRef{
				Provider: p.provider,
				Model:    model_ref.NewModel(m.Name),
			},
			Raw:        modelRaw(m),
			LastSeenAt: now,
		})
	}

	return models, nil
}

func modelRaw(m rpc.RobotModelProviderModel) map[string]any {
	raw := map[string]any{}
	for k, v := range m.Raw {
		raw[k] = v
	}
	if display, ok := m.DisplayName.Get(); ok {
		raw["display_name"] = display
	}
	if contextWindow, ok := m.ContextWindow.Get(); ok {
		raw["context_window"] = contextWindow
	}
	if maxOutputTokens, ok := m.MaxOutputTokens.Get(); ok {
		raw["max_output_tokens"] = maxOutputTokens
	}
	if releasedAt, ok := m.ReleasedAt.Get(); ok {
		raw["released_at"] = releasedAt.Format(time.RFC3339)
	}
	return raw
}

func (p *Provider) GetADKModelLLM(ctx context.Context, ref model_ref.ModelRef) (model.LLM, error) {
	if ref.Provider != p.provider {
		return nil, fmt.Errorf("model ref provider %q does not match plugin provider %q", ref.Provider, p.provider)
	}

	return &pluginModel{
		provider: p.provider,
		model:    ref.Model,
		session:  p.session,
	}, nil
}

func (p *Provider) ValidateModel(ctx context.Context, ref model_ref.ModelRef) error {
	if ref.Provider != p.provider {
		return fmt.Errorf("model ref provider %q does not match plugin provider %q", ref.Provider, p.provider)
	}

	models, err := p.ListModels(ctx)
	if err != nil {
		return err
	}

	for _, model := range models {
		if model.Ref.Model == ref.Model {
			return nil
		}
	}

	return fmt.Errorf("model %q is not available for plugin provider %s", ref.Model, ref.Provider)
}

type pluginModel struct {
	provider model_ref.Provider
	model    model_ref.Model
	session  plugin_runner.Session
}

func (m *pluginModel) Name() string {
	return model_ref.ModelRef{Provider: m.provider, Model: m.model}.String()
}

func (m *pluginModel) GenerateContent(ctx context.Context, req *model.LLMRequest, _ bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		params := convertLLMRequest(m.provider, m.model, req)
		id := xid.New()
		resp, err := m.session.Send(ctx, id, &rpc.RPCRequestRobotModelProviderGenerate{
			ID:      id,
			Jsonrpc: "2.0",
			Method:  "robot_model_provider_generate",
			Params:  params,
		})
		if err != nil {
			yield(nil, err)
			return
		}

		body, ok := resp.HostToPluginResponseUnionUnion.(*rpc.RPCResponseRobotModelProviderGenerate)
		if !ok {
			yield(nil, fmt.Errorf("unexpected robot model provider generate response: %T", resp.HostToPluginResponseUnionUnion))
			return
		}
		if msg, ok := body.Error.Get(); ok && msg != "" {
			yield(nil, fmt.Errorf("plugin model provider returned error: %s", msg))
			return
		}

		yield(convertPluginResponse(body), nil)
	}
}

func convertLLMRequest(provider model_ref.Provider, modelName model_ref.Model, req *model.LLMRequest) rpc.RPCRequestRobotModelProviderGenerateParams {
	out := rpc.RPCRequestRobotModelProviderGenerateParams{
		Provider: provider.String(),
		Model:    modelName.String(),
		Messages: convertContents(req.Contents),
		Tools:    convertTools(req.Config),
	}

	if req.Config != nil && req.Config.SystemInstruction != nil {
		system := contentText(req.Config.SystemInstruction)
		if system != "" {
			out.System = opt.New(system)
		}
	}

	return out
}

func convertContents(contents []*genai.Content) []rpc.RobotModelProviderMessage {
	out := make([]rpc.RobotModelProviderMessage, 0, len(contents))
	for _, content := range contents {
		if content == nil {
			continue
		}

		role := "user"
		if content.Role == string(genai.RoleModel) {
			role = "assistant"
		}

		if msg := contentText(content); msg != "" {
			out = append(out, rpc.RobotModelProviderMessage{
				Role:    role,
				Content: opt.New(msg),
			})
		}

		for _, part := range content.Parts {
			if part == nil {
				continue
			}
			if part.FunctionResponse != nil {
				out = append(out, rpc.RobotModelProviderMessage{
					Role:       "tool",
					Content:    opt.New(mustJSON(part.FunctionResponse.Response)),
					Name:       opt.New(part.FunctionResponse.Name),
					ToolCallID: opt.New(part.FunctionResponse.ID),
				})
			}
		}
	}
	return out
}

func contentText(content *genai.Content) string {
	var parts []string
	for _, part := range content.Parts {
		if part != nil && part.Text != "" {
			parts = append(parts, part.Text)
		}
	}
	return strings.Join(parts, "\n")
}

func convertTools(config *genai.GenerateContentConfig) []rpc.RobotModelProviderTool {
	if config == nil {
		return nil
	}

	var tools []rpc.RobotModelProviderTool
	for _, tool := range config.Tools {
		if tool == nil {
			continue
		}
		for _, declaration := range tool.FunctionDeclarations {
			if declaration == nil {
				continue
			}
			parameters := map[string]any{}
			switch {
			case declaration.ParametersJsonSchema != nil:
				parameters = normaliseObject(declaration.ParametersJsonSchema)
			case declaration.Parameters != nil:
				parameters = normaliseObject(declaration.Parameters)
			}

			t := rpc.RobotModelProviderTool{
				Name:       declaration.Name,
				Parameters: parameters,
			}
			if declaration.Description != "" {
				t.Description = opt.New(declaration.Description)
			}
			tools = append(tools, t)
		}
	}
	return tools
}

func convertPluginResponse(resp *rpc.RPCResponseRobotModelProviderGenerate) *model.LLMResponse {
	parts := []*genai.Part{}

	if content, ok := resp.Content.Get(); ok && content != "" {
		parts = append(parts, genai.NewPartFromText(content))
	}
	for _, call := range resp.ToolCalls {
		part := genai.NewPartFromFunctionCall(call.Name, call.Arguments)
		if id, ok := call.ID.Get(); ok {
			part.FunctionCall.ID = id
		}
		parts = append(parts, part)
	}

	return &model.LLMResponse{
		Content: &genai.Content{
			Role:  genai.RoleModel,
			Parts: parts,
		},
		FinishReason: convertFinishReason(resp.FinishReason.Or("")),
		TurnComplete: true,
	}
}

func convertFinishReason(reason string) genai.FinishReason {
	switch reason {
	case "tool_calls":
		return genai.FinishReasonStop
	case "length":
		return genai.FinishReasonMaxTokens
	case "content_filter":
		return genai.FinishReasonProhibitedContent
	case "error":
		return genai.FinishReasonOther
	default:
		return genai.FinishReasonStop
	}
}

func normaliseObject(v any) map[string]any {
	b, err := json.Marshal(v)
	if err != nil {
		return map[string]any{}
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		return map[string]any{}
	}
	return out
}

func mustJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprint(v)
	}
	return string(b)
}
