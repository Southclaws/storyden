package pluginbuilder

import (
	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type SDKReferenceInput struct {
	Topic string `json:"topic,omitempty" jsonschema:"SDK topic to look up, for example thread replies, events, api client, or robot run"`
}

type SDKReferenceResult struct {
	Content string `json:"content"`
}

func (a *Agent) addSDKReferenceTools(add toolAdder) error {
	return add(functiontool.New(functiontool.Config{
		Name:        "plugin_sdk_reference",
		Description: "Return curated Storyden plugin concepts, gotchas, and common examples. This is not exhaustive API documentation; use Go discovery tools for actual package symbols and methods.",
	}, func(ctx adktool.Context, args SDKReferenceInput) (SDKReferenceResult, error) {
		return SDKReference(args), nil
	}))
}

func SDKReference(SDKReferenceInput) SDKReferenceResult {
	return SDKReferenceResult{Content: sdkReferenceContent()}
}

func sdkReferenceContent() string {
	return `Storyden Go plugin SDK reference

Core setup:
- Import "github.com/Southclaws/storyden/sdk/go/storyden".
- Import "github.com/Southclaws/storyden/lib/plugin/rpc" for event payload types.
- Call pl, err := storyden.New(ctx), register handlers, then pl.Run(ctx).

Event handlers:
- pl.OnThreadPublished(func(ctx context.Context, event *rpc.EventThreadPublished) error { ... })
- EventThreadPublished has field ID post.ID. Use event.ID.String() when a ThreadMark/string is needed.
- There is no event.ThreadID field on EventThreadPublished.

Host API access:
- To call Storyden HTTP APIs from a plugin, call api, err := pl.BuildAPIClient(ctx).
- BuildAPIClient returns *openapi.ClientWithResponses from app/transports/http/openapi.
- The client is already configured with the plugin access key.
- Prefer building the API client inside the event/configuration handler that needs it. Do not fail plugin startup just because a host API client cannot be built before the plugin has connected.

Reply to a published thread:
- There is currently no pl.ThreadReply, pl.Reply, or pl.ReplyToThread convenience method.
- Use BuildAPIClient and ReplyCreateWithResponse:

	import (
		"context"
		"fmt"
		"log"
		"os"
		"os/signal"

		"github.com/Southclaws/storyden/app/transports/http/openapi"
		"github.com/Southclaws/storyden/lib/plugin/rpc"
		"github.com/Southclaws/storyden/sdk/go/storyden"
	)

	func main() {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()

		pl, err := storyden.New(ctx)
		if err != nil {
			log.Fatalf("failed to initialise plugin: %v", err)
		}
		defer pl.Shutdown()

		pl.OnThreadPublished(func(ctx context.Context, event *rpc.EventThreadPublished) error {
			api, err := pl.BuildAPIClient(ctx)
			if err != nil {
				return fmt.Errorf("failed to build API client: %v", err)
			}

			resp, err := api.ReplyCreateWithResponse(ctx, event.ID.String(), openapi.ReplyInitialProps{
				Body: openapi.PostContent("🔥"),
			})
			if err != nil {
				return err
			}
			if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
				return fmt.Errorf("reply create failed: status %d", resp.StatusCode())
			}
			return nil
		})

		if err := pl.Run(ctx); err != nil {
			log.Fatalf("plugin stopped: %v", err)
		}
	}

Robot run:
- pl.RunRobot(ctx, robotID, message) runs a Storyden Robot through plugin RPC.

Configuration:
- pl.GetConfig(ctx, keys...) reads plugin configuration.
- pl.OnConfigure(func(ctx context.Context, config map[string]any) error { ... }) handles configuration changes.
`
}
