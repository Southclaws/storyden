package role

import (
	"math"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/internal/utils"
)

var (
	DefaultRoleGuestID  = RoleID(utils.Must(xid.FromString("0000000000000000000g")))
	DefaultRoleMemberID = RoleID(utils.Must(xid.FromString("000000000000000000m0")))
	DefaultRoleAdminID  = RoleID(utils.Must(xid.FromString("00000000000000000a00")))
)

var DefaultRoleMember = Role{
	ID:     DefaultRoleMemberID,
	Name:   "Member",
	Colour: "green",
	Permissions: rbac.NewList(
		rbac.PermissionCreatePost,
		rbac.PermissionReadPublishedThreads,
		rbac.PermissionCreateReaction,
		rbac.PermissionReadPublishedLibrary,
		rbac.PermissionSubmitLibraryNode,
		rbac.PermissionUploadAsset,
		rbac.PermissionListProfiles,
		rbac.PermissionReadProfile,
		rbac.PermissionCreateCollection,
		rbac.PermissionListCollections,
		rbac.PermissionReadCollection,
		rbac.PermissionCollectionSubmit,
	),
	SortKey: -1, // Always sorts after guest and before custom roles.
}

var DefaultRoleGuest = Role{
	ID:     DefaultRoleGuestID,
	Name:   "Guest",
	Colour: "gray",
	Permissions: rbac.NewList(
		rbac.PermissionReadPublishedThreads,
		rbac.PermissionReadPublishedLibrary,
		rbac.PermissionListProfiles,
		rbac.PermissionReadProfile,
		rbac.PermissionListCollections,
		rbac.PermissionReadCollection,
	),
	SortKey: -2, // Always sorts first.
}

var DefaultRoleAdmin = Role{
	ID:          DefaultRoleAdminID,
	Name:        "Admin",
	Colour:      "red",
	Permissions: rbac.NewList(rbac.PermissionAdministrator),
	SortKey:     math.MaxFloat64, // Always sorts last.
}
