package tui

import tea "charm.land/bubbletea/v2"

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = max(msg.Width, 80)
		m.height = max(msg.Height, 20)
		m.table = m.buildTable()

	case nodePageMsg:
		m = m.handleNodePage(msg)

	case threadPageMsg:
		m = m.handleThreadPage(msg)

	case nodeViewMsg:
		m = m.handleNodeView(msg)

	case threadViewMsg:
		m = m.handleThreadView(msg)

	case tea.KeyPressMsg:
		if m.screen != screenList {
			return m.updateDocument(msg)
		}

		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		case "n":
			m.active = viewNodes
			m.err = nil
			m.table = m.buildTable()
			return m, nil
		case "t":
			m.active = viewThreads
			m.err = nil
			m.table = m.buildTable()
			return m, nil
		case "r":
			return m.startFetch(m.currentPage())
		case "enter":
			return m.startOpenSelected()
		case "c":
			if node, ok := m.selectedNode(); m.active == viewNodes && ok && len(node.Children) > 0 {
				m.nodeStack = append(m.nodeStack, node)
				m.err = nil
				m.table = m.buildTable()
			}
			return m, nil
		case "backspace":
			if m.active == viewNodes && len(m.nodeStack) > 0 {
				m.nodeStack = m.nodeStack[:len(m.nodeStack)-1]
				m.err = nil
				m.table = m.buildTable()
				return m, nil
			}
		}

		if m.inPagedRootView() {
			before := m.currentPaginator().Page
			var cmd tea.Cmd
			switch m.active {
			case viewNodes:
				m.nodePaginator, cmd = m.nodePaginator.Update(msg)
			case viewThreads:
				m.threadPaginator, cmd = m.threadPaginator.Update(msg)
			}
			if m.currentPaginator().Page != before {
				return m.startFetch(m.currentPaginator().Page + 1)
			}
			if cmd != nil {
				return m, cmd
			}
		}
	}

	updatedTable, cmd := m.table.Update(msg)
	m.table = updatedTable

	return m, cmd
}

func (m model) updateDocument(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "esc", "backspace", "enter":
		m = m.closeDocument()
		return m, nil
	case "up", "k":
		m.scroll = max(m.scroll-1, 0)
		return m, nil
	case "down", "j":
		m.scroll++
		return m, nil
	case "pgup", "b":
		m.scroll = max(m.scroll-(m.height-4), 0)
		return m, nil
	case "pgdown", "f", " ":
		m.scroll += max(m.height-4, 1)
		return m, nil
	case "home", "g":
		m.scroll = 0
		return m, nil
	}

	return m, nil
}
