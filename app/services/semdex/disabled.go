package semdex

import (
	"context"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/services/search/searcher"
)

type Disabled struct{}

var _ Semdexer = &Disabled{}

func (*Disabled) Index(ctx context.Context, object datagraph.Item) (int, error) {
	return 0, nil
}

func (*Disabled) Delete(ctx context.Context, object xid.ID) (int, error) {
	return 0, nil
}

func (*Disabled) Search(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) (*pagination.Result[datagraph.Item], error) {
	panic("semdex disabled: searcher switch bug")
}

func (*Disabled) SearchRefs(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) (*pagination.Result[*datagraph.Ref], error) {
	panic("semdex disabled: searcher switch bug")
}

func (*Disabled) SearchChunks(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) ([]*Chunk, error) {
	return nil, nil
}

func (*Disabled) Ask(ctx context.Context, q string) (chan string, chan error) {
	return nil, nil
}

func (*Disabled) Recommend(ctx context.Context, object datagraph.Item) (datagraph.ItemList, error) {
	return nil, nil
}

func (*Disabled) RecommendRefs(ctx context.Context, object datagraph.Item) (datagraph.RefList, error) {
	return nil, nil
}

func (*Disabled) ScoreRelevance(ctx context.Context, object datagraph.Item, idx ...xid.ID) (map[xid.ID]float64, error) {
	return nil, nil
}

func (*Disabled) GetMany(ctx context.Context, limit uint, ids ...xid.ID) (datagraph.RefList, error) {
	return nil, nil
}
