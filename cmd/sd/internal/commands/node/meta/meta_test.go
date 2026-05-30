package meta

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadMetadata(t *testing.T) {
	r := require.New(t)

	metadata, err := readMetadata(`{"source":"cli","count":2}`, "", strings.NewReader(""))
	r.NoError(err)
	r.Equal("cli", metadata["source"])
	r.Equal(float64(2), metadata["count"])
}

func TestReadMetadataFromStdin(t *testing.T) {
	r := require.New(t)

	metadata, err := readMetadata("", "-", strings.NewReader(`{"from":"stdin"}`))
	r.NoError(err)
	r.Equal("stdin", metadata["from"])
}

func TestReadMetadataRequiresObject(t *testing.T) {
	r := require.New(t)

	_, err := readMetadata(`null`, "", strings.NewReader(""))
	r.ErrorContains(err, "metadata must be a JSON object")

	_, err = readMetadata(`[]`, "", strings.NewReader(""))
	r.ErrorContains(err, "invalid metadata JSON")
}
