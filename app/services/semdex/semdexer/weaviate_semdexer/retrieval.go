package weaviate_semdexer

import (
	"context"
	"encoding/json"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate/entities/models"

	"github.com/Southclaws/storyden/app/resources/datagraph"
)

func (o *weaviateSemdexer) GetMany(ctx context.Context, limit uint, ids ...xid.ID) (datagraph.RefList, error) {
	stringIDs := dt.Map(ids, func(x xid.ID) string { return x.String() })

	objects, err := o.wc.
		GraphQL().
		Get().
		WithClassName(o.cn.String()).
		WithWhere(
			filters.Where().
				WithPath([]string{"id"}).
				WithOperator(filters.ContainsAny).
				WithValueString(stringIDs...),
		).
		WithLimit(int(limit)).
		Do(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	data, err := mapResponseObjects(objects.Data)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	refs, err := dt.MapErr(data.Get[string(o.cn)], mapToNodeReference)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return datagraph.RefList(refs), nil
}

func mapWeaviateObject(o *models.Object) (*datagraph.Ref, error) {
	wo, err := unmarshalWeaviateObject(o.Properties)
	if err != nil {
		return nil, err
	}

	ref, err := mapToNodeReference(*wo)
	if err != nil {
		return nil, err
	}

	return ref, nil
}

func unmarshalWeaviateObject(p models.PropertySchema) (*WeaviateObject, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	wo := WeaviateObject{}
	if err := json.Unmarshal(b, &wo); err != nil {
		return nil, err
	}

	return &wo, nil
}
