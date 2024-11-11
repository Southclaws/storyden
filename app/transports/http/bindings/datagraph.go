package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/services/search/searcher"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Datagraph struct {
	searcher searcher.Searcher
}

func NewDatagraph(
	searcher searcher.Searcher,
) Datagraph {
	return Datagraph{
		searcher: searcher,
	}
}

const datagraphSearchPageSize = 50

func (d Datagraph) DatagraphSearch(ctx context.Context, request openapi.DatagraphSearchRequestObject) (openapi.DatagraphSearchResponseObject, error) {
	if request.Params.Q == nil {
		return nil, fault.New("missing query")
	}

	pp := deserialisePageParams(request.Params.Page, datagraphSearchPageSize)

	kindFilter, err := opt.MapErr(opt.NewPtr(request.Params.Kind), deserialiseDatagraphKindList)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	opts := searcher.Options{
		Kinds: kindFilter,
	}

	r, err := d.searcher.Search(ctx, *request.Params.Q, pp, opts)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.DatagraphSearch200JSONResponse{
		DatagraphSearchOKJSONResponse: openapi.DatagraphSearchOKJSONResponse{
			CurrentPage: r.CurrentPage,
			Items:       dt.Map(r.Items, serialiseDatagraphItem),
			NextPage:    r.NextPage.Ptr(),
			PageSize:    r.Size,
			Results:     r.Results,
			TotalPages:  r.TotalPages,
		},
	}, nil
}

func deserialiseDatagraphKindList(ks []openapi.DatagraphItemKind) ([]datagraph.Kind, error) {
	return dt.MapErr(ks, deserialiseDatagraphKind)
}

func deserialiseDatagraphKind(v openapi.DatagraphItemKind) (datagraph.Kind, error) {
	return datagraph.NewKind(string(v))
}

func serialiseDatagraphItem(v datagraph.Item) openapi.DatagraphItem {
	out := openapi.DatagraphItem{}
	var err error

	switch in := v.(type) {
	case *post.Post:
		err = out.FromDatagraphItemPost(serialiseDatagraphItemPost(in))

	case *thread.Thread:
		err = out.FromDatagraphItemThread(serialiseDatagraphItemPostThread(in))

	case *reply.Reply:
		err = out.FromDatagraphItemReply(serialiseDatagraphItemPostReply(in))

	case *library.Node:
		err = out.FromDatagraphItemNode(serialiseDatagraphItemNode(in))

	case *profile.Public:
		err = out.FromDatagraphItemProfile(serialiseDatagraphItemProfile(in))

	default:
		err = fault.Newf("invalid datagraph item type: %T", v)
	}

	if err != nil {
		zap.L().Error("failed to serialise datagraph item", zap.Error(err))
	}

	return out
}

func serialiseDatagraphItemPost(in *post.Post) openapi.DatagraphItemPost {
	return openapi.DatagraphItemPost{
		Kind: openapi.DatagraphItemKindPost,
		Ref:  serialisePost(in),
	}
}

func serialiseDatagraphItemPostThread(in *thread.Thread) openapi.DatagraphItemThread {
	return openapi.DatagraphItemThread{
		Kind: openapi.DatagraphItemKindThread,
		Ref:  serialiseThread(in),
	}
}

func serialiseDatagraphItemPostReply(in *reply.Reply) openapi.DatagraphItemReply {
	return openapi.DatagraphItemReply{
		Kind: openapi.DatagraphItemKindReply,
		Ref:  serialiseReply(in),
	}
}

func serialiseDatagraphItemNode(in *library.Node) openapi.DatagraphItemNode {
	return openapi.DatagraphItemNode{
		Kind: openapi.DatagraphItemKindNode,
		Ref:  serialiseNode(in),
	}
}

func serialiseDatagraphItemProfile(in *profile.Public) openapi.DatagraphItemProfile {
	return openapi.DatagraphItemProfile{
		Kind: openapi.DatagraphItemKindProfile,
		Ref:  serialiseProfile(in),
	}
}

func serialiseDatagraphItemList(in datagraph.ItemList) openapi.DatagraphItemList {
	return dt.Map(in, serialiseDatagraphItem)
}
