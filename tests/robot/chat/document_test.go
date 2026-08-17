package chat_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/resources/settings"
	authsession "github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
	"github.com/Southclaws/storyden/tests/robot"
)

func TestRobotOpensAndNavigatesLibraryDocument(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{LanguageModelProvider: "mock"},
		e2e.Setup(),
		robot.WithRobotSettings(mockModelAck),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
			settingsRepo *settings.SettingsRepository,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminCtx = authsession.WithAccountPermissions(adminCtx, *admin, rbac.NewList(rbac.PermissionReadPublishedLibrary, rbac.PermissionManageLibrary))
				adminSession := sh.WithSession(adminCtx)

				visibility := openapi.VisibilityPublished
				var content strings.Builder
				content.WriteString(`<h2>Architecture</h2><p>` + strings.Repeat("The navigation needle belongs in this paragraph. ", 30) + `</p><h2>Operations</h2><p>` + strings.Repeat("Operational details belong elsewhere. ", 30) + `</p>`)
				for index := 3; index <= 30; index++ {
					content.WriteString(fmt.Sprintf("<h2>Topic %d</h2><p>Details for topic %d.</p>", index, index))
				}
				contentHTML := content.String()
				page := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeCreateJSONRequestBody{
					Name: "Document navigation test", Content: &contentHTML, Visibility: &visibility,
				}, adminSession))(t, http.StatusOK)

				scriptName := "robot-chat-document-navigation-" + xid.New().String() + ".yaml"
				scriptPath := filepath.Join("..", "scripts", scriptName)
				writeScript(t, scriptPath, `steps:
  - match:
      tool_result: document_get
    respond:
      text: "Document inspected."
      finish: "stop"
  - match:
      tool_result: library_page_open
    respond:
      tool_calls:
        - id: call_document_get
          name: document_get
          args:
            page: 2
  - match:
      contains: "inspect the document"
    respond:
      tool_calls:
        - id: call_library_page_open
          name: library_page_open
          args:
            id: "`+string(page.JSON200.Id)+`"
`)
				defer os.Remove(scriptPath)
				require.NoError(t, robot.SetRobotSettings(root, settingsRepo, "mock/../scripts/"+scriptName))

				created := tests.AssertRequest(cl.RobotCreateWithResponse(root, openapi.RobotCreateJSONRequestBody{
					Name:        "document-robot-" + xid.New().String(),
					Description: "robot for document navigation tests",
					Playbook:    "Inspect the requested Library document.",
					Toolsets:    robotToolsetsPtr("system.library", "system.documents"),
				}, adminSession))(t, http.StatusOK)

				stream := doChat(t, root, ts, adminSession, xid.New().String(), string(created.JSON200.Id), "inspect the document")
				assert.Equal(t, []string{"library_page_open", "document_get"}, collectToolCalls(stream))
				assert.Equal(t, "Document inspected.", strings.Join(collectTextDeltas(stream), ""))
				inputs := collectToolInputs(stream)
				var documentGetInput map[string]any
				for _, input := range inputs {
					if input.ToolName == "document_get" {
						documentGetInput, _ = input.Input.(map[string]any)
					}
				}
				require.NotNil(t, documentGetInput)
				assert.Equal(t, float64(2), documentGetInput["page"])

				outputs := collectToolOutputs(stream)
				opened := toolOutputByCallID(t, outputs, "call_library_page_open")
				inspected := toolOutputByCallID(t, outputs, "call_document_get")
				assert.NotEmpty(t, opened["document_id"])
				assert.Equal(t, opened["document_id"], inspected["document_id"])
				assert.Equal(t, float64(1), opened["page"])
				assert.Equal(t, float64(2), opened["total_pages"])
				assert.Equal(t, float64(2), inspected["page"])
				assert.Equal(t, float64(2), inspected["total_pages"])
				assert.Equal(t, float64(26), inspected["item_start"])
				assert.Equal(t, float64(30), inspected["item_end"])
				assert.Equal(t, float64(30), inspected["total_items"])
				assert.Contains(t, opened["projection"], "Architecture")
				assert.NotContains(t, opened["projection"], "navigation needle")
				assert.Contains(t, inspected["projection"], "Topic 26")
				assert.NotContains(t, inspected["projection"], "Architecture")
			}))
		}),
	)
}
