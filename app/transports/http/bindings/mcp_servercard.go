package bindings

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/internal/config"
)

type mcpServerCardRemote struct {
	Type                      string   `json:"type"`
	URL                       string   `json:"url"`
	SupportedProtocolVersions []string `json:"supportedProtocolVersions,omitempty"`
}

type mcpServerCardRepository struct {
	URL    string `json:"url"`
	Source string `json:"source"`
}

type mcpServerCardResponse struct {
	Schema      string                   `json:"$schema"`
	Name        string                   `json:"name"`
	Version     string                   `json:"version"`
	Description string                   `json:"description"`
	Title       string                   `json:"title,omitempty"`
	WebsiteURL  string                   `json:"websiteUrl,omitempty"`
	Repository  *mcpServerCardRepository `json:"repository,omitempty"`
	Remotes     []mcpServerCardRemote    `json:"remotes,omitempty"`
}

func mountMCPServerCard(
	router *echo.Echo,
	cfg config.Config,
	sr *settings.SettingsRepository,
) {
	if !cfg.MCPEnabled {
		return
	}

	router.GET("/.well-known/mcp/server-card.json", func(c echo.Context) error {
		set, err := sr.Get(c.Request().Context())
		if err != nil {
			return err
		}

		title := set.Title.Or("Storyden")
		description := set.Description.Or("A modern community platform combining forum, wiki, and community hub features.")

		card := mcpServerCardResponse{
			Schema:      "https://static.modelcontextprotocol.io/schemas/v1/server-card.schema.json",
			Name:        "org.storyden/storyden",
			Version:     config.Version,
			Description: description,
			Title:       title,
			WebsiteURL:  cfg.PublicWebAddress.String(),
			Repository: &mcpServerCardRepository{
				URL:    "https://github.com/Southclaws/storyden",
				Source: "github",
			},
			Remotes: []mcpServerCardRemote{
				{
					Type:                      "sse",
					URL:                       cfg.PublicAPIAddress.String() + "/mcp/sse",
					SupportedProtocolVersions: []string{"2025-03-26"},
				},
			},
		}

		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set("Access-Control-Allow-Methods", "GET")
		c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Response().Header().Set("Cache-Control", "public, max-age=3600")

		return c.JSON(http.StatusOK, card)
	})
}
