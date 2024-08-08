package weaviate_semdexer

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/weaviate/weaviate/entities/models"
)

// mergeErrors merges GraphQL errors with other errors, because the Weaviate
// client library separates application level errors and transport level errors.
func mergeErrors(g *models.GraphQLResponse, err error) (*models.GraphQLResponse, error) {
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(context.Background()))
	}

	if len(g.Errors) > 0 {
		return nil, fault.Wrap(gqlerror(g.Errors), fctx.With(context.Background()))
	}

	return g, nil
}
