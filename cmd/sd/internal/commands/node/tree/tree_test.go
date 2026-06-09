package tree

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func TestRenderTree(t *testing.T) {
	nodes := []openapi.NodeWithChildren{
		{
			Id:   "node_root_1",
			Name: "Root One",
			Slug: "root-one",
			Children: []openapi.NodeWithChildren{
				{
					Id:   "node_child_1",
					Name: "Child One",
					Slug: "child-one",
				},
				{
					Id:   "node_child_2",
					Name: "Child Two",
					Slug: "child-two",
					Children: []openapi.NodeWithChildren{
						{
							Id:   "node_grandchild_1",
							Name: "Grandchild",
							Slug: "grandchild",
						},
					},
				},
			},
		},
		{
			Id:   "node_root_2",
			Name: "Root Two",
			Slug: "root-two",
		},
	}

	var out bytes.Buffer

	err := renderTree(&out, nodes, "")

	require.NoError(t, err)
	require.Equal(t, `.
├── Root One [slug=root-one id=node_root_1]
│   ├── Child One [slug=child-one id=node_child_1]
│   └── Child Two [slug=child-two id=node_child_2]
│       └── Grandchild [slug=grandchild id=node_grandchild_1]
└── Root Two [slug=root-two id=node_root_2]
`, out.String())
}

func TestValidateVisibilities(t *testing.T) {
	r := require.New(t)

	r.NoError(validateVisibilities(nil))
	r.NoError(validateVisibilities([]string{"published"}))
	r.NoError(validateVisibilities([]string{"draft", "review"}))
	r.ErrorContains(validateVisibilities([]string{"private"}), "invalid --visibility: private")
}

func TestNodeLabelIsSingleLine(t *testing.T) {
	node := openapi.NodeWithChildren{
		Id:   "node_id",
		Name: "Odd\nName\tHere",
		Slug: "odd-slug",
	}

	require.Equal(t, "Odd Name Here [slug=odd-slug id=node_id]", nodeLabel(node))
}
