package tests

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

func TestGuestLibraryAccess(t *testing.T) {
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

			// Create test content
			published := openapi.Published
			content := "<body>Test library node content</body>"

			// Set explicit permissions for Guest - INCLUDE library access
			AssertRequest(cl.RoleUpdateWithResponse(adminCtx,
				role.DefaultRoleGuestID.String(),
				openapi.RoleUpdateJSONRequestBody{
					Permissions: &openapi.PermissionList{
						"READ_PUBLISHED_THREADS",
						"READ_PUBLISHED_LIBRARY", // <-- This should allow guest access
					},
				}, adminSession))(t, http.StatusOK)

			// Create a test node
			nodeResp := AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
				Name:       "guest_access_test_" + uuid.NewString(),
				Content:    &content,
				Visibility: &published,
			}, adminSession))(t, http.StatusOK)

			// Test guest access to NodeList (no authentication)
			listResp, listErr := cl.NodeListWithResponse(root, &openapi.NodeListParams{})
			if listErr != nil {
				t.Logf("NodeList error: %v", listErr)
			} else if listResp.StatusCode() == http.StatusUnauthorized {
				t.Errorf("ERROR: NodeList got 401 Unauthorized instead of 200 OK")
				t.Logf("HTTP response: %+v", listResp.HTTPResponse)
			} else {
				t.Logf("NodeList response status: %d (expected 200)", listResp.StatusCode())
			}

			// Test guest access to NodeGet (no authentication)
			getResp, getErr := cl.NodeGetWithResponse(root, nodeResp.JSON200.Slug, &openapi.NodeGetParams{})
			if getErr != nil {
				t.Logf("NodeGet error: %v", getErr)
			} else if getResp.StatusCode() == http.StatusUnauthorized {
				t.Errorf("ERROR: NodeGet got 401 Unauthorized instead of 200 OK")
				t.Logf("HTTP response: %+v", getResp.HTTPResponse)
			} else {
				t.Logf("NodeGet response status: %d (expected 200)", getResp.StatusCode())
			}
		}))
	}))
}