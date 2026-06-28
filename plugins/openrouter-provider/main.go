package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"time"

	openrouter "github.com/OpenRouterTeam/go-sdk"
	"github.com/OpenRouterTeam/go-sdk/models/components"
	"github.com/OpenRouterTeam/go-sdk/models/operations"
	"github.com/OpenRouterTeam/go-sdk/models/sdkerrors"
	"github.com/OpenRouterTeam/go-sdk/optionalnullable"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/lib/plugin/rpc"
	"github.com/Southclaws/storyden/sdk/go/storyden"
)

const (
	providerName          = "openrouter"
	apiKeyConfig          = "api_key"
	modelIDsConfig        = "model_ids"
	embeddingModelConfig  = "embedding_model"
	defaultEmbeddingModel = "openai/text-embedding-3-small"
)

var defaultModelIDs = []string{
	"openai/gpt-4o-mini",
	"openai/gpt-4.1-mini",
	"anthropic/claude-3.5-sonnet",
	"anthropic/claude-3.5-haiku",
	"google/gemini-2.0-flash-001",
	"meta-llama/llama-3.3-70b-instruct",
}

type provider struct {
	plugin *storyden.Plugin
	logger *slog.Logger
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	plugin, err := storyden.New(ctx)
	if err != nil {
		exitError(logger, "create plugin", err)
	}
	defer func() {
		if err := plugin.Shutdown(); err != nil && !errors.Is(err, context.Canceled) {
			logger.Warn("plugin shutdown returned error", slog.String("error", err.Error()))
		}
	}()

	p := &provider{plugin: plugin, logger: logger}
	plugin.OnRobotModelProviderListModels(p.listModels)
	plugin.OnRobotModelProviderGenerate(p.generate)
	plugin.OnRobotModelProviderStructuredPrompt(p.structuredPrompt)
	plugin.OnRobotModelProviderEmbedText(p.embedText)

	logger.Info("starting OpenRouter provider plugin")
	if err := plugin.Run(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		exitError(logger, "plugin runtime", err)
	}
}

func (p *provider) listModels(ctx context.Context, req rpc.RPCRequestRobotModelProviderListModelsParams) (rpc.RPCResponseRobotModelProviderListModels, error) {
	if req.Provider != providerName {
		return rpc.RPCResponseRobotModelProviderListModels{}, fmt.Errorf("unsupported provider %q", req.Provider)
	}

	config, err := p.pluginConfig(ctx)
	if err != nil {
		return rpc.RPCResponseRobotModelProviderListModels{}, err
	}

	client, err := p.clientFromConfig(config)
	if err != nil {
		return rpc.RPCResponseRobotModelProviderListModels{}, err
	}

	outputModalities := "text"
	res, err := client.Models.List(ctx, &operations.GetModelsRequest{
		OutputModalities: &outputModalities,
	})
	if err != nil {
		return rpc.RPCResponseRobotModelProviderListModels{}, fmt.Errorf("list OpenRouter models: %w", err)
	}
	if res == nil {
		return rpc.RPCResponseRobotModelProviderListModels{}, fmt.Errorf("OpenRouter returned no model response")
	}

	models := filterModels(res.Data, modelIDs(config))

	return rpc.RPCResponseRobotModelProviderListModels{
		Method: "robot_model_provider_list_models",
		Models: models,
	}, nil
}

func (p *provider) generate(ctx context.Context, req rpc.RPCRequestRobotModelProviderGenerateParams) (rpc.RPCResponseRobotModelProviderGenerate, error) {
	if req.Provider != providerName {
		return rpc.RPCResponseRobotModelProviderGenerate{}, fmt.Errorf("unsupported provider %q", req.Provider)
	}

	client, err := p.client(ctx)
	if err != nil {
		return rpc.RPCResponseRobotModelProviderGenerate{}, err
	}

	messages := convertMessages(req)
	tools := convertTools(req.Tools)
	stream := false

	res, err := client.Chat.Send(ctx, components.ChatRequest{
		Model:    &req.Model,
		Messages: messages,
		Tools:    tools,
		Stream:   &stream,
	}, nil)
	if err != nil {
		return rpc.RPCResponseRobotModelProviderGenerate{}, fmt.Errorf("OpenRouter chat completion failed: %s", openRouterErrorMessage(err))
	}
	if res == nil || res.ChatResult == nil {
		return rpc.RPCResponseRobotModelProviderGenerate{}, fmt.Errorf("OpenRouter returned no chat result")
	}
	if len(res.ChatResult.Choices) == 0 {
		return rpc.RPCResponseRobotModelProviderGenerate{}, fmt.Errorf("OpenRouter returned no choices")
	}

	choice := res.ChatResult.Choices[0]
	message := choice.Message
	response := rpc.RPCResponseRobotModelProviderGenerate{
		Method:    "robot_model_provider_generate",
		ToolCalls: convertToolCalls(message.ToolCalls),
	}

	if choice.FinishReason != nil {
		response.FinishReason = opt.New(string(*choice.FinishReason))
	}
	if content, ok := message.Content.GetOrZero(); ok {
		if text := contentText(content); text != "" {
			response.Content = opt.New(text)
		}
	}

	return response, nil
}

func (p *provider) structuredPrompt(ctx context.Context, req rpc.RPCRequestRobotModelProviderStructuredPromptParams) (rpc.RPCResponseRobotModelProviderStructuredPrompt, error) {
	if req.Provider != providerName {
		return rpc.RPCResponseRobotModelProviderStructuredPrompt{}, fmt.Errorf("unsupported provider %q", req.Provider)
	}

	client, err := p.client(ctx)
	if err != nil {
		return rpc.RPCResponseRobotModelProviderStructuredPrompt{}, err
	}

	stream := false
	strict := true
	responseFormat := components.CreateResponseFormatJSONSchema(components.ChatFormatJSONSchemaConfig{
		JSONSchema: components.ChatJSONSchemaConfig{
			Name:        "structured_output",
			Description: &req.Description,
			Schema:      req.Schema,
			Strict:      optionalnullable.From(&strict),
		},
	})

	res, err := client.Chat.Send(ctx, components.ChatRequest{
		Model:          &req.Model,
		Messages:       []components.ChatMessages{userMessage(req.Input)},
		Stream:         &stream,
		ResponseFormat: &responseFormat,
	}, nil)
	if err != nil {
		return rpc.RPCResponseRobotModelProviderStructuredPrompt{}, fmt.Errorf("OpenRouter structured prompt failed: %s", openRouterErrorMessage(err))
	}
	content, err := chatResponseText(res)
	if err != nil {
		return rpc.RPCResponseRobotModelProviderStructuredPrompt{}, err
	}

	return rpc.RPCResponseRobotModelProviderStructuredPrompt{
		Method:  "robot_model_provider_structured_prompt",
		Content: opt.New(content),
	}, nil
}

func (p *provider) embedText(ctx context.Context, req rpc.RPCRequestRobotModelProviderEmbedTextParams) (rpc.RPCResponseRobotModelProviderEmbedText, error) {
	if req.Provider != providerName {
		return rpc.RPCResponseRobotModelProviderEmbedText{}, fmt.Errorf("unsupported provider %q", req.Provider)
	}

	config, err := p.pluginConfig(ctx)
	if err != nil {
		return rpc.RPCResponseRobotModelProviderEmbedText{}, err
	}

	client, err := p.clientFromConfig(config)
	if err != nil {
		return rpc.RPCResponseRobotModelProviderEmbedText{}, err
	}

	encodingFormat := operations.EncodingFormatFloat
	res, err := client.Embeddings.Generate(ctx, operations.CreateEmbeddingsRequest{
		Model:          embeddingModel(config),
		Input:          operations.CreateInputUnionStr(req.Text),
		EncodingFormat: &encodingFormat,
	})
	if err != nil {
		return rpc.RPCResponseRobotModelProviderEmbedText{}, fmt.Errorf("OpenRouter embedding failed: %s", openRouterErrorMessage(err))
	}
	if res == nil || res.CreateEmbeddingsResponseBody == nil {
		return rpc.RPCResponseRobotModelProviderEmbedText{}, fmt.Errorf("OpenRouter returned no embedding response")
	}
	if len(res.CreateEmbeddingsResponseBody.Data) == 0 {
		return rpc.RPCResponseRobotModelProviderEmbedText{}, fmt.Errorf("OpenRouter returned no embeddings")
	}

	embedding := res.CreateEmbeddingsResponseBody.Data[0].Embedding.ArrayOfNumber
	if len(embedding) == 0 {
		return rpc.RPCResponseRobotModelProviderEmbedText{}, fmt.Errorf("OpenRouter returned an empty embedding")
	}

	return rpc.RPCResponseRobotModelProviderEmbedText{
		Method:    "robot_model_provider_embed_text",
		Embedding: embedding,
	}, nil
}

func (p *provider) client(ctx context.Context) (*openrouter.OpenRouter, error) {
	config, err := p.pluginConfig(ctx)
	if err != nil {
		return nil, err
	}

	return p.clientFromConfig(config)
}

func (p *provider) pluginConfig(ctx context.Context) (map[string]any, error) {
	config, err := p.plugin.GetConfig(ctx, apiKeyConfig, modelIDsConfig, embeddingModelConfig)
	if err != nil {
		return nil, fmt.Errorf("get plugin config: %w", err)
	}

	return config, nil
}

func (p *provider) clientFromConfig(config map[string]any) (*openrouter.OpenRouter, error) {
	apiKey, _ := config[apiKeyConfig].(string)
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, fmt.Errorf("OpenRouter API key is not configured")
	}

	return openrouter.New(openrouter.WithSecurity(apiKey)), nil
}

func modelIDs(config map[string]any) []string {
	raw, _ := config[modelIDsConfig].(string)
	ids := parseModelIDs(raw)
	if len(ids) == 0 {
		return defaultModelIDs
	}
	return ids
}

func embeddingModel(config map[string]any) string {
	raw, _ := config[embeddingModelConfig].(string)
	model := strings.TrimSpace(raw)
	if model == "" {
		return defaultEmbeddingModel
	}
	return model
}

func parseModelIDs(raw string) []string {
	fields := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\r' || r == '\t' || r == ' '
	})

	ids := make([]string, 0, len(fields))
	seen := map[string]struct{}{}
	for _, field := range fields {
		id := strings.TrimSpace(field)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}

	return ids
}

func filterModels(all []components.Model, ids []string) []rpc.RobotModelProviderModel {
	byID := make(map[string]components.Model, len(all))
	for _, model := range all {
		byID[model.ID] = model
	}

	models := make([]rpc.RobotModelProviderModel, 0, len(ids))
	for _, id := range ids {
		model, ok := byID[id]
		if !ok {
			continue
		}
		models = append(models, mapModel(model))
	}

	return models
}

func openRouterErrorMessage(err error) string {
	var paymentRequired *sdkerrors.PaymentRequiredResponseError
	if errors.As(err, &paymentRequired) {
		return compactOpenRouterMessage("payment required", paymentRequired.Error_.Message)
	}

	var badRequest *sdkerrors.BadRequestResponseError
	if errors.As(err, &badRequest) {
		return compactOpenRouterMessage("bad request", badRequest.Error_.Message)
	}

	var unauthorized *sdkerrors.UnauthorizedResponseError
	if errors.As(err, &unauthorized) {
		return compactOpenRouterMessage("unauthorized", unauthorized.Error_.Message)
	}

	var tooManyRequests *sdkerrors.TooManyRequestsResponseError
	if errors.As(err, &tooManyRequests) {
		return compactOpenRouterMessage("rate limited", tooManyRequests.Error_.Message)
	}

	var apiError *sdkerrors.APIError
	if errors.As(err, &apiError) {
		if msg := openRouterAPIErrorBodyMessage(apiError.Body); msg != "" {
			return compactOpenRouterMessage(apiError.Message, msg)
		}
		return compactOpenRouterMessage(apiError.Message, "")
	}

	return compactOpenRouterMessage("", err.Error())
}

func openRouterAPIErrorBodyMessage(body string) string {
	var payload struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return ""
	}
	return payload.Error.Message
}

func compactOpenRouterMessage(prefix, message string) string {
	message = strings.TrimSpace(message)
	if message == "" {
		message = "request failed"
	}

	if idx := strings.Index(message, "https://openrouter.ai/"); idx >= 0 {
		message = strings.TrimSpace(message[:idx])
	}
	message = strings.TrimRight(message, " .")
	if len(message) > 240 {
		message = strings.TrimSpace(message[:240]) + "..."
	}

	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return message
	}

	return fmt.Sprintf("%s: %s", prefix, message)
}

func chatResponseText(res *operations.SendChatCompletionRequestResponse) (string, error) {
	if res == nil || res.ChatResult == nil {
		return "", fmt.Errorf("OpenRouter returned no chat result")
	}
	if len(res.ChatResult.Choices) == 0 {
		return "", fmt.Errorf("OpenRouter returned no choices")
	}

	choice := res.ChatResult.Choices[0]
	content, ok := choice.Message.Content.GetOrZero()
	if !ok {
		return "", fmt.Errorf("OpenRouter returned no message content")
	}

	text := strings.TrimSpace(contentText(content))
	if text == "" {
		return "", fmt.Errorf("OpenRouter returned empty message content")
	}

	return text, nil
}

func userMessage(content string) components.ChatMessages {
	return components.CreateChatMessagesUser(components.ChatUserMessage{
		Content: components.CreateChatUserMessageContentStr(content),
	})
}

func mapModel(model components.Model) rpc.RobotModelProviderModel {
	out := rpc.RobotModelProviderModel{
		Name:        model.ID,
		DisplayName: opt.New(model.Name),
	}

	if model.ContextLength != nil {
		out.ContextWindow = opt.New(int(*model.ContextLength))
	}
	if value, ok := model.TopProvider.MaxCompletionTokens.GetOrZero(); ok && value > 0 {
		out.MaxOutputTokens = opt.New(int(value))
	} else if model.PerRequestLimits != nil && model.PerRequestLimits.CompletionTokens > 0 {
		out.MaxOutputTokens = opt.New(int(model.PerRequestLimits.CompletionTokens))
	}
	if model.Created > 0 {
		out.ReleasedAt = opt.New(time.Unix(model.Created, 0))
	}

	return out
}

func convertMessages(req rpc.RPCRequestRobotModelProviderGenerateParams) []components.ChatMessages {
	messages := []components.ChatMessages{}

	if system, ok := req.System.Get(); ok && system != "" {
		messages = append(messages, components.CreateChatMessagesSystem(components.ChatSystemMessage{
			Content: components.CreateChatSystemMessageContentStr(system),
		}))
	}

	for _, msg := range req.Messages {
		content := msg.Content.Or("")
		switch msg.Role {
		case "system":
			messages = append(messages, components.CreateChatMessagesSystem(components.ChatSystemMessage{
				Content: components.CreateChatSystemMessageContentStr(content),
			}))
		case "assistant":
			value := components.CreateChatAssistantMessageContentStr(content)
			messages = append(messages, components.CreateChatMessagesAssistant(components.ChatAssistantMessage{
				Content: optionalnullable.From(&value),
			}))
		case "tool":
			toolCallID := msg.ToolCallID.Or("")
			messages = append(messages, components.CreateChatMessagesTool(components.ChatToolMessage{
				Content:    components.CreateChatToolMessageContentStr(content),
				ToolCallID: toolCallID,
			}))
		default:
			messages = append(messages, userMessage(content))
		}
	}

	return messages
}

func convertTools(tools []rpc.RobotModelProviderTool) []components.ChatFunctionTool {
	out := make([]components.ChatFunctionTool, 0, len(tools))
	for _, tool := range tools {
		description := tool.Description.Or("")
		out = append(out, components.CreateChatFunctionToolChatFunctionToolFunction(components.ChatFunctionToolFunction{
			Type: components.ChatFunctionToolTypeFunction,
			Function: components.ChatFunctionToolFunctionFunction{
				Description: &description,
				Name:        tool.Name,
				Parameters:  tool.Parameters,
			},
		}))
	}
	return out
}

func convertToolCalls(calls []components.ChatToolCall) []rpc.RobotModelProviderToolCall {
	out := make([]rpc.RobotModelProviderToolCall, 0, len(calls))
	for _, call := range calls {
		out = append(out, rpc.RobotModelProviderToolCall{
			ID:        opt.New(call.ID),
			Name:      call.Function.Name,
			Arguments: parseArguments(call.Function.Arguments),
		})
	}
	return out
}

func contentText(content components.ChatAssistantMessageContent) string {
	if content.Str != nil {
		return *content.Str
	}

	if content.Any != nil {
		return fmt.Sprint(content.Any)
	}

	if len(content.ArrayOfChatContentItems) == 0 {
		return ""
	}

	parts := make([]string, 0, len(content.ArrayOfChatContentItems))
	for _, item := range content.ArrayOfChatContentItems {
		if item.ChatContentText != nil {
			parts = append(parts, item.ChatContentText.Text)
		}
	}
	return strings.Join(parts, "\n")
}

func parseArguments(raw string) map[string]any {
	if strings.TrimSpace(raw) == "" {
		return map[string]any{}
	}

	var out map[string]any
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return map[string]any{"raw": raw}
	}
	return out
}

func exitError(logger *slog.Logger, message string, err error) {
	logger.Error(message, slog.String("error", err.Error()))
	os.Exit(1)
}
