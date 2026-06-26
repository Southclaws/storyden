package bindings

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/services/search/searcher"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Datagraph struct {
	searcher       searcher.Searcher
	accountQuerier *account_querier.Querier
}

func NewDatagraph(
	searcher searcher.Searcher,
	accountQuerier *account_querier.Querier,
	router *echo.Echo,
) Datagraph {
	d := Datagraph{
		searcher:       searcher,
		accountQuerier: accountQuerier,
	}

	_ = router

	return d
}

const (
	datagraphSearchPageSize = 50
	datagraphMatchesLimit   = 20
)

func (d Datagraph) DatagraphSearch(ctx context.Context, request openapi.DatagraphSearchRequestObject) (openapi.DatagraphSearchResponseObject, error) {
	pp := deserialisePageParams(request.Params.Page, datagraphSearchPageSize)

	kindFilter, err := opt.MapErr(opt.NewPtr(request.Params.Kind), deserialiseDatagraphKindList)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	authorFilter, err := opt.MapErr(opt.NewPtr(request.Params.Authors), func(ids []openapi.Identifier) ([]account.AccountID, error) {
		return d.resolveAuthorFilter(ctx, ids)
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	categoryFilter, err := opt.MapErr(opt.NewPtr(request.Params.Categories), deserialiseCategoryList)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	tagFilter, err := opt.MapErr(opt.NewPtr(request.Params.Tags), deserialiseTagList)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	opts := searcher.Options{
		Kinds:      kindFilter,
		Authors:    authorFilter,
		Categories: categoryFilter,
		Tags:       tagFilter,
	}

	r, err := d.searcher.Search(ctx, opt.NewPtr(request.Params.Q).Or(""), pp, opts)
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

func (d Datagraph) DatagraphMatches(ctx context.Context, request openapi.DatagraphMatchesRequestObject) (openapi.DatagraphMatchesResponseObject, error) {
	kindFilter, err := opt.MapErr(opt.NewPtr(request.Params.Kind), deserialiseDatagraphKindList)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	opts := searcher.Options{Kinds: kindFilter}

	matches, err := d.searcher.MatchFast(ctx, request.Params.Q, datagraphMatchesLimit, opts)
	if err != nil {
		if errors.Is(err, searcher.ErrFastMatchesUnavailable) {
			return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("datagraph matches are not enabled"))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.DatagraphMatches200JSONResponse{
		DatagraphMatchesOKJSONResponse: openapi.DatagraphMatchesOKJSONResponse{
			Items: serialiseDatagraphMatchList(matches),
		},
	}, nil
}

func deserialiseDatagraphKindList(ks []openapi.DatagraphItemKind) ([]datagraph.Kind, error) {
	return dt.MapErr(ks, deserialiseDatagraphKind)
}

func deserialiseDatagraphKind(v openapi.DatagraphItemKind) (datagraph.Kind, error) {
	return datagraph.NewKind(string(v))
}

func (d *Datagraph) resolveAuthorFilter(ctx context.Context, identifiers []openapi.Identifier) ([]account.AccountID, error) {
	if len(identifiers) == 0 {
		return []account.AccountID{}, nil
	}

	var ids []account.AccountID
	var handles []string

	for _, identifier := range identifiers {
		idStr := string(identifier)
		if parsed, err := xid.FromString(idStr); err == nil {
			ids = append(ids, account.AccountID(parsed))
		} else {
			handles = append(handles, idStr)
		}
	}

	if len(handles) > 0 {
		accounts, err := d.accountQuerier.ProbeMany(ctx, handles...)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		for _, acc := range accounts {
			ids = append(ids, acc.ID)
		}
	}

	return ids, nil
}

func deserialiseCategoryList(ids []openapi.Identifier) ([]category.CategoryID, error) {
	return dt.MapErr(ids, func(id openapi.Identifier) (category.CategoryID, error) {
		parsed, err := xid.FromString(string(id))
		if err != nil {
			return category.CategoryID(xid.NilID()), err
		}
		return category.CategoryID(parsed), nil
	})
}

func deserialiseTagList(names []openapi.TagName) ([]tag_ref.Name, error) {
	return dt.Map(names, func(name openapi.TagName) tag_ref.Name {
		return tag_ref.NewName(string(name))
	}), nil
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
		slog.Error("failed to serialise datagraph item", slog.String("error", err.Error()))
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
		Ref:  serialiseReplyPtr(in),
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

func serialiseDatagraphMatchList(in datagraph.MatchList) openapi.DatagraphMatchList {
	return dt.Map(in, serialiseDatagraphMatch)
}

func serialiseDatagraphMatch(in datagraph.Match) openapi.DatagraphMatch {
	description := opt.NewIf(in.Description, func(s string) bool { return s != "" }).Ptr()

	return openapi.DatagraphMatch{
		Id:          openapi.Identifier(in.ID.String()),
		Kind:        openapi.DatagraphItemKind(in.Kind.String()),
		Slug:        in.Slug,
		Name:        in.Name,
		Description: description,
	}
}
