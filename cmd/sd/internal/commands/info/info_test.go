package info

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func TestValidateFormat(t *testing.T) {
	r := require.New(t)

	r.NoError(validateFormat("plain"))
	r.NoError(validateFormat("json"))
	r.ErrorContains(validateFormat("xml"), "--format must be one of: plain, json")
}

func TestRenderPlain(t *testing.T) {
	r := require.New(t)

	result := instanceInfo{
		Context:  "local",
		Endpoint: "https://example.com",
		BaseURL:  "https://example.com/api",
		Info: openapi.Info{
			AccentColour:       "#ffcc00",
			ApiAddress:         "https://api.example.com",
			AuthenticationMode: "email",
			Capabilities:       openapi.InstanceCapabilityList{"oauth", "plugins"},
			Content:            `<body><p>Welcome <strong>home</strong></p><ul><li>Build things</li></ul></body>`,
			Description:        "Community hub",
			OnboardingStatus:   "requires_category",
			RegistrationMode:   "invitation",
			Title:              "Storyden",
			WebAddress:         "https://example.com",
		},
	}

	var out bytes.Buffer
	r.NoError(renderPlain(&out, result))

	text := out.String()
	r.Contains(text, "Context")
	r.Contains(text, "local")
	r.Contains(text, "Endpoint")
	r.Contains(text, "https://example.com")
	r.Contains(text, "Web address")
	r.Contains(text, "API address")
	r.Contains(text, "https://api.example.com")
	r.Contains(text, "Storyden")
	r.Contains(text, "Capabilities")
	r.Contains(text, "- Oauth")
	r.Contains(text, "- Plugins")
	r.Contains(text, "Overview")
	r.Contains(text, "Welcome **home**")
	r.Contains(text, "- Build things")
	r.NotContains(text, "Onboarding")
	r.NotContains(text, "Metadata keys")
	r.NotContains(text, "<body>")
}

func TestFetchInfo(t *testing.T) {
	r := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		r.Equal("/info", request.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		r.NoError(json.NewEncoder(w).Encode(map[string]any{
			"accent_colour":       "#5b7fff",
			"api_address":         "https://api.test",
			"authentication_mode": "email",
			"capabilities":        []string{"oauth"},
			"content":             "Welcome",
			"description":         "A test instance",
			"onboarding_status":   "requires_more_accounts",
			"registration_mode":   "public",
			"title":               "Testden",
			"web_address":         "https://web.test",
		}))
	}))
	defer server.Close()

	client, err := openapi.NewClientWithResponses(server.URL)
	r.NoError(err)

	info, err := fetchInfo(context.Background(), client)
	r.NoError(err)
	r.Equal("Testden", info.Title)
	r.Equal("A test instance", info.Description)
	r.Equal("https://web.test", info.WebAddress)
	r.Equal("https://api.test", info.ApiAddress)
}

func TestFetchInfoError(t *testing.T) {
	r := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		http.Error(w, "broken", http.StatusInternalServerError)
	}))
	defer server.Close()

	client, err := openapi.NewClientWithResponses(server.URL)
	r.NoError(err)

	_, err = fetchInfo(context.Background(), client)
	r.ErrorContains(err, "get instance info failed")
	r.True(strings.Contains(err.Error(), "broken") || strings.Contains(err.Error(), "500"))
}

func TestNormaliseColour(t *testing.T) {
	r := require.New(t)

	hex, ok := normaliseColour("#fc0")
	r.True(ok)
	r.Equal("#ffcc00", hex)

	hex, ok = normaliseColour("hsl(0, 100%, 50%)")
	r.True(ok)
	r.Equal("#FF0000", hex)
}
