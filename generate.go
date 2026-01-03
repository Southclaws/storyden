package storyden

// Generates the OpenAPI stubs for the API server and client.
//go:generate go run github.com/Southclaws/storyden/internal/tools/configen internal/config/config.yaml internal/config/config.go home/content/docs/operation/configuration.mdx
//go:generate go run github.com/atombender/go-jsonschema@latest -p mcp -o mcp/mcp_schema.go api/robots.yaml
