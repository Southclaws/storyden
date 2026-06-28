package chat_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/app/transports/sse"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
	"github.com/Southclaws/storyden/tests/robot"
)

func TestRobotChatContentSearch(t *testing.T) {
	t.Parallel()
	integration.Test(t,
		&config.Config{},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelContentSearch),
		sse.Build(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				rob := tests.AssertRequest(cl.RobotCreateWithResponse(root, openapi.RobotCreateJSONRequestBody{
					Name:        "content-search-robot-" + xid.New().String(),
					Description: "robot for content search tool tests",
					Playbook:    "you are a test robot",
					Tools:       robotToolsPtr("content_search"),
				}, adminSession))(t, http.StatusOK)
				robotID := string(rob.JSON200.Id)

				vis := openapi.VisibilityPublished
				tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeCreateJSONRequestBody{
					Name:       "Magnolia Library Page",
					Visibility: &vis,
				}, adminSession))(t, http.StatusOK)

				cat := tests.AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
					Colour:      "#fe4efd",
					Description: "search test category",
					Name:        "Category " + uuid.NewString(),
				}, adminSession))(t, http.StatusOK)

				tests.AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      "Magnolia Forum Thread",
					Body:       opt.New("<p>discussion about magnolia</p>").Ptr(),
					Category:   opt.New(cat.JSON200.Id).Ptr(),
					Visibility: opt.New(openapi.VisibilityPublished).Ptr(),
				}, adminSession))(t, http.StatusOK)

				t.Run("triggers_content_search", func(t *testing.T) {
					a := assert.New(t)

					sessionID := xid.New().String()
					stream := doChat(t, root, ts, adminSession, sessionID, robotID, "search for magnolia content")

					toolNames := collectToolCalls(stream)
					textDeltas := collectTextDeltas(stream)
					toolOutputs := collectToolOutputs(stream)

					a.Contains(toolNames, "content_search")
					a.Equal("Content search complete.", strings.Join(textDeltas, ""))
					a.NotEmpty(toolOutputs)
					a.Greater(toolOutputResultCount(toolOutputs[0]), float64(0))
				})
			}))
		}),
	)
}

func TestRobotChatThreadSearch(t *testing.T) {
	t.Parallel()
	integration.Test(t,
		&config.Config{},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelThreadSearch),
		sse.Build(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				rob := tests.AssertRequest(cl.RobotCreateWithResponse(root, openapi.RobotCreateJSONRequestBody{
					Name:        "thread-search-robot-" + xid.New().String(),
					Description: "robot for thread search tool tests",
					Playbook:    "you are a test robot",
					Tools:       robotToolsPtr("thread_search"),
				}, adminSession))(t, http.StatusOK)
				robotID := string(rob.JSON200.Id)

				cat := tests.AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
					Colour:      "#fe4efd",
					Description: "search test category",
					Name:        "Category " + uuid.NewString(),
				}, adminSession))(t, http.StatusOK)

				tests.AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      "Magnolia Thread Discussion",
					Body:       opt.New("<p>talk about magnolia trees</p>").Ptr(),
					Category:   opt.New(cat.JSON200.Id).Ptr(),
					Visibility: opt.New(openapi.VisibilityPublished).Ptr(),
				}, adminSession))(t, http.StatusOK)

				t.Run("triggers_thread_search", func(t *testing.T) {
					a := assert.New(t)

					sessionID := xid.New().String()
					stream := doChat(t, root, ts, adminSession, sessionID, robotID, "search threads about magnolia")

					toolNames := collectToolCalls(stream)
					textDeltas := collectTextDeltas(stream)

					a.Contains(toolNames, "thread_search")
					a.Equal("Thread search complete.", strings.Join(textDeltas, ""))
				})
			}))
		}),
	)
}

func TestRobotChatReplySearch(t *testing.T) {
	t.Parallel()
	integration.Test(t,
		&config.Config{},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelReplySearch),
		sse.Build(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				rob := tests.AssertRequest(cl.RobotCreateWithResponse(root, openapi.RobotCreateJSONRequestBody{
					Name:        "reply-search-robot-" + xid.New().String(),
					Description: "robot for reply search tool tests",
					Playbook:    "you are a test robot",
					Tools:       robotToolsPtr("reply_search"),
				}, adminSession))(t, http.StatusOK)
				robotID := string(rob.JSON200.Id)

				cat := tests.AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
					Colour:      "#fe4efd",
					Description: "search test category",
					Name:        "Category " + uuid.NewString(),
				}, adminSession))(t, http.StatusOK)

				thread := tests.AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      "Reply Search Test Thread",
					Body:       opt.New("<p>base thread for reply tests</p>").Ptr(),
					Category:   opt.New(cat.JSON200.Id).Ptr(),
					Visibility: opt.New(openapi.VisibilityPublished).Ptr(),
				}, adminSession))(t, http.StatusOK)

				tests.AssertRequest(cl.ReplyCreateWithResponse(root, thread.JSON200.Slug, openapi.ReplyInitialProps{
					Body: "<p>magnolia blossom reply content</p>",
				}, adminSession))(t, http.StatusOK)

				t.Run("triggers_reply_search", func(t *testing.T) {
					a := assert.New(t)

					sessionID := xid.New().String()
					stream := doChat(t, root, ts, adminSession, sessionID, robotID, "search replies about magnolia")

					toolNames := collectToolCalls(stream)
					textDeltas := collectTextDeltas(stream)

					a.Contains(toolNames, "reply_search")
					a.Equal("Reply search complete.", strings.Join(textDeltas, ""))
				})
			}))
		}),
	)
}

func TestRobotChatPostSearch(t *testing.T) {
	t.Parallel()
	integration.Test(t,
		&config.Config{},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelPostSearch),
		sse.Build(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				rob := tests.AssertRequest(cl.RobotCreateWithResponse(root, openapi.RobotCreateJSONRequestBody{
					Name:        "post-search-robot-" + xid.New().String(),
					Description: "robot for post search tool tests",
					Playbook:    "you are a test robot",
					Tools:       robotToolsPtr("post_search"),
				}, adminSession))(t, http.StatusOK)
				robotID := string(rob.JSON200.Id)

				cat := tests.AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
					Colour:      "#fe4efd",
					Description: "search test category",
					Name:        "Category " + uuid.NewString(),
				}, adminSession))(t, http.StatusOK)

				tests.AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      "Magnolia Post Thread",
					Body:       opt.New("<p>magnolia post content</p>").Ptr(),
					Category:   opt.New(cat.JSON200.Id).Ptr(),
					Visibility: opt.New(openapi.VisibilityPublished).Ptr(),
				}, adminSession))(t, http.StatusOK)

				t.Run("triggers_post_search", func(t *testing.T) {
					a := assert.New(t)

					sessionID := xid.New().String()
					stream := doChat(t, root, ts, adminSession, sessionID, robotID, "search posts about magnolia")

					toolNames := collectToolCalls(stream)
					textDeltas := collectTextDeltas(stream)

					a.Contains(toolNames, "post_search")
					a.Equal("Post search complete.", strings.Join(textDeltas, ""))
				})
			}))
		}),
	)
}

func TestRobotChatMemberSearch(t *testing.T) {
	t.Parallel()
	integration.Test(t,
		&config.Config{},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelMemberSearch),
		sse.Build(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				rob := tests.AssertRequest(cl.RobotCreateWithResponse(root, openapi.RobotCreateJSONRequestBody{
					Name:        "member-search-robot-" + xid.New().String(),
					Description: "robot for member search tool tests",
					Playbook:    "you are a test robot",
					Tools:       robotToolsPtr("member_search"),
				}, adminSession))(t, http.StatusOK)
				robotID := string(rob.JSON200.Id)

				t.Run("triggers_member_search", func(t *testing.T) {
					a := assert.New(t)

					sessionID := xid.New().String()
					stream := doChat(t, root, ts, adminSession, sessionID, robotID, "find member odin")

					toolNames := collectToolCalls(stream)
					textDeltas := collectTextDeltas(stream)

					a.Contains(toolNames, "member_search")
					a.Equal("Member search complete.", strings.Join(textDeltas, ""))
				})
			}))
		}),
	)
}
