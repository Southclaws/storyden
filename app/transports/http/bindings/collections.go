package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/collection"
	"github.com/Southclaws/storyden/app/resources/collection/collection_item_status"
	"github.com/Southclaws/storyden/app/resources/collection/collection_querier"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/collection/collection_item_manager"
	"github.com/Southclaws/storyden/app/services/collection/collection_manager"
	"github.com/Southclaws/storyden/app/services/collection/collection_read"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Collections struct {
	colQuerier     *collection_querier.Querier
	colReader      *collection_read.Hydrator
	colManager     *collection_manager.Manager
	colItemManager *collection_item_manager.Manager
}

func NewCollections(
	colQuerier *collection_querier.Querier,
	colReader *collection_read.Hydrator,
	colManager *collection_manager.Manager,
	colItemManager *collection_item_manager.Manager,
) Collections {
	return Collections{
		colQuerier:     colQuerier,
		colReader:      colReader,
		colManager:     colManager,
		colItemManager: colItemManager,
	}
}

func (i *Collections) CollectionCreate(ctx context.Context, request openapi.CollectionCreateRequestObject) (openapi.CollectionCreateResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	coll, err := i.colManager.Create(ctx, accountID, request.Body.Name, collection_manager.Partial{
		Description: opt.NewPtr(request.Body.Description),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CollectionCreate200JSONResponse{
		CollectionCreateOKJSONResponse: openapi.CollectionCreateOKJSONResponse(serialiseCollection(&coll.Collection)),
	}, nil
}

func (i *Collections) CollectionList(ctx context.Context, request openapi.CollectionListRequestObject) (openapi.CollectionListResponseObject, error) {
	opts := []collection_querier.Option{}

	if v := request.Params.AccountHandle; v != nil {
		opts = append(opts, collection_querier.WithOwnerHandle(*v))
	}

	itemPresenceQuery := opt.Map(opt.NewPtr(request.Params.HasItem), deserialiseID)
	if v, ok := itemPresenceQuery.Get(); ok {
		opts = append(opts, collection_querier.WithItemPresenceQuery(v))
	}

	colls, err := i.colQuerier.List(ctx, opts...)
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
	coll, err := i.colReader.GetCollection(ctx, collection.NewKey(request.CollectionMark))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CollectionGet200JSONResponse{
		CollectionGetOKJSONResponse: openapi.CollectionGetOKJSONResponse(serialiseCollectionWithItems(coll)),
	}, nil
}

func (i *Collections) CollectionUpdate(ctx context.Context, request openapi.CollectionUpdateRequestObject) (openapi.CollectionUpdateResponseObject, error) {
	c, err := i.colManager.Update(ctx,
		collection.NewKey(request.CollectionMark),
		collection_manager.Partial{
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
	err := i.colManager.Delete(ctx, collection.NewKey(request.CollectionMark))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CollectionDelete200Response{}, nil
}

func (i *Collections) CollectionAddPost(ctx context.Context, request openapi.CollectionAddPostRequestObject) (openapi.CollectionAddPostResponseObject, error) {
	c, err := i.colItemManager.PostAdd(ctx,
		collection.NewKey(request.CollectionMark),
		post.ID(deserialiseID(request.PostId)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CollectionAddPost200JSONResponse{
		CollectionAddPostOKJSONResponse: openapi.CollectionAddPostOKJSONResponse(serialiseCollectionWithItems(c)),
	}, nil
}

func (i *Collections) CollectionRemovePost(ctx context.Context, request openapi.CollectionRemovePostRequestObject) (openapi.CollectionRemovePostResponseObject, error) {
	c, err := i.colItemManager.PostRemove(ctx,
		collection.NewKey(request.CollectionMark),
		post.ID(deserialiseID(request.PostId)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CollectionRemovePost200JSONResponse{
		CollectionRemovePostOKJSONResponse: openapi.CollectionRemovePostOKJSONResponse(serialiseCollectionWithItems(c)),
	}, nil
}

func (i *Collections) CollectionAddNode(ctx context.Context, request openapi.CollectionAddNodeRequestObject) (openapi.CollectionAddNodeResponseObject, error) {
	c, err := i.colItemManager.NodeAdd(ctx,
		collection.NewKey(request.CollectionMark),
		library.NodeID(deserialiseID(request.NodeId)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.CollectionAddNode200JSONResponse{
		CollectionAddNodeOKJSONResponse: openapi.CollectionAddNodeOKJSONResponse(serialiseCollectionWithItems(c)),
	}, nil
}

func (i *Collections) CollectionRemoveNode(ctx context.Context, request openapi.CollectionRemoveNodeRequestObject) (openapi.CollectionRemoveNodeResponseObject, error) {
	c, err := i.colItemManager.NodeRemove(ctx,
		collection.NewKey(request.CollectionMark),
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
		Id:             in.Mark.ID().String(),
		CreatedAt:      in.CreatedAt,
		UpdatedAt:      in.UpdatedAt,
		Name:           in.Name,
		Slug:           in.Mark.String(),
		Description:    in.Description.Ptr(),
		Owner:          serialiseProfileReference(in.Owner),
		ItemCount:      int(in.ItemCount),
		HasQueriedItem: in.HasQueriedItem,
	}
}

func serialiseCollectionWithItems(in *collection.CollectionWithItems) openapi.CollectionWithItems {
	return openapi.CollectionWithItems{
		Id:          in.Mark.ID().String(),
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
		Name:        in.Name,
		Slug:        in.Mark.String(),
		Description: in.Description.Ptr(),
		Owner:       serialiseProfileReference(in.Owner),
		Items:       dt.Map(in.Items, serialiseCollectionItem),
	}
}

func serialiseCollectionItem(in *collection.CollectionItem) openapi.CollectionItem {
	score := opt.PtrMap(in.RelevanceScore, func(s float64) float32 { return float32(s) })

	return openapi.CollectionItem{
		Id:             in.Item.GetID().String(),
		CreatedAt:      in.Item.GetCreated(),
		UpdatedAt:      in.Item.GetUpdated(),
		Owner:          serialiseProfileReference(in.Author), // Invalid, wrong owner
		AddedAt:        in.Added,
		MembershipType: openapi.CollectionItemMembershipType(in.MembershipType.String()),
		RelevanceScore: score,
		Item:           serialiseDatagraphItem(in.Item),
	}
}

func serialiseCollectionStatus(in collection_item_status.Status) openapi.CollectionStatus {
	return openapi.CollectionStatus{
		InCollections: in.Count,
		HasCollected:  in.Status,
	}
}
