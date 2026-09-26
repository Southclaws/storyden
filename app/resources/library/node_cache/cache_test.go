package node_cache

import (
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"

	"github.com/Southclaws/storyden/app/resources/mark"
)

func TestCanonicalKey(t *testing.T) {
	id := xid.New()

	tests := map[string]struct {
		key  mark.Queryable
		want string
	}{
		"slug": {
			key:  mark.NewQueryKey("node-slug"),
			want: "node-slug",
		},
		"id": {
			key:  mark.NewQueryKeyID(id),
			want: id.String(),
		},
		"id and slug": {
			key:  mark.NewQueryKey(id.String() + "-node-slug"),
			want: id.String(),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, CanonicalKey(tt.key))
		})
	}
}
