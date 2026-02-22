package library

import (
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/ent"
)

func RoleHydrationTargetsFromNode(root *ent.Node) []*ent.Account {
	if root == nil {
		return nil
	}

	targets := make([]*ent.Account, 0, 9)
	seenNodes := map[xid.ID]struct{}{}

	var walk func(*ent.Node)
	walk = func(n *ent.Node) {
		if n == nil {
			return
		}

		if _, ok := seenNodes[n.ID]; ok {
			return
		}
		seenNodes[n.ID] = struct{}{}

		if owner := n.Edges.Owner; owner != nil {
			targets = append(targets, owner)
		}

		if parent := n.Edges.Parent; parent != nil {
			walk(parent)
		}

		for _, child := range n.Edges.Nodes {
			walk(child)
		}
	}

	walk(root)

	return targets
}
