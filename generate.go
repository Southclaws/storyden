package storyden

// Generates the OpenAPI stubs for the API server and client.
//go:generate go run -mod=mod github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.3.0 --config api/config.yaml api/openapi.yaml
//go:generate go run github.com/Southclaws/storyden/internal/tools/rbacgen -schema api/openapi.yaml -output app/transports/http/bindings/openapi_rbac/openapi_rbac_gen.go -enum app/resources/rbac/rbac_enum_gen.go -operation-enum app/transports/http/openapi/operation/operation_enum_gen.go -operation-cost app/transports/http/openapi/operation/operation_cost_gen.go
//go:generate go run github.com/Southclaws/storyden/internal/tools/configen internal/config/config.yaml internal/config/config.go home/content/docs/operation/configuration.mdx
//go:generate go run github.com/Southclaws/storyden/internal/tools/rpcdocgen -schema api/rpc -out home/content/docs/extending/rpc
