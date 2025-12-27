package thread_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestThreadModerationLength(t *testing.T) {
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

			updatedSettings, err := settingsRepo.Set(root, settings.Settings{
				Services: opt.New(settings.ServiceSettings{
					Moderation: opt.New(settings.ModerationServiceSettings{
						MaxThreadBodyLength: opt.New(1),
					}),
				}),
			})
			r.NoError(err, "should be able to update settings")

			bus.Publish(root, &message.EventSettingsUpdated{
				Settings: updatedSettings,
			})

			time.Sleep(100 * time.Millisecond)

			longContent := "<p>" + strings.Repeat("a", 100) + "</p>"
			threadCreate, err := cl.ThreadCreateWithResponse(userCtx, openapi.ThreadInitialProps{
				Body:       opt.New(longContent).Ptr(),
				Title:      "Test Thread",
				Visibility: opt.New(openapi.Published).Ptr(),
			}, sessionUser)
			r.NoError(err)
			if threadCreate.StatusCode() != 200 {
				t.Logf("Thread create response: %+v", threadCreate)
				r.Fail("Expected 200 status code")
			}

			r.Equal(openapi.Review, threadCreate.JSON200.Visibility, "thread with too-long content should be flagged for review")
			r.Equal(user.ID.String(), threadCreate.JSON200.Author.Id)
			r.Equal("Test Thread", threadCreate.JSON200.Title)

			updatedSettings2, err := settingsRepo.Set(root, settings.Settings{
				Services: opt.New(settings.ServiceSettings{
					Moderation: opt.New(settings.ModerationServiceSettings{
						MaxThreadBodyLength: opt.New(10000),
					}),
				}),
			})
			r.NoError(err)

			bus.Publish(root, &message.EventSettingsUpdated{
				Settings: updatedSettings2,
			})

			time.Sleep(100 * time.Millisecond)

			normalContent := "<p>Normal sized content</p>"
			threadCreate2, err := cl.ThreadCreateWithResponse(userCtx, openapi.ThreadInitialProps{
				Body:       opt.New(normalContent).Ptr(),
				Title:      "Test Thread 2",
				Visibility: opt.New(openapi.Published).Ptr(),
			}, sessionUser)
			tests.Ok(t, err, threadCreate2)

			r.Equal(openapi.Published, threadCreate2.JSON200.Visibility, "thread with acceptable content should be published")
		}))
	}))
}
