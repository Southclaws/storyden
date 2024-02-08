package weaviate

import (
	"context"
	"errors"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/google/uuid"
	"github.com/rs/xid"
	weaviate_errors "github.com/weaviate/weaviate-go-client/v4/weaviate/fault"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

func (s *weaviateSemdexer) Index(ctx context.Context, object datagraph.Indexable) error {
	content := object.GetText()
	sid := object.GetID()

	wid := GetWeaviateID(object.GetID())

	// Don't bother indexing if the content is too short.
	if len(content) < 30 {
		return nil
	}

	_, err := s.wc.Data().Creator().
		WithClassName(TestClassName).
		WithID(wid).
		WithProperties(map[string]any{
			"datagraph_id":   sid.String(),
			"datagraph_type": object.GetKind(),
			"name":           object.GetName(),
			"content":        content,
			"props":          object.GetProps(),
		}).
		Do(ctx)
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

func GetWeaviateID(id xid.ID) string {
	return uuid.NewSHA1(uuid.NameSpaceOID, id.Bytes()).String()
}
