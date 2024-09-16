package rbac

import (
	"context"
	"fmt"
	"strings"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
)

var ErrPermissions = fault.New("invalid permissions", ftag.With(ftag.PermissionDenied))

//go:generate go run github.com/Southclaws/enumerator

// Actions defined, roughly are: grouped by resource, ordered by CRUD operation.

type permissionEnum string

// -
// NOTE:
// When adding new permissions, ensure that the permission is also included in
// the OpenAPI specification under the `Permission` schema enum values.
// -
const (
	// Posts (Threads, replies and posts)
	permissionCreatePost           permissionEnum = "CREATE_POST"
	permissionReadPublishedThreads permissionEnum = "READ_PUBLISHED_THREADS"
	permissionCreateReaction       permissionEnum = "CREATE_REACTION"
	permissionManagePosts          permissionEnum = "MANAGE_POSTS"
	permissionManageCategories     permissionEnum = "MANAGE_CATEGORIES"

	// Library (Page tree nodes)
	permissionReadPublishedLibrary permissionEnum = "READ_PUBLISHED_LIBRARY"
	permissionManageLibrary        permissionEnum = "MANAGE_LIBRARY"
	permissionSubmitLibraryNode    permissionEnum = "SUBMIT_LIBRARY_NODE"

	// Assets
	permissionUploadAsset permissionEnum = "UPLOAD_ASSET"

	// Profiles (listing and viewing member profiles)
	permissionListProfiles permissionEnum = "LIST_PROFILES"
	permissionReadProfile  permissionEnum = "READ_PROFILE"

	// Collections
	permissionCreateCollection  permissionEnum = "CREATE_COLLECTION"
	permissionListCollections   permissionEnum = "LIST_COLLECTIONS"
	permissionReadCollection    permissionEnum = "READ_COLLECTION"
	permissionManageCollections permissionEnum = "MANAGE_COLLECTIONS"
	permissionCollectionSubmit  permissionEnum = "COLLECTION_SUBMIT"

	// Administrative (Settings, bans, etc)
	// NOTE: Currently all-or-nothing, it will eventually be much more granular.
	permissionManageSettings    permissionEnum = "MANAGE_SETTINGS"
	permissionManageSuspensions permissionEnum = "MANAGE_SUSPENSIONS"
	permissionManageRoles       permissionEnum = "MANAGE_ROLES"

	// Administrator implicitly has all permissions.
	permissionAdministrator permissionEnum = "ADMINISTRATOR"
)

type PermissionList []Permission

func (p PermissionList) String() string {
	s := make([]string, len(p))
	for i, pp := range p {
		s[i] = pp.String()
	}
	return strings.Join(s, ",")
}

type Permissions struct {
	p PermissionList
	m map[Permission]struct{}
}

func NewPermissions(s []string) (*Permissions, error) {
	ps, err := dt.MapErr(s, NewPermission)
	if err != nil {
		return nil, err
	}

	list := NewList(ps...)

	return &list, nil
}

func (p Permissions) List() PermissionList {
	return p.p
}

func NewList(permissions ...Permission) Permissions {
	m := map[Permission]struct{}{}
	for _, p := range permissions {
		m[p] = struct{}{}
	}

	return Permissions{permissions, m}
}

func (p Permissions) HasAll(perms ...Permission) bool {
	for _, pp := range perms {
		if _, ok := p.m[pp]; !ok {
			return false
		}
	}
	return true
}

func (p Permissions) HasAny(perms ...Permission) bool {
	for _, pp := range perms {
		if _, ok := p.m[pp]; ok {
			return true
		}
	}
	return false
}

// Authorise will check if the account holds any of the permissions provided. If
// it does not, it then runs the additional check function in order to apply any
// domain-specific logic such as resource ownership, to determine authorisation.
func (p Permissions) Authorise(ctx context.Context, fn func() error, perms ...Permission) error {
	ctx = fctx.WithMeta(ctx,
		"permissions", PermissionList(perms).String(),
	)

	// PermissionAdministrator can do anything.
	if p.HasAny(PermissionAdministrator) {
		return nil
	}

	// If permissions are valid, do not need to run additional check function.
	if p.HasAny(perms...) {
		return nil
	}

	// Additional check functions are for when permission requires logic, such
	// as checking for resource ownership or non-publishing draft submissions.
	if fn != nil {
		if err := fn(); err != nil {
			return fault.Wrap(
				fmt.Errorf("%w: additional check failed: %w", ErrPermissions, err),
				fctx.With(ctx),
				ftag.With(ftag.PermissionDenied),
			)
		}

		// If the additional check function passes, override the missing perm.
		return nil
	}

	// No additional check and permission is missing.
	return fault.Wrap(ErrPermissions, fctx.With(ctx))
}
