package oauth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/oauthregistrationapproval"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestOAuthRuntimeRegistrationSettings(t *testing.T) {
	if tests.IsSharedPostgresDatabase() {
		t.Skip("skipping global registration policy mutations on shared postgres database")
	}

	t.Parallel()
	cfg := oauthConfig(t)
	cfg.PublicWebAddress = url.URL{Scheme: "https", Host: "frontend.example.com", Path: "/community/"}

	integration.Test(t, cfg, e2e.Setup(), fx.Invoke(func(lc fx.Lifecycle, root context.Context, cl *openapi.ClientWithResponses, sh *e2e.SessionHelper, aw *account_writer.Writer, db *ent.Client, ts *httptest.Server) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)
			memberCtx, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			memberSession := sh.WithSession(memberCtx)
			patch := func(value openapi.OAuthServiceSettings) openapi.AdminSettingsUpdateJSONRequestBody {
				return openapi.AdminSettingsUpdateJSONRequestBody{Services: &openapi.AdminSettingsServiceProps{Oauth: &value}}
			}
			update := func(value openapi.OAuthServiceSettings) *openapi.OAuthServiceSettings {
				t.Helper()
				response := tests.AssertRequest(cl.AdminSettingsUpdateWithResponse(root, patch(value), adminSession))(t, http.StatusOK)
				require.NotNil(t, response.JSON200.Services.Oauth)
				return response.JSON200.Services.Oauth
			}
			discovery := func(enabled, approval bool) {
				t.Helper()
				for _, path := range []string{"/.well-known/openid-configuration", "/.well-known/oauth-authorization-server"} {
					response, err := http.Get(ts.URL + path)
					require.NoError(t, err)
					require.Equal(t, http.StatusOK, response.StatusCode)
					require.Equal(t, "no-cache", response.Header.Get("Cache-Control"))
					var body struct {
						RegistrationEndpoint       string   `json:"registration_endpoint"`
						RegistrationModesSupported []string `json:"registration_modes_supported"`
					}
					err = json.NewDecoder(response.Body).Decode(&body)
					response.Body.Close()
					require.NoError(t, err)
					if enabled {
						require.Equal(t, "http://localhost:8000/api/oauth/register", body.RegistrationEndpoint)
						modes := []string{"immediate"}
						if approval {
							modes = append(modes, "approval")
						}
						require.Equal(t, modes, body.RegistrationModesSupported)
					} else {
						require.Empty(t, body.RegistrationEndpoint)
						require.Empty(t, body.RegistrationModesSupported)
					}
				}
			}

			initial := tests.AssertRequest(cl.AdminSettingsGetWithResponse(root, adminSession))(t, http.StatusOK).JSON200.Services.Oauth
			require.NotNil(t, initial)
			require.False(t, *initial.DynamicRegistrationEnabled)
			require.Equal(t, "disabled", string(*initial.AutonomousRegistrationMode))
			require.Equal(t, 600, *initial.RegistrationApprovalTtl)
			require.Empty(t, *initial.RegistrationApprovalUrl)
			discovery(false, false)
			agent, _ := agentRegistration(t)
			rejected := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, agent))(t, http.StatusBadRequest)
			require.Equal(t, "temporarily_unavailable", rejected.JSON400.Error)

			tests.AssertRequest(cl.AdminSettingsUpdateWithResponse(root, patch(openapi.OAuthServiceSettings{DynamicRegistrationEnabled: new(true)}), memberSession))(t, http.StatusForbidden)
			for _, invalid := range []openapi.OAuthServiceSettings{
				{AutonomousRegistrationMode: new(openapi.OAuthAutonomousRegistrationMode("unknown"))},
				{RegistrationApprovalTtl: new(0)},
				{RegistrationApprovalTtl: new(-1)},
				{RegistrationApprovalTtl: new(9223372037)},
				{RegistrationApprovalUrl: new("/approve")},
				{RegistrationApprovalUrl: new("javascript:alert(1)")},
				{RegistrationApprovalUrl: new("https://user:password@example.com/approve")},
				{RegistrationApprovalUrl: new("https://example.com/approve#fragment")},
			} {
				tests.AssertRequest(cl.AdminSettingsUpdateWithResponse(root, patch(invalid), adminSession))(t, http.StatusBadRequest)
			}
			discovery(false, false)

			update(openapi.OAuthServiceSettings{DynamicRegistrationEnabled: new(true)})
			for _, mode := range []string{"protected", "open", "disabled", "approval"} {
				update(openapi.OAuthServiceSettings{AutonomousRegistrationMode: new(openapi.OAuthAutonomousRegistrationMode(mode))})
				discovery(true, mode == "approval")
				delegated := openapi.OAuthClientRegisterJSONRequestBody{ClientName: new("connector-" + xid.New().String()), RedirectUris: new([]string{"https://client.example/callback"}), TokenEndpointAuthMethod: new("none")}
				tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, delegated))(t, http.StatusCreated)
				exists, err := db.Account.Query().Where(account.Handle(*delegated.ClientName)).Exist(root)
				require.NoError(t, err)
				require.False(t, exists)
				body, _ := agentRegistration(t)
				switch mode {
				case "protected":
					tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body))(t, http.StatusUnauthorized)
				case "open":
					tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body))(t, http.StatusCreated)
				case "disabled":
					tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body))(t, http.StatusBadRequest)
				}
			}

			custom := update(openapi.OAuthServiceSettings{RegistrationApprovalTtl: new(120), RegistrationApprovalUrl: new("https://example.com/approve?from=agent")})
			require.True(t, *custom.DynamicRegistrationEnabled)
			require.Equal(t, "approval", string(*custom.AutonomousRegistrationMode))
			agent.RegistrationMode = new([]string{"approval"})
			pending := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, agent))(t, http.StatusAccepted).JSON202
			require.Equal(t, 120, *pending.ExpiresIn)
			require.Equal(t, "https://example.com/approve?from=agent", *pending.VerificationUri)
			require.Equal(t, "https://example.com/approve?from=agent&verification_code="+*pending.VerificationCode, *pending.VerificationUriComplete)
			record, err := db.OAuthRegistrationApproval.Query().Where(oauthregistrationapproval.VerificationCodeDisplay(*pending.VerificationCode)).Only(root)
			require.NoError(t, err)
			tests.AssertRequest(cl.AdminOAuthRegistrationApprovalSubmitWithResponse(root, openapi.AdminOAuthRegistrationApprovalSubmitJSONRequestBody{VerificationCode: *pending.VerificationCode, Approved: true}, adminSession))(t, http.StatusNoContent)
			update(openapi.OAuthServiceSettings{AutonomousRegistrationMode: new(openapi.OAuthAutonomousRegistrationMode("disabled"))})
			poll := openapi.OAuthClientRegisterJSONRequestBody{RegistrationCode: pending.RegistrationCode}
			denied := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, poll))(t, http.StatusBadRequest)
			require.Equal(t, "access_denied", denied.JSON400.Error)
			exists, err := db.Account.Query().Where(account.Handle(*agent.ClientName)).Exist(root)
			require.NoError(t, err)
			require.False(t, exists)

			update(openapi.OAuthServiceSettings{AutonomousRegistrationMode: new(openapi.OAuthAutonomousRegistrationMode("approval")), RegistrationApprovalTtl: new(300), RegistrationApprovalUrl: new("")})
			existing, err := db.OAuthRegistrationApproval.Get(root, record.ID)
			require.NoError(t, err)
			require.True(t, record.ExpiresAt.Equal(existing.ExpiresAt))
			require.NoError(t, db.OAuthRegistrationApproval.UpdateOneID(record.ID).SetNextPollAt(time.Now().Add(-time.Second)).Exec(root))
			tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, poll))(t, http.StatusCreated)
			newAgent, _ := agentRegistration(t)
			newAgent.RegistrationMode = new([]string{"approval"})
			fresh := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, newAgent))(t, http.StatusAccepted).JSON202
			require.Equal(t, 300, *fresh.ExpiresIn)
			require.Equal(t, "https://frontend.example.com/community/_/resolve/admin/oauth-dcr-approval", *fresh.VerificationUri)
			require.Equal(t, "https://frontend.example.com/community/_/resolve/admin/oauth-dcr-approval?verification_code="+*fresh.VerificationCode, *fresh.VerificationUriComplete)

			update(openapi.OAuthServiceSettings{DynamicRegistrationEnabled: new(false)})
			discovery(false, false)
			rejected = tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{RegistrationCode: fresh.RegistrationCode}))(t, http.StatusBadRequest)
			require.Equal(t, "temporarily_unavailable", rejected.JSON400.Error)
			stored, err := db.Setting.Get(root, settings.StorydenPrimarySettingsKey)
			require.NoError(t, err)
			var persisted settings.Settings
			require.NoError(t, json.Unmarshal([]byte(stored.Value), &persisted))
			value := persisted.Services.OrZero().OAuth.OrZero()
			require.False(t, value.DynamicRegistrationEnabled.OrZero())
			require.Equal(t, "approval", value.AutonomousRegistrationMode.OrZero().String())
			require.Equal(t, 5*time.Minute, value.RegistrationApprovalTTL.OrZero())
			require.Empty(t, value.RegistrationApprovalURL.OrZero())
		}))
	}))
}
