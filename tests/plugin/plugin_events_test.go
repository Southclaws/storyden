package plugin_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	resource_plugin "github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/app/transports/rpc"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/tests"
)

func TestPluginEventSubscription(t *testing.T) {
	t.Parallel()

	outputDir := filepath.Join(os.TempDir(), "plugin_event_test_"+xid.New().String())
	os.Setenv("OUTPUT_DIR", outputDir)
	defer os.RemoveAll(outputDir)

	integration.Test(t, &config.Config{
		QueueType: "amqp",
		AmqpURL:   "amqp://guest:guest@localhost:5672/",
	}, e2e.Setup(), rpc.Build(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		accountWrite *account_writer.Writer,
		runner plugin_runner.Runner,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			type serverURLSetter interface {
				SetServerURL(string)
			}
			if setter, ok := runner.(serverURLSetter); ok {
				setter.SetServerURL(ts.URL)
			}

			adminHandle := "admin-" + xid.New().String()
			admin, err := cl.AuthPasswordSignupWithResponse(root, nil, openapi.AuthPair{
				Identifier: adminHandle,
				Token:      "password",
			})
			r.NoError(err)
			r.Equal(http.StatusOK, admin.StatusCode())
			adminID := account.AccountID(utils.Must(xid.FromString(admin.JSON200.Id)))
			adminSession := sh.WithSession(e2e.WithAccountID(root, adminID))

			accountWrite.Update(root, adminID, account_writer.SetAdmin(true))

			pluginPath := "test_data/event_listener/event_listener.sdx"
			pluginFile, err := os.Open(pluginPath)
			r.NoError(err)
			defer pluginFile.Close()

			addResp, err := cl.PluginAddWithBodyWithResponse(root, "application/zip", pluginFile, adminSession)
			r.NoError(err)
			r.Equal(http.StatusOK, addResp.StatusCode())
			r.NotNil(addResp.JSON200)

			pluginID := string(addResp.JSON200.Id)
			pluginIDxid, _ := xid.FromString(pluginID)
			installationID := resource_plugin.InstallationID(pluginIDxid)

			activeResp, err := cl.PluginSetActiveStateWithResponse(root, pluginID, openapi.PluginSetActiveStateJSONRequestBody{
				Active: "active",
			}, adminSession)
			r.NoError(err)
			r.Equal(http.StatusOK, activeResp.StatusCode())

			time.Sleep(2 * time.Second)

			sess, err := runner.GetSession(root, installationID)
			r.NoError(err)
			r.NotNil(sess)

			state := sess.GetReportedState()
			r.Equal(resource_plugin.ReportedStateActive, state, "plugin should be active")

			threadTitle := "Event Test Thread " + xid.New().String()
			threadCreate := tests.AssertRequest(
				cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Body:       opt.New("<p>test event publishing</p>").Ptr(),
					Visibility: opt.New(openapi.Published).Ptr(),
					Title:      threadTitle,
				}, adminSession),
			)(t, http.StatusOK)

			threadID := threadCreate.JSON200.Id
			expectedFile := filepath.Join(outputDir, fmt.Sprintf("%s.json", threadID))

			found := false
			for i := 0; i < 30; i++ {
				if _, err := os.Stat(expectedFile); err == nil {
					found = true
					break
				}
				time.Sleep(100 * time.Millisecond)
			}

			a.True(found, "event file should be created by plugin")

			if found {
				fileData, err := os.ReadFile(expectedFile)
				r.NoError(err, "should read event file")

				var eventData map[string]interface{}
				err = json.Unmarshal(fileData, &eventData)
				r.NoError(err, "should unmarshal event data")

				idField, ok := eventData["ID"]
				r.True(ok, "event should have ID field")

				var extractedID string
				switch v := idField.(type) {
				case string:
					extractedID = v
				case map[string]interface{}:
					if idStr, ok := v["id"].(string); ok {
						extractedID = idStr
					}
				}

				a.Equal(threadID, extractedID, "event should contain correct thread ID")
			}

			inactiveResp, err := cl.PluginSetActiveStateWithResponse(root, pluginID, openapi.PluginSetActiveStateJSONRequestBody{
				Active: "inactive",
			}, adminSession)
			r.NoError(err)
			r.Equal(http.StatusOK, inactiveResp.StatusCode())

			deleteResp, err := cl.PluginDeleteWithResponse(root, pluginID, adminSession)
			r.NoError(err)
			r.Equal(http.StatusNoContent, deleteResp.StatusCode())
		}))
	}))
}
