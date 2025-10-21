package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
	"github.com/Southclaws/storyden/sdk/go/storyden"
)

const (
	fireEmoji      = "\U0001F525"
	apiCallTimeout = 10 * time.Second
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	plugin, err := storyden.New(ctx)
	if err != nil {
		logger.Error("failed to create plugin", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if err := plugin.Shutdown(); err != nil && !errors.Is(err, context.Canceled) {
			logger.Warn("plugin shutdown returned error", slog.String("error", err.Error()))
		}
	}()

	bot := &reactBot{
		plugin: plugin,
		logger: logger,
	}

	// Declare event subscription handlers before starting the plugin runtime.
	plugin.OnThreadReplyCreated(bot.onThreadReplyCreated)

	// Run opens the RPC connection to Storyden and starts receiving events.
	if err := plugin.Run(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		logger.Error("plugin stopped", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

type reactBot struct {
	plugin *storyden.Plugin
	logger *slog.Logger
}

func (r *reactBot) onThreadReplyCreated(ctx context.Context, event *rpc.EventThreadReplyCreated) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, apiCallTimeout)
	defer cancel()

	// Build an authenticated API client using the plugin access declared in manifest.yaml.
	client, err := r.plugin.BuildAPIClient(timeoutCtx)
	if err != nil {
		return fmt.Errorf("build api client: %w", err)
	}

	// ReplyID is the post ID for the newly created reply, so we react to that post directly.
	postID := openapi.PostIDParam(event.ReplyID.String())
	resp, err := client.PostReactAddWithResponse(timeoutCtx, postID, openapi.PostReactAddJSONRequestBody{
		Emoji: fireEmoji,
	})
	if err != nil {
		return fmt.Errorf("create reaction: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("create reaction returned status %d", resp.StatusCode())
	}

	r.logger.Info(
		"reacted to new reply",
		slog.String("reply_id", event.ReplyID.String()),
		slog.String("thread_id", event.ThreadID.String()),
		slog.String("emoji", fireEmoji),
	)

	return nil
}
