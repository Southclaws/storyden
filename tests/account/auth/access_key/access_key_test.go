package access_key_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

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

func TestAccessKeyAuth(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			// An admin account, who can list or revoke any key
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			// A member account, who is granted use of "personal access keys"
			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			memberSession := sh.WithSession(memberCtx)

			// Grant the member account permission to use personal access keys
			grant(t, cl, adminSession, member.Handle, openapi.PermissionList{openapi.USEPERSONALACCESSKEYS})

			t.Run("revoked_key_cannot_authenticate", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				// Create an access key
				ak := tests.AssertRequest(
					cl.AccessKeyCreateWithResponse(root, openapi.AccessKeyInitialProps{
						Name: "test-revoked",
					}, memberSession),
				)(t, http.StatusOK)

				keySession := createAccessKeyAuth(ak.JSON200.Secret)

				// Verify the key works initially
				tests.AssertRequest(
					cl.AccountGetWithResponse(root, keySession),
				)(t, http.StatusOK)

				// Revoke the key
				tests.AssertRequest(
					cl.AccessKeyDeleteWithResponse(root, ak.JSON200.Id, memberSession),
				)(t, http.StatusNoContent)

				// Verify the key is disabled
				list := tests.AssertRequest(
					cl.AccessKeyListWithResponse(root, memberSession),
				)(t, http.StatusOK)
				r.Len(list.JSON200.Keys, 1)
				ak1 := list.JSON200.Keys[0]
				a.Equal(ak.JSON200.Id, ak1.Id)
				a.False(ak1.Enabled) // Key should be disabled

				// Try to use the revoked key for authentication - should fail
				// First test with no authentication to ensure it fails
				tests.AssertRequest(
					cl.AccountGetWithResponse(root),
				)(t, http.StatusForbidden)

				// Now test with the revoked key - should also fail
				tests.AssertRequest(
					cl.AccountGetWithResponse(root, keySession),
				)(t, http.StatusForbidden)
			})

			t.Run("access_key_authentication", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				// Create an access key
				ak := tests.AssertRequest(
					cl.AccessKeyCreateWithResponse(root, openapi.AccessKeyInitialProps{
						Name: "test-valid",
					}, memberSession),
				)(t, http.StatusOK)

				// Use the access key to authenticate and make an API call
				self := tests.AssertRequest(
					cl.AccountGetWithResponse(root, createAccessKeyAuth(ak.JSON200.Secret)),
				)(t, http.StatusOK)

				r.NotNil(self.JSON200)
				a.Equal(member.Handle, self.JSON200.Handle)
			})

			t.Run("invalid_token_formats", func(t *testing.T) {
				// Test various malformed token formats
				invalidTokens := []string{
					"invalid",   // completely invalid
					"sdpak_",    // missing id and secret
					"sdpak_123", // too short
					"wrong_123456789012abcdefghijklmnopqrstuvwxyz123456", // wrong prefix
					"sdpak_123456789012!@#$%^&*()abcdefghijklmnopqr",     // invalid characters
					"", // empty
				}

				for _, token := range invalidTokens {
					tests.AssertRequest(
						cl.AccountGetWithResponse(root, createAccessKeyAuth(token)),
					)(t, http.StatusForbidden)
				}
			})

			t.Run("browser_only_endpoint_protection", func(t *testing.T) {
				// Create an access key
				ak := tests.AssertRequest(
					cl.AccessKeyCreateWithResponse(root, openapi.AccessKeyInitialProps{
						Name: "test-browser-protection",
					}, memberSession),
				)(t, http.StatusOK)

				// Try to access browser-only endpoints that should reject access key auth
				// These endpoints typically require session cookies for security reasons
				tests.AssertRequest(
					cl.AuthPasswordCreateWithResponse(root, openapi.AuthPasswordInitialProps{
						Password: "test-password",
					}, createAccessKeyAuth(ak.JSON200.Secret)),
				)(t, http.StatusForbidden)
			})

			t.Run("access_keys_cannot_create_more_access_keys", func(t *testing.T) {
				// Create an access key
				ak := tests.AssertRequest(
					cl.AccessKeyCreateWithResponse(root, openapi.AccessKeyInitialProps{
						Name: "test-browser-protection",
					}, memberSession),
				)(t, http.StatusOK)

				// Try to use that access key to create another access key.
				tests.AssertRequest(
					cl.AccessKeyCreateWithResponse(root, openapi.AccessKeyInitialProps{
						Name: "my child...",
					}, createAccessKeyAuth(ak.JSON200.Secret)),
				)(t, http.StatusForbidden)
			})

			t.Run("access_key_expiry", func(t *testing.T) {
				// Note: This test would require creating an expired key or mocking time
				// For now, we'll test that a valid key works and note the limitation
				// TODO: Implement proper expiry testing with time mocking

				now := time.Now().Add(-24 * time.Hour) // 24 hours ago, simulating an expired key

				// Create an access key with no expiry (should work indefinitely)
				ak := tests.AssertRequest(
					cl.AccessKeyCreateWithResponse(root, openapi.AccessKeyInitialProps{
						Name:      "test-no-expiry",
						ExpiresAt: &now,
					}, memberSession),
				)(t, http.StatusOK)

				// Verify the key works for authentication
				tests.AssertRequest(
					cl.AccountGetWithResponse(root, createAccessKeyAuth(ak.JSON200.Secret)),
				)(t, http.StatusForbidden)
			})
		}))
	}))
}

// Creates a request editor function for the OpenAPI client that uses an access
// key as a bearer token in the Authorization header.
func createAccessKeyAuth(accessKeyToken string) openapi.RequestEditorFn {
	authHeader := fmt.Sprintf("Bearer %s", accessKeyToken)
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", authHeader)
		return nil
	}
}
