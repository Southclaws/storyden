package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/go-jose/go-jose/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

type enrollmentServer struct {
	t                    *testing.T
	server               *httptest.Server
	mu                   sync.Mutex
	keys                 jose.JSONWebKeySet
	registrations        int
	polls                int
	tokens               int
	tokenScopes          []string
	grantedScope         string
	approval             bool
	failToken            bool
	truncateRegistration bool
	decision             string
	assertions           map[string]bool
}

func newEnrollmentServer(t *testing.T, approval bool) *enrollmentServer {
	f := &enrollmentServer{t: t, approval: approval, assertions: map[string]bool{}}
	f.server = httptest.NewServer(http.HandlerFunc(f.serve))
	t.Cleanup(f.server.Close)
	return f
}

func (f *enrollmentServer) serve(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	switch r.URL.Path {
	case "/.well-known/openid-configuration":
		json.NewEncoder(w).Encode(map[string]any{"issuer": f.server.URL, "registration_endpoint": f.server.URL + "/api/enroll", "token_endpoint": f.server.URL + "/api/token"})
	case "/api/enroll":
		var body map[string]json.RawMessage
		require.NoError(f.t, json.NewDecoder(r.Body).Decode(&body))
		if _, ok := body["registration_code"]; ok {
			f.polls++
			require.JSONEq(f.t, `"poll-secret"`, string(body["registration_code"]))
			if cancel, ok := body["cancel_registration"]; ok {
				require.JSONEq(f.t, `true`, string(cancel))
				require.Len(f.t, body, 2)
				w.WriteHeader(http.StatusNoContent)
				return
			}
			require.Len(f.t, body, 1)
			switch f.decision {
			case "unavailable":
				w.WriteHeader(http.StatusServiceUnavailable)
				fmt.Fprint(w, `{"error":"temporarily_unavailable"}`)
				return
			case "disconnect":
				conn, _, err := w.(http.Hijacker).Hijack()
				require.NoError(f.t, err)
				conn.Close()
				return
			case "truncated":
				w.WriteHeader(http.StatusAccepted)
				fmt.Fprint(w, `{`)
				return
			case "denied":
				w.WriteHeader(400)
				fmt.Fprint(w, `{"error":"access_denied"}`)
				return
			case "slow":
				w.Header().Set("Retry-After", "25")
				w.WriteHeader(429)
				return
			case "approved":
			default:
				w.WriteHeader(202)
				fmt.Fprint(w, `{}`)
				return
			}
		} else {
			f.registrations++
			require.JSONEq(f.t, `["client_credentials"]`, string(body["grant_types"]))
			require.JSONEq(f.t, `"private_key_jwt"`, string(body["token_endpoint_auth_method"]))
			require.JSONEq(f.t, `["approval"]`, string(body["registration_mode"]))
			require.NoError(f.t, json.Unmarshal(body["jwks"], &f.keys))
			require.Len(f.t, f.keys.Keys, 1)
			require.True(f.t, f.keys.Keys[0].IsPublic())
			if f.approval {
				w.WriteHeader(202)
				fmt.Fprint(w, `{"registration_code":"poll-secret","verification_code":"ABCD-EFGH","verification_uri":"https://example.com/admin/verify","expires_in":600,"interval":5}`)
				return
			}
		}
		w.WriteHeader(201)
		if f.truncateRegistration {
			fmt.Fprint(w, `{`)
			return
		}
		fmt.Fprint(w, `{"client_id":"agent-client","token_endpoint_auth_method":"private_key_jwt"}`)
	case "/api/token":
		f.tokens++
		if f.failToken {
			w.WriteHeader(503)
			fmt.Fprint(w, `{}`)
			return
		}
		require.NoError(f.t, r.ParseForm())
		f.tokenScopes = append(f.tokenScopes, r.Form.Get("scope"))
		require.Equal(f.t, "client_credentials", r.Form.Get("grant_type"))
		require.Equal(f.t, "agent-client", r.Form.Get("client_id"))
		require.Equal(f.t, "urn:ietf:params:oauth:client-assertion-type:jwt-bearer", r.Form.Get("client_assertion_type"))
		require.Empty(f.t, r.Form.Get("refresh_token"))
		claims := &jwt.RegisteredClaims{}
		parsed, err := jwt.ParseWithClaims(r.Form.Get("client_assertion"), claims, func(token *jwt.Token) (any, error) {
			require.Equal(f.t, f.keys.Keys[0].KeyID, token.Header["kid"])
			return f.keys.Keys[0].Key, nil
		}, jwt.WithValidMethods([]string{"RS256"}), jwt.WithAudience(f.server.URL+"/api/token"), jwt.WithIssuer("agent-client"), jwt.WithSubject("agent-client"), jwt.WithExpirationRequired())
		require.NoError(f.t, err)
		require.True(f.t, parsed.Valid)
		require.NotEmpty(f.t, claims.ID)
		require.False(f.t, f.assertions[claims.ID])
		f.assertions[claims.ID] = true
		require.InDelta(f.t, 60, claims.ExpiresAt.Sub(claims.IssuedAt.Time).Seconds(), 1)
		scope := f.grantedScope
		if scope == "" {
			scope = "CREATE_POST"
		}
		fmt.Fprintf(w, `{"access_token":"bot-token-%d","token_type":"Bearer","expires_in":900,"scope":%q}`, f.tokens, scope)
	case "/api/threads":
		require.Equal(f.t, "Bearer bot-token-2", r.Header.Get("Authorization"))
		fmt.Fprint(w, `{"threads":[]}`)
	default:
		http.NotFound(w, r)
	}
}

func registrationStore(t *testing.T) *config.Store {
	store := config.NewFileStoreAt(filepath.Join(t.TempDir(), "config.yaml"))
	cfg := config.New()
	cfg.CurrentContext = "human"
	cfg.UpsertContext("human", config.Context{APIURL: "https://human.example", AuthType: config.AuthStorageFile, Auth: &config.Auth{Method: config.AuthMethodAccessKey, AccessToken: "human-secret"}})
	require.NoError(t, store.Save(cfg))
	store.SelectedContext = "bot"
	return store
}

func enroll(t *testing.T, f *enrollmentServer, store *config.Store) *RegistrationResult {
	result, err := Register(context.Background(), store, RegistrationOptions{Endpoint: f.server.URL, Name: "bot", Handle: "test-bot", Storage: config.AuthStorageFile, Scope: "CREATE_POST"})
	require.NoError(t, err)
	return result
}

func allowPoll(t *testing.T, store *config.Store) {
	require.NoError(t, store.Update(context.Background(), func(cfg *config.Config) error {
		cfg.Contexts["bot"].Auth.Registration.NextPollAt = time.Now().Add(-time.Second)
		return nil
	}))
}

func TestAutonomousRegistrationAndTokenRenewal(t *testing.T) {
	f := newEnrollmentServer(t, false)
	store := registrationStore(t)
	result := enroll(t, f, store)
	require.Equal(t, "ready", result.State)
	cfg, err := store.Load()
	require.NoError(t, err)
	require.Equal(t, "human", cfg.CurrentContext)
	require.Equal(t, "human-secret", cfg.Contexts["human"].Auth.AccessToken)
	require.NotEmpty(t, cfg.Contexts["bot"].Auth.PrivateKey)
	require.NoError(t, store.Update(context.Background(), func(cfg *config.Config) error {
		cfg.Contexts["bot"].Auth.ExpiresAt = time.Now().Add(-time.Second)
		return nil
	}))
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			other := config.NewFileStoreAt(store.Path())
			other.SelectedContext = "bot"
			auth, err := Credentials(context.Background(), other, false)
			require.NoError(t, err)
			require.Equal(t, "bot-token-2", auth.AccessToken)
		})
	}
	wg.Wait()
	require.Equal(t, 2, f.tokens)
	require.Equal(t, 1, f.registrations)
	_, err = Register(context.Background(), store, RegistrationOptions{Endpoint: f.server.URL, Name: "bot", Handle: "test-bot", Storage: config.AuthStorageFile})
	require.ErrorContains(t, err, "already exists")
	require.Equal(t, 1, f.registrations)
	store.SelectedContext = "missing"
	_, err = Credentials(context.Background(), store, false)
	require.ErrorContains(t, err, "missing")
}

func TestApprovalRestartAndThrottling(t *testing.T) {
	f := newEnrollmentServer(t, true)
	store := registrationStore(t)
	result := enroll(t, f, store)
	require.Equal(t, "pending", result.State)
	require.Zero(t, f.tokens)
	data, err := json.Marshal(result)
	require.NoError(t, err)
	require.NotContains(t, string(data), "poll-secret")
	require.NotContains(t, string(data), "PRIVATE KEY")
	store = config.NewFileStoreAt(store.Path())
	store.SelectedContext = "bot"
	result, err = CheckRegistration(context.Background(), store, false)
	require.NoError(t, err)
	require.Equal(t, "pending", result.State)
	require.Zero(t, f.polls)
	allowPoll(t, store)
	f.decision = "slow"
	result, err = CheckRegistration(context.Background(), store, false)
	require.NoError(t, err)
	require.Greater(t, time.Until(result.NextPollAt), 24*time.Second)
	f.decision = "approved"
	allowPoll(t, store)
	result, err = CheckRegistration(context.Background(), store, false)
	require.NoError(t, err)
	require.Equal(t, "ready", result.State)
	require.Equal(t, 1, f.registrations)
	require.Equal(t, 2, f.polls)
	require.Equal(t, 1, f.tokens)
	require.Empty(t, result.Code)
}

func TestApprovalDenialAndExpiry(t *testing.T) {
	for _, decision := range []string{"denied", "expired"} {
		t.Run(decision, func(t *testing.T) {
			f := newEnrollmentServer(t, true)
			store := registrationStore(t)
			enroll(t, f, store)
			f.decision = decision
			allowPoll(t, store)
			if decision == "expired" {
				require.NoError(t, store.Update(context.Background(), func(cfg *config.Config) error {
					cfg.Contexts["bot"].Auth.Registration.ExpiresAt = time.Now().Add(-time.Second)
					return nil
				}))
			}
			result, err := CheckRegistration(context.Background(), store, false)
			require.Error(t, err)
			require.NotEqual(t, "pending", result.State)
			require.Empty(t, result.Code)
			require.Zero(t, f.tokens)
		})
	}
}

func TestRegistrationDoesNotFollowRedirectWithCredentials(t *testing.T) {
	var received bool
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { received = true }))
	defer target.Close()
	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusTemporaryRedirect)
	}))
	defer origin.Close()
	response, err := oauthRequest(context.Background(), origin.URL, origin.URL, "application/json", nil, "secret-iat")
	require.NoError(t, err)
	response.Body.Close()
	require.Equal(t, 307, response.StatusCode)
	require.False(t, received)
	_, err = oauthRequest(context.Background(), origin.URL, target.URL, "application/json", nil, "secret-iat")
	require.Error(t, err)
	require.False(t, received)
}

func TestRegistrationRecoveryBoundaries(t *testing.T) {
	t.Run("token_failure_retains_registered_client", func(t *testing.T) {
		f := newEnrollmentServer(t, false)
		store := registrationStore(t)
		f.failToken = true
		result, err := Register(context.Background(), store, RegistrationOptions{Endpoint: f.server.URL, Name: "bot", Handle: "test-bot", Storage: config.AuthStorageFile})
		require.Error(t, err)
		require.Equal(t, "registered", result.State)
		require.Equal(t, "agent-client", result.ClientID)
		f.failToken = false
		result, err = CheckRegistration(context.Background(), store, false)
		require.NoError(t, err)
		require.Equal(t, "ready", result.State)
		require.Equal(t, 1, f.registrations)
	})
	t.Run("lost_registration_response_is_not_replayed", func(t *testing.T) {
		f := newEnrollmentServer(t, false)
		store := registrationStore(t)
		f.truncateRegistration = true
		result, err := Register(context.Background(), store, RegistrationOptions{Endpoint: f.server.URL, Name: "bot", Handle: "test-bot", Storage: config.AuthStorageFile})
		require.Error(t, err)
		require.Equal(t, "outcome_unknown", result.State)
		_, err = CheckRegistration(context.Background(), store, false)
		require.ErrorContains(t, err, "outcome_unknown")
		require.Equal(t, 1, f.registrations)
		require.Zero(t, f.tokens)
	})
}

func TestApprovalRecoveryAfterFailedPoll(t *testing.T) {
	for _, failure := range []string{"unavailable", "disconnect", "truncated"} {
		for _, action := range []string{"resume", "cancel"} {
			t.Run(failure+"/"+action, func(t *testing.T) {
				f := newEnrollmentServer(t, true)
				store := registrationStore(t)
				enroll(t, f, store)
				allowPoll(t, store)
				f.decision = failure
				result, err := CheckRegistration(context.Background(), store, false)
				require.Error(t, err)
				require.Equal(t, "outcome_unknown", result.State)

				// A new CLI invocation must reuse the original code and key.
				store = config.NewFileStoreAt(store.Path())
				store.SelectedContext = "bot"
				f.decision = "approved"
				result, err = CheckRegistration(context.Background(), store, action == "cancel")
				require.NoError(t, err)
				if action == "cancel" {
					require.Equal(t, "cancelled", result.State)
					require.Zero(t, f.tokens)
				} else {
					require.Equal(t, "ready", result.State)
					require.Equal(t, 1, f.tokens)
				}
				require.Equal(t, 1, f.registrations)
				require.Equal(t, 2, f.polls)
				require.Empty(t, result.Code)
			})
		}
	}
}

func TestClientCredentialsRetainRequestedScope(t *testing.T) {
	for _, requested := range []string{"", "CREATE_POST READ_PUBLISHED_LIBRARY"} {
		t.Run(requested, func(t *testing.T) {
			f := newEnrollmentServer(t, false)
			store := registrationStore(t)
			result, err := Register(context.Background(), store, RegistrationOptions{Endpoint: f.server.URL, Name: "bot", Handle: "test-bot", Storage: config.AuthStorageFile, Scope: requested})
			require.NoError(t, err)
			require.Equal(t, "CREATE_POST", result.Scope)
			f.mu.Lock()
			f.grantedScope = "CREATE_POST READ_PUBLISHED_LIBRARY"
			f.mu.Unlock()
			reloaded := config.NewFileStoreAt(store.Path())
			reloaded.SelectedContext = "bot"
			renewed, err := Credentials(context.Background(), reloaded, true)
			require.NoError(t, err)
			require.Equal(t, "CREATE_POST READ_PUBLISHED_LIBRARY", renewed.Scope)
			require.Equal(t, requested, renewed.RequestedScope)
			require.Equal(t, []string{requested, requested}, f.tokenScopes)
		})
	}
}
