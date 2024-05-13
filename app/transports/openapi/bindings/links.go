package bindings

import (
	"context"
	"strconv"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/datagraph/link"
	"github.com/Southclaws/storyden/app/resources/datagraph/link_graph"
	"github.com/Southclaws/storyden/app/services/hydrator/fetcher"
	"github.com/Southclaws/storyden/app/services/link_getter"
	"github.com/Southclaws/storyden/app/transports/openapi"
)

type Links struct {
	fr fetcher.Service
	lr link.Repository
	lg *link_getter.Getter
}

func NewLinks(
	fr fetcher.Service,
	lr link.Repository,
	lg *link_getter.Getter,
) Links {
	return Links{
		fr: fr,
		lr: lr,
		lg: lg,
	}
}

func (i *Links) LinkCreate(ctx context.Context, request openapi.LinkCreateRequestObject) (openapi.LinkCreateResponseObject, error) {
	link, err := i.fr.Fetch(ctx,
		request.Body.Url,
	)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.LinkCreate200JSONResponse{
		LinkCreateOKJSONResponse: openapi.LinkCreateOKJSONResponse(serialiseLink(link)),
	}, nil
}

func (i *Links) LinkList(ctx context.Context, request openapi.LinkListRequestObject) (openapi.LinkListResponseObject, error) {
	pageSize := 50

	page := opt.NewPtrMap(request.Params.Page, func(s string) int {
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return 0
		}

		return max(1, int(v))
	}).Or(1)

	opts := []link.Filter{}

	if v := request.Params.Q; v != nil {
		opts = append(opts, link.WithKeyword(*v))
	}

	// API is 1-indexed, internally it's 0-indexed.
	page = max(0, page-1)

	result, err := i.lr.Search(ctx, page, pageSize, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// API is 1-indexed, internally it's 0-indexed.
	page = result.CurrentPage + 1

	return openapi.LinkList200JSONResponse{
		LinkListOKJSONResponse: openapi.LinkListOKJSONResponse{
			PageSize:    pageSize,
			Results:     result.Results,
			TotalPages:  result.TotalPages,
			CurrentPage: page,
			NextPage:    result.NextPage.Ptr(),
			Links:       dt.Map(result.Links, serialiseLink),
		},
	}, nil
}

func (i *Links) LinkGet(ctx context.Context, request openapi.LinkGetRequestObject) (openapi.LinkGetResponseObject, error) {
	l, err := i.lg.Get(ctx, request.LinkSlug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.LinkGet200JSONResponse{
		LinkGetOKJSONResponse: openapi.LinkGetOKJSONResponse(serialiseLinkWithRefs(l)),
	}, nil
}

func serialiseLinkWithRefs(in *link_graph.WithRefs) openapi.LinkWithRefs {
	return openapi.LinkWithRefs{
		Url:            in.URL,
		Title:          in.Title.Ptr(),
		Description:    in.Description.Ptr(),
		Slug:           in.Slug,
		Domain:         in.Domain,
		Assets:         dt.Map(in.Assets, serialiseAssetReference),
		Threads:        dt.Map(in.Threads, serialiseThreadReference),
		Posts:          dt.Map(in.Replies, serialisePost),
		Clusters:       dt.Map(in.Clusters, serialiseCluster),
		Recomentations: dt.Map(in.Related, serialiseDatagraphNodeReference),
	}
}
