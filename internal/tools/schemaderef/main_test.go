package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRootOnlyDereferencesReachableFileRefsAndDropsDefs(t *testing.T) {
	dir := t.TempDir()

	require.NoError(t, os.WriteFile(filepath.Join(dir, "root.yaml"), []byte(`$schema: "http://json-schema.org/draft-07/schema#"
title: Root
type: object
additionalProperties: false
properties:
  child:
    $ref: "child.yaml#/$defs/Child"
$defs:
  Unused:
    $ref: "unused.yaml#/$defs/Unused"
`), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "child.yaml"), []byte(`$schema: "http://json-schema.org/draft-07/schema#"
$defs:
  Child:
    type: object
    properties:
      value:
        $ref: "value.yaml#/$defs/Value"
  ChildUnused:
    type: string
`), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "value.yaml"), []byte(`$schema: "http://json-schema.org/draft-07/schema#"
$defs:
  Value:
    type: string
    enum:
      - one
      - two
`), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "unused.yaml"), []byte(`$schema: "http://json-schema.org/draft-07/schema#"
$defs:
  Unused:
    type: number
`), 0o644))

	r := resolver{
		cache: map[string]map[string]any{},
		stack: map[string]struct{}{},
	}

	root, err := r.loadDocument(filepath.Join(dir, "root.yaml"))
	require.NoError(t, err)

	out, err := r.dereferenceRoot(root, filepath.Join(dir, "root.yaml"))
	require.NoError(t, err)
	require.NoError(t, stripNestedDefinitions(out, false, true))

	encoded, err := json.Marshal(out)
	require.NoError(t, err)
	require.NotContains(t, string(encoded), "$ref")
	require.NotContains(t, string(encoded), "$defs")
	require.NotContains(t, string(encoded), "Unused")
	require.Contains(t, string(encoded), `"enum":["one","two"]`)
}
