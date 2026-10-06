package completion

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuggestionsSelectSafeFieldsAndBoundResults(t *testing.T) {
	var body map[string]any
	require.NoError(t, json.Unmarshal([]byte(`{"items":[{"id":"one","name":"Hello\nworld","secret":"PRIVATE","children":[{"id":"two","name":"Nested"}]},{"id":"one","name":"Duplicate"},{"id":"bad\u001bvalue","name":"Bad"},null]}`), &body))
	require.Equal(t, []string{"one", "Hello world", "two", "Nested"}, suggestions(body["items"], source{collection: "items", value: "id", description: "name", children: "children"}))
	var rows []any
	for i := 0; i < 150; i++ {
		rows = append(rows, map[string]any{"id": fmt.Sprintf("id-%d", i)})
	}
	values := suggestions(rows, source{value: "id"})
	require.Len(t, values, 200)
	require.Equal(t, "id-99", values[198])
}

func TestLookupRequiresParentAndUsesGET(t *testing.T) {
	path, ok := lookupPath(source{operation: "NodeVersionList"}, map[string]string{"node_slug": "notes / today"}, "")
	require.True(t, ok)
	require.Equal(t, "/nodes/notes%20%2F%20today/versions", path)
	for _, params := range []map[string]string{nil, {"node_slug": ".."}, {"node_slug": "."}} {
		_, ok := lookupPath(source{operation: "NodeVersionList"}, params, "")
		require.False(t, ok)
	}
	_, ok = lookupPath(source{operation: "TrailRunCreate"}, map[string]string{"trail_id": "trail"}, "")
	require.False(t, ok)
}
