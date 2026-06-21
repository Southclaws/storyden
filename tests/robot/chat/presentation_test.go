package chat_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/semdex/robot/presentation"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/app/transports/sse"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
	robot_tests "github.com/Southclaws/storyden/tests/robot"
)

func TestRobotPresentationMarkupRenderCard(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{
			LanguageModelProvider: "mock",
		},
		e2e.Setup(),
		robot_tests.WithRobotSettings(mockModelAck),
		sse.Build(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
			sessionRepo *robot_session.Repository,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				pageSlug := "robot-presentation-page-" + xid.New().String()
				vis := openapi.VisibilityPublished
				page := tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeCreateJSONRequestBody{
					Name:        "Robot Presentation Page",
					Slug:        &pageSlug,
					Description: ptr("A page rendered as a Robot presentation card."),
					Visibility:  &vis,
				}, adminSession))(t, http.StatusOK)
				require.NotNil(t, page.JSON200)
				pageID := string(page.JSON200.Id)
				missingPageID := xid.New().String()

				scriptName := "robot-chat-presentation-" + xid.New().String() + ".yaml"
				scriptPath := filepath.Join("..", "scripts", scriptName)
				writeScript(t, scriptPath, `steps:
  - match:
      contains: "show page card"
    respond:
      text: |
        Here's the page I found:

        [Robot Presentation Page](sdr:node/`+pageID+`)

        Let me know if you'd like a summary.
      finish: "stop"
  - match:
      contains: "show invalid card"
    respond:
      text: |
        Here's a missing page:

        [Missing page](sdr:node/`+missingPageID+`)
      finish: "stop"
  - match:
      contains: "next message"
    respond:
      text: "Next message handled normally."
      finish: "stop"
`)
				defer os.Remove(scriptPath)

				actor := tests.AssertRequest(cl.RobotCreateWithResponse(root,
					openapi.RobotCreateJSONRequestBody{
						Name:        "presentation-actor-" + xid.New().String(),
						Description: "robot that renders presentation markup",
						Playbook:    "you render Library page cards with presentation markup",
						Model:       robotModelPtr("mock/../scripts/" + scriptName),
					},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, actor.JSON200)

				sessionID := xid.New().String()
				first := doChat(t, root, ts, adminSession, sessionID, string(actor.JSON200.Id), "show page card")
				assert.Empty(t, collectToolInputs(first))
				assert.Empty(t, collectToolOutputs(first))
				assert.Equal(t,
					"Here's the page I found:\n\n\n\nLet me know if you'd like a summary.",
					strings.TrimSuffix(strings.Join(collectTextDeltas(first), ""), "\n"),
				)

				streamCards := collectDataParts(first, presentation.DataRenderCard)
				require.Len(t, streamCards, 1)
				assert.Equal(t, "sdr:node/"+pageID, streamCards[0]["ref"])
				assert.Equal(t, "node", streamCards[0]["kind"])
				assert.Equal(t, pageID, streamCards[0]["id"])
				assert.NotContains(t, streamCards[0], "status")
				assert.NotContains(t, streamCards[0], "page")

				parsedSessionID, err := xid.FromString(sessionID)
				require.NoError(t, err)

				sess, _, err := sessionRepo.Get(root, robot.SessionID(parsedSessionID), robot.NewMessageCursorParams(opt.NewEmpty[robot.MessageID](), 50))
				require.NoError(t, err)
				var rawAssistantText string
				var functionParts int
				for _, message := range sess.Messages {
					if message.Event.LLMResponse.Content == nil {
						continue
					}
					for _, part := range message.Event.LLMResponse.Content.Parts {
						if part == nil {
							continue
						}
						if part.Text != "" && strings.Contains(part.Text, "sdr:node/"+pageID) {
							rawAssistantText = part.Text
						}
						if part.FunctionCall != nil || part.FunctionResponse != nil {
							functionParts++
						}
					}
				}
				assert.Contains(t, rawAssistantText, `[Robot Presentation Page](sdr:node/`+pageID+`)`)
				assert.Zero(t, functionParts)

				sessionResponse := tests.AssertRequest(cl.RobotSessionGetWithResponse(root,
					openapi.RobotSessionIDParam(sessionID),
					&openapi.RobotSessionGetParams{},
					adminSession,
				))(t, http.StatusOK)
				require.NotNil(t, sessionResponse.JSON200)

				hydratedCards := collectSessionDataParts(sessionResponse.JSON200.MessageList.Messages, presentation.DataRenderCard)
				require.Len(t, hydratedCards, 1)
				assert.Equal(t, "sdr:node/"+pageID, hydratedCards[0]["ref"])
				assert.Equal(t, "node", hydratedCards[0]["kind"])
				assert.Equal(t, pageID, hydratedCards[0]["id"])
				assert.NotContains(t, hydratedCards[0], "status")
				assert.NotContains(t, hydratedCards[0], "page")

				second := doChat(t, root, ts, adminSession, sessionID, string(actor.JSON200.Id), "next message")
				assert.Equal(t, "Next message handled normally.", strings.Join(collectTextDeltas(second), ""))

				invalidSessionID := xid.New().String()
				invalid := doChat(t, root, ts, adminSession, invalidSessionID, string(actor.JSON200.Id), "show invalid card")
				invalidCards := collectDataParts(invalid, presentation.DataRenderCard)
				require.Len(t, invalidCards, 1)
				assert.Equal(t, "sdr:node/"+missingPageID, invalidCards[0]["ref"])
				assert.Equal(t, "node", invalidCards[0]["kind"])
				assert.Equal(t, missingPageID, invalidCards[0]["id"])
			}))
		}),
	)
}

func collectDataParts(ev *fullResponse, partType string) []map[string]any {
	var result []map[string]any
	for _, part := range ev.parts {
		if part.Type != partType {
			continue
		}
		dataPart, err := part.AsDataPart()
		if err != nil {
			continue
		}
		if data, ok := dataPart.Data.(map[string]any); ok {
			result = append(result, data)
		}
	}
	return result
}

func collectSessionDataParts(messages []openapi.RobotSessionMessage, partType string) []map[string]any {
	var result []map[string]any
	for _, message := range messages {
		for _, part := range message.Parts {
			if string(part.Type) != partType {
				continue
			}
			dataPart, err := part.AsDataPart()
			if err != nil {
				continue
			}
			if data, ok := dataPart.Data.(map[string]any); ok {
				result = append(result, data)
			}
		}
	}
	return result
}

func ptr[T any](v T) *T {
	return &v
}
