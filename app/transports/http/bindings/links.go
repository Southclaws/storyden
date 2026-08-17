package bindings

import (
	"context"
	"net/url"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/link"
	"github.com/Southclaws/storyden/app/resources/link/link_querier"
	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/app/resources/link/link_writer"

	"github.com/Southclaws/storyden/app/services/link/fetcher"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Links struct {
	fetcher     *fetcher.Fetcher
	linkWriter  *link_writer.LinkWriter
	linkQuerier *link_querier.LinkQuerier
}

func NewLinks(
	fetcher *fetcher.Fetcher,
	linkWriter *link_writer.LinkWriter,
	linkQuerier *link_querier.LinkQuerier,
) Links {
	return Links{
		fetcher:     fetcher,
		linkWriter:  linkWriter,
		linkQuerier: linkQuerier,
	}
}

func (i *Links) LinkCreate(ctx context.Context, request openapi.LinkCreateRequestObject) (openapi.LinkCreateResponseObject, error) {
	u, err := url.Parse(request.Body.Url)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	link, err := i.fetcher.Fetch(ctx, *u)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx),
			fmsg.WithDesc("failed to fetch link",
				"The URL could not be fetched. It may be invalid or the server may be unreachable.",
			), ftag.With(ftag.InvalidArgument))
	}

	return openapi.LinkCreate200JSONResponse{
		LinkCreateOKJSONResponse: openapi.LinkCreateOKJSONResponse(serialiseLinkRef(*link)),
	}, nil
}

func (i *Links) LinkList(ctx context.Context, request openapi.LinkListRequestObject) (openapi.LinkListResponseObject, error) {
	params := deserialisePageParams(request.Params.Page, 50)

	opts := []link_querier.Filter{}

	if v := request.Params.Q; v != nil {
		opts = append(opts, link_querier.WithKeyword(*v))
	}

	result, err := i.linkQuerier.Search(ctx, params, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.LinkList200JSONResponse{
		LinkListOKJSONResponse: openapi.LinkListOKJSONResponse{
			PageSize:    result.Size,
			Results:     result.Results,
			TotalPages:  result.TotalPages,
			CurrentPage: result.CurrentPage,
			NextPage:    result.NextPage.Ptr(),
			Links:       serialiseLinkRefs(result.Items),
		},
	}, nil
}

func (i *Links) LinkGet(ctx context.Context, request openapi.LinkGetRequestObject) (openapi.LinkGetResponseObject, error) {
	l, err := i.linkQuerier.Get(ctx, request.LinkSlug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.LinkGet200JSONResponse{
		LinkGetOKJSONResponse: openapi.LinkGetOKJSONResponse(serialiseLink(l)),
	}, nil
}

func serialiseLink(in *link.Link) openapi.Link {
	return openapi.Link{
		Id:             in.ID.String(),
		CreatedAt:      in.CreatedAt,
		UpdatedAt:      in.UpdatedAt,
		Url:            in.URL,
		Slug:           in.Slug,
		Domain:         in.Domain,
		Title:          in.Title.Ptr(),
		Description:    in.Description.Ptr(),
		FaviconImage:   opt.Map(in.FaviconImage, serialiseAsset).Ptr(),
		PrimaryImage:   opt.Map(in.PrimaryImage, serialiseAsset).Ptr(),
		Assets:         dt.Map(in.Assets, serialiseAssetPtr),
		Nodes:          dt.Map(in.Nodes, serialiseNode),
		Posts:          dt.Map(in.Posts, serialisePostRef),
		Recomentations: dt.Map(in.Related, serialiseDatagraphItem),
	}
}

func serialiseLinkRef(in link_ref.LinkRef) openapi.LinkReference {
	return openapi.LinkReference{
		Id:           in.ID.String(),
		CreatedAt:    in.CreatedAt,
		UpdatedAt:    in.UpdatedAt,
		Slug:         in.Slug,
		Url:          in.URL,
		Domain:       in.Domain,
		Title:        in.Title.Ptr(),
		Description:  in.Description.Ptr(),
		FaviconImage: opt.Map(in.FaviconImage, serialiseAsset).Ptr(),
		PrimaryImage: opt.Map(in.PrimaryImage, serialiseAsset).Ptr(),
	}
}

func serialiseLinkRefPtr(in *link_ref.LinkRef) openapi.LinkReference {
	return serialiseLinkRef(*in)
}

func serialiseLinkRefs(in link_ref.LinkRefs) []openapi.LinkReference {
	return dt.Map(in, serialiseLinkRefPtr)
}
