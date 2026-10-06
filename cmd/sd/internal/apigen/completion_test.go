package main

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
)

func TestCompletionChoicesFromSchemas(t *testing.T) {
	for _, tc := range []struct {
		schema string
		want   []string
	}{
		{`{"type":"string","enum":["draft","published"]}`, []string{"draft", "published"}},
		{`{"type":"array","items":{"type":"string","enum":["one","two"]}}`, []string{"one", "two"}},
		{`{"type":"boolean"}`, []string{"true", "false"}},
		{`{"oneOf":[{"type":"string","enum":["sse"]},{"type":"null"}]}`, []string{"sse"}},
		{`{"type":"integer"}`, nil},
	} {
		var schema openapi3.Schema
		require.NoError(t, schema.UnmarshalJSON([]byte(tc.schema)))
		require.Equal(t, tc.want, schemaChoices(&openapi3.SchemaRef{Value: &schema}))
	}
}

func TestSharedCompletionArgumentsResolve(t *testing.T) {
	commands := []command{{Commands: []command{{Arguments: []parameter{{Ref: "#/components/arguments/page"}}}}}}
	require.NoError(t, resolveArguments(commands, map[string]parameter{"page": {Name: "SLUG", Variadic: true}}))
	require.Equal(t, "SLUG", commands[0].Commands[0].Arguments[0].Name)
	require.True(t, commands[0].Commands[0].Arguments[0].Variadic)
	require.ErrorContains(t, resolveArguments([]command{{Arguments: []parameter{{Ref: "#/components/arguments/missing"}}}}, nil), "unknown argument reference")
}
