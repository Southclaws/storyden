package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

const discoveryTimeout = 15 * time.Second

type Client struct {
	Endpoint  string
	BaseURL   string
	OpenAPI   *openapi.ClientWithResponses
	Discovery *Discovery
}

type Discovery struct {
	Issuer                           string   `json:"issuer"`
	AuthorizationEndpoint            string   `json:"authorization_endpoint"`
	DeviceAuthorizationEndpoint      string   `json:"device_authorization_endpoint"`
	TokenEndpoint                    string   `json:"token_endpoint"`
	UserinfoEndpoint                 string   `json:"userinfo_endpoint"`
	JWKSURI                          string   `json:"jwks_uri"`
	ResponseTypesSupported           []string `json:"response_types_supported"`
	GrantTypesSupported              []string `json:"grant_types_supported"`
	CodeChallengeMethodsSupported    []string `json:"code_challenge_methods_supported"`
	ScopesSupported                  []string `json:"scopes_supported"`
	SubjectTypesSupported            []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`
}

func NewClient(ctx context.Context, rawEndpoint string) (*Client, error) {
	endpoint, explicitBase, err := parseEndpoint(rawEndpoint)
	if err != nil {
		return nil, err
	}

	discovery, err := discoverOAuth(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	var lastErr error
	for _, candidate := range candidateBaseURLs(endpoint, explicitBase) {
		client, err := openapi.NewClientWithResponses(candidate)
		if err != nil {
			lastErr = err
			continue
		}

		if discovery.TokenEndpoint != "" && strings.HasPrefix(discovery.TokenEndpoint, strings.TrimRight(candidate, "/")+"/") {
			return &Client{
				Endpoint:  endpoint,
				BaseURL:   candidate,
				OpenAPI:   client,
				Discovery: discovery,
			}, nil
		}

		lastErr = fmt.Errorf("OAuth discovery endpoint does not match API base %s", candidate)
	}

	if lastErr != nil {
		return nil, lastErr
	}

	return nil, fmt.Errorf("OAuth discovery failed")
}

func NewStaticClient(rawEndpoint string) (*Client, error) {
	endpoint, explicitBase, err := parseEndpoint(rawEndpoint)
	if err != nil {
		return nil, err
	}

	baseURL := strings.TrimRight(endpoint, "/") + "/api"
	if explicitBase != "" {
		baseURL = explicitBase
	}

	client, err := openapi.NewClientWithResponses(baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		Endpoint: endpoint,
		BaseURL:  baseURL,
		OpenAPI:  client,
	}, nil
}

func discoverOAuth(ctx context.Context, endpoint string) (*Discovery, error) {
	u, err := url.Parse(strings.TrimRight(endpoint, "/") + "/.well-known/openid-configuration")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{Timeout: discoveryTimeout}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, rateLimitExceededFromHeaders(resp.Header, time.Now())
	}

	if resp.StatusCode != http.StatusOK {
		var apiError openapi.APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiError); err == nil && apiError.Title != nil {
			message := *apiError.Title
			if apiError.Metadata != nil {
				if suggested, ok := (*apiError.Metadata)["suggested"].(string); ok && suggested != "" {
					message += "\n\n" + suggested
				}
			}

			return nil, fmt.Errorf("%s", message)
		}

		return nil, fmt.Errorf("OAuth discovery failed: %s", resp.Status)
	}

	var discovery Discovery
	if err := json.NewDecoder(resp.Body).Decode(&discovery); err != nil {
		return nil, err
	}

	return &discovery, nil
}

func CanonicalEndpoint(rawEndpoint string) (string, error) {
	endpoint, _, err := parseEndpoint(rawEndpoint)
	return endpoint, err
}

func parseEndpoint(raw string) (string, string, error) {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", "", err
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", "", fmt.Errorf("Storyden API URL must use http or https")
	}

	if parsed.Host == "" {
		return "", "", fmt.Errorf("Storyden API URL must include a host")
	}

	parsed.RawQuery = ""
	parsed.Fragment = ""

	explicitPath := strings.TrimRight(parsed.EscapedPath(), "/")
	parsed.Path = ""
	parsed.RawPath = ""

	endpoint := parsed.String()
	explicitBase := ""
	if explicitPath != "" {
		explicit := *parsed
		explicit.Path = explicitPath
		explicitBase = explicit.String()
	}

	return endpoint, explicitBase, nil
}

func candidateBaseURLs(endpoint string, explicitBase string) []string {
	candidates := []string{}
	if explicitBase != "" {
		candidates = append(candidates, explicitBase)
	}

	candidates = append(candidates, strings.TrimRight(endpoint, "/")+"/api")
	candidates = append(candidates, endpoint)

	return unique(candidates)
}

func unique(values []string) []string {
	seen := map[string]struct{}{}
	result := []string{}

	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}

		seen[value] = struct{}{}
		result = append(result, value)
	}

	return result
}
