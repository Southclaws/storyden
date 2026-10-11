package rbac

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPermissionsIntersect(t *testing.T) {
	t.Parallel()

	token := NewList(PermissionCreatePost, PermissionManageRoles, PermissionReadProfile)
	account := NewList(PermissionCreatePost, PermissionReadProfile)

	granted := token.Intersect(account)

	assert.True(t, granted.HasAll(PermissionCreatePost, PermissionReadProfile))
	assert.False(t, granted.HasAny(PermissionManageRoles), "a permission the account lost must not survive")
	assert.Len(t, granted.List(), 2)
}

func TestPermissionsIntersectDropsAdministrator(t *testing.T) {
	t.Parallel()

	token := NewList(PermissionAdministrator, PermissionCreatePost)
	account := NewList(PermissionCreatePost)

	granted := token.Intersect(account)

	assert.False(t, granted.HasAny(PermissionAdministrator), "Authorise short circuits on administrator, it cannot be inherited from a stale token")
	assert.True(t, granted.HasAll(PermissionCreatePost))
}

func TestPermissionsIntersectWithEmpty(t *testing.T) {
	t.Parallel()

	token := NewList(PermissionCreatePost, PermissionUploadAsset)

	assert.Empty(t, token.Intersect(NewList()).List())
	assert.Empty(t, NewList().Intersect(token).List())
}

func TestHasAnyWriteCountsManagementPermissions(t *testing.T) {
	t.Parallel()

	// these management perms must count as write so the guest role cannot be
	// granted them via validateGuestPermissionsMutation
	for _, p := range []Permission{
		PermissionManageReports,
		PermissionManageModerationNotes,
		PermissionManageRobots,
		PermissionManageTrails,
	} {
		assert.True(t, NewList(p).HasAnyWrite(), p.String())
	}

	readOnly := NewList(PermissionReadPublishedThreads, PermissionReadProfile)
	assert.False(t, readOnly.HasAnyWrite(), "a read-only set is not write")
}
