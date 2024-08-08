package weaviate_semdexer

import (
	"context"

	"github.com/rs/xid"
)

func (s *weaviateRefIndex) GetVectorFor(ctx context.Context, idx ...xid.ID) ([]float64, error) {
	// TODO: pull vectors for all items, compute average and return?
	return nil, nil
}
