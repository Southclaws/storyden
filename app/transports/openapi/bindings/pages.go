package bindings

import (
	"context"

	"github.com/Southclaws/storyden/internal/openapi"
)

type Pages struct{}

func NewPages() Pages {
	return Pages{}
}

func (p *Pages) PageList(ctx context.Context, request openapi.PageListRequestObject) (openapi.PageListResponseObject, error) {
	return nil, nil
}

func (p *Pages) PageUpdateOrder(ctx context.Context, request openapi.PageUpdateOrderRequestObject) (openapi.PageUpdateOrderResponseObject, error) {
	return nil, nil
}

func (p *Pages) PageCreate(ctx context.Context, request openapi.PageCreateRequestObject) (openapi.PageCreateResponseObject, error) {
	return nil, nil
}

func (p *Pages) PageUpdate(ctx context.Context, request openapi.PageUpdateRequestObject) (openapi.PageUpdateResponseObject, error) {
	return nil, nil
}
