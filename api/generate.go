package api

//go:generate go run -mod=mod github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.3.0 --config config.yaml openapi.yaml
//go:generate go run github.com/Southclaws/schemancer@latest plugin.yaml
//go:generate go run github.com/Southclaws/storyden/internal/tools/eventgen
//go:generate go run github.com/Southclaws/storyden/internal/tools/rpcdocgen
//go:generate go run github.com/Southclaws/storyden/internal/tools/permissiondocgen
