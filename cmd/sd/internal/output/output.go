package output

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"charm.land/glamour/v2"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

type fdWriter interface {
	Fd() uintptr
}

// IsTerminal checks whether the provided writer is backed by a terminal.
func IsTerminal(w io.Writer) bool {
	f, ok := w.(fdWriter)
	if !ok {
		return false
	}

	return term.IsTerminal(int(f.Fd()))
}

// TerminalWidth returns the terminal width, defaulting when detection fails.
func TerminalWidth(w io.Writer, fallback int) int {
	f, ok := w.(fdWriter)
	if !ok {
		return fallback
	}

	width, _, err := term.GetSize(int(f.Fd()))
	if err != nil || width <= 0 {
		return fallback
	}

	return width
}

// JSON writes stable, indented JSON for humans and agents.
func JSON(out io.Writer, value interface{}) error {
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")

	return encoder.Encode(value)
}

// YAML writes structured YAML using the real YAML encoder.
func YAML(out io.Writer, value interface{}) error {
	encoder := yaml.NewEncoder(out)
	defer encoder.Close()

	return encoder.Encode(value)
}

// Markdown renders markdown with Glamour for terminals and returns the raw
// markdown for pipes/files.
func Markdown(markdown string, out io.Writer, width int) string {
	if !IsTerminal(out) {
		return markdown
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithEnvironmentConfig(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return markdown
	}

	rendered, err := renderer.Render(markdown)
	if err != nil {
		return markdown
	}

	return rendered
}

type StatusResponse interface {
	Status() string
	StatusCode() int
}

func RequestError(operation string, response StatusResponse, body []byte) error {
	return RequestErrorWithMessages(operation, response, body, nil)
}

func RequestErrorWithMessages(operation string, response StatusResponse, body []byte, messages map[int]string) error {
	if message, ok := messages[response.StatusCode()]; ok {
		return fmt.Errorf("%s", message)
	}

	text := strings.TrimSpace(string(body))
	if text != "" {
		return fmt.Errorf("%s failed: %s: %s", operation, response.Status(), text)
	}

	return fmt.Errorf("%s failed: %s", operation, response.Status())
}

func UnauthorizedMessage(operation string) map[int]string {
	return map[int]string{
		http.StatusUnauthorized: fmt.Sprintf("%s was not authorised; run sd auth login again", operation),
	}
}
