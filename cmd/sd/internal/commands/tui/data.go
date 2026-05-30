package tui

import (
	"context"
	"net/http"
	"strconv"

	tea "charm.land/bubbletea/v2"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/nodeapi"
	outputfmt "github.com/Southclaws/storyden/cmd/sd/internal/output"
	"github.com/Southclaws/storyden/cmd/sd/internal/threadapi"
)

type nodePageMsg struct {
	seq    int
	result *openapi.NodeListResult
	err    error
}

type threadPageMsg struct {
	seq    int
	result *openapi.ThreadListResult
	err    error
}

type nodeViewMsg struct {
	seq  int
	node *openapi.NodeWithChildren
	err  error
}

type threadViewMsg struct {
	seq    int
	thread *openapi.Thread
	err    error
}

func fetchNodePage(ctx context.Context, client *openapi.ClientWithResponses, page int, seq int) tea.Cmd {
	return func() tea.Msg {
		result, err := fetchNodes(ctx, client, page)
		return nodePageMsg{seq: seq, result: result, err: err}
	}
}

func fetchThreadPage(ctx context.Context, client *openapi.ClientWithResponses, page int, seq int) tea.Cmd {
	return func() tea.Msg {
		result, err := fetchThreads(ctx, client, page)
		return threadPageMsg{seq: seq, result: result, err: err}
	}
}

func fetchNodeView(ctx context.Context, client *openapi.ClientWithResponses, slug string, seq int) tea.Cmd {
	return func() tea.Msg {
		node, err := nodeapi.Fetch(ctx, client, slug)
		return nodeViewMsg{seq: seq, node: node, err: err}
	}
}

func fetchThreadView(ctx context.Context, client *openapi.ClientWithResponses, mark string, seq int) tea.Cmd {
	return func() tea.Msg {
		thread, err := threadapi.Fetch(ctx, client, mark)
		return threadViewMsg{seq: seq, thread: thread, err: err}
	}
}

func fetchNodes(ctx context.Context, client *openapi.ClientWithResponses, page int) (*openapi.NodeListResult, error) {
	pageQuery := openapi.PaginationQuery(strconv.Itoa(page))
	response, err := client.NodeListWithResponse(ctx, &openapi.NodeListParams{Page: &pageQuery})
	if err != nil {
		return nil, err
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, outputfmt.RequestErrorWithMessages(
			"node list request",
			response,
			response.Body,
			outputfmt.UnauthorizedMessage("node list request"),
		)
	}
	return response.JSON200, nil
}

func fetchThreads(ctx context.Context, client *openapi.ClientWithResponses, page int) (*openapi.ThreadListResult, error) {
	pageQuery := openapi.PaginationQuery(strconv.Itoa(page))
	response, err := client.ThreadListWithResponse(ctx, &openapi.ThreadListParams{Page: &pageQuery})
	if err != nil {
		return nil, err
	}
	if response.StatusCode() != http.StatusOK || response.JSON200 == nil {
		return nil, outputfmt.RequestErrorWithMessages(
			"thread list request",
			response,
			response.Body,
			outputfmt.UnauthorizedMessage("thread list request"),
		)
	}
	return response.JSON200, nil
}
