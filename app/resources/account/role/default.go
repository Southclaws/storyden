package role

import (
	"math"

	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/rs/xid"
)

var (
	DefaultRoleEveryoneID = RoleID(utils.Must(xid.FromString("00000000000000000010")))
	DefaultRoleAdminID    = RoleID(utils.Must(xid.FromString("00000000000000000020")))
)

var DefaultRoleEveryone = Role{
	ID:     DefaultRoleEveryoneID,
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
	SortKey: 0,
}

var DefaultRoleAdmin = Role{
	ID:          DefaultRoleAdminID,
	Name:        "Admin",
	Colour:      "red",
	Permissions: rbac.NewList(rbac.PermissionAdministrator),
	SortKey:     math.MaxFloat64,
}
