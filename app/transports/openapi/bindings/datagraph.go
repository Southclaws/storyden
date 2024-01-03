package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/services/datagraph_searcher"
	"github.com/Southclaws/storyden/internal/openapi"
)

type Datagraph struct {
	ds datagraph_searcher.Searcher
}

func NewDatagraph(
	ds datagraph_searcher.Searcher,
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

func serialiseDatagraphItem(v *datagraph_searcher.Result) openapi.DatagraphItem {
	return openapi.DatagraphItem{
		Id:   v.Id.String(),
		Type: openapi.DatagraphItemType(v.Type),
		Name: v.Name,
	}
}
