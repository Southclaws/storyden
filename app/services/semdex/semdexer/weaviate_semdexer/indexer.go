package weaviate_semdexer

import (
	"context"
	"errors"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	weaviate_errors "github.com/weaviate/weaviate-go-client/v4/weaviate/fault"
	"github.com/weaviate/weaviate/entities/models"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

func (s *weaviateSemdexer) Index(ctx context.Context, object datagraph.Item) error {
	chunks := object.GetContent().Split()

	if len(chunks) == 0 {
		return fault.New("no text chunks to index", fctx.With(ctx))
	}

	for _, chunk := range chunks {
		err := s.indexChunk(ctx, object, chunk)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	return nil
}

func (s *weaviateSemdexer) indexChunk(ctx context.Context, object datagraph.Item, chunk string) error {
	objectID := object.GetID()
	chunkID := generateChunkID(objectID, chunk).String()

	current, exists, err := s.existsByContent(ctx, objectID, chunk)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	props := map[string]any{
		"datagraph_id":   objectID.String(),
		"datagraph_type": object.GetKind(),
		"name":           object.GetName(),
		"description":    object.GetDesc(),
		"content":        chunk,
	}

	if exists {
		existingProps := current.Properties.(map[string]any)

		isSame := compareIndexedContentProperties(existingProps, props)
		if isSame {
			return nil
		}

		err = s.wc.Data().Updater().
			WithClassName(s.cn.String()).
			WithID(chunkID).
			WithProperties(props).
			Do(ctx)
	} else {
		_, err = s.wc.Data().Creator().
			WithClassName(s.cn.String()).
			WithID(chunkID).
			WithProperties(props).
			Do(ctx)
	}

	if err != nil {
		we := &weaviate_errors.WeaviateClientError{}
		if errors.As(err, &we) {
			if we.StatusCode == 422 {
				return nil
			}
		}

		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (s *weaviateSemdexer) existsByContent(ctx context.Context, objectID xid.ID, chunk string) (*models.Object, bool, error) {
	chunkID := generateChunkID(objectID, chunk)

	result, err := s.wc.Data().ObjectsGetter().
		WithClassName(s.cn.String()).
		WithID(chunkID.String()).
		Do(ctx)

	we := &weaviate_errors.WeaviateClientError{}
	if errors.As(err, &we) {
		if we.StatusCode == 404 {
			return nil, false, nil
		}
	}
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	if len(result) == 0 {
		return nil, false, nil
	}

	return result[0], true, nil
}

func compareIndexedContentProperties(a, b map[string]any) bool {
	if a["name"] != b["name"] {
		return false
	}
	if a["description"] != b["description"] {
		return false
	}
	if a["content"] != b["content"] {
		return false
	}

	return true
}
