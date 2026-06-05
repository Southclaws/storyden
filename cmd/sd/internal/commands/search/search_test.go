package search

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func TestOptionsValidate(t *testing.T) {
	r := require.New(t)

	r.NoError((&options{Query: "docs", Kinds: []string{"node", "THREAD"}}).validate())
	r.ErrorContains((&options{Query: " "}).validate(), "search query must not be empty")
	r.ErrorContains((&options{Query: "docs", Kinds: []string{"garbage"}}).validate(), "invalid --kind")
}

func TestFetchSearchSendsEveryAPIParameter(t *testing.T) {
	r := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		r.Equal("docs", query.Get("q"))
		r.ElementsMatch([]string{"node", "thread"}, query["kind"])
		r.ElementsMatch([]string{"southclaws", "alice"}, query["authors"])
		r.ElementsMatch([]string{"guides", "release-notes"}, query["categories"])
		r.ElementsMatch([]string{"cli", "review"}, query["tags"])
		r.Equal("3", query.Get("page"))

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"current_page":3,"page_size":0,"results":0,"items":[]}`))
		r.NoError(err)
	}))
	defer server.Close()

	client, err := openapi.NewClientWithResponses(server.URL)
	r.NoError(err)

	result, err := fetchSearch(context.Background(), client, &options{
		Query:      "docs",
		Kinds:      []string{"node", "thread"},
		Authors:    []string{"southclaws", "alice"},
		Categories: []string{"guides", "release-notes"},
		Tags:       []string{"cli", "review"},
	}, 3)
	r.NoError(err)
	r.Equal(3, result.CurrentPage)
}

func TestItemRowNode(t *testing.T) {
	r := require.New(t)

	item := openapi.DatagraphItem{}
	r.NoError(item.FromDatagraphItemNode(openapi.DatagraphItemNode{
		Ref: openapi.Node{
			Id:          "node-id",
			Name:        "Documentation Hub",
			Slug:        "documentation-hub",
			Description: "Main docs page",
			Owner: openapi.ProfileReference{
				Name: "Storyden",
			},
			UpdatedAt: time.Date(2026, 5, 29, 10, 56, 0, 0, time.UTC),
		},
	}))

	row := itemRow(item)

	r.Equal("node", row.Kind)
	r.Equal("Documentation Hub", row.Title)
	r.Equal("node-id", row.ID)
	r.Equal("documentation-hub", row.Slug)
	r.Equal("Storyden", row.Author)
	r.Equal("Main docs page", row.Summary)
}
