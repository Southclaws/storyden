package account

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/rs/xid"
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

func TestAccountCreateProvisioning(t *testing.T) {
	if tests.IsSharedPostgresDatabase() {
		t.Skip("skipping account provisioning test on shared postgres database")
	}

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

			accessKey := tests.AssertRequest(
				cl.AccessKeyCreateWithResponse(root, openapi.AccessKeyInitialProps{
					Name: "account-provisioning-" + xid.New().String(),
				}, adminSession),
			)(t, http.StatusOK)
			accessKeySession := createAccessKeyAuth(accessKey.JSON200.Secret)

			unauthenticated, err := cl.AccountManageCreateWithResponse(root, openapi.AccountManageCreateJSONRequestBody{
				Handle: openapi.AccountHandle(xid.New().String()),
			})
			tests.Status(t, err, unauthenticated, http.StatusUnauthorized)

			forbidden, err := cl.AccountManageCreateWithResponse(root, openapi.AccountManageCreateJSONRequestBody{
				Handle: openapi.AccountHandle(xid.New().String()),
			}, memberSession)
			tests.Status(t, err, forbidden, http.StatusForbidden)

			grant(t, cl, adminSession, member.Handle, openapi.PermissionList{openapi.MANAGEACCOUNTS})

			registrationMode := openapi.RegistrationModeDisabled
			tests.AssertRequest(
				cl.AdminSettingsUpdateWithResponse(root, openapi.AdminSettingsUpdateJSONRequestBody{
					RegistrationMode: &registrationMode,
				}, adminSession),
			)(t, http.StatusOK)

			email := openapi.EmailAddress(fmt.Sprintf("%s@example.com", xid.New().String()))
			name := openapi.AccountName("Provisioned Account")
			bio := openapi.AccountBio("<p>Provisioned bio</p>")
			signature := openapi.AccountSignature("<p>Provisioned signature</p>")
			meta := openapi.Metadata{"source": "provisioning-test"}
			links := openapi.ProfileExternalLinkList{
				{
					Text: "Website",
					Url:  "https://example.com",
				},
			}
			handle := openapi.AccountHandle(xid.New().String())
			verifiedStatus := openapi.AccountVerifiedStatusManual

			created := tests.AssertRequest(
				cl.AccountManageCreateWithResponse(root, openapi.AccountManageCreateJSONRequestBody{
					Handle:         handle,
					Name:           &name,
					Bio:            &bio,
					Signature:      &signature,
					Meta:           &meta,
					Links:          &links,
					EmailAddress:   &email,
					VerifiedStatus: &verifiedStatus,
				}, accessKeySession),
			)(t, http.StatusOK)
			require.NotNil(t, created.JSON200)
			assert.Equal(t, name, created.JSON200.Name)
			assert.Equal(t, "<body>"+string(bio)+"</body>", tests.StripBlockIDs(created.JSON200.Bio))
			require.NotNil(t, created.JSON200.Signature)
			assert.Equal(t, "<body>"+string(signature)+"</body>", tests.StripBlockIDs(*created.JSON200.Signature))
			assert.Equal(t, meta, created.JSON200.Meta)
			assert.Equal(t, links, created.JSON200.Links)
			assert.Equal(t, openapi.AccountVerifiedStatusManual, created.JSON200.VerifiedStatus)
			require.Len(t, created.JSON200.EmailAddresses, 1)
			assert.Equal(t, email, created.JSON200.EmailAddresses[0].EmailAddress)
			assert.False(t, created.JSON200.EmailAddresses[0].Verified)

			admin := false
			nonAdminCreated := tests.AssertRequest(
				cl.AccountManageCreateWithResponse(root, openapi.AccountManageCreateJSONRequestBody{
					Handle: openapi.AccountHandle(xid.New().String()),
					Admin:  &admin,
				}, memberSession),
			)(t, http.StatusOK)
			require.NotNil(t, nonAdminCreated.JSON200)
			assert.False(t, nonAdminCreated.JSON200.Admin)

			duplicateHandle, err := cl.AccountManageCreateWithResponse(root, openapi.AccountManageCreateJSONRequestBody{
				Handle: handle,
			}, adminSession)
			tests.Status(t, err, duplicateHandle, http.StatusConflict)

			duplicateEmailHandle := openapi.AccountHandle(xid.New().String())
			duplicateEmail, err := cl.AccountManageCreateWithResponse(root, openapi.AccountManageCreateJSONRequestBody{
				Handle:       duplicateEmailHandle,
				EmailAddress: &email,
			}, adminSession)
			tests.Status(t, err, duplicateEmail, http.StatusConflict)
		}))
	}))
}

func createAccessKeyAuth(accessKeyToken string) openapi.RequestEditorFn {
	authHeader := fmt.Sprintf("Bearer %s", accessKeyToken)

	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", authHeader)
		return nil
	}
}

func grant(
	t *testing.T,
	cl *openapi.ClientWithResponses,
	adminSession openapi.RequestEditorFn,
	targetHandle openapi.AccountHandle,
	permissions openapi.PermissionList,
) {
	t.Helper()

	name := "role-account-provisioning-" + xid.New().String()
	colour := "blue"

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
			adminSession),
	)(t, http.StatusOK)
}
