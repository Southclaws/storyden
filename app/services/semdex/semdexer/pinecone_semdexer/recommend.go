package pinecone_semdexer

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

func (s *pineconeSemdexer) Recommend(ctx context.Context, object datagraph.Item) (datagraph.ItemList, error) {
	refs, err := s.RecommendRefs(ctx, object)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	items, err := s.hydrator.Hydrate(ctx, refs...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return items, nil
}

func (s *pineconeSemdexer) RecommendRefs(ctx context.Context, object datagraph.Item) (datagraph.RefList, error) {
	chunkIDs := chunkIDsForItem(object)
	if len(chunkIDs) == 0 {
		return nil, nil
	}

	response, err := s.index.FetchVectors(ctx, chunkIDs)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	chunkvecs := [][]float32{}

	for _, v := range response.Vectors {
		chunkvecs = append(chunkvecs, *v.Values)
	}

	targetvec := averageVectors(chunkvecs)

	result, err := s.index.QueryByVectorValues(ctx, &pinecone.QueryByVectorValuesRequest{
		Vector:          targetvec,
		TopK:            10,
		IncludeMetadata: true,
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	objects, err := mapScoredVectors(result.Matches)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	results := objects.ToRefs()

	deduped := dedupeChunks(results)

	filtered := filterChunks(deduped)

	// filter out the source of the recommendations query
	withoutSource := dt.Filter(filtered, func(r *datagraph.Ref) bool {
		return r.ID != object.GetID()
	})

	return withoutSource, nil
}

func averageVectors(datagraphID [][]float32) []float32 {
	if len(datagraphID) == 0 {
		return []float32{}
	}

	// Determine the length of vectors
	vectorLength := len(datagraphID[0])
	if vectorLength == 0 {
		return []float32{}
	}

	// Initialize a slice to store the sum of vectors
	sum := make([]float32, vectorLength)

	// Sum all vectors
	for _, vector := range datagraphID {
		if len(vector) != vectorLength {
			panic("Vectors must have the same length")
		}
		for i := 0; i < vectorLength; i++ {
			sum[i] += vector[i]
		}
	}

	// Compute the average
	count := float32(len(datagraphID))
	for i := 0; i < vectorLength; i++ {
		sum[i] /= count
	}

	return sum
}
