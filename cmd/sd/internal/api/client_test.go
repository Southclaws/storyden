package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Run("discovers API base under api path from canonical endpoint", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			if request.URL.Path != "/.well-known/openid-configuration" {
				http.NotFound(w, request)
				return
			}

			writeDiscovery(t, w, request)
		}))
		defer server.Close()

		client, err := NewClient(context.Background(), server.URL)
		r.NoError(err)

		a.Equal(server.URL, client.Endpoint)
		a.Equal(server.URL+"/api", client.BaseURL)
		a.Equal("https://storyden.example", client.Discovery.Issuer)
	})

	t.Run("stores canonical endpoint even when user types api path", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			if request.URL.Path != "/.well-known/openid-configuration" {
				http.NotFound(w, request)
				return
			}

			writeDiscovery(t, w, request)
		}))
		defer server.Close()

		client, err := NewClient(context.Background(), server.URL+"/api")
		r.NoError(err)

		a.Equal(server.URL, client.Endpoint)
		a.Equal(server.URL+"/api", client.BaseURL)
	})

	t.Run("formats rate limit errors during OAuth discovery", func(t *testing.T) {
		r := require.New(t)

		resetAt := time.Now().Add(3 * time.Minute).UTC().Format(time.RFC3339)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			if request.URL.Path != "/.well-known/openid-configuration" {
				http.NotFound(w, request)
				return
			}

			w.Header().Set("X-RateLimit-Limit", "100")
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", resetAt)
			http.Error(w, "too many requests", http.StatusTooManyRequests)
		}))
		defer server.Close()

		_, err := NewClient(context.Background(), server.URL)
		r.Error(err)
		r.ErrorContains(err, "Rate limit exceeded.")
		r.ErrorContains(err, "Limit: 100 requests")
		r.ErrorContains(err, "Reset time:")
		r.ErrorContains(err, "(in about")
	})
}

func TestCanonicalEndpoint(t *testing.T) {
	t.Run("removes path query and fragment", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		endpoint, err := CanonicalEndpoint("https://forum.example.com/api?x=y#fragment")
		r.NoError(err)

		a.Equal("https://forum.example.com", endpoint)
	})
}

func writeDiscovery(t *testing.T, w http.ResponseWriter, request *http.Request) {
	t.Helper()

	w.Header().Set("Content-Type", "application/json")
	require.NoError(t, json.NewEncoder(w).Encode(map[string]any{
		"issuer":                                "https://storyden.example",
		"authorization_endpoint":                serverURL(request) + "/api/oauth/authorize",
		"device_authorization_endpoint":         serverURL(request) + "/api/oauth/device_authorization",
		"token_endpoint":                        serverURL(request) + "/api/oauth/token",
		"userinfo_endpoint":                     serverURL(request) + "/api/oauth/userinfo",
		"jwks_uri":                              serverURL(request) + "/api/oauth/jwks",
		"response_types_supported":              []string{"code"},
		"grant_types_supported":                 []string{"urn:ietf:params:oauth:grant-type:device_code"},
		"scopes_supported":                      []string{"openid", "profile", "offline_access"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
		"code_challenge_methods_supported":      []string{"S256"},
	}))
}

func serverURL(r *http.Request) string {
	return "http://" + r.Host
}
