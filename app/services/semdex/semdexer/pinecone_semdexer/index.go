package pinecone_semdexer

import (
	"context"
	"runtime"
	"sync"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/infrastructure/vector/pinecone"
)

func (c *pineconeSemdexer) Index(ctx context.Context, object datagraph.Item) error {
	chunks := object.GetContent().Split()

	if len(chunks) == 0 {
		return fault.New("no text chunks to index", fctx.With(ctx))
	}

	numWorkers := min(runtime.NumCPU(), len(chunks))
	chunkQueue := make(chan string, len(chunks))
	errChan := make(chan error, len(chunks))
	chunkChan := make(chan *pinecone.Vector, len(chunks))

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for chunk := range chunkQueue {
				vec, err := c.ef(ctx, chunk)
				if err != nil {
					errChan <- err
				}

				objectID := object.GetID()

				metadata, err := structpb.NewStruct(map[string]any{
					"datagraph_id":   objectID.String(),
					"datagraph_type": object.GetKind().String(),
					"name":           object.GetName(),
					"content":        chunk,
				})
				if err != nil {
					errChan <- err
				}

				chunkID := generateChunkID(objectID, chunk).String()

				chunkChan <- &pinecone.Vector{
					Id:       chunkID,
					Values:   vec,
					Metadata: metadata,
				}
			}
		}(i)
	}

	go func() {
		for _, chunk := range chunks {
			chunkQueue <- chunk
		}
		close(chunkQueue)
	}()

	go func() {
		wg.Wait()

		close(errChan)
		close(chunkChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	var vecs []*pinecone.Vector
	for vec := range chunkChan {
		vecs = append(vecs, vec)
	}

	_, err := c.index.UpsertVectors(ctx, vecs)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (c *pineconeSemdexer) Delete(ctx context.Context, object xid.ID) error {
	filter, err := structpb.NewStruct(map[string]any{
		"datagraph_id": object.String(),
	})
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	err = c.index.DeleteVectorsByFilter(ctx, filter)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
