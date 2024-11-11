// Package semdex provides an interface for semantic indexing of the datagraph.
package semdex

import (
	"context"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
)

type Indexer interface {
	Index(ctx context.Context, object datagraph.Item) error
}

type Deleter interface {
	Delete(ctx context.Context, object xid.ID) error
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

type Tagger interface {
	SuggestTags(ctx context.Context, content datagraph.Content, available tag_ref.Names) (tag_ref.Names, error)
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
	GetMany(ctx context.Context, limit uint, ids ...xid.ID) (datagraph.RefList, error)
	// GetVectorFor(ctx context.Context, idx ...xid.ID) ([]float64, error)
}

type RefSemdexer interface {
	Indexer
	Deleter
	RefSearcher
	RefRecommender
	Tagger
	Retriever
	RelevanceScorer
	Summariser
}

type Semdexer interface {
	Indexer
	Deleter
	Searcher
	Recommender
	Tagger
	Retriever
	RelevanceScorer
	Summariser
}

type Disabled struct{}

func (*Disabled) Index(ctx context.Context, object datagraph.Item) error { return nil }
func (*Disabled) Delete(ctx context.Context, object xid.ID) error        { return nil }

func (*Disabled) Search(ctx context.Context, query string) (datagraph.ItemList, error) {
	return nil, nil
}

func (*Disabled) Recommend(ctx context.Context, object datagraph.Item) (datagraph.ItemList, error) {
	return nil, nil
}

func (*Disabled) SuggestTags(ctx context.Context, content datagraph.Content, available tag_ref.Names) (tag_ref.Names, error) {
	return nil, nil
}

func (*Disabled) ScoreRelevance(ctx context.Context, object datagraph.Item, idx ...xid.ID) (map[xid.ID]float64, error) {
	return nil, nil
}

func (*Disabled) Summarise(ctx context.Context, object datagraph.Item) (string, error) {
	return "", nil
}

func (*Disabled) GetMany(ctx context.Context, limit uint, ids ...xid.ID) (datagraph.RefList, error) {
	return nil, nil
}
