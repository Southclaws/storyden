package tui

import (
	"context"
	"io"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/paginator"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	sharedtui "github.com/Southclaws/storyden/cmd/sd/internal/tui"
)

const (
	viewNodes   = "nodes"
	viewThreads = "threads"

	screenList   = "list"
	screenNode   = "node"
	screenThread = "thread"
)

type model struct {
	ctx    context.Context
	client *openapi.ClientWithResponses
	out    io.Writer

	screen string
	active string
	table  table.Model

	nodes     *openapi.NodeListResult
	nodeStack []openapi.NodeWithChildren
	nodeView  *openapi.NodeWithChildren

	threads    *openapi.ThreadListResult
	threadView *openapi.Thread

	nodePaginator   paginator.Model
	threadPaginator paginator.Model

	width   int
	height  int
	loading bool
	err     error
	seq     int
	pending int
	scroll  int
}

func newModel(ctx context.Context, client *openapi.ClientWithResponses, out io.Writer, nodes *openapi.NodeListResult, threads *openapi.ThreadListResult) model {
	m := model{
		ctx:     ctx,
		client:  client,
		out:     out,
		screen:  screenList,
		active:  viewNodes,
		nodes:   nodes,
		threads: threads,
		width:   120,
		height:  32,
	}
	m.nodePaginator = newPaginator(nodes.CurrentPage, nodes.TotalPages, nodes.PageSize)
	m.threadPaginator = newPaginator(threads.CurrentPage, threads.TotalPages, threads.PageSize)
	m.table = m.buildTable()

	return m
}

func newPaginator(currentPage int, totalPages int, pageSize int) paginator.Model {
	p := paginator.New(
		paginator.WithTotalPages(max(totalPages, 1)),
		paginator.WithPerPage(max(pageSize, 1)),
	)
	p.Page = max(currentPage, 1) - 1
	p.Type = paginator.Dots
	p.ActiveDot = sharedtui.Accent.Render("●")
	p.InactiveDot = sharedtui.Muted.Render("·")
	p.KeyMap = paginator.KeyMap{
		PrevPage: key.NewBinding(key.WithKeys("left", "h")),
		NextPage: key.NewBinding(key.WithKeys("right", "l")),
	}

	return p
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) currentNodes() []openapi.NodeWithChildren {
	if len(m.nodeStack) == 0 {
		return m.nodes.Nodes
	}

	return m.nodeStack[len(m.nodeStack)-1].Children
}

func (m model) selectedNode() (openapi.NodeWithChildren, bool) {
	nodes := m.currentNodes()
	idx := m.table.Cursor()
	if idx < 0 || idx >= len(nodes) {
		return openapi.NodeWithChildren{}, false
	}

	return nodes[idx], true
}

func (m model) selectedThread() (openapi.ThreadReference, bool) {
	idx := m.table.Cursor()
	if m.threads == nil || idx < 0 || idx >= len(m.threads.Threads) {
		return openapi.ThreadReference{}, false
	}

	return m.threads.Threads[idx], true
}

func (m model) inPagedRootView() bool {
	return m.active == viewThreads || (m.active == viewNodes && len(m.nodeStack) == 0)
}

func (m model) currentPaginator() paginator.Model {
	if m.active == viewThreads {
		return m.threadPaginator
	}

	return m.nodePaginator
}

func (m model) currentPage() int {
	return m.currentPaginator().Page + 1
}

func (m model) startFetch(page int) (tea.Model, tea.Cmd) {
	m.loading = true
	m.err = nil
	m.seq++
	m.pending = m.seq

	if m.active == viewThreads {
		return m, fetchThreadPage(m.ctx, m.client, page, m.pending)
	}

	m.nodeStack = nil
	return m, fetchNodePage(m.ctx, m.client, page, m.pending)
}

func (m model) startOpenSelected() (tea.Model, tea.Cmd) {
	m.loading = true
	m.err = nil
	m.seq++
	m.pending = m.seq

	switch m.active {
	case viewThreads:
		thread, ok := m.selectedThread()
		if !ok {
			m.loading = false
			m.pending = 0
			return m, nil
		}
		return m, fetchThreadView(m.ctx, m.client, string(thread.Slug), m.pending)
	default:
		node, ok := m.selectedNode()
		if !ok {
			m.loading = false
			m.pending = 0
			return m, nil
		}
		return m, fetchNodeView(m.ctx, m.client, string(node.Slug), m.pending)
	}
}

func (m model) handleNodePage(msg nodePageMsg) model {
	if msg.seq != m.pending {
		return m
	}
	m.loading = false
	m.pending = 0
	if msg.err != nil {
		m.err = msg.err
		return m
	}

	m.nodes = msg.result
	m.nodePaginator = newPaginator(msg.result.CurrentPage, msg.result.TotalPages, msg.result.PageSize)
	m.table = m.buildTable()
	return m
}

func (m model) handleThreadPage(msg threadPageMsg) model {
	if msg.seq != m.pending {
		return m
	}
	m.loading = false
	m.pending = 0
	if msg.err != nil {
		m.err = msg.err
		return m
	}

	m.threads = msg.result
	m.threadPaginator = newPaginator(msg.result.CurrentPage, msg.result.TotalPages, msg.result.PageSize)
	m.table = m.buildTable()
	return m
}

func (m model) handleNodeView(msg nodeViewMsg) model {
	if msg.seq != m.pending {
		return m
	}
	m.loading = false
	m.pending = 0
	if msg.err != nil {
		m.err = msg.err
		return m
	}

	m.screen = screenNode
	m.nodeView = msg.node
	m.threadView = nil
	m.scroll = 0
	return m
}

func (m model) handleThreadView(msg threadViewMsg) model {
	if msg.seq != m.pending {
		return m
	}
	m.loading = false
	m.pending = 0
	if msg.err != nil {
		m.err = msg.err
		return m
	}

	m.screen = screenThread
	m.threadView = msg.thread
	m.nodeView = nil
	m.scroll = 0
	return m
}

func (m model) closeDocument() model {
	m.screen = screenList
	m.nodeView = nil
	m.threadView = nil
	m.scroll = 0
	m.err = nil
	m.table = m.buildTable()
	return m
}
