// Package semdex provides an interface for semantic indexing of the datagraph.
package semdex

import (
	"context"

	"github.com/weaviate/weaviate/entities/models"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

type Service interface {
	Index(ctx context.Context, object datagraph.Indexable) error
}

type Empty struct{}

func (n Empty) Index(ctx context.Context, object datagraph.Indexable) error {
	return nil
}

// NOT PROD READY: Just using local transformers for now.

const TestClassName = "ContentText2vecTransformers"

var TestClassObject = &models.Class{
	Class:      TestClassName,
	Vectorizer: "text2vec-transformers",
	ModuleConfig: map[string]interface{}{
		// "text2vec-openai":   map[string]interface{}{},
		// "generative-openai": map[string]interface{}{},
		"text2vec-transformers": map[string]interface{}{},
	},
}
