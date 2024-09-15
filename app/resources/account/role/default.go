package role

import (
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/rs/xid"
)

var DefaultRoleEveryone = Role{
	ID:     RoleID(utils.Must(xid.FromString("00000000000000000010"))),
	Name:   "Everyone",
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
}

var DefaultRoleAdmin = Role{
	ID:          RoleID(utils.Must(xid.FromString("00000000000000000020"))),
	Name:        "Admin",
	Colour:      "red",
	Permissions: rbac.NewList(rbac.PermissionAdministrator),
}
