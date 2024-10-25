package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

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

	tags := dt.Map(list, serialiseTag)

	return openapi.TagList200JSONResponse{
		TagListOKJSONResponse: openapi.TagListOKJSONResponse{
			Tags: tags,
		},
	}, nil
}
