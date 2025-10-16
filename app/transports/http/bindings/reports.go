package bindings

import (
	"context"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Reports struct{}

func NewReports() Reports {
	return Reports{}
}

func (h *Reports) ReportCreate(ctx context.Context, request openapi.ReportCreateRequestObject) (openapi.ReportCreateResponseObject, error) {
	return nil, nil
}

func (h *Reports) ReportList(ctx context.Context, request openapi.ReportListRequestObject) (openapi.ReportListResponseObject, error) {
	return nil, nil
}

func (h *Reports) ReportUpdate(ctx context.Context, request openapi.ReportUpdateRequestObject) (openapi.ReportUpdateResponseObject, error) {
	return nil, nil
}
