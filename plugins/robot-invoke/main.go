package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Southclaws/storyden/lib/plugin/rpc"
	"github.com/Southclaws/storyden/sdk/go/storyden"
)

const (
	robotIDString = "d8nq6veot5p1b1g4lbhg"
	prompt        = "Search the available Storyden content for anything about robots or agents. Summarise the most relevant result you find."

	connectTimeout = 15 * time.Second
	requestTimeout = 2 * time.Minute
	retryDelay     = 250 * time.Millisecond
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	plugin, err := storyden.New(ctx)
	if err != nil {
		exitError(logger, "create plugin", err)
	}

	runCtx, cancelRun := context.WithCancel(ctx)
	defer cancelRun()

	runErr := make(chan error, 1)
	go func() {
		runErr <- plugin.Run(runCtx)
	}()

	defer func() {
		cancelRun()
		if err := plugin.Shutdown(); err != nil && !errors.Is(err, context.Canceled) {
			logger.Warn("plugin shutdown returned error", slog.String("error", err.Error()))
		}

		select {
		case err := <-runErr:
			if err != nil && !errors.Is(err, context.Canceled) && !strings.Contains(err.Error(), "plugin shutting down") {
				logger.Warn("plugin run loop returned error", slog.String("error", err.Error()))
			}
		case <-time.After(time.Second):
		}
	}()

	req := rpc.RPCRequestRobotRun{
		Jsonrpc: "2.0",
		Method:  "robot_run",
		Params: rpc.RPCRequestRobotRunParams{
			RobotID: robotIDString,
			Message: prompt,
		},
	}

	resp, err := sendWhenConnected(ctx, plugin, req)
	if err != nil {
		exitError(logger, "robot_run", err)
	}

	typed, ok := resp.(*rpc.RPCResponseRobotRun)
	if !ok {
		exitError(logger, "robot_run", fmt.Errorf("unexpected response type %T", resp))
	}

	if err := printJSON(os.Stdout, typed); err != nil {
		exitError(logger, "print response", err)
	}

	if methodErr, ok := typed.Error.Get(); ok && methodErr != "" {
		os.Exit(1)
	}
}

func sendWhenConnected(ctx context.Context, plugin *storyden.Plugin, req rpc.RPCRequestRobotRun) (rpc.PluginToHostResponseUnionUnion, error) {
	connectDeadline := time.Now().Add(connectTimeout)

	for {
		requestCtx, cancelRequest := context.WithTimeout(ctx, requestTimeout)
		resp, err := plugin.Send(requestCtx, req)
		cancelRequest()
		if err == nil {
			return resp, nil
		}
		if !isConnectionPending(err) {
			return nil, err
		}

		if time.Now().After(connectDeadline) {
			return nil, fmt.Errorf("connect to Storyden RPC: %w", err)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(retryDelay):
		}
	}
}

func isConnectionPending(err error) bool {
	if err == nil {
		return false
	}
	message := err.Error()
	return strings.Contains(message, "connection closed") ||
		strings.Contains(message, "failed to send request")
}

func printJSON(out *os.File, value any) error {
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, string(encoded))
	return err
}

func exitError(logger *slog.Logger, action string, err error) {
	logger.Error(action+" failed", slog.String("error", err.Error()))
	os.Exit(1)
}
