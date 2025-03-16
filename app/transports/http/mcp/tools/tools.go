package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/api"
	"github.com/Southclaws/storyden/app/transports/http/mcp/mcp_schema"
)

type Provider struct {
	logger *zap.Logger
	doc    v3.Document
}

func New(logger *zap.Logger) (*Provider, error) {
	spec := api.GetOpenAPISpec()

	document, err := libopenapi.NewDocument(spec)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	docModel, errs := document.BuildV3Model()
	if len(errs) > 0 {
		var errMsgs []string
		for _, e := range errs {
			errMsgs = append(errMsgs, e.Error())
		}
		return nil, fault.New(fmt.Sprintf("errors building OpenAPI model: %s", strings.Join(errMsgs, ", ")))
	}

	return &Provider{
		logger: logger,
		doc:    docModel.Model,
	}, nil
}

func (p *Provider) ListTools(ctx context.Context, req mcp_schema.ListToolsRequest) (mcp_schema.ListToolsResult, error) {
	var tools []mcp_schema.Tool

	for pathPair := p.doc.Paths.PathItems.First(); pathPair != nil; pathPair = pathPair.Next() {
		path := pathPair.Key()
		pathItem := pathPair.Value()

		operations := map[string]*v3.Operation{
			"GET":    pathItem.Get,
			"POST":   pathItem.Post,
			"PUT":    pathItem.Put,
			"DELETE": pathItem.Delete,
			"PATCH":  pathItem.Patch,
		}

		for method, operation := range operations {
			if operation == nil {
				continue
			}

			description := getDescription(operation)

			properties := make(mcp_schema.ToolInputSchemaProperties)

			properties["path"] = map[string]interface{}{
				"type":        "string",
				"description": fmt.Sprintf("The path for the API request: %s", path),
				"const":       path,
			}

			properties["method"] = map[string]interface{}{
				"type":        "string",
				"description": fmt.Sprintf("The HTTP method: %s", method),
				"enum":        []interface{}{method},
			}

			required := []string{"path", "method"}

			pathParams := make(map[string]interface{})
			queryParams := make(map[string]interface{})

			if len(operation.Parameters) > 0 {
				for _, paramRef := range operation.Parameters {
					if paramRef == nil {
						continue
					}

					paramSchema := map[string]interface{}{
						"type":        "string",
						"description": paramRef.Description,
					}

					if paramRef.Schema != nil {
						schema := paramRef.Schema.Schema()
						if schema.Type != nil && len(schema.Type) > 0 {
							paramSchema["type"] = schema.Type[0]
						}
						if schema.Format != "" {
							paramSchema["format"] = schema.Format
						}
						if schema.Enum != nil {
							paramSchema["enum"] = schema.Enum
						}
					}

					paramName := paramRef.Name

					switch paramRef.In {
					case "path":
						pathParams[paramName] = paramSchema
						if paramRef.Required != nil && *paramRef.Required {
							required = append(required, "path_params")
						}
					case "query":
						queryParams[paramName] = paramSchema
						if paramRef.Required != nil && *paramRef.Required {
							required = append(required, "query_params")
						}
					}
				}

				if len(pathParams) > 0 {
					properties["path_params"] = map[string]interface{}{
						"type":        "object",
						"description": "URL path parameters",
						"properties":  pathParams,
					}
				}

				if len(queryParams) > 0 {
					properties["query_params"] = map[string]interface{}{
						"type":        "object",
						"description": "URL query parameters",
						"properties":  queryParams,
					}
				}
			}

			if operation.RequestBody != nil {
				for contentPair := operation.RequestBody.Content.First(); contentPair != nil; contentPair = contentPair.Next() {
					contentType := contentPair.Key()
					mediaType := contentPair.Value()

					if mediaType.Schema != nil {
						schema := mediaType.Schema.Schema()
						bodySchema := map[string]interface{}{
							"type":        "object",
							"description": fmt.Sprintf("The request body (Content-Type: %s)", contentType),
						}

						if schema.Properties != nil && schema.Properties.Len() > 0 {
							props := make(map[string]interface{})
							for pair := schema.Properties.First(); pair != nil; pair = pair.Next() {
								propName := pair.Key()
								propSchema := pair.Value().Schema()

								prop := map[string]interface{}{
									"type":        propSchema.Type[0],
									"description": propSchema.Description,
								}
								if propSchema.Format != "" {
									prop["format"] = propSchema.Format
								}
								if propSchema.Enum != nil {
									prop["enum"] = propSchema.Enum
								}
								props[propName] = prop
							}
							bodySchema["properties"] = props
						}

						properties["body"] = bodySchema

						if operation.RequestBody.Required != nil && *operation.RequestBody.Required {
							required = append(required, "body")
						}

						break
					}
				}
			}

			tool := mcp_schema.Tool{
				Name:        operation.OperationId,
				Description: &description,
				InputSchema: mcp_schema.ToolInputSchema{
					Type:       "object",
					Properties: properties,
					Required:   required,
				},
			}

			tools = append(tools, tool)
		}
	}

	p.logger.Info("Generated tools from OpenAPI spec", zap.Int("count", len(tools)))

	return mcp_schema.ListToolsResult{
		Tools: tools,
	}, nil
}

func getDescription(op *v3.Operation) string {
	if op.Description != "" {
		return op.Description
	}
	if op.Summary != "" {
		return op.Summary
	}
	return fmt.Sprintf("Operation: %s", op.OperationId)
}
