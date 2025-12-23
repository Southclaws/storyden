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

			threadSpanish, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Recetas de Cocina Mediterránea",
				Body:       opt.New("<p>Descubre los secretos de la cocina mediterránea tradicional</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadSpanish)

			threadFrench, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Histoire de l'Architecture Gothique",
				Body:       opt.New("<p>Exploration des cathédrales gothiques européennes et leur influence</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadFrench)

			threadGerman, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Wandern in den Alpen",
				Body:       opt.New("<p>Die besten Wanderwege und Bergtouren in den Alpen</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadGerman)

			threadPortuguese, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Música Tradicional Brasileira",
				Body:       opt.New("<p>Explorando os ritmos e melodias da música brasileira</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadPortuguese)

			threadGreek, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Αρχαία Ελληνική Φιλοσοφία",
				Body:       opt.New("<p>Μελέτη των έργων των αρχαίων Ελλήνων φιλοσόφων</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadGreek)

			threadTurkish, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Geleneksel Türk Mutfağı",
				Body:       opt.New("<p>Türk mutfağının zengin lezzetlerini keşfedin</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadTurkish)

			threadGeorgian, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "ქართული ხალხური სიმღერები",
				Body:       opt.New("<p>ქართული პოლიფონიური სიმღერების ისტორია და მნიშვნელობა</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadGeorgian)

			threadHindi, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "भारतीय शास्त्रीय संगीत",
				Body:       opt.New("<p>भारतीय शास्त्रीय संगीत की परंपराओं का अन्वेषण</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadHindi)

			threadSwahili, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Hadithi za Kiswahili",
				Body:       opt.New("<p>Hadithi na masimulizi ya kitamaduni kutoka Afrika Mashariki</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadSwahili)

			threadArmenian, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Հայկական Ավանդական Խոհանոց",
				Body:       opt.New("<p>Հայկական խոհանոցի պատմությունը և ավանդական բաղադրատոմսերը</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadArmenian)

			threadHebrew, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "ספרות עברית מודרנית",
				Body:       opt.New("<p>סקירה של הספרות העברית המודרנית והסופרים המשפיעים</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadHebrew)

			threadPersian, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "شعر کلاسیک فارسی",
				Body:       opt.New("<p>بررسی شاعران بزرگ فارسی و آثار ماندگار آنها</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadPersian)

			threadUrdu, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "اردو شاعری کی روایت",
				Body:       opt.New("<p>اردو شاعری کی تاریخ اور مشہور شعراء کا تعارف</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadUrdu)

			threadPunjabi, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "ਪੰਜਾਬੀ ਲੋਕ ਗੀਤ",
				Body:       opt.New("<p>ਪੰਜਾਬੀ ਸੱਭਿਆਚਾਰ ਵਿੱਚ ਲੋਕ ਗੀਤਾਂ ਦਾ ਮਹੱਤਵ</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadPunjabi)

			threadNepali, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "नेपाली पर्वतारोहण",
				Body:       opt.New("<p>नेपालको हिमाल र पर्वतारोहणको इतिहास</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadNepali)

			threadYoruba, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Àṣà Yorùbá",
				Body:       opt.New("<p>Ìtàn àti àṣà àwọn ènìyàn Yorùbá ní Nàìjíríà</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadYoruba)

			threadIgbo, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Omenala Igbo",
				Body:       opt.New("<p>Akụkọ ọdịnala na omenala ndị Igbo</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadIgbo)

			threadHausa, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Tarihin Hausawa",
				Body:       opt.New("<p>Tarihi da al'adun Hausawa a Arewacin Najeriya</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadHausa)

			threadAkan, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Akanfo Atetesɛm",
				Body:       opt.New("<p>Akanman mu atetesɛm ne amammerɛ</p>").Ptr(),
				Category:   opt.New(catResp.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
			}, adminSession)
			tests.Ok(t, err, threadAkan)

			ctx1, authorOne := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			session1 := sh.WithSession(ctx1)
			ctx2, authorTwo := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			session2 := sh.WithSession(ctx2)

			cat1, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Name:   "Tech" + uuid.NewString(),
				Colour: "#FF0000",
			}, adminSession)
			tests.Ok(t, err, cat1)

			cat2, err := cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
				Name:   "Food" + uuid.NewString(),
				Colour: "#00FF00",
			}, adminSession)
			tests.Ok(t, err, cat2)

			hot := "<p>searchable keyword content</p>"

			t1, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Thread by Baldur in Tech with sharing",
				Body:       opt.New(hot).Ptr(),
				Category:   opt.New(cat1.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
				Tags:       &[]openapi.TagName{"sharing"},
			}, session1)
			tests.Ok(t, err, t1)

			t2, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Thread by Loki in Food with tips",
				Body:       opt.New(hot).Ptr(),
				Category:   opt.New(cat2.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
				Tags:       &[]openapi.TagName{"tips"},
			}, session2)
			tests.Ok(t, err, t2)

			t3, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
				Title:      "Thread by Baldur in Tech with sharing and tips",
				Body:       opt.New(hot).Ptr(),
				Category:   opt.New(cat1.JSON200.Id).Ptr(),
				Visibility: opt.New(openapi.Published).Ptr(),
				Tags:       &[]openapi.TagName{"sharing", "tips"},
			}, session1)
			tests.Ok(t, err, t3)

			err = idx.ReindexAll(root)
			r.NoError(err, "failed to reindex all items")

			threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}

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

			t.Run("spanish_search", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "Cocina",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Spanish thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadSpanish.JSON200.Id), "should find the Spanish thread about Mediterranean cooking")
			})

			t.Run("french_search", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "Architecture",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find French thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadFrench.JSON200.Id), "should find the French thread about Gothic architecture")
			})

			t.Run("german_search", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "Wandern",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find German thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadGerman.JSON200.Id), "should find the German thread about hiking in the Alps")
			})

			t.Run("portuguese_search", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "Música",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Portuguese thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadPortuguese.JSON200.Id), "should find the Portuguese thread about Brazilian music")
			})

			t.Run("greek_search", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "Φιλοσοφία",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Greek thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadGreek.JSON200.Id), "should find the Greek thread about ancient philosophy")
			})

			t.Run("turkish_search", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "Mutfağı",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Turkish thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadTurkish.JSON200.Id), "should find the Turkish thread about traditional cuisine")
			})

			t.Run("georgian_search", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "სიმღერები",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Georgian thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadGeorgian.JSON200.Id), "should find the Georgian thread about folk songs")
			})

			t.Run("hindi_search", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "संगीत",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Hindi thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadHindi.JSON200.Id), "should find the Hindi thread about classical music")
			})

			t.Run("swahili_search", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "Hadithi",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Swahili thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadSwahili.JSON200.Id), "should find the Swahili thread about stories")
			})

			t.Run("armenian_search", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "Խոհանոց",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Armenian thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadArmenian.JSON200.Id), "should find the Armenian thread about traditional cuisine")
			})

			t.Run("hebrew_search", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "ספרות",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Hebrew thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadHebrew.JSON200.Id), "should find the Hebrew thread about modern literature")
			})

			t.Run("persian_search", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "شعر",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Persian thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadPersian.JSON200.Id), "should find the Persian thread about classical poetry")
			})

			t.Run("urdu_search", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "شاعری",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Urdu thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadUrdu.JSON200.Id), "should find the Urdu thread about poetry tradition")
			})

			t.Run("punjabi_search", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "ਲੋਕ",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Punjabi thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadPunjabi.JSON200.Id), "should find the Punjabi thread about folk songs")
			})

			t.Run("nepali_search", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "पर्वतारोहण",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Nepali thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadNepali.JSON200.Id), "should find the Nepali thread about mountaineering")
			})

			t.Run("yoruba_search", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "Àṣà",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Yoruba thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadYoruba.JSON200.Id), "should find the Yoruba thread about culture")
			})

			t.Run("igbo_search", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "Omenala",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Igbo thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadIgbo.JSON200.Id), "should find the Igbo thread about culture")
			})

			t.Run("hausa_search", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "Tarihin",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Hausa thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadHausa.JSON200.Id), "should find the Hausa thread about history")
			})

			t.Run("akan_search", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "Atetesɛm",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				a.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Akan thread")
				r.NotNil(findThreadItem(resp.JSON200.Items, threadAkan.JSON200.Id), "should find the Akan thread about proverbs")
			})

			// -
			// Filtering tests
			// -

			t.Run("filter_by_author", func(t *testing.T) {
				r := require.New(t)

				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:       "keyword",
					Kind:    &threadKind,
					Authors: &[]openapi.Identifier{openapi.Identifier(authorOne.ID.String())},
				}, session1)
				tests.Ok(t, err, resp)

				r.NotNil(findThreadItem(resp.JSON200.Items, t1.JSON200.Id))
				r.NotNil(findThreadItem(resp.JSON200.Items, t3.JSON200.Id))
				r.Nil(findThreadItem(resp.JSON200.Items, t2.JSON200.Id))
			})

			t.Run("filter_by_category", func(t *testing.T) {
				r := require.New(t)

				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:          "keyword",
					Kind:       &threadKind,
					Categories: &[]openapi.CategorySlug{openapi.CategorySlug(cat1.JSON200.Id)},
				}, session1)
				tests.Ok(t, err, resp)

				r.NotNil(findThreadItem(resp.JSON200.Items, t1.JSON200.Id))
				r.NotNil(findThreadItem(resp.JSON200.Items, t3.JSON200.Id))
				r.Nil(findThreadItem(resp.JSON200.Items, t2.JSON200.Id))
			})

			t.Run("filter_by_single_tag", func(t *testing.T) {
				r := require.New(t)

				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "keyword",
					Kind: &threadKind,
					Tags: &[]openapi.TagName{"sharing"},
				}, session1)
				tests.Ok(t, err, resp)

				r.NotNil(findThreadItem(resp.JSON200.Items, t1.JSON200.Id))
				r.NotNil(findThreadItem(resp.JSON200.Items, t3.JSON200.Id))
				r.Nil(findThreadItem(resp.JSON200.Items, t2.JSON200.Id))
			})

			t.Run("filter_by_multiple_tags_AND", func(t *testing.T) {
				r := require.New(t)

				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "keyword",
					Kind: &threadKind,
					Tags: &[]openapi.TagName{"sharing", "tips"},
				}, session1)
				tests.Ok(t, err, resp)

				r.NotNil(findThreadItem(resp.JSON200.Items, t3.JSON200.Id))
				r.Nil(findThreadItem(resp.JSON200.Items, t1.JSON200.Id))
				r.Nil(findThreadItem(resp.JSON200.Items, t2.JSON200.Id))
			})

			t.Run("filter_by_multiple_authors_OR", func(t *testing.T) {
				r := require.New(t)

				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:    "keyword",
					Kind: &threadKind,
					Authors: &[]openapi.Identifier{
						openapi.Identifier(authorOne.ID.String()),
						openapi.Identifier(authorTwo.ID.String()),
					},
				}, session1)
				tests.Ok(t, err, resp)

				r.NotNil(findThreadItem(resp.JSON200.Items, t1.JSON200.Id))
				r.NotNil(findThreadItem(resp.JSON200.Items, t2.JSON200.Id))
				r.NotNil(findThreadItem(resp.JSON200.Items, t3.JSON200.Id))
			})

			t.Run("filter_combined_author_AND_category_AND_tags", func(t *testing.T) {
				r := require.New(t)

				resp, err := cl.DatagraphSearchWithResponse(root, &openapi.DatagraphSearchParams{
					Q:          "keyword",
					Kind:       &threadKind,
					Authors:    &[]openapi.Identifier{openapi.Identifier(authorOne.ID.String())},
					Categories: &[]openapi.CategorySlug{openapi.CategorySlug(cat1.JSON200.Id)},
					Tags:       &[]openapi.TagName{"sharing"},
				}, session1)
				tests.Ok(t, err, resp)

				r.NotNil(findThreadItem(resp.JSON200.Items, t1.JSON200.Id))
				r.NotNil(findThreadItem(resp.JSON200.Items, t3.JSON200.Id))
				r.Nil(findThreadItem(resp.JSON200.Items, t2.JSON200.Id))
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

			t.Run("spanish_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "Coci",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Spanish thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadSpanish.JSON200.Id), "should find the Spanish thread about Mediterranean cooking")
			})

			t.Run("french_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "Histoir",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find French thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadFrench.JSON200.Id), "should find the French thread about Gothic architecture")
			})

			t.Run("german_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "Wand",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find German thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadGerman.JSON200.Id), "should find the German thread about hiking")
			})

			t.Run("portuguese_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "Músi",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Portuguese thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadPortuguese.JSON200.Id), "should find the Portuguese thread about Brazilian music")
			})

			t.Run("greek_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "Φιλοσ",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Greek thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadGreek.JSON200.Id), "should find the Greek thread about philosophy")
			})

			t.Run("turkish_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "Mutfa",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Turkish thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadTurkish.JSON200.Id), "should find the Turkish thread about cuisine")
			})

			t.Run("georgian_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "სიმღ",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Georgian thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadGeorgian.JSON200.Id), "should find the Georgian thread about folk songs")
			})

			t.Run("hindi_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "संगी",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Hindi thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadHindi.JSON200.Id), "should find the Hindi thread about classical music")
			})

			t.Run("swahili_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "Hadi",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Swahili thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadSwahili.JSON200.Id), "should find the Swahili thread about stories")
			})

			t.Run("armenian_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "Խոհան",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Armenian thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadArmenian.JSON200.Id), "should find the Armenian thread about cuisine")
			})

			t.Run("hebrew_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "ספרו",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Hebrew thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadHebrew.JSON200.Id), "should find the Hebrew thread about literature")
			})

			t.Run("persian_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "شعر",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Persian thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadPersian.JSON200.Id), "should find the Persian thread about poetry")
			})

			t.Run("urdu_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "شاعر",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Urdu thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadUrdu.JSON200.Id), "should find the Urdu thread about poetry")
			})

			t.Run("punjabi_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "ਲੋਕ",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Punjabi thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadPunjabi.JSON200.Id), "should find the Punjabi thread about folk songs")
			})

			t.Run("nepali_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "पर्वत",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Nepali thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadNepali.JSON200.Id), "should find the Nepali thread about mountaineering")
			})

			t.Run("yoruba_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "Àṣà",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Yoruba thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadYoruba.JSON200.Id), "should find the Yoruba thread about culture")
			})

			t.Run("igbo_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "Omen",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Igbo thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadIgbo.JSON200.Id), "should find the Igbo thread about culture")
			})

			t.Run("hausa_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "Tarih",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Hausa thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadHausa.JSON200.Id), "should find the Hausa thread about history")
			})

			t.Run("akan_match", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				threadKind := []openapi.DatagraphItemKind{openapi.DatagraphItemKindThread}
				resp, err := cl.DatagraphMatchesWithResponse(root, &openapi.DatagraphMatchesParams{
					Q:    "Atete",
					Kind: &threadKind,
				}, adminSession)
				tests.Ok(t, err, resp)

				r.GreaterOrEqual(len(resp.JSON200.Items), 1, "should find Akan thread")
				a.NotNil(findMatchItem(resp.JSON200.Items, threadAkan.JSON200.Id), "should find the Akan thread about proverbs")
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
