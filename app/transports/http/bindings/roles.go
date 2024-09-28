package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/app/resources/account/role/role_querier"
	"github.com/Southclaws/storyden/app/resources/account/role/role_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Roles struct {
	accountQuerier *account_querier.Querier
	roleQuerier    *role_querier.Querier
	roleWriter     *role_writer.Writer
}

func NewRoles(
	accountQuerier *account_querier.Querier,
	roleQuerier *role_querier.Querier,
	roleWriter *role_writer.Writer,
) Roles {
	return Roles{
		accountQuerier: accountQuerier,
		roleQuerier:    roleQuerier,
		roleWriter:     roleWriter,
	}
}

func (h *Roles) RoleCreate(ctx context.Context, request openapi.RoleCreateRequestObject) (openapi.RoleCreateResponseObject, error) {
	perms, err := deserialisePermissionList(request.Body.Permissions)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	role, err := h.roleWriter.Create(ctx, request.Body.Name, request.Body.Colour, perms)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RoleCreate200JSONResponse{
		RoleCreateOKJSONResponse: openapi.RoleCreateOKJSONResponse(serialiseRolePtr(role)),
	}, nil
}

func (h *Roles) RoleList(ctx context.Context, request openapi.RoleListRequestObject) (openapi.RoleListResponseObject, error) {
	roles, err := h.roleQuerier.List(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RoleList200JSONResponse{
		RoleListOKJSONResponse: openapi.RoleListOKJSONResponse{
			Roles: serialiseRoleList(roles),
		},
	}, nil
}

func (h *Roles) RoleGet(ctx context.Context, request openapi.RoleGetRequestObject) (openapi.RoleGetResponseObject, error) {
	id := role.RoleID(openapi.ParseID(request.RoleId))

	role, err := h.roleQuerier.Get(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RoleGet200JSONResponse{
		RoleGetOKJSONResponse: openapi.RoleGetOKJSONResponse(serialiseRolePtr(role)),
	}, nil
}

func (h *Roles) RoleUpdate(ctx context.Context, request openapi.RoleUpdateRequestObject) (openapi.RoleUpdateResponseObject, error) {
	id := role.RoleID(openapi.ParseID(request.RoleId))

	opts := []role_writer.Mutation{}

	if request.Body.Name != nil {
		opts = append(opts, role_writer.WithName(*request.Body.Name))
	}

	if request.Body.Colour != nil {
		opts = append(opts, role_writer.WithColour(*request.Body.Colour))
	}

	if request.Body.Permissions != nil {
		perms, err := deserialisePermissionList(*request.Body.Permissions)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		opts = append(opts, role_writer.WithPermissions(perms))
	}

	role, err := h.roleWriter.Update(ctx, id, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RoleUpdate200JSONResponse{
		RoleGetOKJSONResponse: openapi.RoleGetOKJSONResponse(serialiseRolePtr(role)),
	}, nil
}

func (h *Roles) RoleDelete(ctx context.Context, request openapi.RoleDeleteRequestObject) (openapi.RoleDeleteResponseObject, error) {
	err := h.roleWriter.Delete(ctx, role.RoleID(openapi.ParseID(request.RoleId)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return nil, nil
}

func serialiseRole(in role.Role) openapi.Role {
	return openapi.Role{
		Id:          in.ID.String(),
		Name:        in.Name,
		Colour:      in.Colour,
		Permissions: serialisePermissionList(in.Permissions),
	}
}

func serialiseRolePtr(in *role.Role) openapi.Role {
	if in == nil {
		return openapi.Role{}
	}
	return serialiseRole(*in)
}

func serialiseRoleList(in role.Roles) openapi.RoleList {
	return dt.Map(in, serialiseRolePtr)
}

func serialiseHeldRole(in held.Role) openapi.AccountRole {
	return openapi.AccountRole{
		Id:          in.ID.String(),
		Name:        in.Name,
		Colour:      in.Colour,
		Permissions: serialisePermissionList(in.Permissions),
		Badge:       in.Badge,
		Default:     in.Default,
	}
}

func serialiseHeldRolePtr(in *held.Role) openapi.AccountRole {
	if in == nil {
		return openapi.AccountRole{}
	}
	return serialiseHeldRole(*in)
}

func serialiseHeldRoleList(in held.Roles) openapi.AccountRoleList {
	return dt.Map(in, serialiseHeldRolePtr)
}

func serialisePermission(in rbac.Permission) openapi.Permission {
	return openapi.Permission(in.String())
}

func serialisePermissionList(in rbac.Permissions) openapi.PermissionList {
	return dt.Map(in.List(), serialisePermission)
}

func deserialisePermission(p openapi.Permission) (rbac.Permission, error) {
	return rbac.NewPermission(string(p))
}

func deserialisePermissionList(in openapi.PermissionList) (rbac.PermissionList, error) {
	ps, err := dt.MapErr(in, deserialisePermission)
	if err != nil {
		return nil, err
	}

	return ps, nil
}
