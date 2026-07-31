package update

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
)

func TestReadJSONPropsFromStdin(t *testing.T) {
	r := require.New(t)

	props, err := readJSONProps("-", strings.NewReader(`{
		"name": "JSON Name",
		"hide_child_tree": true,
		"url": null,
		"meta": {"source": "test"}
	}`))
	r.NoError(err)
	r.NotNil(props.Name)
	r.Equal("JSON Name", *props.Name)
	r.NotNil(props.HideChildTree)
	r.True(*props.HideChildTree)
	r.True(props.Url.IsNull())
	r.NotNil(props.Meta)
	r.Equal("test", (*props.Meta)["source"])
}

func TestReadJSONPropsWithURLValue(t *testing.T) {
	r := require.New(t)

	props, err := readJSONProps("-", strings.NewReader(`{"url":"https://example.com"}`))
	r.NoError(err)
	r.True(props.Url.IsSpecified())
	r.False(props.Url.IsNull())

	url, err := props.Url.Get()
	r.NoError(err)
	r.Equal("https://example.com", url)
}

func TestReadJSONPropsRequiresObject(t *testing.T) {
	r := require.New(t)

	_, err := readJSONProps("-", strings.NewReader(`null`))
	r.ErrorContains(err, "node update JSON must be an object")

	_, err = readJSONProps("-", strings.NewReader(`[]`))
	r.ErrorContains(err, "invalid node update JSON")
}

func TestBuildMutablePropsRejectsJSONWithFieldFlags(t *testing.T) {
	r := require.New(t)

	_, err := buildMutableProps(cligen.IO{}, cligen.NodeUpdateParams{
		Name: "Flag Name",
		Json: "-",
	})
	r.ErrorContains(err, "cannot combine --json with --name")
}
