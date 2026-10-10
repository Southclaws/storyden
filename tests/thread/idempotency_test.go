package thread_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

// These exercise the HTTP binding and generated client together. Each test has
// an isolated database; the role change must not be run in parallel.
func TestCreateIdempotencyHTTP(t *testing.T) {
	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			memberCtx, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			adminSession := sh.WithSession(adminCtx)
			memberSession := sh.WithSession(memberCtx)
			published := openapi.VisibilityPublished

			t.Run("thread replay, mismatch, principal scope and revoked permission", func(t *testing.T) {
				key := openapi.IdempotencyKey(uuid.NewString())
				params := &openapi.ThreadCreateParams{IdempotencyKey: &key}
				body := openapi.ThreadInitialProps{Title: "idempotent thread " + uuid.NewString(), Visibility: &published}
				first := tests.AssertRequest(cl.ThreadCreateWithResponse(root, params, body, memberSession))(t, http.StatusOK)
				replay := tests.AssertRequest(cl.ThreadCreateWithResponse(root, params, body, memberSession))(t, http.StatusOK)
				r.Equal(first.JSON200.Id, replay.JSON200.Id)

				changed := body
				changed.Title += " changed"
				conflict := tests.AssertRequest(cl.ThreadCreateWithResponse(root, params, changed, memberSession))(t, http.StatusConflict)
				r.Contains(string(conflict.Body), "key_mismatch")

				other := tests.AssertRequest(cl.ThreadCreateWithResponse(root, params, body, adminSession))(t, http.StatusOK)
				r.NotEqual(first.JSON200.Id, other.JSON200.Id)

				// The same receipt cannot bypass a permission removed after creation.
				tests.AssertRequest(cl.RoleUpdateWithResponse(root, role.DefaultRoleMemberID.String(), openapi.RoleUpdateJSONRequestBody{
					Permissions: &openapi.PermissionList{},
				}, adminSession))(t, http.StatusOK)
				denied := tests.AssertRequest(cl.ThreadCreateWithResponse(root, params, body, memberSession))(t, http.StatusForbidden)
				r.NotEmpty(denied.Body)
				tests.AssertRequest(cl.RoleDeleteWithResponse(root, role.DefaultRoleMemberID.String(), adminSession))(t, http.StatusOK)
			})

			t.Run("reply replay loses access to hidden parent", func(t *testing.T) {
				thread := tests.AssertRequest(cl.ThreadCreateWithResponse(root, nil, openapi.ThreadInitialProps{
					Title: "reply parent " + uuid.NewString(), Visibility: &published,
				}, adminSession))(t, http.StatusOK)
				key := openapi.IdempotencyKey(uuid.NewString())
				params := &openapi.ReplyCreateParams{IdempotencyKey: &key}
				body := openapi.ReplyInitialProps{Body: "idempotent reply"}
				first := tests.AssertRequest(cl.ReplyCreateWithResponse(root, thread.JSON200.Slug, params, body, memberSession))(t, http.StatusOK)
				replay := tests.AssertRequest(cl.ReplyCreateWithResponse(root, thread.JSON200.Slug, params, body, memberSession))(t, http.StatusOK)
				r.Equal(first.JSON200.Id, replay.JSON200.Id)

				changed := body
				changed.Body += " changed"
				tests.AssertRequest(cl.ReplyCreateWithResponse(root, thread.JSON200.Slug, params, changed, memberSession))(t, http.StatusConflict)

				// Making the parent private makes a stored reply inaccessible on replay.
				draft := openapi.VisibilityDraft
				tests.AssertRequest(cl.ThreadUpdateWithResponse(root, thread.JSON200.Slug, openapi.ThreadMutableProps{
					Visibility: &draft,
				}, adminSession))(t, http.StatusOK)
				result, err := cl.ReplyCreateWithResponse(root, thread.JSON200.Slug, params, body, memberSession)
				r.NoError(err)
				r.NotEqual(http.StatusOK, result.StatusCode())
			})

			t.Run("library page replay and resource visibility", func(t *testing.T) {
				key := openapi.IdempotencyKey(uuid.NewString())
				params := &openapi.NodeCreateParams{IdempotencyKey: &key}
				body := openapi.NodeInitialProps{Name: "idempotent page " + uuid.NewString(), Visibility: &published}
				first := tests.AssertRequest(cl.NodeCreateWithResponse(root, params, body, adminSession))(t, http.StatusOK)
				replay := tests.AssertRequest(cl.NodeCreateWithResponse(root, params, body, adminSession))(t, http.StatusOK)
				r.Equal(first.JSON200.Id, replay.JSON200.Id)
				changed := body
				changed.Name += " changed"
				tests.AssertRequest(cl.NodeCreateWithResponse(root, params, changed, adminSession))(t, http.StatusConflict)
				tests.AssertRequest(cl.NodeDeleteWithResponse(root, first.JSON200.Slug, nil, adminSession))(t, http.StatusOK)
				result, err := cl.NodeCreateWithResponse(root, params, body, adminSession)
				r.NoError(err)
				r.NotEqual(http.StatusOK, result.StatusCode())
			})
		}))
	}))
}
