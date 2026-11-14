package role_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/admin_guard"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/account/account_suspension"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestAdministratorLockout(t *testing.T) {
	cfg := administratorTestConfig(t)
	integration.Test(t, cfg, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		suspender account_suspension.Service,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)
			adminFlag := false

			t.Run("last_builtin_administrator_cannot_be_demoted_or_suspended", func(t *testing.T) {
				tests.AssertRequest(cl.AccountManageUpdateWithResponse(root, admin.ID.String(), openapi.AccountManageUpdateJSONRequestBody{
					Admin: &adminFlag,
				}, adminSession))(t, http.StatusBadRequest)
				tests.AssertRequest(cl.AdminAccountBanCreateWithResponse(root, admin.Handle, adminSession))(t, http.StatusBadRequest)

				stored := tests.AssertRequest(cl.AccountGetWithResponse(root, adminSession))(t, http.StatusOK)
				require.True(t, stored.JSON200.Admin)
				require.Nil(t, stored.JSON200.Suspended)
				tests.AssertRequest(cl.AdminThemeGetWithResponse(root, adminSession))(t, http.StatusOK)
			})

			created := tests.AssertRequest(cl.RoleCreateWithResponse(root, openapi.RoleCreateJSONRequestBody{
				Name:        "last-admin-" + xid.New().String(),
				Colour:      "red",
				Permissions: openapi.PermissionList{openapi.ADMINISTRATOR},
			}, adminSession))(t, http.StatusOK)
			roleID := created.JSON200.Id
			tests.AssertRequest(cl.AccountAddRoleWithResponse(root, admin.Handle, roleID, adminSession))(t, http.StatusOK)

			t.Run("another_role_on_the_same_account_preserves_administration", func(t *testing.T) {
				tests.AssertRequest(cl.AccountManageUpdateWithResponse(root, admin.ID.String(), openapi.AccountManageUpdateJSONRequestBody{
					Admin: &adminFlag,
				}, adminSession))(t, http.StatusOK)
				tests.AssertRequest(cl.AdminThemeGetWithResponse(root, adminSession))(t, http.StatusOK)
			})

			t.Run("last_custom_administrator_cannot_be_removed", func(t *testing.T) {
				tests.AssertRequest(cl.AccountRemoveRoleWithResponse(root, admin.Handle, roleID, adminSession))(t, http.StatusBadRequest)
				permissions := openapi.PermissionList{openapi.MANAGEROLES}
				tests.AssertRequest(cl.RoleUpdateWithResponse(root, roleID, openapi.RoleUpdateJSONRequestBody{
					Permissions: &permissions,
				}, adminSession))(t, http.StatusBadRequest)
				tests.AssertRequest(cl.RoleDeleteWithResponse(root, roleID, adminSession))(t, http.StatusBadRequest)
				_, err := suspender.Suspend(root, admin.ID)
				require.ErrorIs(t, err, admin_guard.ErrLastAdministrator)

				stored := tests.AssertRequest(cl.RoleGetWithResponse(root, roleID, adminSession))(t, http.StatusOK)
				require.Contains(t, stored.JSON200.Permissions, openapi.ADMINISTRATOR)
				tests.AssertRequest(cl.AdminThemeGetWithResponse(root, adminSession))(t, http.StatusOK)
			})

			otherCtx, other := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			otherSession := sh.WithSession(otherCtx)
			tests.AssertRequest(cl.AccountAddRoleWithResponse(root, other.Handle, roleID, adminSession))(t, http.StatusOK)

			t.Run("shared_role_cannot_demote_all_administrators", func(t *testing.T) {
				tests.AssertRequest(cl.RoleDeleteWithResponse(root, roleID, adminSession))(t, http.StatusBadRequest)
				tests.AssertRequest(cl.AdminThemeGetWithResponse(root, adminSession))(t, http.StatusOK)
				tests.AssertRequest(cl.AdminThemeGetWithResponse(root, otherSession))(t, http.StatusOK)
			})

			t.Run("independent_administrator_allows_revocation", func(t *testing.T) {
				adminFlag = true
				tests.AssertRequest(cl.AccountManageUpdateWithResponse(root, other.ID.String(), openapi.AccountManageUpdateJSONRequestBody{
					Admin: &adminFlag,
				}, adminSession))(t, http.StatusOK)
				tests.AssertRequest(cl.AccountRemoveRoleWithResponse(root, admin.Handle, roleID, adminSession))(t, http.StatusOK)
				tests.AssertRequest(cl.RoleDeleteWithResponse(root, roleID, otherSession))(t, http.StatusOK)
				tests.AssertRequest(cl.AdminThemeGetWithResponse(root, adminSession))(t, http.StatusForbidden)
				tests.AssertRequest(cl.AdminThemeGetWithResponse(root, otherSession))(t, http.StatusOK)
			})
		}))
	}))
}

func TestConcurrentAdministratorRevocation(t *testing.T) {
	cfg := administratorTestConfig(t)
	integration.Test(t, cfg, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			firstCtx, first := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			secondCtx, second := e2e.WithAccount(root, aw, seed.Account_002_Frigg)
			firstSession := sh.WithSession(firstCtx)
			secondSession := sh.WithSession(secondCtx)
			adminFlag := false
			start := make(chan struct{})
			statuses := make([]int, 2)
			errors := make([]error, 2)
			var wg sync.WaitGroup

			for index, target := range []struct {
				id      string
				session openapi.RequestEditorFn
			}{
				{first.ID.String(), firstSession},
				{second.ID.String(), secondSession},
			} {
				wg.Go(func() {
					<-start
					response, err := cl.AccountManageUpdateWithResponse(root, target.id, openapi.AccountManageUpdateJSONRequestBody{
						Admin: &adminFlag,
					}, target.session)
					errors[index] = err
					if response != nil {
						statuses[index] = response.StatusCode()
					}
				})
			}

			close(start)
			wg.Wait()
			for _, err := range errors {
				require.NoError(t, err)
			}
			require.ElementsMatch(t, []int{http.StatusOK, http.StatusBadRequest}, statuses)

			firstStored := tests.AssertRequest(cl.AccountGetWithResponse(root, firstSession))(t, http.StatusOK)
			secondStored := tests.AssertRequest(cl.AccountGetWithResponse(root, secondSession))(t, http.StatusOK)
			require.NotEqual(t, firstStored.JSON200.Admin, secondStored.JSON200.Admin)
		}))
	}))
}

func TestDefaultMemberAdministratorLockout(t *testing.T) {
	integration.Test(t, administratorTestConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		suspender account_suspension.Service,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)
			memberID := role.DefaultRoleMemberID.String()
			permissions := openapi.PermissionList{openapi.ADMINISTRATOR}
			tests.AssertRequest(cl.RoleUpdateWithResponse(root, memberID, openapi.RoleUpdateJSONRequestBody{
				Permissions: &permissions,
			}, adminSession))(t, http.StatusOK)

			adminFlag := false
			tests.AssertRequest(cl.AccountManageUpdateWithResponse(root, admin.ID.String(), openapi.AccountManageUpdateJSONRequestBody{
				Admin: &adminFlag,
			}, adminSession))(t, http.StatusOK)

			_, suspended := e2e.WithAccount(root, aw, seed.Account_002_Frigg)
			_, err := suspender.Suspend(root, suspended.ID)
			require.NoError(t, err)

			permissions = openapi.PermissionList{}
			tests.AssertRequest(cl.RoleUpdateWithResponse(root, memberID, openapi.RoleUpdateJSONRequestBody{
				Permissions: &permissions,
			}, adminSession))(t, http.StatusBadRequest)
			tests.AssertRequest(cl.RoleDeleteWithResponse(root, memberID, adminSession))(t, http.StatusBadRequest)
			tests.AssertRequest(cl.AdminThemeGetWithResponse(root, adminSession))(t, http.StatusOK)

			adminFlag = true
			tests.AssertRequest(cl.AccountManageUpdateWithResponse(root, admin.ID.String(), openapi.AccountManageUpdateJSONRequestBody{
				Admin: &adminFlag,
			}, adminSession))(t, http.StatusOK)
			tests.AssertRequest(cl.RoleUpdateWithResponse(root, memberID, openapi.RoleUpdateJSONRequestBody{
				Permissions: &permissions,
			}, adminSession))(t, http.StatusOK)
			tests.AssertRequest(cl.AdminThemeGetWithResponse(root, adminSession))(t, http.StatusOK)
		}))
	}))
}

func administratorTestConfig(t *testing.T) *config.Config {
	t.Helper()

	databaseURL, err := url.Parse(os.Getenv("DATABASE_URL"))
	require.NoError(t, err)

	if databaseURL.Scheme == "postgres" || databaseURL.Scheme == "postgresql" {
		db, err := sql.Open("pgx", databaseURL.String())
		require.NoError(t, err)
		t.Cleanup(func() { require.NoError(t, db.Close()) })

		schema := "admin_lockout_" + xid.New().String()
		_, err = db.ExecContext(t.Context(), "CREATE SCHEMA "+schema)
		require.NoError(t, err)
		t.Cleanup(func() {
			_, err := db.ExecContext(context.Background(), "DROP SCHEMA "+schema+" CASCADE")
			require.NoError(t, err)
		})

		query := databaseURL.Query()
		query.Set("search_path", schema)
		databaseURL.RawQuery = query.Encode()

		return &config.Config{DatabaseURL: databaseURL.String()}
	}

	scheme := "sqlite"
	if databaseURL.Scheme == "libsql" {
		scheme = "libsql"
	}

	databasePath := filepath.Join(t.TempDir(), "lockout.db")
	isolatedURL := scheme + "://" + databasePath + "?_pragma=foreign_keys(1)&_pragma=busy_timeout(10000)"

	return &config.Config{DatabaseURL: isolatedURL}
}
