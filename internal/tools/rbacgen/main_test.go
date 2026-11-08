package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunPermissionNames(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	schema := filepath.Join(dir, "openapi.yaml")
	err := os.WriteFile(schema, []byte(`openapi: 3.0.3
info:
  title: Permission generation test
  version: "1.0"
paths:
  /example:
    get:
      operationId: Example
      responses:
        "200":
          description: OK
components:
  schemas:
    Permission:
      type: string
      enum:
        - CREATE_POST
        - READ_PUBLISHED_THREADS
        - USE_OAUTH_CLIENTS
        - HTTP2_API
        - "mixed-Case.value"
        - "école/東京"
`), 0o600)
	require.NoError(t, err)

	permissionsFile := filepath.Join(dir, "permissions.go")
	bindingsFile := filepath.Join(dir, "bindings.go")
	operationsFile := filepath.Join(dir, "operations.go")
	costsFile := filepath.Join(dir, "costs.go")
	listFile := filepath.Join(dir, "permission_list_gen.go")

	require.NoError(t, run(schema, bindingsFile, permissionsFile, operationsFile, costsFile))

	for _, path := range []string{bindingsFile, permissionsFile, operationsFile, costsFile, listFile} {
		_, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.AllErrors)
		require.NoError(t, err, "generated file %s must be valid Go", path)
	}

	file, err := parser.ParseFile(token.NewFileSet(), permissionsFile, nil, 0)
	require.NoError(t, err)

	permissions := map[string]string{}
	ast.Inspect(file, func(node ast.Node) bool {
		declaration, ok := node.(*ast.ValueSpec)
		if !ok {
			return true
		}

		if len(declaration.Values) != 1 {
			return true
		}

		literal, ok := declaration.Values[0].(*ast.CompositeLit)
		if !ok {
			return true
		}

		typeName, ok := literal.Type.(*ast.Ident)
		if !ok || typeName.Name != "Permission" {
			return true
		}

		require.Len(t, declaration.Names, 1)
		require.Len(t, literal.Elts, 1)

		value, ok := literal.Elts[0].(*ast.BasicLit)
		require.True(t, ok, "permission value must be a string literal")

		text, err := strconv.Unquote(value.Value)
		require.NoError(t, err)

		permissions[declaration.Names[0].Name] = text

		return false
	})

	require.Equal(t, map[string]string{
		"PermissionCreatePost":           "CREATE_POST",
		"PermissionReadPublishedThreads": "READ_PUBLISHED_THREADS",
		"PermissionUseOauthClients":      "USE_OAUTH_CLIENTS",
		"PermissionHttp2Api":             "HTTP2_API",
		"PermissionMixedCaseValue":       "mixed-Case.value",
		"PermissionÉcole東京":              "école/東京",
	}, permissions)
}
