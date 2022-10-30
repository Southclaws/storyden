package bindings

import (
	"context"
	"encoding/json"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/transports/openapi/openapi"
)

type Spec struct{}

func NewSpec() Spec { return Spec{} }

func (v *Spec) GetSpec(ctx context.Context, request openapi.GetSpecRequestObject) (openapi.GetSpecResponseObject, error) {
	spec, err := openapi.GetSwagger()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	b, err := json.Marshal(spec)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.GetSpec200TextResponse(b), nil
}
