package tools_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/link/scrape"
	robottools "github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/nodeversion"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/lib/mcp"
)

func TestLibraryBatchToolsCreateHierarchyAtomically(t *testing.T) {
	integration.Test(t, nil, fx.Invoke(func(lc fx.Lifecycle, ctx context.Context, registry *robottools.Registry, aw *account_writer.Writer, reader *node_querier.Querier, db *ent.Client) {
		lc.Append(fx.StartHook(func() {
			ctx, account := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			ctx = session.WithAccountPermissions(ctx, *account, rbac.NewList(rbac.PermissionManageLibrary))
			create, err := registry.GetTool(ctx, "library_pages_create")
			require.NoError(t, err)
			result := callLibraryTool(t, create, ctx, `{"items":[{"ref":"leaf","name":"Resource","parent_ref":"category","tags":["new-topic"]},{"ref":"category","name":"Category","parent_ref":"root"},{"ref":"root","name":"Collection"}]}`)
			require.Len(t, result.Results, 3)
			for _, item := range result.Results {
				assert.Equal(t, mcp.LibraryPageBatchResultYamlStatusCreated, item.Status)
				require.NotEmpty(t, item.Id)
				require.NotEmpty(t, item.BrowserUrl)
			}
			leaf, err := reader.Get(ctx, library.NewKey(result.Results[0].Id))
			require.NoError(t, err)
			parent, ok := leaf.Parent.Get()
			require.True(t, ok)
			assert.Equal(t, result.Results[1].Id, parent.GetID().String())
			require.Len(t, leaf.Tags, 1)
			assert.Equal(t, "new-topic", leaf.Tags[0].Name.String())
			category, err := reader.Get(ctx, library.NewKey(result.Results[1].Id))
			require.NoError(t, err)
			root, ok := category.Parent.Get()
			require.True(t, ok)
			assert.Equal(t, result.Results[2].Id, root.GetID().String())

			_, err = create.Handler(ctx, json.RawMessage(fmt.Sprintf(`{"items":[{"ref":"child","name":"Blocked child","parent_ref":"conflict"},{"ref":"conflict","name":"Conflict","slug":%q},{"ref":"independent","name":"Independent"}]}`, result.Results[2].Slug)))
			require.ErrorContains(t, err, "already exists")

			before, err := db.Node.Query().Count(ctx)
			require.NoError(t, err)
			_, err = create.Handler(ctx, json.RawMessage(`{"items":[{"ref":"valid","name":"Must not be written"},{"ref":"bad","name":"Missing parent","parent_ref":"absent"}]}`))
			require.Error(t, err)
			after, err := db.Node.Query().Count(ctx)
			require.NoError(t, err)
			assert.Equal(t, before, after)
			denied := session.WithAccountPermissions(ctx, *account, rbac.Permissions{})
			_, err = create.Handler(denied, json.RawMessage(`{"items":[{"ref":"no","name":"Denied"}]}`))
			require.Error(t, err)
		}))
	}))
}

func TestLibraryToolsUpdateTagsWithoutLosingOtherFields(t *testing.T) {
	integration.Test(t, nil, fx.Invoke(func(lc fx.Lifecycle, ctx context.Context, registry *robottools.Registry, aw *account_writer.Writer, reader *node_querier.Querier) {
		lc.Append(fx.StartHook(func() {
			ctx, account := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			ctx = session.WithAccountPermissions(ctx, *account, rbac.NewList(rbac.PermissionManageLibrary))
			create, err := registry.GetTool(ctx, "library_pages_create")
			require.NoError(t, err)
			created := callLibraryTool(t, create, ctx, `{"items":[{"ref":"one","name":"Keep this name","content":"<p>Keep this content</p>","tags":["original"]},{"ref":"two","name":"Other page"}]}`)
			id := created.Results[0].Id
			update, err := registry.GetTool(ctx, "library_pages_update")
			require.NoError(t, err)
			result := callLibraryTool(t, update, ctx, fmt.Sprintf(`{"items":[{"ref":"one","id":%q,"tags":["new-tag"]},{"ref":"two","id":%q,"name":"Renamed"}]}`, id, created.Results[1].Id))
			assert.Equal(t, mcp.LibraryPageBatchResultYamlStatusUpdated, result.Results[0].Status)
			page, err := reader.Get(ctx, library.NewKey(id))
			require.NoError(t, err)
			assert.Equal(t, "Keep this name", page.Name)
			assert.Contains(t, page.Content.OrZero().Plaintext(), "Keep this content")
			require.Len(t, page.Tags, 1)
			assert.Equal(t, "new-tag", page.Tags[0].Name.String())
			callLibraryTool(t, update, ctx, fmt.Sprintf(`{"items":[{"ref":"one","id":%q,"name":"Still tagged"}]}`, id))
			page, err = reader.Get(ctx, library.NewKey(id))
			require.NoError(t, err)
			require.Len(t, page.Tags, 1)
			callLibraryTool(t, update, ctx, fmt.Sprintf(`{"items":[{"ref":"one","id":%q,"tags":[]}]}`, id))
			page, err = reader.Get(ctx, library.NewKey(id))
			require.NoError(t, err)
			assert.Empty(t, page.Tags)
			single, err := registry.GetTool(ctx, "update_library_page")
			require.NoError(t, err)
			_, err = single.Handler(ctx, json.RawMessage(fmt.Sprintf(`{"id":%q,"tags":["single"]}`, id)))
			require.NoError(t, err)
			_, err = single.Handler(ctx, json.RawMessage(fmt.Sprintf(`{"id":%q,"tags":[]}`, id)))
			require.NoError(t, err)
			page, err = reader.Get(ctx, library.NewKey(id))
			require.NoError(t, err)
			assert.Empty(t, page.Tags)
			_, err = update.Handler(ctx, json.RawMessage(fmt.Sprintf(`{"items":[{"ref":"a","id":%q,"name":"No mutation"},{"ref":"b","id":%q,"name":"Duplicate target"}]}`, id, page.GetSlug())))
			require.ErrorContains(t, err, "multiple updates")
			_, err = update.Handler(ctx, json.RawMessage(fmt.Sprintf(`{"items":[{"ref":"conflict","id":%q,"slug":%q},{"ref":"other","id":%q,"content":"<p>Must not update</p>"}]}`, id, created.Results[1].Slug, created.Results[1].Id)))
			require.ErrorContains(t, err, "already exists")

			unchanged, err := reader.Get(ctx, library.NewKey(id))
			require.NoError(t, err)
			assert.Equal(t, page.GetSlug(), unchanged.GetSlug())
			independent, err := reader.Get(ctx, library.NewKey(created.Results[1].Id))
			require.NoError(t, err)
			assert.False(t, independent.Content.Ok())
		}))
	}))
}

func callLibraryTool(t *testing.T, tool *robottools.Tool, ctx context.Context, input string) mcp.ToolLibraryPagesCreateOutput {
	t.Helper()
	raw, err := tool.Handler(ctx, json.RawMessage(input))
	require.NoError(t, err)
	var result mcp.ToolLibraryPagesCreateOutput
	require.NoError(t, json.Unmarshal(raw, &result))
	return result
}

func TestLibraryURLFallbackAndStorageFailure(t *testing.T) {
	integration.Test(t, nil, fx.Decorate(func(scrape.Scraper) scrape.Scraper { return unavailablePageScraper{} }), fx.Invoke(func(lc fx.Lifecycle, ctx context.Context, registry *robottools.Registry, aw *account_writer.Writer, reader *node_querier.Querier, db *ent.Client) {
		lc.Append(fx.StartHook(func() {
			ctx, account := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			ctx = session.WithAccountPermissions(ctx, *account, rbac.NewList(rbac.PermissionManageLibrary))
			create, err := registry.GetTool(ctx, "library_pages_create")
			require.NoError(t, err)
			created := callLibraryTool(t, create, ctx, `{"items":[{"ref":"fallback","name":"Unavailable source","url":"https://example.com/unavailable"}]}`)
			require.Equal(t, mcp.LibraryPageBatchResultYamlStatusCreated, created.Results[0].Status)
			page, err := reader.Get(ctx, library.NewKey(created.Results[0].Id))
			require.NoError(t, err)
			link, ok := page.WebLink.Get()
			require.True(t, ok)
			assert.Equal(t, "https://example.com/unavailable", link.URL)
			db.Link.Use(func(ent.Mutator) ent.Mutator {
				return ent.MutateFunc(func(context.Context, ent.Mutation) (ent.Value, error) {
					return nil, errors.New("injected link storage failure")
				})
			})
			withoutLink := callLibraryTool(t, create, ctx, `{"items":[{"ref":"failure","name":"Page without link","url":"https://example.com/failed"}]}`)
			require.Equal(t, mcp.LibraryPageBatchResultYamlStatusCreated, withoutLink.Results[0].Status)
			page, err = reader.Get(ctx, library.NewKey(withoutLink.Results[0].Id))
			require.NoError(t, err)
			assert.Equal(t, "Page without link", page.Name)
			assert.False(t, page.WebLink.Ok())
		}))
	}))
}

func TestLibraryBatchHierarchyAndDraftGuards(t *testing.T) {
	integration.Test(t, nil, fx.Invoke(func(lc fx.Lifecycle, ctx context.Context, registry *robottools.Registry, aw *account_writer.Writer, reader *node_querier.Querier, db *ent.Client) {
		lc.Append(fx.StartHook(func() {
			ctx, account := e2e.WithAccount(ctx, aw, seed.Account_001_Odin)
			ctx = session.WithAccountPermissions(ctx, *account, rbac.NewList(rbac.PermissionManageLibrary))
			create, err := registry.GetTool(ctx, "library_pages_create")
			require.NoError(t, err)

			update, err := registry.GetTool(ctx, "library_pages_update")
			require.NoError(t, err)

			for _, input := range []string{
				`{"items":[{"ref":"self","name":"Self","parent_ref":"self"}]}`,
				`{"items":[{"ref":"a","name":"A","parent_ref":"b"},{"ref":"b","name":"B","parent_ref":"a"}]}`,
				`{"items":[{"ref":"same","name":"A"},{"ref":"same","name":"B"}]}`,
				`{"items":[{"ref":"a","name":"Same"},{"ref":"b","name":"Same"}]}`,
				`{"items":[{"ref":"a","name":"A","parent":"missing-parent"}]}`,
				`{"items":[{"ref":"a","name":"A","url":"file:///etc/passwd"}]}`,
			} {
				_, err := create.Handler(ctx, json.RawMessage(input))
				require.Error(t, err, input)
			}

			created := callLibraryTool(t, create, ctx, `{"items":[{"ref":"root","name":"Root"},{"ref":"child","name":"Child","parent_ref":"root","content":"<p data-block-id=\"kept-block\">Keep this paragraph</p>"}]}`)
			root, child := created.Results[0], created.Results[1]

			_, err = update.Handler(ctx, json.RawMessage(fmt.Sprintf(`{"items":[{"ref":"root","id":%q,"parent":%q},{"ref":"child","id":%q,"name":"Must not change"}]}`, root.Id, child.Id, child.Id)))
			require.ErrorContains(t, err, "cycle")

			page, err := reader.Get(ctx, library.NewKey(child.Id))
			require.NoError(t, err)
			assert.Equal(t, "Child", page.Name)

			version := db.NodeVersion.Create().
				SetNodeID(page.GetID()).
				SetAuthorID(page.GetAuthor()).
				SetName("Draft").
				SetSlug("draft").
				SaveX(ctx)

			_, err = update.Handler(ctx, json.RawMessage(fmt.Sprintf(`{"items":[{"ref":"root","id":%q,"name":"Must not change"},{"ref":"child","id":%q,"content":"<p>Direct edit</p>"}]}`, root.Id, child.Id)))
			require.ErrorContains(t, err, "working draft")

			callLibraryTool(t, update, ctx, fmt.Sprintf(`{"items":[{"ref":"child","id":%q,"tags":["draft-tag"]}]}`, child.Id))
			page, err = reader.Get(ctx, library.NewKey(child.Id))
			require.NoError(t, err)
			require.Len(t, page.Tags, 1)

			db.NodeVersion.UpdateOneID(version.ID).SetStatus(nodeversion.StatusApplied).ExecX(ctx)
			db.Node.UpdateOneID(page.GetID()).SetCurrentVersionID(version.ID).ExecX(ctx)
			beforeContent := page.Content.OrZero().HTML()
			callLibraryTool(t, update, ctx, fmt.Sprintf(`{"items":[{"ref":"child","id":%q,"content":%q,"slug":"child-renamed"}]}`, child.Id, beforeContent))
			page, err = reader.Get(ctx, library.NewKey(child.Id))
			require.NoError(t, err)
			assert.Equal(t, beforeContent, page.Content.OrZero().HTML())
			assert.False(t, page.CurrentVersion.Ok())
			assert.Equal(t, "child-renamed", page.GetSlug())

			callLibraryTool(t, update, ctx, fmt.Sprintf(`{"items":[{"ref":"child","id":%q,"parent":%q}]}`, child.Id, root.Slug))
			page, err = reader.Get(ctx, library.NewKey(child.Id))
			require.NoError(t, err)
			parent := page.Parent.OrZero()
			assert.Equal(t, root.Id, parent.GetID().String())
		}))
	}))
}

type unavailablePageScraper struct{}

func (unavailablePageScraper) Scrape(context.Context, url.URL) (*scrape.WebContent, error) {
	return nil, errors.New("page is unavailable")
}
