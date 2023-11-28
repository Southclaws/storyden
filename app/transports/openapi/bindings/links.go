package bindings

import (
	"context"
	"strconv"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/link"
	"github.com/Southclaws/storyden/app/resources/link_graph"
	"github.com/Southclaws/storyden/app/services/hydrator/fetcher"
	"github.com/Southclaws/storyden/internal/openapi"
)

type Links struct {
	fr fetcher.Service
	lr link.Repository
	lg link_graph.Repository
}

func NewLinks(
	fr fetcher.Service,
	lr link.Repository,
	lg link_graph.Repository,
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
	opts := []link.Filter{}

	if v := request.Params.Q; v != nil {
		opts = append(opts, link.WithKeyword(*v))
	}

	if v := request.Params.Page; v != nil {
		pageNumber, err := strconv.ParseInt(*v, 10, 32)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}

		opts = append(opts, link.WithPage(int(pageNumber-1), 50))
	} else {
		opts = append(opts, link.WithPage(0, 50))
	}

	links, err := i.lr.Search(ctx, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.LinkList200JSONResponse{
		LinkListOKJSONResponse: openapi.LinkListOKJSONResponse{
			Links: dt.Map(links, serialiseLink),
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
		Url:         in.URL,
		Title:       in.Title.Ptr(),
		Description: in.Description.Ptr(),
		Slug:        in.Slug,
		Domain:      in.Domain,
		Assets:      dt.Map(in.Assets, serialiseAssetReference),
		Threads:     dt.Map(in.Threads, serialiseThreadReference),
		Posts:       dt.Map(in.Replies, serialisePost),
		Clusters:    dt.Map(in.Clusters, serialiseCluster),
		Items:       dt.Map(in.Items, serialiseItemWithParents),
	}
}
