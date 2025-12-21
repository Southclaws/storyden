package storyden

// Generates the OpenAPI stubs for the API server and client.
//go:generate go run -mod=mod github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.3.0 --config api/config.yaml api/openapi.yaml
//go:generate go run github.com/Southclaws/storyden/internal/tools/rbacgen api/openapi.yaml app/transports/http/bindings/openapi_rbac_gen/openapi_rbac_gen.go
//go:generate go run github.com/Southclaws/storyden/internal/tools/ratelimitgen api/openapi.yaml app/transports/http/middleware/limiter/ratelimit_config_gen.go
//go:generate go run github.com/Southclaws/storyden/internal/tools/configen internal/config/config.yaml internal/config/config.go home/content/docs/operation/configuration.mdx
