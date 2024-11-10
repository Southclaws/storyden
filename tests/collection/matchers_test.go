package collection_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func matchThreadToItem(t *testing.T, thread *openapi.Thread, item openapi.CollectionItem) {
	t.Helper()
	a := assert.New(t)

	itemPost, err := item.Item.AsDatagraphItemPost()
	a.NoError(err)

	a.Equal(openapi.DatagraphItemKindPost, itemPost.Kind)
	a.Equal(thread.Id, itemPost.Ref.Id)
	a.Equal(thread.CreatedAt, itemPost.Ref.CreatedAt)
	a.Equal(thread.Title, itemPost.Ref.Title)
	a.Contains(thread.Slug, itemPost.Ref.Slug)
	a.Equal(thread.Description, itemPost.Ref.Description)
	a.Equal(thread.Author, itemPost.Ref.Author)
}

func matchNodeToItem(t *testing.T, node *openapi.Node, item openapi.CollectionItem) {
	t.Helper()
	a := assert.New(t)

	itemNode, err := item.Item.AsDatagraphItemNode()
	a.NoError(err)

	a.Equal(openapi.DatagraphItemKindNode, itemNode.Kind)
	a.Equal(node.Id, itemNode.Ref.Id)
	a.Equal(node.CreatedAt, itemNode.Ref.CreatedAt)
	a.Equal(node.Name, itemNode.Ref.Name)
	a.Contains(node.Slug, itemNode.Ref.Slug)
	a.Equal(node.Description, itemNode.Ref.Description)
	a.Equal(node.Owner, itemNode.Ref.Owner)
}
