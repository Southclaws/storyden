package set

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func TestPropertiesFetchErrorIncludesHTTPContext(t *testing.T) {
	r := require.New(t)

	err := propertiesFetchError(&openapi.NodeGetResponse{
		Body: []byte("missing node"),
		HTTPResponse: &http.Response{
			StatusCode: http.StatusNotFound,
			Status:     "404 Not Found",
		},
	})

	r.ErrorContains(err, "failed to fetch node properties: 404 Not Found: missing node")
}

func TestParsePropertiesRejectsEmptyExplicitType(t *testing.T) {
	r := require.New(t)

	_, err := parseProperties([]string{"status:=draft"}, nil)

	r.ErrorContains(err, "property type cannot be empty")
}
