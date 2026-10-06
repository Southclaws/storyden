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

const (
	formatPlain = "plain"
	formatJSON  = "json"
)

type instanceInfo struct {
	Context  string       `json:"context,omitempty"`
	Auth     authInfo     `json:"auth"`
	Endpoint string       `json:"endpoint"`
	BaseURL  string       `json:"base_url"`
	Info     openapi.Info `json:"info"`
}

type authInfo struct {
	Status    string            `json:"status"`
	Method    config.AuthMethod `json:"method,omitempty"`
	AccountID string            `json:"account_id,omitempty"`
	Handle    string            `json:"handle,omitempty"`
}

var hslPattern = regexp.MustCompile(`(?i)^hsla?\(\s*([0-9.]+)\s*,\s*([0-9.]+)%\s*,\s*([0-9.]+)%`)

func New(store *config.Store) cligen.InfoHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.InfoParams) error {
		contextName, current, err := store.Current()
		if err != nil {
			return err
		}

		var client *api.Client
		if current.Auth == nil {
			client, err = api.NewStaticClient(current.APIURL)
		} else {
			client, err = api.NewAuthenticatedClient(ctx, store)
		}
		if err != nil {
			return err
		}

		session, err := fetchSession(ctx, client.OpenAPI)
		if err != nil {
			return err
		}

		result := instanceInfo{
			Context:  contextName,
			Auth:     authInfo{Status: "unauthenticated"},
			Endpoint: current.APIURL,
			BaseURL:  client.BaseURL,
			Info:     session.Info,
		}
		if current.Auth != nil {
			result.Auth.Method = current.Auth.MethodOrDefault()
		}
		if session.Account != nil {
			result.Auth.Status = "authenticated"
			result.Auth.AccountID = string(session.Account.Id)
			result.Auth.Handle = session.Account.Handle
		}
		if result.Endpoint == "" {
			result.Endpoint = client.Endpoint
		}

		return renderOutput(io.Out, string(p.Output), result)
	}
}

func NewMetadata(store *config.Store) cligen.InfoMetadataHandler {
	return func(ctx context.Context, cmd *cobra.Command, io cligen.IO, p cligen.InfoMetadataParams) error {
		client, err := api.NewAuthenticatedClient(ctx, store)
		if err != nil {
			return err
		}

		info, err := fetchInfo(ctx, client.OpenAPI)
		if err != nil {
			return err
		}

		if info.Metadata == nil {
			return output.JSON(io.Out, openapi.Metadata{})
		}

		return output.JSON(io.Out, *info.Metadata)
	}
}

func fetchSession(ctx context.Context, client *openapi.ClientWithResponses) (*openapi.SessionInfo, error) {
	response, err := client.GetSessionWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, fmt.Errorf("get session info failed: %s: %s", response.Status(), string(response.Body))
	}
	return response.JSON200, nil
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

func renderOutput(out io.Writer, format string, result instanceInfo) error {
	switch format {
	case formatPlain:
		return renderPlain(out, result)
	case formatJSON:
		return output.JSON(out, result)
	default:
		return fmt.Errorf("unsupported format %q", format)
	}
}

func renderPlain(out io.Writer, result instanceInfo) error {
	terminal := output.IsTerminal(out)
	styles := infoStyles(terminal)

	fmt.Fprintln(out, styles.Title.Render(result.Info.Title))
	if result.Info.Description != "" {
		fmt.Fprintln(out, styles.Description.Render(result.Info.Description))
	}
	fmt.Fprintln(out)

	fields := [][2]string{
		{"Context", result.Context},
		{"Auth status", result.Auth.Status},
		{"Auth method", string(result.Auth.Method)},
		{"Signed in as", result.Auth.Handle},
		{"Endpoint", result.Endpoint},
		{"Web address", result.Info.WebAddress},
		{"API address", result.Info.ApiAddress},
		{"API client base", result.BaseURL},
		{"Instance auth mode", string(result.Info.AuthenticationMode)},
		{"Registration", string(result.Info.RegistrationMode)},
		{"Accent colour", accentColour(result.Info.AccentColour, terminal)},
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
		if motd, err := overviewMarkdown(result.Info.Motd.Content); err == nil && motd != "" {
			fmt.Fprintln(out)
			fmt.Fprintln(out, styles.Section.Render("Message"))
			fmt.Fprintln(out, renderMarkdown(motd, out))
		}
	}

	if overview, err := overviewMarkdown(result.Info.Content); err == nil && overview != "" {
		fmt.Fprintln(out)
		fmt.Fprintln(out, styles.Section.Render("Overview"))
		fmt.Fprintln(out, renderMarkdown(overview, out))
	}

	return nil
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

func prettyCapability(value openapi.InstanceCapability) string {
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

func validateFormat(format string) error {
	switch format {
	case formatPlain, formatJSON:
		return nil
	default:
		return fmt.Errorf("--output must be one of: plain, json")
	}
}
