package logger

import "github.com/charmbracelet/lipgloss"

var (
	// Color palette
	purple  = lipgloss.Color("#7D56F4")
	pink    = lipgloss.Color("#FF6AC1")
	blue    = lipgloss.Color("#00D9FF")
	green   = lipgloss.Color("#32CD32")
	yellow  = lipgloss.Color("#FFD700")
	orange  = lipgloss.Color("#FF8C00")
	gray    = lipgloss.Color("#888888")
	white   = lipgloss.Color("#FAFAFA")
	crimson = lipgloss.Color("#DC143C")

	// Header style for phase announcements
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white).
			Background(purple).
			Padding(0, 1).
			MarginTop(1).
			MarginBottom(0)

	// Success style for completion messages
	SuccessStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white).
			Background(green).
			Padding(0, 1)

	// Error style
	ErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white).
			Background(crimson).
			Padding(0, 1)

	// Resource type label styles
	AccountLabel = lipgloss.NewStyle().
			Bold(true).
			Foreground(pink).
			Width(12).
			Align(lipgloss.Right)

	CategoryLabel = lipgloss.NewStyle().
			Bold(true).
			Foreground(blue).
			Width(12).
			Align(lipgloss.Right)

	PostLabel = lipgloss.NewStyle().
			Bold(true).
			Foreground(orange).
			Width(12).
			Align(lipgloss.Right)

	RoleLabel = lipgloss.NewStyle().
			Bold(true).
			Foreground(purple).
			Width(12).
			Align(lipgloss.Right)

	TagLabel = lipgloss.NewStyle().
			Bold(true).
			Foreground(yellow).
			Width(12).
			Align(lipgloss.Right)

	// Field styles
	FieldKey = lipgloss.NewStyle().
			Foreground(gray).
			Italic(true)

	FieldValue = lipgloss.NewStyle().
			Foreground(white).
			Bold(false)

	// Arrow connector
	Arrow = lipgloss.NewStyle().
		Foreground(gray).
		Render("→")

	// Dimmed style for less important info
	Dim = lipgloss.NewStyle().
		Foreground(gray)
)
