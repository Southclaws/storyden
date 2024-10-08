package collection_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func matchThreadToItem(t *testing.T, thread *openapi.Thread, item openapi.CollectionItem) {
	t.Helper()
	a := assert.New(t)

	a.Equal(openapi.DatagraphItemKindPost, item.Kind)
	a.Equal(thread.Id, item.Id)
	// a.Equal(thread.CreatedAt, item.CreatedAt) // TODO
	a.Equal(thread.Title, item.Name)
	a.Contains(thread.Slug, item.Slug)
	a.Equal(thread.Description, item.Description)
	a.Equal(thread.Author, item.Owner)
}

func matchNodeToItem(t *testing.T, node *openapi.Node, item openapi.CollectionItem) {
	t.Helper()
	a := assert.New(t)

	a.Equal(openapi.DatagraphItemKindNode, item.Kind)
	a.Equal(node.Id, item.Id)
	// a.Equal(node.CreatedAt, item.CreatedAt) // TODO
	a.Equal(node.Name, item.Name)
	a.Contains(node.Slug, item.Slug)
	a.Equal(node.Description, *item.Description)
	a.Equal(node.Owner, item.Owner)
}
