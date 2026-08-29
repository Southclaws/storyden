package chat_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	robot_resource "github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/robot_memory"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry/denbot"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
	robottest "github.com/Southclaws/storyden/tests/robot"
)

const (
	managementParentMemoryID  = "cto7n8ifunp55p1bujv0"
	managementTargetMemoryID  = "cto7nm2funp55p1bujvg"
	managementArchiveMemoryID = "d8818ueot5pfij6bvm90"
)

func TestRobotMemoryAutonomouslyCreatesDurableFact(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{LanguageModelProvider: "mock"},
		e2e.Setup(),
		robottest.WithRobotSettings(mockModelMemoryCreate),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
			memoryRepo *robot_memory.Repository,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)
				created := tests.AssertRequest(cl.RobotCreateWithResponse(root, openapi.RobotCreateJSONRequestBody{
					Name:        "memory-robot-" + xid.New().String(),
					Description: "robot for persistent memory tests",
					Playbook:    "Store durable information in memory.",
					Toolsets:    robotToolsetsPtr("system.memory"),
				}, adminSession))(t, http.StatusOK)
				robotRef := string(created.JSON200.Id)

				stream := doChat(t, root, ts, adminSession, xid.New().String(), robotRef, "Freyja is owned by Southclaws.")

				assert.Equal(t, []string{"memory_create"}, collectToolCalls(stream))
				assert.Equal(t, "Memory stored.", strings.Join(collectTextDeltas(stream), ""))

				results, _, err := memoryRepo.Search(root, robotRef, robot_memory.SearchFilter{
					Subject: opt.New("freyja"), Predicate: opt.New("owned_by"), Object: opt.New("southclaws"),
				}, opt.NewEmpty[robot_resource.MemoryID](), robot_memory.SearchLimit)
				require.NoError(t, err)
				require.Len(t, results, 1)
				assert.Equal(t, "Freyja is owned by Southclaws.", results[0].Memory.Content)
			}))
		}),
	)
}

func TestRobotMemorySavesAndRecallsAcrossSessions(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{LanguageModelProvider: "mock"},
		e2e.Setup(),
		robottest.WithRobotSettings(mockModelMemoryWorkflow),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
			db *ent.Client,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)
				created := tests.AssertRequest(cl.RobotCreateWithResponse(root, openapi.RobotCreateJSONRequestBody{
					Name:        "memory-workflow-robot-" + xid.New().String(),
					Description: "robot for memory workflow tests",
					Playbook:    "Use memory when lasting personal context would help.",
					Toolsets:    robotToolsetsPtr("system.memory"),
				}, adminSession))(t, http.StatusOK)
				robotRef := string(created.JSON200.Id)

				saved := doChat(t, root, ts, adminSession, xid.New().String(), robotRef, "Freyja is owned by Southclaws.")
				assert.Equal(t, []string{"memory_create"}, collectToolCalls(saved))
				assert.Equal(t, "I'll remember that about Freyja.", strings.Join(collectTextDeltas(saved), ""))
				createOutput := toolOutputByCallID(t, collectToolOutputs(saved), "call_memory_workflow_create")
				memoryID, ok := createOutput["memory_id"].(string)
				require.True(t, ok)
				require.NotEmpty(t, memoryID)
				assert.Equal(t, "freyja", createOutput["subject"])
				assert.Equal(t, "owned_by", createOutput["predicate"])
				assert.Equal(t, "southclaws", createOutput["object"])

				recalled := doChat(t, root, ts, adminSession, xid.New().String(), robotRef, "Who owns Freyja?")
				assert.Equal(t, []string{"memory_search"}, collectToolCalls(recalled))
				assert.Equal(t, "Freyja is owned by Southclaws.", strings.Join(collectTextDeltas(recalled), ""))
				searchOutput := toolOutputByCallID(t, collectToolOutputs(recalled), "call_memory_workflow_search")
				assert.Equal(t, float64(1), searchOutput["returned"])
				assert.Equal(t, false, searchOutput["has_more"])
				results, ok := searchOutput["results"].([]any)
				require.True(t, ok)
				require.Len(t, results, 1)
				result, ok := results[0].(map[string]any)
				require.True(t, ok)
				assert.Equal(t, memoryID, result["memory_id"])
				assert.Equal(t, "Freyja is owned by Southclaws.", result["excerpt"])
				assert.Equal(t, "freyja", result["subject"])
				assert.Equal(t, "owned_by", result["predicate"])
				assert.Equal(t, "southclaws", result["object"])

				parsedMemoryID, err := xid.FromString(memoryID)
				require.NoError(t, err)
				stored, err := db.RobotMemory.Get(root, parsedMemoryID)
				require.NoError(t, err)
				assert.Equal(t, robotRef, stored.RobotRef)
				assert.Equal(t, "Freyja is owned by Southclaws.", stored.Content)
				assert.Equal(t, uint64(1), stored.AccessCount)
			}))
		}),
	)
}

func TestDenbotDiscoversAndUsesMemoryManagement(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{LanguageModelProvider: "mock"},
		e2e.Setup(),
		robottest.WithRobotSettings(mockModelMemoryManagement),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
			db *ent.Client,
		) {
			lc.Append(fx.StartHook(func() {
				parentID, err := xid.FromString(managementParentMemoryID)
				require.NoError(t, err)
				targetID, err := xid.FromString(managementTargetMemoryID)
				require.NoError(t, err)

				_, err = db.RobotMemory.Create().
					SetID(parentID).
					SetRobotRef(denbot.ID).
					SetContent("People and pets.").
					Save(root)
				require.NoError(t, err)
				_, err = db.RobotMemory.Create().
					SetID(targetID).
					SetRobotRef(denbot.ID).
					SetContent("Freyja is owned by Southclaws.").
					SetSubject("freyja").
					SetPredicate("owned_by").
					SetObject("southclaws").
					Save(root)
				require.NoError(t, err)

				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)
				stream := doChat(t, root, ts, adminSession, xid.New().String(), "", "Correct Freyja's owner to Barney and organize that memory under people and pets.")

				assert.Equal(t, []string{
					"toolset_search",
					"toolset_get",
					"toolset_load",
					"memory_list",
					"memory_open",
					"memory_update",
					"memory_move",
				}, collectToolCalls(stream))
				assert.Equal(t, "Freyja's memory was corrected and organized.", strings.Join(collectTextDeltas(stream), ""))

				outputs := collectToolOutputs(stream)
				searchOutput := toolOutputByCallID(t, outputs, "call_memory_management_search")
				toolsets, ok := searchOutput["toolsets"].([]any)
				require.True(t, ok)
				management := resultByStringField(t, toolsets, "id", "system.memory_management")
				assert.Equal(t, "Knowledge graph management", management["name"])

				getOutput := toolOutputByCallID(t, outputs, "call_memory_management_get")
				assert.Equal(t, "system.memory_management", getOutput["id"])
				assert.ElementsMatch(t, []any{"memory_list", "memory_move", "memory_open", "memory_update"}, getOutput["tools"])

				listOutput := toolOutputByCallID(t, outputs, "call_memory_management_list")
				assert.Equal(t, float64(2), listOutput["returned"])
				assert.Equal(t, false, listOutput["has_more"])
				memories, ok := listOutput["memories"].([]any)
				require.True(t, ok)
				listedTarget := resultByStringField(t, memories, "id", managementTargetMemoryID)
				assert.Equal(t, "Freyja is owned by Southclaws.", listedTarget["excerpt"])
				assert.Equal(t, "southclaws", listedTarget["object"])

				openOutput := toolOutputByCallID(t, outputs, "call_memory_management_open")
				opened, ok := openOutput["memory"].(map[string]any)
				require.True(t, ok)
				assert.Equal(t, managementTargetMemoryID, opened["id"])
				assert.Equal(t, "Freyja is owned by Southclaws.", opened["content"])

				updateOutput := toolOutputByCallID(t, outputs, "call_memory_management_update")
				updated, ok := updateOutput["memory"].(map[string]any)
				require.True(t, ok)
				assert.Equal(t, "Freyja is owned by Barney.", updated["content"])
				assert.Equal(t, "barney", updated["object"])

				moveOutput := toolOutputByCallID(t, outputs, "call_memory_management_move")
				moved, ok := moveOutput["memory"].(map[string]any)
				require.True(t, ok)
				assert.Equal(t, managementParentMemoryID, moved["parent_id"])

				stored, err := db.RobotMemory.Get(root, targetID)
				require.NoError(t, err)
				require.NotNil(t, stored.ParentID)
				assert.Equal(t, parentID, *stored.ParentID)
				assert.Equal(t, "Freyja is owned by Barney.", stored.Content)
				require.NotNil(t, stored.Subject)
				require.NotNil(t, stored.Predicate)
				require.NotNil(t, stored.Object)
				assert.Equal(t, "freyja", *stored.Subject)
				assert.Equal(t, "owned_by", *stored.Predicate)
				assert.Equal(t, "barney", *stored.Object)
				assert.Equal(t, uint64(1), stored.AccessCount)
			}))
		}),
	)
}

func TestCustomRobotArchivesStaleMemory(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{LanguageModelProvider: "mock"},
		e2e.Setup(),
		robottest.WithRobotSettings(mockModelMemoryArchive),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			ts *httptest.Server,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
			memoryRepo *robot_memory.Repository,
			db *ent.Client,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)
				created := tests.AssertRequest(cl.RobotCreateWithResponse(root, openapi.RobotCreateJSONRequestBody{
					Name:        "memory-archive-robot-" + xid.New().String(),
					Description: "robot for memory archive tests",
					Playbook:    "Recall useful facts and archive stale evidence when asked.",
					Toolsets:    robotToolsetsPtr("system.memory", "system.memory_management"),
				}, adminSession))(t, http.StatusOK)
				robotRef := string(created.JSON200.Id)

				memoryID, err := xid.FromString(managementArchiveMemoryID)
				require.NoError(t, err)
				_, err = db.RobotMemory.Create().
					SetID(memoryID).
					SetRobotRef(robotRef).
					SetContent("Freyja is owned by Southclaws.").
					SetSubject("freyja").
					SetPredicate("owned_by").
					SetObject("southclaws").
					Save(root)
				require.NoError(t, err)

				stream := doChat(t, root, ts, adminSession, xid.New().String(), robotRef, "Archive Freyja's stale owner memory.")
				assert.Equal(t, []string{"memory_search", "memory_open", "memory_update"}, collectToolCalls(stream))
				assert.Equal(t, "The stale memory was archived.", strings.Join(collectTextDeltas(stream), ""))

				outputs := collectToolOutputs(stream)
				searchOutput := toolOutputByCallID(t, outputs, "call_memory_archive_search")
				assert.Equal(t, float64(1), searchOutput["returned"])
				openOutput := toolOutputByCallID(t, outputs, "call_memory_archive_open")
				opened, ok := openOutput["memory"].(map[string]any)
				require.True(t, ok)
				assert.Equal(t, managementArchiveMemoryID, opened["id"])
				assert.Equal(t, "active", opened["state"])
				updateOutput := toolOutputByCallID(t, outputs, "call_memory_archive_update")
				updated, ok := updateOutput["memory"].(map[string]any)
				require.True(t, ok)
				assert.Equal(t, "archived", updated["state"])

				stored, err := db.RobotMemory.Get(root, memoryID)
				require.NoError(t, err)
				assert.Equal(t, "archived", string(stored.State))
				assert.Equal(t, uint64(2), stored.AccessCount)

				active, hasMore, err := memoryRepo.Search(root, robotRef, robot_memory.SearchFilter{
					Subject: opt.New("freyja"), Predicate: opt.New("owned_by"), Object: opt.New("southclaws"),
				}, opt.NewEmpty[robot_resource.MemoryID](), robot_memory.SearchLimit)
				require.NoError(t, err)
				assert.Empty(t, active)
				assert.False(t, hasMore)
			}))
		}),
	)
}

func resultByStringField(t *testing.T, values []any, field, expected string) map[string]any {
	t.Helper()
	for _, value := range values {
		result, ok := value.(map[string]any)
		if !ok {
			continue
		}
		if result[field] == expected {
			return result
		}
	}
	require.FailNow(t, "result not found", "field %q did not contain %q", field, expected)
	return nil
}
