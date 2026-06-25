# mcp

Contains tool definition loading code for Robots.

Note: kinda not really an actual MCP lol, it _is_ used also in some way in the mcp server itself, but it's also used internally for agents mostly as well as frontend codegen.

## Files

- robots.json is a de-referenced schema from ../api/robots.yaml
- embed.go loads the json into go's memory space at build time
- bindings.go wires together the actual schema with ToolDefinition instances
- mcp_schema.go is generated code (structs) from the schema itself

### Why dereferenced?

The actual schema source of truth, the file you want to edit when you change tools, is in `../api/robots.yaml`.

This file uses JSONSchema references to other files, so we can keep it clean and share schema definitions with the OpenAPI schema, less repetition for more complex types etc.

But, the schema loading tools are a bit awkward with file $refs so we run a super basic dereference process (see ./api/generate.go) to remove all the $refs and replace them with the actual content. This yields an awful huge .json version of the schema but it's only used to load into \*jsonschema.Schema objects so it's fine.

Rule is:

- edit api/robots.yaml
- generate code to get:
  - robots.json
  - go types
  - typescript types

The frontend does something similar too, for the same reason: $refs to files are awkward and a lot of tooling just doesn't deal with it well.

## How it's used

You may be wondering, why generate code AND load the schema JSON into memory?

Google's ADK library accepts a \*jsonschema.Schema as part of its tool definition:

```go
return functiontool.New(
	functiontool.Config{
		Name:        toolDef.Name,
		Description: toolDef.Description, // <- here
		InputSchema: toolDef.InputSchema, // <- here
```

These are actual Go objects that hold a representation of the schema.

But the actual function call itself doesn't use a schema, it receives a deserialised structure. ADK actually takes the generic type of the argument and deserialises the tool output into that type:

```go
func (st *searchTools) ExecuteSearch(ctx tool.Context, args mcp.ToolSearchInput) (*mcp.ToolSearchOutput, error) {
```

It uses a similar trick for the output too.

So this requires keeping tabs on two sources of truth:

1. the \*jsonschema.Schema{} instance that defines the schema sent to the LLM
2. the actual Go struct in our code that schema is representing

So, in order to not trip up forgetting to update either, we simply generate the code for 2. and load the JSON in and parse it for 1.

Because the structs are generated from the JSON, it's guaranteed they are the same and we won't pass in a schema that doesn't match the struct.

And, as a last step to ensure it's all mapped together properly, the `bindings.go` file ensures we keep the pairing of each schema to each set of structs in the same folder as the schema, so they can change together. `initTool` links them together, returns a `ToolDefinition` and the actual Robots service imports those ToolDefinitions and uses them in its actual tool function setup.
