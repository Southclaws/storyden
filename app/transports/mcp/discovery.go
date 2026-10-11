package mcp

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/Southclaws/opt"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/fx"
	"golang.org/x/net/idna"

	"github.com/Southclaws/storyden/app/resources/cachecontrol"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/internal/config"
)

const (
	serverCardMediaType = "application/mcp-server-card+json"
	serverCardSchema    = "https://static.modelcontextprotocol.io/schemas/v1/server-card.schema.json"
	apiCatalogMediaType = "application/linkset+json; profile=\"https://www.rfc-editor.org/info/rfc9727\""
)

type serverCard struct {
	Schema      string       `json:"$schema"`
	Name        string       `json:"name"`
	Version     string       `json:"version"`
	Description string       `json:"description"`
	Title       string       `json:"title,omitempty"`
	WebsiteURL  string       `json:"websiteUrl,omitempty"`
	Remotes     []cardRemote `json:"remotes"`
}

type cardRemote struct {
	Type                      string   `json:"type"`
	URL                       string   `json:"url"`
	SupportedProtocolVersions []string `json:"supportedProtocolVersions"`
}

type aiCatalog struct {
	SpecVersion string           `json:"specVersion"`
	Entries     []aiCatalogEntry `json:"entries"`
}

type aiCatalogEntry struct {
	Identifier  string `json:"identifier"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Type        string `json:"type"`
	URL         string `json:"url"`
}

type apiCatalog struct {
	Linkset []apiCatalogLinks `json:"linkset"`
}

type apiCatalogLinks struct {
	Anchor      string           `json:"anchor"`
	Item        []apiCatalogLink `json:"item,omitempty"`
	ServiceDesc []apiCatalogLink `json:"service-desc,omitempty"`
}

type apiCatalogLink struct {
	Href string `json:"href"`
	Type string `json:"type,omitempty"`
}

func serverName(webHost string) string {
	host, err := idna.Lookup.ToASCII(webHost)
	if err != nil {
		return "localhost/mcp"
	}
	// Remove the optional DNS root dot before applying registration constraints.
	host, err = idna.Registration.ToASCII(strings.TrimSuffix(host, "."))
	if err != nil || host == "" || net.ParseIP(host) != nil {
		return "localhost/mcp"
	}

	labels := strings.Split(host, ".")
	for i, j := 0, len(labels)-1; i < j; i, j = i+1, j-1 {
		labels[i], labels[j] = labels[j], labels[i]
	}
	name := strings.Join(labels, ".") + "/mcp"
	if len(name) > 200 {
		return "localhost/mcp"
	}
	return name
}

func discoveryDocuments(cfg config.Config, title string) (serverCard, aiCatalog, apiCatalog) {
	apiURL := cfg.PublicAPIAddress.JoinPath("/api")
	apiSchemaURL := cfg.PublicAPIAddress.JoinPath("/api/openapi.json")
	apiCatalogURL := cfg.PublicAPIAddress.JoinPath("/.well-known/api-catalog")

	api := apiCatalog{Linkset: []apiCatalogLinks{
		{Anchor: apiCatalogURL.String(), Item: []apiCatalogLink{{Href: apiURL.String()}}},
		{Anchor: apiURL.String(), ServiceDesc: []apiCatalogLink{{Href: apiSchemaURL.String(), Type: "application/vnd.oai.openapi+json"}}},
	}}
	card := serverCard{}
	ai := aiCatalog{SpecVersion: "1.0", Entries: []aiCatalogEntry{}}
	if cfg.MCPEnabled {
		mcpURL := cfg.PublicAPIAddress.JoinPath("/mcp")
		cardURL := cfg.PublicAPIAddress.JoinPath("/mcp/server-card")
		card = serverCard{
			Schema:      serverCardSchema,
			Name:        serverName(cfg.PublicWebAddress.Hostname()),
			Version:     config.Version,
			Description: "Access this Storyden community through MCP.",
			WebsiteURL:  cfg.PublicWebAddress.String(),
			Remotes: []cardRemote{{
				Type: "streamable-http", URL: mcpURL.String(),
				SupportedProtocolVersions: sdkmcp.SupportedProtocolVersions(),
			}},
		}
		if utf8.RuneCountInString(title) <= 100 {
			card.Title = title
		}
		ai.Entries = append(ai.Entries, aiCatalogEntry{
			Identifier:  fmt.Sprintf("urn:air:%s:mcp:storyden", cfg.PublicWebAddress.Hostname()),
			DisplayName: title,
			Description: fmt.Sprintf("Explore discussions and Library pages in %s through MCP.", title),
			Type:        serverCardMediaType, URL: cardURL.String(),
		})
		api.Linkset[0].Item = append(api.Linkset[0].Item, apiCatalogLink{Href: mcpURL.String()})
	}
	return card, ai, api
}

func discoveryHandler(document any, mediaType string, cors bool) http.Handler {
	body, err := json.Marshal(document)
	if err != nil {
		panic(err)
	}
	body = append(body, '\n')
	hash := sha256.Sum256(body)
	etag := fmt.Sprintf("\"%x\"", hash)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mediaType)
		w.Header().Set("Cache-Control", "no-cache, must-revalidate")
		w.Header().Set("ETag", etag)
		if cors {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, If-None-Match")
			w.Header().Set("Access-Control-Expose-Headers", "ETag")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if cachecontrol.NewQuery(opt.New(r.Header.Get("If-None-Match"))).MatchesETag(etag) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		if r.Method == http.MethodHead {
			return
		}
		_, _ = w.Write(body)
	})
}

func MountDiscovery(lc fx.Lifecycle, cfg config.Config, repo *settings.SettingsRepository, mux *http.ServeMux) {
	lc.Append(fx.StartHook(func() {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store")

			current, err := repo.Get(r.Context())
			if err != nil {
				http.Error(w, "Discovery settings unavailable", http.StatusServiceUnavailable)
				return
			}

			policy := current.Services.OrZero().Discovery.OrZero()
			card, ai, api := discoveryDocuments(cfg, current.Title.Or("Storyden"))
			if !policy.MCPServerCard.Or(false) {
				ai.Entries = []aiCatalogEntry{}
			}

			var enabled bool
			var document any
			var mediaType string

			switch r.URL.Path {
			case "/mcp/server-card":
				enabled = policy.MCPServerCard.Or(false)
				document, mediaType = card, serverCardMediaType

			case "/.well-known/ai-catalog.json":
				enabled = policy.AICatalog.Or(false)
				document, mediaType = ai, "application/ai-catalog+json"

			case "/.well-known/api-catalog":
				enabled = policy.APICatalog.Or(false)
				document, mediaType = api, apiCatalogMediaType
			}

			if !enabled {
				http.NotFound(w, r)
				return
			}

			if r.URL.Path == "/.well-known/api-catalog" {
				w.Header().Set("Link", fmt.Sprintf("<%s>; rel=api-catalog", cfg.PublicAPIAddress.JoinPath("/.well-known/api-catalog")))
			}

			discoveryHandler(document, mediaType, true).ServeHTTP(w, r)
		})

		paths := []string{"/.well-known/ai-catalog.json", "/.well-known/api-catalog"}
		if cfg.MCPEnabled {
			paths = append(paths, "/mcp/server-card")
		}

		for _, path := range paths {
			mux.Handle("GET "+path, handler)
			mux.Handle("OPTIONS "+path, handler)
		}
	}))
}
