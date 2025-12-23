package bindings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/services/search/searcher"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/app/services/system/instance_info"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Datagraph struct {
	searcher searcher.Searcher
	asker    semdex.Asker
}

func NewDatagraph(
	info *instance_info.Provider,
	searcher searcher.Searcher,
	asker semdex.Asker,
	router *echo.Echo,
) Datagraph {
	d := Datagraph{
		searcher: searcher,
		asker:    asker,
	}

	// The generated OpenAPI code does not expose the underlying ResponseWriter
	// which we need for streaming Q&A responses for that ✨chatgpt✨ effect.
	router.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ii, err := info.Get(context.Background())
			if err != nil {
				return fault.Wrap(err, fctx.With(c.Request().Context()))
			}
			isEnabled := ii.Capabilities.Has(instance_info.CapabilitySemdex)

			path := c.Path()
			if path == "/api/datagraph/ask" {
				if !isEnabled {
					return echo.NewHTTPError(http.StatusNotImplemented, "Semdex is not enabled")
				}

				ctx := c.Request().Context()

				query := c.QueryParam("q")

				p := c.QueryParam("parent_question_id")
				parent, err := opt.MapErr(opt.NewSafe(p, p != ""), xid.FromString)
				if err != nil {
					return fault.Wrap(err, fctx.With(ctx))
				}

				iter, err := d.asker.Ask(ctx, query, parent)
				if err != nil {
					return fault.Wrap(err, fctx.With(ctx))
				}

				w := c.Response().Writer

				w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")

				// if fails to cast, do nothing, flusher nil, just guard clause.
				flusher, ok := w.(http.Flusher)
				if !ok {
					if innerWriter := unwrapWriter(w); innerWriter != nil {
						flusher, _ = innerWriter.(http.Flusher)
					}
				}

				for chunk, err := range iter {
					if err != nil {
						return err
					}

					msg := ""

					switch v := chunk.(type) {
					case *semdex.AskResponseChunkText:
						b, err := json.Marshal(chunk)
						if err != nil {
							return err
						}

						msg = fmt.Sprintf("event: text\ndata: %s\n\n", string(b))

					case *semdex.AskResponseChunkMeta:
						b, err := serialiseAskResponseChunkMeta(*v)
						if err != nil {
							return err
						}

						msg = fmt.Sprintf("event: meta\ndata: %s\n\n", string(b))
					}

					if _, err := w.Write([]byte(msg)); err != nil {
						return err
					}

					if flusher != nil {
						flusher.Flush()
					}
				}

				if _, err := w.Write([]byte("event: end\n\n")); err != nil {
					return err
				}

				if flusher != nil {
					flusher.Flush()
				}

				return nil
			}

			return next(c)
		}
	})

	return d
}

func unwrapWriter(w http.ResponseWriter) http.ResponseWriter {
	switch v := w.(type) {
	case interface{ Unwrap() http.ResponseWriter }:
		return v.Unwrap()
	default:
		return nil
	}
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

	authorFilter, err := opt.MapErr(opt.NewPtr(request.Params.Authors), deserialiseAuthorList)
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

	r, err := d.searcher.Search(ctx, request.Params.Q, pp, opts)
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

func (d Datagraph) DatagraphAsk(ctx context.Context, request openapi.DatagraphAskRequestObject) (openapi.DatagraphAskResponseObject, error) {
	// NOTE: Unused stub, see middleware above.
	return nil, nil
}

func deserialiseDatagraphKindList(ks []openapi.DatagraphItemKind) ([]datagraph.Kind, error) {
	return dt.MapErr(ks, deserialiseDatagraphKind)
}

func deserialiseDatagraphKind(v openapi.DatagraphItemKind) (datagraph.Kind, error) {
	return datagraph.NewKind(string(v))
}

func deserialiseAuthorList(ids []openapi.Identifier) ([]account.AccountID, error) {
	return dt.MapErr(ids, func(id openapi.Identifier) (account.AccountID, error) {
		parsed, err := xid.FromString(string(id))
		if err != nil {
			return account.AccountID(xid.NilID()), err
		}
		return account.AccountID(parsed), nil
	})
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

func serialiseAskResponseChunkMeta(in semdex.AskResponseChunkMeta) ([]byte, error) {
	urls := dt.Map(in.URLs, func(u url.URL) string {
		return u.String()
	})

	refs := dt.Map(in.Refs, func(r *datagraph.Ref) map[string]any {
		return map[string]any{
			"id":   r.ID.String(),
			"kind": r.Kind.String(),
		}
	})

	return json.Marshal(map[string]any{
		"urls": urls,
		"refs": refs,
	})
}
