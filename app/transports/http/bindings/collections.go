package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	collection_svc "github.com/Southclaws/storyden/app/services/collection"
	"github.com/Southclaws/storyden/app/services/collection/collection_read"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Collections struct {
	collection_repo    collection.Repository
	collection_svc     collection_svc.Service
	collection_querier collection_read.CollectionQuerier
}

func NewCollections(
	collection_repo collection.Repository,
	collection_svc collection_svc.Service,
	collection_querier collection_read.CollectionQuerier,
) Collections {
	return Collections{
		collection_repo:    collection_repo,
		collection_svc:     collection_svc,
		collection_querier: collection_querier,
	}
}

func (i *Collections) CollectionCreate(ctx context.Context, request openapi.CollectionCreateRequestObject) (openapi.CollectionCreateResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	coll, err := i.collection_repo.Create(ctx, accountID, request.Body.Name, request.Body.Description)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CollectionCreate200JSONResponse{
		CollectionCreateOKJSONResponse: openapi.CollectionCreateOKJSONResponse(serialiseCollection(&coll.Collection)),
	}, nil
}

func (i *Collections) CollectionList(ctx context.Context, request openapi.CollectionListRequestObject) (openapi.CollectionListResponseObject, error) {
	filters := []collection.Filter{}

	if v := request.Params.AccountHandle; v != nil {
		filters = append(filters, collection.WithOwnerHandle(*v))
	}

	colls, err := i.collection_repo.List(ctx, filters...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	list := dt.Map(colls, serialiseCollection)

	return openapi.CollectionList200JSONResponse{
		CollectionListOKJSONResponse: openapi.CollectionListOKJSONResponse{
			Collections: list,
		},
	}, nil
}

func (i *Collections) CollectionGet(ctx context.Context, request openapi.CollectionGetRequestObject) (openapi.CollectionGetResponseObject, error) {
	coll, err := i.collection_querier.GetCollection(ctx, collection.CollectionID(deserialiseID(request.CollectionId)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CollectionGet200JSONResponse{
		CollectionGetOKJSONResponse: openapi.CollectionGetOKJSONResponse(serialiseCollectionWithItems(coll)),
	}, nil
}

func (i *Collections) CollectionUpdate(ctx context.Context, request openapi.CollectionUpdateRequestObject) (openapi.CollectionUpdateResponseObject, error) {
	c, err := i.collection_svc.Update(ctx,
		collection.CollectionID(deserialiseID(request.CollectionId)),
		collection_svc.Partial{
			Name:        opt.NewPtr(request.Body.Name),
			Description: opt.NewPtr(request.Body.Description),
		})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CollectionUpdate200JSONResponse{
		CollectionUpdateOKJSONResponse: openapi.CollectionUpdateOKJSONResponse(serialiseCollection(&c.Collection)),
	}, nil
}

func (i *Collections) CollectionDelete(ctx context.Context, request openapi.CollectionDeleteRequestObject) (openapi.CollectionDeleteResponseObject, error) {
	err := i.collection_svc.Delete(ctx, collection.CollectionID(deserialiseID(request.CollectionId)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CollectionDelete200Response{}, nil
}

func (i *Collections) CollectionAddPost(ctx context.Context, request openapi.CollectionAddPostRequestObject) (openapi.CollectionAddPostResponseObject, error) {
	c, err := i.collection_svc.PostAdd(ctx,
		collection.CollectionID(deserialiseID(request.CollectionId)),
		post.ID(deserialiseID(request.PostId)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CollectionAddPost200JSONResponse{
		CollectionAddPostOKJSONResponse: openapi.CollectionAddPostOKJSONResponse(serialiseCollectionWithItems(c)),
	}, nil
}

func (i *Collections) CollectionRemovePost(ctx context.Context, request openapi.CollectionRemovePostRequestObject) (openapi.CollectionRemovePostResponseObject, error) {
	c, err := i.collection_svc.PostRemove(ctx,
		collection.CollectionID(deserialiseID(request.CollectionId)),
		post.ID(deserialiseID(request.PostId)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CollectionRemovePost200JSONResponse{
		CollectionRemovePostOKJSONResponse: openapi.CollectionRemovePostOKJSONResponse(serialiseCollectionWithItems(c)),
	}, nil
}

func (i *Collections) CollectionAddNode(ctx context.Context, request openapi.CollectionAddNodeRequestObject) (openapi.CollectionAddNodeResponseObject, error) {
	c, err := i.collection_svc.NodeAdd(ctx,
		collection.CollectionID(deserialiseID(request.CollectionId)),
		library.NodeID(deserialiseID(request.NodeId)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CollectionAddNode200JSONResponse{
		CollectionAddNodeOKJSONResponse: openapi.CollectionAddNodeOKJSONResponse(serialiseCollectionWithItems(c)),
	}, nil
}

func (i *Collections) CollectionRemoveNode(ctx context.Context, request openapi.CollectionRemoveNodeRequestObject) (openapi.CollectionRemoveNodeResponseObject, error) {
	c, err := i.collection_svc.NodeRemove(ctx,
		collection.CollectionID(deserialiseID(request.CollectionId)),
		library.NodeID(deserialiseID(request.NodeId)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CollectionRemoveNode200JSONResponse{
		CollectionRemoveNodeOKJSONResponse: openapi.CollectionRemoveNodeOKJSONResponse(serialiseCollectionWithItems(c)),
	}, nil
}

func serialiseCollection(in *collection.Collection) openapi.Collection {
	return openapi.Collection{
		Id:          in.ID.String(),
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
		Owner:       serialiseProfileReference(in.Owner),
		Name:        in.Name,
		Description: in.Description,
	}
}

func serialiseCollectionWithItems(in *collection.CollectionWithItems) openapi.CollectionWithItems {
	return openapi.CollectionWithItems{
		Id:          in.ID.String(),
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
		Owner:       serialiseProfileReference(in.Owner),
		Name:        in.Name,
		Description: in.Description,
		Items:       dt.Map(in.Items, serialiseCollectionItem),
	}
}

func serialiseCollectionItem(in *collection.CollectionItem) openapi.CollectionItem {
	score := opt.PtrMap(in.RelevanceScore, func(s float64) float32 { return float32(s) })
	meta := in.Item.GetProps()

	return openapi.CollectionItem{
		Id:             in.Item.GetID().String(),
		AddedAt:        in.Added,
		CreatedAt:      in.Item.GetCreated(),
		UpdatedAt:      in.Item.GetUpdated(),
		MembershipType: openapi.CollectionItemMembershipType(in.MembershipType.String()),
		Owner:          serialiseProfileReference(in.Author), // Invalid, wrong owner
		Kind:           serialiseDatagraphKind(in.Item.GetKind()),
		Name:           in.Item.GetName(),
		Slug:           in.Item.GetSlug(),
		Description:    opt.New(in.Item.GetDesc()).Ptr(),
		RelevanceScore: score,
		Meta:           (*openapi.Metadata)(&meta),
		Assets:         dt.Map(in.Item.GetAssets(), serialiseAssetPtr),
	}
}
