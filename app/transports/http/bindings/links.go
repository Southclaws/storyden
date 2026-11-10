package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/link"
	"github.com/Southclaws/storyden/app/resources/link/link_querier"
	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/app/resources/link/link_url"
	"github.com/Southclaws/storyden/app/resources/link/link_writer"

	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/services/authentication/session"
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
	if err := session.Authorise(ctx, nil, rbac.PermissionUseLinks, rbac.PermissionManageLinks); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.WithDesc("missing link permission", "You need the Use Links permission to add a link. Ask an administrator to grant it."))
	}

	u, err := link_url.Parse(request.Body.Url)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	link, err := i.fetcher.Fetch(ctx, u)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.LinkCreate200JSONResponse{
		LinkCreateOKJSONResponse: openapi.LinkCreateOKJSONResponse(serialiseLinkRef(*link)),
	}, nil
}

func (i *Links) LinkList(ctx context.Context, request openapi.LinkListRequestObject) (openapi.LinkListResponseObject, error) {
	if err := session.Authorise(ctx, nil, rbac.PermissionViewLinks, rbac.PermissionUseLinks, rbac.PermissionManageLinks); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

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

func (i *Links) LinkDelete(ctx context.Context, request openapi.LinkDeleteRequestObject) (openapi.LinkDeleteResponseObject, error) {
	if err := i.linkWriter.Delete(ctx, request.LinkSlug); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.LinkDelete204Response{}, nil
}

func serialiseLink(in *link.Link) openapi.Link {
	return openapi.Link{
		Id:                in.ID.String(),
		CreatedAt:         in.CreatedAt,
		UpdatedAt:         in.UpdatedAt,
		Url:               in.URL.String(),
		Slug:              in.Slug,
		Domain:            in.Domain,
		Title:             in.Title.Ptr(),
		Description:       in.Description.Ptr(),
		LivenessStatus:    opt.New(openapi.LinkLivenessStatus(in.LivenessStatus)).Ptr(),
		LivenessCheckedAt: in.LivenessCheckedAt,
		ScrapeStatus:      opt.New(openapi.LinkScrapeStatus(in.ScrapeStatus)).Ptr(),
		ScrapeAttemptedAt: in.ScrapeAttemptedAt,
		ScrapedAt:         in.ScrapedAt,
		FaviconImage:      opt.Map(in.FaviconImage, serialiseAsset).Ptr(),
		PrimaryImage:      opt.Map(in.PrimaryImage, serialiseAsset).Ptr(),
		Assets:            dt.Map(in.Assets, serialiseAssetPtr),
		Nodes:             dt.Map(in.Nodes, serialiseNode),
		Posts:             dt.Map(in.Posts, serialisePostRef),
		Recomentations:    dt.Map(in.Related, serialiseDatagraphItem),
	}
}

func serialiseLinkRef(in link_ref.LinkRef) openapi.LinkReference {
	return openapi.LinkReference{
		Id:                in.ID.String(),
		CreatedAt:         in.CreatedAt,
		UpdatedAt:         in.UpdatedAt,
		Slug:              in.Slug,
		Url:               in.URL.String(),
		Domain:            in.Domain,
		Title:             in.Title.Ptr(),
		Description:       in.Description.Ptr(),
		LivenessStatus:    opt.New(openapi.LinkLivenessStatus(in.LivenessStatus)).Ptr(),
		LivenessCheckedAt: in.LivenessCheckedAt,
		ScrapeStatus:      opt.New(openapi.LinkScrapeStatus(in.ScrapeStatus)).Ptr(),
		ScrapeAttemptedAt: in.ScrapeAttemptedAt,
		ScrapedAt:         in.ScrapedAt,
		FaviconImage:      opt.Map(in.FaviconImage, serialiseAsset).Ptr(),
		PrimaryImage:      opt.Map(in.PrimaryImage, serialiseAsset).Ptr(),
	}
}

func serialiseLinkRefPtr(in *link_ref.LinkRef) openapi.LinkReference {
	return serialiseLinkRef(*in)
}

func serialiseLinkRefs(in link_ref.LinkRefs) []openapi.LinkReference {
	return dt.Map(in, serialiseLinkRefPtr)
}
