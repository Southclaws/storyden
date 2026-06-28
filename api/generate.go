package api

//go:generate go run -mod=mod github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.3.0 --config config.yaml openapi.yaml

//go:generate go run github.com/Southclaws/schemancer@latest plugin.yaml
//go:generate go run github.com/Southclaws/storyden/internal/tools/eventgen
//go:generate go run github.com/Southclaws/storyden/internal/tools/rpcdocgen
//go:generate go run github.com/Southclaws/storyden/internal/tools/permissiondocgen
//go:generate go run github.com/Southclaws/storyden/internal/tools/rbacgen -schema openapi.yaml -output ../app/transports/http/bindings/openapi_rbac/openapi_rbac_gen.go -enum ../app/resources/rbac/rbac_enum_gen.go -operation-enum ../app/transports/http/openapi/operation/operation_enum_gen.go -operation-cost ../app/transports/http/openapi/operation/operation_cost_gen.go
//go:generate go run github.com/Southclaws/storyden/internal/tools/schemaderef robots.yaml ../lib/mcp/robots.json
//go:generate go run github.com/Southclaws/storyden/internal/tools/schemaderef --root-only plugin.yaml ../lib/plugin/plugin.json
//go:generate go run github.com/Southclaws/storyden/internal/tools/toolbindgen ../lib/mcp/robots.json ../lib/mcp/bindings_gen.go
//go:generate go run github.com/atombender/go-jsonschema@latest -p mcp -o ../lib/mcp/mcp_schema.go robots.yaml
