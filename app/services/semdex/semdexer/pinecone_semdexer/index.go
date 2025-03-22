package pinecone_semdexer

import (
	"context"
	"runtime"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/alitto/pond/v2"
	"github.com/rs/xid"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/infrastructure/vector/pinecone"
)

func (s *pineconeSemdexer) Index(ctx context.Context, object datagraph.Item) (int, error) {
	inserts, deletes, err := s.buildIndexOps(ctx, object)
	if err != nil {
		return 0, fault.Wrap(err, fctx.With(ctx))
	}

	if len(inserts) > 0 {
		_, err = s.index.UpsertVectors(ctx, inserts)
		if err != nil {
			return 0, fault.Wrap(err, fctx.With(ctx))
		}
	}

	if len(deletes) > 0 {
		err = s.deleteVectors(ctx, deletes)
		if err != nil {
			return 0, fault.Wrap(err, fctx.With(ctx))
		}
	}

	changes := len(inserts) - len(deletes)

	return changes, nil
}

func (c *pineconeSemdexer) Delete(ctx context.Context, object xid.ID) (int, error) {
	prefix := object.String()
	vectors, err := c.index.ListVectors(ctx, &pinecone.ListVectorsRequest{
		Prefix: &prefix,
	})
	if err != nil {
		return 0, fault.Wrap(err, fctx.With(ctx))
	}

	if len(vectors.VectorIds) == 0 {
		return 0, nil
	}

	ids := dt.Map(vectors.VectorIds, func(id *string) string { return *id })

	err = c.index.DeleteVectorsById(ctx, ids)
	if err != nil {
		return 0, fault.Wrap(err, fctx.With(ctx))
	}

	return len(ids), nil
}

func (s *pineconeSemdexer) deleteVectors(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	err := s.index.DeleteVectorsById(ctx, ids)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (s *pineconeSemdexer) buildIndexOps(ctx context.Context, object datagraph.Item) ([]*pinecone.Vector, []string, error) {
	allChunks := chunksFor(object)
	if len(allChunks) == 0 {
		return nil, nil, nil
	}
	chunkIDs := dt.Map(allChunks, func(c chunk) string { return c.id })

	objectID := object.GetID()

	inputChunkTable := lo.SliceToMap(allChunks, func(c chunk) (string, chunk) {
		return c.id, c
	})

	prefix := objectID.String()
	indexedChunk, err := s.index.ListVectors(ctx, &pinecone.ListVectorsRequest{
		Prefix: &prefix,
	})
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	vecids := dt.Map(indexedChunk.VectorIds, func(id *string) string { return *id })

	indexedChunkTable := map[string]*pinecone.Vector{}
	if len(vecids) > 0 {
		resp, err := s.index.FetchVectors(ctx, vecids)
		if err != nil {
			return nil, nil, fault.Wrap(err, fctx.With(ctx))
		}
		indexedChunkTable = resp.Vectors
	}

	pool := pond.NewResultPool[*pinecone.Vector](min(runtime.NumCPU(), len(chunkIDs)))
	group := pool.NewGroupContext(ctx)

	for id, chunk := range inputChunkTable {
		_, exists := indexedChunkTable[id]
		if exists {
			continue
		}

		group.SubmitErr(func() (*pinecone.Vector, error) {
			vec, err := s.ef(ctx, chunk.content)
			if err != nil {
				return nil, err
			}

			metadata, err := structpb.NewStruct(map[string]any{
				"datagraph_id":   objectID.String(),
				"datagraph_type": object.GetKind().String(),
				"name":           object.GetName(),
				"content":        chunk.content,
			})
			if err != nil {
				return nil, err
			}

			return &pinecone.Vector{
				Id:       id,
				Values:   &vec,
				Metadata: metadata,
			}, nil
		})
	}

	inserts, err := group.Wait()
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	// build a list of vectors to delete by yielding items that are indexed but
	// not present in the input object chunk table.
	deletes := []string{}
	for id := range indexedChunkTable {
		_, exists := inputChunkTable[id]
		if exists {
			continue
		}

		deletes = append(deletes, id)
	}

	return inserts, deletes, nil
}
