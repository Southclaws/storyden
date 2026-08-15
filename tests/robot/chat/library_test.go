package chat_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/seed"
	authsession "github.com/Southclaws/storyden/app/services/authentication/session"
	robot_tools "github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/app/transports/sse"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/lib/mcp"
	"github.com/Southclaws/storyden/tests"
	"github.com/Southclaws/storyden/tests/robot"
)

func TestRobotChatLibraryPageList(t *testing.T) {
	t.Parallel()
	integration.Test(t,
		&config.Config{
			LanguageModelProvider: "mock",
		},
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
			toolRegistry *robot_tools.Registry,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminCtx = authsession.WithAccountPermissions(
					adminCtx,
					*admin,
					rbac.NewList(rbac.PermissionReadPublishedLibrary, rbac.PermissionManageLibrary),
				)
				adminSession := sh.WithSession(adminCtx)

				rob := tests.AssertRequest(cl.RobotCreateWithResponse(root, openapi.RobotCreateJSONRequestBody{
					Name:        "test-robot-" + xid.New().String(),
					Description: "robot for library tool tests",
					Playbook:    "you are a test robot",
					Toolsets:    robotToolsetsPtr("system.library"),
				}, adminSession))(t, http.StatusOK)
				robotID := string(rob.JSON200.Id)

				vis := openapi.VisibilityPublished
				tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeCreateJSONRequestBody{
					Name:       "Robot Test Page Alpha",
					Visibility: &vis,
				}, adminSession))(t, http.StatusOK)

				publishedPage := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeCreateJSONRequestBody{
					Name:       "Robot Test Page Beta",
					Visibility: &vis,
				}, adminSession))(t, http.StatusOK)

				reviewVisibility := openapi.VisibilityReview
				reviewContent := "This page is ready for a substantive review."
				reviewPage := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeCreateJSONRequestBody{
					Name:       "Robot Review Page",
					Content:    &reviewContent,
					Visibility: &reviewVisibility,
				}, adminSession))(t, http.StatusOK)

				t.Run("filters_pages_by_review_visibility", func(t *testing.T) {
					tool, err := toolRegistry.GetTool(adminCtx, "library_page_list")
					require.NoError(t, err)

					input, err := json.Marshal(mcp.ToolLibraryPageTreeInput{
						Visibility: opt.New(mcp.ToolLibraryPageTreeInputVisibilityReview).Ptr(),
					})
					require.NoError(t, err)
					rawOutput, err := tool.Handler(adminCtx, input)
					require.NoError(t, err)

					var output mcp.ToolLibraryPageTreeOutput
					require.NoError(t, json.Unmarshal(rawOutput, &output))

					pageIDs := make([]string, 0, len(output.Pages))
					for _, page := range output.Pages {
						pageIDs = append(pageIDs, page.Id)
						assert.Equal(t, "review", string(page.Visibility))
					}
					assert.Contains(t, pageIDs, string(reviewPage.JSON200.Id))
					assert.NotContains(t, pageIDs, string(publishedPage.JSON200.Id))
					assert.Equal(t, len(output.Pages), output.Results)
				})

				t.Run("gets_page_content_and_visibility", func(t *testing.T) {
					tool, err := toolRegistry.GetTool(adminCtx, "get_library_page")
					require.NoError(t, err)

					input, err := json.Marshal(mcp.ToolLibraryPageGetInput{
						Id: string(reviewPage.JSON200.Id),
					})
					require.NoError(t, err)
					rawOutput, err := tool.Handler(adminCtx, input)
					require.NoError(t, err)

					var output mcp.ToolLibraryPageGetOutput
					require.NoError(t, json.Unmarshal(rawOutput, &output))
					assert.Equal(t, "This page is ready for a substantive review.", output.Content)
					assert.Equal(t, "review", string(output.Visibility))
				})

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
		&config.Config{
			LanguageModelProvider: "mock",
		},
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
					Toolsets:    robotToolsetsPtr("system.library"),
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
