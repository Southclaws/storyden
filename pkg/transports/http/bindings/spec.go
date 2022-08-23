package bindings

import (
	"context"
	"encoding/json"

	"github.com/Southclaws/storyden/pkg/transports/http/openapi"
)

type Spec struct{}

func NewSpec() Spec { return Spec{} }

func (v *Spec) GetSpec(ctx context.Context, request openapi.GetSpecRequestObject) any {
	spec, err := openapi.GetSwagger()
	if err != nil {
		return err
	}

	b, err := json.Marshal(spec)
	if err != nil {
		return err
	}

	return openapi.GetSpec200TextResponse(b)
}
