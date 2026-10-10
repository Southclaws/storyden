package mcp

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/Southclaws/opt"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/fx"

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
	Identifier string `json:"identifier"`
	Type       string `json:"type"`
	URL        string `json:"url"`
}

type apiCatalog struct {
	Linkset []apiCatalogLinks `json:"linkset"`
}

type apiCatalogLinks struct {
	Anchor      string           `json:"anchor"`
	Item        []apiCatalogLink `json:"item"`
	ServiceDesc []apiCatalogLink `json:"service-desc,omitempty"`
}

type apiCatalogLink struct {
	Href string `json:"href"`
	Type string `json:"type,omitempty"`
}

func serverName(webHost string) string {
	if webHost == "" || net.ParseIP(webHost) != nil {
		return "org.storyden/mcp"
	}

	labels := strings.Split(strings.ToLower(webHost), ".")
	for i, label := range labels {
		if label == "" || strings.Trim(label, "abcdefghijklmnopqrstuvwxyz0123456789-") != "" {
			return "org.storyden/mcp"
		}
		labels[i] = label
	}
	for i, j := 0, len(labels)-1; i < j; i, j = i+1, j-1 {
		labels[i], labels[j] = labels[j], labels[i]
	}
	name := strings.Join(labels, ".") + "/mcp"
	if len(name) > 200 {
		return "org.storyden/mcp"
	}
	return name
}

func discoveryDocuments(cfg config.Config, title string) (serverCard, aiCatalog, apiCatalog) {
	apiURL := cfg.PublicAPIAddress.JoinPath("/api")
	apiSchemaURL := cfg.PublicAPIAddress.JoinPath("/api/openapi.json")
	apiCatalogURL := cfg.PublicWebAddress.JoinPath("/.well-known/api-catalog")

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
			Identifier: "urn:air:" + cfg.PublicWebAddress.Hostname() + ":mcp:storyden",
			Type:       serverCardMediaType, URL: cardURL.String(),
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
	etag := "\"" + hex.EncodeToString(hash[:]) + "\""

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mediaType)
		w.Header().Set("Cache-Control", "public, max-age=3600")
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

func MountDiscovery(lc fx.Lifecycle, ctx context.Context, cfg config.Config, settings *settings.SettingsRepository, mux *http.ServeMux) {
	lc.Append(fx.StartHook(func() error {
		title := "Storyden"
		if cfg.MCPEnabled {
			set, err := settings.Get(ctx)
			if err != nil {
				return err
			}
			title = set.Title.Or(title)
		}

		mountDiscovery(mux, cfg, title)
		return nil
	}))
}

func mountDiscovery(mux *http.ServeMux, cfg config.Config, title string) {
	card, ai, api := discoveryDocuments(cfg, title)
	if cfg.MCPEnabled {
		cardHandler := discoveryHandler(card, serverCardMediaType, true)
		mux.Handle("GET /mcp/server-card", cardHandler)
		mux.Handle("OPTIONS /mcp/server-card", cardHandler)
	}
	aiHandler := discoveryHandler(ai, "application/ai-catalog+json", true)
	mux.Handle("GET /.well-known/ai-catalog.json", aiHandler)
	mux.Handle("OPTIONS /.well-known/ai-catalog.json", aiHandler)
	apiHandler := discoveryHandler(api, apiCatalogMediaType, true)
	linkedAPIHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Link", "<"+cfg.PublicWebAddress.JoinPath("/.well-known/api-catalog").String()+">; rel=api-catalog")
		apiHandler.ServeHTTP(w, r)
	})
	mux.Handle("GET /.well-known/api-catalog", linkedAPIHandler)
	mux.Handle("OPTIONS /.well-known/api-catalog", linkedAPIHandler)
}
