package search_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/search/search_indexer"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestBleveThreadSearch(t *testing.T) {
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
		idx *search_indexer.Indexer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)

			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			catResp, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Name:   "test-category-" + uuid.NewString(),
				Colour: "#123456",
			}, adminSession)
			tests.Ok(t, err, catResp)

			threadFox, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "The Quick Brown Fox",
				Body:       opt.New("<p>A thread about a quick brown fox jumping over lazy dogs</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadFox)

			threadQuantum, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Understanding Quantum Computing",
				Body:       opt.New("<p>A deep dive into quantum mechanics and computing principles</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadQuantum)

			threadPancakes, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Recipe for Perfect Pancakes",
				Body:       opt.New("<p>Learn how to make fluffy pancakes with this simple recipe</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadPancakes)

			threadJavaScript, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "JavaScript Tutorial",
				Body:       opt.New("<p>Learn JavaScript basics</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadJavaScript)

			threadJava, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Java Programming Guide",
				Body:       opt.New("<p>Java programming fundamentals</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadJava)

			threadMatlab, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Matlab Data Science",
				Body:       opt.New("<p>Data science with Matlab</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadMatlab)

			threadChinese, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "机器学习入门指南",
				Body:       opt.New("<p>这是一个关于机器学习的基础教程</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadChinese)

			threadChineseDeepLearning, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "深度学习与神经网络",
				Body:       opt.New("<p>深入探讨深度学习和神经网络的原理与应用</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadChineseDeepLearning)

			threadRussian, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Программирование на Python",
				Body:       opt.New("<p>Изучение основ программирования на языке Python</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadRussian)

			threadRussianWeb, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Основы веб-разработки",
				Body:       opt.New("<p>Полное руководство по современной веб-разработке</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadRussianWeb)

			threadArabic, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "تطوير تطبيقات الويب",
				Body:       opt.New("<p>دورة كاملة في تطوير تطبيقات الويب الحديثة</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadArabic)

			err = idx.ReindexAll(root)
			r.NoError(err, "failed to reindex all items")

			t.Run("exact_match", func(t *testing.T) {
				r := require.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "Quantum Computing",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.NotNil(findThreadItem(resp.JSON200.Items, threadQuantum.JSON200.Id), "should find thread with exact match")
			})

			t.Run("prefix_match", func(t *testing.T) {
				r := require.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "Quick",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.NotNil(findThreadItem(resp.JSON200.Items, threadFox.JSON200.Id), "should find thread with prefix match")
			})

			t.Run("nonsense_no_results", func(t *testing.T) {
				r := require.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "xyzabc123impossible",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.Len(resp.JSON200.Items, 0, "should return no results for nonsense query")
			})

			t.Run("verify_all_threads_indexed", func(t *testing.T) {
				r := require.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}

				resp1, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "Quick",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp1)
				r.NotNil(findThreadItem(resp1.JSON200.Items, threadFox.JSON200.Id))

				resp2, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "Quantum",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp2)
				r.NotNil(findThreadItem(resp2.JSON200.Items, threadQuantum.JSON200.Id))

				resp3, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "Pancakes",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp3)
				r.NotNil(findThreadItem(resp3.JSON200.Items, threadPancakes.JSON200.Id))
			})

			t.Run("chinese_search_machine_learning", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "机器学习",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Chinese thread")
				a.NotNil(findThreadItem(resp.JSON200.Items, threadChinese.JSON200.Id), "should find the Chinese thread about machine learning")
			})

			t.Run("chinese_search_deep_learning", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "深度学习",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Chinese thread")
				a.NotNil(findThreadItem(resp.JSON200.Items, threadChineseDeepLearning.JSON200.Id), "should find the Chinese thread about deep learning")
			})

			t.Run("russian_search_programming", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "Программирование",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Russian thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadRussian.JSON200.Id), "should find the Russian thread about Python programming")
			})

			t.Run("russian_search_web", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "веб-разработки",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Russian thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadRussianWeb.JSON200.Id), "should find the Russian thread about web development")
			})

			t.Run("arabic_search", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "تطوير",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Arabic thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadArabic.JSON200.Id), "should find the Arabic thread about web app development")
			})

			// -
			// Match tests - typeahead search endpoint
			// -

			t.Run("prefix_typeahead_java", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "Jav",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 2, "should find both Java and JavaScript")
				ids := matchItemIDs(resp.JSON200.Items)
				a.Contains(ids, threadJavaScript.JSON200.Id, "should find JavaScript")
				a.Contains(ids, threadJava.JSON200.Id, "should find Java")
			})

			t.Run("prefix_typeahead_javascript", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "JavaS",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.Len(resp.JSON200.Items, 1, "should find exactly one result")
				a.Equal(threadJavaScript.JSON200.Id, resp.JSON200.Items[0].Id, "should find JavaScript")
			})

			t.Run("prefix_typeahead_matlab_data", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "Matlab Da",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find at least one result")
				a.Contains(matchItemIDs(resp.JSON200.Items), threadMatlab.JSON200.Id, "should find Matlab Data Science thread")
			})

			t.Run("no_results", func(t *testing.T) {
				r := require.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "Rust",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.Len(resp.JSON200.Items, 0, "should find no results for Rust")
			})

			t.Run("chinese_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "机",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Chinese thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadChinese.JSON200.Id), "should find the Chinese thread about machine learning")
			})

			t.Run("russian_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "програ",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Russian thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadRussian.JSON200.Id), "should find the Russian thread about Python programming")
			})

			t.Run("arabic_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "تطوير",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Arabic thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadArabic.JSON200.Id), "should find the Arabic thread about learning programming")
			})
		}))
	}))
}

func findThreadItem(items []openapi.DatagraphItem, id openapi.Identifier) *openapi.DatagraphItemThread {
	for _, item := range items {
		if threadItem, err := item.AsDatagraphItemThread(); err == nil {
			if threadItem.Ref.Id == id {
				return &threadItem
			}
		}
	}
	return nil
}

func findMatchItem(items []openapi.DatagraphMatch, id openapi.Identifier) *openapi.DatagraphMatch {
	for _, item := range items {
		if item.Id == id {
			return &item
		}
	}
	return nil
}

func matchItemIDs(items []openapi.DatagraphMatch) []string {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.Id
	}
	return ids
}
