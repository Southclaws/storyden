package thread_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
	"github.com/Southclaws/storyden/tests"
)

func TestThreadModerationWordLists(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		settingsRepo *settings.SettingsRepository,
		bus *pubsub.Bus,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)

			userCtx, user := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			sessionUser := sh.WithSession(userCtx)

			updateModerationSettings(t, root, settingsRepo, bus, settings.ModerationServiceSettings{
				WordBlockList:  opt.New([]string{"banned"}),
				WordReportList: opt.New([]string{}),
			})

			blockedThread, err := cl.ThreadCreateWithResponse(userCtx, openapi.ThreadInitialProps{
				Body:       opt.New("<p>This contains a banned topic</p>").Ptr(),
				Title:      "Banned content",
				Visibility: opt.New(openapi.Published).Ptr(),
			}, sessionUser)
			tests.Status(t, err, blockedThread, http.StatusBadRequest)
			r.NotNil(blockedThread.JSONDefault)
			if blockedThread.JSONDefault != nil {
				r.NotNil(blockedThread.JSONDefault.Message)
				r.Equal("Content violates community guidelines", *blockedThread.JSONDefault.Message)
			}

			updateModerationSettings(t, root, settingsRepo, bus, settings.ModerationServiceSettings{
				WordBlockList:  opt.New([]string{}),
				WordReportList: opt.New([]string{"flagged"}),
			})

			reviewThread, err := cl.ThreadCreateWithResponse(userCtx, openapi.ThreadInitialProps{
				Body:       opt.New("<p>This should be flagged for review</p>").Ptr(),
				Title:      "Flagged thread",
				Visibility: opt.New(openapi.Published).Ptr(),
			}, sessionUser)
			tests.Ok(t, err, reviewThread)
			r.Equal(openapi.Review, reviewThread.JSON200.Visibility, "thread with report-listed content should be sent to review")
			r.Equal(user.ID.String(), reviewThread.JSON200.Author.Id)
		}))
	}))
}

func TestReplyModerationWordLists(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		settingsRepo *settings.SettingsRepository,
		bus *pubsub.Bus,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)

			threadCtx, threadAuthor := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			replierCtx, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			sessionThreadAuthor := sh.WithSession(threadCtx)
			sessionReplier := sh.WithSession(replierCtx)

			threadCreate, err := cl.ThreadCreateWithResponse(threadCtx, openapi.ThreadInitialProps{
				Body:       opt.New("<p>Safe thread content</p>").Ptr(),
				Title:      "Reply moderation",
				Visibility: opt.New(openapi.Published).Ptr(),
			}, sessionThreadAuthor)
			tests.Ok(t, err, threadCreate)
			r.Equal(threadAuthor.ID.String(), threadCreate.JSON200.Author.Id)

			updateModerationSettings(t, root, settingsRepo, bus, settings.ModerationServiceSettings{
				WordBlockList:  opt.New([]string{"banned"}),
				WordReportList: opt.New([]string{}),
			})

			blockedReply, err := cl.ReplyCreateWithResponse(root, threadCreate.JSON200.Slug, openapi.ReplyInitialProps{
				Body: "this reply mentions a banned topic",
			}, sessionReplier)
			tests.Status(t, err, blockedReply, http.StatusBadRequest)
			r.NotNil(blockedReply.JSONDefault)
			if blockedReply.JSONDefault != nil {
				r.NotNil(blockedReply.JSONDefault.Message)
				r.Equal("Content violates community guidelines", *blockedReply.JSONDefault.Message)
			}

			updateModerationSettings(t, root, settingsRepo, bus, settings.ModerationServiceSettings{
				WordBlockList:  opt.New([]string{}),
				WordReportList: opt.New([]string{"flagged"}),
			})

			reviewReply, err := cl.ReplyCreateWithResponse(root, threadCreate.JSON200.Slug, openapi.ReplyInitialProps{
				Body: "this reply contains flagged content",
			}, sessionReplier)
			tests.Ok(t, err, reviewReply)
			r.Equal(openapi.Review, reviewReply.JSON200.Visibility, "reply with report-listed content should enter review")
		}))
	}))
}

func updateModerationSettings(
	t testing.TB,
	ctx context.Context,
	repo *settings.SettingsRepository,
	bus *pubsub.Bus,
	moderation settings.ModerationServiceSettings,
) {
	t.Helper()

	updatedSettings, err := repo.Set(ctx, settings.Settings{
		Services: opt.New(settings.ServiceSettings{
			Moderation: opt.New(moderation),
		}),
	})
	require.NoError(t, err, "should be able to update settings")

	bus.Publish(ctx, &rpc.EventSettingsUpdated{
		Settings: rpc.SerialiseSettings(*updatedSettings),
	})

	time.Sleep(100 * time.Millisecond)
}
