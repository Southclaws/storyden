package search_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestBleveSearchBasic(t *testing.T) {
	t.Parallel()

	bleveName := time.Now().Format(time.RFC3339) + t.Name()
	cfg := &config.Config{
		SearchProvider: "bleve",
		BlevePath:      fmt.Sprintf("data/%s.bleve", bleveName),
	}

	integration.Test(t, cfg, e2e.Setup(), fx.Invoke(func(
		root context.Context,
		lc fx.Lifecycle,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			catResp, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Name:   uuid.NewString(),
				Colour: "#123456",
			}, adminSession)
			tests.Ok(t, err, catResp)

			published := openapi.Published
			keyword := "bleve-wired-keyword"
			noMatchKeyword := "bleve-no-results"

			threadMatch, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "match-thread",
				Body:       opt.New(fmt.Sprintf("<p>body %s content</p>", keyword)).Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadMatch)

			threadOther, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "other-thread",
				Body:       opt.New("<p>different content</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadOther)

			nodeContent := fmt.Sprintf("<p>node talks about %s</p>", keyword)
			nodeMatch, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
				Name:       "match-node-" + uuid.NewString(),
				Content:    &nodeContent,
				Visibility: &published,
			}, adminSession)
			tests.Ok(t, err, nodeMatch)

			otherContent := "<p>node without magic words</p>"
			nodeOther, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
				Name:       "other-node-" + uuid.NewString(),
				Content:    &otherContent,
				Visibility: &published,
			}, adminSession)
			tests.Ok(t, err, nodeOther)

			waitForItems := func(currentT *testing.T, params openapi.DatagraphSearchParams, ready func([]openapi.DatagraphItem) bool) []openapi.DatagraphItem {
				currentT.Helper()
				var last []openapi.DatagraphItem
				require.Eventually(currentT, func() bool {
					resp, err := cl.DatagraphSearchWithResponse(root, &params, adminSession)
					tests.Ok(currentT, err, resp)
					last = resp.JSON200.Items
					return ready(last)
				}, 10*time.Second, 200*time.Millisecond, "timed out waiting for bleve results")
				return last
			}

			t.Run("search_all_kinds", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)
				items := waitForItems(t, openapi.DatagraphSearchParams{Q: keyword}, func(items []openapi.DatagraphItem) bool {
					return findItem(items, threadMatch.JSON200.Id) != nil && findItem(items, nodeMatch.JSON200.Id) != nil
				})

				threadItem, err := findItem(items, threadMatch.JSON200.Id).AsDatagraphItemThread()
				r.NoError(err)
				a.Equal(threadMatch.JSON200.Id, threadItem.Ref.Id)
				a.Nil(findItem(items, threadOther.JSON200.Id))

				nodeItem, err := findItem(items, nodeMatch.JSON200.Id).AsDatagraphItemNode()
				r.NoError(err)
				a.Equal(nodeMatch.JSON200.Id, nodeItem.Ref.Id)
				a.Nil(findItem(items, nodeOther.JSON200.Id))
			})

			t.Run("search_threads_only", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)
				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				items := waitForItems(t, openapi.DatagraphSearchParams{
					Q:    keyword,
					Kind: &threadKind,
				}, func(items []openapi.DatagraphItem) bool {
					return findItem(items, threadMatch.JSON200.Id) != nil
				})

				threadItem, err := findItem(items, threadMatch.JSON200.Id).AsDatagraphItemThread()
				r.NoError(err)
				a.Equal(threadMatch.JSON200.Id, threadItem.Ref.Id)
				a.Nil(findItem(items, nodeMatch.JSON200.Id))
			})

			t.Run("search_nodes_only", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)
				nodeKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindNode}
				items := waitForItems(t, openapi.DatagraphSearchParams{
					Q:    keyword,
					Kind: &nodeKind,
				}, func(items []openapi.DatagraphItem) bool {
					return findItem(items, nodeMatch.JSON200.Id) != nil
				})

				nodeItem, err := findItem(items, nodeMatch.JSON200.Id).AsDatagraphItemNode()
				r.NoError(err)
				a.Equal(nodeMatch.JSON200.Id, nodeItem.Ref.Id)
				a.Nil(findItem(items, threadMatch.JSON200.Id))
			})

			t.Run("search_no_results", func(t *testing.T) {
				a := assert.New(t)
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{Q: noMatchKeyword}, adminSession)
				tests.Ok(t, err, resp)
				a.Len(resp.JSON200.Items, 0)
			})
		}))
	}))
}

func findItem(items []openapi.DatagraphItem, id openapi.Identifier) *openapi.DatagraphItem {
	for i := range items {
		if datagraphItemID(items[i]) == id {
			return &items[i]
		}
	}
	return nil
}

func datagraphItemID(item openapi.DatagraphItem) string {
	data, err := item.MarshalJSON()
	if err != nil {
		panic(err)
	}

	var v struct {
		Ref struct {
			ID string `json:"id"`
		} `json:"ref"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		panic(err)
	}

	return v.Ref.ID
}
