// Package semdex provides an interface for semantic indexing of the datagraph.
package semdex

import (
	"context"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

type Indexer interface {
	Index(ctx context.Context, object datagraph.Item) error
}

type Searcher interface {
	Search(ctx context.Context, query string) (datagraph.ItemList, error)
}

type RefSearcher interface {
	Search(ctx context.Context, query string) (datagraph.RefList, error)
}

type Recommender interface {
	Recommend(ctx context.Context, object datagraph.Item) (datagraph.ItemList, error)
}

type RefRecommender interface {
	Recommend(ctx context.Context, object datagraph.Item) (datagraph.RefList, error)
}

type RelevanceScorer interface {
	ScoreRelevance(ctx context.Context, object datagraph.Item, idx ...xid.ID) (map[xid.ID]float64, error)
}

type Summariser interface {
	Summarise(ctx context.Context, object datagraph.Item) (string, error)
}

type Retriever interface {
	GetAll(ctx context.Context) (datagraph.RefList, error)
	// GetVectorFor(ctx context.Context, idx ...xid.ID) ([]float64, error)
}

type RefSemdexer interface {
	Indexer
	RefSearcher
	RefRecommender
	Retriever
	RelevanceScorer
	Summariser
}

type Semdexer interface {
	Indexer
	Searcher
	Recommender
	Retriever
	RelevanceScorer
	Summariser
}

type OnlySearcher struct {
	Searcher
}

func (o *OnlySearcher) Search(ctx context.Context, query string) (datagraph.ItemList, error) {
	return o.Searcher.Search(ctx, query)
}

func (o *OnlySearcher) Index(ctx context.Context, object datagraph.Item) error {
	return nil
}

func (o *OnlySearcher) Recommend(ctx context.Context, object datagraph.Item) (datagraph.ItemList, error) {
	return nil, nil
}

func (o *OnlySearcher) ScoreRelevance(ctx context.Context, object datagraph.Item, idx ...xid.ID) (map[xid.ID]float64, error) {
	return nil, nil
}

func (o *OnlySearcher) Summarise(ctx context.Context, object datagraph.Item) (string, error) {
	return "", nil
}

func (o *OnlySearcher) GetAll(ctx context.Context) (datagraph.RefList, error) {
	return nil, nil
}

func (o *OnlySearcher) GetVectorFor(ctx context.Context, idx ...xid.ID) ([]float64, error) {
	return nil, nil
}

type Empty struct{}

func (n Empty) Index(ctx context.Context, object datagraph.Item) error {
	return nil
}

func (n Empty) Search(ctx context.Context, query string) (datagraph.RefList, error) {
	return nil, nil
}
