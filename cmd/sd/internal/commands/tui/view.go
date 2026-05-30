package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/Southclaws/storyden/cmd/sd/internal/render"
	sharedtui "github.com/Southclaws/storyden/cmd/sd/internal/tui"
)

const sidebarWidth = 18

func (m model) View() tea.View {
	if m.screen != screenList {
		return m.documentView()
	}

	detailWidth := max(m.width/3, 36)
	tableWidth := max(m.width-sidebarWidth-detailWidth-6, 40)
	bodyHeight := max(m.height-7, 10)

	sidebar := m.renderSidebar(sidebarWidth, bodyHeight)
	content := m.renderContent(tableWidth, bodyHeight)
	details := m.renderDetails(detailWidth, bodyHeight)

	header := sharedtui.Title.Render("Storyden")
	if m.loading {
		header += " " + sharedtui.Muted.Render("loading...")
	}
	if m.err != nil {
		header += " " + lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Render(m.err.Error())
	}

	footer := sharedtui.Muted.Render("n/t resources  ↑/↓ move  enter open  c children  backspace up  r refresh  q quit")

	view := tea.NewView(
		header + "\n\n" +
			lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content, details) +
			"\n\n" + footer,
	)
	view.AltScreen = true

	return view
}

func (m model) documentView() tea.View {
	title := "Storyden"
	body := ""
	var err error

	switch m.screen {
	case screenNode:
		if m.nodeView != nil {
			title = string(m.nodeView.Name)
			body, err = render.NodeViewString(m.out, m.nodeView)
		}
	case screenThread:
		if m.threadView != nil {
			title = string(m.threadView.Title)
			body, err = render.ThreadViewString(m.out, m.threadView)
		}
	}
	if err != nil {
		body = err.Error()
	}
	if m.loading {
		body = "Loading..."
	}
	if m.err != nil {
		body = m.err.Error()
	}

	bodyHeight := max(m.height-4, 1)
	body = visibleLines(body, m.scroll, bodyHeight)

	header := sharedtui.Title.Render(title)
	footer := sharedtui.Muted.Render("↑/↓ scroll  pgup/pgdn page  backspace/esc return  q quit")

	view := tea.NewView(header + "\n\n" + body + "\n" + footer)
	view.AltScreen = true

	return view
}

func (m model) renderSidebar(width int, height int) string {
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(lipgloss.Color("238")).
		PaddingRight(1)

	item := func(label string, active bool) string {
		if active {
			return sharedtui.Accent.Render("▸ " + label)
		}
		return sharedtui.Muted.Render("  " + label)
	}

	lines := []string{
		sharedtui.Title.Render("Resources"),
		"",
		item("Nodes", m.active == viewNodes),
		item("Threads", m.active == viewThreads),
	}

	if m.active == viewNodes && len(m.nodeStack) > 0 {
		lines = append(lines, "", sharedtui.Muted.Render("Path"))
		for _, node := range m.nodeStack {
			lines = append(lines, "  "+string(node.Slug))
		}
	}

	return style.Render(strings.Join(lines, "\n"))
}

func (m model) renderContent(width int, height int) string {
	pager := ""
	if m.inPagedRootView() {
		p := m.currentPaginator()
		pager = fmt.Sprintf(" page %d/%d", p.Page+1, max(p.TotalPages, 1))
	}

	title := titleCase(m.active) + pager
	style := lipgloss.NewStyle().Width(width).Height(height).Padding(0, 1)

	return style.Render(sharedtui.Title.Render(title) + "\n\n" + m.table.View())
}

func (m model) renderDetails(width int, height int) string {
	body := ""
	switch m.active {
	case viewNodes:
		if node, ok := m.selectedNode(); ok {
			body = nodeDetails(node)
		}
	case viewThreads:
		if thread, ok := m.selectedThread(); ok {
			body = threadDetails(thread)
		}
	}

	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color("238")).
		PaddingLeft(1)

	return style.Render(sharedtui.Title.Render("Details") + "\n\n" + render.ClampLines(body, height-3))
}

func visibleLines(text string, offset int, height int) string {
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	if offset > len(lines) {
		offset = max(len(lines)-1, 0)
	}

	end := min(offset+height, len(lines))
	if offset >= end {
		return ""
	}

	return strings.Join(lines[offset:end], "\n") + "\n"
}
