package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/datagraph/semdex"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Datagraph struct {
	ds semdex.Searcher
}

func NewDatagraph(
	ds semdex.Searcher,
) Datagraph {
	return Datagraph{
		ds: ds,
	}
}

func (d Datagraph) DatagraphSearch(ctx context.Context, request openapi.DatagraphSearchRequestObject) (openapi.DatagraphSearchResponseObject, error) {
	if request.Params.Q == nil {
		return nil, fault.New("missing query")
	}

	r, err := d.ds.Search(ctx, *request.Params.Q)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	items := dt.Map(r, serialiseDatagraphItem)

	return openapi.DatagraphSearch200JSONResponse{
		DatagraphSearchOKJSONResponse: openapi.DatagraphSearchOKJSONResponse{
			// TODO: pagination
			Items: items,
		},
	}, nil
}

func serialiseDatagraphItem(v datagraph.Item) openapi.DatagraphItem {
	desc := v.GetDesc()
	return openapi.DatagraphItem{
		Kind:        serialiseDatagraphKind(v.GetKind()),
		Id:          v.GetID().String(),
		Name:        v.GetName(),
		Slug:        v.GetSlug(),
		Description: &desc,
	}
}

func serialiseDatagraphKind(in datagraph.Kind) openapi.DatagraphItemKind {
	switch in {
	case datagraph.KindPost:
		return openapi.DatagraphItemKindPost
	case datagraph.KindNode:
		return openapi.DatagraphItemKindNode
	case datagraph.KindProfile:
		return openapi.DatagraphItemKindProfile
	default:
		panic(fault.Newf("unknown kind '%s'", in))
	}
}
