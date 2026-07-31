package info

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/spf13/cobra"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/api"
	"github.com/Southclaws/storyden/cmd/sd/internal/cligen"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/output"
	"github.com/Southclaws/storyden/cmd/sd/internal/render"
	"github.com/Southclaws/storyden/cmd/sd/internal/tui"
)

var hslPattern = regexp.MustCompile(`(?i)^hsla?\(\s*([0-9.]+)\s*,\s*([0-9.]+)%\s*,\s*([0-9.]+)%`)

func New(store *config.Store) cligen.InfoHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.InfoParams) (cligen.InstanceInfo, error) {
		contextName, configuredEndpoint, err := currentContext(store)
		if err != nil {
			return cligen.InstanceInfo{}, err
		}

		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return cligen.InstanceInfo{}, err
		}

		info, err := fetchInfo(ctx, client.OpenAPI)
		if err != nil {
			return cligen.InstanceInfo{}, err
		}

		result := cligen.InstanceInfo{
			Context:  &contextName,
			Endpoint: configuredEndpoint,
			BaseURL:  client.BaseURL,
			Info:     toInfo(*info),
		}
		if result.Endpoint == "" {
			result.Endpoint = client.Endpoint
		}

		if p.Format == cligen.InfoFormatPlain {
			if err := renderPlain(io.Out, result); err != nil {
				return cligen.InstanceInfo{}, err
			}
		}

		return result, nil
	}
}

func NewMetadata(store *config.Store) cligen.InfoMetadataHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.InfoMetadataParams) (cligen.Metadata, error) {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return nil, err
		}

		info, err := fetchInfo(ctx, client.OpenAPI)
		if err != nil {
			return nil, err
		}

		if info.Metadata == nil {
			return nil, nil
		}

		return cligen.Metadata(*info.Metadata), nil
	}
}

// toInfo maps the HTTP API's Info shape to the OpenCLI-schema-generated
// Info type: the two are the same data, but the generated type's enum-like
// fields are plain strings rather than the API client's named types.
func toInfo(info openapi.Info) cligen.Info {
	motd := (*cligen.InfoMotd)(nil)
	if info.Motd != nil {
		motd = &cligen.InfoMotd{Content: &info.Motd.Content}
	}

	capabilities := make([]string, len(info.Capabilities))
	for i, c := range info.Capabilities {
		capabilities[i] = string(c)
	}

	description := info.Description
	authMode := string(info.AuthenticationMode)
	registrationMode := string(info.RegistrationMode)
	accentColour := info.AccentColour
	content := info.Content

	return cligen.Info{
		Title:              info.Title,
		Description:        &description,
		Content:            &content,
		WebAddress:         info.WebAddress,
		APIAddress:         info.ApiAddress,
		AuthenticationMode: &authMode,
		RegistrationMode:   &registrationMode,
		AccentColour:       &accentColour,
		Capabilities:       capabilities,
		Motd:               motd,
	}
}

func currentContext(store *config.Store) (string, string, error) {
	cfg, err := store.Load()
	if err != nil {
		return "", "", err
	}

	if cfg.CurrentContext == "" {
		return "", "", nil
	}

	ctx, ok := cfg.Contexts[cfg.CurrentContext]
	if !ok {
		return cfg.CurrentContext, "", nil
	}

	return cfg.CurrentContext, ctx.APIURL, nil
}

func fetchInfo(ctx context.Context, client *openapi.ClientWithResponses) (*openapi.Info, error) {
	response, err := client.GetInfoWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, output.RequestError("get instance info", response, response.Body)
	}

	return (*openapi.Info)(response.JSON200), nil
}

func renderPlain(out io.Writer, result cligen.InstanceInfo) error {
	terminal := output.IsTerminal(out)
	styles := infoStyles(terminal)

	fmt.Fprintln(out, styles.Title.Render(result.Info.Title))
	if str(result.Info.Description) != "" {
		fmt.Fprintln(out, styles.Description.Render(str(result.Info.Description)))
	}
	fmt.Fprintln(out)

	fields := [][2]string{
		{"Context", str(result.Context)},
		{"Endpoint", result.Endpoint},
		{"Web address", result.Info.WebAddress},
		{"API address", result.Info.APIAddress},
		{"API client base", result.BaseURL},
		{"Authentication", str(result.Info.AuthenticationMode)},
		{"Registration", str(result.Info.RegistrationMode)},
		{"Accent colour", accentColour(str(result.Info.AccentColour), terminal)},
	}
	writeFields(out, fields, styles)

	if len(result.Info.Capabilities) > 0 {
		fmt.Fprintln(out)
		fmt.Fprintln(out, styles.Section.Render("Capabilities"))
		for _, capability := range result.Info.Capabilities {
			fmt.Fprintf(out, "  %s %s\n", styles.Bullet.Render("-"), styles.Value.Render(prettyCapability(capability)))
		}
	}

	if result.Info.Motd != nil {
		if motd, err := overviewMarkdown(str(result.Info.Motd.Content)); err == nil && motd != "" {
			fmt.Fprintln(out)
			fmt.Fprintln(out, styles.Section.Render("Message"))
			fmt.Fprintln(out, renderMarkdown(motd, out))
		}
	}

	if overview, err := overviewMarkdown(str(result.Info.Content)); err == nil && overview != "" {
		fmt.Fprintln(out)
		fmt.Fprintln(out, styles.Section.Render("Overview"))
		fmt.Fprintln(out, renderMarkdown(overview, out))
	}

	return nil
}

func str(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

type styles struct {
	Title       lipgloss.Style
	Description lipgloss.Style
	Section     lipgloss.Style
	Label       lipgloss.Style
	Value       lipgloss.Style
	Bullet      lipgloss.Style
}

func infoStyles(terminal bool) styles {
	if !terminal {
		return styles{}
	}

	return styles{
		Title:       tui.Title.Copy().Bold(true),
		Description: lipgloss.NewStyle().Foreground(lipgloss.Color("252")),
		Section:     lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")),
		Label:       lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
		Value:       lipgloss.NewStyle().Foreground(lipgloss.Color("252")),
		Bullet:      lipgloss.NewStyle().Foreground(lipgloss.Color("86")),
	}
}

func writeFields(out io.Writer, fields [][2]string, styles styles) {
	width := 0
	for _, field := range fields {
		if field[1] == "" {
			continue
		}
		width = max(width, len(field[0]))
	}

	for _, field := range fields {
		if field[1] == "" {
			continue
		}
		label := fmt.Sprintf("%-*s", width, field[0])
		fmt.Fprintf(out, "%s  %s\n", styles.Label.Render(label), styles.Value.Render(field[1]))
	}
}

func prettyCapability(value string) string {
	text := strings.ReplaceAll(string(value), "_", " ")
	if text == "" {
		return ""
	}
	return strings.ToUpper(text[:1]) + text[1:]
}

func overviewMarkdown(value string) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "", nil
	}
	return render.HTMLToMarkdown(value)
}

func renderMarkdown(markdown string, out io.Writer) string {
	width := output.TerminalWidth(out, 88)
	if width > 8 {
		width -= 4
	}
	return strings.TrimSpace(output.Markdown(markdown, out, width))
}

func accentColour(value string, terminal bool) string {
	if value == "" {
		return ""
	}

	hex, ok := normaliseColour(value)
	if !ok || !terminal {
		return value
	}

	swatch := lipgloss.NewStyle().
		Background(lipgloss.Color(hex)).
		Foreground(lipgloss.Color(contrastColour(hex))).
		Render("  ")

	return fmt.Sprintf("%s %s", swatch, value)
}

func normaliseColour(value string) (string, bool) {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "#") {
		switch len(value) {
		case 4:
			return "#" + strings.Repeat(value[1:2], 2) + strings.Repeat(value[2:3], 2) + strings.Repeat(value[3:4], 2), true
		case 7:
			return strings.ToUpper(value), true
		}
	}

	matches := hslPattern.FindStringSubmatch(value)
	if len(matches) != 4 {
		return "", false
	}

	h, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return "", false
	}
	s, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return "", false
	}
	l, err := strconv.ParseFloat(matches[3], 64)
	if err != nil {
		return "", false
	}

	r, g, b := hslToRGB(h, s/100, l/100)
	return fmt.Sprintf("#%02X%02X%02X", r, g, b), true
}

func hslToRGB(h, s, l float64) (int, int, int) {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}
	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - c/2

	var rp, gp, bp float64
	switch {
	case h < 60:
		rp, gp, bp = c, x, 0
	case h < 120:
		rp, gp, bp = x, c, 0
	case h < 180:
		rp, gp, bp = 0, c, x
	case h < 240:
		rp, gp, bp = 0, x, c
	case h < 300:
		rp, gp, bp = x, 0, c
	default:
		rp, gp, bp = c, 0, x
	}

	return clampRGB((rp + m) * 255), clampRGB((gp + m) * 255), clampRGB((bp + m) * 255)
}

func clampRGB(value float64) int {
	return min(255, max(0, int(math.Round(value))))
}

func contrastColour(hex string) string {
	if len(hex) != 7 {
		return "#000000"
	}
	r, _ := strconv.ParseInt(hex[1:3], 16, 64)
	g, _ := strconv.ParseInt(hex[3:5], 16, 64)
	b, _ := strconv.ParseInt(hex[5:7], 16, 64)
	if (float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114) > 150 {
		return "#000000"
	}
	return "#FFFFFF"
}
