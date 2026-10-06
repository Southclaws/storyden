package nodeapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/stretchr/testify/require"
)

func TestMergePropertySchemaPreservesFieldIdentityAndUnmentionedFields(t *testing.T) {
	existing := []openapi.PropertySchema{
		{Fid: "status-id", Name: "status", Type: openapi.PropertyTypeText, Sort: "asc"},
		{Fid: "priority-id", Name: "priority", Type: openapi.PropertyTypeNumber, Sort: "desc"},
	}
	got, err := mergePropertySchema(existing, []openapi.PropertySchemaMutableProps{
		{Name: "status", Type: openapi.PropertyTypeText, Sort: "desc"},
		{Name: "approved", Type: openapi.PropertyTypeBoolean, Sort: "asc"},
	})
	require.NoError(t, err)
	statusID, priorityID := "status-id", "priority-id"
	require.Equal(t, []openapi.PropertySchemaMutableProps{
		{Fid: &statusID, Name: "status", Type: openapi.PropertyTypeText, Sort: "desc"},
		{Fid: &priorityID, Name: "priority", Type: openapi.PropertyTypeNumber, Sort: "desc"},
		{Name: "approved", Type: openapi.PropertyTypeBoolean, Sort: "asc"},
	}, got)
	require.Equal(t, "asc", existing[0].Sort)
}

func TestMergePropertySchemaRejectsDuplicateNames(t *testing.T) {
	_, err := mergePropertySchema(nil, []openapi.PropertySchemaMutableProps{{Name: "status"}, {Name: "status"}})
	require.ErrorContains(t, err, `duplicate schema field "status"`)
}

func TestPreparePropertySchemaRejectsEmptyParent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"slug":"empty","children":[],"child_property_schema":[]}`))
	}))
	defer server.Close()
	client, err := openapi.NewClientWithResponses(server.URL)
	require.NoError(t, err)
	_, err = PreparePropertySchema(context.Background(), client, "empty", nil)
	require.ErrorContains(t, err, "has no children")
}
