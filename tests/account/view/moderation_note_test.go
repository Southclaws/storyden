package view_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestAccountModerationNotes(t *testing.T) {
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

			viewerCtx, viewer := e2e.WithAccount(root, aw, seed.Account_005_Þórr)
			viewerSession := sh.WithSession(viewerCtx)

			managerCtx, manager := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			managerSession := sh.WithSession(managerCtx)

			manageOnlyCtx, manageOnly := e2e.WithAccount(root, aw, seed.Account_007_Freyr)
			manageOnlySession := sh.WithSession(manageOnlyCtx)

			targetCtx, target := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			targetSession := sh.WithSession(targetCtx)

			outsiderCtx, _ := e2e.WithAccount(root, aw, seed.Account_006_Freyja)
			outsiderSession := sh.WithSession(outsiderCtx)

			grantNotePermission(t, cl, adminSession, viewer.Handle, openapi.PermissionList{
				openapi.VIEWACCOUNTS,
				openapi.VIEWMODERATIONNOTES,
			})
			grantNotePermission(t, cl, adminSession, manager.Handle, openapi.PermissionList{
				openapi.VIEWACCOUNTS,
				openapi.MANAGEMODERATIONNOTES,
				openapi.VIEWMODERATIONNOTES,
			})
			grantNotePermission(t, cl, adminSession, manageOnly.Handle, openapi.PermissionList{
				openapi.VIEWACCOUNTS,
				openapi.MANAGEMODERATIONNOTES,
			})

			t.Run("manager_can_create_and_view_notes", func(t *testing.T) {
				r := require.New(t)

				create := tests.AssertRequest(cl.AccountModerationNoteCreateWithResponse(
					root,
					target.ID.String(),
					openapi.AccountModerationNoteCreateJSONRequestBody{Content: "Observed escalation pattern in support tickets."},
					managerSession,
				))(t, http.StatusOK)

				r.Equal(manager.ID.String(), create.JSON200.Author.Id)
				r.Equal("Observed escalation pattern in support tickets.", create.JSON200.Content)

				list := tests.AssertRequest(cl.AccountModerationNoteListWithResponse(
					root,
					target.ID.String(),
					managerSession,
				))(t, http.StatusOK)

				r.Len(list.JSON200.Notes, 1)
				r.Equal(create.JSON200.Id, list.JSON200.Notes[0].Id)

				tests.AssertRequest(cl.AccountModerationNoteDeleteWithResponse(
					root,
					target.ID.String(),
					string(create.JSON200.Id),
					managerSession,
				))(t, http.StatusNoContent)

				listAfterDelete := tests.AssertRequest(cl.AccountModerationNoteListWithResponse(
					root,
					target.ID.String(),
					managerSession,
				))(t, http.StatusOK)
				r.Len(listAfterDelete.JSON200.Notes, 0)
			})

			t.Run("manager_without_view_can_create_and_delete_notes", func(t *testing.T) {
				r := require.New(t)

				create := tests.AssertRequest(cl.AccountModerationNoteCreateWithResponse(
					root,
					target.ID.String(),
					openapi.AccountModerationNoteCreateJSONRequestBody{Content: "Escalated for follow-up by manage-only role."},
					manageOnlySession,
				))(t, http.StatusOK)

				r.Equal(manageOnly.ID.String(), create.JSON200.Author.Id)
				r.Equal("Escalated for follow-up by manage-only role.", create.JSON200.Content)

				tests.AssertRequest(cl.AccountModerationNoteListWithResponse(
					root,
					target.ID.String(),
					manageOnlySession,
				))(t, http.StatusForbidden)

				tests.AssertRequest(cl.AccountModerationNoteDeleteWithResponse(
					root,
					target.ID.String(),
					string(create.JSON200.Id),
					manageOnlySession,
				))(t, http.StatusNoContent)

				listAfterDelete := tests.AssertRequest(cl.AccountModerationNoteListWithResponse(
					root,
					target.ID.String(),
					managerSession,
				))(t, http.StatusOK)
				r.Len(listAfterDelete.JSON200.Notes, 0)
			})

			t.Run("viewer_can_read_but_not_create", func(t *testing.T) {
				r := require.New(t)

				create := tests.AssertRequest(cl.AccountModerationNoteCreateWithResponse(
					root,
					target.ID.String(),
					openapi.AccountModerationNoteCreateJSONRequestBody{Content: "manager-only cleanup note"},
					managerSession,
				))(t, http.StatusOK)

				list := tests.AssertRequest(cl.AccountModerationNoteListWithResponse(
					root,
					target.ID.String(),
					viewerSession,
				))(t, http.StatusOK)
				r.NotEmpty(list.JSON200.Notes)
				r.Truef(
					containsNoteWithID(list.JSON200.Notes, create.JSON200.Id),
					"expected viewer list to include created note id %s",
					create.JSON200.Id,
				)

				tests.AssertRequest(cl.AccountModerationNoteCreateWithResponse(
					root,
					target.ID.String(),
					openapi.AccountModerationNoteCreateJSONRequestBody{Content: "should fail"},
					viewerSession,
				))(t, http.StatusForbidden)

				tests.AssertRequest(cl.AccountModerationNoteDeleteWithResponse(
					root,
					target.ID.String(),
					string(create.JSON200.Id),
					viewerSession,
				))(t, http.StatusForbidden)

				tests.AssertRequest(cl.AccountModerationNoteDeleteWithResponse(
					root,
					target.ID.String(),
					string(create.JSON200.Id),
					managerSession,
				))(t, http.StatusNoContent)
			})

			t.Run("target_member_cannot_read_or_create_notes", func(t *testing.T) {
				create := tests.AssertRequest(cl.AccountModerationNoteCreateWithResponse(
					root,
					target.ID.String(),
					openapi.AccountModerationNoteCreateJSONRequestBody{Content: "staff visibility test note"},
					managerSession,
				))(t, http.StatusOK)

				tests.AssertRequest(cl.AccountModerationNoteListWithResponse(
					root,
					target.ID.String(),
					targetSession,
				))(t, http.StatusForbidden)

				tests.AssertRequest(cl.AccountModerationNoteCreateWithResponse(
					root,
					target.ID.String(),
					openapi.AccountModerationNoteCreateJSONRequestBody{Content: "self-note"},
					targetSession,
				))(t, http.StatusForbidden)

				tests.AssertRequest(cl.AccountModerationNoteDeleteWithResponse(
					root,
					target.ID.String(),
					string(create.JSON200.Id),
					targetSession,
				))(t, http.StatusForbidden)

				tests.AssertRequest(cl.AccountModerationNoteDeleteWithResponse(
					root,
					target.ID.String(),
					string(create.JSON200.Id),
					managerSession,
				))(t, http.StatusNoContent)
			})

			t.Run("non_admin_without_permissions_cannot_operate_crud_endpoints", func(t *testing.T) {
				create := tests.AssertRequest(cl.AccountModerationNoteCreateWithResponse(
					root,
					target.ID.String(),
					openapi.AccountModerationNoteCreateJSONRequestBody{Content: "permission baseline note"},
					managerSession,
				))(t, http.StatusOK)

				tests.AssertRequest(cl.AccountModerationNoteListWithResponse(
					root,
					target.ID.String(),
					outsiderSession,
				))(t, http.StatusForbidden)

				tests.AssertRequest(cl.AccountModerationNoteCreateWithResponse(
					root,
					target.ID.String(),
					openapi.AccountModerationNoteCreateJSONRequestBody{Content: "outsider should fail"},
					outsiderSession,
				))(t, http.StatusForbidden)

				tests.AssertRequest(cl.AccountModerationNoteDeleteWithResponse(
					root,
					target.ID.String(),
					string(create.JSON200.Id),
					outsiderSession,
				))(t, http.StatusForbidden)

				tests.AssertRequest(cl.AccountModerationNoteDeleteWithResponse(
					root,
					target.ID.String(),
					string(create.JSON200.Id),
					managerSession,
				))(t, http.StatusNoContent)
			})
		}))
	}))
}

func grantNotePermission(
	t *testing.T,
	cl *openapi.ClientWithResponses,
	adminSession openapi.RequestEditorFn,
	targetHandle openapi.AccountHandle,
	permissions openapi.PermissionList,
) {
	t.Helper()

	role := tests.AssertRequest(
		cl.RoleCreateWithResponse(
			t.Context(),
			openapi.RoleCreateJSONRequestBody{
				Name:        "role-notes-" + xid.New().String(),
				Colour:      "indigo",
				Permissions: permissions,
			},
			adminSession,
		),
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

func containsNoteWithID(notes []openapi.ModerationNote, id openapi.Identifier) bool {
	for _, note := range notes {
		if note.Id == id {
			return true
		}
	}

	return false
}
