package children

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func TestParseSchemaOmitsSyntheticFieldIDs(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	schema, err := parseSchema([]string{"status:text:asc", "priority:number:desc"})
	r.NoError(err)
	r.Len(schema, 2)

	a.Nil(schema[0].Fid)
	a.Equal("status", schema[0].Name)
	a.Equal(openapi.PropertyTypeText, schema[0].Type)
	a.Equal("asc", schema[0].Sort)

	a.Nil(schema[1].Fid)
	a.Equal("priority", schema[1].Name)
	a.Equal(openapi.PropertyTypeNumber, schema[1].Type)
	a.Equal("desc", schema[1].Sort)
}

func TestParseSchemaValidatesSort(t *testing.T) {
	r := require.New(t)

	schema, err := parseSchema([]string{"status:text:ASC"})
	r.NoError(err)
	r.Equal("asc", schema[0].Sort)

	_, err = parseSchema([]string{"status:text:sideways"})
	r.ErrorContains(err, "invalid sort")
}

func TestSetChildrenSchemaPreservesValuesByKeepingFieldIDs(t *testing.T) {
	var patched bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method + " " + r.URL.Path {
		case "GET /nodes/parent":
			fmt.Fprint(w, `{"children":[{"slug":"child"}],"child_property_schema":[{"fid":"status-id","name":"status","type":"text","sort":"asc"},{"fid":"priority-id","name":"priority","type":"number","sort":"asc"}]}`)
		case "PATCH /nodes/parent/children/property-schema":
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.JSONEq(t, `[{"fid":"status-id","name":"status","type":"text","sort":"desc"},{"fid":"priority-id","name":"priority","type":"number","sort":"asc"}]`, string(body))
			patched = true
			fmt.Fprint(w, `{"properties":[]}`)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	client, err := openapi.NewClientWithResponses(server.URL)
	require.NoError(t, err)
	_, err = setChildrenSchema(context.Background(), client, "parent", []openapi.PropertySchemaMutableProps{{Name: "status", Type: openapi.PropertyTypeText, Sort: "desc"}})
	require.NoError(t, err)
	require.True(t, patched)
}
