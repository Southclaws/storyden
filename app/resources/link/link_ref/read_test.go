package link_ref_test

import (
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/internal/ent"
)

func TestMapOptionalSkipsLegacyURLs(t *testing.T) {
	for _, row := range []*ent.Link{
		nil,
		{URL: ""},
		{URL: "http://127.0.0.1/"},
		{URL: "https://example.com:443/"},
	} {
		require.False(t, link_ref.MapOptional(row).Ok())
	}

	ref, ok := link_ref.MapOptional(&ent.Link{URL: "https://example.com/"}).Get()
	require.True(t, ok)
	require.Equal(t, "https://example.com/", ref.URL.String())
}

func TestLinksFromModelSkipsLegacyURLs(t *testing.T) {
	first := &ent.Link{ID: xid.New(), URL: "https://example.com/first"}
	last := &ent.Link{ID: xid.New(), URL: "https://example.com/last"}
	rows := []*ent.Link{first, {URL: "http://127.0.0.1/"}, last, {URL: "https://example.com:443/"}}
	refs := link_ref.LinksFromModel(rows)
	require.Len(t, refs, 2)
	require.Equal(t, last.ID, refs[0].ID)
	require.Equal(t, first.ID, refs[1].ID)
	require.Empty(t, link_ref.LinksFromModel(rows[1:2]))
}
