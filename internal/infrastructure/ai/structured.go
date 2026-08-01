package ai

import (
	"context"
	"encoding/json"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
)

func PromptObject[T any](ctx context.Context, prompter Prompter, description, input string, schema T) (*T, error) {
	s, ok := prompter.(*OpenAI)
	if !ok {
		return nil, fault.New("structured prompt only supported with OpenAI prompter")
	}

	serialisedSchema, err := schemaFromObjectInstance(schema)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	var schemaMap map[string]any
	encodedSchema, err := json.Marshal(serialisedSchema)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if err := json.Unmarshal(encodedSchema, &schemaMap); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	res, err := s.client.Responses.New(ctx, responses.ResponseNewParams{
		Model: openai.ChatModelGPT4_1,
		Store: param.NewOpt(false),
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(input),
		},
		Text: responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfJSONSchema: &responses.ResponseFormatTextJSONSchemaConfigParam{
					Name:        "json_schema",
					Strict:      param.NewOpt(true),
					Description: param.NewOpt(description),
					Schema:      schemaMap,
				},
			},
		},
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if res.Status != responses.ResponseStatusCompleted {
		return nil, fault.Wrap(openAIResponseError(*res), fctx.With(ctx))
	}
	payload := res.OutputText()
	if payload == "" {
		return nil, fault.New("result json is empty")
	}

	var result T
	err = json.Unmarshal([]byte(payload), &result)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &result, nil
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
