package crud_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

const testModel = "anthropic/claude-3-sonnet"

func robotModel(s string) *openapi.RobotModelRef {
	m := openapi.RobotModelRef(s)
	return &m
}

func robotTools(names ...string) *openapi.RobotToolNameList {
	list := openapi.RobotToolNameList(names)
	return &list
}

func strPtr(s string) *string { return &s }

func TestRobotCRUD(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{},
		e2e.Setup(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, adminAcc := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				t.Run("create", func(t *testing.T) {
					a := assert.New(t)

					name := "helper-bot-" + uuid.NewString()
					create := tests.AssertRequest(cl.RobotCreateWithResponse(root,
						openapi.RobotCreateJSONRequestBody{
							Name:        name,
							Description: "Helps with testing",
							Playbook:    "You are a helpful test assistant.",
							Model:       robotModel(testModel),
						},
						adminSession,
					))(t, http.StatusOK)

					a.Equal(name, create.JSON200.Name)
					a.Equal("Helps with testing", create.JSON200.Description)
					a.Equal("You are a helpful test assistant.", create.JSON200.Playbook)
					a.Equal(testModel, create.JSON200.Model)
					a.Equal(adminAcc.ID.String(), create.JSON200.Author.Id)
					a.NotEmpty(create.JSON200.Id)
				})

				t.Run("create_with_tools", func(t *testing.T) {
					a := assert.New(t)

					create := tests.AssertRequest(cl.RobotCreateWithResponse(root,
						openapi.RobotCreateJSONRequestBody{
							Name:        "tool-bot-" + uuid.NewString(),
							Description: "Bot with tools",
							Playbook:    "Use your tools.",
							Model:       robotModel(testModel),
							Tools:       robotTools("library_search_pages", "thread_search"),
						},
						adminSession,
					))(t, http.StatusOK)

					a.ElementsMatch([]string{"library_search_pages", "thread_search"}, create.JSON200.Tools)
				})

				t.Run("create_with_unknown_tool_returns_400", func(t *testing.T) {
					r := require.New(t)

					create, err := cl.RobotCreateWithResponse(root,
						openapi.RobotCreateJSONRequestBody{
							Name:        "unknown-tool-bot-" + uuid.NewString(),
							Description: "Bot with unknown tool",
							Playbook:    "Use your tools.",
							Model:       robotModel(testModel),
							Tools:       robotTools("unknown_tool"),
						},
						adminSession,
					)
					r.NoError(err)
					r.Equal(http.StatusBadRequest, create.StatusCode())
				})

				t.Run("list_includes_created", func(t *testing.T) {
					a := assert.New(t)

					name := "list-test-bot-" + uuid.NewString()
					created := tests.AssertRequest(cl.RobotCreateWithResponse(root,
						openapi.RobotCreateJSONRequestBody{
							Name:        name,
							Description: "Created for list test",
							Playbook:    "List test playbook.",
							Model:       robotModel(testModel),
						},
						adminSession,
					))(t, http.StatusOK)

					list := tests.AssertRequest(cl.RobotsListWithResponse(root,
						&openapi.RobotsListParams{},
						adminSession,
					))(t, http.StatusOK)

					found := false
					for _, r := range list.JSON200.Robots {
						if r.Id == created.JSON200.Id {
							found = true
							a.Equal(name, r.Name)
						}
					}
					a.True(found, "created robot must appear in list")
				})

				t.Run("get_by_id", func(t *testing.T) {
					a := assert.New(t)

					name := "get-test-bot-" + uuid.NewString()
					created := tests.AssertRequest(cl.RobotCreateWithResponse(root,
						openapi.RobotCreateJSONRequestBody{
							Name:        name,
							Description: "Created for get test",
							Playbook:    "Get test playbook.",
							Model:       robotModel(testModel),
						},
						adminSession,
					))(t, http.StatusOK)

					get := tests.AssertRequest(cl.RobotGetWithResponse(root,
						created.JSON200.Id,
						adminSession,
					))(t, http.StatusOK)

					a.Equal(created.JSON200.Id, get.JSON200.Id)
					a.Equal(name, get.JSON200.Name)
					a.Equal("Created for get test", get.JSON200.Description)
					a.Equal("Get test playbook.", get.JSON200.Playbook)
					a.Equal(testModel, get.JSON200.Model)
					a.Equal(adminAcc.ID.String(), get.JSON200.Author.Id)
				})

				t.Run("get_nonexistent_returns_404", func(t *testing.T) {
					r := require.New(t)

					get, err := cl.RobotGetWithResponse(root, xid.New().String(), adminSession)
					r.NoError(err)
					r.Equal(http.StatusNotFound, get.StatusCode())
				})

				t.Run("update_name_description_playbook", func(t *testing.T) {
					a := assert.New(t)

					created := tests.AssertRequest(cl.RobotCreateWithResponse(root,
						openapi.RobotCreateJSONRequestBody{
							Name:        "pre-update-bot-" + uuid.NewString(),
							Description: "Original description",
							Playbook:    "Original playbook.",
							Model:       robotModel(testModel),
						},
						adminSession,
					))(t, http.StatusOK)

					newName := "post-update-bot-" + uuid.NewString()
					newDesc := "Updated description"
					newPlaybook := "Updated playbook instructions."

					updated := tests.AssertRequest(cl.RobotUpdateWithResponse(root,
						created.JSON200.Id,
						openapi.RobotUpdateJSONRequestBody{
							Name:        &newName,
							Description: &newDesc,
							Playbook:    &newPlaybook,
						},
						adminSession,
					))(t, http.StatusOK)

					a.Equal(newName, updated.JSON200.Name)
					a.Equal(newDesc, updated.JSON200.Description)
					a.Equal(newPlaybook, updated.JSON200.Playbook)
					a.Equal(created.JSON200.Model, updated.JSON200.Model)

					// Verify persistence via get
					get := tests.AssertRequest(cl.RobotGetWithResponse(root,
						created.JSON200.Id,
						adminSession,
					))(t, http.StatusOK)
					a.Equal(newName, get.JSON200.Name)
					a.Equal(newDesc, get.JSON200.Description)
					a.Equal(newPlaybook, get.JSON200.Playbook)
				})

				t.Run("update_tools", func(t *testing.T) {
					a := assert.New(t)

					created := tests.AssertRequest(cl.RobotCreateWithResponse(root,
						openapi.RobotCreateJSONRequestBody{
							Name:        "tool-update-bot-" + uuid.NewString(),
							Description: "Tool update test",
							Playbook:    "Playbook.",
							Model:       robotModel(testModel),
							Tools:       robotTools("library_search_pages"),
						},
						adminSession,
					))(t, http.StatusOK)
					a.ElementsMatch([]string{"library_search_pages"}, created.JSON200.Tools)

					updated := tests.AssertRequest(cl.RobotUpdateWithResponse(root,
						created.JSON200.Id,
						openapi.RobotUpdateJSONRequestBody{
							Tools: robotTools("thread_search", "member_search"),
						},
						adminSession,
					))(t, http.StatusOK)

					a.ElementsMatch([]string{"thread_search", "member_search"}, updated.JSON200.Tools)
				})

				t.Run("update_with_unknown_tool_returns_400", func(t *testing.T) {
					r := require.New(t)

					created := tests.AssertRequest(cl.RobotCreateWithResponse(root,
						openapi.RobotCreateJSONRequestBody{
							Name:        "unknown-tool-update-bot-" + uuid.NewString(),
							Description: "Tool update test",
							Playbook:    "Playbook.",
							Model:       robotModel(testModel),
						},
						adminSession,
					))(t, http.StatusOK)

					updated, err := cl.RobotUpdateWithResponse(root,
						created.JSON200.Id,
						openapi.RobotUpdateJSONRequestBody{
							Tools: robotTools("unknown_tool"),
						},
						adminSession,
					)
					r.NoError(err)
					r.Equal(http.StatusBadRequest, updated.StatusCode())
				})

				t.Run("update_nonexistent_returns_404", func(t *testing.T) {
					r := require.New(t)

					update, err := cl.RobotUpdateWithResponse(root,
						xid.New().String(),
						openapi.RobotUpdateJSONRequestBody{Name: strPtr("new name")},
						adminSession,
					)
					r.NoError(err)
					r.Equal(http.StatusNotFound, update.StatusCode())
				})
			}))
		}),
	)
}
