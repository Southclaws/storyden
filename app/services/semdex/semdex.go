// Package semdex provides an interface for semantic indexing of the datagraph.
package semdex

import (
	"context"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

type Indexer interface {
	Index(ctx context.Context, object datagraph.Indexable) error
}

type Searcher interface {
	Search(ctx context.Context, query string) ([]*Result, error)
}

type Semdexer interface {
	Indexer
	Searcher
}

type Result struct {
	Id   xid.ID
	Name string
	Type string
}

type Empty struct{}

func (n Empty) Index(ctx context.Context, object datagraph.Indexable) error {
	return nil
}

func (n Empty) Search(ctx context.Context, query string) ([]*Result, error) {
	return nil, nil
}
