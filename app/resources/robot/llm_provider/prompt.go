package llm_provider

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/invopop/jsonschema"
	"google.golang.org/adk/model"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/robot/model_ref"
)

type PromptOptions struct {
	System string
}

type StructuredPromptRequest struct {
	Description string
	Input       string
	Schema      any
}

type StructuredPrompter interface {
	SupportsStructuredOutput() bool
	StructuredPrompt(ctx context.Context, ref model_ref.ModelRef, request StructuredPromptRequest) (string, error)
}

func (f *Factory) PromptText(ctx context.Context, input string, options ...PromptOptions) (string, error) {
	llm, err := f.defaultLLM(ctx)
	if err != nil {
		return "", err
	}

	req := &model.LLMRequest{
		Contents: []*genai.Content{
			genai.NewContentFromText(input, genai.RoleUser),
		},
	}
	if len(options) > 0 && strings.TrimSpace(options[0].System) != "" {
		req.Config = &genai.GenerateContentConfig{
			SystemInstruction: genai.NewContentFromText(options[0].System, genai.RoleUser),
		}
	}

	return generateText(ctx, llm, req)
}

func PromptObject[T any](ctx context.Context, factory *Factory, description, input string, schema T) (*T, error) {
	serialisedSchema, err := schemaFromObjectInstance(schema)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	defaultModel, err := factory.DefaultModel(ctx)
	if err != nil {
		return nil, err
	}

	if err := factory.EnsureModelAvailable(ctx, defaultModel); err != nil {
		return nil, err
	}

	_, provider, err := factory.provider(ctx, defaultModel.Provider)
	if err != nil {
		return nil, err
	}

	prompter, ok := provider.(StructuredPrompter)
	if !ok || !prompter.SupportsStructuredOutput() {
		return nil, fault.Newf("robot model provider %q does not support structured output", defaultModel.Provider)
	}

	text, err := prompter.StructuredPrompt(ctx, defaultModel, StructuredPromptRequest{
		Description: description,
		Input:       input,
		Schema:      serialisedSchema,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var result T
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &result, nil
}

func (f *Factory) defaultLLM(ctx context.Context) (model.LLM, error) {
	defaultModel, err := f.DefaultModel(ctx)
	if err != nil {
		return nil, err
	}

	return f.GetADKModelLLM(ctx, defaultModel)
}

func generateText(ctx context.Context, llm model.LLM, req *model.LLMRequest) (string, error) {
	var out strings.Builder

	for response, err := range llm.GenerateContent(ctx, req, false) {
		if err != nil {
			return "", fault.Wrap(err, fctx.With(ctx))
		}
		if response == nil || response.Content == nil {
			continue
		}
		for _, part := range response.Content.Parts {
			if part != nil && part.Text != "" {
				out.WriteString(part.Text)
			}
		}
	}

	text := strings.TrimSpace(out.String())
	if text == "" {
		return "", fault.New("LLM response is empty")
	}

	return text, nil
}

func schemaFromObjectInstance[T any](instance T) (any, error) {
	r := jsonschema.Reflector{
		Anonymous:                 true,
		ExpandedStruct:            true,
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}

	schema := r.Reflect(instance)

	return schema, nil
}
