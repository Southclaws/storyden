package bindings

import (
	"strings"

	"github.com/Southclaws/storyden/internal/config"
)

// This file implements the integrations.sh discovery document, a static
// description of this instance's remote-accessible surfaces (REST API, MCP
// server, etc.) and how an agent authenticates against them.
//
// See: https://integrations.sh

type integrationsDocument struct {
	Version     int                               `json:"version"`
	Summary     string                            `json:"summary,omitempty"`
	Credentials map[string]integrationsCredential `json:"credentials,omitempty"`
	Surfaces    []integrationsSurface             `json:"surfaces,omitempty"`
}

type integrationsCredential struct {
	Type        string `json:"type"`
	Label       string `json:"label"`
	GenerateURL string `json:"generateUrl,omitempty"`
	Setup       string `json:"setup"`
}

type integrationsBasis struct {
	Via    string `json:"via"`
	Source string `json:"source"`
}

type integrationsMechanics struct {
	Source     string `json:"source"`
	In         string `json:"in,omitempty"`
	HeaderName string `json:"headerName,omitempty"`
	Scheme     string `json:"scheme,omitempty"`
	ParamName  string `json:"paramName,omitempty"`
}

type integrationsAuthUse struct {
	ID        string                `json:"id"`
	Mechanics integrationsMechanics `json:"mechanics"`
}

type integrationsAuthEntry struct {
	Use   []integrationsAuthUse `json:"use"`
	Basis integrationsBasis     `json:"basis"`
}

type integrationsAuthStatus struct {
	Status  string                  `json:"status"`
	Basis   *integrationsBasis      `json:"basis,omitempty"`
	Entries []integrationsAuthEntry `json:"entries,omitempty"`
}

type integrationsSurface struct {
	Slug       string                 `json:"slug"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Docs       string                 `json:"docs,omitempty"`
	Basis      integrationsBasis      `json:"basis"`
	Auth       integrationsAuthStatus `json:"auth"`
	Spec       string                 `json:"spec,omitempty"`
	URL        string                 `json:"url,omitempty"`
	Transports []string               `json:"transports,omitempty"`
}

const (
	accessKeyCredentialID = "storyden_access_key"
	oauthCredentialID     = "storyden_oauth"
)

func buildIntegrationsDocument(cfg config.Config, oauthBinding OAuth) integrationsDocument {
	selfURL := trimSlash(cfg.PublicWebAddress.String()) + "/.well-known/integrations.json"
	basis := integrationsBasis{Via: "declared", Source: selfURL}
	apiBase := trimSlash(cfg.PublicAPIAddress.String())
	webBase := trimSlash(cfg.PublicWebAddress.String())

	credentials := map[string]integrationsCredential{
		accessKeyCredentialID: {
			Type:        "bearer",
			Label:       "Storyden access key",
			GenerateURL: webBase + "/settings",
			Setup: "Sign in to your Storyden account, open Settings -> Access Keys and create a " +
				"new key (requires the `USE_PERSONAL_ACCESS_KEYS` permission, granted automatically " +
				"to administrators). Copy the secret immediately, it is only ever shown once.",
		},
	}

	accessKeyEntry := integrationsAuthEntry{
		Use: []integrationsAuthUse{{
			ID: accessKeyCredentialID,
			Mechanics: integrationsMechanics{
				Source:     "http",
				In:         "header",
				HeaderName: "Authorization",
				Scheme:     "Bearer",
			},
		}},
		Basis: basis,
	}

	apiAuthEntries := []integrationsAuthEntry{accessKeyEntry}

	if oauthBinding.oauth.Enabled() {
		credentials[oauthCredentialID] = integrationsCredential{
			Type:        "oauth2",
			Label:       "Storyden OAuth 2.0",
			GenerateURL: trimSlash(oauthBinding.Issuer()) + "/.well-known/oauth-authorization-server",
			Setup: "Storyden runs a standard OAuth 2.0 / OpenID Connect authorisation server. " +
				"Register a client (or use Dynamic Client Registration / a Client ID Metadata " +
				"Document) and run the authorization_code or device_code grant to obtain a " +
				"member-scoped bearer token. See /docs/introduction/oauth for details.",
		}

		apiAuthEntries = append(apiAuthEntries, integrationsAuthEntry{
			Use: []integrationsAuthUse{{
				ID: oauthCredentialID,
				Mechanics: integrationsMechanics{
					Source: "spec",
					Scheme: "oauth_token",
				},
			}},
			Basis: basis,
		})
	}

	surfaces := []integrationsSurface{
		{
			Slug:  "storyden-api",
			Name:  "Storyden REST API",
			Type:  "http",
			Docs:  apiBase + "/api/docs",
			Basis: basis,
			Auth: integrationsAuthStatus{
				Status:  "required",
				Entries: apiAuthEntries,
			},
			Spec: apiBase + "/api/openapi.json",
			URL:  apiBase + "/api",
		},
	}

	if cfg.MCPEnabled {
		surfaces = append(surfaces, integrationsSurface{
			Slug:  "storyden-mcp",
			Name:  "Storyden MCP Server",
			Type:  "mcp",
			Docs:  "https://storyden.org/docs/introduction/mcp",
			Basis: basis,
			Auth: integrationsAuthStatus{
				Status:  "required",
				Entries: []integrationsAuthEntry{accessKeyEntry},
			},
			URL:        apiBase + "/mcp",
			Transports: []string{"streamable-http"},
		})
	}

	return integrationsDocument{
		Version: 3,
		Summary: "Storyden is a self-hosted forum, wiki and community hub. This instance exposes " +
			"a REST API and, when enabled, an MCP server for agent-driven content creation and search.",
		Credentials: credentials,
		Surfaces:    surfaces,
	}
}

func trimSlash(s string) string {
	return strings.TrimSuffix(s, "/")
}
