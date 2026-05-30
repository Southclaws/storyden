package children

import (
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
