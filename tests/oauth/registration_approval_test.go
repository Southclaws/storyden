package oauth_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	entaccount "github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/oauthclient"
	"github.com/Southclaws/storyden/internal/ent/oauthregistrationapproval"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestOAuthRegistrationApproval(t *testing.T) {
	t.Parallel()
	cfg := oauthConfig(t)

	integration.Test(t, cfg, e2e.Setup(), withOAuthRegistration(t, "approval"), fx.Invoke(func(lc fx.Lifecycle, root context.Context, cl *openapi.ClientWithResponses, sh *e2e.SessionHelper, aw *account_writer.Writer, db *ent.Client, ts *httptest.Server) {
		lc.Append(fx.StartHook(func() {
			adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)
			memberCtx, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			memberSession := sh.WithSession(memberCtx)
			template, key := agentRegistration(t)
			template.Scope = new("openid profile CREATE_POST ADMINISTRATOR")
			template.RegistrationMode = &[]string{"unknown-future-mode", "approval"}
			start := func(t *testing.T) (openapi.OAuthClientRegisterProps, *openapi.OAuthRegistrationPending) {
				t.Helper()
				body := template
				body.ClientName = new("agent-" + xid.New().String())
				response := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body))(t, http.StatusAccepted)
				require.Equal(t, "no-store", response.HTTPResponse.Header.Get("Cache-Control"))
				pending := response.JSON202
				require.NotNil(t, pending)
				require.NotEmpty(t, *pending.RegistrationCode)
				require.NotEmpty(t, *pending.VerificationCode)
				require.Equal(t, 600, *pending.ExpiresIn)
				require.Equal(t, 5, *pending.Interval)
				u, err := url.Parse(*pending.VerificationUriComplete)
				require.NoError(t, err)
				require.Equal(t, *pending.VerificationCode, u.Query().Get("verification_code"))
				require.Equal(t, "/_/resolve/admin/oauth-dcr-approval", u.Path)
				exists, err := db.Account.Query().Where(entaccount.Handle(*body.ClientName)).Exist(root)
				require.NoError(t, err)
				require.False(t, exists)
				exists, err = db.OAuthClient.Query().Where(oauthclient.Name(*body.ClientName)).Exist(root)
				require.NoError(t, err)
				require.False(t, exists)
				return body, pending
			}
			record := func(t *testing.T, code string) *ent.OAuthRegistrationApproval {
				t.Helper()
				hash := sha256.Sum256([]byte(code))
				rec, err := db.OAuthRegistrationApproval.Query().Where(oauthregistrationapproval.RegistrationCodeHash(hex.EncodeToString(hash[:]))).Only(root)
				require.NoError(t, err)
				return rec
			}
			allowPoll := func(t *testing.T, code string) {
				t.Helper()
				rec := record(t, code)
				require.NoError(t, db.OAuthRegistrationApproval.UpdateOneID(rec.ID).SetNextPollAt(time.Now().Add(-time.Second)).Exec(root))
			}
			decide := func(t *testing.T, code string, approved bool) {
				t.Helper()
				tests.AssertRequest(cl.AdminOAuthRegistrationApprovalSubmitWithResponse(root, openapi.AdminOAuthRegistrationApprovalSubmitJSONRequestBody{VerificationCode: code,
					Approved: approved}, adminSession))(t, http.StatusNoContent)
			}
			poll := func(t *testing.T, code string, status int) *openapi.OAuthClientRegisterResponse {
				t.Helper()
				return tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{RegistrationCode: &code}))(t, status)
			}

			t.Run("pending_queue_filters_and_paginates", func(t *testing.T) {
				now := time.Now().UTC()
				expected := []string{}
				excluded := []string{}
				for n := range 56 {
					_, pending := start(t)
					rec := record(t, *pending.RegistrationCode)
					update := db.OAuthRegistrationApproval.UpdateOneID(rec.ID)
					switch n {
					case 0:
						update.SetApprovedAt(now)
					case 1:
						update.SetDeniedAt(now)
					case 2:
						update.SetCancelledAt(now)
					case 3:
						update.SetConsumedAt(now)
					default:
						expected = append(expected, *pending.VerificationCode)
					}
					if n < 4 {
						excluded = append(excluded, *pending.VerificationCode)
					}

					require.NoError(t, update.Exec(root))
				}
				expiredID := xid.New().String()
				_, err := db.OAuthRegistrationApproval.Create().SetRegistrationCodeHash(expiredID).SetVerificationCodeHash(expiredID).SetVerificationCodeDisplay(expiredID).SetMetadata([]byte(`{}`)).SetExpiresAt(now.Add(-time.Minute)).SetNextPollAt(now).Save(root)
				require.NoError(t, err)
				excluded = append(excluded, expiredID)
				tests.AssertRequest(cl.AdminOAuthRegistrationApprovalListWithResponse(root, nil))(t, http.StatusUnauthorized)
				tests.AssertRequest(cl.AdminOAuthRegistrationApprovalListWithResponse(root, nil, memberSession))(t, http.StatusForbidden)
				list := func() []openapi.OAuthRegistrationApprovalReview {
					t.Helper()
					var reviews []openapi.OAuthRegistrationApprovalReview
					for pageNumber := 1; ; pageNumber++ {
						response := tests.AssertRequest(cl.AdminOAuthRegistrationApprovalListWithResponse(root, &openapi.AdminOAuthRegistrationApprovalListParams{Page: new(strconv.Itoa(pageNumber))}, adminSession))(t, http.StatusOK)
						require.Equal(t, "no-store", response.HTTPResponse.Header.Get("Cache-Control"))
						require.Equal(t, 50, response.JSON200.PageSize)
						require.LessOrEqual(t, len(response.JSON200.Registrations), 50)
						reviews = append(reviews, response.JSON200.Registrations...)

						if pageNumber >= response.JSON200.TotalPages {
							return reviews
						}
					}
				}

				reviews := list()
				for _, code := range expected {
					review, found := lo.Find(reviews, func(review openapi.OAuthRegistrationApprovalReview) bool {
						return review.VerificationCode == code
					})
					require.True(t, found, "pending registration %s missing", code)
					require.Nil(t, review.Metadata.RegistrationCode)
					require.Equal(t, []string{"client_credentials"}, *review.Metadata.GrantTypes)
				}
				for _, code := range excluded {
					_, found := lo.Find(reviews, func(review openapi.OAuthRegistrationApprovalReview) bool {
						return review.VerificationCode == code
					})
					require.False(t, found, "non-pending registration %s returned", code)
				}

				decide(t, expected[0], true)
				after := list()
				for n, code := range expected {
					_, found := lo.Find(after, func(review openapi.OAuthRegistrationApprovalReview) bool {
						return review.VerificationCode == code
					})
					require.Equal(t, n != 0, found, "unexpected pending state for %s", code)
				}
			})

			t.Run("discovery_advertises_approval", func(t *testing.T) {
				response, err := http.Get(ts.URL + "/.well-known/oauth-authorization-server")
				require.NoError(t, err)
				defer response.Body.Close()
				require.Equal(t, http.StatusOK, response.StatusCode)
				var metadata struct {
					TokenEndpoint              string   `json:"token_endpoint"`
					RegistrationModesSupported []string `json:"registration_modes_supported"`
				}
				require.NoError(t, json.NewDecoder(response.Body).Decode(&metadata))
				require.Equal(t, "http://localhost:8000/api/oauth/token", metadata.TokenEndpoint)
				require.Equal(t, []string{"immediate", "approval"}, metadata.RegistrationModesSupported)
			})

			t.Run("approval_requires_opt_in", func(t *testing.T) {
				for _, modes := range []*[]string{nil, new([]string{"immediate"}), new([]string{"unknown"})} {
					body := template
					body.RegistrationMode = modes
					response := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body))(t, http.StatusBadRequest)
					require.Equal(t, "invalid_client_metadata", response.JSON400.Error)
				}
			})

			t.Run("review_approval_and_token_subject", func(t *testing.T) {
				body, pending := start(t)
				code := *pending.RegistrationCode
				rec := record(t, code)
				require.NotEqual(t, code, rec.RegistrationCodeHash)
				require.NotContains(t, string(rec.Metadata), code)
				require.NotContains(t, string(rec.Metadata), "sdiat_")
				for _, access := range []struct {
					editor openapi.RequestEditorFn
					status int
				}{{memberSession, http.StatusForbidden}, {bearer(code), http.StatusUnauthorized}} {
					tests.AssertRequest(cl.AdminOAuthRegistrationApprovalGetWithResponse(root, &openapi.AdminOAuthRegistrationApprovalGetParams{VerificationCode: *pending.VerificationCode}, access.editor))(t, access.status)
					tests.AssertRequest(cl.AdminOAuthRegistrationApprovalSubmitWithResponse(root, openapi.AdminOAuthRegistrationApprovalSubmitJSONRequestBody{VerificationCode: *pending.VerificationCode,
						Approved: true}, access.editor))(t, access.status)
				}
				review := tests.AssertRequest(cl.AdminOAuthRegistrationApprovalGetWithResponse(root, &openapi.AdminOAuthRegistrationApprovalGetParams{VerificationCode: *pending.VerificationCode}, adminSession))(t, http.StatusOK)
				require.Equal(t, "no-store", review.HTTPResponse.Header.Get("Cache-Control"))
				require.Equal(t, *body.ClientName, *review.JSON200.Metadata.ClientName)
				require.Equal(t, []string{"client_credentials"}, *review.JSON200.Metadata.GrantTypes)
				require.Equal(t, *body.Jwks, *review.JSON200.Metadata.Jwks)
				require.Nil(t, review.JSON200.Metadata.RegistrationCode)
				require.Nil(t, record(t, code).ApprovedAt)
				decide(t, *pending.VerificationCode, true)
				require.Nil(t, record(t, code).ConsumedAt)
				exists, err := db.Account.Query().Where(entaccount.Handle(*body.ClientName)).Exist(root)
				require.NoError(t, err)
				require.False(t, exists)
				allowPoll(t, code)
				response := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
					RegistrationCode:        &code,
					ClientName:              new("tampered-name"),
					Scope:                   new("ADMINISTRATOR"),
					TokenEndpointAuthMethod: new("none"),
				}))(t, http.StatusCreated)
				require.Nil(t, response.JSON201.ClientSecret)
				require.Equal(t, *body.ClientName, *response.JSON201.ClientName)
				require.Equal(t, "private_key_jwt", response.JSON201.TokenEndpointAuthMethod)
				client, err := db.OAuthClient.Query().Where(oauthclient.ClientID(response.JSON201.ClientId)).Only(root)
				require.NoError(t, err)
				require.NotNil(t, client.AccountID)
				require.NotEqual(t, client.ID, *client.AccountID)
				require.NotEqual(t, admin.ID.String(), client.AccountID.String())
				require.Equal(t, admin.ID.String(), client.RegistrationApprovedByAccountID.String())
				require.Nil(t, client.DcrIatID)
				bot, err := db.Account.Get(root, *client.AccountID)
				require.NoError(t, err)
				require.False(t, bot.Admin)
				require.Equal(t, entaccount.KindBot, bot.Kind)
				require.NotNil(t, record(t, code).ConsumedAt)
				require.Equal(t, "invalid_registration", poll(t, code, http.StatusBadRequest).JSON400.Error)
				tests.AssertRequest(cl.AdminOAuthRegistrationApprovalSubmitWithResponse(root, openapi.AdminOAuthRegistrationApprovalSubmitJSONRequestBody{VerificationCode: *pending.VerificationCode,
					Approved: true}, adminSession))(t, http.StatusBadRequest)
				managed := tests.AssertRequest(cl.AdminOAuthClientGetWithResponse(root, openapi.Identifier(client.ID.String()), adminSession))(t, http.StatusOK)
				require.Equal(t, admin.ID.String(), string(*managed.JSON200.RegistrationApprovedByAccountId))

				assertion := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
					"iss": client.ClientID, "sub": client.ClientID, "aud": "http://localhost:8000/api/oauth/token",
					"iat": time.Now().Unix(), "exp": time.Now().Add(time.Minute).Unix(), "jti": uuid.NewString(),
				})
				assertion.Header["kid"] = "agent-key"
				signed, err := assertion.SignedString(key)
				require.NoError(t, err)
				form := url.Values{"grant_type": {"client_credentials"}, "client_id": {client.ClientID}, "client_assertion_type": {"urn:ietf:params:oauth:client-assertion-type:jwt-bearer"}, "client_assertion": {signed}}
				token := tests.AssertRequest(cl.OAuthTokenWithBodyWithResponse(root, "application/x-www-form-urlencoded", strings.NewReader(form.Encode())))(t, http.StatusOK)
				claims := parseClaims(t, *token.JSON200.AccessToken)
				require.Equal(t, bot.ID.String(), claims["sub"])
				require.NotContains(t, claims["scope"], "ADMINISTRATOR")
				require.Contains(t, claims["scope"], "CREATE_POST")
				category := tests.AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{Name: "agent-posts-" + xid.New().String()}, adminSession))(t, http.StatusOK)
				thread := tests.AssertRequest(cl.ThreadCreateWithResponse(root, nil, openapi.ThreadInitialProps{
					Title:      "Posted by an autonomous agent",
					Body:       new("<p>Registered, authenticated, and posted.</p>"),
					Category:   &category.JSON200.Id,
					Visibility: new(openapi.VisibilityPublished),
				}, bearer(*token.JSON200.AccessToken)))(t, http.StatusOK)
				require.NotNil(t, thread.JSON200)
				require.Equal(t, bot.ID.String(), thread.JSON200.Author.Id)
			})

			t.Run("delegated_private_key_tokens_represent_consenting_account", func(t *testing.T) {
				body := template
				body.ClientName = new("Delegated MCP Connector")
				body.GrantTypes = &[]string{"authorization_code", "refresh_token"}
				body.RedirectUris = &[]string{"https://client.example/callback"}
				body.Scope = new("openid profile")
				registered := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body))(t, http.StatusCreated)
				clientID := registered.JSON201.ClientId
				verifier := strings.Repeat("v", 43)
				location := authorizeRedirect(t, root, ts, adminSession, authorizeRequest{
					ClientID:            clientID,
					RedirectURI:         "https://client.example/callback",
					Scope:               "openid profile",
					State:               uuid.NewString(),
					CodeChallenge:       codeChallenge(verifier),
					CodeChallengeMethod: "S256",
				})
				consentURL, err := url.Parse(location)
				require.NoError(t, err)
				requestID := consentURL.Query().Get("request_id")
				require.NotEmpty(t, requestID)
				tests.AssertRequest(cl.OAuthAuthoriseConsentWithResponse(root, &openapi.OAuthAuthoriseConsentParams{RequestId: (*openapi.OAuthAuthorizationRequestIDQuery)(&requestID)}, adminSession))(t, http.StatusOK)
				consent := tests.AssertRequest(cl.OAuthAuthoriseConsentSubmitWithResponse(root, openapi.OAuthAuthoriseConsentSubmitJSONRequestBody{RequestId: requestID,
					Decision: openapi.OAuthAuthoriseDecisionApprove}, adminSession))(t, http.StatusOK)
				redirect, err := url.Parse(consent.JSON200.Location)
				require.NoError(t, err)
				assertion := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
					"iss": clientID, "sub": clientID, "aud": "http://localhost:8000/api/oauth/token",
					"iat": time.Now().Unix(), "exp": time.Now().Add(time.Minute).Unix(), "jti": uuid.NewString(),
				})
				assertion.Header["kid"] = "agent-key"
				signed, err := assertion.SignedString(key)
				require.NoError(t, err)
				form := url.Values{"grant_type": {"authorization_code"}, "client_id": {clientID}, "code": {redirect.Query().Get("code")}, "redirect_uri": {"https://client.example/callback"}, "code_verifier": {verifier}, "client_assertion_type": {"urn:ietf:params:oauth:client-assertion-type:jwt-bearer"}, "client_assertion": {signed}}
				token := tests.AssertRequest(cl.OAuthTokenWithBodyWithResponse(root, "application/x-www-form-urlencoded", strings.NewReader(form.Encode())))(t, http.StatusOK)
				require.Equal(t, admin.ID.String(), parseClaims(t, *token.JSON200.AccessToken)["sub"])
			})

			t.Run("pending_and_increased_poll_interval", func(t *testing.T) {
				_, pending := start(t)
				code := *pending.RegistrationCode
				early := poll(t, code, http.StatusTooManyRequests)
				require.Equal(t, "10", early.HTTPResponse.Header.Get("Retry-After"))
				require.Equal(t, "no-store", early.HTTPResponse.Header.Get("Cache-Control"))
				require.Equal(t, 10, record(t, code).PollIntervalSeconds)
				allowPoll(t, code)
				waiting := poll(t, code, http.StatusAccepted)
				require.NotNil(t, waiting.JSON202)
				require.Nil(t, waiting.JSON202.RegistrationCode)
				require.Equal(t, 10, record(t, code).PollIntervalSeconds)
				require.Equal(t, "15", poll(t, code, http.StatusTooManyRequests).HTTPResponse.Header.Get("Retry-After"))
			})

			t.Run("denied_cancelled_unknown_and_expired", func(t *testing.T) {
				_, denied := start(t)
				decide(t, *denied.VerificationCode, false)
				require.Equal(t, "access_denied", poll(t, *denied.RegistrationCode, http.StatusBadRequest).JSON400.Error)
				_, cancelled := start(t)
				tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{RegistrationCode: cancelled.RegistrationCode,
					CancelRegistration: new(true)}))(t, http.StatusNoContent)
				require.Equal(t, "invalid_registration", poll(t, *cancelled.RegistrationCode, http.StatusBadRequest).JSON400.Error)
				tests.AssertRequest(cl.AdminOAuthRegistrationApprovalSubmitWithResponse(root, openapi.AdminOAuthRegistrationApprovalSubmitJSONRequestBody{VerificationCode: *cancelled.VerificationCode,
					Approved: true}, adminSession))(t, http.StatusBadRequest)
				tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{RegistrationCode: cancelled.RegistrationCode,
					CancelRegistration: new(true)}))(t, http.StatusBadRequest)
				require.Equal(t, "invalid_registration", poll(t, "unknown-registration", http.StatusBadRequest).JSON400.Error)
				expiredCode := "expired-" + xid.New().String()
				hash := sha256.Sum256([]byte(expiredCode))
				_, err := db.OAuthRegistrationApproval.Create().SetRegistrationCodeHash(hex.EncodeToString(hash[:])).SetVerificationCodeHash(expiredCode).SetVerificationCodeDisplay("ABCD-EFGH").SetMetadata([]byte(`{}`)).SetExpiresAt(time.Now().Add(-time.Minute)).SetNextPollAt(time.Now()).Save(root)
				require.NoError(t, err)
				require.Equal(t, "expired_token", poll(t, expiredCode, http.StatusBadRequest).JSON400.Error)
				_, approved := start(t)
				decide(t, *approved.VerificationCode, true)
				tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{RegistrationCode: approved.RegistrationCode,
					CancelRegistration: new(true)}))(t, http.StatusBadRequest)
			})

			t.Run("iat_bypasses_approval", func(t *testing.T) {
				issued := tests.AssertRequest(cl.AdminOAuthDynamicRegistrationAccessTokenCreateWithResponse(root, openapi.AdminOAuthDynamicRegistrationAccessTokenCreateJSONRequestBody{Label: "preapproved",
					ExpiresAt: time.Now().Add(time.Hour)}, adminSession))(t, http.StatusCreated)
				body := template
				body.ClientName = new("iat-" + xid.New().String())
				body.RegistrationMode = nil
				response := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body, bearer(issued.JSON201.Token)))(t, http.StatusCreated)
				client, err := db.OAuthClient.Query().Where(oauthclient.ClientID(response.JSON201.ClientId)).Only(root)
				require.NoError(t, err)
				require.NotNil(t, client.DcrIatID)
				require.Nil(t, client.RegistrationApprovedByAccountID)
				body.ClientName = new("invalid-" + xid.New().String())
				body.RegistrationMode = &[]string{"approval"}
				tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body, bearer(issued.JSON201.Token)))(t, http.StatusUnauthorized)
			})

			t.Run("concurrent_completion_creates_one_identity", func(t *testing.T) {
				body, pending := start(t)
				decide(t, *pending.VerificationCode, true)
				allowPoll(t, *pending.RegistrationCode)
				var wg sync.WaitGroup
				statuses := make(chan int, 6)
				errors := make(chan error, 6)
				for range 6 {
					wg.Add(1)
					go func() {
						defer wg.Done()
						response, err := cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{RegistrationCode: pending.RegistrationCode})
						if err != nil {
							errors <- err
							return
						}
						statuses <- response.StatusCode()
					}()
				}
				wg.Wait()
				close(statuses)
				close(errors)
				for err := range errors {
					require.NoError(t, err)
				}
				created := 0
				for status := range statuses {
					require.Contains(t, []int{201, 400, 429}, status)
					if status == 201 {
						created++
					}
				}
				require.Equal(t, 1, created)
				count, err := db.Account.Query().Where(entaccount.Handle(*body.ClientName)).Count(root)
				require.NoError(t, err)
				require.Equal(t, 1, count)
			})
		}))
	}))
}
