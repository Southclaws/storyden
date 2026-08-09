package tags_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/google/uuid"
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

func TestTagListIsPagedAndOrderedByUse(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			a := assert.New(t)
			r := require.New(t)

			ctx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			session := sh.WithSession(ctx)

			run := uuid.NewString()[:8]

			// popular gets three threads and one node, the rest get one thread
			// each, so the ordering is not the insertion order
			popular := openapi.TagName(fmt.Sprintf("%s-popular", run))
			names := []openapi.TagName{popular}
			const quietTagCount = 55
			for i := range quietTagCount {
				names = append(names, openapi.TagName(fmt.Sprintf("%s-quiet-%d", run, i)))
			}

			create := func(tags []openapi.TagName) {
				resp, err := cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>tagged</p>").Ptr(),
					Visibility: opt.New(openapi.VisibilityPublished).Ptr(),
					Title:      "tag paging " + uuid.NewString(),
					Tags:       &tags,
				}, session)
				tests.Ok(t, err, resp)
			}

			for range 3 {
				create([]openapi.TagName{popular})
			}
			tests.AssertRequest(cl.NodeCreateWithResponse(root, openapi.NodeInitialProps{
				Name:       "tag paging " + uuid.NewString(),
				Visibility: opt.New(openapi.VisibilityPublished).Ptr(),
				Tags:       &[]openapi.TagName{popular},
			}, session))(t, http.StatusOK)
			for _, n := range names[1:] {
				create([]openapi.TagName{n})
			}

			t.Run("page_size_is_respected", func(t *testing.T) {
				page := tests.AssertRequest(cl.TagListWithResponse(root, &openapi.TagListParams{
					Q:    ptr(run),
					Page: ptr("1"),
				}))(t, http.StatusOK)

				r.NotNil(page.JSON200)
				a.Len(page.JSON200.Tags, page.JSON200.PageSize)
				a.Equal(2, page.JSON200.TotalPages)
				a.Equal(popular, openapi.TagName(page.JSON200.Tags[0].Name),
					"the most used tag sorts first")
				a.Equal(4, page.JSON200.Tags[0].ItemCount,
					"posts and nodes are each counted once")
			})

			t.Run("second_page_does_not_repeat_the_first", func(t *testing.T) {
				first := tests.AssertRequest(cl.TagListWithResponse(root, &openapi.TagListParams{
					Q:    ptr(run),
					Page: ptr("1"),
				}))(t, http.StatusOK)
				r.NotNil(first.JSON200)
				a.Equal(2, first.JSON200.TotalPages)
				a.Equal(1, first.JSON200.CurrentPage)

				second := tests.AssertRequest(cl.TagListWithResponse(root, &openapi.TagListParams{
					Q:    ptr(run),
					Page: ptr("2"),
				}))(t, http.StatusOK)
				r.NotNil(second.JSON200)
				a.Equal(2, second.JSON200.CurrentPage)

				seen := map[string]struct{}{}
				for _, tag := range first.JSON200.Tags {
					seen[tag.Name] = struct{}{}
				}
				for _, tag := range second.JSON200.Tags {
					_, repeated := seen[tag.Name]
					a.False(repeated, "a tag must not appear on two pages")
					seen[tag.Name] = struct{}{}
				}
				a.Len(seen, quietTagCount+1, "walking the pages must reach every tag")
			})

			t.Run("unfiltered_list_is_bounded", func(t *testing.T) {
				all := tests.AssertRequest(cl.TagListWithResponse(root, &openapi.TagListParams{}))(t, http.StatusOK)

				r.NotNil(all.JSON200)
				a.LessOrEqual(len(all.JSON200.Tags), all.JSON200.PageSize,
					"an unauthenticated caller must never receive the whole table")
			})
		}))
	}))
}

func ptr[T any](v T) *T { return &v }
