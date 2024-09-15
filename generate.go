package storyden

// Generates the OpenAPI stubs for the API server and client.
//go:generate go run -mod=mod github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.3.0 --config api/config.yaml api/openapi.yaml
//go:generate go run github.com/Southclaws/storyden/internal/tools/rbacgen api/openapi.yaml app/transports/http/bindings/openapi_rbac_gen/openapi_rbac_gen.go
