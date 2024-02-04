package indexer

import (
	"context"
	"errors"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	weaviate_errors "github.com/weaviate/weaviate-go-client/v4/weaviate/fault"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/semdex"
)

func New(lc fx.Lifecycle, wc *weaviate.Client) (semdex.Service, error) {
	lc.Append(fx.StartHook(func(ctx context.Context) error {
		r, err := wc.Schema().
			ClassExistenceChecker().
			WithClassName(semdex.TestClassName).
			Do(ctx)
		if err != nil {
			return fault.Wrap(err)
		}

		if !r {
			err := wc.Schema().
				ClassCreator().
				WithClass(semdex.TestClassObject).
				Do(ctx)
			if err != nil {
				return fault.Wrap(err)
			}
		}

		return nil
	}))

	return &service{wc}, nil
}

type service struct {
	wc *weaviate.Client
}

func (s *service) Index(ctx context.Context, object datagraph.Indexable) error {
	content := object.GetText()
	id := uuid.NewSHA1(uuid.NameSpaceOID, []byte(content)).String()

	// Don't bother indexing if the content is too short.
	if len(content) < 30 {
		return nil
	}

	_, err := s.wc.Data().Creator().
		WithClassName(semdex.TestClassName).
		WithID(id).
		WithProperties(map[string]any{
			"datagraph_id":   object.GetID().String(),
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
