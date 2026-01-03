package api

//go:generate go run -mod=mod github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.3.0 --config config.yaml openapi.yaml
//go:generate go run github.com/Southclaws/storyden/internal/tools/rbacgen openapi.yaml ../app/transports/http/bindings/openapi_rbac_gen/openapi_rbac_gen.go
//go:generate go run github.com/Southclaws/storyden/internal/tools/schemaderef robots.yaml ../mcp/robots.json
//go:generate go run github.com/atombender/go-jsonschema@latest -p mcp -o ../mcp/mcp_schema.go robots.yaml
