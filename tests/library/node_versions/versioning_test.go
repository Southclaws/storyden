package node_versions_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/oapi-codegen/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestNodeVersionProposalLifecycle(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			authorCtx, author := e2e.WithAccount(root, aw, seed.Account_007_Freyr)
			authorSession := sh.WithSession(authorCtx)

			otherCtx, _ := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			otherSession := sh.WithSession(otherCtx)

			content := "<p>Original content.</p>"
			published := openapi.VisibilityPublished
			name := "version-target-" + uuid.NewString()

			node, err := cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
				Name:       name,
				Content:    &content,
				Visibility: &published,
			}, adminSession)
			tests.Ok(t, err, node)

			t.Run("draft_proposal_can_be_applied_by_manager", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				updatedName := "Versioned name " + uuid.NewString()
				updatedSlug := "versioned-name-" + uuid.NewString()
				updatedDescription := "A proposed description"
				updatedContent := "<p>Proposed content.</p>"
				propertyType := openapi.PropertyTypeText
				propertySort := openapi.PropertySortKey("1")
				properties := openapi.PropertyMutationList{
					{Name: "Release year", Value: "1992", Type: &propertyType, Sort: &propertySort},
				}
				meta := openapi.Metadata{"source": "e2e"}

				create, err := cl.NodeVersionCreateWithResponse(root, node.JSON200.Id, openapi.NodeVersionCreateJSONRequestBody{
					Name:        &updatedName,
					Slug:        &updatedSlug,
					Description: nullable.NewNullableWithValue(openapi.NodeDescription(updatedDescription)),
					Content:     nullable.NewNullableWithValue(openapi.PostContent(updatedContent)),
					Properties:  &properties,
					Meta:        &meta,
				}, authorSession)
				tests.Ok(t, err, create)

				version := create.JSON200
				r.NotNil(version)
				a.Equal(openapi.NodeVersionStatusDraft, version.Status)
				a.Equal(author.ID.String(), string(version.Author.Id))
				a.Equal(updatedName, version.Name)
				a.Equal(updatedSlug, string(version.Slug))
				a.Equal(updatedDescription, version.Description.MustGet())
				a.Equal("<body><p>Proposed content.</p></body>", version.Content.MustGet())
				a.Equal(meta, version.Meta)

				proposedProperties := version.Properties
				r.Len(proposedProperties, 1)
				a.Equal("Release year", proposedProperties[0].Name)
				a.Equal("1992", proposedProperties[0].Value)
				r.NotNil(proposedProperties[0].Type)
				a.Equal(openapi.PropertyTypeText, *proposedProperties[0].Type)

				nodeBeforeApply, err := cl.NodeGetWithResponse(root, node.JSON200.Slug, &openapi.NodeGetParams{})
				tests.Ok(t, err, nodeBeforeApply)
				a.Equal(name, nodeBeforeApply.JSON200.Name)
				a.Equal("<body><p>Original content.</p></body>", *nodeBeforeApply.JSON200.Content)

				authorList, err := cl.NodeVersionListWithResponse(root, node.JSON200.Slug, nil, authorSession)
				tests.Ok(t, err, authorList)
				r.Len(authorList.JSON200.Versions, 1)
				a.Equal(version.Id, authorList.JSON200.Versions[0].Id)

				otherList, err := cl.NodeVersionListWithResponse(root, node.JSON200.Slug, nil, otherSession)
				tests.Ok(t, err, otherList)
				a.Empty(otherList.JSON200.Versions)

				managerList, err := cl.NodeVersionListWithResponse(root, node.JSON200.Slug, nil, adminSession)
				tests.Ok(t, err, managerList)
				r.Len(managerList.JSON200.Versions, 1)
				a.Equal(version.Id, managerList.JSON200.Versions[0].Id)

				applyAsAuthor, err := cl.NodeVersionUpdateStatusWithResponse(root, node.JSON200.Slug, version.Id, openapi.NodeVersionUpdateStatusJSONRequestBody{
					Status: openapi.NodeVersionStatusApplied,
				}, authorSession)
				tests.Status(t, err, applyAsAuthor, http.StatusForbidden)

				apply, err := cl.NodeVersionUpdateStatusWithResponse(root, node.JSON200.Id, version.Id, openapi.NodeVersionUpdateStatusJSONRequestBody{
					Status: openapi.NodeVersionStatusApplied,
				}, adminSession)
				tests.Ok(t, err, apply)
				r.NotNil(apply.JSON200)
				a.Equal(openapi.NodeVersionStatusApplied, apply.JSON200.Status)

				var applyNotification openapi.Notification
				require.Eventually(t, func() bool {
					notList, listErr := cl.NotificationListWithResponse(root, &openapi.NotificationListParams{}, authorSession)
					if listErr != nil || notList == nil || notList.JSON200 == nil {
						return false
					}

					for _, n := range notList.JSON200.Notifications {
						if n.Event != openapi.NodeVersionApplied || n.Item == nil {
							continue
						}

						item, itemErr := n.Item.AsDatagraphItemNode()
						if itemErr != nil {
							continue
						}

						if item.Ref.Id == node.JSON200.Id {
							applyNotification = n
							return true
						}
					}

					return false
				}, 5*time.Second, 100*time.Millisecond, "expected node version applied notification with node item")

				r.NotNil(applyNotification.Item)
				notificationItem, err := applyNotification.Item.AsDatagraphItemNode()
				r.NoError(err)
				a.Equal(openapi.DatagraphItemKindNode, notificationItem.Kind)
				a.Equal(node.JSON200.Id, notificationItem.Ref.Id)
				a.Equal(updatedSlug, notificationItem.Ref.Slug)

				nodeAfterApply, err := cl.NodeGetWithResponse(root, updatedSlug, &openapi.NodeGetParams{})
				tests.Ok(t, err, nodeAfterApply)
				a.Equal(updatedName, nodeAfterApply.JSON200.Name)
				a.Equal(updatedSlug, nodeAfterApply.JSON200.Slug)
				r.NotNil(nodeAfterApply.JSON200.CurrentVersionId)
				a.Equal(version.Id, string(*nodeAfterApply.JSON200.CurrentVersionId))
				a.Equal("Proposed content.", nodeAfterApply.JSON200.Description)
				a.Equal("<body><p>Proposed content.</p></body>", *nodeAfterApply.JSON200.Content)
				r.Len(nodeAfterApply.JSON200.Properties, 1)
				a.Equal("Release year", nodeAfterApply.JSON200.Properties[0].Name)
				a.Equal("1992", nodeAfterApply.JSON200.Properties[0].Value)

				updateApplied, err := cl.NodeVersionUpdateWithResponse(root, updatedSlug, version.Id, openapi.NodeVersionUpdateJSONRequestBody{
					Name: &name,
				}, authorSession)
				tests.Status(t, err, updateApplied, http.StatusForbidden)

				publicList, err := cl.NodeVersionListWithResponse(root, node.JSON200.Id, nil)
				tests.Ok(t, err, publicList)
				r.Len(publicList.JSON200.Versions, 1)
				a.Equal(openapi.NodeVersionStatusApplied, publicList.JSON200.Versions[0].Status)
			})
		}))
	}))
}

func TestNodeVersionDraftDeletion(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			authorCtx, _ := e2e.WithAccount(root, aw, seed.Account_007_Freyr)
			authorSession := sh.WithSession(authorCtx)

			t.Run("author_can_discard_own_draft", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				node := createPublishedNode(t, root, cl, adminSession, "author-discard-target")
				version := createDraftVersion(t, root, cl, authorSession, node.Slug, "Discarded proposal "+uuid.NewString())
				r.NotNil(version)

				deleteAsAuthor, err := cl.NodeVersionDeleteWithResponse(root, node.Slug, version.Id, authorSession)
				tests.Status(t, err, deleteAsAuthor, http.StatusOK)

				list, err := cl.NodeVersionListWithResponse(root, node.Slug, nil, adminSession)
				tests.Ok(t, err, list)
				a.Empty(list.JSON200.Versions)
			})

			t.Run("manager_can_discard_another_authors_draft", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				node := createPublishedNode(t, root, cl, adminSession, "manager-discard-target")
				version := createDraftVersion(t, root, cl, authorSession, node.Slug, "Manager discarded proposal "+uuid.NewString())
				r.NotNil(version)

				deleteAsManager, err := cl.NodeVersionDeleteWithResponse(root, node.Slug, version.Id, adminSession)
				tests.Status(t, err, deleteAsManager, http.StatusOK)

				list, err := cl.NodeVersionListWithResponse(root, node.Slug, nil, adminSession)
				tests.Ok(t, err, list)
				a.Empty(list.JSON200.Versions)
			})
		}))
	}))
}

func TestNodeVersionLinearCheckpointRules(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			authorCtx, _ := e2e.WithAccount(root, aw, seed.Account_007_Freyr)
			authorSession := sh.WithSession(authorCtx)

			t.Run("second_draft_is_rejected", func(t *testing.T) {
				node := createPublishedNode(t, root, cl, adminSession, "single-draft-target")
				createDraftVersion(t, root, cl, authorSession, node.Slug, "First draft "+uuid.NewString())

				secondName := "Second draft " + uuid.NewString()
				second, err := cl.NodeVersionCreateWithResponse(root, node.Slug, openapi.NodeVersionCreateJSONRequestBody{
					Name: &secondName,
				}, authorSession)
				tests.Status(t, err, second, http.StatusConflict)
			})

			t.Run("draft_alias_reads_and_updates_the_visible_draft", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				node := createPublishedNode(t, root, cl, adminSession, "draft-alias-target")

				missing, err := cl.NodeVersionDraftGetWithResponse(root, node.Slug, authorSession)
				tests.Status(t, err, missing, http.StatusNotFound)

				version := createDraftVersion(t, root, cl, authorSession, node.Slug, "Alias draft "+uuid.NewString())

				authorGet, err := cl.NodeVersionDraftGetWithResponse(root, node.Slug, authorSession)
				tests.Ok(t, err, authorGet)
				r.NotNil(authorGet.JSON200)
				a.Equal(version.Id, authorGet.JSON200.Id)

				publicGet, err := cl.NodeVersionDraftGetWithResponse(root, node.Slug)
				tests.Status(t, err, publicGet, http.StatusNotFound)

				managerGet, err := cl.NodeVersionDraftGetWithResponse(root, node.Slug, adminSession)
				tests.Ok(t, err, managerGet)
				r.NotNil(managerGet.JSON200)
				a.Equal(version.Id, managerGet.JSON200.Id)

				updatedName := "Alias updated " + uuid.NewString()
				update, err := cl.NodeVersionDraftUpdateWithResponse(root, node.Slug, openapi.NodeVersionDraftUpdateJSONRequestBody{
					Name: &updatedName,
				}, authorSession)
				tests.Ok(t, err, update)
				r.NotNil(update.JSON200)
				a.Equal(version.Id, update.JSON200.Id)
				a.Equal(updatedName, update.JSON200.Name)
			})

			t.Run("direct_versioned_edit_is_rejected_while_draft_exists", func(t *testing.T) {
				node := createPublishedNode(t, root, cl, adminSession, "draft-block-target")
				createDraftVersion(t, root, cl, authorSession, node.Slug, "Blocking draft "+uuid.NewString())

				newName := "Direct edit " + uuid.NewString()
				update, err := cl.NodeUpdateWithResponse(root, node.Slug, openapi.NodeUpdateJSONRequestBody{
					Name: &newName,
				}, adminSession)
				tests.Status(t, err, update, http.StatusConflict)
			})

			t.Run("direct_versioned_edit_clears_current_version_pointer", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				node := createPublishedNode(t, root, cl, adminSession, "pointer-clear-target")
				version := createDraftVersion(t, root, cl, authorSession, node.Slug, "Applied draft "+uuid.NewString())

				apply, err := cl.NodeVersionUpdateStatusWithResponse(root, node.Slug, version.Id, openapi.NodeVersionUpdateStatusJSONRequestBody{
					Status: openapi.NodeVersionStatusApplied,
				}, adminSession)
				tests.Ok(t, err, apply)

				appliedNode, err := cl.NodeGetWithResponse(root, apply.JSON200.Slug, &openapi.NodeGetParams{}, adminSession)
				tests.Ok(t, err, appliedNode)
				r.NotNil(appliedNode.JSON200.CurrentVersionId)
				a.Equal(version.Id, string(*appliedNode.JSON200.CurrentVersionId))

				newName := "Direct update " + uuid.NewString()
				update, err := cl.NodeUpdateWithResponse(root, apply.JSON200.Slug, openapi.NodeUpdateJSONRequestBody{
					Name: &newName,
				}, adminSession)
				tests.Ok(t, err, update)
				a.Nil(update.JSON200.CurrentVersionId)
			})

			t.Run("visibility_change_preserves_current_version_pointer", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				node := createPublishedNode(t, root, cl, adminSession, "pointer-preserve-target")
				version := createDraftVersion(t, root, cl, authorSession, node.Slug, "Visibility draft "+uuid.NewString())

				apply, err := cl.NodeVersionUpdateStatusWithResponse(root, node.Slug, version.Id, openapi.NodeVersionUpdateStatusJSONRequestBody{
					Status: openapi.NodeVersionStatusApplied,
				}, adminSession)
				tests.Ok(t, err, apply)

				unlisted := openapi.VisibilityUnlisted
				update, err := cl.NodeUpdateVisibilityWithResponse(root, apply.JSON200.Slug, openapi.VisibilityMutationProps{
					Visibility: unlisted,
				}, adminSession)
				tests.Ok(t, err, update)

				r.NotNil(update.JSON200.CurrentVersionId)
				a.Equal(version.Id, string(*update.JSON200.CurrentVersionId))
			})

			t.Run("version_get_includes_previous_applied_reference", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				node := createPublishedNode(t, root, cl, adminSession, "previous-version-target")

				first := createDraftVersion(t, root, cl, authorSession, node.Slug, "First applied "+uuid.NewString())
				firstApply, err := cl.NodeVersionUpdateStatusWithResponse(root, node.Slug, first.Id, openapi.NodeVersionUpdateStatusJSONRequestBody{
					Status: openapi.NodeVersionStatusApplied,
				}, adminSession)
				tests.Ok(t, err, firstApply)
				r.NotNil(firstApply.JSON200)

				firstGet, err := cl.NodeVersionGetWithResponse(root, firstApply.JSON200.Slug, firstApply.JSON200.Id, adminSession)
				tests.Ok(t, err, firstGet)
				r.NotNil(firstGet.JSON200)
				a.Nil(firstGet.JSON200.Previous)

				second := createDraftVersion(t, root, cl, authorSession, firstApply.JSON200.Slug, "Second applied "+uuid.NewString())
				secondApply, err := cl.NodeVersionUpdateStatusWithResponse(root, firstApply.JSON200.Slug, second.Id, openapi.NodeVersionUpdateStatusJSONRequestBody{
					Status: openapi.NodeVersionStatusApplied,
				}, adminSession)
				tests.Ok(t, err, secondApply)
				r.NotNil(secondApply.JSON200)

				secondGet, err := cl.NodeVersionGetWithResponse(root, secondApply.JSON200.Slug, secondApply.JSON200.Id, adminSession)
				tests.Ok(t, err, secondGet)
				r.NotNil(secondGet.JSON200)
				r.NotNil(secondGet.JSON200.Previous)
				a.Equal(firstApply.JSON200.Id, secondGet.JSON200.Previous.Id)
				a.Equal(openapi.NodeVersionStatusApplied, secondGet.JSON200.Previous.Status)
				a.Equal(firstApply.JSON200.Author.Id, secondGet.JSON200.Previous.Author.Id)

				list, err := cl.NodeVersionListWithResponse(root, secondApply.JSON200.Slug, nil, adminSession)
				tests.Ok(t, err, list)
				r.NotNil(list.JSON200)
				r.Len(list.JSON200.Versions, 2)
				a.Nil(list.JSON200.Versions[0].Previous)
				a.Nil(list.JSON200.Versions[1].Previous)
			})

			t.Run("version_list_is_paginated_newest_first", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				node := createPublishedNode(t, root, cl, adminSession, "paged-version-target")
				appliedIDs := make([]string, 0, 51)

				for i := 0; i < 51; i++ {
					version := createDraftVersion(
						t,
						root,
						cl,
						authorSession,
						node.Slug,
						fmt.Sprintf("Paged version %02d %s", i, uuid.NewString()),
					)

					apply, err := cl.NodeVersionUpdateStatusWithResponse(root, node.Slug, version.Id, openapi.NodeVersionUpdateStatusJSONRequestBody{
						Status: openapi.NodeVersionStatusApplied,
					}, adminSession)
					tests.Ok(t, err, apply)
					appliedIDs = append(appliedIDs, apply.JSON200.Id)
				}

				firstPage, err := cl.NodeVersionListWithResponse(root, node.Slug, nil, adminSession)
				tests.Ok(t, err, firstPage)
				r.NotNil(firstPage.JSON200)
				r.Len(firstPage.JSON200.Versions, 50)
				a.Equal(1, firstPage.JSON200.CurrentPage)
				a.Equal(50, firstPage.JSON200.PageSize)
				a.Equal(50, firstPage.JSON200.Results)
				a.Equal(2, firstPage.JSON200.TotalPages)
				r.NotNil(firstPage.JSON200.NextPage)
				a.Equal(2, *firstPage.JSON200.NextPage)
				a.Equal(appliedIDs[50], firstPage.JSON200.Versions[0].Id)
				a.Equal(appliedIDs[1], firstPage.JSON200.Versions[49].Id)

				page := openapi.PaginationQuery("2")
				secondPage, err := cl.NodeVersionListWithResponse(root, node.Slug, &openapi.NodeVersionListParams{
					Page: &page,
				}, adminSession)
				tests.Ok(t, err, secondPage)
				r.NotNil(secondPage.JSON200)
				r.Len(secondPage.JSON200.Versions, 1)
				a.Equal(2, secondPage.JSON200.CurrentPage)
				a.Equal(50, secondPage.JSON200.PageSize)
				a.Equal(1, secondPage.JSON200.Results)
				a.Equal(2, secondPage.JSON200.TotalPages)
				a.Nil(secondPage.JSON200.NextPage)
				a.Equal(appliedIDs[0], secondPage.JSON200.Versions[0].Id)
			})
		}))
	}))
}
