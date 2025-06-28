package access_key_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rs/xid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestAccessKeyCRUD(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)
			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			memberSession := sh.WithSession(memberCtx)
			randomCtx, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			randomSession := sh.WithSession(randomCtx)

			// Grand personal access key usage to member
			grant(t, cl, adminSession, member.Handle, openapi.PermissionList{openapi.USEPERSONALACCESSKEYS})

			t.Run("unauthenticated", func(t *testing.T) {
				tests.AssertRequest(
					cl.AccessKeyCreateWithResponse(root, openapi.AccessKeyInitialProps{
						Name: "test",
					}),
				)(t, http.StatusForbidden)
			})

			t.Run("no_permissions", func(t *testing.T) {
				tests.AssertRequest(
					cl.AccessKeyCreateWithResponse(root, openapi.AccessKeyInitialProps{
						Name: "test",
					}, randomSession),
				)(t, http.StatusForbidden)
			})

			t.Run("create_list_delete", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				ak := tests.AssertRequest(
					cl.AccessKeyCreateWithResponse(root, openapi.AccessKeyInitialProps{
						Name: "test",
					}, memberSession),
				)(t, http.StatusOK)
				a.Equal("test", ak.JSON200.Name)

				list1 := tests.AssertRequest(
					cl.AccessKeyListWithResponse(root, memberSession),
				)(t, http.StatusOK)
				r.Len(list1.JSON200.Keys, 1)
				ak1 := list1.JSON200.Keys[0]
				a.Equal(ak.JSON200.Id, ak1.Id)
				a.True(ak1.Enabled)

				tests.AssertRequest(
					cl.AccessKeyDeleteWithResponse(root, ak.JSON200.Id, memberSession),
				)(t, http.StatusNoContent)

				list2 := tests.AssertRequest(
					cl.AccessKeyListWithResponse(root, memberSession),
				)(t, http.StatusOK)
				r.Len(list2.JSON200.Keys, 1)
				ak1 = list2.JSON200.Keys[0]
				a.Equal(ak.JSON200.Id, ak1.Id)
				a.False(ak1.Enabled)
			})

			t.Run("admin_list_permission_denied", func(t *testing.T) {
				tests.AssertRequest(
					cl.AccessKeyCreateWithResponse(root, openapi.AccessKeyInitialProps{
						Name: "test",
					}, memberSession),
				)(t, http.StatusOK)

				tests.AssertRequest(
					cl.AdminAccessKeyListWithResponse(root, memberSession),
				)(t, http.StatusForbidden)
			})

			t.Run("admin_list", func(t *testing.T) {
				a := assert.New(t)

				// member creates a key
				ak := tests.AssertRequest(
					cl.AccessKeyCreateWithResponse(root, openapi.AccessKeyInitialProps{
						Name: "test",
					}, memberSession),
				)(t, http.StatusOK)
				a.Equal("test", ak.JSON200.Name)

				// admin can see it
				list1 := tests.AssertRequest(
					cl.AdminAccessKeyListWithResponse(root, adminSession),
				)(t, http.StatusOK)
				ak1 := findOwnedAccessKey(t, list1.JSON200.Keys, ak.JSON200.Id)
				a.Equal(ak.JSON200.Id, ak1.Id)
				a.Equal(ak.JSON200.Name, ak1.Name)
				a.Equal(member.Handle, ak1.Handle)
			})

			t.Run("admin_revoke", func(t *testing.T) {
				a := assert.New(t)

				// member creates a key
				ak := tests.AssertRequest(
					cl.AccessKeyCreateWithResponse(root, openapi.AccessKeyInitialProps{
						Name: "test",
					}, memberSession),
				)(t, http.StatusOK)
				a.Equal("test", ak.JSON200.Name)

				// revoke the key as admin
				tests.AssertRequest(
					cl.AdminAccessKeyDeleteWithResponse(root, ak.JSON200.Id, adminSession),
				)(t, http.StatusNoContent)

				// list as the member, key is revoked
				list := tests.AssertRequest(
					cl.AccessKeyListWithResponse(root, memberSession),
				)(t, http.StatusOK)
				ak1 := findAccessKey(t, list.JSON200.Keys, ak.JSON200.Id)
				a.Equal(ak.JSON200.Id, ak1.Id)
				a.False(ak1.Enabled)
			})
		}))
	}))
}

func grant(
	t *testing.T,
	cl *openapi.ClientWithResponses,
	adminSession openapi.RequestEditorFn,
	targetHandle openapi.AccountHandle,
	permissions openapi.PermissionList,
) {
	t.Helper()

	name := "role-with-access-keys-" + xid.New().String()
	colour := "red"

	role := tests.AssertRequest(
		cl.RoleCreateWithResponse(
			t.Context(),
			openapi.RoleCreateJSONRequestBody{
				Name:        name,
				Colour:      colour,
				Permissions: permissions,
			}, adminSession),
	)(t, http.StatusOK)

	tests.AssertRequest(
		cl.AccountAddRoleWithResponse(
			t.Context(),
			targetHandle,
			role.JSON200.Id,
			adminSession,
		),
	)(t, http.StatusOK)
}

func findAccessKey(
	t *testing.T,
	list openapi.AccessKeyList,
	akid openapi.Identifier,
) openapi.AccessKey {
	k, found := lo.Find(list, func(k openapi.AccessKey) bool {
		return k.Id == akid
	})
	require.True(t, found)
	return k
}

func findOwnedAccessKey(
	t *testing.T,
	list openapi.OwnedAccessKeyList,
	akid openapi.Identifier,
) openapi.OwnedAccessKey {
	k, found := lo.Find(list, func(k openapi.OwnedAccessKey) bool {
		return k.Id == akid
	})
	require.True(t, found)
	return k
}
