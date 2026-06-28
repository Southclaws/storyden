package pinecone_semdexer

import (
	"context"
	"fmt"
	"hash/fnv"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/google/uuid"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/datagraph/hydrate"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/vector/pinecone"
)

type pineconeSemdexer struct {
	client   *pinecone.Client
	index    *pinecone.Index
	hydrator *hydrate.Hydrator
	ef       semdex.Embedder
}

func New(ctx context.Context, cfg config.Config, pc *pinecone.Client, rh *hydrate.Hydrator, ef semdex.Embedder) (semdex.Semdexer, error) {
	if ef == nil {
		return nil, fault.New("a language model provider must be enabled for the pinecone semdexer to be enabled")
	}

	index, err := pc.GetOrCreateIndex(ctx, cfg.PineconeIndex)
	if err != nil {
		return nil, err
	}

	return &pineconeSemdexer{
		client:   pc,
		index:    index,
		hydrator: rh,
		ef:       ef,
	}, nil
}

func generateChunkID(id xid.ID, chunk string) string {
	// We don't currently support sharing chunks across content nodes, so append
	// the object's ID to the chunk's hash, to ensure it's unique to the object.

	hash := uuid.NewHash(fnv.New128(), uuid.NameSpaceOID, []byte(chunk), 4)

	prefix := id.String()

	return fmt.Sprintf("%s/%s", prefix, hash)
}

func chunkIDsFor(id xid.ID) func(chunk string) string {
	return func(chunk string) string {
		return generateChunkID(id, chunk)
	}
}

func chunkIDsForItem(object datagraph.Item) []string {
	return dt.Map(object.GetContent().Split(), chunkIDsFor(object.GetID()))
}

type chunk struct {
	id      string
	content string
}

func chunksFor(object datagraph.Item) []chunk {
	id := object.GetID()
	chunks := object.GetContent().Split()

	return dt.Map(chunks, func(c string) chunk {
		return chunk{
			id:      generateChunkID(id, c),
			content: c,
		}
	})
}
