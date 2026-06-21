package datagraph

import (
	"net/url"
	"testing"

	"github.com/Southclaws/storyden/internal/utils"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefURL(t *testing.T) {
	id := utils.Must(xid.FromString("crk0gvqfunp7891n7ah0"))
	ref := &Ref{ID: id, Kind: KindNode}

	u := ref.URL()

	assert.Equal(t, "sdr", u.Scheme)
	assert.Equal(t, "node/crk0gvqfunp7891n7ah0", u.Opaque)
}

func TestRefString(t *testing.T) {
	id := utils.Must(xid.FromString("crk0gvqfunp7891n7ah0"))
	ref := &Ref{ID: id, Kind: KindNode}

	assert.Equal(t, "sdr:node/crk0gvqfunp7891n7ah0", ref.String())
}

func TestRefStringRoundtrip(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	id := utils.Must(xid.FromString("crk0gvqfunp7891n7ah0"))
	original := &Ref{ID: id, Kind: KindNode}

	parsed, err := url.Parse(original.String())
	r.NoError(err)

	restored, err := NewRefFromSDR(*parsed)
	r.NoError(err)

	a.Equal(original.ID, restored.ID)
	a.Equal(original.Kind, restored.Kind)
}
