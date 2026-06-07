package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/oapi-codegen/nullable"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library/node_version"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/services/library/node_versioning"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/deletable"
)

func (c *Nodes) NodeDraftList(ctx context.Context, request openapi.NodeDraftListRequestObject) (openapi.NodeDraftListResponseObject, error) {
	page := deserialisePageParams(request.Params.Page, 50)

	drafts, err := c.versioning.ListAllDrafts(ctx, page)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeDraftList200JSONResponse{
		NodeDraftListOKJSONResponse: openapi.NodeDraftListOKJSONResponse{
			CurrentPage: drafts.CurrentPage,
			NextPage:    drafts.NextPage.Ptr(),
			PageSize:    drafts.Size,
			Results:     drafts.Results,
			TotalPages:  drafts.TotalPages,
			Drafts:      dt.Map(drafts.Items, serialiseNodeDraft),
		},
	}, nil
}

func (c *Nodes) NodeVersionList(ctx context.Context, request openapi.NodeVersionListRequestObject) (openapi.NodeVersionListResponseObject, error) {
	qk := deserialiseNodeMark(request.NodeSlug)
	page := deserialisePageParams(request.Params.Page, 50)

	versions, err := c.versioning.ListVisible(ctx, qk, page)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeVersionList200JSONResponse{
		NodeVersionListOKJSONResponse: openapi.NodeVersionListOKJSONResponse{
			CurrentPage: versions.CurrentPage,
			NextPage:    versions.NextPage.Ptr(),
			PageSize:    versions.Size,
			Results:     versions.Results,
			TotalPages:  versions.TotalPages,
			Versions:    dt.Map(versions.Items, serialiseNodeVersion),
		},
	}, nil
}

func (c *Nodes) NodeVersionCreate(ctx context.Context, request openapi.NodeVersionCreateRequestObject) (openapi.NodeVersionCreateResponseObject, error) {
	qk := deserialiseNodeMark(request.NodeSlug)

	partial, err := deserialiseNodeVersionInitialProps(*request.Body)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	version, err := c.versioning.CreateDraft(ctx, qk, partial)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeVersionCreate200JSONResponse{
		NodeVersionCreateOKJSONResponse: openapi.NodeVersionCreateOKJSONResponse(serialiseNodeVersion(version)),
	}, nil
}

func (c *Nodes) NodeVersionDraftGet(ctx context.Context, request openapi.NodeVersionDraftGetRequestObject) (openapi.NodeVersionDraftGetResponseObject, error) {
	qk := deserialiseNodeMark(request.NodeSlug)

	version, err := c.versioning.GetVisibleDraft(ctx, qk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeVersionDraftGet200JSONResponse{
		NodeVersionGetOKJSONResponse: openapi.NodeVersionGetOKJSONResponse(serialiseNodeVersion(version)),
	}, nil
}

func (c *Nodes) NodeVersionDraftUpdate(ctx context.Context, request openapi.NodeVersionDraftUpdateRequestObject) (openapi.NodeVersionDraftUpdateResponseObject, error) {
	qk := deserialiseNodeMark(request.NodeSlug)

	partial, err := deserialiseNodeVersionMutableProps(*request.Body)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	version, err := c.versioning.UpdateVisibleDraft(ctx, qk, partial)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeVersionDraftUpdate200JSONResponse{
		NodeVersionUpdateOKJSONResponse: openapi.NodeVersionUpdateOKJSONResponse(serialiseNodeVersion(version)),
	}, nil
}

func (c *Nodes) NodeVersionGet(ctx context.Context, request openapi.NodeVersionGetRequestObject) (openapi.NodeVersionGetResponseObject, error) {
	qk := deserialiseNodeMark(request.NodeSlug)

	versionID, err := node_version.VersionIDFromString(request.VersionId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	version, err := c.versioning.GetVisible(ctx, qk, versionID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeVersionGet200JSONResponse{
		NodeVersionGetOKJSONResponse: openapi.NodeVersionGetOKJSONResponse(serialiseNodeVersion(version)),
	}, nil
}

func (c *Nodes) NodeVersionUpdate(ctx context.Context, request openapi.NodeVersionUpdateRequestObject) (openapi.NodeVersionUpdateResponseObject, error) {
	qk := deserialiseNodeMark(request.NodeSlug)

	versionID, err := node_version.VersionIDFromString(request.VersionId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	partial, err := deserialiseNodeVersionMutableProps(*request.Body)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	version, err := c.versioning.UpdateDraft(ctx, qk, versionID, partial)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeVersionUpdate200JSONResponse{
		NodeVersionUpdateOKJSONResponse: openapi.NodeVersionUpdateOKJSONResponse(serialiseNodeVersion(version)),
	}, nil
}

func (c *Nodes) NodeVersionDelete(ctx context.Context, request openapi.NodeVersionDeleteRequestObject) (openapi.NodeVersionDeleteResponseObject, error) {
	qk := deserialiseNodeMark(request.NodeSlug)

	versionID, err := node_version.VersionIDFromString(request.VersionId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	if err := c.versioning.DiscardDraft(ctx, qk, versionID); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeVersionDelete200Response{}, nil
}

func (c *Nodes) NodeVersionUpdateStatus(ctx context.Context, request openapi.NodeVersionUpdateStatusRequestObject) (openapi.NodeVersionUpdateStatusResponseObject, error) {
	qk := deserialiseNodeMark(request.NodeSlug)

	if request.Body.Status != openapi.NodeVersionStatusApplied {
		return nil, fault.New("unsupported version status transition",
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("unsupported status", "Only applying a draft version is currently supported."),
		)
	}

	versionID, err := node_version.VersionIDFromString(request.VersionId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	version, err := c.versioning.ApplyVersion(ctx, qk, versionID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.NodeVersionUpdateStatus200JSONResponse{
		NodeVersionUpdateOKJSONResponse: openapi.NodeVersionUpdateOKJSONResponse(serialiseNodeVersion(version)),
	}, nil
}

func deserialiseNodeVersionInitialProps(in openapi.NodeVersionInitialProps) (node_versioning.DraftPartial, error) {
	return deserialiseNodeVersionPartial(
		in.Name,
		in.Slug,
		in.Description,
		in.Content,
		in.Properties,
		in.Meta,
	)
}

func deserialiseNodeVersionMutableProps(in openapi.NodeVersionMutableProps) (node_versioning.DraftPartial, error) {
	return deserialiseNodeVersionPartial(
		in.Name,
		in.Slug,
		in.Description,
		in.Content,
		in.Properties,
		in.Meta,
	)
}

func deserialiseNodeVersionPartial(
	name *openapi.NodeName,
	slug *openapi.NodeSlug,
	description nullable.Nullable[openapi.NodeDescription],
	content nullable.Nullable[openapi.PostContent],
	properties *openapi.PropertyMutationList,
	meta *openapi.Metadata,
) (node_versioning.DraftPartial, error) {
	richContent, err := deletable.NewMapErr(content, datagraph.NewRichText)
	if err != nil {
		return node_versioning.DraftPartial{}, fault.Wrap(err, ftag.With(ftag.InvalidArgument))
	}

	parsedSlug, err := opt.MapErr(opt.NewPtr(slug), func(s openapi.NodeSlug) (mark.Slug, error) {
		v, err := mark.NewSlug(s)
		if err != nil {
			return mark.Slug{}, err
		}
		return *v, nil
	})
	if err != nil {
		return node_versioning.DraftPartial{}, fault.Wrap(err, ftag.With(ftag.InvalidArgument))
	}

	pml, err := opt.MapErr(opt.NewPtr(properties), deserialisePropertyMutationList)
	if err != nil {
		return node_versioning.DraftPartial{}, fault.Wrap(err, ftag.With(ftag.InvalidArgument))
	}

	return node_versioning.DraftPartial{
		Name:               opt.NewPtr(name),
		Slug:               parsedSlug,
		Description:        deletable.New(description),
		Content:            richContent,
		PropertiesSnapshot: pml,
		Metadata:           opt.NewPtr((*map[string]any)(meta)),
	}, nil
}

func serialiseNodeVersion(in *node_version.NodeVersion) openapi.NodeVersion {
	meta := openapi.Metadata(in.Metadata)

	return openapi.NodeVersion{
		Id:          in.ID.String(),
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
		NodeId:      in.NodeID.String(),
		Author:      serialiseProfileReference(in.Author),
		Status:      openapi.NodeVersionStatus(in.Status.String()),
		Name:        in.Name,
		Slug:        openapi.NodeSlug(in.Slug),
		Description: serialiseNullableOpt(in.Description),
		Content: serialiseNullableOpt(opt.Map(in.Content, func(c datagraph.Content) openapi.PostContent {
			return c.HTML()
		})),
		Properties: opt.Map(in.PropertiesSnapshot, serialisePropertySnapshotList).Or(openapi.PropertyMutationList{}),
		Previous:   opt.Map(in.Previous, serialiseNodeVersionReference).Ptr(),
		Meta:       meta,
	}
}

func serialiseNodeVersionReference(in node_version.VersionReference) openapi.NodeVersionReference {
	return openapi.NodeVersionReference{
		Id:        in.ID.String(),
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
		Author:    serialiseProfileReference(in.Author),
		Status:    openapi.NodeVersionStatus(in.Status.String()),
	}
}

func serialiseNodeDraft(in *node_version.NodeVersionWithNode) openapi.NodeDraft {
	version := serialiseNodeVersion(&in.NodeVersion)
	node := serialiseNode(in.Node)

	return openapi.NodeDraft{
		Author:      version.Author,
		Content:     version.Content,
		CreatedAt:   version.CreatedAt,
		Description: version.Description,
		Id:          version.Id,
		Meta:        version.Meta,
		Name:        version.Name,
		Node:        node,
		NodeId:      version.NodeId,
		Properties:  version.Properties,
		Slug:        version.Slug,
		Status:      version.Status,
		UpdatedAt:   version.UpdatedAt,
	}
}

func serialisePropertySnapshotList(in []node_version.PropertySnapshot) openapi.PropertyMutationList {
	return dt.Map(in, serialisePropertySnapshot)
}

func serialisePropertySnapshot(in node_version.PropertySnapshot) openapi.PropertyMutation {
	var fid *openapi.Identifier
	if id, ok := in.ID.Get(); ok {
		v := openapi.Identifier(id.String())
		fid = &v
	}

	var typ *openapi.PropertyType
	if t, ok := in.Type.Get(); ok {
		v := openapi.PropertyType(t.String())
		typ = &v
	}

	var sort *openapi.PropertySortKey
	if s, ok := in.Sort.Get(); ok {
		v := openapi.PropertySortKey(s)
		sort = &v
	}

	return openapi.PropertyMutation{
		Fid:   fid,
		Name:  in.Name,
		Type:  typ,
		Sort:  sort,
		Value: in.Value,
	}
}
