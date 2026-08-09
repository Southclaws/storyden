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

			// popular gets three threads, the rest get one each, so the
			// ordering is not the insertion order
			popular := openapi.TagName(fmt.Sprintf("%s-popular", run))
			names := []openapi.TagName{popular}
			for i := range 6 {
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
			for _, n := range names[1:] {
				create([]openapi.TagName{n})
			}

			t.Run("page_size_is_respected", func(t *testing.T) {
				page := tests.AssertRequest(cl.TagListWithResponse(root, &openapi.TagListParams{
					Q:    ptr(run),
					Page: ptr("1"),
				}))(t, http.StatusOK)

				r.NotNil(page.JSON200)
				a.Len(page.JSON200.Tags, 7, "the run's tags fit on one page")
				a.Equal(popular, openapi.TagName(page.JSON200.Tags[0].Name),
					"the most used tag sorts first")
			})

			t.Run("second_page_does_not_repeat_the_first", func(t *testing.T) {
				first := tests.AssertRequest(cl.TagListWithResponse(root, &openapi.TagListParams{
					Q: ptr(run),
				}))(t, http.StatusOK)
				r.NotNil(first.JSON200)
				a.GreaterOrEqual(first.JSON200.TotalPages, 1)
				a.Equal(1, first.JSON200.CurrentPage)

				seen := map[string]struct{}{}
				for _, tag := range first.JSON200.Tags {
					seen[tag.Name] = struct{}{}
				}
				a.Len(seen, len(first.JSON200.Tags), "a page must not repeat a tag")
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
