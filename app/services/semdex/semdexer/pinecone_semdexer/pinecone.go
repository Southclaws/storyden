package pinecone_semdexer

import (
	"context"
	"hash/fnv"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/google/uuid"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/datagraph/hydrate"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/ai"
	"github.com/Southclaws/storyden/internal/infrastructure/vector/pinecone"
)

type pineconeSemdexer struct {
	client   *pinecone.Client
	index    *pinecone.Index
	hydrator *hydrate.Hydrator
	ef       ai.Embedder
}

func New(ctx context.Context, cfg config.Config, pc *pinecone.Client, rh *hydrate.Hydrator, aip ai.Prompter) (semdex.Semdexer, error) {
	if _, ok := aip.(*ai.Disabled); ok {
		return nil, fault.New("a language model provider must be enabled for the pinecone semdexer to be enabled")
	}

	ef := aip.EmbeddingFunc()

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

func generateChunkID(id xid.ID, chunk string) uuid.UUID {
	// We don't currently support sharing chunks across content nodes, so append
	// the object's ID to the chunk's hash, to ensure it's unique to the object.
	payload := []byte(append(id.Bytes(), chunk...))

	return uuid.NewHash(fnv.New128(), uuid.NameSpaceOID, payload, 4)
}

func chunkIDsFor(id xid.ID) func(chunk string) string {
	return func(chunk string) string {
		return generateChunkID(id, chunk).String()
	}
}

func chunkIDsForItem(object datagraph.Item) []string {
	return dt.Map(object.GetContent().Split(), chunkIDsFor(object.GetID()))
}

func (c *pineconeSemdexer) GetMany(ctx context.Context, limit uint, ids ...xid.ID) (datagraph.RefList, error) {
	return nil, nil
}
