package batch

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadIdentifiersPlainLines(t *testing.T) {
	r := require.New(t)
	ids, err := ReadIdentifiers(strings.NewReader("a\nb\n\nc\n"))
	r.NoError(err)
	r.Equal([]string{"a", "b", "c"}, ids)
}

func TestReadIdentifiersJSONLPrefersSlug(t *testing.T) {
	r := require.New(t)
	in := `{"slug":"the-slug","id":"the-id"}` + "\n" +
		`{"id":"d6o72ret00cs2a2a4ll0"}` + "\n" +
		`{"name":"no-id-here"}` + "\n"
	ids, err := ReadIdentifiers(strings.NewReader(in))
	r.NoError(err)
	r.Equal([]string{"the-slug", "d6o72ret00cs2a2a4ll0"}, ids)
}

func TestReadIdentifiersInvalidJSON(t *testing.T) {
	r := require.New(t)
	_, err := ReadIdentifiers(strings.NewReader(`{not valid json`))
	r.Error(err)
}

func TestReadIdentifiersMixedPlainAndJSON(t *testing.T) {
	r := require.New(t)
	in := "plain-slug\n" + `{"slug":"from-json"}` + "\n"
	ids, err := ReadIdentifiers(strings.NewReader(in))
	r.NoError(err)
	r.Equal([]string{"plain-slug", "from-json"}, ids)
}
