package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/account/role/held"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
)

type Roles struct {
	accountQuerier *account_querier.Querier
	roleQuerier    *role_repo.Repository
	roleWriter     role_repo.Writer
}

func NewRoles(
	accountQuerier *account_querier.Querier,
	roleQuerier *role_repo.Repository,
	roleWriter role_repo.Writer,
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

	opts := []role_repo.Mutation{}
	if request.Body.Meta != nil {
		opts = append(opts, role_repo.WithMeta(map[string]any(*request.Body.Meta)))
	}

	role, err := h.roleWriter.Create(ctx, request.Body.Name, request.Body.Colour, perms, opts...)
	if err != nil {
		if ent.IsConstraintError(err) {
			err = fault.Wrap(err, fmsg.WithDesc("unique", "A role with that name already exists"), ftag.With(ftag.AlreadyExists))
		}
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

	opts := []role_repo.Mutation{}

	if request.Body.Name != nil {
		opts = append(opts, role_repo.WithName(*request.Body.Name))
	}

	if request.Body.Colour != nil {
		opts = append(opts, role_repo.WithColour(*request.Body.Colour))
	}

	if request.Body.Permissions != nil {
		perms, err := deserialisePermissionList(*request.Body.Permissions)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		opts = append(opts, role_repo.WithPermissions(perms))
	}

	if request.Body.Meta != nil {
		opts = append(opts, role_repo.WithMeta(map[string]any(*request.Body.Meta)))
	}

	role, err := h.roleWriter.Update(ctx, id, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RoleUpdate200JSONResponse{
		RoleGetOKJSONResponse: openapi.RoleGetOKJSONResponse(serialiseRolePtr(role)),
	}, nil
}

func (h *Roles) RoleUpdateOrder(ctx context.Context, request openapi.RoleUpdateOrderRequestObject) (openapi.RoleUpdateOrderResponseObject, error) {
	if request.Body == nil {
		return nil, fault.Wrap(
			fault.New("missing role reorder body", ftag.With(ftag.InvalidArgument), fmsg.With("request body is required")),
			fctx.With(ctx),
		)
	}

	ids := dt.Map(request.Body.RoleIds, func(id openapi.Identifier) role.RoleID {
		return role.RoleID(openapi.ParseID(id))
	})

	if err := h.roleWriter.UpdateSortOrder(ctx, ids); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	roles, err := h.roleQuerier.List(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.RoleUpdateOrder200JSONResponse{
		RoleListOKJSONResponse: openapi.RoleListOKJSONResponse{
			Roles: serialiseRoleList(roles),
		},
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
		Meta:        serialiseMetadata(in.Metadata),
		CreatedAt:   in.CreatedAt,
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
		Meta:        serialiseMetadata(in.Metadata),
		Badge:       in.Badge,
		Default:     in.Default,
		CreatedAt:   in.CreatedAt,
	}
}

func serialiseMetadata(meta map[string]any) *openapi.Metadata {
	if meta == nil {
		return nil
	}
	v := openapi.Metadata(meta)
	return &v
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

func serialiseHeldRoleRef(in held.Role) openapi.AccountRoleRef {
	return openapi.AccountRoleRef{
		Id:      in.ID.String(),
		Name:    in.Name,
		Colour:  in.Colour,
		Meta:    serialiseMetadata(in.Metadata),
		Badge:   in.Badge,
		Default: in.Default,
	}
}

func serialiseHeldRoleRefPtr(in *held.Role) openapi.AccountRoleRef {
	return serialiseHeldRoleRef(*in)
}

func serialiseHeldRoleRefList(in held.Roles) openapi.AccountRoleRefList {
	return dt.Map(in, serialiseHeldRoleRefPtr)
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
