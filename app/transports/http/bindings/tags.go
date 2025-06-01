package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/app/resources/tag/tag_querier"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Tags struct {
	tagQuerier *tag_querier.Querier
}

func NewTags(tagQuerier *tag_querier.Querier) Tags {
	return Tags{tagQuerier: tagQuerier}
}

func (h Tags) TagList(ctx context.Context, request openapi.TagListRequestObject) (openapi.TagListResponseObject, error) {
	fn := func() (tag_ref.Tags, error) {
		if request.Params.Q == nil {
			return h.tagQuerier.List(ctx)
		} else {
			return h.tagQuerier.Search(ctx, *request.Params.Q)
		}
	}

	list, err := fn()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.TagList200JSONResponse{
		TagListOKJSONResponse: openapi.TagListOKJSONResponse{
			Tags: serialiseTagReferenceList(list),
		},
	}, nil
}

func (h Tags) TagGet(ctx context.Context, request openapi.TagGetRequestObject) (openapi.TagGetResponseObject, error) {
	name := tag_ref.NewName(request.TagName)

	tag, err := h.tagQuerier.Get(ctx, name)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.TagGet200JSONResponse{
		TagGetOKJSONResponse: openapi.TagGetOKJSONResponse(serialiseTag(tag)),
	}, nil
}

func serialiseTag(in *tag.Tag) openapi.Tag {
	return openapi.Tag{
		Id:        in.ID.String(),
		Name:      in.Name.String(),
		Colour:    in.Colour,
		ItemCount: in.ItemCount,
		Items:     serialiseDatagraphItemList(in.Items),
	}
}

func serialiseTagReference(in *tag_ref.Tag) openapi.TagReference {
	return openapi.TagReference{
		Name:      in.Name.String(),
		Colour:    in.Colour,
		ItemCount: in.ItemCount,
	}
}

func serialiseTagReferenceList(in tag_ref.Tags) []openapi.TagReference {
	return dt.Map(in, serialiseTagReference)
}

func deserialiseTagName(in string) tag_ref.Name {
	return tag_ref.NewName(in)
}

func tagsIDs(i openapi.TagListIDs) []xid.ID {
	return dt.Map(i, deserialiseID)
}
