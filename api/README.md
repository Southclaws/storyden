# Storyden Schemas

This directory contains JSON Schema definitions for:

- OpenAPI: the public API contract for the Storyden API surface
- Robots: MCP style tool definitions for Storyden Robots and MCP integrations
- Common: Shared schemas for common types

## OpenAPI

Storyden makes use of the OpenAPI schema in three ways:

- Backend interface and JSON object structures: oapi-codegen, it's configured using `config.yaml` in this directory
- Frontend client and types: Orval, with server side fetch and client side SWR code generated
- Custom generation: We also generate additional code using custom tooling
