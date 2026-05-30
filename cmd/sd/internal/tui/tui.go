package tui

import (
	"io"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

var (
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	Accent = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86"))

	Muted = lipgloss.NewStyle().
		Foreground(lipgloss.Color("245"))

	URL = lipgloss.NewStyle().
		Underline(true).
		Foreground(lipgloss.Color("81"))
)

func NewForm(in io.Reader, out io.Writer, groups ...*huh.Group) *huh.Form {
	return huh.NewForm(groups...).
		WithInput(in).
		WithOutput(out).
		WithTheme(huh.ThemeFunc(huh.ThemeCatppuccin))
}
