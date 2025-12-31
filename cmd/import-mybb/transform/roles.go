package transform

import (
	"context"
	"fmt"
	"log"

	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/cmd/import-mybb/loader"
	"github.com/Southclaws/storyden/cmd/import-mybb/writer"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/rs/xid"
)

func ImportRoles(ctx context.Context, w *writer.Writer, data *loader.MyBBData) error {
	if len(data.UserGroups) == 0 {
		log.Println("No user groups to import")
		return nil
	}

	builders := make([]*ent.RoleCreate, 0, len(data.UserGroups))

	for _, group := range data.UserGroups {
		id := xid.New()

		permissions := mapMyBBPermissions(group)

		builder := w.Client().Role.Create().
			SetID(id).
			SetName(group.Title).
			SetColour("").
			SetPermissions(permissions).
			SetSortKey(float64(group.GID))

		builders = append(builders, builder)
		w.RoleIDMap[group.GID] = id
	}

	roles, err := w.CreateRoles(ctx, builders)
	if err != nil {
		return fmt.Errorf("create roles: %w", err)
	}

	log.Printf("Imported %d roles", len(roles))
	return nil
}

func mapMyBBPermissions(group loader.MyBBUserGroup) []string {
	var perms []string

	// Administrator - if cancp is set, grant full admin access
	if group.CanCP == 1 {
		return []string{rbac.PermissionAdministrator.String()}
	}

	// Read permissions
	if group.CanViewThreads == 1 {
		perms = append(perms, rbac.PermissionReadPublishedThreads.String())
	}
	if group.CanViewProfiles == 1 {
		perms = append(perms, rbac.PermissionListProfiles.String())
		perms = append(perms, rbac.PermissionReadProfile.String())
	}

	// Write permissions
	if group.CanPostThreads == 1 || group.CanPostReplys == 1 {
		perms = append(perms, rbac.PermissionCreatePost.String())
	}
	if group.CanRateThreads == 1 {
		perms = append(perms, rbac.PermissionCreateReaction.String())
	}
	if group.CanUploadAvatars == 1 {
		perms = append(perms, rbac.PermissionUploadAsset.String())
	}

	// Moderation permissions
	if group.CanEditPosts == 1 || group.CanDeletePosts == 1 || group.CanDeleteThreads == 1 || group.IsSuperMod == 1 {
		perms = append(perms, rbac.PermissionManagePosts.String())
	}
	if group.CanManageAnnounce == 1 || group.CanManageModQueue == 1 {
		perms = append(perms, rbac.PermissionManageCategories.String())
	}
	if group.CanBanUsers == 1 {
		perms = append(perms, rbac.PermissionManageSuspensions.String())
	}

	// Library permissions (grant to all who can view)
	if group.CanViewThreads == 1 {
		perms = append(perms, rbac.PermissionReadPublishedLibrary.String())
	}

	// Collection permissions
	perms = append(perms, rbac.PermissionListCollections.String())
	perms = append(perms, rbac.PermissionReadCollection.String())

	return perms
}
