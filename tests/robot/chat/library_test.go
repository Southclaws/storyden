package chat_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

func TestRobotChatLibraryPageList(t *testing.T) {
	t.Parallel()
	integration.Test(t,
		&config.Config{},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelLibraryPageList),
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
					Name:        "test-robot-" + xid.New().String(),
					Description: "robot for library tool tests",
					Playbook:    "you are a test robot",
					Tools:       robotToolsPtr("library_page_list"),
				}, adminSession))(t, http.StatusOK)
				robotID := string(rob.JSON200.Id)

				vis := openapi.VisibilityPublished
				tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeCreateJSONRequestBody{
					Name:       "Robot Test Page Alpha",
					Visibility: &vis,
				}, adminSession))(t, http.StatusOK)

				tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeCreateJSONRequestBody{
					Name:       "Robot Test Page Beta",
					Visibility: &vis,
				}, adminSession))(t, http.StatusOK)

				t.Run("triggers_library_page_list", func(t *testing.T) {
					a := assert.New(t)

					sessionID := xid.New().String()
					stream := doChat(t, root, ts, adminSession, sessionID, robotID, "list pages")

					toolNames := collectToolCalls(stream)
					textDeltas := collectTextDeltas(stream)

					a.Contains(toolNames, "library_page_list")
					a.Equal("I found the pages in the library.", strings.Join(textDeltas, ""))
				})
			}))
		}),
	)
}

func TestRobotChatLibrarySearchPages(t *testing.T) {
	t.Parallel()
	integration.Test(t,
		&config.Config{},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelLibrarySearchPages),
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
					Name:        "search-robot-" + xid.New().String(),
					Description: "robot for search tool tests",
					Playbook:    "you are a test robot",
					Tools:       robotToolsPtr("library_search_pages"),
				}, adminSession))(t, http.StatusOK)
				robotID := string(rob.JSON200.Id)

				vis := openapi.VisibilityPublished
				tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeCreateJSONRequestBody{
					Name:       "Magnolia Library Page " + xid.New().String(),
					Visibility: &vis,
				}, adminSession))(t, http.StatusOK)

				t.Run("triggers_library_search_pages", func(t *testing.T) {
					a := assert.New(t)

					sessionID := xid.New().String()
					stream := doChat(t, root, ts, adminSession, sessionID, robotID, "search for magnolia pages")

					toolNames := collectToolCalls(stream)
					textDeltas := collectTextDeltas(stream)

					a.Contains(toolNames, "library_search_pages")
					a.Equal("Search for library pages complete.", strings.Join(textDeltas, ""))
				})
			}))
		}),
	)
}
