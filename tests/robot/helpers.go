package robot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Southclaws/opt"
	durablestreams "github.com/durable-streams/durable-streams/packages/client-go"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type requestEditorTransport struct {
	base http.RoundTripper
	edit openapi.RequestEditorFn
}

func (t requestEditorTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	clone := request.Clone(request.Context())
	clone.Header = request.Header.Clone()
	if err := t.edit(request.Context(), clone); err != nil {
		return nil, err
	}
	return t.base.RoundTrip(clone)
}

type DurableJSONRead[T any] struct {
	Items       []T
	UsedLiveSSE bool
}

// ReadDurableJSON reads a complete JSON Durable Stream through the official Go
// client. The request editor is applied to both the conformance HEAD request
// and the catch-up/live requests so authenticated integration tests exercise
// the same protocol path as a real client.
func ReadDurableJSON[T any](ctx context.Context, streamURL string, edit openapi.RequestEditorFn) (*DurableJSONRead[T], error) {
	return readDurableJSON[T](ctx, streamURL, edit, "", nil)
}

func ReadDurableJSONUntil[T any](ctx context.Context, streamURL string, edit openapi.RequestEditorFn, offset string, stop func(T) bool) (*DurableJSONRead[T], error) {
	return readDurableJSON[T](ctx, streamURL, edit, offset, stop)
}

func readDurableJSON[T any](ctx context.Context, streamURL string, edit openapi.RequestEditorFn, offset string, stop func(T) bool) (*DurableJSONRead[T], error) {
	httpClient := &http.Client{Transport: requestEditorTransport{
		base: http.DefaultTransport,
		edit: edit,
	}}
	client := durablestreams.NewClient(durablestreams.WithHTTPClient(httpClient))
	stream := client.Stream(streamURL)

	metadata, err := stream.Head(ctx)
	if err != nil {
		return nil, err
	}
	if metadata.ContentType != "application/json" {
		return nil, fmt.Errorf("durable stream content type: got %q, want application/json", metadata.ContentType)
	}
	if metadata.NextOffset == "" {
		return nil, errors.New("durable stream HEAD response is missing Stream-Next-Offset")
	}

	options := []durablestreams.ReadOption{durablestreams.WithLive(durablestreams.LiveModeSSE)}
	if offset != "" {
		options = append(options, durablestreams.WithOffset(durablestreams.Offset(offset)))
	}
	iterator := stream.Read(ctx, options...)
	defer iterator.Close()

	result := &DurableJSONRead[T]{}
	for {
		chunk, err := iterator.Next()
		if errors.Is(err, durablestreams.Done) {
			return result, nil
		}
		if err != nil {
			return nil, err
		}
		if len(chunk.Data) > 0 {
			var batch []T
			if err := json.Unmarshal(chunk.Data, &batch); err != nil {
				return nil, fmt.Errorf("decode durable stream JSON batch: %w", err)
			}
			for _, item := range batch {
				result.Items = append(result.Items, item)
				if stop != nil && stop(item) {
					return result, nil
				}
			}
		}
		if chunk.Cursor != "" {
			result.UsedLiveSSE = true
		}
		if chunk.StreamClosed {
			return result, nil
		}
	}
}

func WithRobotSettings(model string) fx.Option {
	return fx.Invoke(func(lc fx.Lifecycle, root context.Context, repo *settings.SettingsRepository) {
		lc.Append(fx.StartHook(func() error {
			return SetRobotSettings(root, repo, model)
		}))
	})
}

func SetRobotSettings(ctx context.Context, repo *settings.SettingsRepository, model string) error {
	_, err := repo.Set(ctx, settings.Settings{
		Services: opt.New(settings.ServiceSettings{
			Robots: opt.New(settings.RobotServiceSettings{
				Enabled:      opt.New(true),
				DefaultModel: opt.New(model),
				Providers: opt.New(map[string]settings.RobotProviderSettings{
					"mock": {
						Enabled: opt.New(true),
					},
				}),
			}),
		}),
	})
	return err
}
