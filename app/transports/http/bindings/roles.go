package bindings

import (
	"context"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Roles struct {
	//
}

func NewRoles() Roles {
	return Roles{}
}

func (h *Roles) RoleCreate(ctx context.Context, request openapi.RoleCreateRequestObject) (openapi.RoleCreateResponseObject, error) {
	return nil, nil
}

func (h *Roles) RoleGet(ctx context.Context, request openapi.RoleGetRequestObject) (openapi.RoleGetResponseObject, error) {
	return nil, nil
}

func (h *Roles) RoleUpdate(ctx context.Context, request openapi.RoleUpdateRequestObject) (openapi.RoleUpdateResponseObject, error) {
	return nil, nil
}
