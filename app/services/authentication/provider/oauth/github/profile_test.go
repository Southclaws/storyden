package github

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/config"
)

type profileTransport func(*http.Request) (*http.Response, error)

func (f profileTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestFetchProfile(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, body string
		status     int
		wantErr    bool
	}{
		{"valid", `{"id":9007199254740993,"login":"Scribe","name":null,"email":null}`, 200, false},
		{"unauthorized", `{"message":"secret-token"}`, 401, true},
		{"redirect", ``, 302, true},
		{"missing id", `{"login":"Scribe"}`, 200, true},
		{"missing login", `{"id":1}`, 200, true},
		{"empty login", `{"id":1,"login":" "}`, 200, true},
		{"null", `null`, 200, true},
		{"invalid json", `{`, 200, true},
		{"oversized", strings.Repeat(" ", 1<<20) + `{"id":1,"login":"Scribe"}`, 200, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := &http.Client{Transport: profileTransport(func(r *http.Request) (*http.Response, error) {
				require.Equal(t, "https://api.github.com/user", r.URL.String())
				require.Equal(t, http.MethodGet, r.Method)
				require.Equal(t, "Bearer secret-token", r.Header.Get("Authorization"))
				require.Equal(t, "application/vnd.github+json", r.Header.Get("Accept"))
				require.Equal(t, "2022-11-28", r.Header.Get("X-GitHub-Api-Version"))
				return &http.Response{StatusCode: tc.status, Body: io.NopCloser(strings.NewReader(tc.body)), Header: make(http.Header)}, nil
			})}
			got, err := fetchProfile(t.Context(), client, "secret-token")
			if tc.wantErr {
				require.Error(t, err)
				require.NotContains(t, err.Error(), "secret-token")
				return
			}
			require.NoError(t, err)
			require.Equal(t, int64(9007199254740993), got.ID)
			require.Equal(t, "Scribe", *got.Login)
			require.Nil(t, got.Name)
			require.Nil(t, got.Email)
		})
	}
}

func TestFetchProfileCancellation(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	client := &http.Client{Transport: profileTransport(func(r *http.Request) (*http.Response, error) {
		require.ErrorIs(t, r.Context().Err(), context.Canceled)
		return nil, r.Context().Err()
	})}
	_, err := fetchProfile(ctx, client, "token")
	require.ErrorIs(t, err, context.Canceled)
}

func TestFetchProfileRejectsRedirect(t *testing.T) {
	t.Parallel()
	provider, err := New(config.Config{}, nil, nil)
	require.NoError(t, err)
	requests := 0
	provider.client.Transport = profileTransport(func(r *http.Request) (*http.Response, error) {
		requests++
		require.Equal(t, "https://api.github.com/user", r.URL.String())
		return &http.Response{StatusCode: http.StatusFound, Header: http.Header{"Location": {"https://api.github.com/other"}}, Body: io.NopCloser(strings.NewReader(""))}, nil
	})
	_, err = fetchProfile(t.Context(), provider.client, "token")
	require.ErrorContains(t, err, "HTTP 302")
	require.Equal(t, 1, requests)
}
