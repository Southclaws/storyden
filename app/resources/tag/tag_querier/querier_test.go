package tag_querier

import (
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"

	"github.com/Southclaws/storyden/internal/ent"
)

func TestRoleHydrationTargetsCollectsPostAuthorsAndNodeOwners(t *testing.T) {
	t.Parallel()

	postAuthor := &ent.Account{ID: xid.New()}
	nodeOwner := &ent.Account{ID: xid.New()}
	nestedNodeOwner := &ent.Account{ID: xid.New()}

	tag := &ent.Tag{
		Edges: ent.TagEdges{
			Posts: []*ent.Post{
				{Edges: ent.PostEdges{Author: postAuthor}},
			},
			Nodes: []*ent.Node{
				{
					ID: xid.New(),
					Edges: ent.NodeEdges{
						Owner: nodeOwner,
						Nodes: []*ent.Node{{
							ID: xid.New(),
							Edges: ent.NodeEdges{Owner: nestedNodeOwner},
						}},
					},
				},
			},
		},
	}

	targets := roleHydrationTargets(tag)
	ids := make([]xid.ID, 0, len(targets))
	for _, a := range targets {
		ids = append(ids, a.ID)
	}

	assert.Contains(t, ids, postAuthor.ID)
	assert.Contains(t, ids, nodeOwner.ID)
	assert.Contains(t, ids, nestedNodeOwner.ID)
}
