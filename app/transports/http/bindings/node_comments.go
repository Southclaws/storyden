package bindings

import (
	"context"
	"net/url"
	"strconv"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	library_resources "github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/library/node_comment"
	thread_service "github.com/Southclaws/storyden/app/services/thread"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type NodeComments struct {
	nodeCommentSvc *node_comment.Manager
}

func NewNodeComments(nodeCommentSvc *node_comment.Manager) NodeComments {
	return NodeComments{nodeCommentSvc}
}

func (nc *NodeComments) NodeCommentCreate(ctx context.Context, request openapi.NodeCommentCreateRequestObject) (openapi.NodeCommentCreateResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	status, err := opt.MapErr(opt.NewPtr(request.Body.Visibility), deserialiseThreadStatus)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var meta map[string]any
	if request.Body.Meta != nil {
		meta = *request.Body.Meta
	}

	tags := opt.Map(opt.NewPtr(request.Body.Tags), func(tags []string) tag_ref.Names {
		return dt.Map(tags, deserialiseTagName)
	})

	richContent, err := opt.MapErr(opt.NewPtr(request.Body.Body), datagraph.NewRichText)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	category := opt.NewPtrMap(request.Body.Category, func(cat openapi.Identifier) xid.ID {
		return openapi.ParseID(cat)
	})

	url, err := opt.MapErr(opt.NewPtr(request.Body.Url), func(s string) (url.URL, error) {
		u, err := url.Parse(s)
		if err != nil {
			return url.URL{}, err
		}
		return *u, nil
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	thread, err := nc.nodeCommentSvc.Create(ctx,
		library_resources.NewKey(request.NodeSlug),
		request.Body.Title,
		accountID,
		meta,
		thread_service.Partial{
			Content:    richContent,
			Category:   category,
			Tags:       tags,
			Visibility: status,
			URL:        url,
		},
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeCommentCreate200JSONResponse{
		NodeCommentCreateOKJSONResponse: openapi.NodeCommentCreateOKJSONResponse(serialiseThread(thread)),
	}, nil
}

func (nc *NodeComments) NodeCommentList(ctx context.Context, request openapi.NodeCommentListRequestObject) (openapi.NodeCommentListResponseObject, error) {
	accountID := session.GetOptAccountID(ctx)
	pageSize := 50

	page := opt.NewPtrMap(request.Params.Page, func(s string) int {
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return 0
		}
		return max(1, int(v))
	}).Or(1)

	page = max(0, page-1)

	result, err := nc.nodeCommentSvc.List(ctx,
		deserialiseNodeMark(request.NodeSlug),
		page,
		pageSize,
		accountID,
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	page = result.CurrentPage + 1
	nextPage := opt.Map(result.NextPage, func(i int) int { return i + 1 })

	return openapi.NodeCommentList200JSONResponse{
		NodeCommentListOKJSONResponse: openapi.NodeCommentListOKJSONResponse{
			CurrentPage: page,
			NextPage:    nextPage.Ptr(),
			PageSize:    result.PageSize,
			Results:     result.Results,
			TotalPages:  result.TotalPages,
			Threads:     dt.Map(result.Threads, serialiseThreadReference),
		},
	}, nil
}
