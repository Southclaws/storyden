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

type OnlySearcher struct {
	Searcher
}

func (o *OnlySearcher) Search(ctx context.Context, query string) ([]*Result, error) {
	return o.Searcher.Search(ctx, query)
}

func (o *OnlySearcher) Index(ctx context.Context, object datagraph.Indexable) error {
	return nil
}

type Result struct {
	Id          xid.ID
	Type        datagraph.Kind
	Name        string
	Description string
	Slug        string
	ImageURL    string
}

type Empty struct{}

func (n Empty) Index(ctx context.Context, object datagraph.Indexable) error {
	return nil
}

func (n Empty) Search(ctx context.Context, query string) ([]*Result, error) {
	return nil, nil
}
