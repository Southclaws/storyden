// Package semdex provides an interface for semantic indexing of the datagraph.
package semdex

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/rs/xid"
)

type Indexer interface {
	Index(ctx context.Context, object datagraph.Indexable) error
}

type Searcher interface {
	Search(ctx context.Context, query string) (datagraph.NodeReferenceList, error)
}

type Recommender interface {
	Recommend(ctx context.Context, object datagraph.Indexable) (datagraph.NodeReferenceList, error)
}

type Retriever interface {
	GetAll(ctx context.Context) (datagraph.NodeReferenceList, error)
	// GetVectorFor(ctx context.Context, idx ...xid.ID) ([]float64, error)
}

type Semdexer interface {
	Indexer
	Searcher
	Recommender
	Retriever
}

type OnlySearcher struct {
	Searcher
}

func (o *OnlySearcher) Search(ctx context.Context, query string) (datagraph.NodeReferenceList, error) {
	return o.Searcher.Search(ctx, query) // nolint:wrapcheck
}

func (o *OnlySearcher) Index(ctx context.Context, object datagraph.Indexable) error {
	return nil
}

func (o *OnlySearcher) Recommend(ctx context.Context, object datagraph.Indexable) (datagraph.NodeReferenceList, error) {
	return nil, nil
}

func (o *OnlySearcher) GetAll(ctx context.Context) (datagraph.NodeReferenceList, error) {
	return nil, nil
}

func (o *OnlySearcher) GetVectorFor(ctx context.Context, idx ...xid.ID) ([]float64, error) {
	return nil, nil
}

type Empty struct{}

func (n Empty) Index(ctx context.Context, object datagraph.Indexable) error {
	return nil
}

func (n Empty) Search(ctx context.Context, query string) (datagraph.NodeReferenceList, error) {
	return nil, nil
}
