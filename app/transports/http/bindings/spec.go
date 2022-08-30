package bindings

import (
	"context"
	"encoding/json"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/errtag"
)

type Spec struct{}

func NewSpec() Spec { return Spec{} }

func (v *Spec) GetSpec(ctx context.Context, request openapi.GetSpecRequestObject) any {
	spec, err := openapi.GetSwagger()
	if err != nil {
		return errtag.Wrap(err, errtag.Internal{})
	}

	b, err := json.Marshal(spec)
	if err != nil {
		return errtag.Wrap(err, errtag.Internal{})
	}

	return openapi.GetSpec200TextResponse(b)
}
