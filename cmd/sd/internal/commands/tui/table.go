package tui

import (
	"fmt"
	"strconv"

	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"

	"github.com/Southclaws/storyden/cmd/sd/internal/render"
)

func (m model) buildTable() table.Model {
	width := m.tableWidth()
	height := max(m.height-11, 8)

	styles := table.DefaultStyles()
	styles.Header = styles.Header.Foreground(lipgloss.Color("205")).Bold(true)
	styles.Selected = styles.Selected.Foreground(lipgloss.Color("86")).Bold(true)

	return table.New(
		table.WithColumns(m.columns(width)),
		table.WithRows(m.rows()),
		table.WithFocused(true),
		table.WithHeight(height),
		table.WithWidth(width),
		table.WithStyles(styles),
	)
}

func (m model) columns(width int) []table.Column {
	switch m.active {
	case viewThreads:
		if width < 70 {
			return []table.Column{
				{Title: "Title", Width: max(width-20, 14)},
				{Title: "Updated", Width: 10},
				{Title: "Author", Width: 10},
			}
		}

		return []table.Column{
			{Title: "Title", Width: max(width-49, 24)},
			{Title: "Updated", Width: 16},
			{Title: "Author", Width: 16},
			{Title: "Replies", Width: 7},
		}
	default:
		if width < 70 {
			return []table.Column{
				{Title: "Name", Width: max(width-20, 14)},
				{Title: "Updated", Width: 10},
				{Title: "Author", Width: 10},
			}
		}

		return []table.Column{
			{Title: "Name", Width: max(width-50, 24)},
			{Title: "Updated", Width: 16},
			{Title: "Author", Width: 16},
			{Title: "Visibility", Width: 10},
		}
	}
}

func (m model) tableWidth() int {
	detailWidth := max(m.width/3, 36)
	return max(m.width-sidebarWidth-detailWidth-8, 34)
}

func (m model) rows() []table.Row {
	narrow := m.tableWidth() < 70

	switch m.active {
	case viewThreads:
		rows := make([]table.Row, 0, len(m.threads.Threads))
		for _, thread := range m.threads.Threads {
			if narrow {
				rows = append(rows, table.Row{
					thread.Title,
					render.FormatTime(thread.UpdatedAt),
					render.AuthorName(thread.Author),
				})

				continue
			}

			rows = append(rows, table.Row{
				thread.Title,
				render.FormatTime(thread.UpdatedAt),
				render.AuthorName(thread.Author),
				strconv.Itoa(thread.ReplyStatus.Replies),
			})
		}
		return rows
	default:
		nodes := m.currentNodes()
		rows := make([]table.Row, 0, len(nodes))
		for _, node := range nodes {
			name := node.Name
			if len(node.Children) > 0 {
				name += fmt.Sprintf(" (%d)", len(node.Children))
			}
			if narrow {
				rows = append(rows, table.Row{
					name,
					render.FormatTime(node.UpdatedAt),
					render.AuthorName(node.Owner),
				})

				continue
			}

			rows = append(rows, table.Row{
				name,
				render.FormatTime(node.UpdatedAt),
				render.AuthorName(node.Owner),
				string(node.Visibility),
			})
		}
		return rows
	}
}
